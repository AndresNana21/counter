# Counter

Counter es una aplicación web full-stack diseñada para proporcionar contadores individuales basados en cada usuario. El objetivo principal de este repositorio es establecer una base arquitectónica lista para producción, sirviendo como un entorno educativo para dominar sistemas distribuidos, tecnologías modernas de desarrollo web, despliegues contenedorizados, autenticación segura y flujos de trabajo colaborativos.

---

## Stack Tecnológico

* **Frontend:** Astro (TypeScript)
* **Backend:** Go (Golang)
* **Base de Datos:** MySQL
* **Proxy Inverso y Enrutador de Borde:** Traefik
* **Contenedores:** Docker & Docker Compose
* **Sistema Operativo y Hosting:** Servidor Ubuntu vía VPS de Hostinger
* **Autenticación:** Google OAuth 2.0 & JSON Web Tokens (JWT)

---

## Objetivos Arquitectónicos

### 1. Separación de Conceptos (Arquitectura Desacoplada)

El proyecto separa completamente la capa de presentación del lado del cliente de la lógica de negocio del lado del servidor. El frontend en Astro opera como una capa estática/servidor independiente que se comunica con la API REST de Go únicamente a través de peticiones HTTP, lo que permite la escalabilidad y el mantenimiento independiente de ambos servicios.

### 2. Alto Rendimiento y Eficiencia

Al utilizar Go en el backend, la aplicación garantiza un bajo consumo de memoria, un alto manejo de concurrencia y una rápida ejecución de las peticiones HTTP. MySQL está optimizado para asegurar lecturas y escrituras de baja latencia para las actualizaciones de estado en tiempo real.

### 3. Seguridad Stateless e Aislamiento de Datos Multi-Inquilino

La autenticación se implementa a través de un flujo seguro de identidad federada:

* El frontend autentica al usuario mediante **Google OAuth 2.0** y recibe un Token de Identidad.
* El backend verifica la firma criptográfica del token de Google, extrae el identificador inmutable del sujeto (`sub`) y proporciona un **JSON Web Token (JWT)** local.
* Se prioriza la integridad de los datos particionando los contadores en función del ID de usuario verificado extraído de los claims del JWT local, evitando el acceso no autorizado a datos entre diferentes usuarios.

### 4. Infraestructura y Enrutamiento de Borde en Producción con Traefik

El entorno está completamente contenedorizado mediante **Docker**, lo que garantiza la consistencia entre el desarrollo local y la producción. En el entorno de producción del VPS de Hostinger que ejecuta **Ubuntu** bajo el dominio `xxxxxx`, **Traefik** actúa como el enrutador de borde unificado, gestionando dinámicamente el enrutamiento, la terminación SSL y el reenvío seguro de peticiones a los servicios contenedorizados respectivos.

### 5. Estándares de Colaboración

El diseño del proyecto, la modularidad de los componentes y las migraciones de la base de datos están estructurados para emular un entorno profesional de nivel empresarial, haciendo hincapié en el código limpio, la documentación de la API y los flujos de trabajo de Git estándar de la industria para optimizar el desarrollo en equipo.
