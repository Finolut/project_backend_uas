package model

// MetaInfo contains pagination and filtering information
type MetaInfo struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	Total  int    `json:"total"`
	Pages  int    `json:"pages"`
	SortBy string `json:"sortBy"`
	Order  string `json:"order"`
	Search string `json:"search"`
}

// AlumniResponse represents the response structure for alumni endpoints with pagination
type AlumniResponse struct {
	Data []Alumni `json:"data"`
	Meta MetaInfo `json:"meta"`
}

// PekerjaanResponse represents the response structure for pekerjaan endpoints with pagination
type PekerjaanResponse struct {
	Data []PekerjaanAlumni `json:"data"`
	Meta MetaInfo          `json:"meta"`
}

// PaginationParams represents query parameters for pagination, sorting, and search
type PaginationParams struct {
	Page   int    `json:"page"`
	Limit  int    `json:"limit"`
	SortBy string `json:"sortBy"`
	Order  string `json:"order"`
	Search string `json:"search"`
}
