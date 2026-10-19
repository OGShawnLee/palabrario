# Documentación de Infraestructura y Despliegue - Palabrario

Este documento describe la arquitectura de desarrollo y despliegue híbrida entre **Windows (Host)**, **Vagrant (VM Linux)** y **Docker Compose (Contenedores)**, diseñada para garantizar el máximo rendimiento, mínimo consumo de memoria RAM (`< 15 MiB` combinados para todo el backend en Go), aislamiento de dependencias y compilación estática.

---

## 1. Arquitectura de 3 Capas (Monorepo 100% Go)

```mermaid
flowchart TD
    subgraph Host["1. Windows (Host)"]
        Editor["Editor / IDE (VS Code, Go 1.24, sqlc)"]
        Frontend["Cliente Nativo Wails v2 (Go + Vue 3 + Vite)"]
        RsyncWatcher["vagrant rsync-auto (Opcional con VM)"]
    end

    subgraph VM["2. Entorno Linux (Docker Desktop o Vagrant Debian 12)"]
        SyncedCode["Repositorio Monorepo Go (palabrario)"]
        DockerEngine["Docker Engine + Docker Compose"]
    end

    subgraph Containers["3. Docker Compose (Servicios Conteinerizados)"]
        Postgres["postgres: PostgreSQL 17 Alpine (5432:5432)"]
        Core["core-service: Go 1.24 Alpine (8000:8000)"]
        Realtime["realtime-service: Go 1.24 Alpine (8001:8001)"]
        ImageSvc["image-service: Go 1.24 Alpine (8002:8002)"]
    end

    Editor -->|Modifica archivos / sqlc generate| SyncedCode
    RsyncWatcher -.->|Sincroniza en tiempo real| SyncedCode
    Frontend -->|HTTP / WS a localhost:8000, 8001, 8002| DockerEngine
    DockerEngine -->|Port Forwarding| Containers
    Core -->|pgxpool TCP :5432| Postgres
    Realtime -->|Prepare & Batch Flush| Core
```

---

## 2. Servicios y Puertos Expuestos

Todos los puertos están expuestos hacia la máquina local (`localhost`):

| Servicio | Tecnología / Runtime | Puerto | Huella RAM (Docker) | Descripción |
| :--- | :--- | :---: | :---: | :--- |
| **`postgres`** | PostgreSQL 17 Alpine | `5432` | `~30 MiB` | Base de datos relacional principal (`unaccent`, `pg_trgm`, FTS GIN, Triggers `RN06/RN07`). Persistida en volumen `postgres_data`. |
| **`core-service`** | Go 1.24 + Chi v5 + `sqlc` + `pgxpool` | `8000` | `< 5 MiB` | Monolito modular SOA (Auth JWT, Gobernanza, Diccionario FTS, Quórum, Studio, SM-2, Ledger de XP). |
| **`realtime-service`** | Go 1.24 + Chi v5 + WebSockets | `8001` | `< 5 MiB` | Salas de competencia multijugador en tiempo real en memoria RAM (`ArenaHub` + Goroutines) y *Batch Flush* al cierre (Migración a Go en Fase 4). |
| **`image-service`** | Go 1.24 + Chi v5 | `8002` | `< 5 MiB` | Validación binaria, eliminación EXIF y compresión de imágenes a WebP. |

---

## 3. Estrategia de Construcción en Monorepo Unificado Go (`go.mod` en Raíz)

En [`docker-compose.yml`](../docker-compose.yml), los tres servicios de backend (`core-service`, `realtime-service` e `image-service`) sitúan su `build.context` en la raíz (`.`):

```yaml
  core-service:
    build:
      context: .
      dockerfile: ./core-service/Dockerfile
```

* **Justificación:** Permite que el build multi-stage en Go acceda tanto al código propio de cada servicio (`./core-service`, `./realtime-service`, `./image-service`) como a los paquetes compartidos de modelos y autenticación (`./pkg/models`, `./pkg/auth`, `./pkg/response`) y al manifiesto unificado [`go.mod`](../go.mod), generando binarios estáticos compilados (`CGO_ENABLED=0`) sobre `alpine:3.21` de `~27 MB` y `< 5 MiB` de uso de RAM en ejecución.

---

## 4. Flujo de Base de Datos y Generación de Código (`SQL Schema-First` con `sqlc`)

El proyecto sigue una arquitectura **SQL Schema-First** estricta:

1. **Definición Canónica:** Todo el esquema relacional vive en [`SCHEMA_FINAL.sql`](../SCHEMA_FINAL.sql) en la raíz, incluyendo las 23 tablas del dominio, extensiones (`unaccent`, `pg_trgm`), índices GIN de búsqueda FTS, triggers de integridad (`trg_prevent_self_approval`, `trg_auto_approve_sense`) y semillas determinísticas.
2. **Inicialización Automática en Docker:**
   En [`docker-compose.yml`](../docker-compose.yml), `SCHEMA_FINAL.sql` se monta en `/docker-entrypoint-initdb.d/init.sql:ro` para que PostgreSQL 17 lo ejecute automáticamente al crear el volumen.
3. **Consultas Tipadas en Tiempo de Compilación (`sqlc`):**
   Las consultas SQL de los módulos se escriben en `sql/queries/*.sql` y se compilan a código Go fuertemente tipado (`core-service/internal/db/*.sql.go`) según [`sqlc.yaml`](../sqlc.yaml):
   ```powershell
   sqlc generate
   ```

---

## 5. Guía de Puesta en Marcha y Comandos Frecuentes

### Opción A: Ejecución Directa con Docker Compose (Local)
```powershell
# Levantar base de datos y servicios de backend
docker compose up -d --build

# Verificar salud de los 3 servicios
curl http://localhost:8000/health
curl http://localhost:8001/health
curl http://localhost:8002/health
```

### Opción B: Ejecución dentro de Máquina Virtual Aislada (Vagrant)
```powershell
# En terminal 1 (Windows):
vagrant up
vagrant rsync-auto

# En terminal 2 (Dentro de la VM Debian 12):
vagrant ssh
cd /home/vagrant/palabrio
docker compose up --build
```

### Tabla Resumen de Comandos de Desarrollo

| Acción | Comando (en la raíz del repositorio) |
| :--- | :--- |
| **Generar código Go desde SQL (`sqlc`)** | `sqlc generate` |
| **Compilar todos los servicios Go localmente** | `go build ./core-service/... ./image-service/... ./pkg/...` |
| **Ejecutar pruebas de integración (`core-service`)** | `go test -v ./core-service/...` |
| **Reconstruir contenedor tras cambios en Go** | `docker compose up -d --build core-service` |
| **Reiniciar base de datos desde cero (re-ejecutar `SCHEMA_FINAL.sql`)** | `docker compose down -v && docker compose up -d postgres` |

---

## 6. Procedimiento de Cerrado y Apagado Ordenado

Para evitar corrupción de datos en **PostgreSQL** y liberar recursos limpiamente:

1. **Detener los Contenedores Docker:**
   ```powershell
   docker compose down
   ```
2. **Si se utiliza Vagrant VM:**
   - Detener `vagrant rsync-auto` (`Ctrl + C`).
   - Salir de la sesión SSH (`exit`) y apagar la máquina virtual con `vagrant halt`.
