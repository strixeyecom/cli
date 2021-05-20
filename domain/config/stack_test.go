package config

import (
	`encoding/json`
	`testing`
)

func TestStackConfig_Marshall(t *testing.T) {
	a := stackConfig{
	
	}
	_, err := json.MarshalIndent(a, "", "  ")
	
	if err != nil {
		t.Error(err)
	}
}

func TestStackConfig_UnMarshall(t *testing.T) {
	a := ApiStackResponse{}
	b := ApiErrorResponse{}
	err := json.Unmarshal(
		[]byte(_goldenStackConfig), &a,
	)
	if err != nil {
		t.Error(err)
	}
	
	err = json.Unmarshal(
		[]byte(`{
    "status": "error",
    "data": {
        "agent_id": [
            "The agent id must be a valid UUID."
        ]
    }
}`), &b,
	)
	if err != nil {
		t.Error(err)
	}
}

func Test_addresses_Validate(t *testing.T) {
	type fields struct {
		ConnectorScheme  string
		ConnectorAddress string
		ConnectorPort    string
		SchedulerAddr    string
		ManagerPort      string
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "pass",
			fields: fields{
				ConnectorScheme:  "wss",
				ConnectorAddress: "dashboard.strixeye.com",
				ConnectorPort:    "2118",
				SchedulerAddr:    "2141",
			},
			
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				a := addresses{
					ConnectorScheme:  tt.fields.ConnectorScheme,
					ConnectorAddress: tt.fields.ConnectorAddress,
					ConnectorPort:    tt.fields.ConnectorPort,
				}
				if err := a.validate(); (err != nil) != tt.wantErr {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

func Test_stackConfig_Validate(t *testing.T) {
	type fields struct {
		data []byte
	}
	tests := []struct {
		name    string
		wantErr bool
		fields  fields
	}{
		{
			name:    "good config",
			wantErr: false,
			fields: fields{
				data: []byte(_goldenStackConfig),
			},
		},
	}
	for _, tt := range tests {
		t.Run(
			tt.name, func(t *testing.T) {
				stackResponse := ApiStackResponse{}
				if err := json.Unmarshal(tt.fields.data, &stackResponse); err != nil {
					t.Error(err)
				}
				
				config := stackResponse.Stack.Config
				if err := config.Validate(); (err != nil) != tt.wantErr {
					t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}

const (
	// _goldenStackConfig im not sure if i should tell this,
	// but these tokens and configs and passwords are randomized/masked. Don't get excited.
	_goldenStackConfig = `{
    "status": "ok",
    "data": {
        "id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "company_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "name": " demo server",
        "ip_address": "154.154.154.154",
        "config": {
            "addresses": {
                "connector_scheme": "wss",
                "connector_address": "***REMOVED***",
                "connector_port": "2120"
            },
            "use_https": false,
            "created_at": "2021-05-16T13:25:29.782299Z",
            "updated_at": "2021-05-20T11:39:26.644745Z",
            "deployment": "docker",
            "database": {
                "db_addr": "database",
                "db_user": "aaaaaaaa",
                "db_pass": "aaaaaaaa",
                "db_root_pass": "aaaaaaaa",
                "db_name": "strixeye",
                "db_port": "4354"
            },
            "broker": {
                "broker_addr": "",
                "broker_hostname": "abc.com",
                "broker_username": "aaaaaaaa",
                "broker_prefix": "amqp",
                "broker_password": "aaaaaaaa",
                "broker_port": "32121"
            },
            "scheduler": {
                "scheduler_listen": "2141"
            },
            "engine": {
                "address": "engine",
                "engine_listen": "2130"
            },
            "sensor": {
                "integration_name": "nginx",
                "sensor_listen": "32122"
            },
            "profiler": {
                "profiler_listen": "2142"
            },
            "intervals": {
                "system_stats_interval_second": 1
            },
            "paths": {
                "kube_config": "",
                "tls_keys": {
                    "certificate": "",
                    "key": ""
                }
            }
        },
        "domains": [
            {
                "id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
                "company_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
                "domain": "localhost",
                "deleted_at": null,
                "created_at": "2021-05-16T13:25:28.000000Z",
                "updated_at": "2021-05-16T13:25:28.000000Z",
                "pivot": {
                    "agent_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
                    "domain_id": "0ff40874-41b8-4fee-88fc-5552d0d34033"
                }
            },
            {
                "id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
                "company_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
                "domain": "127.0.0.1",
                "deleted_at": null,
                "created_at": "2021-05-16T13:25:28.000000Z",
                "updated_at": "2021-05-16T13:25:28.000000Z",
                "pivot": {
                    "agent_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
                    "domain_id": "0ff40874-41b8-4fee-88fc-5552d0d34033"
                }
            }
        ]
    }
}`
)
