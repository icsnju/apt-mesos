package registry

type Constraint struct {
	Key   string         `json:"key"`
	Value interface{}    `json:"value"`
	Type  ConstraintType `enum=ConstraintType,json:"type,omitempty"`
}

type ConstraintType int32

const (
	ConstraintTypeOptions ConstraintType = 0
	ConstraintTypeRequire ConstraintType = 0
)
