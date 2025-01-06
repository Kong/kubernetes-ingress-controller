package license

type ListLicenseResponse struct {
	Items []*Item `json:"items"`
	// TODO our APIs generally assume that there are no unary objects. Any object type can have multiple instances,
	// and lists of instances can be paginated. However, the license API doesn't return pagination info, as it is
	// effectively a unary object. We should sort that out, to at least have a guarantee as to whether or not we'll
	// represent unary objects as a collection that coincidentally always only has one page with one entry.
	// Page  *PaginationInfo `json:"page"`
}

// Item is a single license from the upstream license API.
type Item struct {
	License   string `json:"payload,omitempty"`
	UpdatedAt uint64 `json:"updated_at,omitempty"`
	ID        string `json:"id,omitempty"`
}
