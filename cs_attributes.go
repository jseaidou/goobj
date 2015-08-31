package goobj

type Rationality string
type CSTypeName string
type Degree int

const (
	Rational    Rationality = "rat"
	NonRational Rationality = "non-rat"
)

const (
	BasisMatrix CSTypeName = "bmatrix"
	Bezier      CSTypeName = "bezier"
	BSpline     CSTypeName = "bspline"
	Cardinal    CSTypeName = "cardinal"
	Taylor      CSTypeName = "taylor"
)

type CSType struct {
	Rat  Rationality // Optional, defaults to non-rational.
	Name CSTypeName
}

type PDegree struct {
	Degu Degree
	Degv Degree
}

type Matrix struct {
	Elements [][]float64
}

type Step struct {
	StepU int
	StepV int
}

type CSAttributes struct {
	Type   CSType
	Degree PDegree
	BMatu  Matrix
	BMatv  Matrix
	Step   Step
}
