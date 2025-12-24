package common

import "fmt"

type RedisKey string

const (
	RedisKeyUVDay              RedisKey = "uv:%s"
	RedisKeyRateLimit          RedisKey = "rl:%s:%s"
	RedisKeyRateLimitUser      RedisKey = "rlu:%s:%s"
	RedisKeySysConfig          RedisKey = "sysconfig:%s"
	RedisKeyVerify             RedisKey = "verify:%s"
	RedisKeyJWTBlacklist       RedisKey = "jwt:blacklist:%s"
	RedisKeyOAuthState         RedisKey = "oauth:state:%s"
	RedisKeyOAuthOTT           RedisKey = "oauth:ott:%s"
	RedisKeyTemplatesListAll   RedisKey = "templates:list:all"
	RedisKeyViews              RedisKey = "views:%s"
	RedisKeyPublicResume       RedisKey = "public:resume:%s"
	RedisKeyCircuitBreaker     RedisKey = "cb:%s"
	RedisKeyCircuitBreakerFail RedisKey = "cb:%s:fail"
)

func (k RedisKey) F(args ...any) string {
	return fmt.Sprintf(string(k), args...)
}
