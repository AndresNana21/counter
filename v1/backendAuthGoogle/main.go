package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Clave secreta ultra confidencial para firmar los tokens.
// En producción, esto DEBE ir en una variable de entorno (.env)
var miClaveSecreta = []byte("mi_clave_super_secreta_clickforge_2026")

// Estructuras para recibir y enviar datos en JSON
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func main() {
	mux := http.NewServeMux()

	// Rutas de la API
	mux.HandleFunc("/api/login", loginHandler)
	mux.HandleFunc("/api/profile", profileHandler)

	// Aplicamos el Middleware de CORS para que Astro (desde otro puerto) pueda conectarse
	fmt.Println("Servidor Go corriendo en http://localhost:8080 🚀")
	http.ListenAndServe(":8080", enableCORS(mux))
}

// 1. ENDPOINT DE LOGIN: Genera el Token y se lo envía a Astro
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		responderError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	var credenciales LoginRequest
	err := json.NewDecoder(r.Body).Decode(&credenciales)
	if err != nil {
		responderError(w, "JSON inválido", http.StatusBadRequest)
		return
	}

	// Simulación de validación (Aquí validarías contra tu base de datos MySQL)
	if credenciales.Username == "nana" && credenciales.Password == "password123" {
		
		// Creamos los "Claims" (la información que viajará dentro del token)
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"user_id":  123, // ID del usuario
			"username": credenciales.Username,
			"exp":      time.Now().Add(time.Hour * 24).Unix(), // El token expira en 24 horas
		})

		// Firmamos el token con nuestra clave secreta
		tokenString, err := token.SignedString(miClaveSecreta)
		if err != nil {
			responderError(w, "No se pudo generar el token", http.StatusInternalServerError)
			return
		}

		// Enviamos el token de vuelta a Astro
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(LoginResponse{Token: tokenString})
	} else {
		responderError(w, "Credenciales incorrectas", http.StatusUnauthorized)
	}
}

// 2. ENDPOINT PROTEGIDO: Recibe el token de Astro y lo valida
func profileHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		responderError(w, "Método no permitido", http.StatusMethodNotAllowed)
		return
	}

	// Extraemos la cabecera "Authorization"
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		responderError(w, "Falta la cabecera Authorization", http.StatusUnauthorized)
		return
	}

	// El formato es "Bearer <token>", así que separamos el string por el espacio
	partes := strings.Split(authHeader, " ")
	if len(partes) != 2 || partes[0] != "Bearer" {
		responderError(w, "Formato de token inválido (usa Bearer <token>)", http.StatusUnauthorized)
		return
	}

	tokenString := partes[1]

	// Parseamos y verificamos el token
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		// Validamos que el método de firma sea el mismo (HS256)
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("método de firma inesperado: %v", t.Header["alg"])
		}
		return miClaveSecreta, nil
	})

	// Si el token expiró, fue modificado o es falso, dará error
	if err != nil || !token.Valid {
		responderError(w, "Token inválido o expirado", http.StatusUnauthorized)
		return
	}

	// Si todo está bien, extraemos los datos que guardamos dentro del token (Claims)
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		responderError(w, "Error al leer los claims", http.StatusInternalServerError)
		return
	}

	// Responder con datos privados de la base de datos simulados
	w.Header().Set("Content-Type", "application/json")
	respuesta := map[string]interface{}{
		"mensaje":   "¡Bienvenido a la zona secreta!",
		"user_id":   claims["user_id"],
		"username":  claims["username"],
		"clicks":    9999, // Un dato de ejemplo para tu juego de clickers
		"servidor":  "Arch Linux Backend",
	}
	json.NewEncoder(w).Encode(respuesta)
}

// --- MIDDLEWARE PARA EVITAR EL ERROR DE CORS ---
func enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Permitimos que tu frontend de Astro se conecte (ajusta el puerto si Astro usa otro, ej: 4321)
		w.Header().Set("Access-Control-Allow-Origin", "http://localhost:4321") 
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, Authorization")

		// Si es una petición de tipo OPTIONS (preflight), respondemos OK inmediatamente
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Función auxiliar para responder errores en JSON limpiamente
func responderError(w http.ResponseWriter, mensaje string, codigo int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(codigo)
	json.NewEncoder(w).Encode(ErrorResponse{Error: mensaje})
}
