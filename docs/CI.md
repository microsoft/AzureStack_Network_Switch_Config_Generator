# Add Config Gen Tool to CI

## Overall Design Flow

### Nuget CI

```mermaid
flowchart TD

    C[(Switch Config Tool GitHub Repo)] --> |CI Pull Remote Repo| D[Build Tool.exe and save folder SwitchLib]
    D --> |CI Build| E[Sign and Package Nuget]
    E --> |CI Push| B[(Nuget Feed)]

```

### Switch Deploy CI

```mermaid
flowchart TD
    L[Lab Repo] --> |Copy to CI| A[New Switch Input.json]
    A[New Switch Input.json] --> |CI Integrate| B{{Invoke Generate Switch Config}}
    E[Pull the SwitchConfigGen Artifact] --> |CI Integrate| B
    B --> |CI Execute and Log| F[Switch Configuration Files]
    F --> |Copy to azlabslogs Server| G[For Debug Use]
    F --> |Copy to Lab TFTP| H[For Switch Deployment Use]
```

### MileStone1: Generate the Config in CI

- Define the Switch input files with Tools and manually generated to make sure the configuration is accurate.
- Figure out where to put the tool as well as input.json file on CI.
- Update the CI script to trigger the tool run

### MileStone2: Push the configuration to the Switch and validate
