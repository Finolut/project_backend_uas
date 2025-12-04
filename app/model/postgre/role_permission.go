package postgres

type RolePermission struct {
	RoleID       string `db:"role_id" json:"role_id"`
	PermissionID string `db:"permission_id" json:"permission_id"`
	// composite PK (role_id, permission_id) at DB level
}
