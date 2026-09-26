# Documentación de Infraestructura - Palabrario

Este documento describe la arquitectura de desarrollo híbrida entre **Windows (Host)**, **Vagrant (VM Linux)** y **Docker (Contenedores)**, diseñada para garantizar el máximo rendimiento, aislamiento de dependencias nativas y recarga en caliente (Hot-Reload) para el frontend y backend.

---

## 1. Arquitectura de 3 Capas

```mermaid
flowchart TD
    subgraph Host["1. Windows (Host)"]
        Editor["Editor / IDE (VS Code, etc.)"]
        Frontend["Frontend Dev Server (Next.js / Vite / React)"]
        RsyncWatcher["vagrant rsync-auto"]
    end

    subgraph VM["2. Vagrant VM (Debian 12 Bookworm)"]
        SyncedCode["/home/vagrant/palabrio (Excluye node_modules)"]
        DockerEngine["Docker Engine + Docker Compose"]
    end

    subgraph Containers["3. Docker Compose (Microservicios)"]
        Postgres["postgres (5432:5432)"]
        Core["core-service: Bun (8000:8000)"]
        Realtime["realtime-service: Bun (8001:8001)"]
        ImageSvc["image-service: Python (8002:8002)"]
    end

    Editor -->|Modifica archivos| SyncedCode
    RsyncWatcher -->|Sincroniza en tiempo real| SyncedCode
    Frontend -->|HTTP / WS a localhost:8000, 8001, 8002| VM
    VM -->|Port Forwarding| Containers
    SyncedCode -->|Volumen montado /app| Core
    SyncedCode -->|Volumen montado /app| Realtime
    SyncedCode -->|Volumen montado /app| ImageSvc
```

---

## 2. Servicios y Puertos Expuestos

Todos los puertos están redirigidos desde el contenedor hacia la máquina virtual y desde la máquina virtual hacia tu máquina local Windows (`localhost`):

| Servicio | Tecnología | Puerto | Descripción |
| :--- | :--- | :--- | :--- |
| **`postgres`** | PostgreSQL 17 Alpine | `5432` | Base de datos relacional principal. Persistida en volumen Docker `postgres_data`. |
| **`core-service`** | Bun 1.x (TypeScript) | `8000` | API principal / Monolito (lógica de negocio, vocabulario, diccionarios). |
| **`realtime-service`** | Bun 1.x (WebSockets) | `8001` | Salas de competencia multijugador en tiempo real y eventos en vivo. |
| **`image-service`** | Python 3.11 + FastAPI + Pillow | `8002` | Validación, procesamiento y moderación de imágenes (Docs en `/docs`). |

---

## 3. Estrategia de Aislamiento de Dependencias (`node_modules`)

Uno de los problemas más comunes al desarrollar en Windows y ejecutar en Linux son las incompatibilidades de binarios (`.node` compilados para Windows vs Linux) y la lentitud del sistema de archivos compartido.

Para evitar esto se implementó un **doble filtro de aislamiento**:

### Capa A: Windows ➔ Vagrant
En [Vagrantfile](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/Vagrantfile):
```ruby
config.vm.synced_folder ".", "/home/vagrant/palabrio",
  type: "rsync",
  rsync__exclude: ["**/build/", "**/node_modules/", ".git/"]
```
* **Efecto:** Cualquier `node_modules` que tengas instalado en Windows jamás entra a la VM.

### Capa B: Vagrant ➔ Contenedores Docker
En [docker-compose.yml](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/docker-compose.yml):
```yaml
volumes:
  - ./core-service:/app
  - /app/node_modules
```
* **Efecto:** El volumen anónimo `/app/node_modules` le dice a Docker que **preserve** los módulos instalados dentro de la imagen Linux durante el build (`RUN bun install`), impidiendo que el montaje de `./core-service` los sobreescriba.

---

## 4. Guía de Uso Diario

### Paso 1: Levantar la sincronización de archivos (Windows)
Abre una terminal en la raíz del proyecto en Windows:
```powershell
vagrant up
vagrant rsync-auto
```
> Deja esta terminal abierta en segundo plano. Detectará cada archivo modificado en Windows y lo sincronizará a la VM en milisegundos.

### Paso 2: Iniciar los contenedores (Dentro de la VM)
En una segunda terminal:
```powershell
vagrant ssh
```
Una vez dentro de Linux:
```bash
cd /home/vagrant/palabrio
docker compose up --build
```
> Gracias a los flags `--watch` (en Bun) y `--reload` (en Uvicorn), cada cambio guardado en Windows se reflejará instantáneamente sin necesidad de reiniciar contenedores.

### Paso 3: Conectar el Frontend desde Windows
Tu aplicación frontend (Next.js, Vite, React Native, etc.) corriendo en Windows se comunica directamente a `localhost`:
* API Core: `http://localhost:8000`
* Realtime: `ws://localhost:8001` o `http://localhost:8001`
* Moderación de imágenes: `http://localhost:8002/moderate`
* Base de Datos (DBeaver / Prisma / Drizzle): `postgresql://postgres:postgres@localhost:5432/palabrario`

---

## 5. Gestión e Instalación de Dependencias (Sincronización)

Debido al aislamiento de dependencias entre Windows y Linux, cuando agregas una nueva librería (por ejemplo **Elysia** o un nuevo paquete de Python), el flujo consta de dos sencillos pasos:

### Para Servicios Bun (`core-service` / `realtime-service`):

1. **Instalar en Windows (para Intellisense y Tipos en tu IDE):**
   ```powershell
   cd core-service
   bun add elysia
   ```
   * Actualiza tu `package.json` y `bun.lock` en Windows.
   * Proporciona autocompletado y tipos de TypeScript en tu editor.
   * `vagrant rsync-auto` sincronizará inmediatamente estos dos archivos a la VM.

2. **Instalar en el Contenedor Docker (para Ejecución en Linux):**
   * **Opción Rápida (Recomendada, en caliente):**
     ```bash
     # Dentro de la VM Vagrant (o vía docker compose exec):
     docker compose exec core-service bun add elysia
     ```
     *Toma un par de segundos y no detiene los contenedores.*
   * **Opción Reconstrucción Limpia:**
     ```bash
     docker compose up -d --build core-service
     ```

### Para el Servicio Python (`image-service`):

1. **Agregar la librería en Windows:**
   Edita [image-service/requirements.txt](file:///C:/Users/OGSha/OneDrive/Escritorio/Developer/Web%20Development/Palabrario/image-service/requirements.txt) agregando el paquete (ejemplo: `httpx>=0.27.0`).
2. **Aplicar en el Contenedor Docker:**
   ```bash
   docker compose exec image-service pip install -r requirements.txt
   ```
   *(O si prefieres build limpio: `docker compose up -d --build image-service`)*.

### Tabla Resumen de Comandos

| Acción | Comando en Windows (Terminal Host) | Comando en Docker (Dentro de Vagrant) |
| :--- | :--- | :--- |
| **Añadir dependencia Bun** | `cd core-service; bun add <pkg>` | `docker compose exec core-service bun add <pkg>` |
| **Remover dependencia Bun** | `cd core-service; bun remove <pkg>` | `docker compose exec core-service bun remove <pkg>` |
| **Añadir dependencia Python**| Editar `requirements.txt` | `docker compose exec image-service pip install -r requirements.txt` |
