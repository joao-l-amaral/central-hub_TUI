param(
    [Parameter(Mandatory = $true)]
    [string]$TargetPath
)

# Prefer opening a tab in an existing Windows Terminal window.
# If that fails, try opening a new WT tab/window, then fallback to plain PowerShell.
$setLocationCmd = "Set-Location -LiteralPath '$TargetPath'"

try {
    wt -w 0 new-tab pwsh -NoExit -Command $setLocationCmd
    exit 0
} catch {
    # continue to fallback
}

try {
    wt new-tab pwsh -NoExit -Command $setLocationCmd
    exit 0
} catch {
    # continue to fallback
}

Start-Process cmd -ArgumentList "/C", "start", "pwsh", "-NoExit", "-Command", $setLocationCmd | Out-Null
