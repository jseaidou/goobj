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
	Textures []Vertex
}

type Vertex struct {
	coords []float64
	vType  string
}

func LoadOBJ(path string) (*OBJFile, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	reader := bufio.NewReader(file)
	obj := OBJFile{}
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
			err = loadVertex(line[1:], &obj.Vertices, 3, "v")
			if err != nil {
				return nil, err
			}
		}

		// Normals
		if checkToken(line, "vn", true) {
			err = loadVertex(line[2:], &obj.Normals, 3, "vn")
			if err != nil {
				return nil, err
			}
		}

		// Points
		if checkToken(line, "vp", true) {
			err = loadVertex(line[2:], &obj.Points, 2, "vp")
			if err != nil {
				return nil, err
			}
		}

		// Texture
		if checkToken(line, "vt", true) {
			err = loadVertex(line[2:], &obj.Textures, 2, "vt")
			if err != nil {
				return nil, err
			}
		}
	}
	return &obj, nil
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
 * loadVertex breaks apart a string and uses the fields as vertex coordinates in the form of
 *     Vertex: x y z [w]
 *     Normal: i j k
 *     Point : u v [w]
 *     Texture: u v [w]
 * and appends it to the slice.
 */
func loadVertex(line string, slice *[]Vertex, minArgs int, vType string) error {
	if slice == nil {
		fmt.Errorf("Got unintialized slice in loadVertex")
	}
	coords := strings.Fields(line)
	if len(coords) < minArgs {
		fmt.Errorf("Expected %v for vertex type: %s, but got: %v", minArgs, vType, len(coords))
	}
	fCoords, err := parseFloats(coords)
	if err != nil {
		return err
	}
	vertex := NewVertex(fCoords, vType)
	*slice = append(*slice, vertex)
	return nil
}

/**
 * Creates a new vertex struct from a given array.
 * Vertices can be:
 *	1.  v with coords x y z [w]. These represent geometric vertex and its x y z coordinates. Rational
 *		curves and surfaces require a fourth homogeneous coordinate, also
 *	    called the weight (w). w is the weight required for rational curves and surfaces. It is
 *		not required for non-rational curves and surfaces. If you do not
 *		specify a value for w, the default is 1.0.
 *
 *	2. vp with coords u v [w]. These represent a vertex point in the parameter space of a curve or surface.
 *	   The usage determines how many coordinates are required. Special points for curves require a 1D control
 *	   point (u only) in the parameter space of the curve. Special points for surfaces require a 2D point
 *	   (u and v) in the parameter space of the surface. Control points for non-rational trimming curves require
 *	   u and v coordinates. Control points for rational trimming curves require u, v, and w (weight) coordinates.
 *     u is the point in the parameter space of a curve or the first coordinate in the parameter space of a surface.
 *     v is the second coordinate in the parameter space of a surface. w is the weight required for rational
 *     trimming curves. If you do not specify a value for w, it defaults to 1.0.
 *
 *	3. vn with coords i j v. These represent a vertex normal. Vertex normals affect the smooth-shading and rendering
 *	   of geometry. For polygons, vertex normals are used in place of the actual facet normals.  For surfaces, vertex
 *     normals are interpolated over the entire surface and replace the actual analytic surface normal.
 *     When vertex normals are present, they supersede smoothing groups.
 *
 *	4. vt with coords u v [w]. These represent a texture vertex and its coordinates. A 1D texture requires only u
 *     texture coordinates, a 2D texture requires both u and v texture coordinates, and a 3D texture requires all
 *     three coordinates.
 *
 *     u is the value for the horizontal direction of the texture.
 *	   v is an optional argument. v is the value for the vertical direction of the texture. The default is 0.
 *	   w is an optional argument. w is a value for the depth of the texture. The default is 0.
 */
func NewVertex(coords []float64, vType string) Vertex {
	if vType == "v" && len(coords) < 4 {
		coords = append(coords, 1.0)
	} else if vType == "vp" && len(coords) < 3 {
		coords = append(coords, 1.0)
	} else if vType == "vt" && len(coords) < 3 {
		coords = append(coords, 0.0)
	}
	vertex := Vertex{
		coords: coords,
		vType:  vType,
	}
	return vertex
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
