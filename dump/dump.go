package dump

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"
	"os"

	"github.com/xackery/colors"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

const (
	width  = 550
	height = 550
)

type Dump struct {
	img            *image.RGBA
	groups         map[string]*group
	groupKeys      []string
	lastGroupColor int
	lastPos        image.Point
	drw            font.Drawer
	grayImage      *image.Uniform
	rowLabel       int
}

var groupColors = []color.RGBA{
	colors.GreenYellow,
	colors.BlanchedAlmond,
	colors.Aqua,
	colors.OrangeRed,
	colors.Azure,
	colors.Pink,
	colors.Lavender,
	colors.Bisque,
	colors.MediumAquaMarine,
	colors.Wheat,
}

var instance *Dump

type group struct {
	mask *image.RGBA
}

func New() (*Dump, error) {
	e := &Dump{
		img:    image.NewRGBA(image.Rect(0, 0, width, height)),
		groups: make(map[string]*group),
	}
	draw.Draw(e.img, e.img.Bounds(), image.Black, image.Point{X: 0, Y: 0}, draw.Src)
	e.grayImage = image.NewUniform(colors.Gray)

	e.lastPos.Y = 16
	e.drw = font.Drawer{
		Dst:  e.img,
		Src:  e.grayImage,
		Face: inconsolata.Bold8x16,
		Dot:  fixed.P(e.lastPos.X, e.lastPos.Y),
	}
	e.lastPos.X = 1
	e.drw.Dot = fixed.P(e.lastPos.X, e.lastPos.Y)
	e.drw.DrawString("   00 01 02 03 04 05 06 07 08 09 0A 0B 0C 0D 0E 0F  LABELS")
	e.lastPos.Y += 16
	e.lastPos.X = 0
	e.addLabel(e.grayImage, "00 ")
	if instance == nil {
		instance = e
	}
	return e, nil
}

func Hex(data interface{}, format string, a ...interface{}) {
	if instance == nil {
		return
	}
	instance.Hex(data, format, a...)
}

func (e *Dump) Hex(data interface{}, format string, a ...interface{}) {

	buf := bytes.NewBuffer(nil)
	err := binary.Write(buf, binary.LittleEndian, data)
	if err != nil {

	}

	labelName := fmt.Sprintf(format, a...)

	grp, ok := e.groups[labelName]
	if !ok {
		if len(groupColors) <= e.lastGroupColor {
			e.lastGroupColor = 0
		}
		img := image.NewRGBA(image.Rect(0, 0, width, height))

		clr := groupColors[e.lastGroupColor]
		for i := 0; i < width; i++ {
			for j := 0; j < height; j++ {
				img.SetRGBA(i, j, clr)
			}
		}

		grp = &group{
			mask: img,
		}
		e.groupKeys = append(e.groupKeys, labelName)
		e.groups[labelName] = grp
		e.lastGroupColor++
	}
	for _, d := range buf.Bytes() {
		e.addLabel(grp.mask, fmt.Sprintf("%02X ", d))
	}
}

func (e *Dump) addLabel(mask image.Image, text string) {

	e.drw.Src = mask
	e.drw.Dot = fixed.P(e.lastPos.X, e.lastPos.Y)
	e.drw.DrawString(text)

	adv := e.drw.MeasureString(text)
	e.lastPos.X += int(adv) / 62
	if e.lastPos.X >= 400 {
		e.lastPos.X = 0
		e.rowLabel += 16
		e.lastPos.Y += 16
		e.addLabel(e.grayImage, fmt.Sprintf("%02X ", e.rowLabel))
	}
}

func (e *Dump) Save(path string) error {
	out := image.NewRGBA(e.img.Bounds())
	draw.Draw(out, e.img.Bounds(), e.img, e.img.Bounds().Min, draw.Src)
	d := font.Drawer{
		Dst:  out,
		Src:  image.White,
		Face: inconsolata.Bold8x16,
		Dot:  fixed.P(e.lastPos.X, e.lastPos.Y),
	}
	pos := 32

	for _, name := range e.groupKeys {
		d.Src = e.groups[name].mask
		d.Dot = fixed.P(420, pos)
		d.DrawString(name)
		pos += 16
	}

	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create: %w", err)
	}
	defer f.Close()
	err = png.Encode(f, out)
	if err != nil {
		return fmt.Errorf("encode: %w", err)
	}
	return nil
}
