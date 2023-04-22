package geo

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/xackery/quail/helper"
)

// MeshManager is a mesh manager
type MeshManager struct {
	meshes []*Mesh
}

// NewMeshManager creates a new mesh manager
func NewMeshManager() *MeshManager {
	return &MeshManager{}
}

// Mesh returns a mesh by name
func (e *MeshManager) Mesh(name string) (*Mesh, bool) {
	for _, mesh := range e.meshes {
		if mesh.Name == name {
			return mesh, true
		}
	}
	return nil, false
}

// BlenderExport writes all materials to a file
func (e *MeshManager) BlenderExport(dir string) error {
	if len(e.meshes) > 0 {
		mw, err := os.Create(fmt.Sprintf("%s/mesh.txt", dir))
		if err != nil {
			return fmt.Errorf("create file %s/mesh.txt: %w", dir, err)
		}
		defer mw.Close()
		mesh := &Mesh{}
		err = mesh.WriteHeader(mw)
		if err != nil {
			return fmt.Errorf("write header: %w", err)
		}
		for _, mesh = range e.meshes {
			err = mesh.Write(mw)
			if err != nil {
				return fmt.Errorf("write mesh: %w", err)
			}
		}
	}

	for _, mesh := range e.meshes {
		meshName := strings.TrimSuffix(mesh.Name, "_DMSPRITEDEF")
		if len(mesh.Triangles) > 0 {
			trianglePath := fmt.Sprintf("%s/%s_triangle.txt", dir, meshName)
			mw, err := os.Create(trianglePath)
			if err != nil {
				return fmt.Errorf("create file %s: %w", trianglePath, err)
			}
			defer mw.Close()
			triangle := &Triangle{}
			err = triangle.WriteHeader(mw)
			if err != nil {
				return fmt.Errorf("write header: %w", err)
			}

			for _, t := range mesh.Triangles {
				err = t.Write(mw)
				if err != nil {
					return fmt.Errorf("write triangle: %w", err)
				}
			}
		}

		if len(mesh.Vertices) > 0 {
			vertexPath := fmt.Sprintf("%s/%s_vertex.txt", dir, meshName)
			pw, err := os.Create(vertexPath)
			if err != nil {
				return fmt.Errorf("create file %s: %w", vertexPath, err)
			}
			defer pw.Close()
			vertex := &Vertex{}
			err = vertex.WriteHeader(pw)
			if err != nil {
				return fmt.Errorf("write vertex header: %w", err)
			}
			for _, v := range mesh.Vertices {
				err = v.Write(pw)
				if err != nil {
					return fmt.Errorf("write vertex: %w", err)
				}
			}
		}

		if len(mesh.Bones) > 0 {
			bonePath := fmt.Sprintf("%s/%s_bone.txt", dir, meshName)
			bw, err := os.Create(bonePath)
			if err != nil {
				return fmt.Errorf("create file %s: %w", bonePath, err)
			}
			defer bw.Close()
			bone := &Bone{}
			err = bone.WriteHeader(bw)
			if err != nil {
				return fmt.Errorf("write bone header: %w", err)
			}
			for _, b := range mesh.Bones {
				err = b.Write(bw)
				if err != nil {
					return fmt.Errorf("write bone: %w", err)
				}
			}
		}

		if len(mesh.Animations) > 0 {
			animationPath := fmt.Sprintf("%s/%s_animation.txt", dir, meshName)
			bw, err := os.Create(animationPath)
			if err != nil {
				return fmt.Errorf("create file %s: %w", animationPath, err)
			}
			defer bw.Close()
			anim := BoneAnimation{}
			err = anim.WriteHeader(bw)
			if err != nil {
				return fmt.Errorf("write bone animation header: %w", err)
			}
			for _, anim = range mesh.Animations {
				err = anim.Write(bw)
				if err != nil {
					return fmt.Errorf("write bone animation: %w", err)
				}
			}
		}
	}

	return nil
}

// ReadFile reads a material file
func (e *MeshManager) ReadFile(dir string) error {
	var err error
	meshPath := fmt.Sprintf("%s/mesh.txt", dir)
	err = e.meshRead(meshPath)
	if err != nil {
		return fmt.Errorf("meshRead: %w", err)
	}

	for _, mesh := range e.meshes {
		trianglePath := fmt.Sprintf("%s/%s_triangle.txt", dir, mesh.Name)
		err = e.triangleRead(trianglePath, mesh)
		if err != nil {
			return fmt.Errorf("read triangle: %w", err)
		}

		vertexPath := fmt.Sprintf("%s/%s_vertex.txt", dir, mesh.Name)
		err = e.vertexRead(vertexPath, mesh)
		if err != nil {
			return fmt.Errorf("read vertex: %w", err)
		}

		bonePath := fmt.Sprintf("%s/%s_bone.txt", dir, mesh.Name)
		err = e.boneRead(bonePath, mesh)
		if err != nil {
			return fmt.Errorf("read bone: %w", err)
		}

		animationPath := fmt.Sprintf("%s/%s_animation.txt", dir, mesh.Name)
		err = e.animationRead(animationPath, mesh)
		if err != nil {
			return fmt.Errorf("read animation: %w", err)
		}
	}
	return nil
}

// VertexAdd adds a vertex to the mesh manager
func (e *MeshManager) VertexAdd(meshName string, vertex Vertex) error {
	mesh, ok := e.Mesh(meshName)
	if !ok {
		mesh = &Mesh{Name: meshName}
		e.Add(mesh)
	}
	mesh.Vertices = append(mesh.Vertices, vertex)
	return nil
}

// TriangleAdd adds a triangle to the mesh manager
func (e *MeshManager) TriangleAdd(meshName string, triangle Triangle) error {
	mesh, ok := e.Mesh(meshName)
	if !ok {
		mesh = &Mesh{Name: meshName}
		e.Add(mesh)
	}

	triangle.MaterialName = strings.ToLower(triangle.MaterialName)
	mesh.Triangles = append(mesh.Triangles, triangle)
	return nil
}

// BoneAdd adds a bone to the mesh manager
func (e *MeshManager) BoneAdd(meshName string, bone Bone) error {
	mesh, ok := e.Mesh(meshName)
	if !ok {
		mesh = &Mesh{Name: meshName}
		e.Add(mesh)
	}

	mesh.Bones = append(mesh.Bones, bone)
	return nil
}

// TriangleCount returns the number of triangles
func (e *MeshManager) TriangleCount(meshName string) int {
	mesh, ok := e.Mesh(meshName)
	if !ok {
		return 0
	}

	return len(mesh.Triangles)
}

// VertexCount returns the number of vertices
func (e *MeshManager) VertexCount(meshName string) int {
	mesh, ok := e.Mesh(meshName)
	if !ok {
		return 0
	}

	return len(mesh.Vertices)
}

// BoneCount returns the number of bones
func (e *MeshManager) BoneCount(meshName string) int {
	mesh, ok := e.Mesh(meshName)
	if !ok {
		return 0
	}
	return len(mesh.Bones)
}

// Inspect prints out the mesh manager
func (e *MeshManager) Inspect() {
	for _, mesh := range e.meshes {
		fmt.Println(mesh.Name, len(mesh.Bones), "bones:")
		for i, bone := range mesh.Bones {
			fmt.Printf("  %d %s\n", i, bone.Name)
		}

		fmt.Println(len(mesh.Triangles), "triangles")
		fmt.Println(len(mesh.Vertices), "vertices")
	}
}

// Add adds a mesh to the mesh manager
func (e *MeshManager) Add(mesh *Mesh) {
	e.meshes = append(e.meshes, mesh)
}

func (e *MeshManager) triangleRead(path string, mesh *Mesh) error {
	fi, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("stat %s: %w", path, err)
	}
	if fi.IsDir() {
		return fmt.Errorf("expected %s to be a file", path)
	}

	r, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if lineNumber == 1 {
			continue
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 3 {
			return fmt.Errorf("invalid triangle.txt (expected 4 records) line %d: %s", lineNumber, line)
		}

		mesh.Triangles = append(mesh.Triangles, Triangle{
			Index:        AtoUIndex3(parts[0]),
			Flag:         helper.AtoU32(parts[1]),
			MaterialName: parts[2],
		})
	}
	return nil
}

func (e *MeshManager) vertexRead(path string, mesh *Mesh) error {
	r, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("open %s: %w", path, err)
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if lineNumber == 1 {
			continue
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			return fmt.Errorf("invalid vertex.txt (expected 5 records) line %d: %s", lineNumber, line)
		}
		vert := Vertex{
			Position: AtoVector3(parts[0]),
			Normal:   AtoVector3(parts[1]),
			Uv:       AtoVector2(parts[2]),
			Uv2:      AtoVector2(parts[3]),
			Tint:     AtoRGBA(parts[4]),
		}
		vert.Position = Vector3{X: -vert.Position.Y, Y: vert.Position.X, Z: vert.Position.Z}
		vert.Normal = Vector3{X: -vert.Normal.Y, Y: vert.Normal.X, Z: vert.Normal.Z}

		mesh.Vertices = append(mesh.Vertices, vert)
	}
	return nil
}

func (e *MeshManager) boneRead(path string, mesh *Mesh) error {
	r, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("open %s: %w", path, err)
	}
	defer r.Close()

	scanner := bufio.NewScanner(r)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if lineNumber == 1 {
			continue
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 7 {
			return fmt.Errorf("invalid bone.txt (expected 7 records) line %d: %s", lineNumber, line)
		}
		mesh.Bones = append(mesh.Bones, Bone{
			Name:          parts[0],
			ChildIndex:    helper.AtoI32(parts[1]),
			ChildrenCount: helper.AtoU32(parts[2]),
			Next:          helper.AtoI32(parts[3]),
			Pivot:         AtoVector3(parts[4]),
			Rotation:      AtoQuad4(parts[5]),
			Scale:         AtoVector3(parts[6]),
		})
	}
	return nil
}

func (e *MeshManager) animationRead(path string, mesh *Mesh) error {
	r, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("open %s: %w", path, err)
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if lineNumber == 1 {
			continue
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		parts := strings.Split(line, "|")
		if len(parts) < 7 {
			return fmt.Errorf("invalid animation.txt (expected 7 records) line %d: %s", lineNumber, line)
		}
		return fmt.Errorf("TODO: implement animation.txt")
	}
	return nil
}

func (e *MeshManager) meshRead(path string) error {
	r, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}

		return fmt.Errorf("open %s: %w", path, err)
	}
	defer r.Close()
	scanner := bufio.NewScanner(r)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if lineNumber == 1 {
			continue
		}
		line := scanner.Text()
		if line == "" {
			continue
		}
		e.meshes = append(e.meshes, &Mesh{Name: line})
	}
	return nil
}

// Meshes returns a pointer to the meshes slice.
func (e *MeshManager) Meshes() []*Mesh {
	return e.meshes
}

// TriangleTotalCount returns the total triangle count
func (e *MeshManager) TriangleTotalCount() int {
	total := 0

	for _, mesh := range e.meshes {
		total += len(mesh.Triangles)
	}
	return total
}

// VertexTotalCount returns the total vertex count
func (e *MeshManager) VertexTotalCount() int {
	total := 0

	for _, mesh := range e.meshes {
		total += len(mesh.Vertices)
	}
	return total
}

// BoneTotalCount returns the total bone count
func (e *MeshManager) BoneTotalCount() int {
	total := 0

	for _, mesh := range e.meshes {
		total += len(mesh.Bones)
	}
	return total
}
