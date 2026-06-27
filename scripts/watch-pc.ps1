#Requires -Version 5.1
<#
.SYNOPSIS
  Auto-pull the git bridge on the PC and pop a GUI toast when results arrive.
.DESCRIPTION
  Loops: git pull. When a pull brings in new commits (HEAD moves) -- e.g. the
  Mac just pushed a device recon result -- it shows a Windows toast (GUI,
  vanishing). Toasts fire on: ARM (start), each applied pull, and STOP, so you
  can tell at a glance whether the bridge is running.

  No admin, no modules: uses the built-in WinRT toast API, with a tray-balloon
  fallback if toast is unavailable. git's progress on stderr is expected and
  ignored.
.EXAMPLE  .\scripts\watch-pc.ps1
.EXAMPLE  .\scripts\watch-pc.ps1 -IntervalSeconds 30
#>
param([int]$IntervalSeconds = 20)

# git writes normal progress to stderr, which PowerShell flags as errors. Keep going.
$ErrorActionPreference = 'Continue'
$repo = Split-Path $PSScriptRoot -Parent
Set-Location $repo

function Show-Toast {
    param([string]$Text, [string]$Title = 'rmwatch (PC)')
    try {
        $null = [Windows.UI.Notifications.ToastNotificationManager, Windows.UI.Notifications, ContentType = WindowsRuntime]
        $AppId = '{1AC14E77-02E7-4E5D-B744-2EB1AE5198B7}\WindowsPowerShell\v1.0\powershell.exe'
        $xml = [Windows.UI.Notifications.ToastNotificationManager]::GetTemplateContent([Windows.UI.Notifications.ToastTemplateType]::ToastText02)
        $nodes = $xml.GetElementsByTagName('text')
        $nodes.Item(0).AppendChild($xml.CreateTextNode($Title)) | Out-Null
        $nodes.Item(1).AppendChild($xml.CreateTextNode($Text))  | Out-Null
        $toast = [Windows.UI.Notifications.ToastNotification]::new($xml)
        [Windows.UI.Notifications.ToastNotificationManager]::CreateToastNotifier($AppId).Show($toast)
    } catch {
        # Fallback: tray balloon (also GUI, also vanishing).
        try {
            Add-Type -AssemblyName System.Windows.Forms
            Add-Type -AssemblyName System.Drawing
            $ni = New-Object System.Windows.Forms.NotifyIcon
            $ni.Icon = [System.Drawing.SystemIcons]::Information
            $ni.Visible = $true
            $ni.ShowBalloonTip(4000, $Title, $Text, [System.Windows.Forms.ToolTipIcon]::Info)
            Start-Sleep -Milliseconds 250
            $ni.Dispose()
        } catch {
            Write-Host "[toast unavailable] $Text"
        }
    }
}

Show-Toast "ARMED - auto-pulling every $IntervalSeconds s"
Write-Host "rmwatch (PC): pulling every $IntervalSeconds s. Ctrl-C to stop." -ForegroundColor Cyan

try {
    while ($true) {
        $before = (git rev-parse HEAD 2>$null)
        git pull --rebase --autostash 2>&1 | Out-Null
        $after = (git rev-parse HEAD 2>$null)
        if ($before -and $after -and ($before -ne $after)) {
            $subject = (git --no-pager log -1 --format='%s' 2>$null)
            Show-Toast "pulled: $subject"
            Write-Host "  pulled: $subject" -ForegroundColor Green
        }
        Start-Sleep -Seconds $IntervalSeconds
    }
} finally {
    Show-Toast "STOPPED - PC is NO LONGER auto-pulling"
    Write-Host "rmwatch (PC): stopped." -ForegroundColor Yellow
}
