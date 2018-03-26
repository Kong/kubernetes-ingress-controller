package admin

import "fmt"

// APIResponse contains a response returned from the Kong admin API
// +k8s:deepcopy-gen=false
type APIResponse struct {
	err        error
	StatusCode int
	Raw        []byte
}

// Error returns the error from the admin API response
func (r *APIResponse) Error() error {
	return r.err
}

func (r *APIResponse) String() string {
	if r.Raw == nil && r.StatusCode == 0 {
		return r.err.Error()
	}
	if r.Raw != nil {
		return fmt.Sprintf("[%d] %s", r.StatusCode, string(r.Raw))
	}
	return fmt.Sprintf("[%d] %s", r.StatusCode, r.err)
}
