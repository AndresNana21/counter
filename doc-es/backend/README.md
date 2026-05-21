# Especificación de Requerimientos y Arquitectura del Backend (Go)

Este documento detalla la especificación técnica, los requerimientos funcionales y no funcionales, y la arquitectura de componentes para el backend del proyecto **Counter**, desarrollado en Go (Golang). El sistema está diseñado bajo un enfoque stateless y desacoplado, garantizando alta concurrencia, seguridad robusta y aislamiento de datos.

---

## 1. Arquitectura del Backend y Estructura del Proyecto

Para asegurar la escalabilidad, el mantenimiento y la facilidad de trabajo en equipo, el backend se estructura siguiendo el patrón de **Arquitectura en Capas** (Controlador, Servicio, Repositorio). Esto garantiza que la lógica de la base de datos esté completamente separada de las rutas HTTP.

### Estructura de Directorios Propuesta

```text
backend/
├── cmd/
│   └── api/
│       └── main.go          # Punto de entrada de la aplicación
├── internal/
│   ├── config/
│   │   └── config.go        # Carga de variables de entorno (.env)
│   ├── server/
│   │   └── server.go        # Configuración del servidor HTTP y rutas
│   ├── auth/
│   │   ├── handler.go       # Controladores de autenticación (Login HTTP)
│   │   ├── service.go       # Lógica de validación de Google JWT y generación de JWT local
│   │   └── jwt.go           # Funciones auxiliares de firmado y verificación de tokens
│   ├── counter/
│   │   ├── handler.go       # Controladores HTTP (Incrementar, Obtener)
│   │   ├── service.go       # Lógica de negocio del contador
│   │   └── repository.go    # Consultas SQL directas a MySQL
│   └── database/
│       └── mysql.go         # Conexión y pool de la base de datos
├── go.mod
└── go.sum

```

---

## 2. Requerimientos Funcionales (RF)

El backend debe cumplir estrictamente con los siguientes flujos de trabajo:

### RF1: Autenticación Federada mediante Google ID Token

* **Descripción:** El backend debe exponer un endpoint público (`POST /api/auth/google`).
* **Flujo:** 1. Recibe el token criptográfico emitido por Google desde el frontend en Astro.
2. Descarga las llaves públicas de Google y verifica la firma digital del token.
3. Valida que el campo `aud` (Audience) coincida exactamente con el `GOOGLE_CLIENT_ID` del proyecto.
4. Valida que el token no haya expirado (`exp`).

### RF2: Gestión y Registro Automatizado de Usuarios

* **Descripción:** Tras validar el token de Google, el sistema debe identificar al sujeto.
* **Flujo:**
1. Extrae el campo `sub` (ID único de Google), `email` y `name`.
2. Consulta la base de datos MySQL buscando el `sub`.
3. Si el usuario no existe, lo registra automáticamente en la tabla de usuarios con sus valores iniciales.
4. Si el usuario ya existe, recupera su información sin alterar sus datos históricos.



### RF3: Emisión de Sesión Local (JWT Propios)

* **Descripción:** El backend no expone el ciclo de vida de Google; genera su propia sesión.
* **Flujo:**
1. Genera un JSON Web Token (JWT) firmado localmente utilizando el algoritmo HS256 y una variable secreta (`JWT_SECRET`).
2. Introduce dentro de los *claims* del JWT el ID interno de la base de datos y el correo electrónico.
3. Establece un tiempo de expiración corto (ej. 24 horas).
4. Devuelve este JWT al frontend en Astro.



### RF4: Control Extrayente de Contadores por Usuario

* **Descripción:** El usuario autenticado debe poder interactuar con su contador individual.
* **Endpoints:**
* `GET /api/counter`: Recupera el valor actual del contador del usuario autenticado.
* `POST /api/counter/increment`: Incrementa en `+1` el valor del contador en la base de datos.


* **Seguridad:** Ambos endpoints requieren el paso por un Middleware de autenticación que valide el JWT local e inyecte el ID del usuario en el contexto de la petición.

---

## 3. Requerimientos No Funcionales (RNF)

### RNF1: Rendimiento y Concurrencia

* El backend debe aprovechar el modelo de goroutines de Go y utilizar un pool de conexiones optimizado para MySQL (`sql.DB`), evitando la apertura y cierre constante de conexiones de red.

### RNF2: Seguridad y Protección de Endpoints (Aislamiento Completo)

* Queda estrictamente prohibido que el frontend envíe el ID de base de datos o el correo en texto plano para operar el contador. El backend **únicamente** identificará al usuario desencriptando el JWT local. Si el token está ausente o alterado, devolverá un código de estado `410 Unauthorized`.

### RNF3: Stateless (Sin Estado en Servidor)

* El servidor de Go no almacenará sesiones en memoria ni en archivos locales. Toda la verificación se realiza de manera matemática a través de las firmas de los tokens JWT, facilitando el despliegue en contenedores Docker detrás de Traefik sin necesidad de configurar sesiones pegajosas (sticky sessions).

---

## 4. Especificación de la Base de Datos (Estructura MySQL)

Para cumplir con los requerimientos, la base de datos contendrá dos tablas principales con integridad referencial.

### Tabla: `users`

Almacena la identidad inmutable provista por Google.

```sql
CREATE TABLE users (
    id INT AUTO_INCREMENT PRIMARY KEY,
    google_id VARCHAR(255) NOT NULL UNIQUE, -- Aquí se almacena el campo 'sub'
    email VARCHAR(255) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

```

### Tabla: `counters`

Almacena el estado transaccional del contador asociado a cada usuario.

```sql
CREATE TABLE counters (
    user_id INT PRIMARY KEY,
    current_value INT DEFAULT 0 NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

```

---

## 5. Especificación de la API (Endpoints)

### 1. Autenticación de Google

* **Ruta:** `POST /api/auth/google`
* **Acceso:** Público
* **Cuerpo de la Petición (JSON):**
```json
{
  "token": "eyJhbGciOiJSUzI1NiIsImtpZCI6..."
}

```


* **Respuesta Exitosa (`200 OK`):**
```json
{
  "token": "jwt_local_generado_por_go_para_el_frontend",
  "user": {
    "name": "Andrés Gonzalo Barrera Cortes",
    "email": "andresgbarrerac@gmail.com"
  }
}

```



### 2. Obtener Contador

* **Ruta:** `GET /api/counter`
* **Acceso:** Protegido (Requiere Header: `Authorization: Bearer <JWT_LOCAL>`)
* **Respuesta Exitosa (`200 OK`):**
```json
{
  "current_value": 42
}

```



### 3. Incrementar Contador

* **Ruta:** `POST /api/counter/increment`
* **Acceso:** Protegido (Requiere Header: `Authorization: Bearer <JWT_LOCAL>`)
* **Respuesta Exitosa (`200 OK`):**
```json
{
  "current_value": 43,
  "status": "success"
}

```