backend/
├── go.mod
├── go.sum
├── main.go               # Solo arranca el servidor
├── config/
│   └── db.go             # Configura y conecta MySQL
├── models/
│   └── juego.go          # Moldes de datos (Estructuras)
└── handlers/
    └── contador.go       # Lógica de las rutas (GET y POST)
