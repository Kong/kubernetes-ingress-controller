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
					Consumer: []string{"foo", "bar"},
				},
			},
			want: []Rel{
				{
					Consumer: "foo",
				},
				{
					Consumer: "bar",
				},
			},
		},
		{
			name: "plugins on consumer group only",
			args: args{
				relations: ForeignRelations{
					ConsumerGroup: []string{"foo", "bar"},
				},
			},
			want: []Rel{
				{
					ConsumerGroup: "foo",
				},
				{
					ConsumerGroup: "bar",
				},
			},
		},
		{
			name: "plugins on service only",
			args: args{
				relations: ForeignRelations{
					Service: []string{"foo", "bar"},
				},
			},
			want: []Rel{
				{
					Service: "foo",
				},
				{
					Service: "bar",
				},
			},
		},
		{
			name: "plugins on routes only",
			args: args{
				relations: ForeignRelations{
					Route: []string{"foo", "bar"},
				},
			},
			want: []Rel{
				{
					Route: "foo",
				},
				{
					Route: "bar",
				},
			},
		},
		{
			name: "plugins on service and routes only",
			args: args{
				relations: ForeignRelations{
					Route:   []string{"foo", "bar"},
					Service: []string{"foo", "bar"},
				},
			},
			want: []Rel{
				{
					Service: "foo",
				},
				{
					Service: "bar",
				},
				{
					Route: "foo",
				},
				{
					Route: "bar",
				},
			},
		},
		{
			name: "plugins on combination of route and consumer",
			args: args{
				relations: ForeignRelations{
					Route:    []string{"foo", "bar"},
					Consumer: []string{"foo", "bar"},
				},
			},
			want: []Rel{
				{
					Consumer: "foo",
					Route:    "foo",
				},
				{
					Consumer: "bar",
					Route:    "foo",
				},
				{
					Consumer: "foo",
					Route:    "bar",
				},
				{
					Consumer: "bar",
					Route:    "bar",
				},
			},
		},
		{
			name: "plugins on combination of service and consumer",
			args: args{
				relations: ForeignRelations{
					Service:  []string{"foo", "bar"},
					Consumer: []string{"foo", "bar"},
				},
			},
			want: []Rel{
				{
					Consumer: "foo",
					Service:  "foo",
				},
				{
					Consumer: "bar",
					Service:  "foo",
				},
				{
					Consumer: "foo",
					Service:  "bar",
				},
				{
					Consumer: "bar",
					Service:  "bar",
				},
			},
		},
		{
			name: "plugins on combination of service,route and consumer",
			args: args{
				relations: ForeignRelations{
					Consumer: []string{"c1", "c2"},
					Route:    []string{"r1", "r2"},
					Service:  []string{"s1", "s2"},
				},
			},
			want: []Rel{
				{
					Consumer: "c1",
					Service:  "s1",
				},
				{
					Consumer: "c2",
					Service:  "s1",
				},
				{
					Consumer: "c1",
					Service:  "s2",
				},
				{
					Consumer: "c2",
					Service:  "s2",
				},
				{
					Consumer: "c1",
					Route:    "r1",
				},
				{
					Consumer: "c2",
					Route:    "r1",
				},
				{
					Consumer: "c1",
					Route:    "r2",
				},
				{
					Consumer: "c2",
					Route:    "r2",
				},
			},
		},
		{
			name: "plugins on combination of service,route and consumer group",
			args: args{
				relations: ForeignRelations{
					Route:         []string{"r1", "r2"},
					Service:       []string{"s1", "s2"},
					ConsumerGroup: []string{"cg1", "cg2"},
				},
			},
			want: []Rel{
				{
					ConsumerGroup: "cg1",
					Service:       "s1",
				},
				{
					ConsumerGroup: "cg2",
					Service:       "s1",
				},
				{
					ConsumerGroup: "cg1",
					Service:       "s2",
				},
				{
					ConsumerGroup: "cg2",
					Service:       "s2",
				},
				{
					ConsumerGroup: "cg1",
					Route:         "r1",
				},
				{
					ConsumerGroup: "cg2",
					Route:         "r1",
				},
				{
					ConsumerGroup: "cg1",
					Route:         "r2",
				},
				{
					ConsumerGroup: "cg2",
					Route:         "r2",
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
			relationsCG := ForeignRelations{
				Route:         []string{"r1", "r2"},
				Service:       []string{"s1", "s2"},
				ConsumerGroup: []string{"cg1", "cg2"},
			}

			rels := relationsCG.GetCombinations()
			_ = rels
		}
	})
	b.Run("consumers", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			relationsCG := ForeignRelations{
				Route:    []string{"r1", "r2"},
				Service:  []string{"s1", "s2"},
				Consumer: []string{"c1", "c2", "c3"},
			}

			rels := relationsCG.GetCombinations()
			_ = rels
		}
	})
}
