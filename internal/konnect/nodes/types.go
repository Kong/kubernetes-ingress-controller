package nodes

const (
	// NodeTypeIngressController is the type of nodes representing KIC instances.
	NodeTypeIngressController = "ingress-controller"
	// NodeTypeKongProxy is the type of nodes representing kong gateway instances.
	NodeTypeKongProxy = "kong-proxy"
)

type NodeItem struct {
	ID                  string               `json:"id"`
	Version             string               `json:"version"`
	Hostname            string               `json:"hostname"`
	LastPing            int64                `json:"last_ping"`
	Type                string               `json:"type"`
	CreatedAt           int64                `json:"created_at"`
	UpdatedAt           int64                `json:"updated_at"`
	ConfigHash          string               `json:"config_hash"`
	CompatibilityStatus *CompatibilityStatus `json:"compatibility_status,omitempty"`
	Status              string               `json:"status,omitempty"`
}

type CompatibilityState string

const (
	CompatibilityStateUnspecified     CompatibilityState = "COMPATIBILITY_STATE_UNSPECIFIED"
	CompatibilityStateFullyCompatible CompatibilityState = "COMPATIBILITY_STATE_FULLY_COMPATIBLE"
	CompatibilityStateInconpatible    CompatibilityState = "COMPATIBILITY_STATE_INCOMPATIBLE"
	CompatibilityStateUnknown         CompatibilityState = "COMPATIBILITY_STATE_UNKNOWN"
)

type KongResource struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

type CompatibilityIssue struct {
	Code              string          `json:"code"`
	Severity          string          `json:"severity"`
	Description       string          `json:"description"`
	Resolution        string          `json:"resolution"`
	AffectedResources []*KongResource `json:"affected_resources"`
}

type CompatibilityStatus struct {
	State  CompatibilityState    `json:"state"`
	Issues []*CompatibilityIssue `json:"issues,omitempty"`
}

type IngressControllerState string

const (
	IngressControllerStateUnspecified                       IngressControllerState = "INGRESS_CONTROLLER_STATE_UNSPECIFIED"
	IngressControllerStateOperational                       IngressControllerState = "INGRESS_CONTROLLER_STATE_OPERATIONAL"
	IngressControllerStatePartialConfigFail                 IngressControllerState = "INGRESS_CONTROLLER_STATE_PARTIAL_CONFIG_FAIL"
	IngressControllerStateInoperable                        IngressControllerState = "INGRESS_CONTROLLER_STATE_INOPERABLE"
	IngressControllerStateOperationalKonnectOutOfSync       IngressControllerState = "INGRESS_CONTROLLER_STATE_OPERATIONAL_KONNECT_OUT_OF_SYNC"
	IngressControllerStatePartialConfigFailKonnectOutOfSync IngressControllerState = "INGRESS_CONTROLLER_STATE_PARTIAL_CONFIG_FAIL_KONNECT_OUT_OF_SYNC"
	IngressControllerStateInoperableKonnectOutOfSync        IngressControllerState = "INGRESS_CONTROLLER_STATE_INOPERABLE_KONNECT_OUT_OF_SYNC"
	IngressControllerStateUnknown                           IngressControllerState = "INGRESS_CONTROLLER_STATE_UNKNOWN"
)

type CreateNodeRequest struct {
	ID                  string               `json:"id,omitempty"`
	Hostname            string               `json:"hostname"`
	Type                string               `json:"type"`
	LastPing            int64                `json:"last_ping"`
	Version             string               `json:"version"`
	CompatabilityStatus *CompatibilityStatus `json:"compatibility_status,omitempty"`
	Status              string               `json:"status,omitempty"`
	ConfigHash          string               `json:"config_hash,omitempty"`
}

type CreateNodeResponse struct {
	Item *NodeItem `json:"item"`
}

type UpdateNodeRequest struct {
	Hostname            string               `json:"hostname"`
	Type                string               `json:"type"`
	LastPing            int64                `json:"last_ping"`
	Version             string               `json:"version"`
	ConfigHash          string               `json:"config_hash,omitempty"`
	CompatabilityStatus *CompatibilityStatus `json:"compatibility_status,omitempty"`
	Status              string               `json:"status,omitempty"`
}

type UpdateNodeResponse struct {
	Item *NodeItem `json:"item"`
}

type ListNodeResponse struct {
	Items []*NodeItem     `json:"items"`
	Page  *PaginationInfo `json:"page"`
}

type PaginationInfo struct {
	TotalCount  int32  `json:"total_count,omitempty"`
	NextCursor  string `json:"next_cursor,omitempty"`
	HasNextPage bool   `json:"has_next_page,omitempty"`
}
