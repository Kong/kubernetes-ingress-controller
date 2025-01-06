package atc

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewPredicate(t *testing.T) {
	testCases := []struct {
		name          string
		lhs           LHS
		op            BinaryOperator
		rhs           Literal
		expectedError error
		expression    string
	}{
		{
			name:       "predicate for string field (http.path)",
			lhs:        FieldHTTPPath,
			op:         OpEqual,
			rhs:        StringLiteral("/path"),
			expression: `http.path == "/path"`,
		},
		{
			name:       "predicate for integer field (net.dst.port)",
			lhs:        FieldNetDstPort,
			op:         OpGreaterEqual,
			rhs:        IntLiteral(1024),
			expression: `net.dst.port >= 1024`,
		},
		{
			name:          "unmatched types (LHS string RHS integer)",
			lhs:           FieldHTTPPath,
			op:            OpEqual,
			rhs:           IntLiteral(80),
			expectedError: ErrTypeNotMatch,
		},
		{
			name:          "unmatched types (LHS integer RHS string)",
			lhs:           FieldNetDstPort,
			op:            OpEqual,
			rhs:           StringLiteral("/"),
			expectedError: ErrTypeNotMatch,
		},
		{
			name:          "invalid operator (contains for integer)",
			lhs:           FieldNetDstPort,
			op:            OpContains,
			rhs:           IntLiteral(10),
			expectedError: ErrOperatorInvalid,
		},
		{
			name:          "invalid operator (less than for string)",
			lhs:           FieldHTTPPath,
			op:            OpLessThan,
			rhs:           StringLiteral("/v1"),
			expectedError: ErrOperatorInvalid,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			predicate, err := NewPredicate(tc.lhs, tc.op, tc.rhs)
			if tc.expectedError == nil {
				require.NoError(t, err)
				require.Equal(t, tc.expression, predicate.Expression())
			} else {
				require.EqualError(t, err, tc.expectedError.Error())
			}
		})
	}
}
