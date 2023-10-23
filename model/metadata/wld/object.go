package wld

import (
	"encoding/binary"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
)

// ParticleSprite is ParticleSpriteDef in libeq, empty in openzone, PARTICLESPRITEDEF in wld
type ParticleSprite struct {
	NameRef                     int32
	Flags                       uint32
	VerticesCount               uint32
	Unknown                     uint32
	CenterOffset                common.Vector3
	Radius                      float32
	Vertices                    []common.Vector3
	RenderMethod                uint32
	RenderFlags                 uint32
	RenderPen                   uint32
	RenderBrightness            float32
	RenderScaledAmbient         float32
	RenderSimpleSpriteReference uint32
	RenderUVInfoOrigin          common.Vector3
	RenderUVInfoUAxis           common.Vector3
	RenderUVInfoVAxis           common.Vector3
	RenderUVMapEntryCount       uint32
	RenderUVMapEntries          []common.Vector2
	Pen                         []uint32
}

func (e *ParticleSprite) FragCode() int {
	return 0x0C
}

func (e *ParticleSprite) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.VerticesCount)
	enc.Uint32(e.Unknown)
	if e.Flags&0x01 != 0 { // has center offset
		enc.Float32(e.CenterOffset.X)
		enc.Float32(e.CenterOffset.Y)
		enc.Float32(e.CenterOffset.Z)
	}
	if e.Flags&0x02 != 0 { // has radius
		enc.Float32(e.Radius)
	}
	if e.VerticesCount > 0 { // has vertices
		for _, vertex := range e.Vertices {
			enc.Float32(vertex.X)
			enc.Float32(vertex.Y)
			enc.Float32(vertex.Z)
		}
	}
	enc.Uint32(e.RenderMethod)
	enc.Uint32(e.RenderFlags)
	enc.Uint32(e.RenderPen)
	enc.Float32(e.RenderBrightness)
	enc.Float32(e.RenderScaledAmbient)
	enc.Uint32(e.RenderSimpleSpriteReference)
	enc.Float32(e.RenderUVInfoOrigin.X)
	enc.Float32(e.RenderUVInfoOrigin.Y)
	enc.Float32(e.RenderUVInfoOrigin.Z)
	enc.Float32(e.RenderUVInfoUAxis.X)
	enc.Float32(e.RenderUVInfoUAxis.Y)
	enc.Float32(e.RenderUVInfoUAxis.Z)
	enc.Float32(e.RenderUVInfoVAxis.X)
	enc.Float32(e.RenderUVInfoVAxis.Y)
	enc.Float32(e.RenderUVInfoVAxis.Z)
	enc.Uint32(e.RenderUVMapEntryCount)
	for _, entry := range e.RenderUVMapEntries {
		enc.Float32(entry.X)
		enc.Float32(entry.Y)
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeParticleSprite(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &ParticleSprite{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.VerticesCount = dec.Uint32()
	d.Unknown = dec.Uint32()
	if d.Flags&0x01 != 0 { // has center offset
		d.CenterOffset.X = dec.Float32()
		d.CenterOffset.Y = dec.Float32()
		d.CenterOffset.Z = dec.Float32()
	}
	if d.Flags&0x02 != 0 { // has radius
		d.Radius = dec.Float32()
	}
	if d.VerticesCount > 0 { // has vertices
		for i := uint32(0); i < d.VerticesCount; i++ {
			var vertex common.Vector3
			vertex.X = dec.Float32()
			vertex.Y = dec.Float32()
			vertex.Z = dec.Float32()
			d.Vertices = append(d.Vertices, vertex)
		}
	}
	d.RenderMethod = dec.Uint32()
	d.RenderFlags = dec.Uint32()
	d.RenderPen = dec.Uint32()
	d.RenderBrightness = dec.Float32()
	d.RenderScaledAmbient = dec.Float32()
	d.RenderSimpleSpriteReference = dec.Uint32()
	d.RenderUVInfoOrigin.X = dec.Float32()
	d.RenderUVInfoOrigin.Y = dec.Float32()
	d.RenderUVInfoOrigin.Z = dec.Float32()
	d.RenderUVInfoUAxis.X = dec.Float32()
	d.RenderUVInfoUAxis.Y = dec.Float32()
	d.RenderUVInfoUAxis.Z = dec.Float32()
	d.RenderUVInfoVAxis.X = dec.Float32()
	d.RenderUVInfoVAxis.Y = dec.Float32()
	d.RenderUVInfoVAxis.Z = dec.Float32()
	d.RenderUVMapEntryCount = dec.Uint32()
	for i := uint32(0); i < d.RenderUVMapEntryCount; i++ {
		var entry common.Vector2
		entry.X = dec.Float32()
		entry.Y = dec.Float32()
		d.RenderUVMapEntries = append(d.RenderUVMapEntries, entry)
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// ParticleSpriteRef is ParticleSprite in libeq, empty in openzone, PARTICLESPRITE (ref) in wld
type ParticleSpriteRef struct {
	NameRef              int32
	ParticleSpriteDefRef int32
	Flags                uint32
}

func (e *ParticleSpriteRef) FragCode() int {
	return 0x0D
}

func (e *ParticleSpriteRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.ParticleSpriteDefRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeParticleSpriteRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &ParticleSpriteRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.ParticleSpriteDefRef = dec.Int32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// CompositeSprite is empty in libeq, empty in openzone, COMPOSITESPRITEDEF in wld, Actor in lantern
type CompositeSprite struct {
	NameRef int32
	Flags   uint32
}

func (e *CompositeSprite) FragCode() int {
	return 0x0E
}

func (e *CompositeSprite) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeCompositeSprite(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &CompositeSprite{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// CompositeSpriteRef is empty in libeq, empty in openzone, COMPOSITESPRITE (ref) in wld
type CompositeSpriteRef struct {
	NameRef               int32
	CompositeSpriteDefRef int32
	Flags                 uint32
}

func (e *CompositeSpriteRef) FragCode() int {
	return 0x0F
}

func (e *CompositeSpriteRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.CompositeSpriteDefRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeCompositeSpriteRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &CompositeSpriteRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.CompositeSpriteDefRef = dec.Int32()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// Model is ActorDef in libeq, Static in openzone, ACTORDEF in wld
type Model struct {
	NameRef          int32
	Flags            uint32
	CallbackNameRef  int32
	ActionCount      uint32
	FragmentRefCount uint32
	BoundsRef        int32
	CurrentAction    uint32
	Offset           common.Vector3
	Rotation         common.Vector3
	Unk1             uint32
	Actions          []Action
	FragmentRefs     []uint32
	Unk2             uint32
}

type Action struct {
	LodCount uint32
	Unk1     uint32
	Lods     []float32
}

func (e *Model) FragCode() int {
	return 0x14
}

func (e *Model) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)

	enc.Int32(e.CallbackNameRef)
	enc.Uint32(e.ActionCount)
	enc.Uint32(e.FragmentRefCount)
	enc.Int32(e.BoundsRef)
	if e.Flags&0x1 == 0x1 {
		enc.Uint32(e.CurrentAction)
	}
	if e.Flags&0x2 == 0x2 {
		enc.Float32(e.Offset.X)
		enc.Float32(e.Offset.Y)
		enc.Float32(e.Offset.Z)
		enc.Float32(e.Rotation.X)
		enc.Float32(e.Rotation.Y)
		enc.Float32(e.Rotation.Z)
		enc.Uint32(e.Unk1)
	}
	for _, action := range e.Actions {
		enc.Uint32(action.LodCount)
		enc.Uint32(action.Unk1)
		for _, lod := range action.Lods {
			enc.Float32(lod)
		}
	}
	for _, fragmentRef := range e.FragmentRefs {
		enc.Uint32(fragmentRef)
	}
	enc.Uint32(e.Unk2)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeModel(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &Model{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.CallbackNameRef = dec.Int32()
	d.ActionCount = dec.Uint32()
	d.FragmentRefCount = dec.Uint32()
	d.BoundsRef = dec.Int32()
	if d.Flags&0x1 == 0x1 {
		d.CurrentAction = dec.Uint32()
	}
	if d.Flags&0x2 == 0x2 {
		d.Offset.X = dec.Float32()
		d.Offset.Y = dec.Float32()
		d.Offset.Z = dec.Float32()
		d.Rotation.X = dec.Float32()
		d.Rotation.Y = dec.Float32()
		d.Rotation.Z = dec.Float32()
		d.Unk1 = dec.Uint32()
	}
	for i := uint32(0); i < d.ActionCount; i++ {
		var action Action
		action.LodCount = dec.Uint32()
		action.Unk1 = dec.Uint32()
		for j := uint32(0); j < action.LodCount; j++ {
			action.Lods = append(action.Lods, dec.Float32())
		}
		d.Actions = append(d.Actions, action)
	}
	for i := uint32(0); i < d.FragmentRefCount; i++ {
		d.FragmentRefs = append(d.FragmentRefs, dec.Uint32())
	}
	d.Unk2 = dec.Uint32()

	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// ModelRef is Actor in libeq, Object Location in openzone, ACTORINST in wld, ObjectInstance in lantern
type ModelRef struct {
	NameRef        int32
	ActorDefRef    int32
	Flags          uint32
	SphereRef      uint32
	CurrentAction  uint32
	Offset         common.Vector3
	Rotation       common.Vector3
	Unk1           uint32
	BoundingRadius float32
	Scale          float32
	SoundNameRef   int32
	Unk2           int32
}

func (e *ModelRef) FragCode() int {
	return 0x15
}

func (e *ModelRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.ActorDefRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.SphereRef)
	if e.Flags&0x1 == 0x1 {
		enc.Uint32(e.CurrentAction)
	}
	if e.Flags&0x2 == 0x2 {
		enc.Float32(e.Offset.X)
		enc.Float32(e.Offset.Y)
		enc.Float32(e.Offset.Z)
		enc.Float32(e.Rotation.X)
		enc.Float32(e.Rotation.Y)
		enc.Float32(e.Rotation.Z)
		enc.Uint32(e.Unk1)
	}
	if e.Flags&0x4 == 0x4 {
		enc.Float32(e.BoundingRadius)
	}
	if e.Flags&0x8 == 0x8 {
		enc.Float32(e.Scale)
	}
	if e.Flags&0x10 == 0x10 {
		enc.Int32(e.SoundNameRef)
	}
	enc.Int32(e.Unk2)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeModelRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &ModelRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.ActorDefRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.SphereRef = dec.Uint32()
	if d.Flags&0x1 == 0x1 {
		d.CurrentAction = dec.Uint32()
	}
	if d.Flags&0x2 == 0x2 {
		d.Offset.X = dec.Float32()
		d.Offset.Y = dec.Float32()
		d.Offset.Z = dec.Float32()
		d.Rotation.X = dec.Float32()
		d.Rotation.Y = dec.Float32()
		d.Rotation.Z = dec.Float32()
		d.Unk1 = dec.Uint32()
	}
	if d.Flags&0x4 == 0x4 {
		d.BoundingRadius = dec.Float32()
	}
	if d.Flags&0x8 == 0x8 {
		d.Scale = dec.Float32()
	}
	if d.Flags&0x10 == 0x10 {
		d.SoundNameRef = dec.Int32()
	}
	d.Unk2 = dec.Int32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// Sphere is Sphere in libeq, Zone Unknown in openzone, SPHERE (ref) in wld, Fragment16 in lantern
type Sphere struct {
	NameRef int32
	Radius  float32
}

func (e *Sphere) FragCode() int {
	return 0x16
}

func (e *Sphere) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Float32(e.Radius)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeSphere(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &Sphere{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Radius = dec.Float32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// SphereList is SphereListDef in libeq, empty in openzone, SPHERELISTDEFINITION in wld
type SphereList struct {
	NameRef     int32
	Flags       uint32
	SphereCount uint32
	Radius      float32
	Scale       float32
	Spheres     []common.Quad4
}

func (e *SphereList) FragCode() int {
	return 0x19
}

func (e *SphereList) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.SphereCount)
	enc.Float32(e.Radius)
	enc.Float32(e.Scale)
	for _, sphere := range e.Spheres {
		enc.Float32(sphere.X)
		enc.Float32(sphere.Y)
		enc.Float32(sphere.Z)
		enc.Float32(sphere.W)
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeSphereList(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &SphereList{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.SphereCount = dec.Uint32()
	d.Radius = dec.Float32()
	d.Scale = dec.Float32()
	for i := uint32(0); i < d.SphereCount; i++ {
		var sphere common.Quad4
		sphere.X = dec.Float32()
		sphere.Y = dec.Float32()
		sphere.Z = dec.Float32()
		sphere.W = dec.Float32()
		d.Spheres = append(d.Spheres, sphere)
	}

	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// SphereListRef is SphereList in libeq, empty in openzone, SPHERELIST (ref) in wld
type SphereListRef struct {
	NameRef          int32
	SphereListDefRef int32
	Params1          uint32
}

func (e *SphereListRef) FragCode() int {
	return 0x1A
}

func (e *SphereListRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(e.SphereListDefRef)
	enc.Uint32(e.Params1)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeSphereListRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &SphereListRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.SphereListDefRef = dec.Int32()
	d.Params1 = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}
