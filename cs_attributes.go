package goobj

import (
	"fmt"
	"strings"
)

type Rationality string
type CSTypeName string
type Degree int

const (
	Rational    Rationality = "rat"
	NonRational             = "non-rat"
)

const (
	BasisMatrix CSTypeName = "bmatrix"
	Bezier                 = "bezier"
	BSpline                = "bspline"
	Cardinal               = "cardinal"
	Taylor                 = "taylor"
)

const (
	CSTypeStatement    = "cstype"
	CSDegStatement     = "deg"
	CSBMatrixStatement = "bmat"
	CSStepStatement    = "step"
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

func LoadCSStatement(statement string, line string) (interface{}, error) {
	switch statement {
	case CSTypeStatement:
		return loadCSTypeStatement(line)
	case CSDegStatement:
		return loadCSDegStatement(line)
	case CSBMatrixStatement:
		return loadCSBMatStatement(line)
	case CSStepStatement:
		return loadCSStepStatement(line)
	}
	return nil, fmt.Errorf("Unrecognized Curve/Surface statement: %v", statement)
}

func loadCSTypeStatement(line string) (CSType, error) {
	cstype := strings.Fields(line)
	ctype := CSType{}
	if len(cstype) == 1 {
		ctype.Rat = NonRational
		ctype.Name = CSTypeName(cstype[0])
	} else if len(cstype) == 2 {
		ctype.Rat = Rationality(cstype[0])
		ctype.Name = CSTypeName(cstype[1])
	} else {
		return ctype, fmt.Errorf("Incorrect number of arguments for statement: cstype. Expected 'cstype [rat] type' got: %v", cstype)
	}

	return ctype, nil
}

func loadCSDegStatement(line string) (PDegree, error) {
	return PDegree{}, nil
}

func loadCSBMatStatement(line string) (Matrix, error) {
	return Matrix{}, nil
}

func loadCSStepStatement(line string) (Step, error) {
	return Step{}, nil
}
