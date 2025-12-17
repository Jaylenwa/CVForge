package handlers

import (
	"testing"

	"openresume/config"
	"openresume/db"
	"openresume/models"
)

func TestMakeWeChatLoginURL(t *testing.T) {
	u := makeWeChatLoginURL(config.Config{WeChatAppID: "appid", WeChatRedirectURI: "https://example.com/api/v1/auth/wechat/callback"}, "state123")
	if !contains(u, "open.weixin.qq.com/connect/qrconnect") ||
		!contains(u, "appid=appid") ||
		!contains(u, "redirect_uri=") ||
		!contains(u, "state=state123") {
		t.Fatalf("unexpected url: %s", u)
	}
}

func TestFindOrCreateWeChatUser(t *testing.T) {
	_, err := db.InitMySQL(config.Config{SQLitePath: ":memory:", MySQLDSN: ""})
	if err != nil {
		t.Fatalf("init db: %v", err)
	}
	g := db.Gorm()
	ui := wechatUserInfo{OpenID: "o1", UnionID: "u1", Nickname: "Nick", HeadImgURL: "http://img"}
	tr := wechatTokenResponse{OpenID: "o1", AccessToken: "at"}
	u, err := findOrCreateWeChatUser(g, ui, tr)
	if err != nil || u.ID == 0 {
		t.Fatalf("create user err=%v id=%d", err, u.ID)
	}
	var oa models.OAuthAccount
	if e := g.Where("provider = ? AND provider_union_id = ?", "wechat", "u1").First(&oa).Error; e != nil {
		t.Fatalf("find oauth account: %v", e)
	}
	if u.ID != oa.UserID {
		t.Fatalf("user mismatch: %d vs %d", u.ID, oa.UserID)
	}
	// second call with same unionid returns same user
	u2, err := findOrCreateWeChatUser(g, ui, tr)
	if err != nil || u2.ID != u.ID {
		t.Fatalf("second call returns different user: %d vs %d err=%v", u2.ID, u.ID, err)
	}
	// call with no unionid but same openid also returns same user
	ui2 := wechatUserInfo{OpenID: "o1", UnionID: "", Nickname: "Nick2"}
	u3, err := findOrCreateWeChatUser(g, ui2, tr)
	if err != nil || u3.ID != u.ID {
		t.Fatalf("openid match returns different user: %d vs %d err=%v", u3.ID, u.ID, err)
	}
}

func contains(s, sub string) bool { return len(s) >= len(sub) && (s == sub || (len(s) > len(sub) && (stringContains(s, sub)))) }
func stringContains(s, sub string) bool { return indexOf(s, sub) >= 0 }
func indexOf(s, sub string) int {
	n := len(sub)
	if n == 0 {
		return 0
	}
	for i := 0; i+n <= len(s); i++ {
		if s[i:i+n] == sub {
			return i
		}
	}
	return -1
}
