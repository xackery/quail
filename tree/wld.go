package tree

import (
	"github.com/xackery/quail/raw/rawfrag"
)

func fragRefs(fragment interface{}) []int32 {
	var refs []int32 // Initialize an empty slice for references

	switch frag := fragment.(type) {
	case *rawfrag.WldFragActor:
		refs = append(refs, frag.ActorDefRef)
		refs = append(refs, int32(frag.SphereRef)) // Cast uint32 to int32
		refs = append(refs, frag.SoundNameRef)
		refs = append(refs, frag.DMRGBTrackRef)

	case *rawfrag.WldFragActorDef:
		refs = append(refs, frag.BoundsRef)
		for _, ref := range frag.FragmentRefs {
			refs = append(refs, int32(ref)) // Cast each uint32 to int32
		}

	case *rawfrag.WldFragBlitSprite:
		refs = append(refs, frag.BlitSpriteRef)

	case *rawfrag.WldFragBlitSpriteDef:
		refs = append(refs, int32(frag.SpriteInstanceRef)) // Cast uint32 to int32

	case *rawfrag.WldFragDmRGBTrack:
		refs = append(refs, frag.TrackRef)

	case *rawfrag.WldFragDMSprite:
		refs = append(refs, frag.DMSpriteRef)

	case *rawfrag.WldFragDMSpriteDef:
		refs = append(refs, int32(frag.MaterialPaletteRef)) // Cast uint32 to int32

	case *rawfrag.WldFragDmSpriteDef2:
		refs = append(refs, int32(frag.MaterialPaletteRef)) // Cast uint32 to int32
		refs = append(refs, frag.DMTrackRef)
		refs = append(refs, frag.Fragment4Ref)

	case *rawfrag.WldFragDMTrack:
		refs = append(refs, frag.TrackRef)

	case *rawfrag.WldFragLight:
		refs = append(refs, frag.LightDefRef)

	case *rawfrag.WldFragHierarchicalSprite:
		refs = append(refs, int32(frag.HierarchicalSpriteRef)) // Cast uint32 to int32

	case *rawfrag.WldFragHierarchicalSpriteDef:
		refs = append(refs, int32(frag.CollisionVolumeRef)) // Cast uint32 to int32
		for _, ref := range frag.DMSprites {
			refs = append(refs, int32(ref)) // Cast each uint32 to int32
		}
		for _, dag := range frag.Dags { // Iterate over the Dags slice
			refs = append(refs, int32(dag.TrackRef))                  // Cast uint32 to int32
			refs = append(refs, int32(dag.MeshOrSpriteOrParticleRef)) // Cast uint32 to int32
		}

	case *rawfrag.WldFragMaterialDef:
		refs = append(refs, int32(frag.SimpleSpriteRef)) // Cast uint32 to int32

	case *rawfrag.WldFragMaterialPalette:
		for _, ref := range frag.MaterialRefs {
			refs = append(refs, int32(ref)) // Cast each uint32 to int32
		}

	case *rawfrag.WldFragParticleCloudDef:
		refs = append(refs, int32(frag.BlitSpriteRef)) // Cast uint32 to int32

	case *rawfrag.WldFragParticleSprite:
		refs = append(refs, frag.ParticleSpriteDefRef)

	case *rawfrag.WldFragParticleSpriteDef:
		refs = append(refs, int32(frag.RenderSimpleSpriteReference)) // Cast uint32 to int32

	case *rawfrag.WldFragPointLight:
		refs = append(refs, frag.LightRef)

	case *rawfrag.WldFragPointLightOldDef:
		refs = append(refs, frag.PointLightRef)

	case *rawfrag.WldFragPolyhedron:
		refs = append(refs, frag.FragmentRef)

	case *rawfrag.WldFragRegion:
		refs = append(refs, frag.AmbientLightRef)
		refs = append(refs, frag.MeshReference)

	case *rawfrag.WldFragSimpleSprite:
		refs = append(refs, int32(frag.SpriteRef)) // Cast uint32 to int32

	case *rawfrag.WldFragSimpleSpriteDef:
		for _, ref := range frag.BitmapRefs {
			refs = append(refs, int32(ref)) // Cast each uint32 to int32
		}

	case *rawfrag.WldFragSphereList:
		refs = append(refs, frag.SphereListDefRef)

	case *rawfrag.WldFragSprite2D:
		refs = append(refs, int32(frag.TwoDSpriteRef)) // Cast uint32 to int32

	case *rawfrag.WldFragSprite2DDef:
		refs = append(refs, int32(frag.SphereListRef)) // Cast uint32 to int32
		for _, pitch := range frag.Pitches {
			for _, heading := range pitch.Headings {
				refs = append(refs, heading.FrameRefs...)
			}
		}
		refs = append(refs, int32(frag.RenderSimpleSpriteReference)) // Cast uint32 to int32

	case *rawfrag.WldFragSprite3D:
		refs = append(refs, frag.Sprite3DDefRef)

	case *rawfrag.WldFragSprite3DDef:
		refs = append(refs, int32(frag.SphereListRef)) // Cast uint32 to int32
		for _, bspnode := range frag.BspNodes {
			refs = append(refs, int32(bspnode.RenderSimpleSpriteReference)) // Cast uint32 to int32
		}

	case *rawfrag.WldFragSprite4D:
		refs = append(refs, frag.FourDRef)

	case *rawfrag.WldFragSprite4DDef:
		refs = append(refs, frag.PolyRef)
		for _, ref := range frag.SpriteFragments {
			refs = append(refs, int32(ref)) // Cast each uint32 to int32
		}

	case *rawfrag.WldFragTrack:
		refs = append(refs, frag.TrackRef)
	}

	return refs
}
