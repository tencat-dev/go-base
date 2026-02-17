package biz

// AuthzBiz is a Auth usecase.
type AuthzBiz struct {
	pm PermissionManager
}

// NewAuthzBiz new a Auth usecase.
func NewAuthzBiz(pm PermissionManager) *AuthzBiz {
	return &AuthzBiz{
		pm: pm,
	}
}

func (b *AuthzBiz) GrantRole(userID, role string) error {
	return b.pm.GrantRole(userID, role)
}

func (b *AuthzBiz) RevokeRole(userID, role string) error {
	return b.pm.RevokeRole(userID, role)
}

func (b *AuthzBiz) GrantPermission(
	subject, object, action string,
) error {
	return b.pm.GrantPermission(subject, object, action)
}
