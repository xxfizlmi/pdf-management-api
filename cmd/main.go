package main

import (
	"log"
	"net/http"
	"pdf-management-api/config"
	"pdf-management-api/handlers"
)

func main() {
	config.InitDB()
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("PDF Management API Running"))
	})

	mux.HandleFunc("/api/pdf/generate", handlers.GeneratePDF)
	mux.HandleFunc("/api/pdf/upload", handlers.UploadPDF)
	mux.HandleFunc("/api/pdf/list", handlers.ListPDF)
	mux.HandleFunc("/api/pdf/", handlers.DeletePDF)

	log.Println("Server running at :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
