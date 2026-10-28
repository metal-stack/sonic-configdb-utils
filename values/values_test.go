package values

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_UnmarshalValues_acl(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    ACL
		wantErr string
	}{
		{
			name: "known fields, unquoted priority",
			input: `
acl:
  ACL_TABLE:
    PXE_IN:
      policy_desc: pxe
      type: L3
      stage: ingress
      ports: [Vlan4000]
  ACL_RULE:
    PXE_IN|DENY_1:
      PRIORITY: 7999
      PACKET_ACTION: DROP
      DST_IP: 10.0.0.0/16
`,
			want: ACL{
				Tables: map[string]ACLTable{"PXE_IN": {PolicyDesc: "pxe", Ports: []string{"Vlan4000"}, Stage: "ingress", Type: "L3"}},
				Rules:  map[string]ACLRule{"PXE_IN|DENY_1": {DstIP: "10.0.0.0/16", PacketAction: "DROP", Priority: "7999"}},
			},
		},
		{
			name: "unknown rule field",
			input: `
acl:
  ACL_RULE:
    PXE_IN|DENY_1:
      PRIORITY: "7999"
      DST_IPS: 10.0.0.0/16
`,
			wantErr: `unknown field "DST_IPS" in acl rule, use one of DST_IP, ETHER_TYPE, IP_TYPE, PACKET_ACTION, PRIORITY, SRC_IP`,
		},
		{
			name: "unknown table field",
			input: `
acl:
  ACL_TABLE:
    PXE_IN:
      type: L3
      port: [Vlan4000]
`,
			wantErr: `unknown field "port" in acl table, use one of policy_desc, ports, services, stage, type`,
		},
		{
			name: "unknown acl section",
			input: `
acl:
  ACL_RULES: {}
`,
			wantErr: `unknown field "ACL_RULES" in acl, use one of ACL_RULE, ACL_TABLE`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := UnmarshalValues([]byte(tt.input))
			if tt.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
					t.Fatalf("UnmarshalValues() error = %v, want it to contain %q", err, tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("UnmarshalValues() unexpected error = %v", err)
			}
			if diff := cmp.Diff(tt.want, got.ACL); diff != "" {
				t.Errorf("UnmarshalValues() acl diff = %s", diff)
			}
		})
	}
}
