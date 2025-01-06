package kongconfig

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	dpconf "github.com/kong/kubernetes-ingress-controller/v3/internal/dataplane/config"
	"github.com/kong/kubernetes-ingress-controller/v3/internal/versions"
)

func TestRoot(t *testing.T) {
	var root Root
	require.NoError(t, json.Unmarshal([]byte(dblessConfigJSON3_4_1), &root))
	require.NoError(t, root.Validate(false))
	require.EqualError(t, root.Validate(true), "--skip-ca-certificates is not available for use with DB-less Kong instances")
}

func TestValidateRoots(t *testing.T) {
	testCases := []struct {
		name                 string
		configStr            string
		expectedDBMode       dpconf.DBMode
		expectedRouterFlavor dpconf.RouterFlavor
		expectedKongVersion  string
	}{
		{
			name:                 "dbless config with version 3.4.1",
			configStr:            dblessConfigJSON3_4_1,
			expectedDBMode:       dpconf.DBModeOff,
			expectedRouterFlavor: dpconf.RouterFlavorTraditionalCompatible,
			expectedKongVersion:  versions.KICv3VersionCutoff.String(),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var root Root
			require.NoError(t, json.Unmarshal([]byte(tc.configStr), &root))
			kongOptions, err := ValidateRoots([]Root{root, root}, false)
			require.NoError(t, err)
			assert.Equal(t, tc.expectedDBMode, kongOptions.DBMode)
			assert.Equal(t, tc.expectedRouterFlavor, kongOptions.RouterFlavor)
			assert.Equal(t, tc.expectedKongVersion, kongOptions.Version.String())
		})
	}
}

const dblessConfigJSON3_4_1 = `
{
	"node_id": "69d063c5-761b-4bab-a426-c89da49a9409",
	"lua_version": "LuaJIT 2.1.0-20220411",
	"hostname": "ingress-controller-kong-579b457597-2vcq8",
	"plugins": {
	  "available_on_server": {
		"prometheus": {
		  "version": "3.4.1",
		  "priority": 13
		},
		"proxy-cache": {
		  "version": "3.4.1",
		  "priority": 100
		},
		"session": {
		  "version": "3.4.1",
		  "priority": 1900
		},
		"acme": {
		  "version": "3.4.1",
		  "priority": 1705
		},
		"grpc-gateway": {
		  "version": "3.4.1",
		  "priority": 998
		},
		"grpc-web": {
		  "version": "3.4.1",
		  "priority": 3
		},
		"pre-function": {
		  "version": "3.4.1",
		  "priority": 1000000
		},
		"post-function": {
		  "version": "3.4.1",
		  "priority": -1000
		},
		"azure-functions": {
		  "version": "3.4.1",
		  "priority": 749
		},
		"zipkin": {
		  "version": "3.4.1",
		  "priority": 100000
		},
		"opentelemetry": {
		  "version": "0.1.0",
		  "priority": 14
		},
		"jwt": {
		  "version": "3.4.1",
		  "priority": 1450
		},
		"acl": {
		  "version": "3.4.1",
		  "priority": 950
		},
		"correlation-id": {
		  "version": "3.4.1",
		  "priority": 1
		},
		"cors": {
		  "version": "3.4.1",
		  "priority": 2000
		},
		"oauth2": {
		  "version": "3.4.1",
		  "priority": 1400
		},
		"tcp-log": {
		  "version": "3.4.1",
		  "priority": 7
		},
		"udp-log": {
		  "version": "3.4.1",
		  "priority": 8
		},
		"file-log": {
		  "version": "3.4.1",
		  "priority": 9
		},
		"http-log": {
		  "version": "3.4.1",
		  "priority": 12
		},
		"key-auth": {
		  "version": "3.4.1",
		  "priority": 1250
		},
		"hmac-auth": {
		  "version": "3.4.1",
		  "priority": 1030
		},
		"basic-auth": {
		  "version": "3.4.1",
		  "priority": 1100
		},
		"ip-restriction": {
		  "version": "3.4.1",
		  "priority": 990
		},
		"request-transformer": {
		  "version": "3.4.1",
		  "priority": 801
		},
		"response-transformer": {
		  "version": "3.4.1",
		  "priority": 800
		},
		"request-size-limiting": {
		  "version": "3.4.1",
		  "priority": 951
		},
		"rate-limiting": {
		  "version": "3.4.1",
		  "priority": 910
		},
		"response-ratelimiting": {
		  "version": "3.4.1",
		  "priority": 900
		},
		"syslog": {
		  "version": "3.4.1",
		  "priority": 4
		},
		"loggly": {
		  "version": "3.4.1",
		  "priority": 6
		},
		"datadog": {
		  "version": "3.4.1",
		  "priority": 10
		},
		"ldap-auth": {
		  "version": "3.4.1",
		  "priority": 1200
		},
		"statsd": {
		  "version": "3.4.1",
		  "priority": 11
		},
		"bot-detection": {
		  "version": "3.4.1",
		  "priority": 2500
		},
		"aws-lambda": {
		  "version": "3.4.1",
		  "priority": 750
		},
		"request-termination": {
		  "version": "3.4.1",
		  "priority": 2
		}
	  },
	  "enabled_in_cluster": []
	},
	"tagline": "Welcome to kong",
	"version": "3.4.1",
	"pids": {
	  "workers": [
		1272,
		1273
	  ],
	  "master": 1
	},
	"timers": {
	  "pending": 1,
	  "running": 79
	},
	"configuration": {
	  "log_level": "notice",
	  "admin_gui_path": "/",
	  "worker_events_max_payload": 65535,
	  "admin_gui_listeners": [
		{
		  "listener": "0.0.0.0:8002",
		  "port": 8002,
		  "ip": "0.0.0.0",
		  "http2": false,
		  "proxy_protocol": false,
		  "deferred": false,
		  "reuseport": false,
		  "backlog=%d+": false,
		  "ssl": false,
		  "ipv6only=off": false,
		  "so_keepalive=on": false,
		  "bind": false,
		  "so_keepalive=%w*:%w*:%d*": false,
		  "so_keepalive=off": false,
		  "ipv6only=on": false
		},
		{
		  "listener": "0.0.0.0:8445 ssl",
		  "port": 8445,
		  "ip": "0.0.0.0",
		  "http2": false,
		  "proxy_protocol": false,
		  "deferred": false,
		  "reuseport": false,
		  "backlog=%d+": false,
		  "ssl": true,
		  "ipv6only=off": false,
		  "so_keepalive=on": false,
		  "bind": false,
		  "so_keepalive=%w*:%w*:%d*": false,
		  "so_keepalive=off": false,
		  "ipv6only=on": false
		}
	  ],
	  "db_update_frequency": 5,
	  "db_update_propagation": 0,
	  "kic": false,
	  "ssl_cert": [
		"/kong_prefix/ssl/kong-default.crt",
		"/kong_prefix/ssl/kong-default-ecdsa.crt"
	  ],
	  "lua_max_req_headers": 100,
	  "lua_package_path": "/opt/?.lua;/opt/?/init.lua;;",
	  "host_ports": {
		"8443": 443,
		"8443": 443,
		"8000": 80,
		"8000": 80
	  },
	  "lua_ssl_trusted_certificate": [
		"/etc/ssl/certs/ca-certificates.crt"
	  ],
	  "lua_ssl_verify_depth": 1,
	  "lua_max_uri_args": 100,
	  "worker_state_update_frequency": 5,
	  "lua_max_post_args": 100,
	  "tracing_instrumentations": [
		"off"
	  ],
	  "tracing_sampling_rate": 0.01,
	  "enabled_headers": {
		"X-Kong-Proxy-Latency": true,
		"Via": true,
		"X-Kong-Response-Latency": true,
		"X-Kong-Admin-Latency": true,
		"X-Kong-Upstream-Latency": true,
		"server_tokens": true,
		"X-Kong-Upstream-Status": false,
		"Server": true,
		"latency_tokens": true
	  },
	  "cluster_listeners": {},
	  "status_ssl_enabled": false,
	  "nginx_events_multi_accept": "on",
	  "ssl_cert_csr_default": "/kong_prefix/ssl/kong-default.csr",
	  "admin_ssl_enabled": false,
	  "upstream_keepalive_pool_size": 512,
	  "ssl_cert_default": "/kong_prefix/ssl/kong-default.crt",
	  "stream_listeners": [
		{
		  "listener": "0.0.0.0:8888",
		  "port": 8888,
		  "ip": "0.0.0.0",
		  "proxy_protocol": false,
		  "so_keepalive=%w*:%w*:%d*": false,
		  "reuseport": false,
		  "backlog=%d+": false,
		  "ssl": false,
		  "ipv6only=off": false,
		  "so_keepalive=on": false,
		  "bind": false,
		  "udp": false,
		  "so_keepalive=off": false,
		  "ipv6only=on": false
		},
		{
		  "listener": "0.0.0.0:8899 ssl reuseport",
		  "port": 8899,
		  "ip": "0.0.0.0",
		  "proxy_protocol": false,
		  "so_keepalive=%w*:%w*:%d*": false,
		  "reuseport": true,
		  "backlog=%d+": false,
		  "ssl": true,
		  "ipv6only=off": false,
		  "so_keepalive=on": false,
		  "bind": false,
		  "udp": false,
		  "so_keepalive=off": false,
		  "ipv6only=on": false
		},
		{
		  "listener": "0.0.0.0:9999 udp reuseport",
		  "port": 9999,
		  "ip": "0.0.0.0",
		  "proxy_protocol": false,
		  "so_keepalive=%w*:%w*:%d*": false,
		  "reuseport": true,
		  "backlog=%d+": false,
		  "ssl": false,
		  "ipv6only=off": false,
		  "so_keepalive=on": false,
		  "bind": false,
		  "udp": true,
		  "so_keepalive=off": false,
		  "ipv6only=on": false
		}
	  ],
	  "proxy_ssl_enabled": true,
	  "cluster_data_plane_purge_delay": 1209600,
	  "db_cache_ttl": 0,
	  "nginx_wasm_main_directives": {},
	  "nginx_wasm_main_shm_directives": {},
	  "nginx_wasm_wasmer_directives": {},
	  "nginx_wasm_v8_directives": {},
	  "worker_consistency": "eventual",
	  "dns_resolver": {},
	  "nginx_wasm_wasmtime_directives": {},
	  "dns_hostsfile": "/etc/hosts",
	  "nginx_sproxy_directives": {},
	  "nginx_supstream_directives": {},
	  "nginx_stream_directives": [
		{
		  "name": "lua_shared_dict",
		  "value": "stream_prometheus_metrics 5m"
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
	  "nginx_http_lua_ssl_protocols": "TLSv1.1 TLSv1.2 TLSv1.3",
	  "nginx_stream_lua_ssl_protocols": "TLSv1.1 TLSv1.2 TLSv1.3",
	  "dns_error_ttl": 1,
	  "cluster_control_plane": "127.0.0.1:8005",
	  "dns_not_found_ttl": 30,
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
	  "dns_stale_ttl": 4,
	  "nginx_status_directives": {},
	  "dns_cache_size": 10000,
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
	  "dns_order": [
		"LAST",
		"SRV",
		"A",
		"CNAME"
	  ],
	  "ssl_cert_key": "******",
	  "dns_no_sync": false,
	  "cluster_use_proxy": false,
	  "cluster_dp_labels": {},
	  "untrusted_lua": "sandbox",
	  "nginx_kong_stream_inject_conf": "/kong_prefix/nginx-kong-stream-inject.conf",
	  "untrusted_lua_sandbox_requires": {},
	  "untrusted_lua_sandbox_environment": {},
	  "lmdb_environment_path": "dbless.lmdb",
	  "lmdb_map_size": "2048m",
	  "opentelemetry_tracing": [
		"off"
	  ],
	  "opentelemetry_tracing_sampling_rate": 0.01,
	  "pluginserver_names": {},
	  "proxy_server_ssl_verify": true,
	  "ssl_cert_key_default": "/kong_prefix/ssl/kong-default.key",
	  "nginx_upstream_directives": {},
	  "ssl_cert_default_ecdsa": "/kong_prefix/ssl/kong-default-ecdsa.crt",
	  "nginx_kong_stream_conf": "/kong_prefix/nginx-kong-stream.conf",
	  "ssl_cert_key_default_ecdsa": "/kong_prefix/ssl/kong-default-ecdsa.key",
	  "ssl_cipher_suite": "intermediate",
	  "client_ssl_cert_default": "/kong_prefix/ssl/kong-default.crt",
	  "client_ssl_cert_key_default": "/kong_prefix/ssl/kong-default.key",
	  "db_cache_warmup_entities": [
		"services"
	  ],
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
	  "admin_ssl_cert_key_default": "/kong_prefix/ssl/admin-kong-default.key",
	  "nginx_kong_gui_include_conf": "/kong_prefix/nginx-kong-gui-include.conf",
	  "admin_ssl_cert_default_ecdsa": "/kong_prefix/ssl/admin-kong-default-ecdsa.crt",
	  "nginx_conf": "/kong_prefix/nginx.conf",
	  "headers": [
		"server_tokens",
		"latency_tokens"
	  ],
	  "lua_ssl_trusted_certificate_combined": "/kong_prefix/.ca_combined",
	  "status_ssl_cert_default": "/kong_prefix/ssl/status-kong-default.crt",
	  "nginx_main_worker_rlimit_nofile": "auto",
	  "status_ssl_cert_key_default": "/kong_prefix/ssl/status-kong-default.key",
	  "nginx_events_worker_connections": "auto",
	  "status_ssl_cert_default_ecdsa": "/kong_prefix/ssl/status-kong-default-ecdsa.crt",
	  "prefix": "/kong_prefix",
	  "client_ssl": false,
	  "nginx_http_charset": "UTF-8",
	  "router_flavor": "traditional_compatible",
	  "nginx_http_client_max_body_size": "0",
	  "admin_gui_ssl_cert_key_default": "/kong_prefix/ssl/admin-gui-kong-default.key",
	  "nginx_http_client_body_buffer_size": "8k",
	  "admin_gui_ssl_cert_default_ecdsa": "/kong_prefix/ssl/admin-gui-kong-default-ecdsa.crt",
	  "admin_acc_logs": "/kong_prefix/logs/admin_access.log",
	  "admin_gui_ssl_cert_key_default_ecdsa": "/kong_prefix/ssl/admin-gui-kong-default-ecdsa.key",
	  "admin_ssl_cert_default": "/kong_prefix/ssl/admin-kong-default.crt",
	  "nginx_acc_logs": "/kong_prefix/logs/access.log",
	  "pg_ro_ssl_verify": false,
	  "port_maps": [
		"80:8000",
		"443:8443"
	  ],
	  "proxy_listen": [
		"0.0.0.0:8000",
		"0.0.0.0:8443 http2 ssl"
	  ],
	  "admin_listen": [
		"0.0.0.0:8001"
	  ],
	  "admin_gui_listen": [
		"0.0.0.0:8002",
		"0.0.0.0:8445 ssl"
	  ],
	  "status_listen": [
		"0.0.0.0:8100"
	  ],
	  "stream_listen": [
		"0.0.0.0:8888",
		"0.0.0.0:8899 ssl reuseport",
		"0.0.0.0:9999 udp reuseport"
	  ],
	  "cluster_listen": [
		"off"
	  ],
	  "admin_ssl_cert": {},
	  "admin_ssl_cert_key": "******",
	  "admin_gui_ssl_cert": [
		"/kong_prefix/ssl/admin-gui-kong-default.crt",
		"/kong_prefix/ssl/admin-gui-kong-default-ecdsa.crt"
	  ],
	  "admin_gui_ssl_cert_key": "******",
	  "status_ssl_cert": {},
	  "status_ssl_cert_key": "******",
	  "db_resurrect_ttl": 30,
	  "nginx_user": "kong kong",
	  "nginx_http_keepalive_requests": "1000",
	  "nginx_main_user": "kong kong",
	  "nginx_daemon": "off",
	  "mem_cache_size": "128m",
	  "nginx_main_daemon": "off",
	  "nginx_worker_processes": "2",
	  "nginx_main_worker_processes": "2",
	  "trusted_ips": {},
	  "real_ip_header": "X-Real-IP",
	  "nginx_proxy_real_ip_header": "X-Real-IP",
	  "real_ip_recursive": "off",
	  "nginx_proxy_real_ip_recursive": "off",
	  "stream_proxy_ssl_enabled": true,
	  "pg_port": 5432,
	  "admin_gui_ssl_enabled": true,
	  "pg_ssl": false,
	  "pg_ssl_verify": false,
	  "pg_max_concurrent_queries": 0,
	  "pg_semaphore_timeout": 60000,
	  "kong_process_secrets": "/kong_prefix/.kong_process_secrets",
	  "kong_env": "/kong_prefix/.kong_env",
	  "allow_debug_header": false,
	  "_debug_pg_ttl_cleanup_interval": 300,
	  "admin_gui_ssl_cert_default": "/kong_prefix/ssl/admin-gui-kong-default.crt",
	  "role": "traditional",
	  "status_ssl_cert_key_default_ecdsa": "/kong_prefix/ssl/status-kong-default-ecdsa.key",
	  "pg_ro_ssl": false,
	  "database": "off",
	  "nginx_kong_inject_conf": "/kong_prefix/nginx-kong-inject.conf",
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
	  "nginx_inject_conf": "/kong_prefix/nginx-inject.conf",
	  "nginx_http_directives": [
		{
		  "name": "charset",
		  "value": "UTF-8"
		},
		{
		  "name": "client_body_buffer_size",
		  "value": "8k"
		},
		{
		  "name": "client_max_body_size",
		  "value": "0"
		},
		{
		  "name": "keepalive_requests",
		  "value": "1000"
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
	  "nginx_kong_conf": "/kong_prefix/nginx-kong.conf",
	  "admin_ssl_cert_key_default_ecdsa": "/kong_prefix/ssl/admin-kong-default-ecdsa.key",
	  "loaded_vaults": {
		"env": true
	  },
	  "lua_ssl_protocols": "TLSv1.1 TLSv1.2 TLSv1.3",
	  "ssl_protocols": "TLSv1.1 TLSv1.2 TLSv1.3",
	  "anonymous_reports": true,
	  "ssl_ciphers": "ECDHE-ECDSA-AES128-GCM-SHA256:ECDHE-RSA-AES128-GCM-SHA256:ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384:ECDHE-ECDSA-CHACHA20-POLY1305:ECDHE-RSA-CHACHA20-POLY1305:DHE-RSA-AES128-GCM-SHA256:DHE-RSA-AES256-GCM-SHA384",
	  "client_body_buffer_size": "8k",
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
	  "ssl_session_cache_size": "10m",
	  "proxy_access_log": "/dev/stdout",
	  "proxy_error_log": "/dev/stderr",
	  "proxy_stream_access_log": "/dev/stdout basic",
	  "proxy_stream_error_log": "/dev/stderr",
	  "admin_access_log": "/dev/stdout",
	  "admin_error_log": "/dev/stderr",
	  "admin_gui_access_log": "/dev/stdout",
	  "admin_gui_error_log": "/dev/stderr",
	  "status_access_log": "off",
	  "status_error_log": "/dev/stderr",
	  "nginx_pid": "/kong_prefix/pids/nginx.pid",
	  "nginx_http_lua_regex_cache_max_entries": "8192",
	  "cluster_ocsp": "off",
	  "nginx_err_logs": "/kong_prefix/logs/error.log",
	  "plugins": [
		"bundled"
	  ],
	  "upstream_keepalive_idle_timeout": 60,
	  "pg_database": "kong",
	  "nginx_admin_client_body_buffer_size": "10m",
	  "nginx_http_ssl_protocols": "TLSv1.2 TLSv1.3",
	  "nginx_stream_ssl_protocols": "TLSv1.2 TLSv1.3",
	  "vaults": [
		"bundled"
	  ],
	  "lua_package_cpath": "",
	  "pg_user": "kong",
	  "error_default_type": "text/plain",
	  "loaded_plugins": {
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
		"hmac-auth": true,
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
		"request-termination": true
	  },
	  "pg_host": "127.0.0.1",
	  "upstream_keepalive_max_requests": 1000,
	  "proxy_listeners": [
		{
		  "listener": "0.0.0.0:8000",
		  "port": 8000,
		  "ip": "0.0.0.0",
		  "http2": false,
		  "proxy_protocol": false,
		  "deferred": false,
		  "reuseport": false,
		  "backlog=%d+": false,
		  "ssl": false,
		  "ipv6only=off": false,
		  "so_keepalive=on": false,
		  "bind": false,
		  "so_keepalive=%w*:%w*:%d*": false,
		  "so_keepalive=off": false,
		  "ipv6only=on": false
		},
		{
		  "listener": "0.0.0.0:8443 ssl http2",
		  "port": 8443,
		  "ip": "0.0.0.0",
		  "http2": true,
		  "proxy_protocol": false,
		  "deferred": false,
		  "reuseport": false,
		  "backlog=%d+": false,
		  "ssl": true,
		  "ipv6only=off": false,
		  "so_keepalive=on": false,
		  "bind": false,
		  "so_keepalive=%w*:%w*:%d*": false,
		  "so_keepalive=off": false,
		  "ipv6only=on": false
		}
	  ],
	  "admin_listeners": [
		{
		  "listener": "0.0.0.0:8001",
		  "port": 8001,
		  "ip": "0.0.0.0",
		  "http2": false,
		  "proxy_protocol": false,
		  "deferred": false,
		  "reuseport": false,
		  "backlog=%d+": false,
		  "ssl": false,
		  "ipv6only=off": false,
		  "so_keepalive=on": false,
		  "bind": false,
		  "so_keepalive=%w*:%w*:%d*": false,
		  "so_keepalive=off": false,
		  "ipv6only=on": false
		}
	  ],
	  "status_listeners": [
		{
		  "listener": "0.0.0.0:8100",
		  "port": 8100,
		  "ip": "0.0.0.0",
		  "http2": false,
		  "proxy_protocol": false,
		  "deferred": false,
		  "reuseport": false,
		  "backlog=%d+": false,
		  "ssl": false,
		  "ipv6only=off": false,
		  "so_keepalive=on": false,
		  "bind": false,
		  "so_keepalive=%w*:%w*:%d*": false,
		  "so_keepalive=off": false,
		  "ipv6only=on": false
		}
	  ],
	  "pg_timeout": 5000,
	  "lua_socket_pool_size": 30,
	  "nginx_http_lua_regex_match_limit": "100000",
	  "wasm": false,
	  "cluster_max_payload": 16777216,
	  "privileged_agent": false,
	  "cluster_mtls": "shared",
	  "nginx_admin_client_max_body_size": "10m",
	  "lua_max_resp_headers": 100
	}
  }
`
