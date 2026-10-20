<#
.SYNOPSIS
    Orquestador automatizado de pruebas para Palabrario (PostgreSQL 17 Efímero en RAM + Go Tests 80/20).

.DESCRIPTION
    Soporta tanto Docker Desktop nativo en Windows como Docker Engine dentro de la VM Vagrant Debian 12.
    - Si 'postgres-test' ya está corriendo en :5433 (por ejemplo en tu entorno de desarrollo o Vagrant),
      reutiliza instantáneamente esa instancia en lugar de reiniciar ni apagar el contenedor.
    - Si 'postgres-test' estaba apagado, lo levanta en RAM, espera el healthcheck, corre las pruebas y
      lo apaga al finalizar (salvo que pases -KeepAlive).

.PARAMETER KeepAlive
    Mantiene el contenedor 'postgres-test' encendido después de correr las pruebas para
    iteraciones rápidas de desarrollo.

.PARAMETER Down
    Apaga y destruye inmediatamente el contenedor 'postgres-test'.

.PARAMETER Target
    Paquetes de Go a probar (por defecto: "./pkg/... ./core-service/...").
#>
param (
    [switch]$KeepAlive,
    [switch]$Down,
    [string]$Target = "./pkg/... ./core-service/..."
)

$ErrorActionPreference = "Continue"

function Invoke-DockerCompose {
    param([string]$ComposeArgs)

    if (Get-Command docker -ErrorAction SilentlyContinue) {
        & cmd /c "docker compose $ComposeArgs"
        return $LASTEXITCODE
    }
    elseif (Get-Command vagrant -ErrorAction SilentlyContinue) {
        Write-Host "🐧 Ejecutando Docker Compose vía Vagrant VM (/home/vagrant/palabrio)..." -ForegroundColor DarkCyan
        & cmd /c "vagrant ssh -c ""cd /home/vagrant/palabrio && docker compose $ComposeArgs"""
        return $LASTEXITCODE
    }
    else {
        Write-Host "❌ No se encontró ni 'docker' ni 'vagrant' en el PATH." -ForegroundColor Red
        return 1
    }
}

function Test-Port5433Open {
    $client = New-Object System.Net.Sockets.TcpClient
    try {
        $async = $client.BeginConnect("127.0.0.1", 5433, $null, $null)
        $wait = $async.AsyncWaitHandle.WaitOne(800, $false)
        if ($wait -and $client.Connected) {
            return $true
        }
        return $false
    }
    catch {
        return $false
    }
    finally {
        $client.Close()
    }
}

if ($Down) {
    Write-Host "🛑 Apagando contenedor efímero postgres-test (:5433)..." -ForegroundColor Yellow
    $code = Invoke-DockerCompose "--profile test rm -s -f -v postgres-test"
    if ($code -eq 0) {
        Write-Host "✅ Contenedor postgres-test detenido y destruido." -ForegroundColor Green
    }
    exit $code
}

$wasAlreadyRunning = Test-Port5433Open

if ($wasAlreadyRunning) {
    Write-Host "⚡ [1/3] Contenedor postgres-test detectado activo en localhost:5433 (Reutilizando instancia en vivo)..." -ForegroundColor Green
} else {
    Write-Host "🚀 [1/3] Levantando PostgreSQL 17 efímero en RAM (postgres-test :5433)..." -ForegroundColor Cyan
    if (-not (Get-Command docker -ErrorAction SilentlyContinue) -and (Get-Command vagrant -ErrorAction SilentlyContinue)) {
        Write-Host "🔄 Sincronizando SCHEMA_FINAL.sql y docker-compose.yml hacia Vagrant VM..." -ForegroundColor DarkCyan
        & cmd /c "vagrant rsync"
    }

    $upCode = Invoke-DockerCompose "--profile test up -d --wait postgres-test"
    if ($upCode -ne 0) {
        Write-Host "❌ Error al levantar postgres-test. Asegúrate de que Docker o la VM de Vagrant ('vagrant up') estén encendidos." -ForegroundColor Red
        exit $upCode
    }
}

$env:TEST_DATABASE_URL = "postgresql://postgres:postgres@localhost:5433/postgres?sslmode=disable"
$testExitCode = 0

try {
    Write-Host "🧪 [2/3] Ejecutando batería de pruebas Go (Clonado por TEMPLATE en RAM)..." -ForegroundColor Cyan
    $targets = $Target -split " "
    & go test -v -count=1 @targets
    $testExitCode = $LASTEXITCODE
}
finally {
    if (-not $KeepAlive -and -not $wasAlreadyRunning) {
        Write-Host "🧹 [3/3] Apagando contenedor efímero postgres-test..." -ForegroundColor Cyan
        [void](Invoke-DockerCompose "--profile test rm -s -f -v postgres-test")
    } else {
        Write-Host "⚡ [3/3] postgres-test se mantiene activo en :5433. (Usa '.\scripts\test.ps1 -Down' si deseas apagarlo)." -ForegroundColor Yellow
    }
}

if ($testExitCode -eq 0) {
    Write-Host "✅ TODAS LAS PRUEBAS PASARON EXITOSAMENTE." -ForegroundColor Green
} else {
    Write-Host "❌ ALGUNAS PRUEBAS FALLARON (Código: $testExitCode)." -ForegroundColor Red
}

exit $testExitCode
