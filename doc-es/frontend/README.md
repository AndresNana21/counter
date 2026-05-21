# Especificación de Requerimientos y Arquitectura del Frontend (Astro)

Este documento detalla la especificación técnica, los requerimientos funcionales y no funcionales, y la estructura de componentes para el frontend del proyecto **Counter**, desarrollado con Astro y TypeScript. El frontend está diseñado como una interfaz ligera, reactiva y optimizada para consumir la API stateless de Go.

---

## 1. Arquitectura del Frontend y Estructura del Proyecto

Astro maneja un sistema de enrutamiento basado en archivos dentro del directorio `src/pages/`. Para mantener el proyecto escalable, profesional y facilitar el trabajo en equipo, se separan los componentes visuales de la lógica de servicios que interactúa con el backend.

### Estructura de Directorios Propuesta

```text
frontend/
├── public/
│   └── favicon.svg          # Recursos estáticos globales
├── src/
│   ├── components/
│   │   ├── CounterButton.astro # Componente visual del botón del contador
│   │   ├── GoogleButton.astro  # Componente del botón oficial de Google
│   │   └── Navbar.astro        # Barra de navegación superior
│   ├── layouts/
│   │   └── Layout.astro        # Plantilla base HTML (Head, Scripts de Google)
│   ├── pages/
│   │   ├── index.astro         # Página de inicio / Landing de presentación
│   │   ├── login.astro         # Página de autenticación
│   │   └── dashboard.astro     # Panel privado del contador por usuario
│   ├── services/
│   │   └── api.ts              # Cliente de peticiones HTTP (Fetch hacia Go)
│   └── utils/
│       └── auth.ts             # Gestión del almacenamiento local del JWT
├── astro.config.mjs
├── package.json
└── tsconfig.json

```

---

## 2. Requerimientos Funcionales (RF)

El frontend debe cumplir estrictamente con los siguientes flujos de interacción con el usuario y el backend:

### RF1: Integración del SDK Oficial de Google Authentication

* **Descripción:** La aplicación debe cargar de manera asíncrona la librería oficial `https://accounts.google.com/gsi/client` en las vistas correspondientes.
* **Flujo:** 1. Renderizar el botón nativo de Google utilizando la API de JavaScript (`google.accounts.id.renderButton`).
2. Capturar el `CredentialResponse` generado de forma segura por los servidores de Google tras la interacción del usuario.

### RF2: Intercepción y Envío del ID Token de Google

* **Descripción:** El frontend no debe intentar procesar ni almacenar los datos del perfil de Google de manera local para iniciar sesión.
* **Flujo:**
1. El script del cliente captura el token JWT nativo de Google.
2. Invoca inmediatamente una función del servicio (`services/api.ts`) para enviar dicho token mediante una petición `POST /api/auth/google` hacia el backend en Go.



### RF3: Gestión del Ciclo de Vida de la Sesión Local (JWT)

* **Descripción:** El frontend debe administrar la sesión propia emitida por la API de Go.
* **Flujo:**
1. Tras recibir la respuesta exitosa del backend con el JWT local, este se almacena de forma segura en el `localStorage` o en una `Cookie` de sesión.
2. El estado de la interfaz cambia para redirigir al usuario hacia la ruta protegida (`/dashboard`).
3. Proveer una función de "Cierre de Sesión" (Logout) que elimine el token local y redirija al usuario a la pantalla de inicio.



### RF4: Consumo e Incremento Transaccional del Contador

* **Descripción:** La página del panel de usuario (`/dashboard`) debe interactuar en tiempo real con el estado del contador almacenado en la base de datos.
* **Flujo:**
1. Al cargar la página `/dashboard`, se realiza una petición `GET /api/counter` incluyendo el JWT local en las cabeceras.
2. Al presionar el botón del contador, se despacha una petición `POST /api/counter/increment`.
3. La interfaz actualiza el número en pantalla basándose estrictamente en el valor entero retornado por la API de Go, garantizando la consistencia de los datos.



---

## 3. Requerimientos No Funcionales (RNF)

### RNF1: Seguridad en las Cabeceras de Petición

* El frontend tiene prohibido enviar el ID del usuario o datos sensibles en el cuerpo de las peticiones para operar el contador. Toda petición hacia rutas protegidas del backend debe adjuntar el token local en la cabecera estándar de autorización:
```text
Authorization: Bearer <JWT_LOCAL>

```



### RNF2: Optimización del Lado del Cliente (Scripts de Astro)

* Para evitar penalizaciones en el rendimiento y mantener la velocidad nativa de Astro, el SDK de Google y los scripts de inicialización de autenticación se ejecutarán del lado del cliente utilizando etiquetas `<script>` optimizadas y tipadas con TypeScript, evitando bloquear el renderizado inicial del HTML estructurado en el servidor.

### RNF3: Control de Acceso en Rutas (Protección de Vistas)

* La página `/dashboard` debe validar la existencia del JWT local antes de renderizar la información del usuario. Si un cliente no autenticado intenta acceder directamente a la URL del panel, el script del lado del cliente debe interceptar la acción y redirigir inmediatamente a la página `/login`.

---

## 4. Especificación de las Vistas del Frontend

### 1. Vista de Login (`/login`)

* **Elementos:** Título del proyecto, descripción breve y contenedor único para el botón dinámico de Google.
* **Comportamiento:** Bloquea la interfaz con un estado de carga mientras se valida el token con el backend tras el clic del usuario.

### 2. Vista del Panel Principal (`/dashboard`)

* **Elementos:** Barra de navegación con el nombre del usuario y botón de cerrar sesión, contenedor numérico que muestra el valor actual del contador, y un botón de acción principal para incrementar.
* **Comportamiento:** Realiza las llamadas asíncronas hacia la API de Go y maneja visualmente los posibles errores de red o tokens expirados (redirigiendo al login si el backend devuelve un código `401 Unauthorized`).