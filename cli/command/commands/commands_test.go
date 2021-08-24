package commands

import (
	`testing`
	
	`github.com/spf13/viper`
)

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

func TestNewStrixeyeCommand(t *testing.T) {
	tests := []struct {
		name string
	}{
		{
			"basic configuration test",
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				got := NewStrixeyeCommand()
				_ = got
				if viper.Get("API_DOMAIN") == "" {
					t.Fatalf("bad config initialization")
				}
			},
		)
	}
}
