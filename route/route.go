package route

import (
	"context"
	"fmt"
	"os"
	"time"

	mongoModel "clean-arch-copy/app/model/mongo"
	pgModel "clean-arch-copy/app/model/postgre"
	"clean-arch-copy/app/service"
	"clean-arch-copy/middleware"
	"clean-arch-copy/utils"

	"github.com/gofiber/fiber/v2"
)

// RegisterRoutes mendaftarkan semua endpoint API ke dalam Fiber App
func RegisterRoutes(app *fiber.App, s *service.Services) {

	// 1. Global Middleware
	app.Use(middleware.RequestID())
	app.Use(middleware.Helmet())
	app.Use(middleware.RateLimiter())
	// app.Use(middleware.Logger()) // Opsional, sudah ada di config/app.go

	// Helper untuk context dengan timeout standar
	timeoutContext := func(c *fiber.Ctx) (context.Context, context.CancelFunc) {
		return context.WithTimeout(c.Context(), 10*time.Second)
	}

	// Wrapper untuk RBAC Permission Checker agar sesuai signature middleware
	rbacCheck := func(roleID string, permission string) (bool, error) {
		// Gunakan context background karena pengecekan permission biasanya cepat/cached
		return s.RBAC.HasPermissionByRoleID(context.Background(), roleID, permission)
	}

	// API Group Base
	api := app.Group("/api/v1")

	// =========================================================================
	// 5.1 AUTHENTICATION
	// =========================================================================
	authGroup := api.Group("/auth")

	// POST /auth/login
	authGroup.Post("/login", middleware.LoginRateLimiter(), func(c *fiber.Ctx) error {
		var req pgModel.LoginRequest // Pastikan struct ini ada di model, atau pakai struct inline
		if err := c.BodyParser(&req); err != nil {
			return utils.JSONError(c, fiber.StatusBadRequest, "Invalid request body")
		}

		ctx, cancel := timeoutContext(c)
		defer cancel()

		token, user, err := s.Auth.Login(ctx, req.Username, req.Password)
		if err != nil {
			return utils.JSONError(c, fiber.StatusUnauthorized, err.Error())
		}

		return utils.JSONSuccess(c, fiber.StatusOK, fiber.Map{
			"token": token,
			"user":  user,
		})
	})

	// POST /auth/refresh
	authGroup.Post("/refresh", middleware.NewJWTMiddleware(), func(c *fiber.Ctx) error {
		userID := c.Locals(middleware.LocalsUserID).(string)
		ctx, cancel := timeoutContext(c)
		defer cancel()

		newToken, err := s.Auth.Refresh(ctx, userID)
		if err != nil {
			return utils.JSONError(c, fiber.StatusUnauthorized, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, fiber.Map{"token": newToken})
	})

	// POST /auth/logout
	authGroup.Post("/logout", middleware.NewJWTMiddleware(), func(c *fiber.Ctx) error {
		// Ambil token mentah dari header untuk diblacklist
		authHeader := c.Get("Authorization")
		if len(authHeader) < 7 {
			return utils.JSONError(c, fiber.StatusBadRequest, "Invalid header")
		}
		tokenString := authHeader[7:] // Remove "Bearer "

		ctx, cancel := timeoutContext(c)
		defer cancel()

		if err := s.Auth.Logout(ctx, tokenString); err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, "Logged out successfully")
	})

	// GET /auth/profile
	authGroup.Get("/profile", middleware.NewJWTMiddleware(), func(c *fiber.Ctx) error {
		userID := c.Locals(middleware.LocalsUserID).(string)
		ctx, cancel := timeoutContext(c)
		defer cancel()

		user, err := s.User.GetByID(ctx, userID)
		if err != nil {
			return utils.JSONError(c, fiber.StatusNotFound, "User not found")
		}
		return utils.JSONSuccess(c, fiber.StatusOK, user)
	})

	// =========================================================================
	// 5.2 USERS (ADMIN)
	// =========================================================================
	// Group ini dilindungi Auth & RBAC (misal permission: 'user:manage')
	userGroup := api.Group("/users", middleware.NewJWTMiddleware())

	// GET /users
	userGroup.Get("/", middleware.RequirePermission(rbacCheck, "user:read"), func(c *fiber.Ctx) error {
		ctx, cancel := timeoutContext(c)
		defer cancel()
		users, err := s.User.ListAll(ctx)
		if err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, users)
	})

	// POST /users
	userGroup.Post("/", middleware.RequirePermission(rbacCheck, "user:create"), func(c *fiber.Ctx) error {
		var u pgModel.User
		if err := c.BodyParser(&u); err != nil {
			return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
		}

		// Hash password manual disini atau di service (di code sebelumnya service hash password tidak dipanggil di Register)
		// Kita asumsikan hash dilakukan di Handler atau Service memanggil utils.HashPassword
		hashed, _ := s.Auth.HashPassword(u.PasswordHash) // field ini biasanya menampung plain text saat request
		u.PasswordHash = hashed

		ctx, cancel := timeoutContext(c)
		defer cancel()

		if err := s.User.Register(ctx, &u); err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusCreated, "User created")
	})

	// PUT /users/:id
	userGroup.Put("/:id", middleware.RequirePermission(rbacCheck, "user:update"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		var u pgModel.User
		if err := c.BodyParser(&u); err != nil {
			return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
		}
		u.ID = id

		ctx, cancel := timeoutContext(c)
		defer cancel()

		if err := s.User.Update(ctx, &u); err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, "User updated")
	})

	// PUT /users/:id/role (Assign Role)
	userGroup.Put("/:id/role", middleware.RequirePermission(rbacCheck, "user:assign_role"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		var req struct {
			RoleID string `json:"role_id"`
		}
		if err := c.BodyParser(&req); err != nil {
			return utils.JSONError(c, fiber.StatusBadRequest, "Invalid body")
		}

		ctx, cancel := timeoutContext(c)
		defer cancel()

		if err := s.User.UpdateRole(ctx, id, req.RoleID); err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, "Role updated")
	})

	// DELETE /users/:id
	userGroup.Delete("/:id", middleware.RequirePermission(rbacCheck, "user:delete"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		ctx, cancel := timeoutContext(c)
		defer cancel()

		if err := s.User.Delete(ctx, id); err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, "User deleted")
	})

	// =========================================================================
	// 5.5 STUDENTS & LECTURERS
	// =========================================================================
	studentGroup := api.Group("/students", middleware.NewJWTMiddleware())
	lecturerGroup := api.Group("/lecturers", middleware.NewJWTMiddleware())

	// GET /students (List)
	studentGroup.Get("/", func(c *fiber.Ctx) error {
		ctx, cancel := timeoutContext(c)
		defer cancel()
		list, err := s.Student.ListAll(ctx)
		if err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, list)
	})

	// GET /students/:id
	studentGroup.Get("/:id", func(c *fiber.Ctx) error {
		// Note: ID disini bisa StudentID (UUID) atau NIM, sesuaikan dengan logic service
		// Service GetByID expect UUID table ID
		id := c.Params("id")
		ctx, cancel := timeoutContext(c)
		defer cancel()
		st, err := s.Student.GetByID(ctx, id)
		if err != nil {
			return utils.JSONError(c, fiber.StatusNotFound, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, st)
	})

	// PUT /students/:id/advisor (Set Advisor) - Admin Only
	studentGroup.Put("/:id/advisor", middleware.RequirePermission(rbacCheck, "student:manage"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		var req struct {
			AdvisorID string `json:"advisor_id"`
		}
		if err := c.BodyParser(&req); err != nil {
			return utils.JSONError(c, fiber.StatusBadRequest, "Invalid body")
		}

		ctx, cancel := timeoutContext(c)
		defer cancel()

		if err := s.Student.UpdateAdvisor(ctx, id, &req.AdvisorID); err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, "Advisor updated")
	})

	// GET /lecturers
	lecturerGroup.Get("/", func(c *fiber.Ctx) error {
		ctx, cancel := timeoutContext(c)
		defer cancel()
		list, err := s.Lecturer.ListAll(ctx)
		if err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, list)
	})

	// GET /lecturers/:id/advisees (Mahasiswa Bimbingan)
	lecturerGroup.Get("/:id/advisees", func(c *fiber.Ctx) error {
		id := c.Params("id") // Lecturer ID
		ctx, cancel := timeoutContext(c)
		defer cancel()
		list, err := s.Lecturer.GetAdvisees(ctx, id)
		if err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, list)
	})

	// =========================================================================
	// 5.4 ACHIEVEMENTS (CORE)
	// =========================================================================
	achGroup := api.Group("/achievements", middleware.NewJWTMiddleware())

	// GET /achievements (List All - Filtered by Service logic)
	// Permission: Admin atau Lecturer (lihat semua/bimbingan), Student (lihat punya sendiri biasanya via endpoint profile)
	achGroup.Get("/", func(c *fiber.Ctx) error {
		// TODO: Parse query params for filtering
		ctx, cancel := timeoutContext(c)
		defer cancel()
		list, err := s.Achievement.GetAllAchievements(ctx, nil)
		if err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, list)
	})

	// POST /achievements (Create Draft - Mahasiswa)
	achGroup.Post("/", middleware.RequirePermission(rbacCheck, "achievement:create"), func(c *fiber.Ctx) error {
		userID := c.Locals(middleware.LocalsUserID).(string)
		var doc mongoModel.Achievement
		if err := c.BodyParser(&doc); err != nil {
			return utils.JSONError(c, fiber.StatusBadRequest, "Invalid body")
		}

		ctx, cancel := timeoutContext(c)
		defer cancel()

		result, err := s.Achievement.CreateDraft(ctx, userID, &doc)
		if err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusCreated, result)
	})

	// GET /achievements/:id (Detail)
	achGroup.Get("/:id", func(c *fiber.Ctx) error {
		id := c.Params("id")
		ctx, cancel := timeoutContext(c)
		defer cancel()

		mongoData, pgRef, err := s.Achievement.GetDetail(ctx, id)
		if err != nil {
			return utils.JSONError(c, fiber.StatusNotFound, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, fiber.Map{
			"reference": pgRef,
			"detail":    mongoData,
		})
	})

	// PUT /achievements/:id (Update Draft - Mahasiswa)
	achGroup.Put("/:id", middleware.RequirePermission(rbacCheck, "achievement:update"), func(c *fiber.Ctx) error {
		// ... isi handler update ...
		return utils.JSONSuccess(c, fiber.StatusOK, "Draft updated")
	})

	// DELETE /achievements/:id (Delete Draft - Mahasiswa)
	achGroup.Delete("/:id", middleware.RequirePermission(rbacCheck, "achievement:delete"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals(middleware.LocalsUserID).(string)

		ctx, cancel := timeoutContext(c)
		defer cancel()

		if err := s.Achievement.DeleteDraft(ctx, id, userID); err != nil {
			return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, "Draft deleted")
	})

	// POST /achievements/:id/submit (Submit for Verification - Mahasiswa)
	achGroup.Post("/:id/attachments", middleware.RequirePermission(rbacCheck, "achievement:update"), func(c *fiber.Ctx) error {
		refID := c.Params("id")
		userID := c.Locals(middleware.LocalsUserID).(string)

		// 1. Ambil File
		file, err := c.FormFile("file")
		if err != nil {
			return utils.JSONError(c, fiber.StatusBadRequest, "File upload failed: "+err.Error())
		}

		// 2. Simpan File (Pastikan folder ada)
		uniqueName := fmt.Sprintf("%d-%s", time.Now().Unix(), file.Filename)
		savePath := fmt.Sprintf("./uploads/%s", uniqueName)

		if _, err := os.Stat("./uploads"); os.IsNotExist(err) {
			os.Mkdir("./uploads", 0755)
		}

		if err := c.SaveFile(file, savePath); err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, "Cannot save file")
		}

		// 3. Siapkan Data
		attachmentData := mongoModel.Attachment{
			FileName: file.Filename,
			URL:      "/uploads/" + uniqueName,
			MimeType: file.Header.Get("Content-Type"),
			Size:     file.Size,
		}

		// 4. Panggil Service
		ctx, cancel := timeoutContext(c)
		defer cancel()

		if err := s.Achievement.AddAttachment(ctx, refID, userID, attachmentData); err != nil {
			os.Remove(savePath) // Hapus file jika gagal simpan DB
			return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
		}

		return utils.JSONSuccess(c, fiber.StatusOK, attachmentData)
	})

achGroup.Post("/:id/submit", middleware.RequirePermission(rbacCheck, "achievement:submit"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		userID := c.Locals(middleware.LocalsUserID).(string)

		ctx, cancel := timeoutContext(c)
		defer cancel()

		// Logic validasi "Hanya Mahasiswa" terjadi di dalam fungsi s.Achievement.Submit ini
		if err := s.Achievement.Submit(ctx, id, userID); err != nil {
			return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, "Achievement submitted successfully")
	})

	// POST /achievements/:id/verify (Verify - Dosen Wali)
	achGroup.Post("/:id/verify", middleware.RequirePermission(rbacCheck, "achievement:verify"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		verifierID := c.Locals(middleware.LocalsUserID).(string)

		ctx, cancel := timeoutContext(c)
		defer cancel()

		if err := s.Achievement.Verify(ctx, id, verifierID); err != nil {
			return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, "Achievement verified")
	})

	// POST /achievements/:id/reject (Reject - Dosen Wali)
	achGroup.Post("/:id/reject", middleware.RequirePermission(rbacCheck, "achievement:verify"), func(c *fiber.Ctx) error {
		id := c.Params("id")
		verifierID := c.Locals(middleware.LocalsUserID).(string)

		var req struct {
			Note string `json:"note"`
		}
		if err := c.BodyParser(&req); err != nil {
			return utils.JSONError(c, fiber.StatusBadRequest, "Note is required")
		}

		ctx, cancel := timeoutContext(c)
		defer cancel()

		if err := s.Achievement.Reject(ctx, id, verifierID, req.Note); err != nil {
			return utils.JSONError(c, fiber.StatusBadRequest, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, "Achievement rejected")
	})

	// GET /achievements/:id/history (History Log)
	achGroup.Get("/:id/history", func(c *fiber.Ctx) error {
		id := c.Params("id")
		ctx, cancel := timeoutContext(c)
		defer cancel()

		hist, err := s.Report.GetAchievementHistory(ctx, id)
		if err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, hist)
	})

	// =========================================================================
	// 5.8 REPORTS & ANALYTICS
	// =========================================================================
	reportGroup := api.Group("/reports", middleware.NewJWTMiddleware())

	// GET /reports/statistics (Global Stats - Admin/Dosen)
	reportGroup.Get("/statistics", middleware.RequirePermission(rbacCheck, "report:view"), func(c *fiber.Ctx) error {
		ctx, cancel := timeoutContext(c)
		defer cancel()

		stats, err := s.Report.GetAllAchievementsStatistics(ctx)
		if err != nil {
			return utils.JSONError(c, fiber.StatusInternalServerError, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, stats)
	})

	// GET /reports/student/:id (Individual Stats)
	reportGroup.Get("/student/:id", func(c *fiber.Ctx) error {
		studentID := c.Params("id") // User ID or Student ID logic depends on implementation
		ctx, cancel := timeoutContext(c)
		defer cancel()

		stats, err := s.Report.GetStudentStatistics(ctx, studentID)
		if err != nil {
			return utils.JSONError(c, fiber.StatusNotFound, err.Error())
		}
		return utils.JSONSuccess(c, fiber.StatusOK, stats)
	})
}
