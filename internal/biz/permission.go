package biz

type PermissionChecker interface {
	Can(sub, obj, act string) (bool, error)
}

type PermissionManager interface {
	GrantRole(userID, role string) error
	RevokeRole(userID, role string) error
	GrantPermission(subject, object, action string) error
	RevokePermission(subject, object, action string) error
}
