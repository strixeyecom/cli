package trip

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

func TestQueryArgs_String(t *testing.T) {
	type fields struct {
		Limit      int
		SinceTime  int64
		TripsIds   []string
		SuspectIds []string
		Endpoints  []string
		Verbose    bool
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				q := QueryArgs{
					Limit:      tt.fields.Limit,
					SinceTime:  tt.fields.SinceTime,
					TripsIds:   tt.fields.TripsIds,
					SuspectIds: tt.fields.SuspectIds,
					Endpoints:  tt.fields.Endpoints,
					Verbose:    tt.fields.Verbose,
				}
				if got := q.String(); got != tt.want {
					t.Errorf("String() = %v, want %v", got, tt.want)
				}
			},
		)
	}
}
