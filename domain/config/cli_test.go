package config

import (
	`encoding/json`
	`testing`
)

/*
	Created by aomerk at 5/20/21 for project strixeye
*/

/*
	INSERT FILE DESCRIPTION HERE
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

func TestCli_UnmarshalJSON(t *testing.T) {
	type fields struct {
		userAPIToken   string
		currentAgentId string
	}
	type args struct {
		bytes []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "testing bad cli object",
			fields: fields{
				userAPIToken:   "bad_token",
				currentAgentId: "bad_id",
			},
			args: args{
				bytes: []byte(`{"user_api_token": "bad_token", "current_agent_id": "bad_id"}`),
			},
			wantErr: false,
		},
		{
			name: "testing worse cli object",
			fields: fields{
				userAPIToken:   "bad_token",
				currentAgentId: "bad_id",
			},
			args: args{
				bytes: []byte(`"user_api_": "bad_token", "current_agent_id": "bad_id"}`),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &Cli{}
				if err := json.Unmarshal(tt.args.bytes, c); (err != nil) != tt.wantErr {
					t.Errorf("UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				}
				
			},
		)
	}
}

func TestCli_Validate(t *testing.T) {
	type fields struct {
		userAPIToken   string
		currentAgentId string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "Empty field exists",
			fields: fields{
				userAPIToken: "blah",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				c := &Cli{
					UserAPIToken:   tt.fields.userAPIToken,
					CurrentAgentId: tt.fields.currentAgentId,
				}
				if err := c.Validate(); (err != nil) != tt.wantErr {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
