package model

import "go-starter-project/pkg/auth"

func GetDefaultPermission(roleLevel int32) auth.StaffPermissionsMap {
	m := make(auth.StaffPermissionsMap)

	// TODO: grant permission following role level
	switch roleLevel {
	case auth.RoleSuperAdmin:
		m.GrantGlobalPermission(auth.GlobalPermCreateTag)
		m.GrantGlobalPermission(auth.GlobalPermUpdateTag)
		m.GrantGlobalPermission(auth.GlobalPermCreateStaffUser)
	case auth.RoleShopStaff:
		m.GrantPermission(auth.PermCreateTag)
		m.GrantPermission(auth.PermUpdateTag)
	}

	return m
}
