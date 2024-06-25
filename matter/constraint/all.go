package constraint

import (
	"encoding/json"

	"github.com/project-chip/alchemy/matter/types"
)

type AllConstraint struct {
	Value string `json:"value"`
}

func NewAllConstraint(value string) *AllConstraint {
	return &AllConstraint{Value: value}
}

func (c *AllConstraint) Type() Type {
	return ConstraintTypeAll
}

func (c *AllConstraint) ASCIIDocString(dataType *types.DataType) string {
	return c.Value
}

func (c *AllConstraint) Equal(o Constraint) bool {
	_, ok := o.(*AllConstraint)
	return ok
}

func (c *AllConstraint) Min(cc Context) (min types.DataTypeExtreme) {
	return
}

func (c *AllConstraint) Max(cc Context) (max types.DataTypeExtreme) {
	return
}

func (c *AllConstraint) Default(cc Context) (max types.DataTypeExtreme) {
	return
}

func (c *AllConstraint) Clone() Constraint {
	return &AllConstraint{}
}

func (c *AllConstraint) MarshalJSON() ([]byte, error) {
	js := map[string]any{
		"type":  "all",
		"value": c.Value,
	}
	return json.Marshal(js)
}
