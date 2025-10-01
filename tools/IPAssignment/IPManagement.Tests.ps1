$here = if ($PSScriptRoot) { $PSScriptRoot } else { Split-Path -Parent $MyInvocation.MyCommand.Path }
$modulePath = Join-Path $here 'IPManagement.psm1'
Import-Module $modulePath -Force

Describe 'New-SubnetPlanFromConfig' {
    It 'emits named assignments and ranges for IPAssignments entries' {
        $jsonConfig = @'
{
  "network": "10.0.0.0/24",
  "subnets": [
    {
      "name": "Mgmt",
      "vlan": 110,
      "cidr": "28",
      "IPAssignments": [
        { "Name": "Gateway", "Position": 1 },
        { "Name": "VMM", "Position": 3 }
      ]
    }
  ]
}
'@

        $resultsJson = New-SubnetPlanFromConfig -JsonConfig $jsonConfig -AsJson
        $results = $resultsJson | ConvertFrom-Json

        $subnetRows = $results | Where-Object { $_.Subnet -eq '10.0.0.0/28' }
        $subnetRows | Should -Not -BeNullOrEmpty

        $networkRow = $subnetRows | Where-Object { $_.Label -eq 'Network' } | Select-Object -First 1
        $networkRow.IP | Should -Be '10.0.0.0'
        $networkRow.Mask | Should -Be '255.255.255.240'
        $networkRow.Name | Should -Be 'Mgmt'
        $networkRow.Vlan | Should -Be 110

        $gatewayRow = $subnetRows | Where-Object { $_.Label -eq 'Gateway' } | Select-Object -First 1
        $gatewayRow | Should -Not -BeNullOrEmpty
        $gatewayRow.IP | Should -Be '10.0.0.1'
        $gatewayRow.Category | Should -Be 'Assignment'
        $gatewayRow.Name | Should -Be 'Mgmt'

        $unusedRows = $subnetRows | Where-Object { $_.Label -eq 'Unused Range' -and $_.Category -eq 'Unused' }
        $unusedRows | Should -Not -BeNullOrEmpty
        ($unusedRows | Where-Object { $_.IP -eq '10.0.0.2 - 10.0.0.2' }) | Should -Not -BeNullOrEmpty
        ($unusedRows | Where-Object { $_.IP -eq '10.0.0.4 - 10.0.0.14' }) | Should -Not -BeNullOrEmpty

        $broadcastRow = $subnetRows | Where-Object { $_.Label -eq 'Broadcast' } | Select-Object -First 1
        $broadcastRow.IP | Should -Be '10.0.0.15'
        $broadcastRow.Category | Should -Be 'Broadcast'
        $broadcastRow.Name | Should -Be 'Mgmt'

        $availableRows = $results | Where-Object { $_.Category -eq 'Available' }
        $availableRows | Should -Not -BeNullOrEmpty
        $availableRows | ForEach-Object { $_.Name } | Should -Contain 'Available'
    }

    It 'includes every IPAssignments entry from network_by_cidr.json' {
        $jsonConfig = @'
{
  "network": "192.168.1.0/24",
  "subnets": [
    {
      "name": "csu-edge-transport-compute",
      "vlan": "203",
      "cidr": "28",
      "IPAssignments": [
        { "Name": "Gateway", "Position": 1 },
        { "Name": "TOR1", "Position": 2 },
        { "Name": "TOR2", "Position": 3 },
        { "Name": "SMTPServer", "Position": 5 },
        { "Name": "VIP", "Position": 4 }
      ]
    },
    {
      "name": "csu-exchange-compute",
      "vlan": "102",
      "cidr": "27",
      "IPAssignments": [
        { "Name": "Gateway", "Position": 1 },
        { "Name": "TOR1", "Position": 2 },
        { "Name": "TOR2", "Position": 3 },
        { "Name": "Mailbox1", "Position": 5 },
        { "Name": "Mailbox2", "Position": 6 },
        { "Name": "Mailbox3", "Position": 7 },
        { "Name": "Mailbox4", "Position": 8 },
        { "Name": "VIP", "Position": 4 }
      ]
    },
    { "name": "msu-compute", "vlan": "302", "cidr": "27" },
    { "name": "csu-edge-transport-management", "vlan": "101", "cidr": "27" },
    { "name": "msu-management", "vlan": "201", "cidr": "27" },
    { "name": "csu-exchange-management", "vlan": "301", "cidr": "26" }
  ]
}
'@

        $resultsJson = New-SubnetPlanFromConfig -JsonConfig $jsonConfig -AsJson
        $results = $resultsJson | ConvertFrom-Json

        $computeAssignments = $results |
            Where-Object { $_.Name -eq 'csu-edge-transport-compute' -and $_.Category -eq 'Assignment' }

        $computeAssignments.Count | Should -Be 5
        ($computeAssignments.Label | Sort-Object) | Should -Be (@('Gateway','SMTPServer','TOR1','TOR2','VIP') | Sort-Object)

        $exchangeAssignments = $results |
            Where-Object { $_.Name -eq 'csu-exchange-compute' -and $_.Category -eq 'Assignment' }

        $exchangeAssignments.Count | Should -Be 8
        ($exchangeAssignments.Label | Sort-Object) | Should -Be (@('Gateway','Mailbox1','Mailbox2','Mailbox3','Mailbox4','TOR1','TOR2','VIP') | Sort-Object)
    }

    It 'exports results to JSON, CSV, and Markdown files' {
        $jsonConfig = @'
{
  "network": "10.0.0.0/24",
  "subnets": [
    {
      "name": "Mgmt",
      "vlan": 110,
      "cidr": "28",
      "IPAssignments": [
        { "Name": "Gateway", "Position": 1 },
        { "Name": "VMM", "Position": 3 }
      ]
    }
  ]
}
'@

        $exportRoot = Join-Path ([System.IO.Path]::GetTempPath()) ([System.Guid]::NewGuid().ToString())
        New-Item -ItemType Directory -Path $exportRoot | Out-Null
        $jsonPath = Join-Path $exportRoot 'plan.json'
        $csvPath = Join-Path $exportRoot 'plan.csv'
        $markdownPath = Join-Path $exportRoot 'plan.md'

        try {
            New-SubnetPlanFromConfig -JsonConfig $jsonConfig -ExportJsonPath $jsonPath -ExportCsvPath $csvPath -ExportMarkdownPath $markdownPath | Out-Null

            Test-Path $jsonPath | Should -BeTrue
            Test-Path $csvPath | Should -BeTrue
            Test-Path $markdownPath | Should -BeTrue

            $jsonData = Get-Content $jsonPath -Raw | ConvertFrom-Json
            ($jsonData | Where-Object { $_.Label -eq 'Gateway' }).IP | Should -Contain '10.0.0.1'

            $csvData = Import-Csv $csvPath
            $csvData[0].Subnet | Should -Be '10.0.0.0/28'

            $markdownContent = Get-Content $markdownPath -Raw
            $markdownContent | Should -Match '^Subnet \|'
            $markdownContent | Should -Match 'Gateway'
        }
        finally {
            if (Test-Path $exportRoot) { Remove-Item $exportRoot -Recurse -Force }
        }
    }

    It 'processes multiple primary networks independently' {
        $jsonConfig = @'
[
  {
    "network": "192.168.1.0/24",
    "subnets": [
      {
        "name": "csu-edge-transport-compute",
        "vlan": "203",
        "cidr": "28",
        "IPAssignments": [
          { "Name": "Gateway", "Position": 1 },
          { "Name": "TOR1", "Position": 2 },
          { "Name": "TOR2", "Position": 3 }
        ]
      },
      {
        "name": "csu-exchange-compute",
        "vlan": "102",
        "cidr": "27"
      }
    ]
  },
  {
    "network": "10.50.1.0/24",
    "subnets": [
      { "name": "csu-edge-transport-management", "vlan": "101", "cidr": "27" },
      { "name": "msu-management", "vlan": "201", "cidr": "27" }
    ]
  }
]
'@

        $resultsJson = New-SubnetPlanFromConfig -JsonConfig $jsonConfig -AsJson
        $results = $resultsJson | ConvertFrom-Json

        # Verify first network (192.168.1.0/24) subnets are present
        $computeSubnets = $results | Where-Object { $_.Subnet -match '^192\.168\.1\.' }
        $computeSubnets | Should -Not -BeNullOrEmpty
        
        $edgeComputeRows = $results | Where-Object { $_.Name -eq 'csu-edge-transport-compute' }
        $edgeComputeRows | Should -Not -BeNullOrEmpty
        
        $exchangeComputeRows = $results | Where-Object { $_.Name -eq 'csu-exchange-compute' }
        $exchangeComputeRows | Should -Not -BeNullOrEmpty

        # Verify IP assignments from first network
        $edgeGateway = $results | Where-Object { $_.Name -eq 'csu-edge-transport-compute' -and $_.Label -eq 'Gateway' }
        $edgeGateway | Should -Not -BeNullOrEmpty
        $edgeGateway.IP | Should -Match '^192\.168\.1\.'

        # Verify second network (10.50.1.0/24) subnets are present
        $managementSubnets = $results | Where-Object { $_.Subnet -match '^10\.50\.1\.' }
        $managementSubnets | Should -Not -BeNullOrEmpty
        
        $edgeMgmtRows = $results | Where-Object { $_.Name -eq 'csu-edge-transport-management' }
        $edgeMgmtRows | Should -Not -BeNullOrEmpty
        
        $msuMgmtRows = $results | Where-Object { $_.Name -eq 'msu-management' }
        $msuMgmtRows | Should -Not -BeNullOrEmpty

        # Verify results include both networks
        ($results | Where-Object { $_.Subnet -match '^192\.168\.1\.' }).Count | Should -BeGreaterThan 0
        ($results | Where-Object { $_.Subnet -match '^10\.50\.1\.' }).Count | Should -BeGreaterThan 0
    }

    It 'processes multiple networks with IPAssignments correctly' {
        $jsonConfig = @'
[
  {
    "network": "10.0.0.0/24",
    "subnets": [
      {
        "name": "Network1-Subnet1",
        "vlan": 100,
        "cidr": "28",
        "IPAssignments": [
          { "Name": "Gateway", "Position": 1 },
          { "Name": "Server1", "Position": 2 }
        ]
      }
    ]
  },
  {
    "network": "172.16.0.0/24",
    "subnets": [
      {
        "name": "Network2-Subnet1",
        "vlan": 200,
        "cidr": "28",
        "IPAssignments": [
          { "Name": "Gateway", "Position": 1 },
          { "Name": "Server2", "Position": 2 }
        ]
      }
    ]
  }
]
'@

        $resultsJson = New-SubnetPlanFromConfig -JsonConfig $jsonConfig -AsJson
        $results = $resultsJson | ConvertFrom-Json

        # Check Network1 assignments
        $net1Gateway = $results | Where-Object { $_.Name -eq 'Network1-Subnet1' -and $_.Label -eq 'Gateway' }
        $net1Gateway | Should -Not -BeNullOrEmpty
        $net1Gateway.IP | Should -Be '10.0.0.1'
        
        $net1Server = $results | Where-Object { $_.Name -eq 'Network1-Subnet1' -and $_.Label -eq 'Server1' }
        $net1Server | Should -Not -BeNullOrEmpty
        $net1Server.IP | Should -Be '10.0.0.2'

        # Check Network2 assignments
        $net2Gateway = $results | Where-Object { $_.Name -eq 'Network2-Subnet1' -and $_.Label -eq 'Gateway' }
        $net2Gateway | Should -Not -BeNullOrEmpty
        $net2Gateway.IP | Should -Be '172.16.0.1'
        
        $net2Server = $results | Where-Object { $_.Name -eq 'Network2-Subnet1' -and $_.Label -eq 'Server2' }
        $net2Server | Should -Not -BeNullOrEmpty
        $net2Server.IP | Should -Be '172.16.0.2'
    }
}
