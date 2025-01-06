package diagnostics

import (
	"container/list"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestDiffMap_Update(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		initial    diffMap
		newDiff    ConfigDiff
		expected   diffMap
		shouldFail bool
	}{
		{
			name: "add new diff to empty diffMap",
			initial: diffMap{
				diffs:     map[string]ConfigDiff{},
				times:     map[string]string{},
				length:    2,
				hashQueue: list.New(),
			},
			newDiff: ConfigDiff{
				Hash:      "hash1",
				Timestamp: now.Format(time.RFC3339),
				Entities:  []EntityDiff{},
			},
			expected: diffMap{
				diffs: map[string]ConfigDiff{
					"hash1": {
						Hash:      "hash1",
						Timestamp: now.Format(time.RFC3339),
						Entities:  []EntityDiff{},
					},
				},
				times: map[string]string{
					"hash1": now.Format(time.RFC3339),
				},
				length:    2,
				hashQueue: list.New(),
				latest:    "hash1",
			},
			shouldFail: false,
		},
		{
			name: "add new diff to non empty diffMap",
			initial: func() diffMap {
				dm := diffMap{
					diffs:     map[string]ConfigDiff{},
					times:     map[string]string{},
					length:    2,
					hashQueue: list.New(),
				}
				dm.Update(ConfigDiff{
					Hash:      "hash1",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				return dm
			}(),
			newDiff: ConfigDiff{
				Hash:      "hash2",
				Timestamp: now.Format(time.RFC3339),
				Entities:  []EntityDiff{},
			},
			expected: diffMap{
				diffs: map[string]ConfigDiff{
					"hash1": {
						Hash:      "hash1",
						Timestamp: now.Format(time.RFC3339),
						Entities:  []EntityDiff{},
					},
					"hash2": {
						Hash:      "hash2",
						Timestamp: now.Format(time.RFC3339),
						Entities:  []EntityDiff{},
					},
				},
				times: map[string]string{
					"hash1": now.Format(time.RFC3339),
					"hash2": now.Format(time.RFC3339),
				},
				length:    2,
				hashQueue: list.New(),
				latest:    "hash2",
			},
			shouldFail: false,
		},
		{
			name: "add new diff to full diffMap",
			initial: func() diffMap {
				dm := diffMap{
					diffs:     map[string]ConfigDiff{},
					times:     map[string]string{},
					length:    2,
					hashQueue: list.New(),
				}
				dm.Update(ConfigDiff{
					Hash:      "hash1",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				dm.Update(ConfigDiff{
					Hash:      "hash2",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				return dm
			}(),
			newDiff: ConfigDiff{
				Hash:      "hash3",
				Timestamp: now.Format(time.RFC3339),
				Entities:  []EntityDiff{},
			},
			expected: func() diffMap {
				dm := diffMap{
					diffs:     map[string]ConfigDiff{},
					times:     map[string]string{},
					length:    2,
					hashQueue: list.New(),
				}
				dm.Update(ConfigDiff{
					Hash:      "hash2",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				dm.Update(ConfigDiff{
					Hash:      "hash3",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				return dm
			}(),
			shouldFail: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.initial.Update(tt.newDiff)
			assert.Equal(t, tt.expected.diffs, tt.initial.diffs)
			assert.Equal(t, tt.expected.times, tt.initial.times)
			assert.Equal(t, tt.expected.length, tt.initial.length)
			assert.Equal(t, tt.expected.latest, tt.initial.latest)
		})
	}
}

func TestDiffMap_Latest(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		initial  diffMap
		expected string
	}{
		{
			name: "latest diff in empty diffMap",
			initial: diffMap{
				diffs:     map[string]ConfigDiff{},
				times:     map[string]string{},
				length:    2,
				hashQueue: list.New(),
			},
			expected: "",
		},
		{
			name: "latest diff in non empty diffMap",
			initial: func() diffMap {
				dm := diffMap{
					diffs:     map[string]ConfigDiff{},
					times:     map[string]string{},
					length:    2,
					hashQueue: list.New(),
				}
				dm.Update(ConfigDiff{
					Hash:      "hash1",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				return dm
			}(),
			expected: "hash1",
		},
		{
			name: "latest diff in full diffMap",
			initial: func() diffMap {
				dm := diffMap{
					diffs:     map[string]ConfigDiff{},
					times:     map[string]string{},
					length:    2,
					hashQueue: list.New(),
				}
				dm.Update(ConfigDiff{
					Hash:      "hash1",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				dm.Update(ConfigDiff{
					Hash:      "hash2",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				return dm
			}(),
			expected: "hash2",
		},
		{
			name: "latest diff after adding to full diffMap",
			initial: func() diffMap {
				dm := diffMap{
					diffs:     map[string]ConfigDiff{},
					times:     map[string]string{},
					length:    2,
					hashQueue: list.New(),
				}
				dm.Update(ConfigDiff{
					Hash:      "hash1",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				dm.Update(ConfigDiff{
					Hash:      "hash2",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				dm.Update(ConfigDiff{
					Hash:      "hash3",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				return dm
			}(),
			expected: "hash3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.initial.Latest())
		})
	}
}

func TestDiffMap_Len(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name     string
		initial  diffMap
		expected int
	}{
		{
			name: "length of empty diffMap",
			initial: diffMap{
				diffs:     map[string]ConfigDiff{},
				times:     map[string]string{},
				length:    2,
				hashQueue: list.New(),
			},
			expected: 0,
		},
		{
			name: "length of non empty diffMap",
			initial: func() diffMap {
				dm := diffMap{
					diffs:     map[string]ConfigDiff{},
					times:     map[string]string{},
					length:    2,
					hashQueue: list.New(),
				}
				dm.Update(ConfigDiff{
					Hash:      "hash1",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				return dm
			}(),
			expected: 1,
		},
		{
			name: "length of full diffMap",
			initial: func() diffMap {
				dm := diffMap{
					diffs:     map[string]ConfigDiff{},
					times:     map[string]string{},
					length:    2,
					hashQueue: list.New(),
				}
				dm.Update(ConfigDiff{
					Hash:      "hash1",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				dm.Update(ConfigDiff{
					Hash:      "hash2",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				return dm
			}(),
			expected: 2,
		},
		{
			name: "length after adding to full diffMap",
			initial: func() diffMap {
				dm := diffMap{
					diffs:     map[string]ConfigDiff{},
					times:     map[string]string{},
					length:    2,
					hashQueue: list.New(),
				}
				dm.Update(ConfigDiff{
					Hash:      "hash1",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				dm.Update(ConfigDiff{
					Hash:      "hash2",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				dm.Update(ConfigDiff{
					Hash:      "hash3",
					Timestamp: now.Format(time.RFC3339),
					Entities:  []EntityDiff{},
				})
				return dm
			}(),
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.initial.Len())
		})
	}
}

func TestDiffMap_ByHash(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name       string
		initial    diffMap
		hash       string
		expected   []EntityDiff
		shouldFail bool
	}{
		{
			name: "retrieve diff by hash from empty diffMap",
			initial: diffMap{
				diffs:     map[string]ConfigDiff{},
				times:     map[string]string{},
				length:    2,
				hashQueue: list.New(),
			},
			hash:       "hash1",
			expected:   []EntityDiff{},
			shouldFail: true,
		},
		{
			name: "retrieve diff by hash from non empty diffMap",
			initial: func() diffMap {
				dm := diffMap{
					diffs:     map[string]ConfigDiff{},
					times:     map[string]string{},
					length:    2,
					hashQueue: list.New(),
				}
				dm.Update(ConfigDiff{
					Hash:      "hash1",
					Timestamp: now.Format(time.RFC3339),
					Entities: []EntityDiff{
						{
							Generated: generatedEntity{Name: "entity1", Kind: "kind1"},
							Action:    "create",
							Diff:      "diff1",
						},
					},
				})
				return dm
			}(),
			hash: "hash1",
			expected: []EntityDiff{
				{
					Generated: generatedEntity{Name: "entity1", Kind: "kind1"},
					Action:    "create",
					Diff:      "diff1",
				},
			},
			shouldFail: false,
		},
		{
			name: "retrieve diff by hash from full diffMap",
			initial: func() diffMap {
				dm := diffMap{
					diffs:     map[string]ConfigDiff{},
					times:     map[string]string{},
					length:    2,
					hashQueue: list.New(),
				}
				dm.Update(ConfigDiff{
					Hash:      "hash1",
					Timestamp: now.Format(time.RFC3339),
					Entities: []EntityDiff{
						{
							Generated: generatedEntity{Name: "entity1", Kind: "kind1"},
							Action:    "create",
							Diff:      "diff1",
						},
					},
				})
				dm.Update(ConfigDiff{
					Hash:      "hash2",
					Timestamp: now.Format(time.RFC3339),
					Entities: []EntityDiff{
						{
							Generated: generatedEntity{Name: "entity2", Kind: "kind2"},
							Action:    "update",
							Diff:      "diff2",
						},
					},
				})
				return dm
			}(),
			hash: "hash2",
			expected: []EntityDiff{
				{
					Generated: generatedEntity{Name: "entity2", Kind: "kind2"},
					Action:    "update",
					Diff:      "diff2",
				},
			},
			shouldFail: false,
		},
		{
			name: "retrieve non-existent diff by hash from full diffMap",
			initial: func() diffMap {
				dm := diffMap{
					diffs:     map[string]ConfigDiff{},
					times:     map[string]string{},
					length:    2,
					hashQueue: list.New(),
				}
				dm.Update(ConfigDiff{
					Hash:      "hash1",
					Timestamp: now.Format(time.RFC3339),
					Entities: []EntityDiff{
						{
							Generated: generatedEntity{Name: "entity1", Kind: "kind1"},
							Action:    "create",
							Diff:      "diff1",
						},
					},
				})
				dm.Update(ConfigDiff{
					Hash:      "hash2",
					Timestamp: now.Format(time.RFC3339),
					Entities: []EntityDiff{
						{
							Generated: generatedEntity{Name: "entity2", Kind: "kind2"},
							Action:    "update",
							Diff:      "diff2",
						},
					},
				})
				return dm
			}(),
			hash:       "hash3",
			expected:   []EntityDiff{},
			shouldFail: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.initial.ByHash(tt.hash)
			if tt.shouldFail {
				assert.Error(t, err)
				return
			}
			assert.NoError(t, err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
