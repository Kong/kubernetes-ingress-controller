package util

import (
	"reflect"
	"testing"
)

func Test_GetCombinations(t *testing.T) {
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.args.relations.GetCombinations(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetCombinations() = %v, want %v", got, tt.want)
			}
		})
	}
}
