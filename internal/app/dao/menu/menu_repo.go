package menu

import (
	"github.com/google/wire"
	"gorm.io/gorm"
)

var MenuSet = wire.NewSet(wire.Struct(new(MenuRepo), "*"))

type MenuRepo struct {
	DB *gorm.DB
}
