package quail

import (
	"io"

	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/mesh/wld"
	"github.com/xackery/quail/pfs/archive"
	"github.com/xackery/quail/quail/def"
)

// Decode decodes a WLD file
func WLDDecode(r io.ReadSeeker, pfs archive.ReadWriter) ([]*def.Mesh, error) {
	meshes := make([]*def.Mesh, 0)

	e, err := wld.New("test", pfs)
	if err != nil {
		return nil, err
	}

	err = e.Decode(r)
	if err != nil {
		return nil, err
	}

	//names := e.Names()
	for _, f := range e.Fragments {
		switch d := f.(type) {
		case *def.Mesh:
			meshes = append(meshes, d)
			//ref := e.Fragments[d.MaterialListRef].(*wld.MaterialList)
		case *wld.MaterialList:
			log.Debugf("material list")
		}
	}

	return meshes, nil
}
