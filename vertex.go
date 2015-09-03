package goobj

import (
	"fmt"
	"strings"
)

type VertexData struct {
	Vertices *[]Vertex
	Normals  *[]Vertex
	Points   *[]Vertex
	Textures *[]Vertex
}

type Vertex struct {
	Coords []float64
	VType  string
}

const (
	// Vertex Types
	Geometric      = "v"
	Normal         = "vn"
	ParameterSpace = "vp"
	Texture        = "vt"
)

/**
 * Creates a new vertex struct from a given array.
 * Vertices can be:
 *	1.  v with coords x y z [w]. These represent geometric vertex and its x y z coordinates. Rational
 *		curves and surfaces require a fourth homogeneous coordinate, also
 *	    called the weight (w). w is the weight required for rational curves and surfaces. It is
 *		not require for non-rational curves and surfaces. If you do not
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
	if vType == Geometric && len(coords) < 4 {
		coords = append(coords, 1.0)
	} else if vType == ParameterSpace && len(coords) < 3 {
		coords = append(coords, 1.0)
	} else if vType == Texture && len(coords) < 3 {
		coords = append(coords, 0.0)
	}
	return Vertex{
		Coords: coords,
		VType:  vType,
	}
}

/**
 * loadVertex breaks apart a string and uses the fields as vertex coordinates in the form of
 *     Vertex: x y z [w]
 *     Normal: i j k
 *     Point : u v [w]
 *     Texture: u v [w]
 * and appends it to the slice.
 */
func LoadVertex(line string, minArgs int, vType string) (Vertex, error) {
	coords := strings.Fields(line)
	if len(coords) < minArgs {
		fmt.Errorf("Expected %v for vertex type: %s, but got: %v", minArgs, vType, len(coords))
	}
	fCoords, err := parseFloats(coords)
	if err != nil {
		return Vertex{}, err
	}
	vertex := NewVertex(fCoords, vType)
	return vertex, nil
}
