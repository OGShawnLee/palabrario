<#
.SYNOPSIS
    Script de métricas y conteo de líneas de código (LOC) para Palabrario.
.DESCRIPTION
    Inspecciona los archivos del proyecto excluyendo dependencias, cachés y
    artefactos autogenerados, emitiendo un desglose por lenguaje y una radiografía
    detallada del código en Go.
#>

[CmdletBinding()]
param()

$ErrorActionPreference = 'Stop'

# Extensiones reconocidas y categorizadas
$categoryMap = [ordered]@{
    'Go'            = @('*.go')
    'Vue 3'         = @('*.vue')
    'Markdown/Docs' = @('*.md')
    'SQL'           = @('*.sql')
    'TypeScript'    = @('*.ts')
    'DevOps/Config' = @('*.yml', '*.yaml', '*.json', 'Dockerfile', 'Vagrantfile', '*.ps1', '*.toml')
    'HTML/Otros'    = @('*.html', '*.example', '.gitignore')
}

# Patrones a excluir
$excludePatterns = @(
    '\.git\',
    'node_modules\',
    'wailsjs\',
    'dist\',
    '\.vscode\',
    'build\bin\',
    '\.tmp\',
    '\.work',
    '\.sum',
    'bun\.lock',
    'package-lock\.json',
    '\.md5'
)

$files = Get-ChildItem -Path . -Recurse -File | Where-Object {
    $filePath = $_.FullName
    $skip = $false
    foreach ($pat in $excludePatterns) {
        if ($filePath -match [regex]::Escape($pat)) {
            $skip = $true
            break
        }
    }
    -not $skip
}

# Estructuras para agregación
$stats = [ordered]@{}
foreach ($cat in $categoryMap.Keys) {
    $stats[$cat] = [PSCustomObject]@{
        Files = 0
        Lines = 0
    }
}

$goBreakdown = [ordered]@{
    'Backend SOA (Core/Realtime/Image/Pkg)' = [PSCustomObject]@{ Files = 0; Lines = 0 }
    'Suites de Pruebas Go (*_test.go)'      = [PSCustomObject]@{ Files = 0; Lines = 0 }
    'Puente Wails Go (local-app/app.go)'     = [PSCustomObject]@{ Files = 0; Lines = 0 }
}

$totalLines = 0
$totalFiles = 0

foreach ($file in $files) {
    $fileName = $file.Name.ToLower()
    $matchedCategory = $null

    foreach ($cat in $categoryMap.Keys) {
        foreach ($glob in $categoryMap[$cat]) {
            if ($fileName -like $glob.ToLower()) {
                $matchedCategory = $cat
                break
            }
        }
        if ($matchedCategory) { break }
    }

    if (-not $matchedCategory) { continue }

    try {
        $lineCount = (Get-Content -Path $file.FullName -ErrorAction SilentlyContinue | Measure-Object -Line).Lines
    } catch {
        $lineCount = 0
    }

    $stats[$matchedCategory].Files++
    $stats[$matchedCategory].Lines += $lineCount
    $totalLines += $lineCount
    $totalFiles++

    # Desglose específico para Go
    if ($matchedCategory -eq 'Go') {
        $relPath = $file.FullName.Replace((Get-Location).Path, '').TrimStart('\/')
        if ($fileName -like '*_test.go') {
            $goBreakdown['Suites de Pruebas Go (*_test.go)'].Files++
            $goBreakdown['Suites de Pruebas Go (*_test.go)'].Lines += $lineCount
        } elseif ($relPath -like 'local-app*') {
            $goBreakdown['Puente Wails Go (local-app/app.go)'].Files++
            $goBreakdown['Puente Wails Go (local-app/app.go)'].Lines += $lineCount
        } else {
            $goBreakdown['Backend SOA (Core/Realtime/Image/Pkg)'].Files++
            $goBreakdown['Backend SOA (Core/Realtime/Image/Pkg)'].Lines += $lineCount
        }
    }
}

Write-Host ""
Write-Host "================================================================================" -ForegroundColor Cyan
Write-Host "            PALABRARIO - REPORTE DE VOLUMETRIA Y LINEAS DE CODIGO (LOC)        " -ForegroundColor Yellow
Write-Host "================================================================================" -ForegroundColor Cyan
Write-Host ""

$table = [System.Collections.Generic.List[PSCustomObject]]::new()
foreach ($cat in $stats.Keys) {
    $fCount = $stats[$cat].Files
    $lCount = $stats[$cat].Lines
    $pct = if ($totalLines -gt 0) { [math]::Round(($lCount / $totalLines) * 100, 1) } else { 0 }
    $table.Add([PSCustomObject]@{
        "Lenguaje / Rol" = $cat
        "Archivos"       = $fCount
        "Lineas"         = ("{0:N0}" -f $lCount)
        "% Total"        = "$pct %"
    })
}

$table | Format-Table -AutoSize | Out-String | Write-Host -ForegroundColor White

Write-Host "--------------------------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host (" TOTAL PROYECTO: {0} archivos analizados | {1:N0} lineas de codigo y docs" -f $totalFiles, $totalLines) -ForegroundColor Green
Write-Host "--------------------------------------------------------------------------------" -ForegroundColor DarkGray
Write-Host ""

Write-Host ">> RADIOGRAFIA DEL BACKEND GO (100% Go 1.24):" -ForegroundColor Yellow
$goTable = [System.Collections.Generic.List[PSCustomObject]]::new()
$totalGoLines = $stats['Go'].Lines

foreach ($key in $goBreakdown.Keys) {
    $fCount = $goBreakdown[$key].Files
    $lCount = $goBreakdown[$key].Lines
    $pct = if ($totalGoLines -gt 0) { [math]::Round(($lCount / $totalGoLines) * 100, 1) } else { 0 }
    $goTable.Add([PSCustomObject]@{
        "Capa / Componente" = $key
        "Archivos"          = $fCount
        "Lineas"            = ("{0:N0}" -f $lCount)
        "% de Go"           = "$pct %"
    })
}

$goTable | Format-Table -AutoSize | Out-String | Write-Host -ForegroundColor Cyan
