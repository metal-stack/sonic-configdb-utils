# SONiC config_db Utils

Commandline tool for generating a `config_db.json` file for SONiC.

## Usage

Create a `sonic-config.yaml` file with config variables.
Then run

```bash
sonic-confidb-utils generate
```

Following files and directories are expected to be in place:

1. The `/usr/share/sonic/device` directory holds all device specific information. `sonic-configdb-utils` needs to read `usr/share/sonic/device/<platform-identifier>/platform.json` to set and validate the ports' breakout configurations.
2. To retrieve the platform identifier and HWSKU the `/etc/sonic/sonic-environment` file is read which looks like this:

```bash
SONIC_VERSION=sonic-123
PLATFORM=x86_64-accton_as7726_32x-r0
HWSKU=Accton-AS7726-32X
DEVICE_TYPE=LeafRouter
ASIC_TYPE=broadcom
```

## Configuration Parameters

### bgp_ports

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

### breakouts

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

### docker_routing_config_mode

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

### features

Enabling a feature makes SONiC check the status of the corresponding container.
If a feature is enabled and the container stops running `container_checker` will complain and `show system-health summary` will show

```
Status: Not OK
  Reasons: Container 'your-feature' is not running
```

Furthermore, enabling `auto_restart` on a service will tell SONiC to automatically restart the container if it stops unexpectedly.

Example:

```yaml
features:
  metal-core.service:
    enabled: true
    auto_restart: true
```

Result:

```json
{
  "FEATURE": {
    "metal-core.service": {
      "auto_restart": "enabled",
      "state": "enabled"
    }
  }
}
```

### frr_mgmt_framework_config

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

### interconnects

Example:

```yaml
interconnects:
  mpls:
    vni: 104010
    vrf: VrfMpls
    unnumbered_interfaces:
      - Ethernet0
      - Ethernet1
```

Result:

```json
{
  "INTERFACE": {
    "Ethernet0": {
      "ipv6_use_link_local_only": "enable",
      "vrf_name": "VrfMpls"
    },
    "Ethernet1": {
      "ipv6_use_link_local_only": "enable",
      "vrf_name": "VrfMpls"
    }
  }
}
```

### lldp_hello_time

This configuration slightly differs across SONiC versions and might not always have an effect.
As this tool is tailored for Edgecore SONiC versions 202111.x the LLDP configuration for those versions should work.

Example:

```yaml
lldp_hello_time: 10
```

Result:

```json
{
  "LLDP": {
    "Global": {
      "hello_timer": "10"
    }
  }
}
```

### looback_address

Example:

```yaml
loopback_address: 10.7.7.7
```

Result:

```json
{
  "LOOPBACK_INTERFACE": {
    "Loopback0": {},
    "Loopback0|10.7.7.7/32": {}
  }
}
```

### mclag

Example:

```yaml
mclag:
  keepalive_vlan: 1000
  member_port_channels:
    - 11
    - 22
  peer_ip: 192.168.255.1
  peer_link: PortChannel01
  source_ip: 192.168.255.2
  system_mac: aa:aa:aa:aa:aa:aa
```

Result:

```json
{
  "MCLAG_DOMAIN": {
    "1": {
      "mclag_system_id": "aa:aa:aa:aa:aa:aa",
      "peer_ip": "192.168.255.1",
      "peer_link": "PortChannel01",
      "source_ip": "192.168.255.2"
    }
  },
  "MCLAG_INTERFACE": {
    "1|PortChannel11": {
      "if_type": "PortChannel"
    },
    "1|PortChannel22": {
      "if_type": "PortChannel"
    }
  },
  "MCLAG_UNIQUE_IP": {
    "Vlan1000": {
      "unique_ip": "enable"
    }
  }
}
```

### mgmt_interface

Example:

```yaml
mgmt_interface:
  gateway_address: 10.7.10.1
  ip: 10.7.10.2
```

Result:

```json
{
  "MGMT_INTERFACE": {
    "eth0|10.7.10.2": {
      "gwaddr": "10.7.10.1"
    }
  }
}
```

### mgmt_vrf

Example:

```yaml
mgmt_vrf: false
```

Result:

```json
{
  "MGMT_VRF_CONFIG": {
    "vrf_global": {
      "mgmtVrfEnabled": "false"
    }
  }
}
```

### nameservers

Example:

```yaml
nameservers:
  - 1.1.1.1
  - 8.8.8.8
```

Result:

```json
{
  "DNS_NAMESERVER": {
    "1.1.1.1": {},
    "8.8.8.8": {}
  }
}
```

### ntp

Example:

```yaml
ntp:
  src_interface: Loopback0
  vrf: default
  servers:
    - 0.europe.pool.ntp.org
    - 1.europe.pool.ntp.org
    - 2.europe.pool.ntp.org
    - 3.europe.pool.ntp.org
```

Result:

```json
{
  "NTP": {
    "global": {
      "src_intf": "Loopback0",
      "vrf": "default"
    }
  },
  "NTP_SERVER": {
    "0.europe.pool.ntp.org": {},
    "1.europe.pool.ntp.org": {},
    "2.europe.pool.ntp.org": {},
    "3.europe.pool.ntp.org": {}
  }
}
```

### portchannels

Example:

```yaml
portchannels:
  default_mtu: 9000
  list:
    - number: 01
      mtu: 1500
      fallback: true
      members:
        - Ethernet4
        - Ethernet5
```

Result:

```json
{
  "PORTCHANNEL": {
    "PortChannel01": {
      "admin_status": "up",
      "fallback": "true",
      "fast_rate": "false",
      "lacp_key": "auto",
      "min_links": "1",
      "mix_speed": "false",
      "mtu": "1500"
    }
  },
  "PORTCHANNEL_MEMBER": {
    "PortChannel01|Ethernet4": {},
    "PortChannel01|Ethernet5": {}
  }
}
```

### ports

Example:

```yaml
breakouts:
  Ethernet0: 4x25G

ports:
  default_fec: none
  default_mtu: 9000
  default_autoneg: on
  list:
    - name: Ethernet0
      ips:
        - 10.4.3.2
      fec: rs
      mtu: 1500
      vrf: VrfMpls
```

Result:

```json
{
  "PORT": {
    "Ethernet0": {
      "admin_status": "up",
      "alias": "Eth1/1(Port1)",
      "autoneg": "on",
      "fec": "rs",
      "index": "1",
      "lanes": "1",
      "mtu": "1500",
      "speed": "25000"
    }
  }
}
```

and

```json
{
  "INTERFACE": {
    "Ethernet0": {
      "vrf_name": "VrfMpls"
    },
    "Ethernet0|10.4.3.2": {}
  }
}
```

The speed of a port is determined by its breakout configuration.
If no breakout configuration for a port is provided its default breakout is assumed.
If a breakout configuration allows more than one speed options, e.g. `1x100G[40G]`, the first speed option is used as a default (`100G` in the example).
A `speed` option can be added to the port config to specify an alternative speed, e.g.

```yaml
ports:
  - name: Ethernet120
    speed: 40000
```

For each port that is not explicitly configured in `breakouts` and `ports` an entry with defaults will be added to the `"PORT"` dictionary.

### sag

Example:

```yaml
sag:
  mac: bb:bb:bb:bb:bb:bb
```

Result:

```json
{
  "SAG": {
    "GLOBAL": {
      "gateway_mac": "bb:bb:bb:bb:bb:bb"
    }
  }
}
```

### ssh_sourceranges

Example:

```yaml
ssh_sourceranges:
  - 10.1.23.1/30
```

Result:

```json
{
  "ACL_RULE": {
    "ALLOW_NTP|DEFAULT_RULE": {
      "ETHER_TYPE": "2048",
      "PACKET_ACTION": "DROP",
      "PRIORITY": "1"
    },
    "ALLOW_NTP|RULE_1": {
      "PACKET_ACTION": "ACCEPT",
      "PRIORITY": "99",
      "SRC_IP": "0.0.0.0/0"
    },
    "ALLOW_SSH|DEFAULT_RULE": {
      "ETHER_TYPE": "2048",
      "PACKET_ACTION": "DROP",
      "PRIORITY": "1"
    },
    "ALLOW_SSH|RULE_1": {
      "PACKET_ACTION": "ACCEPT",
      "PRIORITY": "91",
      "SRC_IP": "10.1.23.1/30"
    }
  },
  "ACL_TABLE": {
    "ALLOW_NTP": {
      "policy_desc": "Allow NTP",
      "services": ["NTP"],
      "stage": "ingress",
      "type": "CTRLPLANE"
    },
    "ALLOW_SSH": {
      "policy_desc": "Allow SSH access",
      "services": ["SSH"],
      "stage": "ingress",
      "type": "CTRLPLANE"
    }
  }
}
```

### vlan_subinterfaces

Example:

```yaml
vlan_subinterfaces:
  - cidr: 1.2.3.0/24
    port: Ethernet0
    vlan: 1000
    vrf: Vrf45
```

Result:

```json
{
  "VLAN_SUB_INTERFACE": {
    "Ethernet0.1000": {
      "admin_status": "up",
      "vrf_name": "Vrf45"
    },
    "Ethernet0.1000|1.2.3.0/24": {}
  }
}
```

### vlans

Example:

```yaml
vlans:
  - id: 4000
    dhcp_servers:
      - 10.9.8.7
      - 10.9.8.6
    ip: 10.9.7.0
    sag: true
    tagged_ports:
      - PortChannel01
    untagged_ports:
      - PortChannel11
    vrf: Vrf45
    vrrp:
      group: 1
      priority: 66
      ip: 10.255.1.1/24
```

Result:

```json
{
  "VLAN": {
    "Vlan4000": {
      "dhcp_servers": ["10.9.8.7", "10.9.8.6"],
      "vlanid": "4000"
    }
  },
  "VLAN_INTERFACE": {
    "Vlan4000": {
      "static_anycast_gateway": "true",
      "vrf_name": "Vrf45"
    },
    "Vlan4000|10.9.7.0": {}
  },
  "VLAN_MEMBER": {
    "Vlan4000|PortChannel01": {
      "tagging_mode": "tagged"
    },
    "Vlan4000|PortChannel11": {
      "tagging_mode": "untagged"
    }
  },
  "VRRP_INTERFACE": {
      "Vrrp1-v4": {
          "parent_interface": "Vlan4000"
      },
      "Vrrp1-v4|10.255.1.1/24": {}
  }
}
```

### vtep

Example:

```yaml
loopback_address: 10.7.7.7

vtep:
  vxlan_tunnel_maps:
    - vni: 103999
      vlan: Vlan3999
```

Result:

```json
{
  "VXLAN_EVPN_NVO": {
    "nvo": {
      "source_vtep": "vtep"
    }
  },
  "VXLAN_TUNNEL": {
    "vtep": {
      "src_ip": "10.7.7.7"
    }
  },
  "VXLAN_TUNNEL_MAP": {
    "vtep|map_103999_Vlan3999": {
      "vlan": "Vlan3999",
      "vni": "103999"
    }
  }
}
```

If only `VXLAN_EVPN_NVO` and `VXLAN_TUNNEL` are needed with no tunnel maps:

```yaml
loopback_address: 10.7.7.7

vtep:
  enabled: true
```

Result:

```json
{
  "VXLAN_EVPN_NVO": {
    "nvo": {
      "source_vtep": "vtep"
    }
  },
  "VXLAN_TUNNEL": {
    "vtep": {
      "src_ip": "10.7.7.7"
    }
  }
}
```
