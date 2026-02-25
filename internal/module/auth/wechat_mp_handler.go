package auth

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"net/http"
	"strings"
	"time"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/pkg/logger"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type weChatMPEvent struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   int64    `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Event        string   `xml:"Event"`
	EventKey     string   `xml:"EventKey"`
	Ticket       string   `xml:"Ticket"`
}

type weChatMPEncryptWrap struct {
	XMLName xml.Name `xml:"xml"`
	Encrypt string   `xml:"Encrypt"`
}

func (h *Handler) WeChatMPCreateScene(c *gin.Context) {
	if !h.svc.FeatureWeChatMPEnabled() {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
		return
	}
	scene, qrURL, expireSeconds, err := h.svc.CreateWeChatMPLoginScene(context.Background())
	if err != nil {
		logger.WithCtx(c).Error("auth.wechatmp create scene failed", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"scene": scene, "qrUrl": qrURL, "expiresIn": expireSeconds})
}

func (h *Handler) WeChatMPSceneStatus(c *gin.Context) {
	if !h.svc.FeatureWeChatMPEnabled() {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
		return
	}
	scene := c.Param("scene")
	val, err := h.svc.GetDelWeChatMPScenePayload(context.Background(), scene)
	if err != nil || val == "" {
		c.JSON(http.StatusNotFound, gin.H{"status": "expired"})
		return
	}
	var payload weChatMPScenePayload
	_ = json.Unmarshal([]byte(val), &payload)
	if payload.Status != "ok" || payload.OTT == "" {
		ttl := time.Until(time.Unix(payload.ExpireAt, 0))
		if ttl > 0 {
			_ = cache.RDB.Set(context.Background(), common.RedisKeyWeChatMPScene.F(scene), val, ttl).Err()
		}
		c.JSON(http.StatusOK, gin.H{"status": "pending"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "ok", "ott": payload.OTT})
}

func (h *Handler) WeChatMPCallbackGet(c *gin.Context) {
	if !h.svc.FeatureWeChatMPEnabled() {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
		return
	}
	msgSignature := c.Query("msg_signature")
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")
	echostr := c.Query("echostr")
	if echostr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}
	if msgSignature != "" {
		if !h.svc.VerifyWeChatMPMsgSignature(msgSignature, timestamp, nonce, echostr) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}
		plain, err := h.svc.decryptWeChatMPText(echostr)
		if err != nil || plain == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
			return
		}
		c.String(http.StatusOK, plain)
		return
	}
	if !h.svc.VerifyWeChatMPSignature(signature, timestamp, nonce) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signature"})
		return
	}
	c.String(http.StatusOK, echostr)
}

func (h *Handler) WeChatMPCallbackPost(c *gin.Context) {
	if !h.svc.FeatureWeChatMPEnabled() {
		c.JSON(http.StatusNotImplemented, gin.H{"error": "feature disabled"})
		return
	}
	msgSignature := c.Query("msg_signature")
	signature := c.Query("signature")
	timestamp := c.Query("timestamp")
	nonce := c.Query("nonce")

	raw, err := c.GetRawData()
	if err != nil || len(raw) == 0 {
		c.String(http.StatusOK, "success")
		return
	}
	if msgSignature != "" {
		var wrap weChatMPEncryptWrap
		if e := xml.Unmarshal(raw, &wrap); e != nil || strings.TrimSpace(wrap.Encrypt) == "" {
			c.String(http.StatusOK, "success")
			return
		}
		enc := strings.TrimSpace(wrap.Encrypt)
		if !h.svc.VerifyWeChatMPMsgSignature(msgSignature, timestamp, nonce, enc) {
			c.String(http.StatusOK, "success")
			return
		}
		plain, e := h.svc.decryptWeChatMPText(enc)
		if e != nil || plain == "" {
			c.String(http.StatusOK, "success")
			return
		}
		raw = []byte(plain)
	} else {
		if !h.svc.VerifyWeChatMPSignature(signature, timestamp, nonce) {
			c.String(http.StatusOK, "success")
			return
		}
	}

	var ev weChatMPEvent
	if e := xml.Unmarshal(raw, &ev); e != nil {
		c.String(http.StatusOK, "success")
		return
	}

	openid := strings.TrimSpace(ev.FromUserName)
	scene := ""
	switch strings.ToUpper(strings.TrimSpace(ev.Event)) {
	case "SUBSCRIBE":
		if strings.HasPrefix(ev.EventKey, "qrscene_") {
			scene = strings.TrimPrefix(ev.EventKey, "qrscene_")
		}
	case "SCAN":
		scene = strings.TrimSpace(ev.EventKey)
	}
	if openid == "" || scene == "" {
		c.String(http.StatusOK, "success")
		return
	}

	p, e := h.svc.GetWeChatMPScenePayload(context.Background(), scene)
	if e != nil || p.Scene == "" || p.Status != "pending" {
		c.String(http.StatusOK, "success")
		return
	}

	ui, _ := h.svc.FetchWeChatMPUserInfo(context.Background(), openid)
	user, err := h.svc.FindOrCreateWeChatMPUser(ui, openid)
	if err != nil {
		logger.WithCtx(c).Error("auth.wechatmp account error", zap.Error(err))
		c.String(http.StatusOK, "success")
		return
	}
	_, _ = h.svc.MarkWeChatMPSceneReady(context.Background(), scene, openid, user)
	c.String(http.StatusOK, "success")
}
