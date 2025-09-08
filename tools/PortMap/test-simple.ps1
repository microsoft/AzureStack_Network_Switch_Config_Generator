#Requires -Version 5.1

[CmdletBinding()]
param(
    [Parameter(Mandatory = $true)]
    [string]$InputFile,
    
    [Parameter(Mandatory = $true)]
    [string]$OutputFormat
)

Write-Host "Script started successfully"
Write-Host "InputFile: $InputFile"
Write-Host "OutputFormat: $OutputFormat" 
Write-Host "Testing file access..."

if (Test-Path $InputFile) {
    Write-Host "File exists: $InputFile"
    $content = Get-Content $InputFile -Raw
    Write-Host "File content length: $($content.Length)"
} else {
    Write-Host "File does not exist: $InputFile"
}

Write-Host "Script completed successfully"
