package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"pdf-management-api/config"
)

func ListPDF(w http.ResponseWriter, r *http.Request) {
	status := r.URL.Query().Get("status")
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	query := "SELECT id, filename, original_name, size, status, created_at FROM pdf_files"
	args := []interface{}{}

	if status != "" {
		query += " WHERE status = ?"
		args = append(args, status)
	}

	query += " ORDER BY created_at DESC LIMIT ? OFFSET ?"
	args = append(args, limit, offset)

	rows, err := config.DB.Query(query, args...)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer rows.Close()

	var result []map[string]interface{}

	for rows.Next() {
		var id, size int
		var filename, originalName, status string
		var createdAt string

		rows.Scan(&id, &filename, &originalName, &size, &status, &createdAt)

		result = append(result, map[string]interface{}{
			"id":            id,
			"filename":      filename,
			"original_name": originalName,
			"size":          size,
			"status":        status,
			"created_at":    createdAt,
		})
	}

	resp := map[string]interface{}{
		"success": true,
		"data":    result,
		"pagination": map[string]interface{}{
			"page":  page,
			"limit": limit,
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
