package credentials

// -----------------------------------------------------------------------------
// Private
// -----------------------------------------------------------------------------

// uniqueKeyConstraints a map of unique key constraints for any given credential type.
// This map is the crux of all unique key constraint validation and is derived from
// the relevant Lua code for the types in the backend Kong Admin API.
//
// Example: https://github.com/kong/kong/blob/master/kong/plugins/basic-auth/daos.lua
//
// So if you're in here doing maintenance and you need to add/remove constraints due
// to upstream changes, check the "kong/plugins" directory of the upstream repo
// for every given type.
var uniqueKeyConstraints = map[string][]string{
	"basic-auth": {"username"},
	"hmac-auth":  {"username"},
	"jwt":        {"key"},
	"key-auth":   {"key"},
	"oauth2":     {"client_id"},
}
