package api

import (
	"encoding/json"
	"html/template"
	"log"
	"net/http"

	"integritypos/internal/hardware"
	"integritypos/internal/models"
	"integritypos/internal/pos"

	"github.com/jackc/pgx/v5/pgxpool"
)

func HandlePOS(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/pos.html")
	if err != nil {
		log.Printf("ERROR: Failed to parse template: %v", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.Execute(w, nil); err != nil {
		log.Printf("ERROR: Failed to execute template: %v", err)
	}
}

func HandleCheckout(db *pgxpool.Pool) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Only allow POST
		if r.Method != http.MethodPost {
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
			return
		}

		var items []models.OrderItem
		if err := json.NewDecoder(r.Body).Decode(&items); err != nil {
			log.Printf("ERROR: Failed to decode checkout payload: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Invalid request payload"})
			return
		}

		if len(items) == 0 {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			json.NewEncoder(w).Encode(map[string]string{"error": "Ticket is empty"})
			return
		}

		err := pos.ProcessOrder(r.Context(), db, items)
		if err != nil {
			log.Printf("ERROR: Transaction failed: %v", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
			return
		}

		var total float64
		for _, item := range items {
			total += item.Subtotal
		}

		// Enviar ticket físico (falla de impresión no aborta transacción)
		// Utilizamos orderID = 0 ya que ProcessOrder actualmente no lo retorna, en un flujo completo lo retornaríamos.
		if err := hardware.PrintTicket(items, total, 0); err != nil {
			log.Printf("ERROR: Hardware printer failed for order: %v", err)
		}

		log.Printf("INFO: Order processed successfully with %d items", len(items))
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"message": "Transaction successful"})
	}
}
