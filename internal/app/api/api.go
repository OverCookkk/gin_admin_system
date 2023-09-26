package api

import "github.com/google/wire"

var ApiSet = wire.NewSet(
	MenuApiSet,
	RoleApiSet,
	UserApiSet,
	LoginApiSet,
)
