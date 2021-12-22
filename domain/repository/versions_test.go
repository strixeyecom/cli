package repository

import (
	"encoding/json"
	"testing"
	"time"
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
					Installer: tt.fields.InstallVersion,
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
						Key       string  `json:"key"`
						Value     Version `json:"value"`
						CreatedAt string  `json:"created_at"`
						UpdatedAt string  `json:"updated_at"`
						DeletedAt string  `json:"deleted_at"`
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
	goldenVersionsAPIResponseGood = `
{
   "status":"ok",
   "data":[
      {
         "key":"installer_version",
         "value":{
            "version":"1.2.3",
            "checksum":"hash",
            "size":123456
         },
         "created_at":"2021-06-14 20:07:02 (Europe\/Istanbul)",
         "updated_at":"2021-07-04 18:58:54 (Europe\/Istanbul)",
         "deleted_at":null
      },
      {
         "key":"dashboard_version",
         "value":{
            "version":"1.2.3",
            "checksum":"hash",
            "size":123456
         },
         "created_at":"2021-06-14 20:07:02 (Europe\/Istanbul)",
         "updated_at":"2021-07-04 18:58:54 (Europe\/Istanbul)",
         "deleted_at":null
      },
      {
         "key":"static_engine_version",
         "value":{
            "version":"v1.5.66-emptytrip-66a9622",
            "size":17476080,
            "checksum":"4c1561ae2f3f705d7533ea22e94f4445"
         },
         "created_at":"2021-11-17 12:17:20 (Europe\/Istanbul)",
         "updated_at":"2021-11-30 20:48:33 (Europe\/Istanbul)",
         "deleted_at":null
      },
      {
         "key":"database_version",
         "value":{
            "version":"v0.0.82-pipeline-3c436cb",
            "size":1,
            "checksum":"4d1191fc8b7f56e836509133cbe5bb04"
         },
         "created_at":"2021-06-14 20:07:02 (Europe\/Istanbul)",
         "updated_at":"2021-11-13 11:11:41 (Europe\/Istanbul)",
         "deleted_at":null
      },
      {
         "key":"queue_version",
         "value":{
            "version":"v0.0.6-rc1-fbeb9f0",
            "size":60071104,
            "checksum":"d41d8cd98f00b204e9800998ecf8427e"
         },
         "created_at":"2021-06-14 20:07:02 (Europe\/Istanbul)",
         "updated_at":"2021-08-11 19:09:11 (Europe\/Istanbul)",
         "deleted_at":null
      },
      {
         "key":"profiling_database_version",
         "value":{
            "version":"v0.1.0",
            "size":24254216,
            "checksum":"b73a81f5c248cc857abb163060831220"
         },
         "created_at":"2021-06-14 20:07:02 (Europe\/Istanbul)",
         "updated_at":"2021-09-24 09:52:19 (Europe\/Istanbul)",
         "deleted_at":null
      },
      {
         "key":"scheduler_version",
         "value":{
            "version":"v0.2.63-dep-bc5c6a4",
            "size":11952128,
            "checksum":"e86358fa69fc7b0f2b639e0c0045327a"
         },
         "created_at":"2021-06-14 20:07:02 (Europe\/Istanbul)",
         "updated_at":"2021-11-10 12:24:46 (Europe\/Istanbul)",
         "deleted_at":null
      },
      {
         "key":"sensor_version",
         "value":{
            "version":"v0.2.98-alpha3-9a28f10",
            "size":15646984,
            "checksum":"1d5b4857ceafc222fa66d39fdef92c88"
         },
         "created_at":"2021-06-14 20:07:02 (Europe\/Istanbul)",
         "updated_at":"2021-12-01 09:49:31 (Europe\/Istanbul)",
         "deleted_at":null
      },
      {
         "key":"profiler_version",
         "value":{
            "version":"v0.3.92-alpha11-c277e55",
            "size":24981256,
            "checksum":"404339fad11175a69939a3c409c9475b"
         },
         "created_at":"2021-06-14 20:07:02 (Europe\/Istanbul)",
         "updated_at":"2021-11-25 08:35:29 (Europe\/Istanbul)",
         "deleted_at":null
      },
      {
         "key":"manager_version",
         "value":{
            "version":"v0.3.9625-sniff-5115b21",
            "size":40674360,
            "checksum":"385a46c20cb88a41edec7c845ebf7ea4"
         },
         "created_at":"2021-06-14 20:07:02 (Europe\/Istanbul)",
         "updated_at":"2021-12-03 12:48:45 (Europe\/Istanbul)",
         "deleted_at":null
      },
      {
         "key":"cache_version",
         "value":{
            "version":"latest",
            "size":1111111,
            "checksum":"b959bf146d2ec0ba2ebfa06c8bf42107"
         },
         "created_at":"2021-11-17 13:08:09 (Europe\/Istanbul)",
         "updated_at":"2021-11-17 13:08:09 (Europe\/Istanbul)",
         "deleted_at":null
      },
      {
         "key":"engine_version",
         "value":{
            "version":"v1.5.883-alpha11-b8acd81",
            "size":23692456,
            "checksum":"85285180299596f666e2db92b488097a"
         },
         "created_at":"2021-06-14 20:07:02 (Europe\/Istanbul)",
         "updated_at":"2021-11-27 08:14:58 (Europe\/Istanbul)",
         "deleted_at":null
      }
   ]
}`
)

func TestUnmarshal(t *testing.T) {
	var s string
	const layout = "2006-01-2 15:04:05 (MST)"

	x, err := time.Parse(layout, `2021-06-14 20:07:02 (CST)`)
	if err != nil {
		t.Fatal(err)
	}
	_ = x
	err = json.Unmarshal([]byte(`{"time":"2021-06-14 20:07:02 (Europe\/Istanbul)`), &s)
	if err != nil {
		t.Fatal(err)
	}
	var mes APIVersionsMessage
	err = json.Unmarshal([]byte(goldenVersionsAPIResponseGood), &mes)
	if err != nil {
		t.Fatal(err)
	}
}
