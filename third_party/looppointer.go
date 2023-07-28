//go:build third_party

package third_party

import _ "github.com/kyoh86/looppointer/cmd/looppointer"

//go:generate go install -modfile go.mod github.com/kyoh86/looppointer/cmd/looppointer
