package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"pdf-management-api/config"
)

func DeletePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// ambil id dari URL
	idStr := strings.TrimPrefix(r.URL.Path, "/api/pdf/")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	var status string
	err = config.DB.QueryRow(
		"SELECT status FROM pdf_files WHERE id = ?",
		id,
	).Scan(&status)

	if err != nil {
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}

	if status == "DELETED" {
		http.Error(w, "File already deleted", http.StatusBadRequest)
		return
	}

	_, err = config.DB.Exec(
		`UPDATE pdf_files
		 SET status = ?, deleted_at = ?
		 WHERE id = ?`,
		"DELETED", time.Now(), id,
	)

	if err != nil {
		http.Error(w, "Failed to delete file", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"success": true,
		"message": "PDF deleted successfully",
		"data": map[string]interface{}{
			"id":         id,
			"status":     "DELETED",
			"deleted_at": time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
