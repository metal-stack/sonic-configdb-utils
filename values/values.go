package values

import "gopkg.in/yaml.v2"

type DockerRoutingConfigMode string

func (m DockerRoutingConfigMode) String() string {
	return string(m)
}

const (
	DockerRoutingConfigModeSeparated DockerRoutingConfigMode = "SEPARATED"
	DockerRoutingConfigModeSplit     DockerRoutingConfigMode = "SPLIT"
	DockerRoutingConfigModeUnified   DockerRoutingConfigMode = "UNIFIED"
)

type FECMode string

func (m FECMode) String() string {
	return string(m)
}

const (
	FECModeNone FECMode = "NONE"
	FECModeRS   FECMode = "RS"
)

type Interconnect struct {
	VNI int    `yaml:"vni"`
	VRF string `yaml:"vrf"`
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

type Portchannel struct {
	Number   string   `yaml:"number"`
	MTU      int      `yaml:"mtu"`
	Fallback bool     `yaml:"fallback"`
	Members  []string `yaml:"members"`
}

type SAG struct {
	MAC string `yaml:"mac"`
}

type Values struct {
	BGPPorts                []string          `yaml:"sonic_bgp_ports"`
	Breakouts               map[string]string `yaml:"sonic_breakouts"`
	DockerRoutingConfigMode `yaml:"sonic_docker_routing_config_mode"`
	FRRMgmtFrameworkConfig  bool                    `yaml:"sonic_frr_mgmt_framework_config"`
	Hostname                string                  `yaml:"sonic_hostname"`
	Interconnects           map[string]Interconnect `yaml:"sonic_interconnects"`
	LLDPHelloTimer          int                     `yaml:"sonic_lldp_hello_timer"`
	LoopbackAddress         string                  `yaml:"sonic_loopback_address"`
	MCLAG                   `yaml:"sonic_mclag"`
	MgmtIfGateway           string        `yaml:"sonic_mgmtif_gateway"`
	MgmtIfIP                string        `yaml:"sonic_mgmtif_ip"`
	MgmtVRF                 bool          `yaml:"sonic_mgmt_vrf"`
	NTPServers              []string      `yaml:"sonic_ntpservers"`
	Portchannels            []Portchannel `yaml:"sonic_portchannels"`
	PortchannelsDefaultMTU  int           `yaml:"sonic_portchannels_default_mtu"`
	Ports                   []Port        `yaml:"sonic_ports"`
	PortsDefaultFEC         FECMode       `yaml:"sonic_ports_default_fec"`
	PortsDefaultMTU         int           `yaml:"sonic_ports_default_mtu"`
	PortsDefaultSpeed       int           `yaml:"sonic_ports_default_speed"`
	SAG                     `yaml:"sonic_sag"`
	SSHSourceranges         []string `yaml:"sonic_ssh_sourceranges"`
	VLANs                   []VLAN   `yaml:"sonic_vlans"`
	VTEPs                   []VTEP   `yaml:"sonic_vteps"`
}

type VLAN struct {
	DHCPServers   []string `yaml:"dhcp_servers"`
	ID            string   `yaml:"id"`
	IP            string   `yaml:"ip"`
	TaggedPorts   []string `yaml:"tagged_ports"`
	UntaggedPorts []string `yaml:"untagged_ports"`
	VRF           string   `yaml:"vrf"`
}

type VTEP struct {
	Comment string `yaml:"comment"`
	VNI     int    `yaml:"vni"`
	VLAN    string `yaml:"vlan"`
}

func UnmarshalValues(input []byte) (*Values, error) {
	values := Values{}
	err := yaml.Unmarshal(input, &values)
	if err != nil {
		return nil, err
	}

	return &values, nil
}
