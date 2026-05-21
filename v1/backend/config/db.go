package config

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

// DB es una variable global accesible desde otros paquetes
var DB *sql.DB

func ConectarDB() {
	var err error
	dsn := "clickforge_user:clickforge_user_password@tcp(127.0.0.1:3306)/clickforge_db"

	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("❌ Error al configurar la base de datos: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("❌ No se pudo conectar a MySQL en Docker: %v", err)
	}

	fmt.Println("⚡ Conectado con éxito a MySQL (clickforge_db)")
}
