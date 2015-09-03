package goobj

import (
	"bufio"
	"io"
	"os"
	"strconv"
	"strings"
)

type OBJFile struct {
	Shapes []Shape
}

type Shape struct {
	Vertices   *VertexData
	Attributes *CSAttributes
}

func makeEmptyCSAtt() *CSAttributes {
	return &CSAttributes{
		BMatu: Matrix{
			Elements: [][]float64{},
		},
		BMatv: Matrix{
			Elements: [][]float64{},
		},
	}
}

func makeEmptyVData() *VertexData {
	return &VertexData{
		Normals:  &[]Vertex{},
		Vertices: &[]Vertex{},
		Points:   &[]Vertex{},
		Textures: &[]Vertex{},
	}
}

func LoadOBJ(path string) (*OBJFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	vdata := makeEmptyVData()
	atts := makeEmptyCSAtt()
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return nil, err
		}
		if err == io.EOF {
			break
		}

		// trim any whitespace
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if line[0] == '#' {
			// Skip comment
			continue
		}

		// Geometric Vertex
		if checkToken(line, Geometric, true) {
			vertex, err := LoadVertex(line[len(Geometric):], 3, Geometric)
			if err != nil {
				return nil, err
			}
			appendVertex(vdata.Vertices, vertex)
		}

		// Vertex Normals
		if checkToken(line, Normal, true) {
			vertex, err := LoadVertex(line[len(Normal):], 3, Normal)
			if err != nil {
				return nil, err
			}
			appendVertex(vdata.Normals, vertex)
		}

		// Parameter Space Vertex
		if checkToken(line, ParameterSpace, true) {
			vertex, err := LoadVertex(line[len(ParameterSpace):], 2, ParameterSpace)
			if err != nil {
				return nil, err
			}
			appendVertex(vdata.Points, vertex)
		}

		// Texture Vertex
		if checkToken(line, Texture, true) {
			vertex, err := LoadVertex(line[len(Texture):], 2, Texture)
			if err != nil {
				return nil, err
			}
			appendVertex(vdata.Textures, vertex)
		}

		// Curve/Surface Type
		if checkToken(line, CSTypeStatement, true) {
			ctype, err := LoadCSStatement(CSTypeStatement, line[len(CSTypeStatement):])
			if err != nil {
				return nil, err
			}
			atts.Type = ctype.(CSType)
		}
	}
	return &OBJFile{
		Shapes: []Shape{
			Shape{
				Attributes: atts,
				Vertices:   vdata,
			},
		},
	}, nil
}

func appendVertex(slice *[]Vertex, vertex Vertex) {
	*slice = append(*slice, vertex)
}

/**
 * Given a token and a compare string, checkToken will check whether:
 *		1. token[len(compare)] == compare
 *		2. if spaceAfter == true then checkToken will check whether
 *			token[len(compare) + 1] == compare + " "
 */
func checkToken(token, compare string, spaceAfter bool) bool {
	if len(token) >= len(compare) {
		indexToCheck := len(compare)
		if spaceAfter {
			indexToCheck += 1
			compare = compare + " "
		}
		t := token[:indexToCheck]
		return t == compare
	}
	return false
}

/**
 * A wrapper function around strconv.ParseFloat. The purpose is to
 * take in a array of float numbers in string format and to convert
 * them into floating point representation.
 * Returns a error (if any) during the string -> floating point conversion
 */
func parseFloats(floats []string) ([]float64, error) {
	ret := []float64{}
	for _, float := range floats {
		p, err := strconv.ParseFloat(float, 64)
		if err != nil {
			return ret, err
		}
		ret = append(ret, p)
	}
	return ret, nil

}
