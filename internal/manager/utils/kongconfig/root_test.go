package kongconfig

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRoot(t *testing.T) {
	var root Root
	require.NoError(t, json.Unmarshal([]byte(dblessConfigJSON), &root))
	require.NoError(t, root.Validate(false))
	require.EqualError(t, root.Validate(true), "--skip-ca-certificates is not available for use with DB-less Kong instances")
}

func TestValidateRoots(t *testing.T) {
	var root Root
	require.NoError(t, json.Unmarshal([]byte(dblessConfigJSON), &root))
	dbmode, v, err := ValidateRoots([]Root{root, root}, false)
	require.NoError(t, err)
	assert.Equal(t, "off", dbmode)
	assert.Equal(t, "3.1.1", v.String())
}

const dblessConfigJSON = `{
	"plugins": {
		"enabled_in_cluster": [],
		"available_on_server": {
			"basic-auth": {
				"version": "3.1.1",
				"priority": 1100
			},
			"ip-restriction": {
				"version": "3.1.1",
				"priority": 990
			},
			"request-transformer": {
				"version": "3.1.1",
				"priority": 801
			},
			"response-transformer": {
				"version": "3.1.1",
				"priority": 800
			},
			"request-size-limiting": {
				"version": "3.1.1",
				"priority": 951
			},
			"rate-limiting": {
				"version": "3.1.1",
				"priority": 910
			},
			"response-ratelimiting": {
				"version": "3.1.1",
				"priority": 900
			},
			"syslog": {
				"version": "3.1.1",
				"priority": 4
			},
			"loggly": {
				"version": "3.1.1",
				"priority": 6
			},
			"datadog": {
				"version": "3.1.1",
				"priority": 10
			},
			"ldap-auth": {
				"version": "3.1.1",
				"priority": 1200
			},
			"statsd": {
				"version": "3.1.1",
				"priority": 11
			},
			"bot-detection": {
				"version": "3.1.1",
				"priority": 2500
			},
			"aws-lambda": {
				"version": "3.1.1",
				"priority": 750
			},
			"request-termination": {
				"version": "3.1.1",
				"priority": 2
			},
			"prometheus": {
				"version": "3.1.1",
				"priority": 13
			},
			"proxy-cache": {
				"version": "3.1.1",
				"priority": 100
			},
			"session": {
				"version": "3.1.1",
				"priority": 1900
			},
			"acme": {
				"version": "3.1.1",
				"priority": 1705
			},
			"grpc-gateway": {
				"version": "3.1.1",
				"priority": 998
			},
			"grpc-web": {
				"version": "3.1.1",
				"priority": 3
			},
			"pre-function": {
				"version": "3.1.1",
				"priority": 1000000
			},
			"post-function": {
				"version": "3.1.1",
				"priority": -1000
			},
			"azure-functions": {
				"version": "3.1.1",
				"priority": 749
			},
			"zipkin": {
				"version": "3.1.1",
				"priority": 100000
			},
			"opentelemetry": {
				"version": "0.1.0",
				"priority": 14
			},
			"jwt": {
				"version": "3.1.1",
				"priority": 1450
			},
			"acl": {
				"version": "3.1.1",
				"priority": 950
			},
			"correlation-id": {
				"version": "3.1.1",
				"priority": 1
			},
			"cors": {
				"version": "3.1.1",
				"priority": 2000
			},
			"oauth2": {
				"version": "3.1.1",
				"priority": 1400
			},
			"tcp-log": {
				"version": "3.1.1",
				"priority": 7
			},
			"udp-log": {
				"version": "3.1.1",
				"priority": 8
			},
			"file-log": {
				"version": "3.1.1",
				"priority": 9
			},
			"http-log": {
				"version": "3.1.1",
				"priority": 12
			},
			"key-auth": {
				"version": "3.1.1",
				"priority": 1250
			},
			"hmac-auth": {
				"version": "3.1.1",
				"priority": 1030
			}
		}
	},
	"tagline": "Welcome to kong",
	"lua_version": "LuaJIT 2.1.0-20220411",
	"version": "3.1.1",
	"pids": {
		"master": 1,
		"workers": [
			1133,
			1134
		]
	},
	"configuration": {
		"nginx_http_client_body_buffer_size": "8k",
		"proxy_ssl_enabled": true,
		"cluster_listeners": {},
		"role": "traditional",
		"vaults": [
			"bundled"
		],
		"ssl_cert": [
			"/kong_prefix/ssl/kong-default.crt",
			"/kong_prefix/ssl/kong-default-ecdsa.crt"
		],
		"loaded_vaults": {
			"env": true
		},
		"real_ip_header": "X-Real-IP",
		"nginx_proxy_real_ip_header": "X-Real-IP",
		"real_ip_recursive": "off",
		"nginx_proxy_real_ip_recursive": "off",
		"pg_port": 5432,
		"pg_ssl": false,
		"pg_ssl_verify": false,
		"pg_max_concurrent_queries": 0,
		"pg_semaphore_timeout": 60000,
		"log_level": "notice",
		"pg_ro_ssl_verify": false,
		"cassandra_contact_points": [
			"127.0.0.1"
		],
		"cassandra_port": 9042,
		"cassandra_ssl": false,
		"cassandra_ssl_verify": false,
		"cassandra_write_consistency": "ONE",
		"cassandra_read_consistency": "ONE",
		"cassandra_lb_policy": "RequestRoundRobin",
		"cassandra_refresh_frequency": 60,
		"cassandra_repl_strategy": "SimpleStrategy",
		"prefix": "/kong_prefix",
		"cassandra_repl_factor": 1,
		"cassandra_data_centers": [
			"dc1:2",
			"dc2:3"
		],
		"status_ssl_enabled": false,
		"ssl_protocols": "TLSv1.1 TLSv1.2 TLSv1.3",
		"lua_ssl_trusted_certificate_combined": "/kong_prefix/.ca_combined",
		"nginx_http_ssl_protocols": "TLSv1.2 TLSv1.3",
		"nginx_stream_ssl_protocols": "TLSv1.2 TLSv1.3",
		"ssl_prefer_server_ciphers": "on",
		"nginx_http_ssl_prefer_server_ciphers": "off",
		"nginx_stream_ssl_prefer_server_ciphers": "off",
		"ssl_dhparam": "ffdhe2048",
		"nginx_http_ssl_dhparam": "ffdhe2048",
		"nginx_stream_ssl_dhparam": "ffdhe2048",
		"ssl_session_tickets": "on",
		"nginx_http_ssl_session_tickets": "on",
		"nginx_stream_ssl_session_tickets": "on",
		"ssl_session_timeout": "1d",
		"nginx_http_ssl_session_timeout": "1d",
		"nginx_stream_ssl_session_timeout": "1d",
		"proxy_access_log": "/dev/stdout",
		"proxy_error_log": "/dev/stderr",
		"proxy_stream_access_log": "/dev/stdout basic",
		"proxy_stream_error_log": "/dev/stderr",
		"admin_access_log": "/dev/stdout",
		"admin_error_log": "/dev/stderr",
		"status_access_log": "off",
		"status_error_log": "/dev/stderr",
		"lua_ssl_trusted_certificate": [
			"/etc/ssl/certs/ca-certificates.crt"
		],
		"lua_ssl_verify_depth": 1,
		"lua_ssl_protocols": "TLSv1.1 TLSv1.2 TLSv1.3",
		"nginx_http_lua_ssl_protocols": "TLSv1.1 TLSv1.2 TLSv1.3",
		"nginx_stream_lua_ssl_protocols": "TLSv1.1 TLSv1.2 TLSv1.3",
		"lua_socket_pool_size": 30,
		"nginx_admin_client_max_body_size": "10m",
		"cluster_mtls": "shared",
		"nginx_http_lua_regex_match_limit": "100000",
		"cassandra_timeout": 5000,
		"pg_timeout": 5000,
		"worker_state_update_frequency": 5,
		"cluster_max_payload": 4194304,
		"client_body_buffer_size": "8k",
		"untrusted_lua": "sandbox",
		"untrusted_lua_sandbox_requires": {},
		"untrusted_lua_sandbox_environment": {},
		"pg_host": "127.0.0.1",
		"lmdb_map_size": "128m",
		"proxy_server_ssl_verify": true,
		"pg_database": "kong",
		"nginx_events_directives": [
			{
				"name": "multi_accept",
				"value": "on"
			},
			{
				"name": "worker_connections",
				"value": "auto"
			}
		],
		"nginx_http_directives": [
			{
				"name": "client_body_buffer_size",
				"value": "8k"
			},
			{
				"name": "client_max_body_size",
				"value": "0"
			},
			{
				"name": "lua_regex_cache_max_entries",
				"value": "8192"
			},
			{
				"name": "lua_regex_match_limit",
				"value": "100000"
			},
			{
				"name": "lua_shared_dict",
				"value": "prometheus_metrics 5m"
			},
			{
				"name": "lua_ssl_protocols",
				"value": "TLSv1.1 TLSv1.2 TLSv1.3"
			},
			{
				"name": "ssl_dhparam",
				"value": "/kong_prefix/ssl/ffdhe2048.pem"
			},
			{
				"name": "ssl_prefer_server_ciphers",
				"value": "off"
			},
			{
				"name": "ssl_protocols",
				"value": "TLSv1.2 TLSv1.3"
			},
			{
				"name": "ssl_session_tickets",
				"value": "on"
			},
			{
				"name": "ssl_session_timeout",
				"value": "1d"
			}
		],
		"pg_user": "kong",
		"nginx_upstream_directives": {},
		"upstream_keepalive_idle_timeout": 60,
		"upstream_keepalive_max_requests": 100,
		"nginx_status_directives": {},
		"nginx_admin_directives": [
			{
				"name": "client_body_buffer_size",
				"value": "10m"
			},
			{
				"name": "client_max_body_size",
				"value": "10m"
			}
		],
		"nginx_stream_directives": [
			{
				"name": "lua_shared_dict",
				"value": "stream_prometheus_metrics 5m"
			},
			{
				"name": "lua_ssl_protocols",
				"value": "TLSv1.1 TLSv1.2 TLSv1.3"
			},
			{
				"name": "ssl_dhparam",
				"value": "/kong_prefix/ssl/ffdhe2048.pem"
			},
			{
				"name": "ssl_prefer_server_ciphers",
				"value": "off"
			},
			{
				"name": "ssl_protocols",
				"value": "TLSv1.2 TLSv1.3"
			},
			{
				"name": "ssl_session_tickets",
				"value": "on"
			},
			{
				"name": "ssl_session_timeout",
				"value": "1d"
			}
		],
		"nginx_supstream_directives": {},
		"nginx_sproxy_directives": {},
		"opentelemetry_tracing": [
			"off"
		],
		"nginx_pid": "/kong_prefix/pids/nginx.pid",
		"kic": false,
		"pluginserver_names": {},
		"nginx_err_logs": "/kong_prefix/logs/error.log",
		"nginx_events_multi_accept": "on",
		"pg_ro_ssl": false,
		"cassandra_keyspace": "kong",
		"nginx_events_worker_connections": "auto",
		"declarative_config": "/kong_dbless/kong.yml",
		"admin_ssl_cert": {},
		"nginx_conf": "/kong_prefix/nginx.conf",
		"cassandra_username": "kong",
		"nginx_kong_conf": "/kong_prefix/nginx-kong.conf",
		"nginx_main_worker_rlimit_nofile": "auto",
		"nginx_kong_stream_conf": "/kong_prefix/nginx-kong-stream.conf",
		"plugins": [
			"bundled"
		],
		"dns_hostsfile": "/etc/hosts",
		"kong_process_secrets": "/kong_prefix/.kong_process_secrets",
		"stream_listeners": {},
		"ssl_cert_csr_default": "/kong_prefix/ssl/kong-default.csr",
		"error_default_type": "text/plain",
		"ssl_cert_default": "/kong_prefix/ssl/kong-default.crt",
		"dns_error_ttl": 1,
		"ssl_cert_key_default": "/kong_prefix/ssl/kong-default.key",
		"dns_not_found_ttl": 30,
		"ssl_cert_default_ecdsa": "/kong_prefix/ssl/kong-default-ecdsa.crt",
		"dns_stale_ttl": 4,
		"mem_cache_size": "128m",
		"dns_cache_size": 10000,
		"client_ssl_cert_default": "/kong_prefix/ssl/kong-default.crt",
		"dns_order": [
			"LAST",
			"SRV",
			"A",
			"CNAME"
		],
		"admin_ssl_cert_default": "/kong_prefix/ssl/admin-kong-default.crt",
		"dns_no_sync": false,
		"admin_ssl_cert_key_default": "/kong_prefix/ssl/admin-kong-default.key",
		"client_ssl": false,
		"enabled_headers": {
			"Via": true,
			"Server": true,
			"X-Kong-Admin-Latency": true,
			"X-Kong-Upstream-Latency": true,
			"X-Kong-Upstream-Status": false,
			"server_tokens": true,
			"latency_tokens": true,
			"X-Kong-Response-Latency": true,
			"X-Kong-Proxy-Latency": true
		},
		"dns_resolver": {},
		"admin_ssl_cert_key_default_ecdsa": "/kong_prefix/ssl/admin-kong-default-ecdsa.key",
		"ssl_ciphers": "ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384",
		"status_ssl_cert_default": "/kong_prefix/ssl/status-kong-default.crt",
		"proxy_listen": [
			"0.0.0.0:8000",
			"0.0.0.0:8443 http2 ssl"
		],
		"status_ssl_cert_key_default": "/kong_prefix/ssl/status-kong-default.key",
		"opentelemetry_tracing_sampling_rate": 1,
		"status_ssl_cert_default_ecdsa": "/kong_prefix/ssl/status-kong-default-ecdsa.crt",
		"stream_proxy_ssl_enabled": false,
		"status_ssl_cert_key_default_ecdsa": "/kong_prefix/ssl/status-kong-default-ecdsa.key",
		"loaded_plugins": {
			"basic-auth": true,
			"ip-restriction": true,
			"request-transformer": true,
			"response-transformer": true,
			"request-size-limiting": true,
			"rate-limiting": true,
			"response-ratelimiting": true,
			"syslog": true,
			"loggly": true,
			"datadog": true,
			"ldap-auth": true,
			"statsd": true,
			"bot-detection": true,
			"aws-lambda": true,
			"request-termination": true,
			"prometheus": true,
			"proxy-cache": true,
			"session": true,
			"acme": true,
			"grpc-gateway": true,
			"grpc-web": true,
			"pre-function": true,
			"post-function": true,
			"azure-functions": true,
			"zipkin": true,
			"opentelemetry": true,
			"jwt": true,
			"acl": true,
			"correlation-id": true,
			"cors": true,
			"oauth2": true,
			"tcp-log": true,
			"udp-log": true,
			"file-log": true,
			"http-log": true,
			"key-auth": true,
			"hmac-auth": true
		},
		"ssl_cert_key": [
			"/kong_prefix/ssl/kong-default.key",
			"/kong_prefix/ssl/kong-default-ecdsa.key"
		],
		"lua_package_cpath": "",
		"port_maps": [
			"80:8000",
			"443:8443"
		],
		"lua_package_path": "/opt/?.lua;/opt/?/init.lua;;",
		"admin_listen": [
			"0.0.0.0:8001"
		],
		"status_listen": [
			"0.0.0.0:8100"
		],
		"stream_listen": [
			"off"
		],
		"cluster_listen": [
			"off"
		],
		"db_cache_warmup_entities": [
			"services"
		],
		"admin_ssl_cert_key": {},
		"status_ssl_cert": {},
		"status_ssl_cert_key": {},
		"db_resurrect_ttl": 30,
		"nginx_user": "kong kong",
		"headers": [
			"server_tokens",
			"latency_tokens"
		],
		"nginx_daemon": "off",
		"cluster_ocsp": "off",
		"nginx_main_daemon": "off",
		"nginx_worker_processes": "2",
		"nginx_main_worker_processes": "2",
		"trusted_ips": {},
		"upstream_keepalive_pool_size": 60,
		"anonymous_reports": true,
		"host_ports": {
			"8000": 80,
			"8443": 443
		},
		"cluster_data_plane_purge_delay": 1209600,
		"lmdb_environment_path": "dbless.lmdb",
		"cluster_use_proxy": false,
		"database": "off",
		"router_flavor": "traditional",
		"legacy_worker_events": false,
		"cluster_control_plane": "127.0.0.1:8005",
		"admin_listeners": [
			{
				"backlog=%d+": false,
				"ipv6only=on": false,
				"ipv6only=off": false,
				"ssl": false,
				"so_keepalive=off": false,
				"so_keepalive=%w*:%w*:%d*": false,
				"listener": "0.0.0.0:8001",
				"bind": false,
				"port": 8001,
				"deferred": false,
				"so_keepalive=on": false,
				"http2": false,
				"proxy_protocol": false,
				"ip": "0.0.0.0",
				"reuseport": false
			}
		],
		"nginx_proxy_directives": [
			{
				"name": "real_ip_header",
				"value": "X-Real-IP"
			},
			{
				"name": "real_ip_recursive",
				"value": "off"
			}
		],
		"nginx_admin_client_body_buffer_size": "10m",
		"ssl_cert_key_default_ecdsa": "/kong_prefix/ssl/kong-default-ecdsa.key",
		"nginx_main_user": "kong kong",
		"nginx_main_directives": [
			{
				"name": "daemon",
				"value": "off"
			},
			{
				"name": "user",
				"value": "kong kong"
			},
			{
				"name": "worker_processes",
				"value": "2"
			},
			{
				"name": "worker_rlimit_nofile",
				"value": "auto"
			}
		],
		"nginx_http_lua_regex_cache_max_entries": "8192",
		"worker_consistency": "eventual",
		"admin_ssl_enabled": false,
		"admin_ssl_cert_default_ecdsa": "/kong_prefix/ssl/admin-kong-default-ecdsa.crt",
		"nginx_acc_logs": "/kong_prefix/logs/access.log",
		"status_listeners": [
			{
				"port": 8100,
				"ip": "0.0.0.0",
				"listener": "0.0.0.0:8100",
				"ssl": false
			}
		],
		"admin_acc_logs": "/kong_prefix/logs/admin_access.log",
		"cassandra_schema_consensus_timeout": 10000,
		"db_update_frequency": 5,
		"proxy_listeners": [
			{
				"backlog=%d+": false,
				"ipv6only=on": false,
				"ipv6only=off": false,
				"ssl": false,
				"so_keepalive=off": false,
				"so_keepalive=%w*:%w*:%d*": false,
				"listener": "0.0.0.0:8000",
				"bind": false,
				"port": 8000,
				"deferred": false,
				"so_keepalive=on": false,
				"http2": false,
				"proxy_protocol": false,
				"ip": "0.0.0.0",
				"reuseport": false
			},
			{
				"backlog=%d+": false,
				"ipv6only=on": false,
				"ipv6only=off": false,
				"ssl": true,
				"so_keepalive=off": false,
				"so_keepalive=%w*:%w*:%d*": false,
				"listener": "0.0.0.0:8443 ssl http2",
				"bind": false,
				"port": 8443,
				"deferred": false,
				"so_keepalive=on": false,
				"http2": true,
				"proxy_protocol": false,
				"ip": "0.0.0.0",
				"reuseport": false
			}
		],
		"db_update_propagation": 0,
		"client_ssl_cert_key_default": "/kong_prefix/ssl/kong-default.key",
		"db_cache_ttl": 0,
		"kong_env": "/kong_prefix/.kong_env",
		"nginx_http_client_max_body_size": "0",
		"ssl_cipher_suite": "intermediate"
	},
	"node_id": "10932f25-0601-4c18-a8e2-811c64995df2",
	"timers": {
		"running": 37,
		"pending": 1
	},
	"hostname": "kong-kong-8d76d469d-zdjhb"
}`
