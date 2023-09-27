package app

import (
	"errors"
	"gin_admin_system/internal/app/config"
	"gin_admin_system/pkg/auth"
	"gin_admin_system/pkg/auth/store/buntdb"
	"gin_admin_system/pkg/auth/store/redis"
	"github.com/dgrijalva/jwt-go"
)

func InitAuth() (auth.JWTAuth, func(), error) {
	cfg := config.C.JWTAuth

	// 构造options
	var opts []auth.Option
	opts = append(opts, auth.SetExpired(cfg.Expired))
	opts = append(opts, auth.SetSigningKey([]byte(cfg.SigningKey)))
	opts = append(opts, auth.SetKeyfunc(func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token")
		}
		return []byte(cfg.SigningKey), nil
	}))

	var method jwt.SigningMethod
	switch cfg.SigningMethod {
	case "HS256":
		method = jwt.SigningMethodHS256
	case "HS384":
		method = jwt.SigningMethodHS384
	default:
		method = jwt.SigningMethodHS512
	}
	opts = append(opts, auth.SetSigningMethod(method))

	// 构造存储器
	var store auth.Storer
	switch cfg.Store {
	case "redis":
		rcfg := config.C.Redis
		store = redis.NewStore(&redis.Config{
			Addr:      rcfg.Addr,
			Password:  rcfg.Password,
			DB:        cfg.RedisDB,
			KeyPrefix: cfg.RedisPrefix,
		})
	default:
		s, err := buntdb.NewStore(cfg.FilePath)
		if err != nil {
			return auth.JWTAuth{}, nil, err
		}
		store = s
	}

	auth := auth.New(store, opts...)

	cleanFunc := func() {
		auth.Release()
	}
	return auth, cleanFunc, nil
}
