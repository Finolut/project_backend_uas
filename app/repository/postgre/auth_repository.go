package repository

import (
	model "clean-arch/app/model/postgre"
	"database/sql"
	"errors"
	"strings"

	"github.com/google/uuid"
)

// GetUserByUsernameOrEmail fetches a User and their password hash along with role name.
func GetUserByUsernameOrEmail(db *sql.DB, identifier string) (*model.User, string, string, error) {
	var user model.User
	var roleName string

	query := `SELECT 
	          u.id, u.username, u.email, u.password_hash, u.full_name, u.role_id, u.is_active, r.name 
	          FROM users u
	          JOIN roles r ON u.role_id = r.id
	          WHERE u.username = $1 OR u.email = $1`

	err := db.QueryRow(query, identifier).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash,
		&user.FullName, &user.RoleID, &user.IsActive, &roleName,
	)

	if err != nil {
		return nil, "", "", err
	}

	return &user, user.PasswordHash, roleName, nil
}

// Mengganti GetAlumniByNIM
func GetStudentByStudentID(db *sql.DB, studentID string) (*model.User, *model.Student, string, error) {
	var student model.Student
	var user model.User
	var roleName string

	query := `SELECT 
	          u.id, u.username, u.email, u.password_hash, u.full_name, r.name,
	          s.id, s.student_id, s.program_study, s.academic_year, s.advisor_id 
	          FROM students s
	          JOIN users u ON s.user_id = u.id
			  JOIN roles r ON u.role_id = r.id
	          WHERE s.student_id = $1`

	err := db.QueryRow(query, studentID).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName, &roleName,
		&student.ID, &student.StudentID, &student.ProgramStudy, &student.AcademicYear, &student.AdvisorID,
	)

	if err != nil {
		return nil, nil, "", err
	}

	return &user, &student, roleName, nil
}

// GetRoleByName fetches a Role by its name.
func GetRoleByName(db *sql.DB, name string) (*model.Role, error) {
	var role model.Role
	query := `SELECT id, name, description, created_at FROM roles WHERE name = $1`
	err := db.QueryRow(query, name).Scan(&role.ID, &role.Name, &role.Description, &role.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &role, nil
}

// CreateUserAndStudent creates a new User and a linked Student record.
func CreateUserAndStudent(db *sql.DB, req model.RegisterStudentRequest, hashedPassword string, roleID string) (*model.User, *model.Student, error) {
	tx, err := db.Begin()
	if err != nil {
		return nil, nil, err
	}
	defer tx.Rollback()

	// 1. Create User
	userUUID := uuid.New().String()
	userQuery := `INSERT INTO users (id, username, email, password_hash, full_name, role_id, is_active) 
	              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING created_at, updated_at`

	var user model.User
	err = tx.QueryRow(userQuery, userUUID, req.Username, req.Email, hashedPassword, req.FullName, roleID, true).Scan(
		&user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		// Periksa jika error unik constraint (username, email)
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			return nil, nil, errors.New("username atau email sudah terdaftar")
		}
		return nil, nil, err
	}
	user.ID = userUUID
	user.Username = req.Username
	user.Email = req.Email
	user.PasswordHash = hashedPassword
	user.FullName = req.FullName
	user.RoleID = roleID
	user.IsActive = true

	// 2. Create Student
	studentUUID := uuid.New().String()
	studentQuery := `INSERT INTO students (id, user_id, student_id, program_study, academic_year) 
	                 VALUES ($1, $2, $3, $4, $5) RETURNING created_at`

	var student model.Student
	err = tx.QueryRow(studentQuery, studentUUID, user.ID, req.StudentID, req.ProgramStudy, req.AcademicYear).Scan(
		&student.CreatedAt,
	)
	if err != nil {
		// Periksa jika error unik constraint (student_id)
		if strings.Contains(err.Error(), "unique constraint") || strings.Contains(err.Error(), "duplicate key") {
			return nil, nil, errors.New("NIM sudah terdaftar")
		}
		return nil, nil, err
	}
	student.ID = studentUUID
	student.UserID = user.ID
	student.StudentID = req.StudentID
	student.ProgramStudy = req.ProgramStudy
	student.AcademicYear = req.AcademicYear

	if err := tx.Commit(); err != nil {
		return nil, nil, err
	}

	return &user, &student, nil
}

// GetUserByID is a utility function to fetch a user by ID
func GetUserByID(db *sql.DB, userID string) (*model.User, error) {
	var user model.User
	query := `SELECT id, username, email, password_hash, full_name, role_id, is_active, created_at, updated_at 
	          FROM users WHERE id = $1`
	err := db.QueryRow(query, userID).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.FullName,
		&user.RoleID, &user.IsActive, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return &user, nil
}
