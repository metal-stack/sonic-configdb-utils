package configdb

type ACLRule struct {
	EtherType    string `json:"ETHER_TYPE"`
	PacketAction `json:"PACKET_ACTION"`
	Priority     string `json:"PRIORITY"`
	SrcIP        string `json:"SRC_IP"`
}

type ACLTable struct {
	PolicyDesc string   `json:"policy_desc"`
	Ports      []string `json:"ports"`
	Services   []string `json:"services"`
	Stage      string   `json:"stage"`
	Type       string   `json:"type"`
}

type AdminStatus string

func (s AdminStatus) String() string {
	return string(s)
}

const (
	AdminStatusUp   AdminStatus = "UP"
	AdminStatusDown AdminStatus = "DOWN"
)

type AutonegMode string

func (m AutonegMode) String() string {
	return string(m)
}

const (
	AutonegModeOn  AutonegMode = "ON"
	AutonegModeOff AutonegMode = "OFF"
)

type BreakoutConfig struct {
	BreakoutMode `json:"brkout_mode"`
}

type BreakoutMode string

func (m BreakoutMode) String() string {
	return string(m)
}

const (
	BreakoutMode1x100G BreakoutMode = "1x100G[40G]"
	BreakoutMode4x25G  BreakoutMode = "4x25G"
	BreakoutMode4x10G  BreakoutMode = "4x10G"
)

type DeviceMetadata struct {
	Localhost Metadata `json:"localhost"`
}

type DockerRoutingConfigMode string

func (m DockerRoutingConfigMode) String() string {
	return string(m)
}

const (
	DockerRoutingConfigModeSeparated DockerRoutingConfigMode = "SEPARATED"
	DockerRoutingConfigModeSplit     DockerRoutingConfigMode = "SPLIT"
	DockerRoutingConfigModeUnified   DockerRoutingConfigMode = "UNIFIED"
)

type Feature struct {
	AutoRestart FeatureMode `json:"auto_restart"`
	State       FeatureMode `json:"state"`
}

type FeatureMode string

func (m FeatureMode) String() string {
	return string(m)
}

const (
	FeatureModeEnabled  FeatureMode = "ENABLED"
	FeatureModeDisabled FeatureMode = "DISABLED"
)

type FECMode string

func (m FECMode) String() string {
	return string(m)
}

const (
	FECModeNone FECMode = "NONE"
	FECModeRS   FECMode = "RS"
)

type LLDP struct {
	Global LLDPGlobal `json:"GLOBAL"`
}

type LLDPGlobal struct {
	HelloTime int `json:"hello_time"`
}

type MgmtInterface struct {
	GWAddr string `json:"gwaddr"`
}

type MgmtPort struct {
	AdminStatus `json:"admin_status"`
	Alias       string `json:"alias"`
	Description string `json:"description"`
}

type MgmtVRFConfig struct {
	VRFGlobal `json:"vrf_global"`
}

type Metadata struct {
	DockerRoutingConfigMode `json:"docker_routing_config_mode"`
	FRRMgmtFrameworkConfig  bool   `json:"frr_mgmt_framework_config"`
	Hostname                string `json:"hostname"`
	RouterType              `json:"type"`
}

type NTP struct {
	NTPGlobal `json:"global"`
}

type NTPGlobal struct {
	SrcIntf string `json:"src_intf"`
}

type PacketAction string

func (pa PacketAction) String() string {
	return string(pa)
}

const (
	PacketActionDrop   PacketAction = "DROP"
	PacketActionAccept PacketAction = "ACCEPT"
)

type Port struct {
	AdminStatus `json:"admin_status"`
	Alias       string      `json:"alias"`
	Autoneg     AutonegMode `json:"autoneg"`
	FEC         FECMode     `json:"fec"`
	Index       int         `json:"index"`
	Lanes       string      `json:"lanes"`
	MTU         int         `json:"mtu"`
	ParentPort  string      `json:"parent_port"`
	Speed       int         `json:"speed"`
}

type RouterType string

func (t RouterType) String() string {
	return string(t)
}

const (
	RouterTypeLeafRouter RouterType = "LEAF_ROUTER"
	RouterTypeDualToR    RouterType = "DUAL_TOR"
)

type VLAN struct {
	DHCPServers []string `json:"dhcp_servers"`
	VLANID      int      `json:"vlanid"`
}

type VRFGlobal struct {
	MgmtVRFEnabled bool `json:"mgmtVrfEnabled"`
}

type VXLANEVPN struct {
	VXLANEVPNNVO `json:"nvo"`
}

type VXLANEVPNNVO struct {
	SourceVTEP string `json:"source_vtep"`
}

type VXLANTunnel struct {
	SrcIP string `json:"src_ip"`
}

type VXLANTunnelMap struct {
	VLAN string `json:"vlan"`
	VNI  int    `json:"vni"`
}