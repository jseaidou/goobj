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

type Shape struct {
	Name string
	Mesh Mesh
}

type Mesh struct {
	Positions   []float64
	Normals     []float64
	Texcoords   []float64
	Indices     []int
	MaterialIds []int
}

type Vertex struct {
	x float64
	y float64
	z float64
	w float64
}

func LoadOBJ(path string) ([]Vertex, error) {
	file, err := os.Open(path)
	if err != nil {
		return []Vertex{}, err
	}
	reader := bufio.NewReader(file)
	vertices := []Vertex{}
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return vertices, err
		}
		if err == io.EOF {
			return vertices, nil
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
				return vertices, nil
			}
			vertex, err := NewVertex(fCoords)
			if err != nil {
				return vertices, err
			}
			vertices = append(vertices, vertex)
		}
	}
	return vertices, nil
}

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
