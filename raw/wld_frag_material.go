package raw

import (
	"encoding/binary"
	"fmt"
	"io"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
)

// WldFragPaletteFile is DefaultPaletteFile in libeq, empty in openzone, DEFAULTPALETTEFILE in wld
type WldFragPaletteFile struct {
	FragName   string `yaml:"frag_name"`
	NameRef    int32  `yaml:"name_ref"`
	NameLength uint16 `yaml:"name_length"`
	FileName   string `yaml:"file_name"`
}

func (e *WldFragPaletteFile) FragCode() int {
	return 0x01
}

func (e *WldFragPaletteFile) Write(w io.Writer) error {
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

func (e *WldFragPaletteFile) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
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

// WldFragTextureList is BmInfo in libeq, Texture Bitmap Names in openzone, FRAME and BMINFO in wld, BitmapName in lantern
type WldFragTextureList struct {
	FragName     string   `yaml:"frag_name"`
	NameRef      int32    `yaml:"name_ref"`
	TextureNames []string `yaml:"texture_names"`
}

func (e *WldFragTextureList) FragCode() int {
	return 0x03
}

func (e *WldFragTextureList) Write(w io.Writer) error {
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

func (e *WldFragTextureList) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
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

// WldFragTexture is SimpleSpriteDef in libeq, WldFragTexture Bitmap Info in openzone, SIMPLESPRITEDEF in wld, BitmapInfo in lantern
type WldFragTexture struct {
	FragName       string   `yaml:"frag_name"`
	NameRef        int32    `yaml:"name_ref"`
	Flags          uint32   `yaml:"flags"`
	TextureCurrent uint32   `yaml:"texture_current"`
	Sleep          uint32   `yaml:"sleep"`
	TextureRefs    []uint32 `yaml:"texture_refs"`
}

func (e *WldFragTexture) FragCode() int {
	return 0x04
}

func (e *WldFragTexture) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.TextureRefs)))
	if e.Flags&0x20 != 0 {
		enc.Uint32(e.TextureCurrent)
	}
	if e.Flags&0x08 != 0 && e.Flags&0x10 != 0 {
		enc.Uint32(e.Sleep)
	}
	for _, textureRef := range e.TextureRefs {
		enc.Uint32(textureRef)
	}
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragTexture) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	textureRefCount := dec.Uint32()
	if e.Flags&0x20 != 0 {
		e.TextureCurrent = dec.Uint32()
	}
	if e.Flags&0x08 != 0 && e.Flags&0x10 != 0 {
		e.Sleep = dec.Uint32()
	}
	for i := 0; i < int(textureRefCount); i++ {
		e.TextureRefs = append(e.TextureRefs, dec.Uint32())
	}
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragTextureRef is SimpleSprite in libeq, Texture Bitmap Info Reference in openzone, SIMPLESPRITEINST in wld, BitmapInfoReference in lantern
type WldFragTextureRef struct {
	FragName   string `yaml:"frag_name"`
	NameRef    int32  `yaml:"name_ref"`
	TextureRef int16  `yaml:"texture_ref"`
	Flags      uint32 `yaml:"flags"`
}

func (e *WldFragTextureRef) FragCode() int {
	return 0x05
}

func (e *WldFragTextureRef) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int16(e.TextureRef)
	enc.Uint32(e.Flags)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragTextureRef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.TextureRef = dec.Int16()
	e.Flags = dec.Uint32()
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragBlitSprite is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSprite in lantern
type WldFragBlitSprite struct {
	FragName      string `yaml:"frag_name"`
	NameRef       int32  `yaml:"name_ref"`
	Flags         uint32 `yaml:"flags"`
	BlitSpriteRef uint32 `yaml:"blit_sprite_ref"`
	Unk1          int32  `yaml:"unk1"`
}

func (e *WldFragBlitSprite) FragCode() int {
	return 0x26
}

func (e *WldFragBlitSprite) Write(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.BlitSpriteRef)
	enc.Int32(e.Unk1)
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}

	return nil
}

func (e *WldFragBlitSprite) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.BlitSpriteRef = dec.Uint32()
	e.Unk1 = dec.Int32()

	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragBlitSpriteRef is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSpriteReference in lantern
type WldFragBlitSpriteRef struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragBlitSpriteRef) FragCode() int {
	return 0x27
}

func (e *WldFragBlitSpriteRef) Write(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func (e *WldFragBlitSpriteRef) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	err := dec.Error()
	if err != nil {
		return fmt.Errorf("read: %w", err)
	}
	return nil
}

// WldFragMaterial is MaterialDef in libeq, Texture in openzone, MATERIALDEFINITION in wld, Material in lantern
type WldFragMaterial struct {
	FragName      string    `yaml:"frag_name"`
	NameRef       int32     `yaml:"name_ref"`
	Flags         uint32    `yaml:"flags"`
	RenderMethod  uint32    `yaml:"render_method"`
	RGBPen        uint32    `yaml:"rgb_pen"`
	Brightness    float32   `yaml:"brightness"`
	ScaledAmbient float32   `yaml:"scaled_ambient"`
	TextureRef    uint32    `yaml:"texture_ref"`
	Pairs         [2]uint32 `yaml:"pairs"`
}

func (e *WldFragMaterial) FragCode() int {
	return 0x30
}

func (e *WldFragMaterial) Write(w io.Writer) error {
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
	err := enc.Error()
	if err != nil {
		return fmt.Errorf("write: %w", err)
	}
	return nil
}

func (e *WldFragMaterial) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	e.NameRef = dec.Int32()
	e.Flags = dec.Uint32()
	e.RenderMethod = dec.Uint32()
	e.RGBPen = dec.Uint32()
	e.Brightness = dec.Float32()
	e.ScaledAmbient = dec.Float32()
	e.TextureRef = dec.Uint32()
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

// WldFragMaterialList is MaterialPalette in libeq, TextureList in openzone, MATERIALPALETTE in wld, WldFragMaterialList in lantern
type WldFragMaterialList struct {
	FragName     string
	NameRef      int32
	Flags        uint32
	MaterialRefs []uint32
}

func (e *WldFragMaterialList) FragCode() int {
	return 0x31
}

func (e *WldFragMaterialList) Write(w io.Writer) error {
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

func (e *WldFragMaterialList) Read(r io.ReadSeeker) error {
	e.FragName = FragName(e.FragCode())
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
