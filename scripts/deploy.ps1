# scripts/deploy.ps1 -- Cross-build rmkbd (ARMv7 static) on Windows.
# The ThinkPad cannot reach the device (Always-On VPN), so this script only
# builds the binary. Device deployment requires the Mac -- run deploy.sh there.
#
# Usage (from repo root):
#   .\scripts\deploy.ps1

Set-StrictMode -Version Latest
$ErrorActionPreference = "Stop"

$REPO = Split-Path -Parent $PSScriptRoot

Write-Host "=== rmkbd: cross-build (ARMv7 static) ==="
Push-Location "$REPO\daemon"
try {
    $env:GOOS       = "linux"
    $env:GOARCH     = "arm"
    $env:GOARM      = "7"
    $env:CGO_ENABLED = "0"
    go build -trimpath -o "$REPO\rmkbd" .
    Write-Host "  built: $REPO\rmkbd"
} finally {
    Remove-Item Env:\GOOS, Env:\GOARCH, Env:\GOARM, Env:\CGO_ENABLED -ErrorAction SilentlyContinue
    Pop-Location
}

Write-Host ""
Write-Host "ThinkPad cannot reach the device (VPN)."
Write-Host "To deploy, commit + push, then on the Mac run:"
Write-Host "  bash scripts/deploy.sh"
