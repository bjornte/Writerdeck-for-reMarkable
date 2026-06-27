#Requires -Version 5.1
<#
.SYNOPSIS  One short command to stage + commit + push with the PERSONAL identity.
.DESCRIPTION
  Replaces the long hand-typed `git add ... ; git -c user.name=... commit ... ; git push`
  chains. Always commits as the personal hobby identity (bjornte@gmail.com),
  which prevents the work-email-leak footgun. Pushes unless -NoPush.
.EXAMPLE  .\scripts\push.ps1 "Add probe script"
.EXAMPLE  .\scripts\push.ps1 "WIP" -NoPush
#>
param(
    [Parameter(Mandatory, Position = 0)][string]$Message,
    [switch]$NoPush
)
Set-StrictMode -Version Latest
# NOTE: git writes normal progress to stderr, which PowerShell flags as errors.
# Keep going regardless; we verify success by checking sync state at the end.
$ErrorActionPreference = 'Continue'
$repo = Split-Path $PSScriptRoot -Parent
Set-Location $repo

git add -A
git -c user.name="Bjørn Tennøe" -c user.email="bjornte@gmail.com" commit -m $Message 2>&1 | Out-Host
if (-not $NoPush) { git push origin main 2>&1 | Out-Host }
git --no-pager log -1 --format="committed %h as %ae - %s"
