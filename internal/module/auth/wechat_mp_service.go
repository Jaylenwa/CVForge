package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"time"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/infra/database"

	"github.com/google/uuid"
)

type weChatMPAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	ErrCode     int    `json:"errcode"`
	ErrMsg      string `json:"errmsg"`
}

type weChatMPQRCodeCreateResponse struct {
	Ticket        string `json:"ticket"`
	ExpireSeconds int    `json:"expire_seconds"`
	URL           string `json:"url"`
	ErrCode       int    `json:"errcode"`
	ErrMsg        string `json:"errmsg"`
}

type weChatMPUserInfoResponse struct {
	Subscribe  int    `json:"subscribe"`
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid"`
	ErrCode    int    `json:"errcode"`
	ErrMsg     string `json:"errmsg"`
}

type weChatMPScenePayload struct {
	Status   string `json:"status"`
	Scene    string `json:"scene"`
	OpenID   string `json:"openid,omitempty"`
	UnionID  string `json:"unionid,omitempty"`
	OTT      string `json:"ott,omitempty"`
	UID      uint   `json:"uid,omitempty"`
	Created  int64  `json:"created"`
	ReadyAt  int64  `json:"readyAt,omitempty"`
	ExpireAt int64  `json:"expireAt"`
}

func (s *Service) FeatureWeChatMPEnabled() bool {
	return s.sysConfig.GetBool(string(common.ConfigKeyEnabledWechatMPLogin), true)
}

func (s *Service) weChatMPConfigured() bool {
	if !s.FeatureWeChatMPEnabled() {
		return false
	}
	return s.sysConfig.Get(string(common.ConfigKeyWeChatMPAppID)) != "" &&
		s.sysConfig.Get(string(common.ConfigKeyWeChatMPAppSecret)) != "" &&
		s.sysConfig.Get(string(common.ConfigKeyWeChatMPToken)) != ""
}

func (s *Service) VerifyWeChatMPSignature(signature, timestamp, nonce string) bool {
	token := s.sysConfig.Get(string(common.ConfigKeyWeChatMPToken))
	if token == "" || signature == "" || timestamp == "" || nonce == "" {
		return false
	}
	parts := []string{token, timestamp, nonce}
	sort.Strings(parts)
	h := sha1.Sum([]byte(strings.Join(parts, "")))
	return hex.EncodeToString(h[:]) == signature
}

func (s *Service) VerifyWeChatMPMsgSignature(msgSignature, timestamp, nonce, encrypt string) bool {
	token := s.sysConfig.Get(string(common.ConfigKeyWeChatMPToken))
	if token == "" || msgSignature == "" || timestamp == "" || nonce == "" || encrypt == "" {
		return false
	}
	parts := []string{token, timestamp, nonce, encrypt}
	sort.Strings(parts)
	h := sha1.Sum([]byte(strings.Join(parts, "")))
	return hex.EncodeToString(h[:]) == msgSignature
}

func (s *Service) decryptWeChatMPText(encrypt string) (string, error) {
	aesKey, ok := s.weChatMPAESKeyBytes()
	if !ok {
		return "", http.ErrNotSupported
	}
	wantAppID := s.sysConfig.Get(string(common.ConfigKeyWeChatMPAppID))
	return decryptWeChatMPTextWithKey(aesKey, wantAppID, encrypt)
}

func (s *Service) weChatMPAESKeyBytes() ([]byte, bool) {
	k := strings.TrimSpace(s.sysConfig.Get(string(common.ConfigKeyWeChatMPAESKey)))
	if k == "" {
		return nil, false
	}
	if m := len(k) % 4; m != 0 {
		k += strings.Repeat("=", 4-m)
	}
	b, err := base64.StdEncoding.DecodeString(k)
	if err != nil || len(b) != 32 {
		return nil, false
	}
	return b, true
}

func decryptWeChatMPTextWithKey(aesKey []byte, appID string, encrypt string) (string, error) {
	if len(aesKey) != 32 || encrypt == "" {
		return "", http.ErrNotSupported
	}
	rawCipher, err := base64.StdEncoding.DecodeString(encrypt)
	if err != nil {
		return "", err
	}
	if len(rawCipher)%aes.BlockSize != 0 || len(rawCipher) == 0 {
		return "", http.ErrNotSupported
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	plain := make([]byte, len(rawCipher))
	iv := aesKey[:aes.BlockSize]
	cipher.NewCBCDecrypter(block, iv).CryptBlocks(plain, rawCipher)
	plain, err = pkcs7Unpad(plain, aes.BlockSize)
	if err != nil {
		return "", err
	}
	if len(plain) < 20 {
		return "", http.ErrNotSupported
	}
	msgLen := int(binary.BigEndian.Uint32(plain[16:20]))
	if msgLen <= 0 || 20+msgLen > len(plain) {
		return "", http.ErrNotSupported
	}
	msg := plain[20 : 20+msgLen]
	appid := string(plain[20+msgLen:])
	if appID != "" && appid != appID {
		return "", http.ErrNotSupported
	}
	return string(msg), nil
}

func pkcs7Unpad(b []byte, blockSize int) ([]byte, error) {
	if len(b) == 0 || len(b)%blockSize != 0 {
		return nil, http.ErrNotSupported
	}
	pad := int(b[len(b)-1])
	if pad <= 0 || pad > blockSize || pad > len(b) {
		return nil, http.ErrNotSupported
	}
	for i := 0; i < pad; i++ {
		if b[len(b)-1-i] != byte(pad) {
			return nil, http.ErrNotSupported
		}
	}
	return b[:len(b)-pad], nil
}

func (s *Service) getWeChatMPAccessToken(ctx context.Context) (string, error) {
	if cache.RDB != nil {
		if val, err := cache.RDB.Get(ctx, common.RedisKeyWeChatMPAccessToken.F()).Result(); err == nil && val != "" {
			return val, nil
		}
	}

	appID := s.sysConfig.Get(string(common.ConfigKeyWeChatMPAppID))
	appSecret := s.sysConfig.Get(string(common.ConfigKeyWeChatMPAppSecret))
	if appID == "" || appSecret == "" {
		return "", http.ErrNotSupported
	}

	u := "https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=" + url.QueryEscape(appID) + "&secret=" + url.QueryEscape(appSecret)
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var out weChatMPAccessTokenResponse
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out.ErrCode != 0 || out.AccessToken == "" || out.ExpiresIn <= 0 {
		return "", http.ErrHandlerTimeout
	}

	if cache.RDB != nil {
		ttl := time.Duration(out.ExpiresIn-60) * time.Second
		if ttl < time.Minute {
			ttl = time.Minute
		}
		_ = cache.RDB.Set(ctx, common.RedisKeyWeChatMPAccessToken.F(), out.AccessToken, ttl).Err()
	}
	return out.AccessToken, nil
}

func (s *Service) CreateWeChatMPLoginScene(ctx context.Context) (scene string, qrURL string, expireSeconds int, err error) {
	if !s.weChatMPConfigured() {
		return "", "", 0, http.ErrNotSupported
	}
	scene = uuidNoDash()
	expireSeconds = 600

	expireAt := time.Now().Add(time.Duration(expireSeconds) * time.Second).Unix()
	payload := weChatMPScenePayload{
		Status:   "pending",
		Scene:    scene,
		Created:  time.Now().Unix(),
		ExpireAt: expireAt,
	}
	b, _ := json.Marshal(payload)
	if cache.RDB != nil {
		if e := cache.RDB.Set(ctx, common.RedisKeyWeChatMPScene.F(scene), string(b), time.Duration(expireSeconds)*time.Second).Err(); e != nil {
			return "", "", 0, e
		}
	} else {
		return "", "", 0, errors.New("redis not configured")
	}

	accessToken, err := s.getWeChatMPAccessToken(ctx)
	if err != nil {
		return "", "", 0, err
	}

	body := map[string]any{
		"expire_seconds": expireSeconds,
		"action_name":    "QR_STR_SCENE",
		"action_info": map[string]any{
			"scene": map[string]any{
				"scene_str": scene,
			},
		},
	}
	raw, _ := json.Marshal(body)
	u := "https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=" + url.QueryEscape(accessToken)
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost, u, strings.NewReader(string(raw)))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", "", 0, err
	}
	defer resp.Body.Close()

	var out weChatMPQRCodeCreateResponse
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out.ErrCode != 0 || out.Ticket == "" {
		return "", "", 0, http.ErrHandlerTimeout
	}
	qrURL = "https://mp.weixin.qq.com/cgi-bin/showqrcode?ticket=" + url.QueryEscape(out.Ticket)
	if out.ExpireSeconds > 0 {
		expireSeconds = out.ExpireSeconds
	}
	return scene, qrURL, expireSeconds, nil
}

func (s *Service) GetWeChatMPScenePayload(ctx context.Context, scene string) (weChatMPScenePayload, error) {
	var out weChatMPScenePayload
	if cache.RDB == nil || scene == "" {
		return out, http.ErrNotSupported
	}
	val, err := cache.RDB.Get(ctx, common.RedisKeyWeChatMPScene.F(scene)).Result()
	if err != nil || val == "" {
		return out, err
	}
	_ = json.Unmarshal([]byte(val), &out)
	return out, nil
}

func (s *Service) GetDelWeChatMPScenePayload(ctx context.Context, scene string) (string, error) {
	if cache.RDB == nil || scene == "" {
		return "", http.ErrNotSupported
	}
	val, err := cache.RDB.GetDel(ctx, common.RedisKeyWeChatMPScene.F(scene)).Result()
	if err == nil {
		return val, nil
	}
	val, err2 := cache.RDB.Get(ctx, common.RedisKeyWeChatMPScene.F(scene)).Result()
	if err2 != nil {
		return "", err2
	}
	_ = cache.RDB.Del(ctx, common.RedisKeyWeChatMPScene.F(scene)).Err()
	return val, nil
}

func (s *Service) FetchWeChatMPUserInfo(ctx context.Context, openid string) (weChatMPUserInfoResponse, error) {
	var out weChatMPUserInfoResponse
	if openid == "" {
		return out, http.ErrNotSupported
	}
	accessToken, err := s.getWeChatMPAccessToken(ctx)
	if err != nil {
		return out, err
	}
	u := "https://api.weixin.qq.com/cgi-bin/user/info?access_token=" + url.QueryEscape(accessToken) + "&openid=" + url.QueryEscape(openid) + "&lang=zh_CN"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	client := &http.Client{Timeout: 8 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out.ErrCode != 0 {
		return out, http.ErrHandlerTimeout
	}
	return out, nil
}

func (s *Service) FindOrCreateWeChatMPUser(ui weChatMPUserInfoResponse, openid string) (User, error) {
	var user User
	var oa OAuthAccount

	if openid == "" {
		openid = ui.OpenID
	}
	if openid == "" {
		return User{}, http.ErrNotSupported
	}

	if ui.UnionID != "" {
		if err := database.DB.Where("provider_union_id <> '' AND provider_union_id = ? AND provider = ?", ui.UnionID, common.WechatMP).First(&oa).Error; err == nil {
			if err2 := database.DB.First(&user, oa.UserID).Error; err2 == nil {
				return user, nil
			}
		}
	}
	if err := database.DB.Where("provider = ? AND provider_open_id = ?", common.WechatMP, openid).First(&oa).Error; err == nil {
		if err2 := database.DB.First(&user, oa.UserID).Error; err2 == nil {
			return user, nil
		}
	}

	name := ui.Nickname
	if name == "" {
		name = "WeChat User"
	}
	h := sha1.Sum([]byte(openid))
	email := "wechatmp_" + hex.EncodeToString(h[:]) + "@oauth.invalid"
	user = User{
		Email:     email,
		Name:      name,
		AvatarURL: ui.HeadImgURL,
		IsActive:  true,
		Role:      s.initialRole(),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return User{}, err
	}

	record := OAuthAccount{
		UserID:          user.ID,
		Provider:        common.WechatMP,
		ProviderOpenID:  openid,
		ProviderUnionID: ui.UnionID,
	}
	if b, e := json.Marshal(ui); e == nil {
		record.RawProfileJSON = string(b)
	}
	_ = database.DB.Create(&record).Error
	return user, nil
}

func (s *Service) MarkWeChatMPSceneReady(ctx context.Context, scene, openid string, user User) (string, error) {
	if cache.RDB == nil || scene == "" || openid == "" {
		return "", http.ErrNotSupported
	}

	access, refresh := s.IssueTokens(user.ID)
	ott := uuidNoDash()
	payload := map[string]any{
		"accessToken":  access,
		"refreshToken": refresh,
		"user":         s.SanitizeUser(user),
	}
	_ = s.SaveOTT(ott, payload)

	cur, _ := s.GetWeChatMPScenePayload(ctx, scene)
	expireAt := cur.ExpireAt
	if expireAt == 0 {
		expireAt = time.Now().Add(time.Minute).Unix()
	}
	ttl := time.Until(time.Unix(expireAt, 0))
	if ttl <= 0 {
		ttl = time.Minute
	}
	cur.Status = "ok"
	cur.Scene = scene
	cur.OpenID = openid
	cur.UnionID = ""
	cur.OTT = ott
	cur.UID = user.ID
	cur.ReadyAt = time.Now().Unix()
	b, _ := json.Marshal(cur)
	if err := cache.RDB.Set(ctx, common.RedisKeyWeChatMPScene.F(scene), string(b), ttl).Err(); err != nil {
		return "", err
	}
	return ott, nil
}

func uuidNoDash() string {
	u := uuid.NewString()
	u = strings.ReplaceAll(u, "-", "")
	return u
}

func readAllBody(r io.Reader) []byte {
	b, _ := io.ReadAll(io.LimitReader(r, 1024*1024))
	return b
}
