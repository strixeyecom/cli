package cli

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/spf13/viper"

	"github.com/strixeyecom/cli/domain/repository"
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

func TestCli_Load(t *testing.T) {
	var (
		cliConfig2 Cli
		err        error
	)

	// get good keys
	viper.SetConfigFile("../../.env")
	viper.AutomaticEnv()

	// Try to read from file, but use env variables if non exists. it's fine
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatal(err)
	}
	err = viper.Unmarshal(&cliConfig2)
	if err != nil {
		t.Fatalf("unable to decode into map, %v", err)
	}

	err = cliConfig2.Validate()
	if err != nil {
		t.Fatalf("test failed while validating cli config %s", err)
	}
}

func TestCli_UnmarshalJSON(t *testing.T) {
	type fields struct {
		userAPIToken string
		agentID      string
		DBConfig     repository.Database
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
				userAPIToken: "bad_token",
				agentID:      "bad_id",
				DBConfig:     repository.Database{},
			},
			args: args{
				bytes: []byte(`{"user_api_token": "bad_token", "agent_id": "bad_id"}`),
			},
			wantErr: false,
		},
		{
			name: "testing worse cli object",
			fields: fields{
				userAPIToken: "bad_token",
				agentID:      "bad_id",
			},
			args: args{
				bytes: []byte(`"user_api_": "bad_token", "agent_id": "bad_id"}`),
			},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]

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

const testConfig = `
agent_id: ""
api_domain: api.strixeye.com
db_addr: 127.0.0.1
db_name: strixeye
db_override: true
db_pass: "sdsss"
db_port: ""
db_user: ""
docker_registry: docker.strixeye.com
download_domain: downloads.strixeye.com
pretty_output: false
user_api_token: ""`

func TestCli_MarshalJSON(t *testing.T) {
	viper.SetConfigType("yaml")
	_ = viper.ReadConfig(strings.NewReader(testConfig))
	a := Cli{}
	err := viper.Unmarshal(&a)
	if err != nil {
		t.Error(err)
		return
	}

}

func TestCli_Validate(t *testing.T) {
	type fields struct {
		userAPIToken string
		agentID      string
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
		}, {
			name: "Good object",
			fields: fields{
				userAPIToken: "blah",
				agentID:      "blah",
			},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]

		t.Run(
			tt.name, func(t *testing.T) {
				c := &Cli{
					UserAPIToken: tt.fields.userAPIToken,
					AgentID:      tt.fields.agentID,
				}
				if err := c.Validate(); (err != nil) != tt.wantErr {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestCli_Save(t *testing.T) {
	type fields struct {
		UserAPIToken string
		agentID      string
	}
	type args struct {
		filePath string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:    "empty stack config",
			args:    args{filePath: "/tmp/stack-config-tmp.json"},
			wantErr: false,
		}, {
			name:    "bad path",
			args:    args{filePath: "/root/bad/"},
			wantErr: true,
		},
	}
	for i := range tests {
		tt := tests[i]

		t.Run(
			tt.name, func(t *testing.T) {
				c := &Cli{
					UserAPIToken: tt.fields.UserAPIToken,
					AgentID:      tt.fields.agentID,
				}
				if err := c.Save(tt.args.filePath); (err != nil) != tt.wantErr {
					t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
