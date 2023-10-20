package common

import (
	"encoding/json"
	"reflect"
	"testing"
)

func TestParseDoc(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *Doc
		wantErr bool
	}{
		{
			name: "doc",
			args: args{
				s: `
asd(fas)dfsa
fsadfa, sdfsadf
fasdf"safd
@say hello work
@an(/123/456)
@enum("hello:world","foo:bar","say:if", "hello:fsafsadf", boweian,
"num:1")
asdfafasdfasdfasd
 @copy("fsdafasf:fdsafas", "fsdfa:fasdfa")
dsafs"dafasdfasd
asdfsa"dfsadf
`,
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "null",
			args: args{
				s: "",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDoc(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDoc() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ParseDoc() got = %v, want %v", got, tt.want)
				b, _ := json.Marshal(got)
				t.Log(string(b))
			}
		})
	}
}
