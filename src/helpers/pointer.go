package helpers

import (
	"strconv"
	"time"

	"github.com/jackc/pgtype"
)

func PointerString(s string) *string {
	return &s
}
func PointerInt(i int) *int {
	return &i
}

func PointerUint(i uint) *uint {
	return &i
}

func PtrInt64(i int64) *int64 {
	return &i
}

func PtrFloat32(f float32) *float32 {
	return &f
}

func PtrFloat64(f float64) *float64 {
	return &f
}

func PtrBool(b bool) *bool {
	return &b
}

func PtrTime(t time.Time) *time.Time {
	return &t
}

func StringPtr(s *string) string {
	if s == nil {
		return ""
	}
	return *s
}

func IntPtr(i *int) int {
	if i == nil {
		return 0
	}
	return *i
}

func UintPtr(i *uint) uint {
	if i == nil {
		return 0
	}
	return *i
}

func Int64Ptr(i *int64) int64 {
	if i == nil {
		return 0
	}
	return *i
}

func Float32Ptr(f *float32) float32 {
	if f == nil {
		return 0
	}
	return *f
}

func Float64Ptr(f *float64) float64 {
	if f == nil {
		return 0
	}
	return *f
}

func BoolPtr(b *bool) bool {
	if b == nil {
		return false
	}
	return *b
}

func TimePtr(t *time.Time) time.Time {
	if t == nil {
		return time.Time{}
	}
	return *t
}

func StringToUint(s string) (*uint, error) {
	i, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return nil, err
	}
	return PointerUint(uint(i)), err
}

func BytesPtr(b []byte) *[]byte {
	return &b
}
func PtrBytes(b *[]byte) []byte {
	if b == nil {
		return []byte{}
	}
	return *b
}

func PGtypeJSONPtr(b pgtype.JSON) *pgtype.JSON {
	return &b
}

func PtrPGTypeJSON(b *pgtype.JSON) pgtype.JSON {
	return *b
}
