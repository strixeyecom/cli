package config

import `testing`

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

func TestVersions_Validate(t *testing.T) {
	type fields struct {
		ManagerVersion   string
		DatabaseVersion  string
		EngineVersion    string
		ProfilerVersion  string
		QueueVersion     string
		SchedulerVersion string
		SensorVersion    string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty fields exist",
			fields: fields{
				ManagerVersion: "bah",
			},
			wantErr: true,
		}, {
			name: "version with commit hash",
			fields: fields{
				ManagerVersion:   "v0.2.5-rc1.1-2505350",
				DatabaseVersion:  "v0.2.5-rc1.1-2505350",
				EngineVersion:    "v0.2.5-rc1.1-2505350",
				ProfilerVersion:  "v0.2.5-rc1.1-2505350",
				QueueVersion:     "v0.2.5-rc1.1-2505350",
				SchedulerVersion: "v0.2.5-rc1.1-2505350",
				SensorVersion:    "v0.2.5-rc1.1-2505350",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				m := Versions{
					ManagerVersion:   tt.fields.ManagerVersion,
					DatabaseVersion:  tt.fields.DatabaseVersion,
					EngineVersion:    tt.fields.EngineVersion,
					ProfilerVersion:  tt.fields.ProfilerVersion,
					QueueVersion:     tt.fields.QueueVersion,
					SchedulerVersion: tt.fields.SchedulerVersion,
					SensorVersion:    tt.fields.SensorVersion,
				}
				if err := m.Validate(); (err != nil) != tt.wantErr {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
