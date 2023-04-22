package zon

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"strings"

	"github.com/xackery/encdec"
	"github.com/xackery/quail/helper"
	"github.com/xackery/quail/log"
	"github.com/xackery/quail/model/geo"
)

func (e *ZON) eqgtzpDecode(r io.ReadSeeker) error {
	scanner := bufio.NewScanner(r)
	scanner.Split(func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		if atEOF && len(data) == 0 {
			return 0, nil, nil
		}
		if i := bytes.IndexByte(data, '\n'); i >= 0 {
			// We have a full newline-terminated line.
			return i + 1, data[0 : i+1], nil
		}
		// If we're at EOF, we have a final, non-terminated line. Return it.
		if atEOF {
			return len(data), data, nil
		}
		// Request more data.
		return 0, nil, nil
	})

	minLongitude := int32(0)
	maxLongitude := int32(0)
	minLatitude := int32(0)
	maxLatitude := int32(0)
	minExtents := geo.Vector3{}
	maxExtents := geo.Vector3{}
	unitsPerVert := float32(0)
	quadsPerTile := int32(0)
	coverMapInputSize := int32(0)
	layeringMapInputSize := int32(0)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "*") {
			continue
		}
		line = strings.TrimPrefix(line, "*")
		line = strings.TrimSpace(line)
		parts := strings.Split(line, " ")
		switch parts[0] {
		case "NAME":
			e.models = append(e.models, model{name: parts[1], baseName: parts[1]})
		case "VERSION":
			e.version = helper.AtoU32(parts[1])
			if e.version == 0 {
				return fmt.Errorf("invalid version on eqtzp: %s", parts[1])
			}
		case "MINLNG":
			minLongitude = helper.AtoI32(parts[1])
		case "MAXLNG":
			maxLongitude = helper.AtoI32(parts[1])
		case "MINLAT":
			minLatitude = helper.AtoI32(parts[1])
		case "MAXLAT":
			maxLatitude = helper.AtoI32(parts[1])
		case "MIN_EXTENTS":
			minExtents.X = helper.AtoF32(parts[1])
			minExtents.Y = helper.AtoF32(parts[2])
			minExtents.Z = helper.AtoF32(parts[3])
		case "MAX_EXTENTS":
			maxExtents.X = helper.AtoF32(parts[1])
			maxExtents.Y = helper.AtoF32(parts[2])
			maxExtents.Z = helper.AtoF32(parts[3])
		case "UNITSPERVERT":
			unitsPerVert = helper.AtoF32(parts[1])
		case "QUADSPERTILE":
			quadsPerTile = helper.AtoI32(parts[1])
		case "COVERMAPINPUTSIZE":
			coverMapInputSize = helper.AtoI32(parts[1])
		case "LAYERINGMAPINPUTSIZE":
			layeringMapInputSize = helper.AtoI32(parts[1])
		}
	}

	log.Debugf("minLongitude: %d, maxLongitude: %d, minLatitude: %d, maxLatitude: %d, unitsPerVert: %f, quadsPerTile: %d, coverMapInputSize: %d, layeringMapInputSize: %d", minLongitude, maxLongitude, minLatitude, maxLatitude, unitsPerVert, quadsPerTile, coverMapInputSize, layeringMapInputSize)
	if e.pfs == nil {
		return fmt.Errorf("pfs required for newer zones")
	}

	datName := fmt.Sprintf("%s.dat", strings.TrimSuffix(e.name, ".zon"))

	df, err := e.pfs.File(datName)
	if err != nil {
		return fmt.Errorf("pfs file %s: %w", datName, err)
	}

	dec := encdec.NewDecoder(bytes.NewReader(df), binary.LittleEndian)

	header1 := dec.Uint32()
	header2 := dec.Uint32()
	header3 := dec.Uint32()
	log.Debugf("header1 to 3: %d, %d, %d", header1, header2, header3)
	/*
		baseTileTextureName := dec.StringLenPrefixUint32()
		tileCount := dec.Uint32()

		zonMin := geo.Vector2{
			X: float32(minLatitude) * quadsPerTile * unitsPerVert,
			Y: float32(minLongitude) * quadsPerTile * unitsPerVert,
		}

		quadCount := quadsPerTile * quadsPerTile
		vertCount := (quadsPerTile + 1) * (quadsPerTile + 1)

		type tile struct {
			longitude      uint32
			latitude       uint32
			unknown        uint32
			startX         float32
			startY         float32
			floats         []float32
			colors         []uint32
			colors2        []uint32
			flags          []uint32
			baseWaterLevel float32
		}

		for i := 0; i < int(tileCount); i++ {
			tile := tile{}
			tileLongitute := dec.Uint32()
			tileLatitude := dec.Uint32()
			tileUnknown := dec.Uint32()
			tileStartY := zonMin.Y + float32(tileLongitute-100000)*unitsPerVert*float32(quadsPerTile)
			tileStartX := zonMin.X + float32(tileLatitude-100000)*unitsPerVert*float32(quadsPerTile)

			floatsAllTheSame := false
			currentAvg := double(0.0)
			for j := 0; j < int(vertCount); j++ {
				val := dec.Float32()
				if j == 0 {
					currentAvg = float
				}
				if float != currentAvg {
					floatsAllTheSame = false
				}
				tile.floats = append(tile.floats, float)
			}

			for j := 0; j < int(vertCount); j++ {
				tile.colors = append(tile.colors, dec.Uint32())
			}

			for j := 0; j < int(vertCount); j++ {
				tile.colors2 = append(tile.colors2, dec.Uint32())
			}

			for j := 0; j < int(quadCount); j++ {
				tile.flags = append(tile.flags, dec.Uint32())
				if tile.flags[j] != 0 {
					floatsAllTheSame = false
				}
			}

			tile.baseWaterLevel = dec.Float32()
			unk := dec.Bytes(4)
			if unk[0] != 0 || unk[1] != 0 || unk[2] != 0 || unk[3] != 0 {
				unkByte := dec.Uint8()
				if unkByte != 0 {
					unkFloat := dec.Float32()
					unkFloat2 := dec.Float32()
					unkFloat3 := dec.Float32()
					unkFloat4 := dec.Float32()
				}
				unkFloat := dec.Float32()
			}

			layerCount := dec.Uint32()
			if layerCount > 0 {
				baseMaterial := dec.StringLenPrefixUint32()
				for j := 1; j < int(layerCount); j++ {
					material := dec.String()
					detailMaskDim := dec.Uint32()
					detailMaskSize := detailMaskDim * detailMaskDim
					detailMask := dec.Bytes(int(detailMaskSize))
					unk := dec.Bytes(4)
					if unk[0] != 0 || unk[1] != 0 || unk[2] != 0 || unk[3] != 0 {
						unkByte := dec.Uint8()
						if unkByte != 0 {
							unkFloat := dec.Float32()
							unkFloat2 := dec.Float32()
							unkFloat3 := dec.Float32()
							unkFloat4 := dec.Float32()
						}
						unkFloat := dec.Float32()
					}
				}
			}

			overlayCount := dec.Uint32()
			for j := 0; j < int(layerCount); j++ {
				materialName := dec.StringLenPrefixUint32()
				detailMaskDim := dec.Uint32()
				szM := detailMaskDim * detailMaskDim
				for k := 0; k < int(szM); k++ {
					detailMaskByte := dec.Uint8()
				}
				overlayCount++
			}
		}

		singlePlacableCount := dec.Uint32()
		for i := 0; i < int(singlePlacableCount); i++ {
			modelName := dec.StringLenPrefixUint32()
			s := dec.StringLenPrefixUint32()
			longitude := dec.Uint32()
			latitude := dec.Uint32()
			x := dec.Float32()
			y := dec.Float32()
			z := dec.Float32()
			rotX := dec.Float32()
			rotY := dec.Float32()
			rotZ := dec.Float32()
			scaleX := dec.Float32()
			scaleY := dec.Float32()
			scaleZ := dec.Float32()
			unk := dec.Bytes(4)
			if unk[0] != 0 || unk[1] != 0 || unk[2] != 0 || unk[3] != 0 {
				unkByte := dec.Uint32()
			}

			/*if(terrain->GetModels().count(model_name) == 0) {
				EQGModelLoader model_loader;
				std::shared_ptr<EQG::Geometry> m(new EQG::Geometry());
				m->SetName(model_name);
				if (model_loader.Load(archive, model_name + ".mod", m)) {
					terrain->GetModels()[model_name] = m;
				}
				else if (model_loader.Load(archive, model_name, m)) {
					terrain->GetModels()[model_name] = m;
				}
				else {
					m->GetMaterials().clear();
					m->GetPolygons().clear();
					m->GetVertices().clear();
					terrain->GetModels()[model_name] = m;
				}
			}*/

	/*std::shared_ptr<Placeable> p(new Placeable());
	p->SetName(model_name);
	p->SetFileName(model_name);
	p->SetLocation(0.0f, 0.0f, 0.0f);
	p->SetRotation(rot_x, rot_y, rot_z);
	p->SetScale(scale_x, scale_y, scale_z);*/

	//There's a lot of work with offsets here =/
	/*std::shared_ptr<PlaceableGroup> pg(new PlaceableGroup());
	pg->SetFromTOG(false);
	pg->SetLocation(x, y, z);

	float terrain_height = 0.0f;
	float adjusted_x = x;
	float adjusted_y = y;

	if(adjusted_x < 0)
		adjusted_x = adjusted_x + (-(int)(adjusted_x / (terrain->GetOpts().units_per_vert * terrain->GetOpts().quads_per_tile)) + 1) * (terrain->GetOpts().units_per_vert * terrain->GetOpts().quads_per_tile);
	else
		adjusted_x = fmod(adjusted_x, terrain->GetOpts().units_per_vert * terrain->GetOpts().quads_per_tile);

	if(adjusted_y < 0)
		adjusted_y = adjusted_y + (-(int)(adjusted_y / (terrain->GetOpts().units_per_vert * terrain->GetOpts().quads_per_tile)) + 1) * (terrain->GetOpts().units_per_vert * terrain->GetOpts().quads_per_tile);
	else
		adjusted_y = fmod(adjusted_y, terrain->GetOpts().units_per_vert * terrain->GetOpts().quads_per_tile);

	int row_number = (int)(adjusted_y / terrain->GetOpts().units_per_vert);
	int column = (int)(adjusted_x / terrain->GetOpts().units_per_vert);
	int quad = row_number * terrain->GetOpts().quads_per_tile + column;

	float quad_vertex1Z = tile->GetFloats()[quad + row_number];
	float quad_vertex2Z = tile->GetFloats()[quad + row_number + terrain->GetOpts().quads_per_tile + 1];
	float quad_vertex3Z = tile->GetFloats()[quad + row_number + terrain->GetOpts().quads_per_tile + 2];
	float quad_vertex4Z = tile->GetFloats()[quad + row_number + 1];

	glm::vec3 p1(row_number * terrain->GetOpts().units_per_vert, (quad % terrain->GetOpts().quads_per_tile) * terrain->GetOpts().units_per_vert, quad_vertex1Z);
	glm::vec3 p2(p1.x + terrain->GetOpts().units_per_vert, p1.y, quad_vertex2Z);
	glm::vec3 p3(p1.x + terrain->GetOpts().units_per_vert, p1.y + terrain->GetOpts().units_per_vert, quad_vertex3Z);
	glm::vec3 p4(p1.x, p1.y + terrain->GetOpts().units_per_vert, quad_vertex4Z);

	terrain_height = HeightWithinQuad(p1, p2, p3, p4, adjusted_y, adjusted_x);

	pg->SetTileLocation(tile_start_y, tile_start_x, terrain_height);
	pg->SetRotation(0.0f, 0.0f, 0.0f);
	pg->SetScale(1.0f, 1.0f, 1.0f);
	pg->AddPlaceable(p);
	terrain->AddPlaceableGroup(pg);*/
	/*}

	areaCount := dec.Uint32()
	for i := uint32(0); i < areaCount; i++ {
		areaName := dec.StringLenPrefixUint32()
		areaType := dec.Int32()
		areaName2 := dec.StringLenPrefixUint32()

		longitude := dec.Uint32()
		latitude := dec.Uint32()

		x := dec.Float32()
		y := dec.Float32()
		z := dec.Float32()

		rotX := dec.Float32()
		rotY := dec.Float32()
		rotZ := dec.Float32()

		scaleX := dec.Float32()
		scaleY := dec.Float32()
		scaleZ := dec.Float32()

		sizeX := dec.Float32()
		sizeY := dec.Float32()
		sizeZ := dec.Float32()

		// region

		terrainHeight := float32(0)
		adjustedX := x
		adjustedY := y
	}
	/*
		if adjustedX < 0 {
			adjustedX = adjustedX + (-(int32(adjustedX/(terrain.GetOpts().UnitsPerVert*terrain.GetOpts().QuadsPerTile)))+1)*(terrain.GetOpts().UnitsPerVert*terrain.GetOpts().QuadsPerTile)
		} else {
			adjustedX = math.Mod(adjustedX, terrain.GetOpts().UnitsPerVert*terrain.GetOpts().QuadsPerTile)
		}

		if adjustedY < 0 {
			adjustedY = adjustedY + (-(int32(adjustedY/(terrain.GetOpts().UnitsPerVert*terrain.GetOpts().QuadsPerTile)))+1)*(terrain.GetOpts().UnitsPerVert*terrain.GetOpts().QuadsPerTile)
		} else {
			adjustedY = math.Mod(adjustedY, terrain.GetOpts().UnitsPerVert*terrain.GetOpts().QuadsPerTile)
		}

		rowNumber := int32(adjustedY / terrain.GetOpts().UnitsPerVert)
		column := int32(adjustedX / terrain.GetOpts().UnitsPerVert)
		quad := rowNumber*terrain.GetOpts().QuadsPerTile + column

		quadVertex1Z := tile.GetFloats()[quad+rowNumber]
		quadVertex2Z := tile.GetFloats()[quad+rowNumber+terrain.GetOpts().QuadsPerTile+1]
		quadVertex3Z := tile.GetFloats()[quad+rowNumber+terrain.GetOpts().QuadsPerTile+2]
		quadVertex4Z := tile.GetFloats()[quad+rowNumber+1]

			p1 := glm.Vec3{rowNumber * terrain.GetOpts().UnitsPerVert, float32(quad%terrain.GetOpts().QuadsPerTile) * terrain.GetOpts().UnitsPerVert, quadVertex1Z}
			p2 := glm.Vec3{p1.X + terrain.GetOpts().UnitsPerVert, p1.Y, quadVertex2Z}
			p3 := glm.Vec3{p1.X + terrain.GetOpts().UnitsPerVert, p1.Y + terrain.GetOpts().UnitsPerVert, quadVertex3Z}
			p4 := glm.Vec3{p1.X, p1.Y + terrain.GetOpts().UnitsPerVert, quadVertex4Z}

			terrainHeight = HeightWithinQuad(p1, p2, p3, p4, adjustedY, adjustedX)

			region := NewRegion()
			region.SetName(areaName)
			region.SetAlternateName(areaName2)
			region.SetLocation(x+tileStartY, y+tileStartX, z+terrainHeight)
			region.SetRotation(rotX, rotY, rotZ)
			region.SetScale(scaleX, scaleY, scaleZ)
			region.SetExtents(sizeX/2.0, sizeY/2.0, sizeZ/2.0)
			region.SetFlags(areaType, 0)

			terrain.AddRegion(region)
		}

			SafeVarAllocParse(uint32_t, Light_effects_count);
			for (uint32_t j = 0; j < Light_effects_count; ++j) {
				SafeStringAllocParse(s);
				SafeStringAllocParse(s2);

				SafeVarAllocParse(int8_t, unk);

				SafeVarAllocParse(uint32_t, longitude);
				SafeVarAllocParse(uint32_t, latitude);

				SafeVarAllocParse(float, x);
				SafeVarAllocParse(float, y);
				SafeVarAllocParse(float, z);

				SafeVarAllocParse(float, rot_x);
				SafeVarAllocParse(float, rot_y);
				SafeVarAllocParse(float, rot_z);

				SafeVarAllocParse(float, scale_x);
				SafeVarAllocParse(float, scale_y);
				SafeVarAllocParse(float, scale_z);

				SafeVarAllocParse(float, unk_float);
			}

			SafeVarAllocParse(uint32_t, tog_ref_count);
			for (uint32_t j = 0; j < tog_ref_count; ++j) {
				SafeStringAllocParse(tog_name);

				SafeVarAllocParse(uint32_t, longitude);
				SafeVarAllocParse(uint32_t, latitude);

				SafeVarAllocParse(float, x);
				SafeVarAllocParse(float, y);
				SafeVarAllocParse(float, z);

				SafeVarAllocParse(float, rot_x);
				SafeVarAllocParse(float, rot_y);
				SafeVarAllocParse(float, rot_z);

				SafeVarAllocParse(float, scale_x);
				SafeVarAllocParse(float, scale_y);
				SafeVarAllocParse(float, scale_z);

				SafeVarAllocParse(float, z_adjust);

				std::vector<char> tog_buffer;
				if(!archive.Get(tog_name + ".tog", tog_buffer))
				{
					eqLogMessage(LogWarn, "Failed to load tog file %s.tog.", tog_name.c_str());
					continue;
				} else {
					eqLogMessage(LogTrace, "Loaded tog file %s.tog.", tog_name.c_str());
				}

				std::shared_ptr<PlaceableGroup> pg(new PlaceableGroup());
				pg->SetFromTOG(true);
				pg->SetLocation(x, y, z + (scale_z * z_adjust));
				pg->SetRotation(rot_x, rot_y, rot_z);
				pg->SetScale(scale_x, scale_y, scale_z);
				pg->SetTileLocation(tile_start_y, tile_start_x, 0.0f);

				std::vector<std::string> tokens;
				std::shared_ptr<Placeable> p;
				ParseConfigFile(tog_buffer, tokens);
				for (size_t k = 0; k < tokens.size();) {
					auto token = tokens[k];
					if (token.compare("*BEGIN_OBJECT") == 0) {
						p.reset(new Placeable());
						++k;
					}
					else if (token.compare("*NAME") == 0) {
						if (k + 1 >= tokens.size()) {
							break;
						}

						std::string model_name = tokens[k + 1];
						std::transform(model_name.begin(), model_name.end(), model_name.begin(), ::tolower);

						if (terrain->GetModels().count(model_name) == 0) {
							EQGModelLoader model_loader;
							std::shared_ptr<EQG::Geometry> m(new EQG::Geometry());
							m->SetName(model_name);
							if (model_loader.Load(archive, model_name + ".mod", m)) {
								terrain->GetModels()[model_name] = m;
							}
							else if (model_loader.Load(archive, model_name, m)) {
								terrain->GetModels()[model_name] = m;
							}
							else {
								m->GetMaterials().clear();
								m->GetPolygons().clear();
								m->GetVertices().clear();
								terrain->GetModels()[model_name] = m;
							}
						}

						p->SetName(model_name);
						p->SetFileName(model_name);
						k += 2;
					}
					else if (token.compare("*POSITION") == 0) {
						if (k + 3 >= tokens.size()) {
							break;
						}

						p->SetLocation(std::stof(tokens[k + 1]), std::stof(tokens[k + 2]), std::stof(tokens[k + 3]));
						k += 4;
					}
					else if (token.compare("*ROTATION") == 0) {
						if (k + 3 >= tokens.size()) {
							break;
						}

						p->SetRotation(std::stof(tokens[k + 1]), std::stof(tokens[k + 2]), std::stof(tokens[k + 3]));
						k += 4;
					}
					else if (token.compare("*SCALE") == 0) {
						if (k + 1 >= tokens.size()) {
							break;
						}

						p->SetScale(std::stof(tokens[k + 1]), std::stof(tokens[k + 1]), std::stof(tokens[k + 1]));
						k += 2;
					}
					else if (token.compare("*END_OBJECT") == 0) {
						pg->AddPlaceable(p);
						++k;
					}
					else {
						++k;
					}
				}

				terrain->AddPlaceableGroup(pg);
			}

			tile->SetLocation(tile_start_x, tile_start_y);

			}
	*/
	if dec.Error() != nil {
		return fmt.Errorf("decode %s: %w", datName, dec.Error())
	}

	log.Debugf("%s is version %d and has %d objects, %d lights", e.name, e.version, len(e.objectManager.Objects()), len(e.lights))
	return nil
}
