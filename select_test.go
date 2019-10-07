package genji

import (
	"testing"

	"github.com/asdine/genji/query/expr"
	"github.com/asdine/genji/query/q"
	"github.com/stretchr/testify/require"
)

func TestParserSelect(t *testing.T) {
	tests := []struct {
		name     string
		s        string
		expected Statement
		mustFail bool
	}{
		{"NoCond", "SELECT * FROM test",
			selectStmt{
				tableSelector: q.Table("test"),
			}, false},
		{"WithFields", "SELECT a, b FROM test",
			selectStmt{
				FieldSelectors: []FieldSelector{q.Field("a"), q.Field("b")},
				tableSelector:  q.Table("test"),
			}, false},
		{"WithCond", "SELECT * FROM test WHERE age = 10",
			selectStmt{
				tableSelector: q.Table("test"),
				whereExpr:     expr.Eq(q.Field("age"), expr.Int64Value(10)),
			}, false},
		{"WithLimit", "SELECT * FROM test WHERE age = 10 LIMIT 20",
			selectStmt{
				tableSelector: q.Table("test"),
				whereExpr:     expr.Eq(q.Field("age"), expr.Int64Value(10)),
				limitExpr:     expr.Int64Value(20),
			}, false},
		{"WithOffset", "SELECT * FROM test WHERE age = 10 OFFSET 20",
			selectStmt{
				tableSelector: q.Table("test"),
				whereExpr:     expr.Eq(q.Field("age"), expr.Int64Value(10)),
				offsetExpr:    expr.Int64Value(20),
			}, false},
		{"WithLimitThenOffset", "SELECT * FROM test WHERE age = 10 LIMIT 10 OFFSET 20",
			selectStmt{
				tableSelector: q.Table("test"),
				whereExpr:     expr.Eq(q.Field("age"), expr.Int64Value(10)),
				offsetExpr:    expr.Int64Value(20),
				limitExpr:     expr.Int64Value(10),
			}, false},
		{"WithOffsetThenLimit", "SELECT * FROM test WHERE age = 10 OFFSET 20 LIMIT 10", nil, true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			q, err := ParseQuery(test.s)
			if !test.mustFail {
				require.NoError(t, err)
				require.Len(t, q.Statements, 1)
				require.EqualValues(t, test.expected, q.Statements[0])
			} else {
				require.Error(t, err)
			}
		})
	}
}
