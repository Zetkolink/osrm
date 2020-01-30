package rest

type UploadRequest struct {
	Comment  string `json:"comment"`
	UploadBy int    `json:"upload_by"`
}
