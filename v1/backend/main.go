package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql" // El guion bajo indica que importamos el driver para que se registre en Go
)

// Las estructuras para los JSON se mantienen exactamente igual
type ContadorResponse struct {
	Valor int `json:"valor"`
}

type ActualizarRequest struct {
	Valor int `json:"valor"`
}

// Variable global para controlar la conexión a la Base de Datos
var db *sql.DB

func main() {
	var err error

	// 1. CONFIGURAR LA CONEXIÓN A MYSQL
	// Estructura del DSN: "usuario:contraseña@tcp(servidor:puerto)/nombre_base_de_datos"
	// CAMBIA "root" y "tu_password" por tus credenciales reales de MySQL
	// Estructura: "usuario:contraseña@tcp(servidor:puerto)/nombre_bd"
dsn := "clickforge_user:clickforge_user_password@tcp(127.0.0.1:3306)/clickforge_db"
	
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("Error al configurar la base de datos: %v", err)
	}
	defer db.Close() // Esto asegura que la BD se cierre cuando el programa termine

	// Verificar si la conexión física con MySQL realmente funciona
	err = db.Ping()
	if err != nil {
		log.Fatalf("No se pudo conectar a MySQL. ¿Está encendido el servicio?: %v", err)
	}
	fmt.Println("⚡ Conectado con éxito a MySQL (clickforge_db)")

	// 2. RUTAS
	http.HandleFunc("/api/contador", handleContador)
	http.HandleFunc("/api/contador/actualizar", handleActualizar)

	// 3. ARRANCAR EL SERVIDOR
	fmt.Println("🚀 Servidor Go corriendo en http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}

// --- MANEJADORES DE LAS RUTAS CON CONSULTAS SQL ---

// Ruta: GET http://localhost:3000/api/contador
func handleContador(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321")
	w.Header().Set("Content-Type", "application/json")

	if r.Method != http.MethodGet {
		http.Error(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// CONSULTA SQL: Traer el valor actual del contador donde id = 1
	var valorActual int
	query := "SELECT contador FROM juego WHERE id = 1"
	
	// QueryRow ejecuta la consulta y Scan guarda el resultado directamente en nuestra variable
	err := db.QueryRow(query).Scan(&valorActual)
	if err != nil {
		log.Printf("Error al consultar MySQL: %v", err)
		http.Error(w, "Error al leer los datos de la BD", http.StatusInternalServerError)
		return
	}

	// Responder a Astro con el valor real traído de MySQL
	respuesta := ContadorResponse{Valor: valorActual}
	json.NewEncoder(w).Encode(respuesta)
}

// Ruta: POST http://localhost:3000/api/contador/actualizar
func handleActualizar(w http.ResponseWriter, r *http.Request) {
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

	// Leer el JSON enviado desde el botón de Astro
	var peticion ActualizarRequest
	err := json.NewDecoder(r.Body).Decode(&peticion)
	if err != nil {
		http.Error(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// ACTUALIZACIÓN SQL: Guardar el nuevo número enviado en la base de datos
	query := "UPDATE juego SET contador = ? WHERE id = 1"
	
	// Exec ejecuta comandos que no devuelven filas (como UPDATE, INSERT o DELETE)
	_, err = db.Exec(query, peticion.Valor)
	if err != nil {
		log.Printf("Error al actualizar MySQL: %v", err)
		http.Error(w, "Error al guardar los datos en la BD", http.StatusInternalServerError)
		return
	}

	fmt.Printf("💾 Guardado en MySQL con éxito: %d\n", peticion.Valor)
	w.WriteHeader(http.StatusOK)
}
