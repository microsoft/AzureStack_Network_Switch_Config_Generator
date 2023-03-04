# Azure Stack Switch Configuration Generator

## Project Overview

### Background

This is a tool to generate network switch deployment configuration for Azure Stack, which:

- Offers Network Switch Deployment Automation to Deployment Engineer
- Supports Multiple ASZ Network Design Use Cases by Customized Input Template Variables.

### Workflow

```mermaid
flowchart TD
    A[SwitchLib: Framework JSON + Go Template]
    B[User Input Template]
    C(Generator Tool)
    D(Switch Output Object)
    E[Switch Object JSON Files]
    F[Switch Configuration Files]
    B --> C
    A --> C
    C --> D
    D -.-> |For Debug| E
    D --> |For Deploy| F
```

## Project Design

- [User Input JSON](docs/User_Input_Json.md)

- [Customized SwitchLib](docs/Customized_SwitchLib.md)

## Get Start

- [Quick Start](docs/Quick_Start.md)

### Preparation

### Use Cases

- [Generate Vlan Configuration](docs/Generate_Vlan_Config.md)
