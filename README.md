# 📚 Palabrario

> **Plataforma Gamificada de Aprendizaje y Preservación del Léxico Dialectal Panhispánico.**  
> *Proyecto desarrollado para la cátedra de Despliegue de Software (Software Deployment).*

---

## 🧭 Mapa de Documentación Arquitectónica (Single Source of Truth)

| Documento | Propósito |
| :--- | :--- |
| **[`.context/CONTEXT.md`](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/.context/CONTEXT.md)** | **Punto de Entrada Maestro:** Enrutador de contexto, reglas anti-alucinación (*Contract-First*) y resumen ejecutivo. |
| **[`.context/PLAN.md`](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/.context/PLAN.md)** | Diagnóstico del estado actual de implementación y hoja de ruta paso a paso por fases. |
| **[`DOMAIN.md`](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/DOMAIN.md)** | Modelo de Dominio (DDD), contextos delimitados, invariantes, máquina de estados editorial y especificación de las 23 entidades. |
| **[`INTERFACES.md`](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/INTERFACES.md)** | Arquitectura Orientada a Servicios (SOA) en Go 1.24, DTOs canónicos, consultas `sqlc`, eventos WebSocket y mapa completo de endpoints HTTP. |
| **[`SCHEMA_FINAL.sql`](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/SCHEMA_FINAL.sql)** | Esquema DDL canónico en PostgreSQL 17 (`unaccent`, `pg_trgm`, `tsvector` GIN, triggers `RN06`/`RN07` y datos semilla). |
| **[`INFRASTRUCTURE.md`](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/INFRASTRUCTURE.md)** | Guía de contenedores Docker Compose, monorepo 100% Go + `sqlc` y máquina virtual Vagrant. |
| **[`UI_SPECIFICATION.md`](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/UI_SPECIFICATION.md)** | Mapa de enrutamiento, espacios de navegación y componentes globales del cliente Wails v2 + Vue 3 + UnoCSS (para Fase 5). |
| **[`.context/academic/`](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/.context/academic)** | Entregables académicos formales autocontenidos (`1-PROJECT-PROPOSAL.md`, `2-PROJECT-SPECIFICATION.md`, `3-FINAL-USE-CASES.md`). |

---

## 🎯 1. Visión y Propósito del Proyecto

El español cuenta con más de 500 millones de hablantes nativos distribuidos en más de 20 países. A pesar de compartir una misma base lingüística, la fragmentación léxica y los choques semánticos (*false friends* dialectales) son un desafío cotidiano:

- Una **guagua** es un autobús en Cuba y Canarias, pero un bebé de brazos en Chile y la región andina.
- Un **popote** (México) es un **pitillo** (Colombia/Venezuela), una **bombilla** (Chile) o un **sorbete/pajita** (Argentina/España).
- Expresiones de afecto profundo como **apapachar** (*acariciar con el alma*) carecen de traducción directa y son patrimonio dialectal vivo.

**Palabrario** resuelve esta brecha mediante un ecosistema tecnológico moderno, combinando:
1. **Diccionario dialectal inteligente:** Búsqueda fonética y por similitud de trigramas ajustada al dialecto ancla del usuario.
2. **Motor de repetición espaciada (SM-2 Dialectal):** Algoritmo de retención de memoria a largo plazo optimizado por intervalo y factor de facilidad (*easiness factor*).
3. **Arena multijugador en tiempo real:** Duelos sincronizados para competir identificando variantes regionales y expresiones idiomáticas.
4. **Taller de contribución y revisión por pares:** Flujo de moderación lingüística con validación de integridad (anti-autoaprobación y quórum tripartito).

---

## 🏛️ 2. Arquitectura del Sistema y Decisiones de Ingeniería (Backend 100% Go 1.24)

Palabrario está diseñado bajo una arquitectura de **Monolito Modular Orientado a Servicios (SOA) con Servicios Satélite Especializados en un Monorepo 100% Go**, priorizando el consumo mínimo de recursos (`< 15 MiB` de RAM combinada en los 3 servicios Go), la rapidez de despliegue y la integridad de datos:

```mermaid
graph TD
    subgraph Cliente Nativo
        Desktop["🖥️ Palabrario Desktop (Wails v2 + Go + Vue 3)"]
    end

    subgraph Backend Services (Monorepo Go 1.24)
        Core["⚡ core-service (Go 1.24 + Chi + sqlc)<br/>Monolito Modular SOA :8000"]
        RT["📡 realtime-service (Go 1.24 + Chi + WebSockets)<br/>ArenaHub en RAM :8001"]
        Img["🖼️ image-service (Go 1.24 + Chi)<br/>Procesamiento WebP :8002"]
    end

    subgraph Capa de Persistencia
        PG[("🐘 PostgreSQL 17 Alpine<br/>pg_trgm · unaccent · Triggers RN06/RN07<br/>:5432")]
    end

    Desktop -->|"REST / JSON (JWT Auth)"| Core
    Desktop <-->|"WebSockets (Duelos en vivo)"| RT
    Desktop -->|"Multipart / Upload"| Img
    Core -->|"pgxpool (Consultas tipadas en compile-time)"| PG
    RT -->|"Prepare & Batch Flush al cierre"| Core
```

### Justificación Tecnológica y Trade-offs de Despliegue

| Componente | Elección Tecnológica | Alternativa Descartada | Justificación Técnica de Despliegue |
| :--- | :--- | :--- | :--- |
| **Backend Unificado** | **Golang (Go 1.24)** | Node.js / Bun / Python | Binarios estáticos compilados sin dependencias externas en runtime. Consumo de memoria reducido a **< 5 MiB** por servicio en Docker y compartición nativa de `pkg/models` y `pkg/auth` en un solo `go.mod`. |
| **Acceso a Datos** | **`sqlc` + `pgxpool`** | ORMs pesados (GORM, Prisma) | Consultas SQL 100% tipadas en tiempo de compilación. Cero costo de reflexión en runtime, rendimiento nativo de base de datos y prevención absoluta de SQL Injections. |
| **Gestión de UUID** | **`github.com/google/uuid`** | Funciones caseras / strings | Estandarización RFC 4122, soporte nativo como `driver.Valuer`/`sql.Scanner` con `pgx` y mapeo transparente en `sqlc`. |
| **Cliente Desktop** | **Wails v2 (Go + Vue 3)** | Electron | Empaqueta un binario nativo que reutiliza el webview del sistema operativo. Huella de RAM de **~30 MiB** (frente a los >250 MiB del runtime Chromium de Electron). |
| **Persistencia** | **PostgreSQL 17** | NoSQL / SQLite | Índices trigramáticos (`pg_trgm`), soporte de diacríticos (`unaccent`), funciones almacenadas determinísticas y triggers relacionales estrictos (RN06 quórum 3 y RN07 anti-autoaprobación). |

---

## ⚙️ 3. Topología de Servicios y Puertos

| Contenedor / Servicio | Lenguaje / Runtime | Puerto Interno | Puerto Host | Descripción |
| :--- | :--- | :---: | :---: | :--- |
| **`postgres`** | PostgreSQL 17 Alpine | `5432` | `5432` | Base de datos relacional (23 tablas), FTS y triggers |
| **`core-service`** | Go 1.24 (Alpine) | `8000` | `8000` | Monolito Modular SOA: Auth, Gobernanza, Diccionario, Quórum, Studio, SM-2 |
| **`realtime-service`**| Go 1.24 (Alpine) | `8001` | `8001` | Sincronización WebSockets en RAM para la Arena multijugador |
| **`image-service`** | Go 1.24 (Alpine) | `8002` | `8002` | Procesamiento, compresión y validación multimedia |
| **`local-app`** | Wails v2 (Go + Vue 3) | N/A | Nativo | Interfaz gráfica de usuario para escritorio |

---

## 🚀 4. Guía de Puesta en Marcha (Getting Started)

### Requisitos Previos
- **Docker** (v24.0+) y **Docker Compose** (v2.20+)
- **Go** (v1.24+) *(para desarrollo local y generación con sqlc)*
- **Vagrant** & **VirtualBox** *(opcional, si se levanta la infraestructura en máquina virtual aislada)*

### 1. Configuración del Entorno (Principio *Fail-Fast*)
Palabrario implementa el principio **Fail-Fast** de la metodología *12-Factor App*. Los servicios **se niegan a iniciar** si variables críticas (`DATABASE_URL`, `JWT_SECRET`) faltan o son inseguras.

Copia la plantilla de configuración:
```bash
cp .env.example .env
```

> [!IMPORTANT]
> Revisa el archivo [`.env.example`](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/.env.example) para entender el propósito y restricciones de cada variable. En particular, `JWT_SECRET` debe tener una longitud mínima de 32 caracteres.

### 2. Despliegue con Docker Compose
Para levantar la base de datos y el servicio central:
```bash
docker compose up -d --build postgres core-service
```

### 3. Verificación de Salud (*Health Check*)
Comprueba el estado de la API y su conectividad con PostgreSQL:
```bash
curl http://localhost:8000/health
```
Respuesta esperada:
```json
{
  "status": "ok",
  "database": "connected",
  "engine": "PostgreSQL 17 (pgxpool)",
  "service": "palabrario-core-service"
}
```

---

## 👥 5. Cuentas y Datos Semilla Determinísticos

La base de datos se inicializa automáticamente con datos semilla representativos de diversas variantes panhispánicas (`MX`, `ES`, `AR`, `CO`, `CL`, `CU`, `CAM`, `PAN`):

| Rol | Correo Electrónico | Contraseña | Dialecto Ancla | XP Inicial | Nivel |
| :--- | :--- | :--- | :---: | :---: | :---: |
| **ADMIN** | `admin@palabrario.org` | `password123` | MX | 15,000 | 10 |
| **EDITOR** | `silas@palabrario.org` | `password123` | MX | 8,500 | 7 |
| **LEARNER** | `gael@palabrario.org` | `password123` | MX | 4,850 | 4 |

### Ejemplos de Palabras Precargadas
- **`popote`** (Sustantivo, MX): Incluye comparativa dialectal (*pitillo*, *bombilla*, *sorbete*).
- **`apapachar`** (Verbo, MX): Con etimología náhuatl y definición afectiva.
- **`guagua`** (Sustantivo, Poligenético): Ejemplos de uso simultáneos para transporte (Caribe) y bebé (Andes).

---

## 🛠️ 6. Flujo de Desarrollo

### Generación de Código Tipado con `sqlc`
Cualquier modificación o adición de consultas SQL en `sql/queries/*.sql` se compila directamente a estructuras Go fuertemente tipadas:
```bash
sqlc generate
```

### Estrategia de Testing y Aislamiento de Entornos
> [!NOTE]
> Siguiendo las mejores prácticas de la industria, las pruebas automatizadas de integración **nunca se ejecutan contra bases de datos productivas**. Utilizan la variable `TEST_DATABASE_URL` y limpian sus registros automáticamente mediante hooks `t.Cleanup()`.

Para ejecutar las pruebas en desarrollo local:
```bash
go test -v ./core-service/...
```

---

## 🗺️ 7. Estado del Despliegue (Ver detalle en [`.context/PLAN.md`](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/.context/PLAN.md))

- [x] **Fase 1: Modelo Relacional Completo y Reglas de Integridad** (`SCHEMA_FINAL.sql` + `DOMAIN.md`)
- [x] **Fase 2: Espacio de Trabajo y Contratos de Dominio** (`go.mod` + `pkg/models/` + `INTERFACES.md`)
- [x] **Fase 3.1 & 3.2: Arquitectura Base de Servicios y Autenticación JWT**
- [x] **Fase 3.3: Módulo de Diccionario, Autoría Léxica y Curaduría por Quórum (`CU05–CU08`)**
- [x] **Fase 3.4: Studio Personal, Motor de Repetición Espaciada (SM-2) y Ledger de XP (`CU10–CU13`)**
- [ ] **Fase 3.5: Gobernanza Institucional, Perfiles y Moderación de Reportes (`CU01–CU04, CU09`)**
- [ ] **Fase 4: Migración de `realtime-service` a Go 1.24 + Arena Multijugador (WebSockets + Batch Flush) (`CU14–CU16`)**
- [ ] **Fase 5: Integración End-to-End Desktop (Wails v2) y Pipeline CI/CD**
- [ ] **Fase 6 (Incremento a Largo Plazo): Sistema de Eventos Culturales de Temporada**

---

## 📄 Licencia

Proyecto académico desarrollado bajo fines educativos para la materia de **Despliegue de Software**.
