package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
)

// WldFragDefaultPaletteFile is DefaultPaletteFile in libeq, empty in openzone, DEFAULTPALETTEFILE in wld
type WldFragDefaultPaletteFile struct {
	NameRef    int32  `yaml:"name_ref"`
	NameLength uint16 `yaml:"name_length"`
	FileName   string `yaml:"file_name"`
}

func (e *WldFragDefaultPaletteFile) FragCode() int {
	return FragCodeDefaultPaletteFile
}

func (e *WldFragDefaultPaletteFile) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint16(e.NameLength)
	enc.String(e.FileName)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragDefaultPaletteFile) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.NameLength = dec.Uint16()
	e.FileName = dec.StringFixed(int(e.NameLength))
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragBMInfo is BmInfo in libeq, Texture Bitmap Names in openzone, FRAME and BMINFO in wld, BitmapName in lantern
type WldFragBMInfo struct {
	NameRef      int32    `yaml:"name_ref"`
	TextureNames []string `yaml:"texture_names"`
}

func (e *WldFragBMInfo) FragCode() int {
	return FragCodeBMInfo
}

func (e *WldFragBMInfo) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(int32(len(e.TextureNames) - 1))
	for _, textureName := range e.TextureNames {
		enc.StringLenPrefixUint16(string(helper.WriteStringHash(textureName + "\x00")))
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragBMInfo) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	textureCount := dec.Int32()

	for i := 0; i < int(textureCount+1); i++ {
		nameLength := dec.Uint16()
		e.TextureNames = append(e.TextureNames, helper.ReadStringHash((dec.Bytes(int(nameLength)))))
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragSimpleSpriteDef is SimpleSpriteDef in libeq, WldFragSimpleSpriteDef Bitmap Info in openzone, SIMPLESPRITEDEF in wld, BitmapInfo in lantern
type WldFragSimpleSpriteDef struct {
	NameRef      int32    `yaml:"name_ref"`
	Flags        uint32   `yaml:"flags"`
	CurrentFrame int32    `yaml:"current_frame"`
	Sleep        uint32   `yaml:"sleep"`
	BitmapRefs   []uint32 `yaml:"bitmap_refs"`
}

func (e *WldFragSimpleSpriteDef) FragCode() int {
	return FragCodeSimpleSpriteDef
}

func (e *WldFragSimpleSpriteDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.BitmapRefs)))
	if e.Flags&0x20 != 0 {
		enc.Int32(e.CurrentFrame)
	}
	if e.Flags&0x08 != 0 && e.Flags&0x10 != 0 {
		enc.Uint32(e.Sleep)
	}
	for _, textureRef := range e.BitmapRefs {
		enc.Uint32(textureRef)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSimpleSpriteDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	textureRefCount := dec.Uint32()
	if e.Flags&0x20 != 0 {
		e.CurrentFrame = dec.Int32()
	}
	if e.Flags&0x08 != 0 && e.Flags&0x10 != 0 {
		e.Sleep = dec.Uint32()
	}
	for i := 0; i < int(textureRefCount); i++ {
		e.BitmapRefs = append(e.BitmapRefs, dec.Uint32())
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragSimpleSprite is SimpleSprite in libeq, Texture Bitmap Info Reference in openzone, SIMPLESPRITEINST in wld, BitmapInfoReference in lantern
type WldFragSimpleSprite struct {
	NameRef   int32  `yaml:"name_ref"`
	SpriteRef int16  `yaml:"sprite_ref"`
	Flags     uint32 `yaml:"flags"`
}

func (e *WldFragSimpleSprite) FragCode() int {
	return FragCodeSimpleSprite
}

func (e *WldFragSimpleSprite) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int16(e.SpriteRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragSimpleSprite) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.SpriteRef = dec.Int16()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragBlitSpriteDef is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSprite in lantern
type WldFragBlitSpriteDef struct {
	NameRef           int32  `yaml:"name_ref"`
	Flags             uint32 `yaml:"flags"`
	SpriteInstanceRef uint32 `yaml:"sprite_instance_ref"`
	Unknown           int32  `yaml:"unknown"`
}

func (e *WldFragBlitSpriteDef) FragCode() int {
	return FragCodeBlitSpriteDef
}

func (e *WldFragBlitSpriteDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.SpriteInstanceRef)
	enc.Int32(e.Unknown)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragBlitSpriteDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.SpriteInstanceRef = dec.Uint32()
	e.Unknown = dec.Int32()

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragBlitSprite is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSpriteReference in lantern
type WldFragBlitSprite struct {
}

func (e *WldFragBlitSprite) FragCode() int {
	return FragCodeBlitSprite
}

func (e *WldFragBlitSprite) Write(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func (e *WldFragBlitSprite) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragMaterialDef is MaterialDef in libeq, Texture in openzone, MATERIALDEFINITION in wld, Material in lantern
type WldFragMaterialDef struct {
	NameRef           int32     `yaml:"name_ref"`
	Flags             uint32    `yaml:"flags"`
	RenderMethod      uint32    `yaml:"render_method"`
	RGBPen            uint32    `yaml:"rgb_pen"`
	Brightness        float32   `yaml:"brightness"`
	ScaledAmbient     float32   `yaml:"scaled_ambient"`
	SpriteInstanceRef uint32    `yaml:"sprite_instance_ref"`
	Pairs             [2]uint32 `yaml:"pairs"`
}

func (e *WldFragMaterialDef) FragCode() int {
	return FragCodeMaterialDef
}

func (e *WldFragMaterialDef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.RenderMethod)
	enc.Uint32(e.RGBPen)
	enc.Float32(e.Brightness)
	enc.Float32(e.ScaledAmbient)
	enc.Uint32(e.SpriteInstanceRef)
	if e.Flags&0x1 != 0 {
		enc.Uint32(e.Pairs[0])
		enc.Uint32(e.Pairs[1])
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragMaterialDef) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.RenderMethod = dec.Uint32()
	e.RGBPen = dec.Uint32()
	e.Brightness = dec.Float32()
	e.ScaledAmbient = dec.Float32()
	e.SpriteInstanceRef = dec.Uint32()
	if e.Flags&0x1 != 0 {
		e.Pairs[0] = dec.Uint32()
		e.Pairs[1] = dec.Uint32()
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragMaterialPalette is MaterialPalette in libeq, TextureList in openzone, MATERIALPALETTE in wld, WldFragMaterialPalette in lantern
type WldFragMaterialPalette struct {
	NameRef      int32
	Flags        uint32
	MaterialRefs []uint32
}

func (e *WldFragMaterialPalette) FragCode() int {
	return FragCodeMaterialPalette
}

func (e *WldFragMaterialPalette) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.MaterialRefs)))
	for _, materialRef := range e.MaterialRefs {
		enc.Uint32(materialRef)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragMaterialPalette) Read(r io.ReadSeeker) error {
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	materialCount := dec.Uint32()
	for i := 0; i < int(materialCount); i++ {
		e.MaterialRefs = append(e.MaterialRefs, dec.Uint32())
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}
