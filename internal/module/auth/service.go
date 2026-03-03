package auth

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha1"
	"encoding/base64"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"cvforge/internal/common"
	"cvforge/internal/infra/cache"
	"cvforge/internal/infra/config"
	conf "cvforge/internal/module/config"
	"cvforge/internal/pkg/mailer"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Service struct {
	repo      *Repo
	sysConfig *conf.Service
}

var (
	ErrInvalidEmail      = errors.New("invalid_email")
	ErrInvalidVerifyCode = errors.New("invalid_verify_code")
)

func NewService() *Service {
	return &Service{
		repo:      DefaultRepo(),
		sysConfig: conf.NewService(),
	}
}

func (s *Service) FeatureGithubEnabled() bool {
	return s.sysConfig.GetBool(string(common.ConfigKeyEnabledGithubLogin), true)
}

func (s *Service) FrontendBase() string {
	return s.sysConfig.Get(string(common.ConfigKeyFrontendBaseURL))
}

func (s *Service) IssueTokens(uid uint) (string, string) {
	var ver int = 1
	if v, err := s.repo.GetUserTokenVersion(uid); err == nil && v > 0 {
		ver = v
	}
	mk := func(exp time.Duration) string {
		t := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{"uid": uid, "ver": ver, "exp": time.Now().Add(exp).Unix(), "jti": uuid.NewString()})
		signed, _ := t.SignedString([]byte(config.CF.JWTSecret))
		return signed
	}
	return mk(2 * time.Hour), mk(7 * 24 * time.Hour)
}

func (s *Service) TouchLastLoginAt(uid uint) {
	if uid == 0 {
		return
	}
	now := time.Now()
	_ = s.repo.SetUserLastLoginAt(uid, now)
}

func (s *Service) initialRole() common.Role {
	n, err := s.repo.CountUsers()
	if err == nil && n == 0 {
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
	email = strings.TrimSpace(email)
	if email == "" {
		return "", "", ErrInvalidEmail
	}
	if s.sysConfig.GetBool(string(common.ConfigKeyEnableEmailVerification), true) {
		if !s.ValidateVerifyCode(email, code) {
			return "", "", ErrInvalidVerifyCode
		}
	}
	hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	now := time.Now()
	u := User{Email: &email, PasswordHash: string(hash), Name: name, Role: s.initialRole(), LastLoginAt: &now}
	if err := s.repo.CreateUser(&u); err != nil {
		return "", "", err
	}
	access, refresh := s.IssueTokens(u.ID)
	return access, refresh, nil
}

func (s *Service) Login(email, password string) (string, string, error) {
	email = strings.TrimSpace(email)
	if email == "" {
		return "", "", gorm.ErrRecordNotFound
	}
	u, err := s.repo.FindUserByEmail(email)
	if err != nil {
		return "", "", gorm.ErrRecordNotFound
	}
	if bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", "", gorm.ErrInvalidData
	}
	if !u.IsActive {
		return "", "", gorm.ErrInvalidTransaction
	}
	s.TouchLastLoginAt(u.ID)
	access, refresh := s.IssueTokens(u.ID)
	return access, refresh, nil
}

func (s *Service) Refresh(refreshToken string) (string, string, error) {
	t, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) { return []byte(config.CF.JWTSecret), nil })
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
	u, err := s.repo.FindUserByID(uid)
	if err != nil {
		return "", "", gorm.ErrRecordNotFound
	}
	if !u.IsActive {
		return "", "", gorm.ErrInvalidTransaction
	}
	access, refresh := s.IssueTokens(uid)
	return access, refresh, nil
}

func (s *Service) Logout(refreshToken string) error {
	t, err := jwt.Parse(refreshToken, func(t *jwt.Token) (interface{}, error) { return []byte(config.CF.JWTSecret), nil })
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

func (s *Service) LogoutAccess(accessToken string) error {
	t, err := jwt.Parse(accessToken, func(t *jwt.Token) (interface{}, error) { return []byte(config.CF.JWTSecret), nil })
	if err != nil || !t.Valid {
		return gorm.ErrInvalidData
	}
	claims, _ := t.Claims.(jwt.MapClaims)
	jti, _ := claims["jti"].(string)
	var ttl time.Duration = 2 * time.Hour
	if exp, ok := claims["exp"].(float64); ok {
		remain := time.Until(time.Unix(int64(exp), 0))
		if remain > 0 {
			ttl = remain
		}
	}
	if jti != "" {
		return cache.RDB.Set(context.Background(), common.RedisKeyJWTBlacklist.F(jti), "1", ttl).Err()
	}
	return nil
}

func (s *Service) GlobalLogout(uid uint) error {
	if err := s.repo.IncrementUserTokenVersion(uid); err != nil {
		return err
	}
	if cache.RDB != nil {
		_ = cache.RDB.Set(context.Background(), common.RedisKeyUserTokenVersion.F(uid), "inc", time.Minute).Err()
	}
	return nil
}

func (s *Service) LogoutSession(accessToken, refreshToken string) error {
	accessToken = strings.TrimSpace(accessToken)
	if accessToken != "" {
		if uid := s.uidFromAccessToken(accessToken); uid > 0 {
			_ = s.GlobalLogout(uid)
		}
		_ = s.LogoutAccess(accessToken)
	}
	return s.Logout(refreshToken)
}

func (s *Service) uidFromAccessToken(accessToken string) uint {
	t, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, errors.New("invalid signing algorithm")
		}
		return []byte(config.CF.JWTSecret), nil
	})
	if err != nil || !t.Valid {
		return 0
	}
	claims, ok := t.Claims.(jwt.MapClaims)
	if !ok {
		return 0
	}
	if uidF, ok := claims["uid"].(float64); ok && uidF > 0 {
		return uint(uidF)
	}
	return 0
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
	email := ""
	if u.Email != nil {
		email = strings.TrimSpace(*u.Email)
	}
	return map[string]any{
		"id":        u.ID,
		"email":     email,
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
	providerOpenID := strconv.Itoa(ui.ID)

	// Try to find by oauth account
	if oa, err := s.repo.FindOAuthAccountByProviderOpenID(common.Github, providerOpenID); err == nil {
		if u, err2 := s.repo.FindUserByID(oa.UserID); err2 == nil {
			user = u
			return user, nil
		}
	}

	// Create new user
	var email *string
	if e := strings.TrimSpace(ui.Email); e != "" {
		email = &e
	}
	user = User{
		Name:      ui.Name,
		Email:     email,
		AvatarURL: ui.AvatarURL,
		IsActive:  true,
		Role:      s.initialRole(),
	}
	if user.Name == "" {
		user.Name = ui.Login
	}

	if err := s.repo.CreateUser(&user); err != nil {
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
	_ = s.repo.CreateOAuthAccount(&record)
	return user, nil
}

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
	if openid == "" {
		openid = ui.OpenID
	}
	if openid == "" {
		return User{}, http.ErrNotSupported
	}

	if ui.UnionID != "" {
		if oa, err := s.repo.FindOAuthAccountByProviderUnionID(common.WechatMP, ui.UnionID); err == nil {
			if user, err2 := s.repo.FindUserByID(oa.UserID); err2 == nil {
				return user, nil
			}
		}
	}
	if oa, err := s.repo.FindOAuthAccountByProviderOpenID(common.WechatMP, openid); err == nil {
		if user, err2 := s.repo.FindUserByID(oa.UserID); err2 == nil {
			return user, nil
		}
	}

	name := strings.TrimSpace(ui.Nickname)
	if name == "" {
		name = "用户_" + randAlphaNum(5)
	}
	user := User{
		Name:      name,
		AvatarURL: ui.HeadImgURL,
		IsActive:  true,
		Role:      s.initialRole(),
	}
	if err := s.repo.CreateUser(&user); err != nil {
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
	_ = s.repo.CreateOAuthAccount(&record)
	return user, nil
}

func (s *Service) MarkWeChatMPSceneReady(ctx context.Context, scene, openid string, user User) (string, error) {
	if cache.RDB == nil || scene == "" || openid == "" {
		return "", http.ErrNotSupported
	}

	s.TouchLastLoginAt(user.ID)
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

func randAlphaNum(n int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if n <= 0 {
		return ""
	}
	out := make([]byte, n)
	for i := 0; i < n; i++ {
		v, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		out[i] = charset[v.Int64()]
	}
	return string(out)
}
