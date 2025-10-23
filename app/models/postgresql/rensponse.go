package models

// MetaInfo -> informasi pagination, sorting, dan search
type MetaInfo struct {
    Page   int    `json:"page"`
    Limit  int    `json:"limit"`
    Total  int    `json:"total"`
    Pages  int    `json:"pages"`
    SortBy string `json:"sortBy"`
    Order  string `json:"order"`
    Search string `json:"search"`
}

// AlumniResponse -> response untuk endpoint /alumni
type AlumniResponse struct {
    Data []Alumni          `json:"data"`
    Meta *MetaInfo         `json:"meta"`
}

// PekerjaanResponse -> response untuk endpoint /pekerjaan
type PekerjaanResponse struct {
    Data []Pekerjaan       `json:"data"`
    Meta *MetaInfo         `json:"meta"`
}