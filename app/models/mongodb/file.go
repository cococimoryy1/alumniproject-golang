package models

import "time"

type File struct {
    ID           int64     `json:"id" bson:"id"`
    FileName     string    `json:"file_name" bson:"file_name"`
    OriginalName string    `json:"original_name" bson:"original_name"`
    FilePath     string    `json:"file_path" bson:"file_path"`
    FileSize     int64     `json:"file_size" bson:"file_size"`
    FileType     string    `json:"file_type" bson:"file_type"`
    UploadedAt   time.Time `json:"uploaded_at" bson:"uploaded_at"`
}

type FileResponse struct {
    ID           int64     `json:"id"`
    FileName     string    `json:"file_name"`
    OriginalName string    `json:"original_name"`
    FilePath     string    `json:"file_path"`
    FileSize     int64     `json:"file_size"`
    FileType     string    `json:"file_type"`
    UploadedAt   time.Time `json:"uploaded_at"`
}
