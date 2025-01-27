package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Zon is a zone EQGZ header
type Zon struct {
	MetaFileName string
	Version      uint32
	Models       []string
	Objects      []ZonInstance
	Regions      []ZonArea
	Lights       []Light
	V4Info       V4Info
	V4Dat        V4Dat
	name         *eqgName
}

func (zon *Zon) Identity() string {
	return "zon"
}

type V4Info struct {
	Name                 string
	MinLng               int
	MinLat               int
	MaxLng               int
	MaxLat               int
	MinExtents           [3]float32
	MaxExtents           [3]float32
	UnitsPerVert         float32
	QuadsPerTile         int
	CoverMapInputSize    int
	LayeringMapInputSize int
}

type V4Dat struct {
	Unk1            uint32
	Unk2            uint32
	Unk3            uint32
	BaseTileTexture string
	Tiles           []V4DatTile
}

type V4DatTile struct {
	Lng     int32
	Lat     int32
	Unk     uint32
	Colors  []uint32
	Colors2 []uint32
}

// ZonInstance is an object
type ZonInstance struct {
	MeshName     string
	InstanceName string
	Translation  [3]float32
	Rotation     [3]float32
	Scale        float32
	Lits         []uint32
}

// ZonArea is a region
type ZonArea struct {
	Name     string
	Position [3]float32
	Color    [3]float32
	Radius   float32
	Unk1     uint32
	Unk2     uint32
}

// Light is a light
type Light struct {
	Name     string
	Position [3]float32
	Color    [3]float32
	Radius   float32
}

// Decode reads a ZON file
func (zon *Zon) Read(r io.ReadSeeker) error {
	if zon.name == nil {
		zon.name = &eqgName{}
	}
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	// Read header
	header := dec.StringFixed(4)
	if header != "EQGZ" && header != "EQTZ" {
		return fmt.Errorf("header %s does not match EQGZ or EQTZ", header)
	}

	if header == "EQTZ" {
		return zon.ReadV4(r)
	}

	zon.Version = dec.Uint32()

	nameLength := dec.Uint32()
	meshCount := dec.Uint32()
	instanceCount := dec.Uint32()
	areaCount := dec.Uint32()
	lightCount := dec.Uint32()

	nameData := dec.Bytes(int(nameLength))
	zon.name.parse(nameData)

	//os.WriteFile("src.txt", []byte(fmt.Sprintf("%+v", names)), 0644)

	for i := 0; i < int(meshCount); i++ {
		name := zon.name.byOffset(dec.Int32())
		zon.Models = append(zon.Models, name)
	}

	for i := 0; i < int(instanceCount); i++ {
		instance := ZonInstance{}
		meshIndex := dec.Int32()

		if meshIndex >= int32(len(zon.Models)) {
			return fmt.Errorf("%d object nameIndex %d out of range (%d)", i, meshIndex, len(zon.Models))
		}
		if meshIndex == -1 {
			continue
		}

		instance.MeshName = zon.Models[meshIndex]

		instance.InstanceName = zon.name.byOffset(dec.Int32())

		instance.Translation[0] = dec.Float32()
		instance.Translation[1] = dec.Float32()
		instance.Translation[2] = dec.Float32()

		instance.Rotation[0] = dec.Float32()
		instance.Rotation[1] = dec.Float32()
		instance.Rotation[2] = dec.Float32()

		instance.Scale = dec.Float32()
		if zon.Version >= 2 {
			litCount := dec.Uint32()
			for j := 0; j < int(litCount); j++ {
				instance.Lits = append(instance.Lits, dec.Uint32())
			}
		}
		zon.Objects = append(zon.Objects, instance)
	}

	for i := 0; i < int(areaCount); i++ {
		area := ZonArea{}

		area.Name = zon.name.byOffset(dec.Int32())

		area.Position[0] = dec.Float32()
		area.Position[1] = dec.Float32()
		area.Position[2] = dec.Float32()

		area.Color[0] = dec.Float32()
		area.Color[1] = dec.Float32()
		area.Color[2] = dec.Float32()

		area.Radius = dec.Float32()
		//		area.Radius[1] = dec.Float32()
		//		area.Radius[2] = dec.Float32()

		//region.Unk1 = dec.Uint32()
		//region.Unk2 = dec.Uint32()

		zon.Regions = append(zon.Regions, area)
	}

	for i := 0; i < int(lightCount); i++ {
		light := Light{}

		light.Name = zon.name.byOffset(dec.Int32())

		light.Position[0] = dec.Float32()
		light.Position[1] = dec.Float32()
		light.Position[2] = dec.Float32()

		light.Color[0] = dec.Float32()
		light.Color[1] = dec.Float32()
		light.Color[2] = dec.Float32()

		light.Radius = dec.Float32()

		zon.Lights = append(zon.Lights, light)
	}

	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}

	return nil
}

// SetFileName sets the name of the file
func (zon *Zon) SetFileName(name string) {
	zon.MetaFileName = name
}

// FileName returns the name of the file
func (zon *Zon) FileName() string {
	return zon.MetaFileName
}
