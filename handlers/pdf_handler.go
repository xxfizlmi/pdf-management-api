package handlers

import (
	"encoding/json"
	"net/http"
	"os"
	"time"
	"pdf-management-api/config"
	"github.com/google/uuid"
	"github.com/jung-kurt/gofpdf"
)

type GeneratePDFRequest struct {
	Title       string `json:"title"`
	Institution string `json:"institution"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Content     string `json:"content"`
}

func GeneratePDF(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req GeneratePDFRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	os.MkdirAll("uploads/pdf", os.ModePerm)
	filename := "report_" + uuid.New().String() + ".pdf"
	filepath := "uploads/pdf/" + filename

	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, req.Title)
	pdf.Ln(12)

	pdf.SetFont("Arial", "", 12)
	pdf.MultiCell(0, 8, req.Content, "", "", false)

	err = pdf.OutputFileAndClose(filepath)
	if err != nil {
		http.Error(w, "Failed to generate PDF", http.StatusInternalServerError)
		return
	}

	resp := map[string]interface{}{
		"status":  "success",
		"message": "PDF generated successfully",
		"data": map[string]interface{}{
			"filename":  filename,
			"filePath":  filepath,
			"status":    "CREATED",
			"createdAt": time.Now(),
		},
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)

	config.DB.Exec(
	`INSERT INTO pdf_files (filename, filepath, status, created_at)
	 VALUES (?, ?, ?, ?)`,
	filename, filepath, "CREATED", time.Now(),
)

}
