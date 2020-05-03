package expr

import (
	"strconv"
	"strings"

	"github.com/asdine/genji/document"
)

// A FieldSelector is a ResultField that extracts a field from a document at a given path.
type FieldSelector []string

// Name joins the chunks of the fields selector with the . separator.
func (f FieldSelector) Name() string {
	return strings.Join(f, ".")
}

// Eval extracts the document from the context and selects the right field.
// It implements the Expr interface.
func (f FieldSelector) Eval(stack EvalStack) (document.Value, error) {
	if stack.Document == nil {
		return nullLitteral, document.ErrFieldNotFound
	}

	var v document.Value
	var a document.Array
	var err error

	for i, chunk := range f {
		if stack.Document != nil {
			v, err = stack.Document.GetByField(chunk)
		} else {
			idx, err := strconv.Atoi(chunk)
			if err != nil {
				return nullLitteral, nil
			}
			v, err = a.GetByIndex(idx)
		}
		if err == document.ErrFieldNotFound {
			return nullLitteral, nil
		}
		if err != nil {
			return nullLitteral, err
		}

		if i+1 == len(f) {
			break
		}

		stack.Document = nil
		a = nil

		switch v.Type {
		case document.DocumentValue:
			stack.Document, err = v.ConvertToDocument()
		case document.ArrayValue:
			a, err = v.ConvertToArray()
		default:
			return nullLitteral, nil
		}
		if err != nil {
			return nullLitteral, err
		}
	}

	return v, nil
}
