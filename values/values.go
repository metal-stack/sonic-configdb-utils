package values

import "gopkg.in/yaml.v3"

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

type Port struct {
	IPs     []string `yaml:"ips"`
	FECMode `yaml:"fec"`
	MTU     int    `yaml:"mtu"`
	Name    string `yaml:"name"`
	Speed   int    `yaml:"speed"`
	VRF     string `yaml:"vrf"`
}

type PortChannel struct {
	Number   string   `yaml:"number"`
	MTU      int      `yaml:"mtu"`
	Fallback bool     `yaml:"fallback"`
	Members  []string `yaml:"members"`
}

type SAG struct {
	MAC string `yaml:"mac"`
}

type Values struct {
	BGPPorts                []string          `yaml:"bgp_ports"`
	Breakouts               map[string]string `yaml:"breakouts"`
	DockerRoutingConfigMode `yaml:"docker_routing_config_mode"`
	Features                map[string]Feature      `yaml:"features"`
	FRRMgmtFrameworkConfig  bool                    `yaml:"frr_mgmt_framework_config"`
	Hostname                string                  `yaml:"hostname"`
	Interconnects           map[string]Interconnect `yaml:"interconnects"`
	LLDPHelloTime           int                     `yaml:"lldp_hello_time"`
	LoopbackAddress         string                  `yaml:"loopback_address"`
	MCLAG                   MCLAG                   `yaml:"mclag"`
	MgmtIfGateway           string                  `yaml:"mgmtif_gateway"`
	MgmtIfIP                string                  `yaml:"mgmtif_ip"`
	MgmtVRF                 bool                    `yaml:"mgmt_vrf"`
	Nameservers             []string                `yaml:"nameservers"`
	NTPServers              []string                `yaml:"ntpservers"`
	PortChannels            []PortChannel           `yaml:"portchannels"`
	PortChannelsDefaultMTU  int                     `yaml:"portchannels_default_mtu"`
	Ports                   []Port                  `yaml:"ports"`
	PortsDefaultFEC         FECMode                 `yaml:"ports_default_fec"`
	PortsDefaultMTU         int                     `yaml:"ports_default_mtu"`
	SAG                     `yaml:"sag"`
	SSHSourceranges         []string `yaml:"ssh_sourceranges"`
	VLANMembers             bool     `yaml:"vlan_members"`
	VLANs                   []VLAN   `yaml:"vlans"`
	VTEPs                   []VTEP   `yaml:"vteps"`
}

type VLAN struct {
	DHCPServers   []string `yaml:"dhcp_servers"`
	ID            string   `yaml:"id"`
	IP            string   `yaml:"ip"`
	SAG           bool     `yaml:"sag"`
	TaggedPorts   []string `yaml:"tagged_ports"`
	UntaggedPorts []string `yaml:"untagged_ports"`
	VRF           string   `yaml:"vrf"`
}

type VTEP struct {
	Comment string `yaml:"comment"`
	VNI     string `yaml:"vni"`
	VLAN    string `yaml:"vlan"`
}

func UnmarshalValues(in []byte) (*Values, error) {
	var values Values
	err := yaml.Unmarshal(in, &values)
	if err != nil {
		return nil, err
	}

	return &values, nil
}
