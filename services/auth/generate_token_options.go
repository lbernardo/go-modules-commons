package auth

import "time"

type generateTokenOption struct {
	Expiration    time.Time
	DoesNotExpire bool
}

// GenerateTokenOption is the generate token options, like expire date and token without expire date
type GenerateTokenOption func(*generateTokenOption)

// WithExpirationOption set expiration time to token (default is 24 hours)
func WithExpirationOption(expiration time.Time) GenerateTokenOption {
	return func(o *generateTokenOption) {
		o.Expiration = expiration
	}
}

// WithDoesNotExpireOption set token without expire date
func WithDoesNotExpireOption() GenerateTokenOption {
	return func(o *generateTokenOption) {
		o.DoesNotExpire = true
	}
}
