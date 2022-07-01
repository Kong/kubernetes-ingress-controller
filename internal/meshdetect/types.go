package meshdetect

// DeploymentResults is the result of detecting signals of whether a certain
// mesh is deployed in kubernetes cluster.
type DeploymentResults struct {
	ServiceExists bool `json:"serviceExists"`
}

// RunUnderResults is the result of detecting signals of whether KIC is
// running under a certain service mesh.
type RunUnderResults struct {
	PodOrServiceAnnotation   bool `json:"podOrServiceAnnotation"`
	SidecarContainerInjected bool `json:"sidecarContainerInjected"`
	InitContainerInjected    bool `json:"initContainerInjected"`
}

// ServiceDistributionResults contains number of total services and number of
// services running under each mesh.
type ServiceDistributionResults struct {
	TotalServices int `json:"totalServices"`
	// MeshDistribution is the number of services running under each kind of mesh.
	// We decided to directly use number here instead of ratio in total services in a
	// floating number because using floating number needs extra work on calculating
	// and serializing.
	MeshDistribution map[MeshKind]int `json:"meshDistribution"`
}
