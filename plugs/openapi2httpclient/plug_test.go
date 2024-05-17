package openapi2httpclient

import "testing"

func TestLoad(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			name: "TestLoad",
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Load()
		})
	}
}
