package auth

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"openresume/internal/common"
	"openresume/internal/infra/cache"
	"openresume/internal/infra/config"
	"openresume/internal/infra/database"
	conf "openresume/internal/module/config"
	"openresume/internal/pkg/mailer"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	cfg       config.Config
	sysConfig *conf.Service
}

func NewService(cfg config.Config, sysConfig *conf.Service) *Service {
	return &Service{cfg: cfg, sysConfig: sysConfig}
}

func (s *Service) FeatureWeChatEnabled() bool {
	return s.sysConfig.GetBool(string(common.ConfigKeyEnabledWechatLogin), true)
}

func (s *Service) FeatureGithubEnabled() bool {
	return s.sysConfig.GetBool(string(common.ConfigKeyEnabledGithubLogin), true)
}

func (s *Service) FrontendBase() string {
	return s.sysConfig.Get(string(common.ConfigKeyFrontendBaseURL))
}

func (s *Service) IssueTokens(uid uint) (string, string) {
	mk := func(exp time.Duration) string {
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"uid": uid, "exp": time.Now().Add(exp).Unix(), "jti": uuid.NewString()})
		signed, _ := t.SignedString([]byte(s.cfg.JWTSecret))
		return signed
	}
	return mk(2 * time.Hour), mk(7 * 24 * time.Hour)
}

func (s *Service) initialRole() common.Role {
	var n int64
	database.DB.Model(&User{}).Count(&n)
	if n == 0 {
		return common.RoleAdmin
	}
	return common.RoleUser
}

func (s *Service) GenerateVerifyCode() string {
	n := 6
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		v, _ := rand.Int(rand.Reader, big.NewInt(10))
		out[i] = byte('0' + v.Int64())
	}
	return string(out)
}

func (s *Service) SaveVerifyCode(email, code string) error {
	return cache.RDB.Set(context.Background(), common.RedisKeyVerify.F(email), code, 10*time.Minute).Err()
}

func (s *Service) ValidateVerifyCode(email, code string) bool {
	val, err := cache.RDB.Get(context.Background(), common.RedisKeyVerify.F(email)).Result()
	return err == nil && val == code
}

func (s *Service) SendCode(email string, code string) error {
	smtpCfg := mailer.SMTPSettings{
		Host:     s.sysConfig.Get(string(common.ConfigKeySMTPHost)),
		Port:     s.sysConfig.Get(string(common.ConfigKeySMTPPort)),
		Username: s.sysConfig.Get(string(common.ConfigKeySMTPUsername)),
		Password: s.sysConfig.Get(string(common.ConfigKeySMTPPassword)),
		FromName: s.sysConfig.Get(string(common.ConfigKeySMTPFromName)),
	}
	return mailer.SendVerificationCode(smtpCfg, email, code)
}

func (s *Service) Register(email, code, password, name string) (string, string, error) {
	if s.sysConfig.GetBool(string(common.ConfigKeyEnableEmailVerification), true) {
		if !s.ValidateVerifyCode(email, code) {
			return "", "", http.ErrBodyNotAllowed
		}
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	u := User{Email: email, PasswordHash: string(hash), Name: name, Role: s.initialRole()}
	if err := database.DB.Create(&u).Error; err != nil {
		return "", "", gorm.ErrDuplicatedKey
	}
	access, refresh := s.IssueTokens(u.ID)
	return access, refresh, nil
}

func (s *Service) Login(email, password string) (string, string, error) {
	var u User
	if err := database.DB.Where("email = ?", email).First(&u).Error; err != nil {
		return "", "", gorm.ErrRecordNotFound
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", "", gorm.ErrInvalidData
	}
	access, refresh := s.IssueTokens(u.ID)
	return access, refresh, nil
}

func (s *Service) Refresh(refreshToken string) (string, string, error) {
	t, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) { return []byte(s.cfg.JWTSecret), nil })
	if err != nil || !t.Valid {
		return "", "", gorm.ErrInvalidData
	}
	claims, _ := t.Claims.(jwt.MapClaims)
	jti, _ := claims["jti"].(string)
	if jti != "" {
		if cache.RDB.Get(context.Background(), common.RedisKeyJWTBlacklist.F(jti)).Val() == "1" {
			return "", "", gorm.ErrInvalidTransaction
		}
	}
	uid := uint(claims["uid"].(float64))
	access, refresh := s.IssueTokens(uid)
	return access, refresh, nil
}

func (s *Service) Logout(refreshToken string) error {
	t, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) { return []byte(s.cfg.JWTSecret), nil })
	if err != nil || !t.Valid {
		return gorm.ErrInvalidData
	}
	claims, _ := t.Claims.(jwt.MapClaims)
	jti, _ := claims["jti"].(string)
	if jti != "" {
		return cache.RDB.Set(context.Background(), common.RedisKeyJWTBlacklist.F(jti), "1", time.Hour*24*7).Err()
	}
	return nil
}

type WechatTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
	OpenID       string `json:"openid"`
	Scope        string `json:"scope"`
	UnionID      string `json:"unionid"`
	ErrCode      int    `json:"errcode"`
	ErrMsg       string `json:"errmsg"`
}
type WechatUserInfo struct {
	OpenID     string `json:"openid"`
	Nickname   string `json:"nickname"`
	HeadImgURL string `json:"headimgurl"`
	UnionID    string `json:"unionid"`
}

func (s *Service) MakeWeChatLoginURL(state string) string {
	appID := s.sysConfig.Get(string(common.ConfigKeyWeChatAppID))
	redirectURI := s.sysConfig.Get(string(common.ConfigKeyWeChatRedirectURI))
	if !s.sysConfig.GetBool(string(common.ConfigKeyEnabledWechatLogin), true) || appID == "" || redirectURI == "" {
		return ""
	}
	params := url.Values{}
	params.Set("appid", appID)
	params.Set("redirect_uri", redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "snsapi_login")
	params.Set("state", state)
	return "https://open.weixin.qq.com/connect/qrconnect?" + params.Encode() + "#wechat_redirect"
}

func (s *Service) ExchangeCode(code string) (WechatTokenResponse, error) {
	var out WechatTokenResponse
	appID := s.sysConfig.Get(string(common.ConfigKeyWeChatAppID))
	appSecret := s.sysConfig.Get(string(common.ConfigKeyWeChatAppSecret))

	if appID == "" || appSecret == "" {
		return out, http.ErrNotSupported
	}
	u := "https://api.weixin.qq.com/sns/oauth2/access_token?appid=" + url.QueryEscape(appID) +
		"&secret=" + url.QueryEscape(appSecret) +
		"&code=" + url.QueryEscape(code) +
		"&grant_type=authorization_code"
	resp, err := http.Get(u)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out.ErrCode != 0 || out.OpenID == "" || out.AccessToken == "" {
		return out, http.ErrHandlerTimeout
	}
	return out, nil
}

func (s *Service) FetchUserInfo(accessToken, openid string) (WechatUserInfo, error) {
	var out WechatUserInfo
	if accessToken == "" || openid == "" {
		return out, http.ErrNotSupported
	}
	u := "https://api.weixin.qq.com/sns/userinfo?access_token=" + url.QueryEscape(accessToken) + "&openid=" + url.QueryEscape(openid)
	resp, err := http.Get(u)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out.OpenID == "" {
		return out, http.ErrHandlerTimeout
	}
	return out, nil
}

func (s *Service) FindOrCreateWeChatUser(ui WechatUserInfo, tr WechatTokenResponse) (User, error) {
	var user User
	var oa OAuthAccount
	q := database.DB.Where("provider = ? AND provider_union_id <> '' AND provider_union_id = ?", "wechat", ui.UnionID).First(&oa)
	if ui.UnionID != "" && q.Error == nil {
		if err := database.DB.First(&user, oa.UserID).Error; err == nil {
			return user, nil
		}
	}
	if err := database.DB.Where("provider = ? AND provider_open_id = ?", "wechat", tr.OpenID).First(&oa).Error; err == nil {
		if err := database.DB.First(&user, oa.UserID).Error; err == nil {
			return user, nil
		}
	}
	user = User{
		Name:      ui.Nickname,
		AvatarURL: ui.HeadImgURL,
		IsActive:  true,
		Role:      s.initialRole(),
	}
	if err := database.DB.Create(&user).Error; err != nil {
		return User{}, err
	}
	record := OAuthAccount{
		UserID:          user.ID,
		Provider:        common.Wechat,
		ProviderOpenID:  tr.OpenID,
		ProviderUnionID: ui.UnionID,
	}
	if b, e := json.Marshal(ui); e == nil {
		record.RawProfileJSON = string(b)
	}
	_ = database.DB.Create(&record).Error
	return user, nil
}

func (s *Service) SaveOAuthState(state string, data map[string]any) error {
	b, _ := json.Marshal(data)
	return cache.RDB.Set(context.Background(), common.RedisKeyOAuthState.F(state), string(b), 10*time.Minute).Err()
}

func (s *Service) GetOAuthState(state string) (string, error) {
	return cache.RDB.Get(context.Background(), common.RedisKeyOAuthState.F(state)).Result()
}

func (s *Service) DelOAuthState(state string) {
	_ = cache.RDB.Del(context.Background(), common.RedisKeyOAuthState.F(state)).Err()
}

func (s *Service) SaveOTT(ott string, payload map[string]any) error {
	b, _ := json.Marshal(payload)
	return cache.RDB.Set(context.Background(), common.RedisKeyOAuthOTT.F(ott), string(b), time.Minute).Err()
}

func (s *Service) GetOTT(ott string) (string, error) {
	return cache.RDB.Get(context.Background(), common.RedisKeyOAuthOTT.F(ott)).Result()
}

func (s *Service) DelOTT(ott string) {
	_ = cache.RDB.Set(context.Background(), common.RedisKeyOAuthOTT.F(ott), "", time.Second).Err()
	_ = cache.RDB.Del(context.Background(), common.RedisKeyOAuthOTT.F(ott)).Err()
}

func (s *Service) SanitizeUser(u User) map[string]any {
	return map[string]any{
		"id":        u.ID,
		"email":     u.Email,
		"name":      u.Name,
		"avatarUrl": u.AvatarURL,
		"language":  u.Language,
		"role":      u.Role,
	}
}

func (s *Service) IsAllowedOrigin(origins string, origin string) bool {
	if origin == "" {
		return false
	}
	for _, o := range splitCSV(origins) {
		if o == origin {
			return true
		}
	}
	return false
}

func splitCSV(s string) []string {
	var out []string
	start := 0
	for i := 0; i <= len(s); i++ {
		if i == len(s) || s[i] == ',' {
			if start < i {
				out = append(out, s[start:i])
			}
			start = i + 1
		}
	}
	return out
}

type GithubTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	Scope       string `json:"scope"`
}

type GithubUserInfo struct {
	ID        int    `json:"id"`
	Login     string `json:"login"`
	Name      string `json:"name"`
	Email     string `json:"email"`
	AvatarURL string `json:"avatar_url"`
}

func (s *Service) MakeGithubLoginURL(state string) string {
	clientID := s.sysConfig.Get(string(common.ConfigKeyGithubClientID))
	redirectURI := s.sysConfig.Get(string(common.ConfigKeyGithubRedirectURI))
	if !s.sysConfig.GetBool(string(common.ConfigKeyEnabledGithubLogin), true) || clientID == "" || redirectURI == "" {
		return ""
	}
	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("redirect_uri", redirectURI)
	params.Set("scope", "user:email")
	params.Set("state", state)
	return "https://github.com/login/oauth/authorize?" + params.Encode()
}

func (s *Service) ExchangeGithubCode(code string) (GithubTokenResponse, error) {
	var out GithubTokenResponse
	clientID := s.sysConfig.Get(string(common.ConfigKeyGithubClientID))
	clientSecret := s.sysConfig.Get(string(common.ConfigKeyGithubClientSecret))
	redirectURI := s.sysConfig.Get(string(common.ConfigKeyGithubRedirectURI))
	if clientID == "" || clientSecret == "" {
		return out, http.ErrNotSupported
	}
	u := "https://github.com/login/oauth/access_token"
	params := url.Values{}
	params.Set("client_id", clientID)
	params.Set("client_secret", clientSecret)
	params.Set("code", code)
	params.Set("redirect_uri", redirectURI)

	req, _ := http.NewRequest("POST", u, nil)
	req.URL.RawQuery = params.Encode()
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out.AccessToken == "" {
		return out, http.ErrHandlerTimeout
	}
	return out, nil
}

func (s *Service) FetchGithubUserInfo(accessToken string) (GithubUserInfo, error) {
	var out GithubUserInfo
	if accessToken == "" {
		return out, http.ErrNotSupported
	}
	u := "https://api.github.com/user"
	req, _ := http.NewRequest("GET", u, nil)
	req.Header.Set("Authorization", "Bearer "+accessToken)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return out, err
	}
	defer resp.Body.Close()
	_ = json.NewDecoder(resp.Body).Decode(&out)
	if out.ID == 0 {
		return out, http.ErrHandlerTimeout
	}
	return out, nil
}

func (s *Service) FindOrCreateGithubUser(ui GithubUserInfo) (User, error) {
	var user User
	var oa OAuthAccount
	providerOpenID := strconv.Itoa(ui.ID)

	// Try to find by oauth account
	if err := database.DB.Where("provider = ? AND provider_open_id = ?", "github", providerOpenID).First(&oa).Error; err == nil {
		if err := database.DB.First(&user, oa.UserID).Error; err == nil {
			return user, nil
		}
	}

	// Create new user
	user = User{
		Name:      ui.Name,
		Email:     ui.Email, // GitHub might return empty email if private, but we try
		AvatarURL: ui.AvatarURL,
		IsActive:  true,
		Role:      s.initialRole(),
	}
	if user.Name == "" {
		user.Name = ui.Login
	}

	if err := database.DB.Create(&user).Error; err != nil {
		return User{}, err
	}

	record := OAuthAccount{
		UserID:         user.ID,
		Provider:       common.Github,
		ProviderOpenID: providerOpenID,
	}
	if b, e := json.Marshal(ui); e == nil {
		record.RawProfileJSON = string(b)
	}
	_ = database.DB.Create(&record).Error
	return user, nil
}
