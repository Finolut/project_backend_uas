package service

import (
	"context"

	pgRepo "app/repository/postgre"
)

// RBACService checks role permissions
type RBACService struct {
	rolePermRepo pgRepo.RolePermissionRepository
	permRepo     pgRepo.PermissionRepository
	roleRepo     pgRepo.RoleRepository
}

func NewRBACService(rp pgRepo.RolePermissionRepository, pr pgRepo.PermissionRepository, rr pgRepo.RoleRepository) *RBACService {
	return &RBACService{
		rolePermRepo: rp,
		permRepo:     pr,
		roleRepo:     rr,
	}
}

// HasPermissionByRoleID returns true if the role has permission name (e.g. "achievement:verify")
func (s *RBACService) HasPermissionByRoleID(ctx context.Context, roleID string, permName string) (bool, error) {
	perms, err := s.rolePermRepo.ListByRole(ctx, roleID)
	if err != nil {
		return false, err
	}
	for _, p := range perms {
		if p.Name == permName {
			return true, nil
		}
	}
	return false, nil
}
