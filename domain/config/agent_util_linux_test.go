// +build linux

package config

import `testing`

/*
	Created by aomerk at 5/23/21 for project cli
*/

/*
	INSERT FILE DESCRIPTION HERE
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

func TestAgentInformation_CheckIfHostSupports(t *testing.T) {
	type fields struct {
		ID        string
		CompanyID string
		Name      string
		IPAddress string
		Config    stackConfig
		Domains   []domains
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Linux Docker Compose",
			fields: fields{
				Config: stackConfig{Deployment: "docker"},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				a := AgentInformation{
					ID:        tt.fields.ID,
					CompanyID: tt.fields.CompanyID,
					Name:      tt.fields.Name,
					IPAddress: tt.fields.IPAddress,
					Config:    tt.fields.Config,
					Domains:   tt.fields.Domains,
				}
				if err := a.CheckIfHostSupports(); (err != nil) != tt.wantErr {
					t.Errorf("CheckIfHostSupports() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
