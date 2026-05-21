# Counter

Counter is a full-stack web application designed to provide individual user-based counters. The primary objective of this repository is to establish a production-ready architectural foundation, serving as an educational environment for mastering distributed systems, modern web development tech stacks, containerized deployments, secure authentication, and collaborative development workflows.

---

## Technical Stack

* **Frontend:** Astro (TypeScript)
* **Backend:** Go (Golang)
* **Database:** MySQL
* **Reverse Proxy & Edge Router:** Traefik
* **Containerization:** Docker & Docker Compose
* **Operating System & Hosting:** Ubuntu Server via Hostinger VPS
* **Authentication:** Google OAuth 2.0 & JSON Web Tokens (JWT)

---

## Architectural Objectives

### 1. Separation of Concerns (Decoupled Architecture)

The project completely separates the client-side presentation layer from the server-side business logic. The Astro frontend operates as an independent static/server layer that communicates with the Go REST API solely through HTTP requests, enabling independent scalability and maintenance of both services.

### 2. High Performance and Efficiency

By utilizing Go for the backend, the application guarantees low memory consumption, high concurrency handling, and rapid HTTP request execution. MySQL is optimized to ensure low-latency reads and writes for real-time state updates.

### 3. Stateless Security & Multi-Tenant Data Isolation

Authentication is implemented via a secure federated identity flow:

* The frontend authenticates the user through **Google OAuth 2.0** and receives an Identity Token.
* The backend verifies the cryptographic signature of the Google token, extracts the immutable subject identifier (`sub`), and provisions a localized **JSON Web Token (JWT)**.
* Data integrity is prioritized by partitioning counters based on the verified user ID extracted from the local JWT claims, preventing unauthorized cross-tenant data access.

### 4. Infrastructure & Production Edge Routing with Traefik

The environment is fully containerized via **Docker**, ensuring consistency across local development and production. In the Hostinger VPS production environment running **Ubuntu** under the domain `xxxxxx`, **Traefik** acts as the unified edge router, dynamically managing routing, SSL termination, and secure request forwarding to the respective containerized services.

### 5. Collaboration Standards

The project layout, component modularity, and database migrations are structured to emulate a professional enterprise environment, emphasizing clean code, API documentation, and industry-standard Git workflows to optimize team-based development.
