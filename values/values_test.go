package values

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func Test_UnmarshalValues_acls(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []ACL
		wantErr string
	}{
		{
			name: "known fields",
			input: `
acls:
  - name: PXE_IN
    description: pxe
    ports: [Vlan4000]
    rules:
      - action: drop
        dst: 10.0.0.0/16
      - action: forward
`,
			want: []ACL{{
				Name:        "PXE_IN",
				Description: "pxe",
				Ports:       []string{"Vlan4000"},
				Rules:       []ACLRule{{Action: "drop", Dst: "10.0.0.0/16"}, {Action: "forward"}},
			}},
		},
		{
			name: "unknown rule field",
			input: `
acls:
  - name: PXE_IN
    rules:
      - action: drop
        dst_ip: 10.0.0.0/16
`,
			wantErr: `unknown field "dst_ip" in acl rule, use one of action, dst, src`,
		},
		{
			name: "fields from a merge key",
			input: `
common: &common
  ports: [Vlan4000]
acls:
  - <<: *common
    name: PXE_IN
    rules:
      - action: drop
`,
			want: []ACL{{
				Name:  "PXE_IN",
				Ports: []string{"Vlan4000"},
				Rules: []ACLRule{{Action: "drop"}},
			}},
		},
		{
			name: "unknown field behind a merge key",
			input: `
common: &common
  port: [Vlan4000]
acls:
  - <<: *common
    name: PXE_IN
`,
			wantErr: `line 3: unknown field "port" in acl, use one of description, name, ports, rules, stage`,
		},
		{
			name: "unknown acl field",
			input: `
acls:
  - name: PXE_IN
    port: [Vlan4000]
`,
			wantErr: `unknown field "port" in acl, use one of description, name, ports, rules, stage`,
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
			if diff := cmp.Diff(tt.want, got.ACLs); diff != "" {
				t.Errorf("UnmarshalValues() acls diff = %s", diff)
			}
		})
	}
}
