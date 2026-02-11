package common

type ConfigKey string

const (
	ConfigKeyEnableEmailVerification ConfigKey = "enable_email_verification"
	ConfigKeyEnabledWechatLogin      ConfigKey = "enabled_wechat_login"
	ConfigKeyEnabledWechatMPLogin    ConfigKey = "enabled_wechat_mp_login"
	ConfigKeyEnabledGithubLogin      ConfigKey = "enabled_github_login"
	ConfigKeyEnabledPricingPage      ConfigKey = "enabled_pricing_page"
	ConfigKeySMTPHost                ConfigKey = "smtp_host"
	ConfigKeySMTPPort                ConfigKey = "smtp_port"
	ConfigKeySMTPUsername            ConfigKey = "smtp_username"
	ConfigKeySMTPPassword            ConfigKey = "smtp_password"
	ConfigKeySMTPFromName            ConfigKey = "smtp_from_name"
	ConfigKeyWeChatAppID             ConfigKey = "wechat_app_id"
	ConfigKeyWeChatAppSecret         ConfigKey = "wechat_app_secret"
	ConfigKeyWeChatRedirectURI       ConfigKey = "wechat_redirect_uri"
	ConfigKeyWeChatMPAppID           ConfigKey = "wechat_mp_app_id"
	ConfigKeyWeChatMPAppSecret       ConfigKey = "wechat_mp_app_secret"
	ConfigKeyWeChatMPToken           ConfigKey = "wechat_mp_token"
	ConfigKeyWeChatMPAESKey          ConfigKey = "wechat_mp_aes_key"
	ConfigKeyGithubClientID          ConfigKey = "github_client_id"
	ConfigKeyGithubClientSecret      ConfigKey = "github_client_secret"
	ConfigKeyGithubRedirectURI       ConfigKey = "github_redirect_uri"
	ConfigKeyEnabledStorageS3        ConfigKey = "enabled_storage_s3"
	ConfigKeyStorageS3Bucket         ConfigKey = "storage_s3_bucket"
	ConfigKeyStorageS3Region         ConfigKey = "storage_s3_region"
	ConfigKeyStorageS3Endpoint       ConfigKey = "storage_s3_endpoint"
	ConfigKeyStorageS3AccessKey      ConfigKey = "storage_s3_access_key"
	ConfigKeyStorageS3SecretKey      ConfigKey = "storage_s3_secret_key"
	ConfigKeyChromeAPIURL            ConfigKey = "chrome_api_url"
	ConfigKeyCORSOrigins             ConfigKey = "cors_origins"
	ConfigKeyFrontendBaseURL         ConfigKey = "frontend_base_url"
)

type Role string

const (
	RoleAdmin Role = "admin"
	RoleUser  Role = "user"
)

type ProviderType string

const (
	Github   ProviderType = "github"
	Wechat   ProviderType = "wechat"
	WechatMP ProviderType = "wechat_mp"
)
