package app

import (
	"context"
	"errors"
	"fmt"
	"gin_admin_system/internal/app/dao/menu"
	"gin_admin_system/internal/app/dao/role"
	"gin_admin_system/internal/app/types"
	"gin_admin_system/pkg/logger"
	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	"github.com/google/wire"
)

var CasbinAdapterSet = wire.NewSet(wire.Struct(new(CasbinAdapter), "*"), wire.Bind(new(persist.Adapter), new(*CasbinAdapter)))

// 实现Adapter接口
type CasbinAdapter struct {
	RoleRepo         *role.RoleRepo
	RoleMenuRepo     *role.RoleMenuRepo
	MenuResourceRepo *menu.MenuActionResourceRepo
}

func (a *CasbinAdapter) LoadPolicy(model casbinModel.Model) error {
	ctx := context.Background()
	err := a.loadRolePolicy(ctx, model)
	if err != nil {
		logger.WithContext(ctx).Errorf("Load casbin role policy error: %s", err.Error())
		return err
	}

	// TODO:用户、角色权限加载
	// err = a.loadUserPolicy(ctx, model)
	// if err != nil {
	// 	logger.WithContext(ctx).Errorf("Load casbin user policy error: %s", err.Error())
	// 	return err
	// }
	return nil
}

func (a *CasbinAdapter) loadRolePolicy(ctx context.Context, m casbinModel.Model) error {
	roleResult, err := a.RoleRepo.Query(ctx, types.RoleQueryReq{
		Status: 1,
	})
	if err != nil {
		return err
	} else if len(roleResult.Data) == 0 {
		return nil
	}

	roleMenuResult, err := a.RoleMenuRepo.Query(ctx, types.RoleMenuQueryReq{})
	if err != nil {
		return err
	}
	roleMenuMap := make(map[uint64]types.RoleMenus) // role_id,  菜单id跟actionID的组合
	for _, v := range roleMenuResult.Data {
		roleMenuMap[v.RoleID] = append(roleMenuMap[v.RoleID], &v)
	}

	menuResourceResult, err := a.MenuResourceRepo.Query(ctx, types.MenuActionResourceQueryReq{})
	if err != nil {
		return err
	}
	menuResourceMap := make(map[uint64]types.MenuActionResources) // action_id,  method和path的组合
	for _, v := range menuResourceResult.Data {
		menuResourceMap[v.ActionID] = append(menuResourceMap[v.ActionID], &v)
	}

	for _, item := range roleResult.Data { // 所有role信息
		mcache := make(map[string]struct{})
		if rms, ok := roleMenuMap[item.ID]; ok {
			for _, actionID := range toActionIDs(rms) { // toActionIDs把actionID全部拎出来形成一个切片
				if mrs, ok := menuResourceMap[actionID]; ok {
					for _, mr := range mrs {
						if mr.Path == "" || mr.Method == "" {
							continue
						} else if _, ok := mcache[mr.Path+mr.Method]; ok { // method和path的去重
							continue
						}
						mcache[mr.Path+mr.Method] = struct{}{}
						line := fmt.Sprintf("p,%d,%s,%s", item.ID, mr.Path, mr.Method)
						err := persist.LoadPolicyLine(line, m)
						if err != nil {
							return err
						}
					}
				}
			}
		}
	}

	return nil
}

// ToActionIDs 转换为动作ID列表
func toActionIDs(a types.RoleMenus) []uint64 {
	idList := make([]uint64, len(a))
	m := make(map[uint64]struct{})
	for i, item := range a {
		if _, ok := m[item.ActionID]; ok {
			continue
		}
		idList[i] = item.ActionID
		m[item.ActionID] = struct{}{}
	}
	return idList
}

// SavePolicy 报错所有的policy到持久层
func (a *CasbinAdapter) SavePolicy(model casbinModel.Model) error {
	return nil
}

// AddPolicy 添加一个policy规则至持久层
func (a *CasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

// RemovePolicy 从持久层删除单条policy规则
func (a *CasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	return errors.New("not implemented")
}

//  RemoveFilteredPolicy 从持久层删除符合筛选条件的policy规则
func (a *CasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	return errors.New("not implemented")
}
