package wce

import (
	"database/sql"
	"fmt"
	"strconv"
)

type NullInt8 struct {
	Int8  int8
	Valid bool
}
type NullInt16 sql.NullInt16
type NullInt32 sql.NullInt32
type NullInt64 sql.NullInt64
type NullUint8 struct {
	Uint8 uint8
	Valid bool
}
type NullUint16 struct {
	Uint16 uint16
	Valid  bool
}
type NullUint32 struct {
	Uint32 uint32
	Valid  bool
}
type NullUint64 struct {
	Uint64 uint64
	Valid  bool
}

type NullFloat32 struct {
	Float32 float32
	Valid   bool
}
type NullFloat64 sql.NullFloat64
type NullString sql.NullString
type NullFloat32Slice3 struct {
	Float32Slice3 [3]float32
	Valid         bool
}
type NullFloat32Slice4 struct {
	Float32Slice4 [4]float32
	Valid         bool
}
type NullFloat32Slice6 struct {
	Float32Slice6 [6]float32
	Valid         bool
}

type NullInt16Slice3 struct {
	Int16Slice3 [3]int16
	Valid       bool
}

func wcVal(inVal interface{}) string {
	switch val := inVal.(type) {
	case int:
		return strconv.Itoa(val)
	case int8:
		return strconv.Itoa(int(val))
	case int16:
		return strconv.Itoa(int(val))
	case int32:
		return strconv.Itoa(int(val))
	case int64:
		return fmt.Sprintf("%d", val)
	case NullInt8:
		if !val.Valid {
			return "NULL"
		}
		return strconv.Itoa(int(val.Int8))
	case NullInt16:
		if !val.Valid {
			return "NULL"
		}
		return strconv.Itoa(int(val.Int16))
	case NullInt32:
		if !val.Valid {
			return "NULL"
		}
		return strconv.Itoa(int(val.Int32))
	case NullInt64:
		if !val.Valid {
			return "NULL"
		}
		return fmt.Sprintf("%d", val.Int64)
	case NullUint8:
		if !val.Valid {
			return "NULL"
		}
		return strconv.Itoa(int(val.Uint8))
	case NullUint16:
		if !val.Valid {
			return "NULL"
		}
		return strconv.Itoa(int(val.Uint16))
	case NullUint32:
		if !val.Valid {
			return "NULL"
		}
		return strconv.Itoa(int(val.Uint32))
	case NullUint64:
		if !val.Valid {
			return "NULL"
		}
		return fmt.Sprintf("%d", val.Uint64)
	case NullFloat32:
		if !val.Valid {
			return "NULL"
		}
		return fmt.Sprintf("%0.8e", val.Float32)
	case NullFloat64:
		if !val.Valid {
			return "NULL"
		}
		return fmt.Sprintf("%0.8e", val.Float64)
	case NullString:
		if !val.Valid {
			return "NULL"
		}
		return val.String
	case NullFloat32Slice3:
		if !val.Valid {
			return "NULL NULL NULL"
		}
		return fmt.Sprintf("%0.8e %0.8e %0.8e", val.Float32Slice3[0], val.Float32Slice3[1], val.Float32Slice3[2])
	case NullFloat32Slice4:
		if !val.Valid {
			return "NULL NULL NULL NULL"
		}
		return fmt.Sprintf("%0.8e %0.8e %0.8e %0.8e", val.Float32Slice4[0], val.Float32Slice4[1], val.Float32Slice4[2], val.Float32Slice4[3])
	case NullFloat32Slice6:
		if !val.Valid {
			return "NULL NULL NULL NULL NULL NULL"
		}
		return fmt.Sprintf("%0.8e %0.8e %0.8e %0.8e %0.8e %0.8e", val.Float32Slice6[0], val.Float32Slice6[1], val.Float32Slice6[2], val.Float32Slice6[3], val.Float32Slice6[4], val.Float32Slice6[5])
	case NullInt16Slice3:
		if !val.Valid {
			return "NULL NULL NULL NULL"
		}
		return fmt.Sprintf("%d %d %d", val.Int16Slice3[0], val.Int16Slice3[1], val.Int16Slice3[2])
	default:
		return fmt.Sprintf("INVALID_%v", inVal)
	}
}

func parse(inVal interface{}, src ...string) error {
	if len(src) < 1 {
		return fmt.Errorf("need at least 1 argument: %v", src)
	}
	switch val := inVal.(type) {
	case *int:
		i, err := strconv.Atoi(src[0])
		if err != nil {
			return err
		}
		*val = i
		return nil
	case *int8:
		i, err := strconv.ParseInt(src[0], 10, 8)
		if err != nil {
			return err
		}
		*val = int8(i)
		return nil
	case *int16:
		i, err := strconv.ParseInt(src[0], 10, 16)
		if err != nil {
			return err
		}
		*val = int16(i)
		return nil
	case *int32:
		i, err := strconv.ParseInt(src[0], 10, 32)
		if err != nil {
			return err
		}
		*val = int32(i)
		return nil
	case *int64:
		i, err := strconv.ParseInt(src[0], 10, 64)
		if err != nil {
			return err
		}
		*val = i
		return nil
	case *uint8:
		i, err := strconv.ParseUint(src[0], 10, 8)
		if err != nil {
			return err
		}
		*val = uint8(i)
		return nil
	case *uint16:
		i, err := strconv.ParseUint(src[0], 10, 16)
		if err != nil {
			return err
		}
		*val = uint16(i)
		return nil
	case *uint32:
		i, err := strconv.ParseUint(src[0], 10, 32)
		if err != nil {
			return err
		}
		*val = uint32(i)
		return nil
	case *uint64:
		i, err := strconv.ParseUint(src[0], 10, 64)
		if err != nil {
			return err
		}
		*val = i
		return nil
	case *float32:
		f, err := strconv.ParseFloat(src[0], 32)
		if err != nil {
			return err
		}
		*val = float32(f)
		return nil
	case *float64:
		f, err := strconv.ParseFloat(src[0], 64)
		if err != nil {
			return err
		}
		*val = f
		return nil

	case *NullInt8:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		i, err := strconv.ParseInt(src[0], 10, 8)
		if err != nil {
			return err
		}
		val.Int8 = int8(i)
		val.Valid = true
		return nil
	case *NullInt16:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		i, err := strconv.ParseInt(src[0], 10, 16)
		if err != nil {
			return err
		}
		val.Int16 = int16(i)
		val.Valid = true
		return nil
	case *NullInt32:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		i, err := strconv.ParseInt(src[0], 10, 32)
		if err != nil {
			return err
		}
		val.Int32 = int32(i)
		val.Valid = true
		return nil
	case *NullInt64:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		i, err := strconv.ParseInt(src[0], 10, 64)
		if err != nil {
			return err
		}
		val.Int64 = i
		val.Valid = true
		return nil
	case *NullUint8:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		i, err := strconv.ParseUint(src[0], 10, 8)
		if err != nil {
			return err
		}

		val.Uint8 = uint8(i)
		val.Valid = true
		return nil
	case *NullUint16:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		i, err := strconv.ParseUint(src[0], 10, 16)
		if err != nil {
			return err
		}
		val.Uint16 = uint16(i)
		val.Valid = true
		return nil
	case *NullUint32:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		i, err := strconv.ParseUint(src[0], 10, 32)
		if err != nil {
			return err
		}
		val.Uint32 = uint32(i)
		val.Valid = true
		return nil
	case *NullUint64:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		i, err := strconv.ParseUint(src[0], 10, 64)
		if err != nil {
			return err
		}
		val.Uint64 = i
		val.Valid = true
		return nil
	case *NullFloat32:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		f, err := strconv.ParseFloat(src[0], 32)
		if err != nil {
			return err
		}
		val.Float32 = float32(f)
		val.Valid = true
		return nil
	case *NullFloat64:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		f, err := strconv.ParseFloat(src[0], 64)
		if err != nil {
			return err
		}
		val.Float64 = f
		val.Valid = true
		return nil
	case *NullString:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		val.String = src[0]
		val.Valid = true
		return nil
	case *NullFloat32Slice3:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}

		for i := 0; i < 3; i++ {
			v, err := strconv.ParseFloat(src[i], 32)
			if err != nil {
				return err
			}
			val.Float32Slice3[i] = float32(v)
		}
		val.Valid = true
		return nil
	case *NullFloat32Slice4:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}

		for i := 0; i < 4; i++ {
			v, err := strconv.ParseFloat(src[i], 32)
			if err != nil {
				return err
			}
			val.Float32Slice4[i] = float32(v)
		}

		val.Valid = true
		return nil
	case *NullFloat32Slice6:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		if len(src) < 6 {
			return fmt.Errorf("need 6 arguments: %v", src)
		}

		for i := 0; i < 6; i++ {
			v, err := strconv.ParseFloat(src[i], 32)
			if err != nil {
				return err
			}
			val.Float32Slice6[i] = float32(v)
		}
		val.Valid = true

		return nil
	case NullInt16Slice3:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		if len(src) < 3 {
			return fmt.Errorf("need 3 arguments: %v", src)
		}

		for i := 0; i < 3; i++ {
			v, err := strconv.ParseInt(src[i], 10, 16)
			if err != nil {
				return err
			}
			val.Int16Slice3[i] = int16(v)
		}
		val.Valid = true
		return nil
	case *[3]int16:
		if len(src) < 3 {
			return fmt.Errorf("need 3 arguments: %v", src)
		}
		for i := 0; i < 3; i++ {
			v, err := strconv.ParseInt(src[i], 10, 16)
			if err != nil {
				return err
			}
			val[i] = int16(v)
		}
		return nil

	case *[4]uint8:
		if len(src) < 4 {
			return fmt.Errorf("need 4 arguments: %v", src)
		}
		for i := 0; i < 4; i++ {
			v, err := strconv.ParseUint(src[i], 10, 8)
			if err != nil {
				return err
			}
			val[i] = uint8(v)
		}
		return nil
	case *[6]float32:
		if len(src) < 6 {
			return fmt.Errorf("need 6 arguments: %v", src)
		}
		for i := 0; i < 6; i++ {
			v, err := strconv.ParseFloat(src[i], 32)

			if err != nil {
				return err
			}
			val[i] = float32(v)
		}
		return nil

	case *[4]float32:
		if len(src) < 4 {
			return fmt.Errorf("need 4 arguments: %v", src)
		}
		for i := 0; i < 4; i++ {

			v, err := strconv.ParseFloat(src[i], 32)
			if err != nil {
				return err
			}
			val[i] = float32(v)
		}
		return nil

	case *[3]float32:
		if len(src) < 3 {
			return fmt.Errorf("need 3 arguments: %v", src)
		}
		for i := 0; i < 3; i++ {
			v, err := strconv.ParseFloat(src[i], 32)
			if err != nil {
				return err
			}
			val[i] = float32(v)
		}
		return nil
	case *[2]float32:
		if len(src) < 2 {
			return fmt.Errorf("need 2 arguments: %v", src)
		}
		for i := 0; i < 2; i++ {
			v, err := strconv.ParseFloat(src[i], 32)
			if err != nil {
				return err
			}
			val[i] = float32(v)
		}
		return nil
	case *[3]uint16:
		if len(src) < 3 {
			return fmt.Errorf("need 3 arguments: %v", src)
		}
		for i := 0; i < 3; i++ {
			v, err := strconv.ParseUint(src[i], 10, 16)
			if err != nil {
				return err
			}
			val[i] = uint16(v)
		}
		return nil
	case *NullInt16Slice3:
		if src[0] == "NULL" {
			val.Valid = false
			return nil
		}
		if len(src) < 3 {
			return fmt.Errorf("need 3 arguments: %v", src)
		}

		for i := 0; i < 3; i++ {
			v, err := strconv.ParseInt(src[i], 10, 16)
			if err != nil {
				return err
			}
			val.Int16Slice3[i] = int16(v)
		}
		val.Valid = true
		return nil

	default:
		return fmt.Errorf("unknown type: %T", inVal)
	}
}
