package values

import "gopkg.in/yaml.v3"

type AutonegMode string

const (
	AutonegModeOn  AutonegMode = "on"
	AutonegModeOff AutonegMode = "off"
)

type DockerRoutingConfigMode string

const (
	DockerRoutingConfigModeSeparated    DockerRoutingConfigMode = "separated"
	DockerRoutingConfigModeSplit        DockerRoutingConfigMode = "split"
	DockerRoutingConfigModeSplitUnified DockerRoutingConfigMode = "split-unified"
	DockerRoutingConfigModeUnified      DockerRoutingConfigMode = "unified"
)

type Feature struct {
	AutoRestart bool `yaml:"auto_restart"`
	Enabled     bool `yaml:"enabled"`
}

type FECMode string

const (
	FECModeNone FECMode = "none"
	FECModeRS   FECMode = "rs"
)

type Interconnect struct {
	UnnumberedInterfaces []string `yaml:"unnumbered_interfaces"`
	VNI                  string   `yaml:"vni"`
	VRF                  string   `yaml:"vrf"`
}

type MCLAG struct {
	KeepaliveVLAN      string   `yaml:"keepalive_vlan"`
	MemberPortChannels []string `yaml:"member_port_channels"`
	PeerIP             string   `yaml:"peer_ip"`
	PeerLink           string   `yaml:"peer_link"`
	SourceIP           string   `yaml:"source_ip"`
	SystemMAC          string   `yaml:"system_mac"`
}

type MgmtInterface struct {
	GatewayAddress string `yaml:"gateway_address"`
	IP             string `yaml:"ip"`
}

type NTP struct {
	Servers      []string `yaml:"servers"`
	SrcInterface string   `yaml:"src_interface"`
	VRF          string   `yaml:"vrf"`
}

type Port struct {
	Autoneg AutonegMode `yaml:"autoneg"`
	FECMode FECMode     `yaml:"fec"`
	IPs     []string    `yaml:"ips"`
	MTU     int         `yaml:"mtu"`
	Name    string      `yaml:"name"`
	Speed   int         `yaml:"speed"`
	VRF     string      `yaml:"vrf"`
}

type Ports struct {
	DefaultAutoneg AutonegMode `yaml:"default_autoneg"`
	DefaultFEC     FECMode     `yaml:"default_fec"`
	DefaultMTU     int         `yaml:"default_mtu"`
	List           []Port      `yaml:"list"`
}

type PortChannel struct {
	Number   string   `yaml:"number"`
	MTU      int      `yaml:"mtu"`
	Fallback bool     `yaml:"fallback"`
	Members  []string `yaml:"members"`
}

type PortChannels struct {
	DefaultMTU int           `yaml:"default_mtu"`
	List       []PortChannel `yaml:"list"`
}

type SAG struct {
	MAC string `yaml:"mac"`
}

type Values struct {
	BGPPorts                []string                `yaml:"bgp_ports"`
	Breakouts               map[string]string       `yaml:"breakouts"`
	DockerRoutingConfigMode DockerRoutingConfigMode `yaml:"docker_routing_config_mode"`
	Features                map[string]Feature      `yaml:"features"`
	FRRMgmtFrameworkConfig  bool                    `yaml:"frr_mgmt_framework_config"`
	Hostname                string                  `yaml:"hostname"`
	Interconnects           map[string]Interconnect `yaml:"interconnects"`
	LLDPHelloTime           int                     `yaml:"lldp_hello_time"`
	LoopbackAddress         string                  `yaml:"loopback_address"`
	MCLAG                   MCLAG                   `yaml:"mclag"`
	MgmtInterface           MgmtInterface           `yaml:"mgmt_interface"`
	MgmtVRF                 bool                    `yaml:"mgmt_vrf"`
	Nameservers             []string                `yaml:"nameservers"`
	NTP                     NTP                     `yaml:"ntp"`
	PortChannels            PortChannels            `yaml:"portchannels"`
	Ports                   Ports                   `yaml:"ports"`
	SAG                     *SAG                    `yaml:"sag"`
	SSHSourceranges         []string                `yaml:"ssh_sourceranges"`
	VLANs                   []VLAN                  `yaml:"vlans"`
	VLANSubinterfaces       []VLANSubinterface      `yaml:"vlan_subinterfaces"`
	VTEP                    VTEP                    `yaml:"vtep"`
}

type VLANSubinterface struct {
	CIDR string `yaml:"cidr"`
	Port string `yaml:"port"`
	VLAN string `yaml:"vlan"`
	VRF  string `yaml:"vrf"`
}

type VLAN struct {
	DHCPServers   []string `yaml:"dhcp_servers"`
	ID            string   `yaml:"id"`
	IP            string   `yaml:"ip"`
	SAG           *bool    `yaml:"sag"`
	TaggedPorts   []string `yaml:"tagged_ports"`
	UntaggedPorts []string `yaml:"untagged_ports"`
	VRF           string   `yaml:"vrf"`
	VRRP          VRRP     `yaml:"vrrp"`
}

type VRRP struct {
	Group string `yaml:"group"`
	IP    string `yaml:"ip"`
}

type VTEP struct {
	Enabled         bool             `yaml:"enabled"`
	VXLANTunnelMaps []VXLANTunnelMap `yaml:"vxlan_tunnel_maps"`
}

type VXLANTunnelMap struct {
	VNI  string `yaml:"vni"`
	VLAN string `yaml:"vlan"`
}

func UnmarshalValues(in []byte) (*Values, error) {
	var values Values
	err := yaml.Unmarshal(in, &values)
	if err != nil {
		return nil, err
	}

	return &values, nil
}
