package diagnostics

// TRR TODO this holds guts for diff endpoints. need to reorg the types.go into something that splits out the
// subcomponents into separate files and move the base stuff into server.go probably.
// the /config/{successful|failed} API endpoint structure isn't great for anything other than "current failure, last
// succeeded" state (which is something we still care about). simplest option to handle the diffs within that landscape
// is to have /config/{successful|failed}/diff, but that's a bit stuck on associating the correct one via the hashes.
// we do now _have_ the hashes, but getting the config and diff separate and linking them together after in the
// success/fail structure is maybe a bit wonky. a chronological sequence where the latest maybe has a failure attempt
// attached is maybe simpler, but doesn't _quite_ solve the binding problem--would need to convert the hashes into
// counters somehow

// TRR upstream type
//// EntityAction describes an entity processed by the diff engine and the action taken on it.
//type EntityAction struct {
//	// Action is the ReconcileAction taken on the entity.
//	Action ReconcileAction `json:"action"` // string
//	// Entity holds the processed entity.
//	Entity Entity `json:"entity"`
//	// Diff is diff string describing the modifications made to an entity.
//	Diff string `json:"-"`
//	// Error is the error encountered processing and entity, if any.
//	Error error `json:"error,omitempty"`
//}
//
//// Entity is an entity processed by the diff engine.
//type Entity struct {
//	// Name is the name of the entity.
//	Name string `json:"name"`
//	// Kind is the type of entity.
//	Kind string `json:"kind"`
//	// Old is the original entity in the current state, if any.
//	Old any `json:"old,omitempty"`
//	// New is the new entity in the target state, if any.
//	New any `json:"new,omitempty"`
//}

type sourceResource struct {
	Name             string
	Namespace        string
	GroupVersionKind string
	UID              string
}

type generatedEntity struct {
	// Name is the name of the entity.
	Name string `json:"name"`
	// Kind is the type of entity.
	Kind string `json:"kind"`
}

type ConfigDiff struct {
	Hash     string       `json:"hash"`
	Entities []EntityDiff `json:"entities"`
}

type EntityDiff struct {
	Source    sourceResource  `json:"kubernetesResource"`
	Generated generatedEntity `json:"kongEntity"`
	Action    string          `json:"action"`
	Diff      string          `json:"diff,omitempty"`
}

// TRR TODO this is stolen from the error event builder, which parses regurgitated entity tags into k8s parents.
// we want the same here, minus the additional error info. could probably make it a function in
// internal/util/k8s.go along with sourceResource, but for now it's just sitting here for reference.
// the end function probably won't keep the GVK entries separate? dunno--we want that for the event builder, but for
// the diag server we just want to get the string version. can probably use the upstream GVK type in the generic
// function and call its String().

//func parseTags([]string) {
//	for _, tag := range raw.Tags {
//		if strings.HasPrefix(tag, util.K8sNameTagPrefix) {
//			re.Name = strings.TrimPrefix(tag, util.K8sNameTagPrefix)
//		}
//		if strings.HasPrefix(tag, util.K8sNamespaceTagPrefix) {
//			re.Namespace = strings.TrimPrefix(tag, util.K8sNamespaceTagPrefix)
//		}
//		if strings.HasPrefix(tag, util.K8sKindTagPrefix) {
//			gvk.Kind = strings.TrimPrefix(tag, util.K8sKindTagPrefix)
//		}
//		if strings.HasPrefix(tag, util.K8sVersionTagPrefix) {
//			gvk.Version = strings.TrimPrefix(tag, util.K8sVersionTagPrefix)
//		}
//		// this will not set anything for core resources
//		if strings.HasPrefix(tag, util.K8sGroupTagPrefix) {
//			gvk.Group = strings.TrimPrefix(tag, util.K8sGroupTagPrefix)
//		}
//		if strings.HasPrefix(tag, util.K8sUIDTagPrefix) {
//			re.UID = strings.TrimPrefix(tag, util.K8sUIDTagPrefix)
//		}
//	}
//}
