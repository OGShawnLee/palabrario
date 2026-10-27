# 🖥️ Guía de Aprendizaje Arquitectónico: Wails v3 + Bun + Go Workspaces (`learning/WAILS3-STACK.md`)

> **Manual de Integración y Referencia Operativa para Agentes de Palabrario**  
> **Propósito:** Documentar los aprendizajes, trampas resueltas y el flujo de trabajo automatizado para desarrollar y verificar el cliente de escritorio con **Wails v3**, **Bun** y **Go Workspaces**.

---

## 1. Los Tres Grandes Desafíos y sus Soluciones

### 1.1. Go Workspaces (`go.work`) y Bindings Vacíos (`0 Services, 0 Methods`)
* **Problema:**  
  Al ejecutar `wails3 generate bindings`, el generador advertía:
  ```text
  WARNING [warn] -: main module (palabrario) does not contain package palabrario/local-app
  INFO Processed: 1 Package, 0 Services, 0 Methods, 0 Enums, 0 Models, 0 Events
  ```
  Esto provocaba que no se generaran los archivos en `frontend/bindings/` y Vite fallara con `[UNRESOLVED_IMPORT] Could not resolve '../bindings/...'`.
* **Causa:**  
  En un monorepo gobernado por `go.work`, si el subdirectorio de Wails (`local-app`) tiene su propio `go.mod` pero no está registrado en `go.work`, las herramientas de Go asumen el módulo raíz y no encuentran los paquetes del subdirectorio.
* **Solución Obligatoria:**  
  Registrar el módulo en `go.work`:
  ```powershell
  go work use ./local-app
  ```
  Una vez registrado, Wails3 indexa instantáneamente todos los servicios, métodos y eventos.

---

### 1.2. Versión de TypeScript vs `tsconfig.json` (`moduleResolution: "bundler"`)
* **Problema:**  
  El chequeo de tipos (`tsc --noEmit`) fallaba con errores TS6046, TS5070 y TS5023 (`verbatimModuleSyntax` desconocido y `moduleResolution: "bundler"` inválido).
* **Causa:**  
  El template inicial traía `"typescript": "^4.9.3"`, pero `tsconfig.json` usaba opciones modernas de TypeScript 5.0+ (`bundler`, `verbatimModuleSyntax`).
* **Solución Obligatoria:**  
  Actualizar `typescript` a `>= 5.7` en `frontend/package.json` y reinstalar con Bun (`bun install`).

---

### 1.3. `go build ./...` vs `go build .` en Wails v3
* **Problema:**  
  Al ejecutar `go build ./...` dentro del proyecto Wails3, Go intentaba compilar `build/ios` y fallaba con:
  ```text
  # changeme/build/ios
  runtime.main_main·f: function main is undeclared in the main package
  ```
* **Causa:**  
  La carpeta `build/` traía plantillas de empaquetado para plataformas no usadas en Wails (iOS, Android, Docker cross-compilation). En Palabrario, el soporte móvil se resuelve mediante **Capacitor PWA**, por lo que esos archivos son código muerto.
* **Solución Obligatoria:**  
  1. Purgar las carpetas muertas `build/ios`, `build/android` y `build/docker`.
  2. Para validar la compilación Go del cliente, usar siempre:
  ```powershell
  go build .
  ```

---

## 2. Diferencias Críticas entre Wails v2 y Wails v3

| Aspecto | Wails v2 | Wails v3 |
| :--- | :--- | :--- |
| **Motor de Tareas** | Scripts libres en `wails.json` | **`Taskfile.yml`** nativo (`wails3 task ...`) |
| **Package Manager** | Principalmente `npm` | **Soporte nativo para `bun`** (`PACKAGE_MANAGER: "bun"`) |
| **Runtime JS/TS** | `@wailsjs/runtime` / `window.runtime` | **`@wailsio/runtime`** (`v3.0.0`) |
| **Directorio de Bindings** | `wailsjs/go/<pkg>/<struct>.js` | **`frontend/bindings/<module>/<struct>.ts`** |
| **Eventos** | `EventsOn` / `EventsEmit` en runtime | `application.RegisterEvent[T]("name")` en Go + `Events.On` tipado |
| **Recursos Windows** | Ocultos en memoria | Genera **`wails_windows_amd64.syso`** con manifiesto, íconos y metadata |

---

## 3. Workflow Encadenado y Automatizado (Desde la Raíz)

Hemos configurado un **`Taskfile.yml` raíz** y el script **[`.\scripts\client.ps1`](../scripts/client.ps1)** para operar todo el monorepo sin fricción:

### Matriz de Comandos Encadenados

| Tarea | Opción A: Con `wails3 task` | Opción B: Con PowerShell | Qué ejecuta por debajo en cascada |
| :--- | :--- | :--- | :--- |
| **Validar cliente (Check rápido)** | `wails3 task client:check` | `.\scripts\client.ps1 check` | 1. Regenera bindings (`wails3 generate bindings`)<br>2. Chequea tipos (`bun run --cwd frontend type-check`)<br>3. Compila frontend (`bun run --cwd frontend build`)<br>4. Compila Go (`go build .`) |
| **Compilar binario nativo** | `wails3 task client:build` | `.\scripts\client.ps1 build` | Genera íconos, bindings, build de Vite, `.syso` de Windows y compila `local-app/bin/palabrario.exe` (11.2 MB). |
| **Compilar y abrir app en vivo** | `wails3 task client:start` | `.\scripts\client.ps1 start` | Ejecuta `client:build` + abre la ventana WebView2 nativa. |
| **Modo desarrollo (Hot-Reload)** | `wails3 task client:dev` | `.\scripts\client.ps1 dev` | Inicia Vite dev server + Wails dev watcher en el puerto 9245. |
| **Regenerar solo bindings TS** | `wails3 task client:bindings` | `.\scripts\client.ps1 bindings` | `wails3 generate bindings -clean=true -ts` |
| **Correr suite de backend** | `wails3 task test` | `.\scripts\test.ps1` | Postgres 17 RAM efímero + tests de servicios 80/20. |

---

## 4. Regla de Oro para Agentes Constructores de Fase 5

> [!TIP]
> **Criterio de Validación Automática para Sub-fases Frontend:**  
> Ninguna tarea de Fase 5 se da por terminada hasta que ejecutes desde la raíz:
> ```powershell
> wails3 task client:check
> ```
> *(o `.\scripts\client.ps1 check`)* y retorne **cero errores** con código de salida `0`.
