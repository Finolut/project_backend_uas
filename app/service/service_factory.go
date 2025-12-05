package service

import (
	"database/sql"

	mongodriver "go.mongodb.org/mongo-driver/mongo"

	mongoRepo "app/repository/mongo"
	pgRepo "app/repository/postgre"
)

// Repos set of repo interfaces needed to create services
type Repos struct {
	UserRepo           pgRepo.UserRepository
	RoleRepo           pgRepo.RoleRepository
	PermissionRepo     pgRepo.PermissionRepository
	RolePermissionRepo pgRepo.RolePermissionRepository
	StudentRepo        pgRepo.StudentRepository
	LecturerRepo       pgRepo.LecturerRepository
	AchievementRefRepo pgRepo.AchievementRefRepository
	AchievementRepo    mongoRepo.AchievementRepository
}

type Services struct {
	Achievement *AchievementService
	User        *UserService
	Auth        *AuthService
	RBAC        *RBACService
	Student     *StudentService
	Lecturer    *LecturerService
}

func NewServices(db *sql.DB, mongoDB *mongodriver.Database, repos *Repos) *Services {
	// if repos nil, create default repo impls elsewhere (not included)
	achSvc := NewAchievementService(repos.AchievementRepo, repos.AchievementRefRepo, repos.StudentRepo, repos.UserRepo)
	userSvc := NewUserService(repos.UserRepo)
	authSvc := NewAuthService(repos.UserRepo) // needs user repo for login
	rbacSvc := NewRBACService(repos.RolePermissionRepo, repos.PermissionRepo, repos.RoleRepo)
	studentSvc := NewStudentService(repos.StudentRepo)
	lecturerSvc := NewLecturerService(repos.LecturerRepo)

	return &Services{
		Achievement: achSvc,
		User:        userSvc,
		Auth:        authSvc,
		RBAC:        rbacSvc,
		Student:     studentSvc,
		Lecturer:    lecturerSvc,
	}
}
