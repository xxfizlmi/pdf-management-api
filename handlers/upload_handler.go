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
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	r.ParseMultipartForm(10 << 20) // 10 MB

	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if !strings.HasSuffix(handler.Filename, ".pdf") {
		http.Error(w, "Only PDF files are allowed", http.StatusBadRequest)
		return
	}

	buffer := make([]byte, 512)
	file.Read(buffer)
	filetype := http.DetectContentType(buffer)
	file.Seek(0, 0)

	if filetype != "application/pdf" {
		http.Error(w, "Invalid file type", http.StatusBadRequest)
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
