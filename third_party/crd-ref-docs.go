//go:build third_party

package third_party

import (
	_ "github.com/elastic/crd-ref-docs"
)

//go:generate go install -modfile go.mod github.com/elastic/crd-ref-docs
