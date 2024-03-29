### Switch Framework and Template

This part is the core of the project. Each switch need to have paired `framework` and `template` files to be able generate configuration accordingly.

### Switch Configuration Files

Switch configuration is generated by using Go native package: [text/template](https://pkg.go.dev/text/template)

#### Logic Diagram

```mermaid
flowchart
    A[Switch Output Object]
    B(allConfig.go.tmpl)
    C[Final Switch Configuration]
    D(header.go.tmpl)
    E(vlan.go.tmpl)
    F(bgp.go.tmpl)
    G(xxx.go.tmpl)

    A <-.-> |Parse| D
    A <-.-> |Parse| E
    A <-.-> |Parse| F
    A <-.-> |Parse| G
    D --> |Merge| B
    E --> |Merge| B
    F --> |Merge| B
    G --> |Merge| B
    B --> C
```

#### Template Structure

| Config     | Template           | Source                       |
| ---------- | ------------------ | ---------------------------- |
| All Config | allConfig.go.tmpl  | All templates below          |
| Header     | header.go.tmpl     | OutputObj.Device             |
| VLAN       | vlan.go.tmpl       | OutputObj.Vlan               |
| InBandPort | inBandPort.go.tmpl | OutputObj.Port               |
| BGP        | bgp.go.tmpl        | OutputObj.Routing.Router.Bgp |
