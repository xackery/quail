package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/tag"
)

// Zon is a zone
type Zon struct {
	MetaFileName string   `yaml:"file_name"`
	Version      uint32   `yaml:"version"`
	Models       []string `yaml:"models"`
	Objects      []Object `yaml:"objects"`
	Regions      []Region `yaml:"regions"`
	Lights       []Light  `yaml:"lights"`
	V4Info       V4Info   `yaml:"v4info"`
	V4Dat        V4Dat    `yaml:"v4dat"`
}

func (zon *Zon) Identity() string {
	return "zon"
}

type V4Info struct {
	Name                 string  `yaml:"name"`
	MinLng               int     `yaml:"min_lng"`
	MinLat               int     `yaml:"min_lat"`
	MaxLng               int     `yaml:"max_lng"`
	MaxLat               int     `yaml:"max_lat"`
	MinExtents           Vector3 `yaml:"min_extents"`
	MaxExtents           Vector3 `yaml:"max_extents"`
	UnitsPerVert         float32 `yaml:"units_per_vert"`
	QuadsPerTile         int     `yaml:"quads_per_tile"`
	CoverMapInputSize    int     `yaml:"cover_map_input_size"`
	LayeringMapInputSize int     `yaml:"layering_map_input_size"`
}

type V4Dat struct {
	Unk1            uint32 `yaml:"unk1"`
	Unk2            uint32 `yaml:"unk2"`
	Unk3            uint32 `yaml:"unk3"`
	BaseTileTexture string `yaml:"base_tile_texture"`
	Tiles           []V4DatTile
}

type V4DatTile struct {
	Lng     int32    `yaml:"lng"`
	Lat     int32    `yaml:"lat"`
	Unk     uint32   `yaml:"unk"`
	Colors  []uint32 `yaml:"colors"`
	Colors2 []uint32 `yaml:"colors2"`
}

// Object is an object
type Object struct {
	ModelName    string  `yaml:"model_name"`
	InstanceName string  `yaml:"instance_name"`
	Position     Vector3 `yaml:"position"`
	Rotation     Vector3 `yaml:"rotation"`
	Scale        float32 `yaml:"scale"`
	Lits         []*RGBA `yaml:"-"` // used in v2+ zones, omitted since it's huge
}

// Region is a region
type Region struct {
	Name    string  `yaml:"name"`
	Center  Vector3 `yaml:"center"`
	Unknown Vector3 `yaml:"unknown"`
	Extent  Vector3 `yaml:"extent"`
	Unk1    uint32  `yaml:"unk1"`
	Unk2    uint32  `yaml:"unk2"`
}

// Light is a light
type Light struct {
	Name     string  `yaml:"name"`
	Position Vector3 `yaml:"position"`
	Color    Vector3 `yaml:"color"`
	Radius   float32 `yaml:"radius"`
}

// Decode reads a ZON file
func (zon *Zon) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)

	tag.New()

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

	tag.Add(0, dec.Pos(), "red", "header")
	nameData := dec.Bytes(int(nameLength))

	names := make(map[int32]string)

	chunk := []byte{}
	lastOffset := 0
	for i, b := range nameData {
		if b == 0 {
			names[int32(lastOffset)] = string(chunk)
			chunk = []byte{}
			lastOffset = i + 1
			continue
		}
		chunk = append(chunk, b)
	}
	tag.Add(tag.LastPos(), dec.Pos(), "green", fmt.Sprintf("names (%d total)", len(names)))

	NamesSet(names)
	//os.WriteFile("src.txt", []byte(fmt.Sprintf("%+v", names)), 0644)

	for i := 0; i < int(modelCount); i++ {
		name := Name(dec.Int32())
		zon.Models = append(zon.Models, name)
	}
	tag.AddRand(tag.LastPos(), dec.Pos(), fmt.Sprintf("modelNames (%d total)", modelCount))

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

		object.InstanceName = Name(dec.Int32())

		object.Position.Y = dec.Float32() // y before x
		object.Position.X = dec.Float32()
		object.Position.Z = dec.Float32()

		object.Rotation.X = dec.Float32()
		object.Rotation.Y = dec.Float32()
		object.Rotation.Z = dec.Float32()

		object.Scale = dec.Float32()
		tag.AddRand(tag.LastPos(), dec.Pos(), fmt.Sprintf("%d|%s", i, object.InstanceName))
		if zon.Version >= 2 {
			litCount := dec.Uint32()
			for j := 0; j < int(litCount); j++ {
				lit := RGBA{}
				lit.R = dec.Uint8()
				lit.G = dec.Uint8()
				lit.B = dec.Uint8()
				lit.A = dec.Uint8()
				object.Lits = append(object.Lits, &lit)
			}
			tag.AddRand(tag.LastPos(), dec.Pos(), fmt.Sprintf("%d|%s|lit_data", i, object.InstanceName))
		}
		zon.Objects = append(zon.Objects, object)
	}

	for i := 0; i < int(regionCount); i++ {
		region := Region{}

		region.Name = Name(dec.Int32())

		region.Center.X = dec.Float32()
		region.Center.Y = dec.Float32()
		region.Center.Z = dec.Float32()

		region.Unknown.X = dec.Float32()
		region.Unknown.Y = dec.Float32()
		region.Unknown.Z = dec.Float32()

		region.Extent.X = dec.Float32()
		region.Extent.Y = dec.Float32()
		region.Extent.Z = dec.Float32()

		//region.Unk1 = dec.Uint32()
		//region.Unk2 = dec.Uint32()

		zon.Regions = append(zon.Regions, region)
		tag.AddRand(tag.LastPos(), dec.Pos(), fmt.Sprintf("%d|%s", i, region.Name))
	}

	for i := 0; i < int(lightCount); i++ {
		light := Light{}

		light.Name = Name(dec.Int32())

		light.Position.X = dec.Float32()
		light.Position.Y = dec.Float32()
		light.Position.Z = dec.Float32()

		light.Color.X = dec.Float32()
		light.Color.Y = dec.Float32()
		light.Color.Z = dec.Float32()

		light.Radius = dec.Float32()

		zon.Lights = append(zon.Lights, light)
		tag.AddRand(tag.LastPos(), dec.Pos(), fmt.Sprintf("%d|%s", i, light.Name))
	}

	if dec.Error() != nil {
		return fmt.Errorf("read: %w", dec.Error())
	}

	//log.Debugf("%s (zon) readd %d model refs, %d object refs, %d regions, %d lights", zon.Header.Name, len(zon.Models), len(zon.Objects), len(zon.Regions), len(zon.Lights))
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
