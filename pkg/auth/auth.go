package auth

import (
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const defaultKey = "gin-admin-system"

type options struct {
	signingMethod jwt.SigningMethod
	signingKey    interface{}
	keyfunc       jwt.Keyfunc
	expired       int
	tokenType     string
}

var defaultOptions = options{
	tokenType:     "Bearer",
	expired:       7200, // 单位s
	signingMethod: jwt.SigningMethodHS512,
	signingKey:    []byte(defaultKey),
	keyfunc: func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(defaultKey), nil
	},
}

type Option func(*options)

// SetSigningMethod 设定签名方式
func SetSigningMethod(method jwt.SigningMethod) Option {
	return func(o *options) {
		o.signingMethod = method
	}
}

// SetSigningKey 设定签名key
func SetSigningKey(key interface{}) Option {
	return func(o *options) {
		o.signingKey = key
	}
}

// SetKeyfunc 设定验证key的回调函数
func SetKeyfunc(keyFunc jwt.Keyfunc) Option {
	return func(o *options) {
		o.keyfunc = keyFunc
	}
}

// SetExpired 设定令牌过期时长(单位秒，默认7200)
func SetExpired(expired int) Option {
	return func(o *options) {
		o.expired = expired
	}
}

func New(opts ...Option) *JWTAuth {
	o := options{}
	for _, opt := range opts {
		opt(&o)
	}
	return &JWTAuth{opts: &o}
}

type JWTAuth struct {
	opts *options
}

// GenerateToken 生成令牌
func (a *JWTAuth) GenerateToken(ctx context.Context, userID string) (string, error) {
	expiresAt := time.Now().Add(time.Duration(a.opts.expired) * time.Second).Unix()
	token := jwt.NewWithClaims(a.opts.signingMethod, &jwt.StandardClaims{
		ExpiresAt: expiresAt,
		IssuedAt:  time.Now().Unix(),
		NotBefore: time.Now().Unix(),
		Subject:   userID, // userID保存在这个字段，验证的时候会取出来
	})

	// 加盐
	tokenString, err := token.SignedString(a.opts.signingKey)
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

// 解析令牌
func (a *JWTAuth) parseToken(tokenString string) (*jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, a.opts.keyfunc)
	if err != nil || !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims.(*jwt.StandardClaims), nil
}

func (a *JWTAuth) DestroyToken(ctx context.Context, tokenString string) error {
	claims, err := a.parseToken(tokenString)
	if err != nil {
		return err
	}

}
