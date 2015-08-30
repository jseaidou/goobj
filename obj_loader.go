package goobj

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
)

type OBJFile struct {
	Vertices []Vertex
	Normals  []Vertex
	Points   []Vertex
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
	points := []Vertex{}
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
		if checkToken(line, "v", true) {
			err = loadVertex(line[1:], &vertices)
			if err != nil {
				return nil, err
			}
		}

		// Normals
		if checkToken(line, "vn", true) {
			err = loadVertex(line[2:], &normals)
			if err != nil {
				return nil, err
			}
		}

		// Points
		if checkToken(line, "vp", true) {
			err = loadVertex(line[2:], &points)
			if err != nil {
				return nil, err
			}
		}
	}
	return &OBJFile{
		Vertices: vertices,
		Normals:  normals,
		Points:   points,
	}, nil
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
 * loadVertex breaks apart a string and uses them as vertex coordinates in the form of
 *     Vertex: x y z [w]
 *     Normal: i j k
 *     Point : u v w
 * and appends it to the slice.
 */
func loadVertex(line string, slice *[]Vertex) error {
	if slice == nil {
		fmt.Errorf("Got unintialized slice in loadVertex")
	}
	coords := strings.Fields(line)
	fCoords, err := parseFloats(coords)
	if err != nil {
		return err
	}
	vertex, err := NewVertex(fCoords)
	if err != nil {
		return err
	}
	*slice = append(*slice, vertex)
	return nil
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
