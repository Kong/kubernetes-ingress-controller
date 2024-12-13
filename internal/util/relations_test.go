package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetCombinations(t *testing.T) {
	type args struct {
		relations ForeignRelations
	}
	tests := []struct {
		name string
		args args
		want []Rel
	}{
		{
			name: "empty",
			args: args{
				relations: ForeignRelations{},
			},
			want: nil,
		},
		{
			name: "plugins on consumer only",
			args: args{
				relations: ForeignRelations{
					Consumer: []FR{{Identifier: "foo"}, {Identifier: "bar"}},
				},
			},
			want: []Rel{
				{
					Consumer: FR{Identifier: "foo"},
				},
				{
					Consumer: FR{Identifier: "bar"},
				},
			},
		},
		{
			name: "plugins on consumer group only",
			args: args{
				relations: ForeignRelations{
					ConsumerGroup: []FR{{Identifier: "foo"}, {Identifier: "bar"}},
				},
			},
			want: []Rel{
				{
					ConsumerGroup: FR{Identifier: "foo"},
				},
				{
					ConsumerGroup: FR{Identifier: "bar"},
				},
			},
		},
		{
			name: "plugins on service only",
			args: args{
				relations: ForeignRelations{
					Service: []FR{{Identifier: "foo"}, {Identifier: "bar"}},
				},
			},
			want: []Rel{
				{
					Service: FR{Identifier: "foo"},
				},
				{
					Service: FR{Identifier: "bar"},
				},
			},
		},
		{
			name: "plugins on routes only",
			args: args{
				relations: ForeignRelations{
					Route: []FR{{Identifier: "foo"}, {Identifier: "bar"}},
				},
			},
			want: []Rel{
				{
					Route: FR{Identifier: "foo"},
				},
				{
					Route: FR{Identifier: "bar"},
				},
			},
		},
		{
			name: "plugins on service and routes only",
			args: args{
				relations: ForeignRelations{
					Route:   []FR{{Identifier: "foo"}, {Identifier: "bar"}},
					Service: []FR{{Identifier: "foo"}, {Identifier: "bar"}},
				},
			},
			want: []Rel{
				{
					Service: FR{Identifier: "foo"},
				},
				{
					Service: FR{Identifier: "bar"},
				},
				{
					Route: FR{Identifier: "foo"},
				},
				{
					Route: FR{Identifier: "bar"},
				},
			},
		},
		{
			name: "plugins on combination of route and consumer",
			args: args{
				relations: ForeignRelations{
					Route:    []FR{{Identifier: "foo"}, {Identifier: "bar"}},
					Consumer: []FR{{Identifier: "foo"}, {Identifier: "bar"}},
				},
			},
			want: []Rel{
				{
					Consumer: FR{Identifier: "foo"},
					Route:    FR{Identifier: "foo"},
				},
				{
					Consumer: FR{Identifier: "foo"},
					Route:    FR{Identifier: "bar"},
				},
				{
					Consumer: FR{Identifier: "bar"},
					Route:    FR{Identifier: "foo"},
				},
				{
					Consumer: FR{Identifier: "bar"},
					Route:    FR{Identifier: "bar"},
				},
			},
		},
		{
			name: "plugins on combination of service and consumer",
			args: args{
				relations: ForeignRelations{
					Service:  []FR{{Identifier: "foo"}, {Identifier: "bar"}},
					Consumer: []FR{{Identifier: "foo"}, {Identifier: "bar"}},
				},
			},
			want: []Rel{
				{
					Consumer: FR{Identifier: "foo"},
					Service:  FR{Identifier: "foo"},
				},
				{
					Consumer: FR{Identifier: "foo"},
					Service:  FR{Identifier: "bar"},
				},
				{
					Consumer: FR{Identifier: "bar"},
					Service:  FR{Identifier: "foo"},
				},
				{
					Consumer: FR{Identifier: "bar"},
					Service:  FR{Identifier: "bar"},
				},
			},
		},
		{
			name: "plugins on combination of service,route and consumer",
			args: args{
				relations: ForeignRelations{
					Consumer: []FR{{Identifier: "c1"}, {Identifier: "c2"}},
					Route:    []FR{{Identifier: "r1"}, {Identifier: "r2"}},
					Service:  []FR{{Identifier: "s1"}, {Identifier: "s2"}},
				},
			},
			want: []Rel{
				{
					Consumer: FR{Identifier: "c1"},
					Service:  FR{Identifier: "s1"},
				},
				{
					Consumer: FR{Identifier: "c1"},
					Service:  FR{Identifier: "s2"},
				},
				{
					Consumer: FR{Identifier: "c1"},
					Route:    FR{Identifier: "r1"},
				},
				{
					Consumer: FR{Identifier: "c1"},
					Route:    FR{Identifier: "r2"},
				},
				{
					Consumer: FR{Identifier: "c2"},
					Service:  FR{Identifier: "s1"},
				},
				{
					Consumer: FR{Identifier: "c2"},
					Service:  FR{Identifier: "s2"},
				},
				{
					Consumer: FR{Identifier: "c2"},
					Route:    FR{Identifier: "r1"},
				},
				{
					Consumer: FR{Identifier: "c2"},
					Route:    FR{Identifier: "r2"},
				},
			},
		},
		{
			name: "plugins on combination of service,route and consumer group",
			args: args{
				relations: ForeignRelations{
					Route:         []FR{{Identifier: "r1"}, {Identifier: "r2"}},
					Service:       []FR{{Identifier: "s1"}, {Identifier: "s2"}},
					ConsumerGroup: []FR{{Identifier: "cg1"}, {Identifier: "cg2"}},
				},
			},
			want: []Rel{
				{
					ConsumerGroup: FR{Identifier: "cg1"},
					Service:       FR{Identifier: "s1"},
				},
				{
					ConsumerGroup: FR{Identifier: "cg1"},
					Service:       FR{Identifier: "s2"},
				},
				{
					ConsumerGroup: FR{Identifier: "cg1"},
					Route:         FR{Identifier: "r1"},
				},
				{
					ConsumerGroup: FR{Identifier: "cg1"},
					Route:         FR{Identifier: "r2"},
				},
				{
					ConsumerGroup: FR{Identifier: "cg2"},
					Service:       FR{Identifier: "s1"},
				},
				{
					ConsumerGroup: FR{Identifier: "cg2"},
					Service:       FR{Identifier: "s2"},
				},
				{
					ConsumerGroup: FR{Identifier: "cg2"},
					Route:         FR{Identifier: "r1"},
				},
				{
					ConsumerGroup: FR{Identifier: "cg2"},
					Route:         FR{Identifier: "r2"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, tt.args.relations.GetCombinations())
		})
	}
}

func BenchmarkGetCombinations(b *testing.B) {
	b.Run("consumer groups", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			relations := ForeignRelations{
				Route:         []FR{{Identifier: "r1"}, {Identifier: "r2"}},
				Service:       []FR{{Identifier: "s1"}, {Identifier: "s2"}},
				ConsumerGroup: []FR{{Identifier: "cg1"}, {Identifier: "cg2"}},
			}

			rels := relations.GetCombinations()
			_ = rels
		}
	})
	b.Run("consumers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			relations := ForeignRelations{
				Route:    []FR{{Identifier: "r1"}, {Identifier: "r2"}},
				Service:  []FR{{Identifier: "s1"}, {Identifier: "s2"}},
				Consumer: []FR{{Identifier: "c1"}, {Identifier: "c2"}, {Identifier: "c3"}},
			}

			rels := relations.GetCombinations()
			_ = rels
		}
	})
}
