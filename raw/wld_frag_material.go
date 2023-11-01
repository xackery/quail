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

func (e *WldFragPaletteFile) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint16(e.NameLength)
	enc.String(e.FileName)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func readPaletteFile(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragPaletteFile{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.NameLength = dec.Uint16()
	d.FileName = dec.StringFixed(int(d.NameLength))
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
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

func (e *WldFragTextureList) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int32(int32(len(e.TextureNames) - 1))
	for _, textureName := range e.TextureNames {
		enc.StringLenPrefixUint16(string(helper.WriteStringHash(textureName + "\x00")))
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func readTextureList(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragTextureList{}
	d.FragName = FragName(d.FragCode())
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

func (e *WldFragTexture) Encode(w io.Writer) error {
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
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func readTexture(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragTexture{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	textureCount := dec.Uint32()
	if d.Flags&0x20 != 0 {
		d.TextureCurrent = dec.Uint32()
	}
	if d.Flags&0x08 != 0 && d.Flags&0x10 != 0 {
		d.Sleep = dec.Uint32()
	}
	for i := 0; i < int(textureCount); i++ {
		d.TextureRefs = append(d.TextureRefs, dec.Uint32())
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
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

func (e *WldFragTextureRef) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Int16(e.TextureRef)
	enc.Uint32(e.Flags)
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func readTextureRef(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragTextureRef{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.TextureRef = dec.Int16()
	d.Flags = dec.Uint32()
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// WldFragBlitSprite is WldFragBlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSprite in lantern
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

func (e *WldFragBlitSprite) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(e.BlitSpriteRef)
	enc.Int32(e.Unk1)
	if enc.Error() != nil {
		return enc.Error()
	}

	return nil
}

func readBlitSprite(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragBlitSprite{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	d.BlitSpriteRef = dec.Uint32()
	d.Unk1 = dec.Int32()

	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}

// WldFragBlitSpriteRef is BlitSprite in libeq, empty in openzone, BLITSPRITE (ref) in wld, ParticleSpriteReference in lantern
type WldFragBlitSpriteRef struct {
	FragName string `yaml:"frag_name"`
}

func (e *WldFragBlitSpriteRef) FragCode() int {
	return 0x27
}

func (e *WldFragBlitSpriteRef) Encode(w io.Writer) error {
	return fmt.Errorf("not implemented")
}

func readBlitSpriteRef(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragBlitSpriteRef{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
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

func (e *WldFragMaterial) Encode(w io.Writer) error {
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

func readMaterial(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragMaterial{}
	d.FragName = FragName(d.FragCode())
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

func (e *WldFragMaterialList) Encode(w io.Writer) error {
	enc := encdec.NewEncoder(w, binary.LittleEndian)
	enc.Int32(e.NameRef)
	enc.Uint32(e.Flags)
	enc.Uint32(uint32(len(e.MaterialRefs)))
	for _, materialRef := range e.MaterialRefs {
		enc.Uint32(materialRef)
	}
	if enc.Error() != nil {
		return enc.Error()
	}
	return nil
}

func readMaterialList(r io.ReadSeeker) (FragmentReader, error) {
	d := &WldFragMaterialList{}
	d.FragName = FragName(d.FragCode())
	dec := encdec.NewDecoder(r, binary.LittleEndian)
	d.NameRef = dec.Int32()
	d.Flags = dec.Uint32()
	materialCount := dec.Uint32()
	for i := 0; i < int(materialCount); i++ {
		d.MaterialRefs = append(d.MaterialRefs, dec.Uint32())
	}
	if dec.Error() != nil {
		return nil, dec.Error()
	}
	return d, nil
}
