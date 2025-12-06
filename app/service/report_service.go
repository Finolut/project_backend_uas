package service

import (
	"context"

	pgRepo "clean-arch-copy/app/repository/postgre"
)

// ReportService handles statistics and reporting functionality
type ReportService struct {
	achievementRefRepo pgRepo.AchievementRefRepository
	studentRepo        pgRepo.StudentRepository
	lecturerRepo       pgRepo.LecturerRepository
}

func NewReportService(
	achievementRefRepo pgRepo.AchievementRefRepository,
	studentRepo pgRepo.StudentRepository,
	lecturerRepo pgRepo.LecturerRepository,
) *ReportService {
	return &ReportService{
		achievementRefRepo: achievementRefRepo,
		studentRepo:        studentRepo,
		lecturerRepo:       lecturerRepo,
	}
}

// AchievementStatistics holds statistics data
type AchievementStatistics struct {
	TotalAchievements    int              `json:"total_achievements"`
	AchievementsByStatus map[string]int   `json:"achievements_by_status"`
	AchievementsByType   map[string]int   `json:"achievements_by_type"`
	TopStudents          []TopStudentData `json:"top_students"`
	VerificationRate     float64          `json:"verification_rate"`
}

type TopStudentData struct {
	StudentID        string `json:"student_id"`
	StudentName      string `json:"student_name"`
	AchievementCount int    `json:"achievement_count"`
}

// GetAllAchievementsStatistics returns overall statistics for all achievements
func (s *ReportService) GetAllAchievementsStatistics(ctx context.Context) (*AchievementStatistics, error) {
	stats := &AchievementStatistics{
		AchievementsByStatus: make(map[string]int),
		AchievementsByType:   make(map[string]int),
		TopStudents:          []TopStudentData{},
	}

	// TODO: Implement actual statistics calculation using database queries
	// This requires repository methods to fetch and aggregate data
	// For now, returning empty statistics structure

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

func (s *ReportService) GetAchievementHistory(ctx context.Context, refID string) (map[string]interface{}, error) {
    // Gunakan activityLogRepo (pastikan struct ReportService punya field activityLogRepo)
    // logs, err := s.activityLogRepo.ListByEntity(ctx, "achievement_reference", refID, 100, 0)
    // return map[string]interface{}{"history": logs}, err
    return map[string]interface{}{}, nil // Placeholder sampai repo di-inject
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
