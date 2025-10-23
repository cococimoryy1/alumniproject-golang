package models

// MetaInfo -> informasi pagination, sorting, dan search
type MetaInfo struct {
	Page   int    `json:"page" bson:"page"`
	Limit  int    `json:"limit" bson:"limit"`
	Total  int    `json:"total" bson:"total"`
	Pages  int    `json:"pages" bson:"pages"`
	SortBy string `json:"sortBy" bson:"sortBy"`
	Order  string `json:"order" bson:"order"`
	Search string `json:"search" bson:"search"`
}

// AlumniResponse -> response untuk endpoint /alumni
type AlumniResponse struct {
	Data []*Alumni  `json:"data" bson:"data"` // gunakan pointer slice
	Meta *MetaInfo  `json:"meta" bson:"meta"`
}

// PekerjaanResponse -> response untuk endpoint /pekerjaan
type PekerjaanResponse struct {
	Data []*Pekerjaan `json:"data" bson:"data"` // gunakan pointer slice
	Meta *MetaInfo    `json:"meta" bson:"meta"`
}
