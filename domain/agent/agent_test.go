package agent

import (
	"testing"
	"time"
)

/*
	Created by aomerk at 6/14/21 for project cli
*/

/*
	INSERT FILE DESCRIPTION HERE
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

func Test_decode(t *testing.T) {
	type args struct {
		s APIVersionsMessage
		i interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "success",
			args: args{
				i: &Versions{},
				s: APIVersionsMessage{
					Data: []struct {
						Key       string      `json:"key"`
						Value     *version    `json:"value"`
						CreatedAt interface{} `json:"created_at"`
						UpdatedAt *time.Time  `json:"updated_at"`
						DeletedAt interface{} `json:"deleted_at"`
					}{
						{
							Key: "Manager", Value: &version{
								Version: "test-version", Size: 12345,
								Checksum: "test-checksum",
							},
						}, {
							Key: "Installer", Value: &version{
								Version: "test-installer-version", Size: 54322,
								Checksum: "test-checksum",
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				if err := decode(tt.args.s, tt.args.i); (err != nil) != tt.wantErr {
					t.Errorf("decode() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
