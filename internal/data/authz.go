package data

import (
	"fmt"

	"github.com/casbin/casbin/v3"
	pgxadapter "github.com/noho-digital/casbin-pgx-adapter"

	"github.com/tencat-dev/go-base/internal/biz"
)

var _ biz.PermissionChecker = (*CasbinAuthz)(nil)
var _ biz.PermissionManager = (*CasbinAuthz)(nil)

func NewPermissionChecker(c *CasbinAuthz) biz.PermissionChecker {
	return c
}
func NewPermissionManager(c *CasbinAuthz) biz.PermissionManager {
	return c
}

type CasbinAuthz struct {
	enforcer casbin.IEnforcer
}

func NewCasbinEnforcer(data *Data) (casbin.IEnforcer, error) {
	adapter, err := pgxadapter.NewAdapterWithPool(data.db.Pool,
		pgxadapter.WithTableName("casbin_rules"),        // Optional: custom table name
		pgxadapter.WithIndex("ptype", "v0", "v1", "v2"), // policy: sub, obj, act
		pgxadapter.WithIndex("ptype", "v0", "v1"),       // grouping: user -> role
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create adapter: %v", err)
	}

	e, err := casbin.NewSyncedCachedEnforcer("configs/rbac_model.conf", adapter)
	if err != nil {
		return nil, err
	}

	if err = e.LoadPolicy(); err != nil {
		return nil, err
	}

	e.EnableAutoSave(true)
	return e, nil
}

func NewCasbinAuthz(enforcer casbin.IEnforcer) (*CasbinAuthz, error) {
	return &CasbinAuthz{enforcer: enforcer}, nil
}

func (c *CasbinAuthz) Can(sub, obj, act string) (bool, error) {
	return c.enforcer.Enforce(sub, obj, act)
}

func (c *CasbinAuthz) GrantRole(userID, role string) error {
	_, err := c.enforcer.AddGroupingPolicy(userID, role)
	return err
}

func (c *CasbinAuthz) RevokeRole(userID, role string) error {
	_, err := c.enforcer.RemoveGroupingPolicy(userID, role)
	return err
}

func (c *CasbinAuthz) GrantPermission(sub, obj, act string) error {
	_, err := c.enforcer.AddPolicy(sub, obj, act)
	return err
}

func (c *CasbinAuthz) RevokePermission(sub, obj, act string) error {
	_, err := c.enforcer.RemovePolicy(sub, obj, act)
	return err
}
