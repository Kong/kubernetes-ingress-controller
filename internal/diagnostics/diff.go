package diagnostics

import (
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/kong/go-database-reconciler/pkg/diff"
)

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
	Hash      string       `json:"hash"`
	Entities  []EntityDiff `json:"entities"`
	Timestamp string       `json:"timestamp"`
}

type EntityDiff struct {
	//Source    sourceResource  `json:"kubernetesResource"`
	Generated generatedEntity `json:"kongEntity"`
	Action    string          `json:"action"`
	Diff      string          `json:"diff,omitempty"`
}

// NewEntityDiff creates a diagnostic entity diff.
func NewEntityDiff(diff string, action string, entity diff.Entity) EntityDiff {
	return EntityDiff{
		// TODO this is mostly a stub at present. Need to either derive the source from tags or just omit it for now with
		// a nice to have feature issue, or a simpler YAGNI but if someone asks add it TODO here.
		//Source: sourceResource{},
		Generated: generatedEntity{
			Name: entity.Name,
			Kind: entity.Kind,
		},
		Action: action,
		Diff:   diff,
	}
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

// diffMap holds DB mode diff history.
type diffMap struct {
	diffs     map[string]ConfigDiff
	hashQueue *queue.Queue
	length    int
	times     map[string]string
	latest    string
}

func newDiffMap(length int) diffMap {
	return diffMap{
		diffs:     map[string]ConfigDiff{},
		times:     map[string]string{},
		length:    length,
		hashQueue: queue.New(),
	}
}

// Update adds a diff to the diffMap. If the diffMap holds the maximum number of diffs in history, it removes the
// oldest diff.
func (d *diffMap) Update(diff ConfigDiff) {
	if d.hashQueue.Len() == d.length {
		oldest := d.hashQueue.Dequeue().(string)
		delete(d.diffs, oldest)
	}
	d.hashQueue.Enqueue(diff.Hash)
	d.diffs[diff.Hash] = diff
	d.times[diff.Hash] = diff.Timestamp
	d.latest = diff.Hash
	return
}

// Latest returns the newest diff hash.
func (d *diffMap) Latest() string {
	return d.latest
}

// ByHash returns the diff array matching the given hash.
func (d *diffMap) ByHash(hash string) ([]EntityDiff, error) {
	if diff, ok := d.diffs[hash]; ok {
		return diff.Entities, nil
	}
	return []EntityDiff{}, fmt.Errorf("no diff found for hash %s", hash)
}

// TimeByHash returns the diff timestamp matching the given hash.
func (d *diffMap) TimeByHash(hash string) string {
	if time, ok := d.times[hash]; ok {
		return time
	}
	return "not found"
}

// DiffIndex maps a hash to its timestamp.
type DiffIndex struct {
	// ConfigHash is the config hash for the associated diff.
	ConfigHash string `json:"hash"`
	Timestamp  string `json:"timestamp"`
}

// Available returns a list of cached diff hashes and their associated timestamps.
func (d *diffMap) Available() []DiffIndex {
	index := []DiffIndex{}
	for hash, diff := range d.diffs {
		index = append(index, DiffIndex{ConfigHash: hash, Timestamp: diff.Timestamp})
	}
	return index
}

// Len returns the number of cached diffs.
func (d *diffMap) Len() int {
	return len(d.diffs)
}
