# 🖥️ Palabrario — Capa de Escritorio Nativa (`local-app`) & Shell Frontend

Esta carpeta contiene la aplicación de escritorio nativa construida con **Wails v3** (Go 1.24) y el shell frontend multiplataforma en **Vue 3 + TypeScript + UnoCSS**, integrado en el monorepo de Bun Workspaces junto al paquete de diseño compartido [`@palabrario/ui`](../packages/ui).

---

## 🏛️ 1. Arquitectura del Monorepo Frontend

El frontend de Palabrario opera bajo una arquitectura desacoplada y fuertemente tipada:

```text
Palabrario/
├── package.json                      # Workspace Root Bun: ["packages/*", "local-app/frontend"]
├── tsconfig.base.json                # Opciones estrictas TypeScript y alias @bindings/*
│
├── packages/
│   └── ui/                           # 📦 @palabrario/ui (Sistema de Diseño SSR-Safe)
│       ├── uno.config.ts             # SSOT Visual (PALETTE.md) + escaneo filesystem
│       ├── src/assets/main.css       # Tipografía Satoshi + .tabular-nums
│       ├── src/types/models.ts       # Re-export de tipos desde @bindings/palabrario/pkg/models
│       └── src/components/           # Componentes puros (ui/, dictionary/, studio/, arena/, editorial/)
│
└── local-app/                        # 🖥️ Shell Nativo de Escritorio y Móvil SPA
    ├── app.go                        # Puente Go BFF (Custodia JWT en RAM con sync.RWMutex)
    ├── main.go                       # Punto de entrada Wails v3 (Ventana 1280x820 #FFFDF5)
    ├── Taskfile.yml                  # Tareas de compilación y verificación (wails3 task check)
    └── frontend/                     # Shell Vue 3 SPA (Capacitor + Wails)
        ├── package.json              # Depende de "@palabrario/ui": "workspace:*"
        ├── uno.config.ts             # Re-exporta @palabrario/ui/uno.config
        ├── vite.config.ts            # Vue + UnoCSS + plugin Wails v3
        ├── capacitor.config.ts       # Configuración Capacitor (com.palabrario.app)
        └── src/                      # Vistas, router, api adapters y composables
```

---

## 🔄 2. Flujo Dual de Transporte (`IPalabrarioAPI`)

Para garantizar que el 100% de la lógica de negocio y las vistas Vue funcionen de forma idéntica en escritorio nativo y en móvil/web sin bifurcaciones condicionales en los componentes:

1. **Wails Desktop (`WailsDesktopAdapter`):**
   - **Sin llamadas directas a red en JS:** El WebView delega en métodos Go expuestos por `local-app/app.go` vía IPC binario.
   - **Bóveda de Sesión Segura:** El token JWT vive exclusivamente en la memoria RAM del proceso Go (`sync.RWMutex`), siendo inmune a ataques de XSS en el navegador/webview.
   - **Backend-For-Frontend (BFF):** El proceso Go realiza las peticiones HTTP (`:8000`, `:8002`) y WebSocket (`:8001`) hacia los microservicios, inyectando de forma transparente las credenciales.

2. **Capacitor Móvil / Web SPA (`CapacitorMobileAdapter`):**
   - **Implementación Idéntica:** Satisface exactamente el contrato `IPalabrarioAPI = typeof WailsAppBridge`.
   - **Almacenamiento Aislado:** El token JWT se custodia mediante `@capacitor/preferences` (Keychain en iOS, Keystore en Android).
   - **Transporte HTTP/WS Estándar:** Ejecuta `fetch()` tipado hacia `:8000` y abre WebSockets nativos del navegador hacia `:8001/ws/arena/{roomCode}?token={jwt}`.

---

## 🛠️ 3. Comandos de Verificación y Desarrollo

Todos los comandos se orquestan a través de `wails3 task` o el script de conveniencia PowerShell:

### Verificación Total (Gate de Sub-Fase)
Ejecuta la cadena completa: generación de bindings Go $\to$ TypeScript $\to$ verificación de tipos en `@palabrario/ui` $\to$ verificación de tipos en `frontend` $\to$ bundle de Vite $\to$ compilación Go:

```powershell
# Usando el script raíz (Recomendado)
.\scripts\client.ps1 check

# O directamente con el CLI de Wails v3 dentro de local-app/
cd local-app
wails3 task check
```

### Desarrollo en Caliente (Hot-Reload)
Inicia la ventana de escritorio conectada al servidor Vite en modo desarrollo:

```powershell
.\scripts\client.ps1 dev
# o: cd local-app && wails3 task dev
```

### Compilación de Producción
Genera el ejecutable nativo optimizado en `local-app/bin/palabrario.exe`:

```powershell
.\scripts\client.ps1 build
```

### Regeneración de Bindings TypeScript
Si modificas structs en `pkg/models` o métodos en `local-app/app.go`:

```powershell
.\scripts\client.ps1 bindings
```
