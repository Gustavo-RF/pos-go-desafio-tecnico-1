package interfaces

type RateLimiterInterface interface {
	SetKey(key, value string, expireTime int) error
	GetKey(key string) (*string, error)
	IncKey(key string) error
	KeyExists(key string) bool
	IsBlocked(key string) bool
	ShouldBlock(key string) bool
}
