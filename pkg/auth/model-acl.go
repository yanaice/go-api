package auth

type StaffPermission interface {
	GetPermissionKey() string
}

// STAFF
type StaffPermissionsMap map[string]bool
type StaffGlobalPermission string
type StaffShopPermission string

const (
	PermCreateTag StaffShopPermission = "shop_create_tag"
	PermUpdateTag StaffShopPermission = "shop_update_tag"
)

const (
	GlobalPermCreateTag       StaffGlobalPermission = "global_create_tag"
	GlobalPermUpdateTag       StaffGlobalPermission = "global_update_tag"
	GlobalPermCreateStaffUser StaffGlobalPermission = "global_create_staff_user"
)

func (u StaffPermissionsMap) GrantPermission(permission StaffPermission) {
	u[permission.GetPermissionKey()] = true
}

func (u StaffPermissionsMap) GrantGlobalPermission(permission StaffGlobalPermission) {
	u.GrantPermission(permission)
}

func (u StaffShopPermission) GetPermissionKey() string {
	return string(u)
}

func (u StaffGlobalPermission) GetPermissionKey() string {
	return string(u)
}
