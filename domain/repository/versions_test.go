package repository

import (
	`encoding/json`
	`testing`
	`time`
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

func TestVersions_Validate(t *testing.T) {
	type fields struct {
	ManagerVersion   Version
	DatabaseVersion  Version
	EngineVersion    Version
	ProfilerVersion  Version
	QueueVersion     Version
	SchedulerVersion Version
	SensorVersion    Version
	InstallVersion   Version
}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "empty fields exist",
			fields: fields{
				ManagerVersion: Version{Version: "v0.2.5-rc1.1-2505350", Size: 123456, Checksum: "checksum"},
			},
			wantErr: true,
		}, {
			name: "version with commit hash",
			fields: fields{
				ManagerVersion: Version{
					Version: "v0.2.5-rc1.1", Size: 123456, Checksum: "checksum",
				},
				DatabaseVersion: Version{
					Version: "v0.2.5-rc1.1", Size: 123456, Checksum: "checksum",
				},
				EngineVersion: Version{
					Version: "v0.2.5-rc1.1", Size: 123456, Checksum: "checksum",
				},
				ProfilerVersion: Version{
					Version: "v0.2.5-rc1.1", Size: 123456, Checksum: "checksum",
				},
				QueueVersion: Version{
					Version: "v0.2.5-rc1.1", Size: 123456, Checksum: "checksum",
				},
				SchedulerVersion: Version{
					Version: "v0.2.5-rc1.1", Size: 123456, Checksum: "checksum",
				},
				SensorVersion: Version{
					Version: "v0.2.5-rc1.1", Size: 123456, Checksum: "checksum",
				},
				InstallVersion: Version{
					Version: "v0.2.5-rc1.1", Size: 123456, Checksum: "checksum",
				},
			},
			wantErr: false,
		},
	}
	for i := range tests {
		tt := tests[i]
		
		t.Run(
			tt.name, func(t *testing.T) {
				m := Versions{
					Manager:   tt.fields.ManagerVersion,
					Database:  tt.fields.DatabaseVersion,
					Engine:    tt.fields.EngineVersion,
					Profiler:  tt.fields.ProfilerVersion,
					Queue:     tt.fields.QueueVersion,
					Scheduler: tt.fields.SchedulerVersion,
					Sensor:    tt.fields.SensorVersion,
					Installer:    tt.fields.InstallVersion,
				}
				if err := m.Validate(); (err != nil) != tt.wantErr {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func TestVersions_ToVersions(t *testing.T) {
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
				v := APIVersionsMessage{}
				err := json.Unmarshal(tt.args.rawData, &v)
				if err != nil {
					t.Fatal(err)
				}
				if _, err := v.ToVersions(); (err != nil) != tt.wantErr {
					t.Errorf("FromRawAPIResponse() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

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
						Key       string     `json:"key"`
						Value     Version    `json:"value"`
						CreatedAt *time.Time `json:"created_at"`
						UpdatedAt *time.Time `json:"updated_at"`
						DeletedAt *time.Time `json:"deleted_at"`
					}{
						{
							Key: "Manager", Value: Version{
							Version: "test-version", Size: 12345,
							Checksum: "test-checksum",
						},
						}, {
							Key: "Installer", Value: Version{
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

const (
	goldenVersionsAPIResponseGood = `{
    "status": "ok",
    "data": [
        {
            "key": "dashboard_version",
            "value": {
                "version": "staging",
                "checksum": "hash",
                "size": 123456
            },
            "created_at": "2021-06-14T17:07:02.000000Z",
            "updated_at": "2021-06-14T17:07:02.000000Z",
            "deleted_at": null
        },
        {
            "key": "installer_version",
            "value": {
                "version": "v0.2.5-rc1.1-2505350",
                "checksum": "hash",
                "size": 123456
            },
            "created_at": "2021-06-14T17:07:02.000000Z",
            "updated_at": "2021-06-14T17:07:02.000000Z",
            "deleted_at": null
        },
        {
            "key": "scheduler_version",
            "value": {
                "version": "staging",
                "checksum": "hash",
                "size": 123456
            },
            "created_at": "2021-06-14T17:07:02.000000Z",
            "updated_at": "2021-06-14T17:07:02.000000Z",
            "deleted_at": null
        },
        {
            "key": "queue_version",
            "value": {
                "version": "staging",
                "checksum": "hash",
                "size": 123456
            },
            "created_at": "2021-06-14T17:07:02.000000Z",
            "updated_at": "2021-06-14T17:07:02.000000Z",
            "deleted_at": null
        },
        {
            "key": "profiler_version",
            "value": {
                "version": "staging",
                "checksum": "hash",
                "size": 123456
            },
            "created_at": "2021-06-14T17:07:02.000000Z",
            "updated_at": "2021-06-14T17:07:02.000000Z",
            "deleted_at": null
        },
        {
            "key": "database_version",
            "value": {
                "version": "staging",
                "checksum": "hash",
                "size": 123456
            },
            "created_at": "2021-06-14T17:07:02.000000Z",
            "updated_at": "2021-06-14T17:07:02.000000Z",
            "deleted_at": null
        },
        {
            "key": "sensor_version",
            "value": {
                "version": "staging",
                "checksum": "hash",
                "size": 123456
            },
            "created_at": "2021-06-14T17:07:02.000000Z",
            "updated_at": "2021-06-14T17:07:02.000000Z",
            "deleted_at": null
        },
        {
            "key": "engine_version",
            "value": {
                "version": "staging",
                "checksum": "hash",
                "size": 111
            },
            "created_at": "2021-06-14T17:07:02.000000Z",
            "updated_at": "2021-06-14T17:07:02.000000Z",
            "deleted_at": null
        },
        {
            "key": "manager_version",
            "value": {
                "version": "v0.2.5-rc1.1-f2c7578",
                "size": 111,
                "hash": "111"
            },
            "created_at": "2021-06-14T17:07:02.000000Z",
            "updated_at": "2021-06-15T18:15:57.000000Z",
            "deleted_at": null
        }
    ]
}`
)
