package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	// Reemplaza "clickforge-backend" por el nombre exacto que pusiste en tu go.mod
	"clickforge-backend/config"
	"clickforge-backend/models"
)

// HandleContador maneja: GET /api/contador
func HandleContador(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var valorActual int
	query := "SELECT contador FROM juego WHERE id = 1"
	
	// Usamos config.DB que es la conexión global que creamos en config/db.go
	err := config.DB.QueryRow(query).Scan(&valorActual)
	if err != nil {
		log.Printf("Error al consultar MySQL: %v", err)
		http.Error(w, "Error al leer la BD", http.StatusInternalServerError)
		return
	}

	respuesta := models.ContadorResponse{Valor: valorActual}
	json.NewEncoder(w).Encode(respuesta)
}

// HandleActualizar maneja: POST /api/contador/actualizar
func HandleActualizar(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var peticion models.ActualizarRequest
	err := json.NewDecoder(r.Body).Decode(&peticion)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	query := "UPDATE juego SET contador = ? WHERE id = 1"
	_, err = config.DB.Exec(query, peticion.Valor)
	if err != nil {
		log.Printf("Error al actualizar MySQL: %v", err)
		http.Error(w, "Error al guardar en la BD", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
