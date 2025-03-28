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
	Instances    []ZonInstance
	Areas        []ZonArea
	Lights       []ZonLight
	V4Info       V4Info
	V4Dat        V4Dat
	name         *eqgName
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

// ZonInstance is an Instance
type ZonInstance struct {
	ModelTag    string
	InstanceTag string
	Translation [3]float32
	Rotation    [3]float32
	Scale       float32
	Lits        []uint32
}

// ZonArea is an area
type ZonArea struct {
	Name        string
	Center      [3]float32
	Orientation [3]float32
	Extents     [3]float32
}

// ZonLight is a light
type ZonLight struct {
	Name     string
	Position [3]float32
	Color    [3]float32
	Radius   float32
}

func (zon *Zon) Identity() string {
	return "zon"
}

func (zon *Zon) String() string {
	out := ""
	out += fmt.Sprintf("Version: %d\n", zon.Version)
	out += fmt.Sprintf("Models: %d\n", len(zon.Models))
	out += fmt.Sprintf("Instances: %d\n", len(zon.Instances))
	out += fmt.Sprintf("Areas: %d\n", len(zon.Areas))
	out += fmt.Sprintf("Lights: %d\n", len(zon.Lights))
	return out
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
	instanceCount := dec.Uint32()
	areaCount := dec.Uint32()
	lightCount := dec.Uint32()

	nameData := dec.Bytes(int(nameLength))
	zon.name.parse(nameData)

	for i := 0; i < int(modelCount); i++ {
		name := zon.name.byOffset(dec.Int32())
		zon.Models = append(zon.Models, name)
	}
	//dec.SetDebugMode(true)

	for i := 0; i < int(instanceCount); i++ {
		instance := ZonInstance{}

		modelID := dec.Int32()
		if modelID >= 0 && modelID < int32(len(zon.Models)) {
			instance.ModelTag = zon.Models[modelID]
		}

		instanceID := dec.Int32()
		// for v1 zones, this is lit name, for v2+ it's instance name
		if instanceID < int32(len(zon.name.nameBuf)) {
			instance.InstanceTag = zon.name.byOffset(instanceID)
		}

		instance.Translation[0] = dec.Float32()
		instance.Translation[1] = dec.Float32()
		instance.Translation[2] = dec.Float32()

		instance.Rotation[0] = dec.Float32()
		instance.Rotation[1] = dec.Float32()
		instance.Rotation[2] = dec.Float32()

		instance.Scale = dec.Float32()
		if zon.Version > 1 {
			litCount := dec.Uint32()
			for j := 0; j < int(litCount); j++ {
				instance.Lits = append(instance.Lits, dec.Uint32())
			}
		}

		zon.Instances = append(zon.Instances, instance)

	}
	//os.WriteFile("/src/quail/test/src.bin", dec.DebugBuf(), 0644)

	for i := 0; i < int(areaCount); i++ {
		area := ZonArea{}

		area.Name = zon.name.byOffset(dec.Int32())

		area.Center[0] = dec.Float32()
		area.Center[1] = dec.Float32()
		area.Center[2] = dec.Float32()

		area.Orientation[0] = dec.Float32()
		area.Orientation[1] = dec.Float32()
		area.Orientation[2] = dec.Float32()

		area.Extents[0] = dec.Float32()
		area.Extents[1] = dec.Float32()
		area.Extents[2] = dec.Float32()

		zon.Areas = append(zon.Areas, area)
	}

	for i := 0; i < int(lightCount); i++ {
		light := ZonLight{}

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

	pos := dec.Pos()
	endPos, err := r.Seek(0, io.SeekEnd)
	if err != nil {
		return fmt.Errorf("seek end: %w", err)
	}
	if pos != endPos {
		if pos < endPos {
			return fmt.Errorf("%d bytes remaining (%d total)", endPos-pos, endPos)
		}

		return fmt.Errorf("read past end of file")
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
