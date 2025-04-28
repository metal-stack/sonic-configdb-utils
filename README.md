# SONiC config_db Utils

Commandline tool for generating a `config_db.json` file for SONiC.

## Usage

Create a `sonic-config.yaml` file with config variables.
Then run

```bash
sonic-confidb-utils generate
```

## Configuration Parameters

**bgp_ports**

Example:

```yaml
bgp_ports:
  - Ethernet120
  - Ethernet124
```

Result:

```json
{
  "INTERFACE": {
    "Ethernet120": {
      "ipv6_use_link_local_only": "enable"
    },
    "Ethernet124": {
      "ipv6_use_link_local_only": "enable"
    }
  }
}
```

**breakouts**

Example:

```yaml
breakouts:
  Ethernet0: 4x25G
  Ethernet120: 1x100G[40G]
```

Result:

```json
{
  "BREAKOUT_CFG": {
    "Ethernet0": {
      "brkout_mode": "4x25G"
    },
    "Ethernet120": {
      "brkout_mode": "100G[40G]"
    }
  }
}
```

For each breakout also the correspondig ports entries are added.

**docker_routing_config_mode**

Can be one of `separated`, `split`, `split-unified`, `unified`.

Example:

```yaml
docker_routing_config_mode: split
```

Result:

```json
{
  "DEVICE_METADATA": {
    "localhost": {
      "docker_routing_config_mode": "split"
    }
  }
}
```

**frr_mgmt_framework_config**

Example:

```yaml
frr_mgmt_framework_config: false
```

Result:

```json
{
  "DEVICE_METADATA": {
    "localhost": {
      "frr_mgmt_framework_config": "false"
    }
  }
}
```

**interconnects**

Example:

```yaml
interconnects:
  mpls:
    vni: 104010
    vrf: VrfMpls
```
