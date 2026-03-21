param(
[string]$p,
[string]$t = "serve"
)

# Get the list of projects with the specified tag (if any)
if ($p) {
  $projectNames = npx nx show projects -p "$p" --withTarget $t
} else {
  $projectNames = npx nx show projects --withTarget $t
}

if (-not $projectNames) {
  Write-Host "No projects found for the given parameters. Please check and try again."
  exit 1
}

# Start each project in a new Windows Terminal window
foreach ($projectName in $projectNames) {
  Write-Host "Starting project '$projectName'..."
  wt -w 0 -d $(Get-Location) --title "$projectName" --suppressApplicationTitle pwsh -c "npx nx run ${projectName}:$t"
}
 