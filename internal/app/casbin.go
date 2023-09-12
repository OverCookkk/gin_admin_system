package app

import (
	"gin_admin_system/internal/app/config"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/persist"
	"time"
)

// InitCasbin 入参：接口函数
func InitCasbin(adapter persist.Adapter) (*casbin.SyncedEnforcer, func(), error) {
	cfg := config.C.Casbin
	if cfg.Model == "" {
		return new(casbin.SyncedEnforcer), nil, nil
	}

	e, err := casbin.NewSyncedEnforcer(cfg.Model)
	if err != nil {
		return nil, nil, err
	}

	err = e.InitWithModelAndAdapter(e.GetModel(), adapter)
	if err != nil {
		return nil, nil, err
	}
	e.EnableEnforce(cfg.Enable) // 根据配置决定是否开启casbin
	cleanFunc := func() {}
	if cfg.AutoLoad {
		// 启用定期自动加载策略规则
		e.StartAutoLoadPolicy(time.Duration(cfg.AutoLoadInternal) * time.Second)
		cleanFunc = func() {
			e.StopAutoLoadPolicy()
		}
	}

	return e, cleanFunc, nil
}
