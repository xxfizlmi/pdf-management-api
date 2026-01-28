package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"pdf-management-api/config"
)

func UploadPDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		jsonError(w, http.StatusMethodNotAllowed, "Method not allowed", "METHOD_NOT_ALLOWED")
		return
	}

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		jsonError(w, http.StatusBadRequest, "File size exceeds maximum limit (10MB)", "FILE_TOO_LARGE")
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		jsonError(w, http.StatusBadRequest, "File is required", "FILE_REQUIRED")
		return
	}

	if !strings.HasSuffix(strings.ToLower(handler.Filename), ".pdf") {
		jsonError(w, http.StatusBadRequest, "Only PDF files are allowed", "INVALID_EXTENSION")
		return
	}

	buffer := make([]byte, 512)
	file.Read(buffer)
	file.Seek(0, 0)

	if http.DetectContentType(buffer) != "application/pdf" {
		jsonError(w, http.StatusBadRequest, "Invalid MIME type, must be application/pdf", "INVALID_MIME")
		return
	}

	os.MkdirAll("uploads/pdf", os.ModePerm)
	filename := "uploaded_" + uuid.New().String() + ".pdf"
	filepath := filepath.Join("uploads/pdf", filename)

	dst, err := os.Create(filepath)
	if err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	defer dst.Close()

	_, err = dst.ReadFrom(file)
	if err != nil {
		http.Error(w, "Failed to write file", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"status":  "success",
		"message": "File uploaded successfully",
		"data": map[string]interface{}{
			"originalName": handler.Filename,
			"filename":     filename,
			"filePath":     filepath,
			"size":         handler.Size,
			"status":       "UPLOADED",
			"uploadedAt":   time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)

	config.DB.Exec(
		`INSERT INTO pdf_files (filename, original_name, filepath, size, status, created_at)
	 VALUES (?, ?, ?, ?, ?, ?)`,
		filename, handler.Filename, filepath, handler.Size, "UPLOADED", time.Now(),
	)

}

func jsonError(w http.ResponseWriter, status int, msg string, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"success":    false,
		"message":    msg,
		"error_code": code,
	})
}
