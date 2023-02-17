package fragment

const (
	// Used for boundaries that are not rendered. TextInfoReference can be null or have reference.
	MaterialTypeBoundary = 0x0
	// Standard diffuse shader
	MaterialTypeDiffuse = 0x01
	// Diffuse variant
	MaterialTypeDiffuse2 = 0x02
	// Transparent with 0.5 blend strength
	MaterialTypeTransparent50 = 0x05
	// Transparent with 0.25 blend strength
	MaterialTypeTransparent25 = 0x09
	// Transparent with 0.75 blend strength
	MaterialTypeTransparent75 = 0x0A
	// Non solid surfaces that shouldn't really be masked
	MaterialTypeTransparentMaskedPassable       = 0x07
	MaterialTypeTransparentAdditiveUnlit        = 0x0B
	MaterialTypeTransparentMasked               = 0x13
	MaterialTypeDiffuse3                        = 0x14
	MaterialTypeDiffuse4                        = 0x15
	MaterialTypeTransparentAdditive             = 0x17
	MaterialTypeDiffuse5                        = 0x19
	MaterialTypeInvisibleUnknown                = 0x53
	MaterialTypeDiffuse6                        = 0x553
	MaterialTypeCompleteUnknown                 = 0x1A // TODO: Analyze this
	MaterialTypeDiffuse7                        = 0x12
	MaterialTypeDiffuse8                        = 0x31
	MaterialTypeInvisibleUnknown2               = 0x4B
	MaterialTypeDiffuseSkydome                  = 0x0D // Need to confirm
	MaterialTypeTransparentSkydome              = 0x0F // Need to confirm
	MaterialTypeTransparentAdditiveUnlitSkydome = 0x10
	MaterialTypeInvisibleUnknown3               = 0x03
)
