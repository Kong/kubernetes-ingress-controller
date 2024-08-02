package diagnostics

import (
	"fmt"

	"github.com/golang-collections/collections/queue"
	"github.com/kong/go-database-reconciler/pkg/diff"
)

// generatedEntity is basic name and type metadata about a Kong gateway entity.
type generatedEntity struct {
	// Name is the name of the entity.
	Name string `json:"name"`
	// Kind is the type of entity.
	Kind string `json:"kind"`
}

// ConfigDiff holds a config update, including its config hash, the rough timestamp when the controller completed
// sending it to the gateway, and the entities that changed in the course of reconciling the config state.
type ConfigDiff struct {
	Hash      string       `json:"hash"`
	Entities  []EntityDiff `json:"entities"`
	Timestamp string       `json:"timestamp"`
}

// EntityDiff is an individual entity change. It includes the entity metadata, the action performed during
// reconciliation, and the diff string for update actions.
type EntityDiff struct {
	Generated generatedEntity `json:"kongEntity"`
	Action    string          `json:"action"`
	Diff      string          `json:"diff,omitempty"`
}

// NewEntityDiff creates a diagnostic entity diff.
func NewEntityDiff(diff string, action string, entity diff.Entity) EntityDiff {
	return EntityDiff{
		Generated: generatedEntity{
			Name: entity.Name,
			Kind: entity.Kind,
		},
		Action: action,
		Diff:   diff,
	}
}

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
