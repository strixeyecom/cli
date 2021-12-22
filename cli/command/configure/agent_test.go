package configure

import (
	"testing"

	"github.com/strixeyecom/cli/domain/agent"
)

/*
	Created by aomerk at 5/24/21 for project cli
*/

/*
	INSERT FILE DESCRIPTION HERE
*/

// global constants for file
const ()

// global variables (not cool) for this file
var ()

func Test_selectAgent(t *testing.T) {
	type args struct {
		agents []agent.AgentInformation
	}
	tests := []struct {
		name    string
		args    args
		want    agent.AgentInformation
		wantErr bool
	}{
		{
			name: "Test From 2 agents",
			args: args{
				agents: []agent.AgentInformation{
					{
						Name: "MyLittleMechanicalOwl",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				// This is a UI/UX test, so I am disabling it for now.
				// _, err := selectAgent(tt.args.agents)
				// if (err != nil) != tt.wantErr {
				// 	t.Errorf("selectAgent() error = %v, wantErr %v", err, tt.wantErr)
				// 	return
				// }
			},
		)
	}
}
