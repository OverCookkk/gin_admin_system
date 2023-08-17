package app

import (
	"gin_admin_system/internal/app/config"
	"gin_admin_system/internal/app/dao"
	"gin_admin_system/pkg/mgorm"
	"gorm.io/gorm"
)

func InitGormDB() (*gorm.DB, func(), error) {
	cfg := config.C
	var dsn string
	switch cfg.Gorm.DBType {
	case "mysql":
		dsn = cfg.MySQL.DSN()
	}

	cleanFunc := func() {}

	db, err := mgorm.NewGormDB(&mgorm.Config{
		Debug:        cfg.Gorm.Debug,
		DBType:       cfg.Gorm.DBType,
		DSN:          dsn,
		MaxIdleConns: cfg.Gorm.MaxIdleConns,
		MaxLifetime:  cfg.Gorm.MaxLifetime,
		MaxOpenConns: cfg.Gorm.MaxOpenConns,
		TablePrefix:  cfg.Gorm.TablePrefix,
	})
	if err != nil {
		return nil, cleanFunc, err
	}

	// 自动迁移创建表
	if cfg.Gorm.EnableAutoMigrate {
		err = dao.AutoMigrate(db)
		if err != nil {
			return nil, cleanFunc, err
		}

		// TODO: test
		/*m1 := menu.Menu{
		  	Name: "menu1",
		  }
		  m2 := menu.Menu{
		  	Name: "menu2",
		  }
		  mc1 := &menu.MenuAction{
		  	MenuID: 0,
		  	Code:   "mc1",
		  	Name:   "mc1",
		  	Menus:  []menu.Menu{m1, m2},
		  }
		  mc2 := &menu.MenuAction{
		  	MenuID: 0,
		  	Code:   "mc2",
		  	Name:   "mc2",
		  	Menus:  []menu.Menu{m2},
		  }
		  db.Create(mc1)
		  db.Create(mc2)*/
	}
	return db, cleanFunc, nil
}
