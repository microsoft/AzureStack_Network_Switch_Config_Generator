<#
.SYNOPSIS
    Test suite for EnvironmentDetailAssignment Tool

.DESCRIPTION
    This script provides comprehensive testing for the EnvironmentDetailAssignment tool,
    validating configuration processing, output generation, and error handling.

.NOTES
    File Name      : Test-EnvironmentDetailAssignment.ps1
    Author         : Network Engineering Team
    Prerequisite   : PowerShell 5.1 or later
    Copyright      : (c) 2025 Azure Stack Framework. All rights reserved.
    Version        : 1.0.0
#>

[CmdletBinding()]
param(
    [Parameter(Mandatory = $false)]
    [switch]$Detailed,
    
    [Parameter(Mandatory = $false)]
    [switch]$KeepOutputFiles
)

Set-StrictMode -Version Latest

# Test tracking variables
$Script:TestResults = @()
$Script:TestsPassed = 0
$Script:TestsFailed = 0
$Script:TestsSkipped = 0

function Write-TestHeader {
    param([string]$Title)
    Write-Host ""
    Write-Host "🧪 $Title" -ForegroundColor Cyan
    Write-Host ("=" * ($Title.Length + 3)) -ForegroundColor Cyan
}

function Write-TestResult {
    param(
        [string]$TestName,
        [string]$Result,
        [string]$Message = ""
    )
    
    $testResult = @{
        TestName = $TestName
        Result = $Result
        Message = $Message
        Duration = 0
    }
    
    $Script:TestResults += $testResult
    
    switch ($Result) {
        "Pass" { 
            Write-Host "  ✅ $TestName" -ForegroundColor Green
            $Script:TestsPassed++
        }
        "Fail" { 
            Write-Host "  ❌ $TestName" -ForegroundColor Red
            if ($Message) { Write-Host "     $Message" -ForegroundColor Yellow }
            $Script:TestsFailed++
        }
        "Skip" { 
            Write-Host "  ⏭️  $TestName" -ForegroundColor Yellow
            if ($Message) { Write-Host "     $Message" -ForegroundColor Gray }
            $Script:TestsSkipped++
        }
    }
}

function Test-FileExists {
    param([string]$Path, [string]$Description)
    
    if (Test-Path $Path) {
        Write-TestResult -TestName "File exists: $Description" -Result "Pass"
        return $true
    } else {
        Write-TestResult -TestName "File exists: $Description" -Result "Fail" -Message "File not found: $Path"
        return $false
    }
}

function Test-JsonSyntax {
    param([string]$Path, [string]$Description)
    
    try {
        $null = Get-Content -Path $Path -Raw | ConvertFrom-Json
        Write-TestResult -TestName "JSON syntax: $Description" -Result "Pass"
        return $true
    } catch {
        Write-TestResult -TestName "JSON syntax: $Description" -Result "Fail" -Message $_.Exception.Message
        return $false
    }
}

function Test-PowerShellSyntax {
    param([string]$Path, [string]$Description)
    
    try {
        $null = [System.Management.Automation.PSParser]::Tokenize((Get-Content -Path $Path -Raw), [ref]$null)
        Write-TestResult -TestName "PowerShell syntax: $Description" -Result "Pass"
        return $true
    } catch {
        Write-TestResult -TestName "PowerShell syntax: $Description" -Result "Fail" -Message $_.Exception.Message
        return $false
    }
}

function Test-OutputContent {
    param(
        [string]$Path,
        [string]$ExpectedContent,
        [string]$Description
    )
    
    try {
        $content = Get-Content -Path $Path -Raw
        if ($content -like "*$ExpectedContent*") {
            Write-TestResult -TestName "Output content: $Description" -Result "Pass"
            return $true
        } else {
            Write-TestResult -TestName "Output content: $Description" -Result "Fail" -Message "Expected content not found: $ExpectedContent"
            return $false
        }
    } catch {
        Write-TestResult -TestName "Output content: $Description" -Result "Fail" -Message $_.Exception.Message
        return $false
    }
}

function Cleanup-TestFiles {
    param([string[]]$Files)
    
    if (-not $KeepOutputFiles) {
        foreach ($file in $Files) {
            if (Test-Path $file) {
                Remove-Item $file -Force -ErrorAction SilentlyContinue
            }
        }
    }
}

# Main test execution
try {
    Write-Host "🧪 EnvironmentDetailAssignment Tool - Test Suite" -ForegroundColor Cyan
    Write-Host "=============================================" -ForegroundColor Cyan
    Write-Host ""
    
    $ScriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
    $TestOutputDir = Join-Path $ScriptDir "test-output"
    
    # Create test output directory
    if (-not (Test-Path $TestOutputDir)) {
        New-Item -ItemType Directory -Path $TestOutputDir -Force | Out-Null
    }
    
    # Test 1: File Structure Tests
    Write-TestHeader "File Structure Tests"
    
    $RequiredFiles = @(
        @{ Path = Join-Path $ScriptDir "EnvironmentDetailAssignment.ps1"; Desc = "Main script" },
        @{ Path = Join-Path $ScriptDir "EnvironmentDetailAssignment.psm1"; Desc = "PowerShell module" },
        @{ Path = Join-Path $ScriptDir "EnvironmentDetailAssignment.psd1"; Desc = "Module manifest" },
        @{ Path = Join-Path $ScriptDir "sample-environment-config.json"; Desc = "Sample configuration" },
        @{ Path = Join-Path $ScriptDir "multi-environment-config.json"; Desc = "Multi-environment configuration" }
    )
    
    foreach ($file in $RequiredFiles) {
        Test-FileExists -Path $file.Path -Description $file.Desc
    }
    
    # Test 2: Configuration File Validation
    Write-TestHeader "Configuration File Validation"
    
    $ConfigFiles = @(
        @{ Path = Join-Path $ScriptDir "sample-environment-config.json"; Desc = "Sample config" },
        @{ Path = Join-Path $ScriptDir "multi-environment-config.json"; Desc = "Multi-environment config" }
    )
    
    foreach ($config in $ConfigFiles) {
        Test-JsonSyntax -Path $config.Path -Description $config.Desc
    }
    
    # Test 3: Module Loading Tests
    Write-TestHeader "Module Loading Tests"
    
    try {
        $ModulePath = Join-Path $ScriptDir "EnvironmentDetailAssignment.psm1"
        Import-Module $ModulePath -Force
        Write-TestResult -TestName "Module import" -Result "Pass"
        
        # Test function availability
        $Functions = Get-Command -Module EnvironmentDetailAssignment -ErrorAction SilentlyContinue
        if ($Functions -and @($Functions).Count -gt 0) {
            Write-TestResult -TestName "Module functions available" -Result "Pass"
        } else {
            Write-TestResult -TestName "Module functions available" -Result "Fail" -Message "No functions exported"
        }
        
    } catch {
        Write-TestResult -TestName "Module import" -Result "Fail" -Message $_.Exception.Message
    }
    
    # Test 4: Configuration Validation Tests
    Write-TestHeader "Configuration Validation Tests"
    
    $SampleConfig = Join-Path $ScriptDir "sample-environment-config.json"
    $ValidationScript = Join-Path $ScriptDir "EnvironmentDetailAssignment.ps1"
    
    if ((Test-Path $SampleConfig) -and (Test-Path $ValidationScript)) {
        try {
            & $ValidationScript -ConfigPath $SampleConfig -OutputFormat JSON -Validate
            Write-TestResult -TestName "Configuration validation" -Result "Pass"
        } catch {
            Write-TestResult -TestName "Configuration validation" -Result "Fail" -Message $_.Exception.Message
        }
    } else {
        Write-TestResult -TestName "Configuration validation" -Result "Skip" -Message "Required files not found"
    }
    
    # Test 5: JSON Output Generation Tests
    Write-TestHeader "JSON Output Generation Tests"
    
    if ((Test-Path $SampleConfig) -and (Test-Path $ValidationScript)) {
        try {
            $JsonOutput = Join-Path $TestOutputDir "test-output.json"
            & $ValidationScript -ConfigPath $SampleConfig -OutputFormat JSON -OutputFile $JsonOutput
            
            if (Test-Path $JsonOutput) {
                Write-TestResult -TestName "JSON output file creation" -Result "Pass"
                Test-JsonSyntax -Path $JsonOutput -Description "Generated JSON output"
                Test-OutputContent -Path $JsonOutput -ExpectedContent "EnvironmentDetailAssignment" -Description "JSON metadata"
                Test-OutputContent -Path $JsonOutput -ExpectedContent "Production" -Description "Production environment"
            } else {
                Write-TestResult -TestName "JSON output file creation" -Result "Fail" -Message "Output file not created"
            }
        } catch {
            Write-TestResult -TestName "JSON output generation" -Result "Fail" -Message $_.Exception.Message
        }
    } else {
        Write-TestResult -TestName "JSON output generation" -Result "Skip" -Message "Required files not found"
    }
    
    # Test 6: CSV Output Generation Tests
    Write-TestHeader "CSV Output Generation Tests"
    
    if ((Test-Path $SampleConfig) -and (Test-Path $ValidationScript)) {
        try {
            $CsvOutput = Join-Path $TestOutputDir "test-output.csv"
            & $ValidationScript -ConfigPath $SampleConfig -OutputFormat CSV -OutputFile $CsvOutput
            
            if (Test-Path $CsvOutput) {
                Write-TestResult -TestName "CSV output file creation" -Result "Pass"
                Test-OutputContent -Path $CsvOutput -ExpectedContent "Parameter" -Description "CSV header"
                Test-OutputContent -Path $CsvOutput -ExpectedContent "NetworkPrefix" -Description "Network parameters"
            } else {
                Write-TestResult -TestName "CSV output file creation" -Result "Fail" -Message "Output file not created"
            }
        } catch {
            Write-TestResult -TestName "CSV output generation" -Result "Fail" -Message $_.Exception.Message
        }
    } else {
        Write-TestResult -TestName "CSV output generation" -Result "Skip" -Message "Required files not found"
    }
    
    # Test 7: PowerShell Output Generation Tests
    Write-TestHeader "PowerShell Output Generation Tests"
    
    if ((Test-Path $SampleConfig) -and (Test-Path $ValidationScript)) {
        try {
            $PsOutput = Join-Path $TestOutputDir "test-output.ps1"
            & $ValidationScript -ConfigPath $SampleConfig -Environment "Production" -OutputFormat PowerShell -OutputFile $PsOutput
            
            if (Test-Path $PsOutput) {
                Write-TestResult -TestName "PowerShell output file creation" -Result "Pass"
                Test-PowerShellSyntax -Path $PsOutput -Description "Generated PowerShell output"
                Test-OutputContent -Path $PsOutput -ExpectedContent "`$EnvironmentName" -Description "Environment variables"
                Test-OutputContent -Path $PsOutput -ExpectedContent "Production" -Description "Production environment"
            } else {
                Write-TestResult -TestName "PowerShell output file creation" -Result "Fail" -Message "Output file not created"
            }
        } catch {
            Write-TestResult -TestName "PowerShell output generation" -Result "Fail" -Message $_.Exception.Message
        }
    } else {
        Write-TestResult -TestName "PowerShell output generation" -Result "Skip" -Message "Required files not found"
    }
    
    # Test 8: Multi-Environment Processing Tests
    Write-TestHeader "Multi-Environment Processing Tests"
    
    $MultiConfig = Join-Path $ScriptDir "multi-environment-config.json"
    if ((Test-Path $MultiConfig) -and (Test-Path $ValidationScript)) {
        try {
            $MultiOutput = Join-Path $TestOutputDir "multi-env-test.json"
            & $ValidationScript -ConfigPath $MultiConfig -Environments @("Development", "Production") -OutputFormat JSON -OutputFile $MultiOutput
            
            if (Test-Path $MultiOutput) {
                Write-TestResult -TestName "Multi-environment processing" -Result "Pass"
                Test-OutputContent -Path $MultiOutput -ExpectedContent "Development" -Description "Development environment"
                Test-OutputContent -Path $MultiOutput -ExpectedContent "Production" -Description "Production environment"
            } else {
                Write-TestResult -TestName "Multi-environment processing" -Result "Fail" -Message "Output file not created"
            }
        } catch {
            Write-TestResult -TestName "Multi-environment processing" -Result "Fail" -Message $_.Exception.Message
        }
    } else {
        Write-TestResult -TestName "Multi-environment processing" -Result "Skip" -Message "Required files not found"
    }
    
    # Test 9: Error Handling Tests
    Write-TestHeader "Error Handling Tests"
    
    # Test with invalid JSON
    try {
        $InvalidJson = Join-Path $TestOutputDir "invalid-config.json"
        '{ "invalid": json syntax }' | Out-File -FilePath $InvalidJson -Encoding UTF8
        
        $process = Start-Process -FilePath "pwsh" -ArgumentList "-NoProfile", "-Command", "`"& '$ValidationScript' -ConfigPath '$InvalidJson' -OutputFormat JSON -Validate`"" -Wait -PassThru -WindowStyle Hidden -RedirectStandardError "$TestOutputDir\error.log"
        
        if ($process.ExitCode -ne 0) {
            Write-TestResult -TestName "Invalid JSON handling" -Result "Pass"
        } else {
            Write-TestResult -TestName "Invalid JSON handling" -Result "Fail" -Message "Should have failed with invalid JSON"
        }
    } catch {
        Write-TestResult -TestName "Invalid JSON handling" -Result "Pass" -Message "Correctly caught invalid JSON error"
    }
    
    # Test with missing file
    try {
        $process = Start-Process -FilePath "pwsh" -ArgumentList "-NoProfile", "-Command", "`"& '$ValidationScript' -ConfigPath 'nonexistent-file.json' -OutputFormat JSON -Validate`"" -Wait -PassThru -WindowStyle Hidden -RedirectStandardError "$TestOutputDir\error2.log"
        
        if ($process.ExitCode -ne 0) {
            Write-TestResult -TestName "Missing file handling" -Result "Pass"
        } else {
            Write-TestResult -TestName "Missing file handling" -Result "Fail" -Message "Should have failed with missing file"
        }
    } catch {
        Write-TestResult -TestName "Missing file handling" -Result "Pass" -Message "Correctly caught missing file error"
    }
    
    # Cleanup test files
    $TestFiles = @(
        (Join-Path $TestOutputDir "test-output.json"),
        (Join-Path $TestOutputDir "test-output.csv"),
        (Join-Path $TestOutputDir "test-output.ps1"),
        (Join-Path $TestOutputDir "multi-env-test.json"),
        (Join-Path $TestOutputDir "invalid-config.json")
    )
    
    Cleanup-TestFiles -Files $TestFiles
    
    # Summary
    Write-Host ""
    Write-Host "📊 Test Summary" -ForegroundColor Cyan
    Write-Host "===============" -ForegroundColor Cyan
    
    $TotalTests = $Script:TestsPassed + $Script:TestsFailed + $Script:TestsSkipped
    $SuccessRate = if ($TotalTests -gt 0) { [Math]::Round(($Script:TestsPassed / $TotalTests) * 100, 1) } else { 0 }
    
    Write-Host "Total Tests: $TotalTests" -ForegroundColor White
    Write-Host "Passed: $Script:TestsPassed" -ForegroundColor Green
    Write-Host "Failed: $Script:TestsFailed" -ForegroundColor Red
    Write-Host "Skipped: $Script:TestsSkipped" -ForegroundColor Yellow
    Write-Host "Success Rate: $SuccessRate%" -ForegroundColor $(if ($SuccessRate -ge 90) { "Green" } elseif ($SuccessRate -ge 75) { "Yellow" } else { "Red" })
    
    # Detailed results
    if ($Detailed -and $Script:TestsFailed -gt 0) {
        Write-Host ""
        Write-Host "❌ Failed Tests:" -ForegroundColor Red
        $FailedTests = $Script:TestResults | Where-Object Result -eq "Fail"
        foreach ($test in $FailedTests) {
            Write-Host "  • $($test.TestName): $($test.Message)" -ForegroundColor Red
        }
    }
    
    Write-Host ""
    if ($Script:TestsFailed -eq 0) {
        Write-Host "🎉 All tests passed!" -ForegroundColor Green
        exit 0
    } else {
        Write-Host "⚠️  Some tests failed. Review the results above." -ForegroundColor Yellow
        exit 1
    }
    
} catch {
    Write-Host ""
    Write-Host "❌ Test execution failed: $($_.Exception.Message)" -ForegroundColor Red
    exit 1
} finally {
    # Clean up module
    Remove-Module EnvironmentDetailAssignment -ErrorAction SilentlyContinue
}