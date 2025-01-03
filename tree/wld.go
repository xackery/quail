package tree

import (
	"fmt"
	"io"

	"github.com/xackery/quail/raw"
	"github.com/xackery/quail/raw/rawfrag"
)

func wldDump(wld *raw.Wld, w io.Writer) error {
	dumpedFrags := make(map[int]bool)

	fmt.Println("wldDump")
	for i, frag := range wld.Fragments {
		switch val := frag.(type) {
		case *rawfrag.WldFragActorDef:
			fmt.Fprintf(w, "ID: %d %s (actordef)\n", i, wld.Name(val.NameRef))
			dumpedFrags[i] = true
			for _, fRef := range val.FragmentRefs {
				spriteRef := wld.Fragments[fRef]
				switch val2 := spriteRef.(type) {
				case *rawfrag.WldFragSprite2D:
					fmt.Fprintf(w, "  ID: %d %s (sprite2d)\n", fRef, wld.Name(val2.NameRef))
				case *rawfrag.WldFragSprite3D:
					fmt.Fprintf(w, "  ID: %d %s (sprite3d)\n", fRef, wld.Name(val2.NameRef))
				case *rawfrag.WldFragBlitSprite:
					fmt.Fprintf(w, "  ID: %d %s (blitsprite)\n", fRef, wld.Name(val2.NameRef))
				case *rawfrag.WldFragDMSprite:
					fmt.Fprintf(w, "  ID: %d %s (dmsprite)\n", fRef, wld.Name(val2.NameRef))
				default:
					fmt.Fprintf(w, "  ID: %d (unknown)\n", fRef)
				}
				dumpedFrags[int(fRef)] = true

			}
		}
	}

	return nil
}
