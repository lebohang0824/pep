$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
$Binary = Join-Path $Root "pep.exe"

Write-Host "  • Building pep..." -NoNewline
go build -o $Binary $Root
Write-Host " ✓" -ForegroundColor Green

$LocalBin = Join-Path $env:USERPROFILE ".local\bin"
if (-not (Test-Path $LocalBin)) {
    New-Item -ItemType Directory -Path $LocalBin -Force | Out-Null
}

$Symlink = Join-Path $LocalBin "pep.exe"
if (-not (Test-Path $Symlink)) {
    New-Item -ItemType SymbolicLink -Path $Symlink -Target $Binary -Force | Out-Null
    Write-Host "    ✓ Symlinked to $Symlink" -ForegroundColor Green
}
else {
    Write-Host "    ! $Symlink already exists" -ForegroundColor Yellow
}

$VscodeExt = Join-Path $env:USERPROFILE ".vscode\extensions\pep-lang.pep-lang"
if (Test-Path $VscodeExt) {
    Write-Host "  • Updating VS Code extension..." -NoNewline
}
else {
    Write-Host "  • Installing VS Code extension..." -NoNewline
}
Remove-Item -Recurse -Force $VscodeExt -ErrorAction SilentlyContinue
Copy-Item -Recurse (Join-Path $Root "vscode-pep") $VscodeExt
Write-Host " ✓" -ForegroundColor Green

Write-Host ""
Write-Host "  ✓ Pep installed. Reload VS Code to activate the extension." -ForegroundColor Green
Write-Host "  ✓ Run 'pep --help' to get started." -ForegroundColor Green
Write-Host ""
