<#
.SYNOPSIS
    Orquestador automatizado del cliente de escritorio Wails v3 + Bun para Palabrario.

.DESCRIPTION
    Encadena las tareas de desarrollo, verificación de tipos, compilación y ejecución de la aplicación cliente.
    Garantiza automáticamente que el módulo esté sincronizado en go.work antes de ejecutar.

.PARAMETER Action
    Acción a ejecutar:
    - check: Verificación rápida total (Bindings Go->TS + Chequeo de tipos + Build Vite + Build Go). Recomendado para validar sub-fases.
    - build: Compilación completa de producción (Genera el binario .exe nativo en bin/palabrario.exe).
    - run: Ejecuta el binario compilado.
    - start: Compila y ejecuta la aplicación de escritorio en vivo.
    - dev: Inicia el servidor de desarrollo en caliente (Wails dev mode + Vite dev).
    - bindings: Regenera exclusivamente los bindings TypeScript desde los structs y servicios Go.

.EXAMPLE
    .\scripts\client.ps1 check
    .\scripts\client.ps1 start
    .\scripts\client.ps1 dev
#>
param (
    [ValidateSet("check", "build", "run", "start", "dev", "bindings")]
    [string]$Action = "check"
)

$ErrorActionPreference = "Stop"

$ProjectRoot = Split-Path -Parent $PSScriptRoot
$ClientDir = Join-Path $ProjectRoot "local-app"

# 1. Asegurar que go.work incluye local-app
Write-Host "==> Verificando go.work en la raiz..." -ForegroundColor Cyan
Push-Location $ProjectRoot
try {
    $workContent = Get-Content "go.work" -Raw -ErrorAction SilentlyContinue
    if ($workContent -notmatch "\./local-app") {
        Write-Host "==> Registrando ./local-app en go.work..." -ForegroundColor Yellow
        go work use ./local-app
    }
} finally {
    Pop-Location
}

# 2. Ejecutar la tarea solicitada en local-app
Write-Host "==> Ejecutando tarea '$Action' en Wails v3..." -ForegroundColor Green
Push-Location $ClientDir
try {
    switch ($Action) {
        "check"    { wails3 task check }
        "build"    { wails3 task build }
        "run"      { wails3 task run }
        "start"    { wails3 task start }
        "dev"      { wails3 task dev }
        "bindings" { wails3 task bindings }
    }
} finally {
    Pop-Location
}
