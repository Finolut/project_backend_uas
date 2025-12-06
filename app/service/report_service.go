package service

import (
	"context"
	"sort"

	pgRepo "clean-arch-copy/app/repository/postgre"
)

// ReportService handles statistics and reporting functionality
type ReportService struct {
	achievementRefRepo pgRepo.AchievementRefRepository
	studentRepo        pgRepo.StudentRepository
	lecturerRepo       pgRepo.LecturerRepository
	activityLogRepo    pgRepo.ActivityLogRepository // <-- Tambahkan ini
}

// Update Constructor: Tambahkan parameter activityLogRepo
func NewReportService(
	achievementRefRepo pgRepo.AchievementRefRepository,
	studentRepo pgRepo.StudentRepository,
	lecturerRepo pgRepo.LecturerRepository,
	activityLogRepo pgRepo.ActivityLogRepository, // <-- Tambahkan parameter
) *ReportService {
	return &ReportService{
		achievementRefRepo: achievementRefRepo,
		studentRepo:        studentRepo,
		lecturerRepo:       lecturerRepo,
		activityLogRepo:    activityLogRepo, // <-- Assign
	}
}

// AchievementStatistics holds statistics data
type AchievementStatistics struct {
	TotalAchievements    int              `json:"total_achievements"`
	AchievementsByStatus map[string]int   `json:"achievements_by_status"`
	TopStudents          []TopStudentData `json:"top_students"`
	VerificationRate     float64          `json:"verification_rate"`
}

type TopStudentData struct {
	StudentID        string `json:"student_id"`
	StudentName      string `json:"student_name"` // Note: Nama mungkin butuh fetch terpisah jika tidak join
	AchievementCount int    `json:"achievement_count"`
}

// GetAllAchievementsStatistics returns overall statistics for all achievements
func (s *ReportService) GetAllAchievementsStatistics(ctx context.Context) (*AchievementStatistics, error) {
	// 1. Ambil semua data (untuk skala besar, sebaiknya gunakan Query COUNT/GROUP BY di repository)
	refs, err := s.achievementRefRepo.ListAll(ctx)
	if err != nil {
		return nil, err
	}

	stats := &AchievementStatistics{
		AchievementsByStatus: make(map[string]int),
		TopStudents:          []TopStudentData{},
	}

	stats.TotalAchievements = len(refs)
	verifiedCount := 0
	studentCounts := make(map[string]int)

	// 2. Agregasi Data in-memory
	for _, ref := range refs {
		// Hitung per status
		stats.AchievementsByStatus[ref.Status]++

		// Hitung verified
		if ref.Status == "verified" {
			verifiedCount++
		}

		// Hitung per mahasiswa
		studentCounts[ref.StudentID]++
	}

	// 3. Hitung Rate
	if stats.TotalAchievements > 0 {
		stats.VerificationRate = float64(verifiedCount) / float64(stats.TotalAchievements)
	}

	// 4. Cari Top 5 Students
	// Konversi map ke slice untuk sorting
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range studentCounts {
		ss = append(ss, kv{k, v})
	}

	// Sort descending berdasarkan jumlah prestasi
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})

	// Ambil top 5
	limit := 5
	if len(ss) < limit {
		limit = len(ss)
	}

	for i := 0; i < limit; i++ {
		// Optional: Fetch nama student jika diperlukan
		// st, _ := s.studentRepo.GetByID(ctx, ss[i].Key)
		// name := "Unknown"
		// if st != nil { name = st.StudentID } // Atau fetch user untuk nama asli

		stats.TopStudents = append(stats.TopStudents, TopStudentData{
			StudentID:        ss[i].Key,
			StudentName:      ss[i].Key, // Gunakan ID dulu untuk efisiensi
			AchievementCount: ss[i].Value,
		})
	}

	return stats, nil
}

// GetStudentStatistics returns statistics for a specific student
func (s *ReportService) GetStudentStatistics(ctx context.Context, studentID string) (map[string]interface{}, error) {
	result := make(map[string]interface{})

	// Get student basic info
	student, err := s.studentRepo.GetByID(ctx, studentID)
	if err != nil {
		return nil, err
	}

	if student == nil {
		return nil, ErrNotFound
	}

	result["student_id"] = student.ID
	result["student_code"] = student.StudentID
	result["program_study"] = student.Program
	result["academic_year"] = student.AcademicYear

	// Get student achievements
	achievements, err := s.achievementRefRepo.ListByStudent(ctx, studentID)
	if err != nil {
		return nil, err
	}

	totalAchievements := len(achievements)
	statusCount := make(map[string]int)
	verifiedCount := 0

	for _, ach := range achievements {
		statusCount[ach.Status]++
		if ach.Status == "verified" {
			verifiedCount++
		}
	}

	result["total_achievements"] = totalAchievements
	result["achievements_by_status"] = statusCount
	result["verified_count"] = verifiedCount
	result["draft_count"] = statusCount["draft"]
	result["submitted_count"] = statusCount["submitted"]
	result["rejected_count"] = statusCount["rejected"]

	if totalAchievements > 0 {
		result["verification_rate"] = float64(verifiedCount) / float64(totalAchievements)
	} else {
		result["verification_rate"] = 0.0
	}

	return result, nil
}

// GetAchievementHistory retrieves activity logs for a specific achievement reference
func (s *ReportService) GetAchievementHistory(ctx context.Context, refID string) (map[string]interface{}, error) {
	// Panggil repository activity log
	logs, err := s.activityLogRepo.ListByEntity(ctx, "achievement_reference", refID, 100, 0)
	if err != nil {
		return nil, err
	}

	// Bungkus dalam map agar format JSON rapi: { "history": [...] }
	return map[string]interface{}{
		"entity_id": refID,
		"history":   logs,
	}, nil
}

var ErrNotFound = &CustomError{"resource_not_found", "resource not found", 404}

type CustomError struct {
	Code    string
	Message string
	Status  int
}

func (e *CustomError) Error() string {
	return e.Message
}
