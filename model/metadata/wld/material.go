package wld

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/common"
	"github.com/xackery/quail/helper"
)

// PaletteFile is DefaultPaletteFile in libeq, empty in openzone, DEFAULTPALETTEFILE in wld
type PaletteFile struct {
	NameRef    int32
	NameLength uint16
	FileName   string
}

func (e *PaletteFile) FragCode() int {
	return 0x01
}

func (e *PaletteFile) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint16(e.NameLength)
	enc.String(e.FileName)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodePaletteFile(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &PaletteFile{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.NameLength = dec.Uint16()
	d.FileName = dec.StringFixed(int(d.NameLength))
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// TextureList is BmInfo in libeq, Texture Bitmap Names in openzone, FRAME and BMINFO in wld, BitmapName in lantern
type TextureList struct {
	NameRef      int32
	TextureNames []string
}

func (e *TextureList) FragCode() int {
	return 0x03
}

func (e *TextureList) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(int32(len(e.TextureNames) - 1))
	for _, textureName := range e.TextureNames {
		enc.StringLenPrefixUint16(string(helper.WriteStringHash(textureName)))
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeTextureList(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &TextureList{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	textureCount := dec.Int32()

	for i := 0; i < int(textureCount+1); i++ {
		nameLength := dec.Uint16()
		d.TextureNames = append(d.TextureNames, helper.ReadStringHash((dec.Bytes(int(nameLength)))))
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// Texture is SimpleSpriteDef in libeq, Texture Bitmap Info in openzone, SIMPLESPRITEDEF in wld, BitmapInfo in lantern
type Texture struct {
	NameRef        int32
	Flags          uint32
	TextureCount   uint32
	TextureCurrent uint32
	Sleep          uint32
	TextureRefs    []uint32
}

func (e *Texture) FragCode() int {
	return 0x04
}

func (e *Texture) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.TextureCount)
	if e.Flags&0x20 != 0 {
		enc.Uint32(e.TextureCurrent)
	}
	if e.Flags&0x08 != 0 && e.Flags&0x10 != 0 {
		enc.Uint32(e.Sleep)
	}
	for _, textureRef := range e.TextureRefs {
		enc.Uint32(textureRef)
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeTexture(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &Texture{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.TextureCount = dec.Uint32()
	if d.Flags&0x20 != 0 {
		d.TextureCurrent = dec.Uint32()
	}
	if d.Flags&0x08 != 0 && d.Flags&0x10 != 0 {
		d.Sleep = dec.Uint32()
	}
	for i := 0; i < int(d.TextureCount); i++ {
		d.TextureRefs = append(d.TextureRefs, dec.Uint32())
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// TextureRef is SimpleSprite in libeq, Texture Bitmap Info Reference in openzone, SIMPLESPRITEINST in wld, BitmapInfoReference in lantern
type TextureRef struct {
	NameRef    int32
	TextureRef int16
	Flags      uint32
}

func (e *TextureRef) FragCode() int {
	return 0x05
}

func (e *TextureRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int16(e.TextureRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeTextureRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &TextureRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.TextureRef = dec.Int16()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// BlitSprite is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSprite in lantern
type BlitSprite struct {
}

func (e *BlitSprite) FragCode() int {
	return 0x26
}

func (e *BlitSprite) Encode(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func decodeBlitSprite(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &BlitSprite{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// BlitSpriteRef is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSpriteReference in lantern
type BlitSpriteRef struct {
}

func (e *BlitSpriteRef) FragCode() int {
	return 0x27
}

func (e *BlitSpriteRef) Encode(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func decodeBlitSpriteRef(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &BlitSpriteRef{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// Material is MaterialDef in libeq, Texture in openzone, MATERIALDEFINITION in wld, Material in lantern
type Material struct {
	NameRef       int32
	Flags         uint32
	RenderMethod  uint32
	RGBPen        uint32
	Brightness    float32
	ScaledAmbient float32
	TextureRef    uint32
	Pairs         [2]uint32
}

func (e *Material) FragCode() int {
	return 0x30
}

func (e *Material) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.RenderMethod)
	enc.Uint32(e.RGBPen)
	enc.Float32(e.Brightness)
	enc.Float32(e.ScaledAmbient)
	enc.Uint32(e.TextureRef)
	if e.Flags&0x1 != 0 {
		enc.Uint32(e.Pairs[0])
		enc.Uint32(e.Pairs[1])
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeMaterial(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &Material{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.RenderMethod = dec.Uint32()
	d.RGBPen = dec.Uint32()
	d.Brightness = dec.Float32()
	d.ScaledAmbient = dec.Float32()
	d.TextureRef = dec.Uint32()
	if d.Flags&0x1 != 0 {
		d.Pairs[0] = dec.Uint32()
		d.Pairs[1] = dec.Uint32()
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// MaterialList is MaterialPalette in libeq, TextureList in openzone, MATERIALPALETTE in wld, MaterialList in lantern
type MaterialList struct {
	NameRef       int32
	Flags         uint32
	MaterialCount uint32
	MaterialRefs  []uint32
}

func (e *MaterialList) FragCode() int {
	return 0x31
}

func (e *MaterialList) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.MaterialCount)
	for _, materialRef := range e.MaterialRefs {
		enc.Uint32(materialRef)
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func decodeMaterialList(r io.ReadSeeker) (common.FragmentReader, error) {
	d := &MaterialList{}
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.MaterialCount = dec.Uint32()
	for i := 0; i < int(d.MaterialCount); i++ {
		d.MaterialRefs = append(d.MaterialRefs, dec.Uint32())
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}
