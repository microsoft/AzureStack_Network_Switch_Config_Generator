@{
    # Module metadata
    RootModule = 'EnvironmentDetailAssignment.psm1'
    ModuleVersion = '1.0.0'
    GUID = 'e1f2g3h4-i5j6-7890-1234-567890abcdef'
    Author = 'Network Engineering Team'
    CompanyName = 'Azure Stack Framework'
    Copyright = '(c) 2025 Azure Stack Framework. All rights reserved.'
    Description = 'Environment configuration management utilities for network infrastructure deployment across multiple environments'
    
    # Minimum PowerShell version
    PowerShellVersion = '5.1'
    
    # Supported PSEditions
    CompatiblePSEditions = @('Desktop', 'Core')
    
    # Functions to export
    FunctionsToExport = @('New-EnvironmentAssignment')
    
    # Cmdlets to export
    CmdletsToExport = @()
    
    # Variables to export
    VariablesToExport = @()
    
    # Aliases to export
    AliasesToExport = @()
    
    # Private data
    PrivateData = @{
        PSData = @{
            Tags = @('Network', 'Environment', 'Configuration', 'Azure', 'Infrastructure', 'Deployment', 'DevOps')
            LicenseUri = ''
            ProjectUri = ''
            IconUri = ''
            ReleaseNotes = 'Initial release of environment configuration management and parameter assignment functions'
        }
    }
}