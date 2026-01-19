package preset

import "testing"

func TestValidatePresetDataJSON(t *testing.T) {
	cases := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{name: "empty ok", input: "", wantErr: false},
		{name: "invalid json", input: "{", wantErr: true},
		{name: "non object", input: `[]`, wantErr: true},
		{name: "minimal ok", input: `{}`, wantErr: false},
		{name: "sections must be array", input: `{"sections":{}}`, wantErr: true},
		{name: "section items must be array", input: `{"sections":[{"items":{}}]}`, wantErr: true},
		{name: "section id must be string", input: `{"sections":[{"id":1}]}`, wantErr: true},
		{name: "ok with sections", input: `{"sections":[{"id":"exp","type":"Experience","items":[{"id":"1"}]}]}`, wantErr: false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidatePresetDataJSON(tc.input)
			if tc.wantErr && err == nil {
				t.Fatalf("expected error, got nil")
			}
			if !tc.wantErr && err != nil {
				t.Fatalf("expected nil, got %v", err)
			}
		})
	}
}

