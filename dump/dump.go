// dump is used to generate png files based on hex data of a file
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
	"path/filepath"

	"github.com/xackery/colors"
	"golang.org/x/image/font"
	"golang.org/x/image/font/inconsolata"
	"golang.org/x/image/math/fixed"
)

const (
	width  = 800
	height = 4000
)

type Dump struct {
	name           string
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

func New(name string) (*Dump, error) {
	e := &Dump{
		name:   filepath.Base(name),
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
	instance.Hex(data, 0, format, a...)
}

func HexRange(data interface{}, size int, format string, a ...interface{}) {
	if instance == nil {
		return
	}
	instance.Hex(data, size, format, a...)
}

func IsActive() bool {
	return instance != nil
}

func (e *Dump) Hex(data interface{}, size int, format string, a ...interface{}) {
	buf := bytes.NewBuffer(nil)
	binary.Write(buf, binary.LittleEndian, data)

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

	if len(buf.Bytes()) >= 16 || size >= 16 {
		e.addLabel(grp.mask, fmt.Sprintf("%02X ", buf.Bytes()[0]))
		e.addLabel(grp.mask, ".. ")
		e.addLabel(grp.mask, fmt.Sprintf("%02X ", buf.Bytes()[buf.Len()-1]))

		fmt.Printf("inspect: dump %d bytes (truncated to %d) for %s\n", size, 3, labelName)
		return
	}
	for _, d := range buf.Bytes() {
		e.addLabel(grp.mask, fmt.Sprintf("%02X ", d))
	}
	if size == 0 {
		size = buf.Len()
	}
	fmt.Printf("inspect: dump %d bytes for %s\n", size, labelName)
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

	if pos < e.lastPos.Y {
		pos = e.lastPos.Y
	}

	d.Src = e.grayImage
	d.Dot = fixed.P(200, e.lastPos.Y+16)
	d.DrawString(e.name)

	pos += 16
	if pos < height {
		b := e.img.Bounds()
		b.Max.Y = pos
		newOut := image.NewRGBA(b)
		draw.Draw(newOut, e.img.Bounds(), out, newOut.Bounds().Min, draw.Src)
		out = newOut
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
	fmt.Println("saved dump", path)
	return nil
}

func Close() {
	if instance != nil {
		instance = nil
	}
}
