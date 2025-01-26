package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
)

// Zon is a zone
type Zon struct {
	MetaFileName string
	Version      uint32
	Models       []string
	Objects      []Object
	Regions      []Region
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

// Object is an object
type Object struct {
	ModelName    string
	InstanceName string
	Position     [3]float32
	Rotation     [3]float32
	Scale        float32
	Lits         [][4]uint8
}

// Region is a region
type Region struct {
	Name    string
	Center  [3]float32
	Unknown [3]float32
	Extent  [3]float32
	Unk1    uint32
	Unk2    uint32
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
	modelCount := dec.Uint32()
	objectCount := dec.Uint32()
	regionCount := dec.Uint32()
	lightCount := dec.Uint32()

	nameData := dec.Bytes(int(nameLength))
	zon.name.parse(nameData)

	//os.WriteFile("src.txt", []byte(fmt.Sprintf("%+v", names)), 0644)

	for i := 0; i < int(modelCount); i++ {
		name := zon.name.byOffset(dec.Int32())
		zon.Models = append(zon.Models, name)
	}

	for i := 0; i < int(objectCount); i++ {
		object := Object{}
		nameIndex := dec.Int32()

		if nameIndex >= int32(len(zon.Models)) {
			return fmt.Errorf("%d object nameIndex %d out of range (%d)", i, nameIndex, len(zon.Models))
		}
		if nameIndex < 0 {
			return fmt.Errorf("%d object nameIndex %d less than 0", i, nameIndex)
		}

		object.ModelName = zon.Models[nameIndex]

		object.InstanceName = zon.name.byOffset(dec.Int32())

		object.Position[1] = dec.Float32() // y before x
		object.Position[0] = dec.Float32()
		object.Position[2] = dec.Float32()

		object.Rotation[0] = dec.Float32()
		object.Rotation[1] = dec.Float32()
		object.Rotation[2] = dec.Float32()

		object.Scale = dec.Float32()
		if zon.Version >= 2 {
			litCount := dec.Uint32()
			for j := 0; j < int(litCount); j++ {
				object.Lits = append(object.Lits, [4]uint8{dec.Uint8(), dec.Uint8(), dec.Uint8(), dec.Uint8()})
			}
		}
		zon.Objects = append(zon.Objects, object)
	}

	for i := 0; i < int(regionCount); i++ {
		region := Region{}

		region.Name = zon.name.byOffset(dec.Int32())

		region.Center[0] = dec.Float32()
		region.Center[1] = dec.Float32()
		region.Center[2] = dec.Float32()

		region.Unknown[0] = dec.Float32()
		region.Unknown[1] = dec.Float32()
		region.Unknown[2] = dec.Float32()

		region.Extent[0] = dec.Float32()
		region.Extent[1] = dec.Float32()
		region.Extent[2] = dec.Float32()

		//region.Unk1 = dec.Uint32()
		//region.Unk2 = dec.Uint32()

		zon.Regions = append(zon.Regions, region)
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
