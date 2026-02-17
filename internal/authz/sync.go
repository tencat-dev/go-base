package authz

import (
	"github.com/casbin/casbin/v3"

	authzv1 "github.com/tencat-dev/go-base/api/authz/v1"
)

func SyncFromRegistry(e casbin.IEnforcer, r *AuthzRegistry) error {
	m := r.data.Load().(map[string]*authzv1.PermissionOption)

	// Ước lượng capacity để giảm re-alloc
	policies := make([][]string, 0, len(m)*2)

	// dedup trong memory (tránh duplicate trong cùng 1 proto load)
	seen := make(map[string]struct{})

	for _, perm := range m {
		for _, role := range perm.Roles {

			key := role + "|" + perm.Object + "|" + perm.Action
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}

			policies = append(policies, []string{
				role,
				perm.Object,
				perm.Action,
			})
		}
	}

	if len(policies) == 0 {
		return nil
	}

	// 🔥 Batch + ignore duplicate DB entries
	_, err := e.AddPolicies(policies)
	if err != nil {
		return err
	}

	return nil
}
