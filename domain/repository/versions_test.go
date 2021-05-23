package repository

import "testing"

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
	for i := range tests {
		tt := tests[i]

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

func TestVersions_FromRawAPIResponse(t *testing.T) {
	type fields struct {
		ManagerVersion   string
		DatabaseVersion  string
		EngineVersion    string
		ProfilerVersion  string
		QueueVersion     string
		SchedulerVersion string
		SensorVersion    string
	}
	type args struct {
		rawData []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "good api response",
			fields: fields{},
			args: args{
				rawData: []byte(goldenVersionsAPIResponseGood),
			},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]

		t.Run(
			tt.name, func(t *testing.T) {
				v := Versions{}
				if err := v.FromRawAPIResponse(tt.args.rawData); (err != nil) != tt.wantErr {
					t.Errorf("FromRawAPIResponse() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

const (
	goldenVersionsAPIResponseGood = `{
    "status": "ok",
    "data": [
        {
            "key": "installer_version",
            "value": "v0.1.3-rc1-8deb131",
            "created_at": null,
            "updated_at": null,
            "deleted_at": null
        },
        {
            "key": "sensor_version",
            "value": "staging",
            "created_at": null,
            "updated_at": null,
            "deleted_at": null
        },
        {
            "key": "dashboard_version",
            "value": "staging",
            "created_at": null,
            "updated_at": null,
            "deleted_at": null
        },
        {
            "key": "scheduler_version",
            "value": "v0.1.2-rc2-0e4ef28",
            "created_at": null,
            "updated_at": "2021-05-19T10:21:41.000000Z",
            "deleted_at": null
        },
        {
            "key": "queue_version",
            "value": "staging",
            "created_at": null,
            "updated_at": null,
            "deleted_at": null
        },
        {
            "key": "profiler_version",
            "value": "v0.2.4-rc6-d39fa3f",
            "created_at": null,
            "updated_at": "2021-05-19T10:06:44.000000Z",
            "deleted_at": null
        },
        {
            "key": "database_version",
            "value": "staging",
            "created_at": null,
            "updated_at": null,
            "deleted_at": null
        },
        {
            "key": "engine_version",
            "value": "v0.2.9-rc1.3-b81395e",
            "created_at": null,
            "updated_at": "2021-05-19T09:40:09.000000Z",
            "deleted_at": null
        },
        {
            "key": "manager_version",
            "value": "0.2.5-rc1.1",
            "created_at": null,
            "updated_at": null,
            "deleted_at": null
        },
        {
            "key": "backoffice_version",
            "value": "staging",
            "created_at": null,
            "updated_at": null,
            "deleted_at": null
        }
    ]
}`
)
