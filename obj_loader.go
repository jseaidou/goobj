package goobj

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
)

type OBJFile struct {
	Vertices []Vertex
	Normals  []Vertex
}

type Vertex struct {
	x float64
	y float64
	z float64
	w float64
}

func LoadOBJ(path string) (*OBJFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	vertices := []Vertex{}
	normals := []Vertex{}
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

		// Vertex
		if line[0] == 'v' && len(line) > 1 && unicode.IsSpace(int32(line[1])) {
			coords := strings.Fields(line[1:])
			fCoords, err := parseFloats(coords)
			if err != nil {
				return nil, err
			}
			vertex, err := NewVertex(fCoords)
			if err != nil {
				return nil, err
			}
			vertices = append(vertices, vertex)
		}

		// Normals
		if line[0] == 'v' && len(line) > 2 && line[1] == 'n' && unicode.IsSpace(int32(line[2])) {
			coords := strings.Fields(line[2:])
			fCoords, err := parseFloats(coords)
			if err != nil {
				return nil, err
			}

			vertex, err := NewVertex(fCoords)
			if err != nil {
				return nil, err
			}
			normals = append(normals, vertex)
		}
	}
	return &OBJFile{
		Vertices: vertices,
		Normals:  normals,
	}, nil
}

/**
 * Creates a new vertex struct from a given array containing the x,y,z coordinates.
 * They define the position of the vertex in 3 dimensions
 * The w coordinate is the weight required for rational curves and surfaces. It is
 * not required for non-rational curves and surfaces. If you do notspecify a value
 * for w, the default is 1.0.
 *
 * The function returns a error when coords is less than 3 elements big.
 */
func NewVertex(coords []float64) (Vertex, error) {
	if len(coords) < 3 {
		return Vertex{}, fmt.Errorf("Error creating Vertex. Expected \"v x y z [w]\" for coordinates. Got: %v", coords)
	}

	w := 1.0
	if len(coords) > 3 {
		w = coords[3]
	}
	vertex := Vertex{
		x: coords[0],
		y: coords[1],
		z: coords[2],
		w: w,
	}
	return vertex, nil
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
