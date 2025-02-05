package metadata

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestProjectNameFromRepo(t *testing.T) {
	testCases := []struct {
		Name     string
		Repo     string
		Expected string
	}{
		{
			Name:     "valid Git SSH repo",
			Repo:     "git@github.com:Kong/kubernetes-ingress-controller.git",
			Expected: "kubernetes-ingress-controller",
		},
		{
			Name:     "valid Git HTTPS repo",
			Repo:     "https://github.com/Kong/kubernetes-ingress-controller.git",
			Expected: "kubernetes-ingress-controller",
		},
		{
			Name:     "not set",
			Repo:     NotSet,
			Expected: "NOT_SET",
		},
		{
			Name:     "empty repo",
			Repo:     "",
			Expected: "NOT_SET",
		},
		{
			Name:     "empty repo (whitespace)",
			Repo:     "    ",
			Expected: "NOT_SET",
		},
		{
			Name:     "repo set to arbitral string",
			Repo:     "this-is-my-project",
			Expected: "this-is-my-project",
		},
	}

	for _, tC := range testCases {
		t.Run(tC.Name, func(t *testing.T) {
			require.Equal(t, tC.Expected, projectNameFromRepo(tC.Repo))
		})
	}
}
