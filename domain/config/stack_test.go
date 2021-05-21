// nolint: testpackage
package config

import (
	"encoding/json"
	"testing"
	"time"
)

func TestStackConfig_Marshall(t *testing.T) {
	t.Parallel()

	var a stackConfig
	_, err := json.MarshalIndent(a, "", "  ")
	if err != nil {
		t.Error(err)
	}
}

func TestStackConfig_UnMarshall(t *testing.T) {
	t.Parallel()

	var (
		a ApiStackResponse

		b ApiErrorResponse
	)

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
	t.Parallel()
	type fields struct {
		ConnectorScheme  string
		ConnectorAddress string
		ConnectorPort    string
		SchedulerAddr    string
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
	for i := range tests {
		tt := tests[i]

		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()

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

//nolint:funlen
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
		}, {
			name:    "bad broker",
			wantErr: true,
			fields: fields{
				data: []byte(_goldenStackConfigBadBroker),
			},
		}, {
			name:    "bad database",
			wantErr: true,
			fields: fields{
				data: []byte(_goldenStackConfigBadDatabase),
			},
		}, {
			name:    "bad interval",
			wantErr: true,
			fields: fields{
				data: []byte(_goldenStackConfigBadInterval),
			},
		}, {
			name:    "bad engine",
			wantErr: true,
			fields: fields{
				data: []byte(_goldenStackConfigBadEngine),
			},
		}, {
			name:    "bad scheduler",
			wantErr: true,
			fields: fields{
				data: []byte(_goldenStackConfigBadScheduler),
			},
		}, {
			name:    "bad profiler",
			wantErr: true,
			fields: fields{
				data: []byte(_goldenStackConfigBadProfiler),
			},
		}, {
			name:    "bad addresses",
			wantErr: true,
			fields: fields{
				data: []byte(_goldenStackConfigBadAddresses),
			},
		}, {
			name:    "bad sensor",
			wantErr: true,
			fields: fields{
				data: []byte(_goldenStackConfigBadSensor),
			},
		},
	}

	for i := range tests {
		tt := tests[i]
		t.Run(
			tt.name, func(t *testing.T) {
				t.Parallel()
				var stackResponse ApiStackResponse

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
                "connector_address": "dashboard.***REMOVED***",
                "connector_port": "2120"
            },
            "use_https": false,
            "created_at": "2021-05-16T13:25:29.782299Z",
            "updated_at": "2021-05-20T11:39:26.644745Z",
            "deployment": "docker",
            "database": {
                "db_addr": "database",
                "db_user": "database",
                "db_pass": "database",
                "db_root_pass": "database",
                "db_name": "strixeye",
                "db_port": "4354"
            },
            "broker": {
                "broker_addr": "",
                "broker_hostname": "abc.com",
                "broker_username": "broker",
                "broker_prefix": "amqp",
                "broker_password": "broker",
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
	_goldenStackConfigBadSensor = `{
    "status": "ok",
    "data": {
        "id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "company_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "name": " demo server",
        "ip_address": "154.154.154.154",
        "config": {
            "addresses": {
                "connector_scheme": "wss",
                "connector_address": "dashboard.***REMOVED***",
                "connector_port": "2120"
            },
            "use_https": false,
            "created_at": "2021-05-16T13:25:29.782299Z",
            "updated_at": "2021-05-20T11:39:26.644745Z",
            "deployment": "docker",
            "database": {
                "db_addr": "database",
                "db_user": "strixeye",
                "db_pass": "strixeye",
                "db_root_pass": "strixeye",
                "db_name": "strixeye",
                "db_port": "4354"
            },
            "broker": {
                "broker_addr": "",
                "broker_hostname": "abc.com",
                "broker_username": "strixeye",
                "broker_prefix": "amqp",
                "broker_password": "strixeye",
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
                "integration_name": "cloudflare",
                "sensor_listen": ":32122"
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
	_goldenStackConfigBadAddresses = `{
    "status": "ok",
    "data": {
        "id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "company_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "name": " demo server",
        "ip_address": "154.154.154.154",
        "config": {
            "addresses": {
                "connector_scheme": "ass",
                "connector_address": "dashboard.***REMOVED***",
                "connector_port": ":2120"
            },
            "use_https": false,
            "created_at": "2021-05-16T13:25:29.782299Z",
            "updated_at": "2021-05-20T11:39:26.644745Z",
            "deployment": "docker",
            "database": {
                "db_addr": "database",
                "db_user": "strixeye",
                "db_pass": "strixeye",
                "db_root_pass": "strixeye",
                "db_name": "strixeye",
                "db_port": "4354"
            },
            "broker": {
                "broker_addr": "",
                "broker_hostname": "abc.com",
                "broker_username": "strixeye",
                "broker_prefix": "amqp",
                "broker_password": "strixeye",
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
	_goldenStackConfigBadEngine = `{
    "status": "ok",
    "data": {
        "id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "company_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "name": " demo server",
        "ip_address": "154.154.154.154",
        "config": {
            "addresses": {
                "connector_scheme": "wss",
                "connector_address": "dashboard.***REMOVED***",
                "connector_port": "2120"
            },
            "use_https": false,
            "created_at": "2021-05-16T13:25:29.782299Z",
            "updated_at": "2021-05-20T11:39:26.644745Z",
            "deployment": "docker",
            "database": {
                "db_addr": "database",
                "db_user": "strixeye",
                "db_pass": "strixeye",
                "db_root_pass": "strixeye",
                "db_name": "strixeye",
                "db_port": "4354"
            },
            "broker": {
                "broker_addr": "",
                "broker_hostname": "abc.com",
                "broker_username": "strixeye",
                "broker_prefix": "amqp",
                "broker_password": "strixeye",
                "broker_port": "32121"
            },
            "scheduler": {
                "scheduler_listen": "2141"
            },
            "engine": {
                "address": "engine",
                "engine_listen": ":2130"
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
	_goldenStackConfigBadScheduler = `{
    "status": "ok",
    "data": {
        "id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "company_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "name": " demo server",
        "ip_address": "154.154.154.154",
        "config": {
            "addresses": {
                "connector_scheme": "wss",
                "connector_address": "dashboard.***REMOVED***",
                "connector_port": "2120"
            },
            "use_https": false,
            "created_at": "2021-05-16T13:25:29.782299Z",
            "updated_at": "2021-05-20T11:39:26.644745Z",
            "deployment": "docker",
            "database": {
                "db_addr": "database",
                "db_user": "strixeye",
                "db_pass": "strixeye",
                "db_root_pass": "strixeye",
                "db_name": "strixeye",
                "db_port": "4354"
            },
            "broker": {
                "broker_addr": "",
                "broker_hostname": "abc.com",
                "broker_username": "strixeye",
                "broker_prefix": "amqp",
                "broker_password": "strixeye",
                "broker_port": "32121"
            },
            "scheduler": {
                "scheduler_listen": ":2141"
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
	_goldenStackConfigBadProfiler = `{
    "status": "ok",
    "data": {
        "id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "company_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "name": " demo server",
        "ip_address": "154.154.154.154",
        "config": {
            "addresses": {
                "connector_scheme": "wss",
                "connector_address": "dashboard.***REMOVED***",
                "connector_port": "2120"
            },
            "use_https": false,
            "created_at": "2021-05-16T13:25:29.782299Z",
            "updated_at": "2021-05-20T11:39:26.644745Z",
            "deployment": "docker",
            "database": {
                "db_addr": "database",
                "db_user": "strixeye",
                "db_pass": "strixeye",
                "db_root_pass": "strixeye",
                "db_name": "strixeye",
                "db_port": "4354"
            },
            "broker": {
                "broker_addr": "",
                "broker_hostname": "abc.com",
                "broker_username": "strixeye",
                "broker_prefix": "amqp",
                "broker_password": "strixeye",
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
                "profiler_listen": ":2142"
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
	_goldenStackConfigBadInterval = `{
    "status": "ok",
    "data": {
        "id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "company_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "name": " demo server",
        "ip_address": "154.154.154.154",
        "config": {
            "addresses": {
                "connector_scheme": "wss",
                "connector_address": "dashboard.***REMOVED***",
                "connector_port": "2120"
            },
            "use_https": false,
            "created_at": "2021-05-16T13:25:29.782299Z",
            "updated_at": "2021-05-20T11:39:26.644745Z",
            "deployment": "docker",
            "database": {
                "db_addr": "database",
                "db_user": "strixeye",
                "db_pass": "strixeye",
                "db_root_pass": "strixeye",
                "db_name": "strixeye",
                "db_port": "4354"
            },
            "broker": {
                "broker_addr": "",
                "broker_hostname": "abc.com",
                "broker_username": "strixeye",
                "broker_prefix": "amqp",
                "broker_password": "strixeye",
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
                "system_stats_interval_second": 0
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
	_goldenStackConfigBadDatabase = `{
    "status": "ok",
    "data": {
        "id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "company_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "name": " demo server",
        "ip_address": "154.154.154.154",
        "config": {
            "addresses": {
                "connector_scheme": "wss",
                "connector_address": "dashboard.***REMOVED***",
                "connector_port": "2120"
            },
            "use_https": false,
            "created_at": "2021-05-16T13:25:29.782299Z",
            "updated_at": "2021-05-20T11:39:26.644745Z",
            "deployment": "docker",
            "database": {
                "db_addr": "database",
                "db_user": "strixeye",
                "db_pass": "strixeye",
                "db_root_pass": "strixeye",
                "db_name": "strixeye",
                "db_port": "strixeye"
            },
            "broker": {
                "broker_addr": "",
                "broker_hostname": "abc.com",
                "broker_username": "strixeye",
                "broker_prefix": "amqp",
                "broker_password": "strixeye",
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
	_goldenStackConfigBadBroker = `{
    "status": "ok",
    "data": {
        "id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "company_id": "0ff40874-41b8-4fee-88fc-5552d0d34033",
        "name": " demo server",
        "ip_address": "154.154.154.154",
        "config": {
            "addresses": {
                "connector_scheme": "wss",
                "connector_address": "dashboard.***REMOVED***",
                "connector_port": "2120"
            },
            "use_https": false,
            "created_at": "2021-05-16T13:25:29.782299Z",
            "updated_at": "2021-05-20T11:39:26.644745Z",
            "deployment": "docker",
            "database": {
                "db_addr": "database",
                "db_user": "strixeye",
                "db_pass": "strixeye",
                "db_root_pass": "strixeye",
                "db_name": "strixeye",
                "db_port": "4354"
            },
            "broker": {
                "broker_addr": "",
                "broker_hostname": "abc.com",
                "broker_username": "strixeye",
                "broker_prefix": "strixeye",
                "broker_password": "",
                "broker_port": "bbb"
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

//nolint:funlen
func Test_stackConfig_Save(t *testing.T) {
	type fields struct {
		Addresses  addresses
		UseHTTPS   bool
		CreatedAt  time.Time
		UpdatedAt  time.Time
		Deployment string
		Database   database
		Broker     broker
		Scheduler  scheduler
		Engine     engine
		Sensor     sensor
		Profiler   profiler
		Intervals  intervals
		Paths      paths
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
				t.Parallel()

				config := stackConfig{
					Addresses:  tt.fields.Addresses,
					UseHTTPS:   tt.fields.UseHTTPS,
					CreatedAt:  tt.fields.CreatedAt,
					UpdatedAt:  tt.fields.UpdatedAt,
					Deployment: tt.fields.Deployment,
					Database:   tt.fields.Database,
					Broker:     tt.fields.Broker,
					Scheduler:  tt.fields.Scheduler,
					Engine:     tt.fields.Engine,
					Sensor:     tt.fields.Sensor,
					Profiler:   tt.fields.Profiler,
					Intervals:  tt.fields.Intervals,
					Paths:      tt.fields.Paths,
				}
				if err := config.Save(tt.args.filePath); (err != nil) != tt.wantErr {
					t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
				}
			},
		)
	}
}
