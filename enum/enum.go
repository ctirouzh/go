package enum

import (
	"fmt"
	"reflect"
	"strings"
)

type numeric interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Enum interface {
	numeric
	// EnumUid must be implemented to be an Enum.
	// It must return a unique identifier for the Enum type.
	EnumUid() string
}

// typeName returns the type name of the given Enum type using
//
//   - reflect.TypeOf(e).PkgPath() if descriptive is true
//
//   - fmt.Sprintf and %T verb in case of descriptive is false.
func typeName[E Enum](e E, descriptive bool) string {
	if !descriptive {
		return fmt.Sprintf("%T", e)
	}
	reflectType := reflect.TypeOf(e)
	return fmt.Sprintf("%s.%s", reflectType.PkgPath(), reflectType.Name())
}

type enum[T Enum] struct {
	uid           string
	shortName     string
	fullName      string
	values        []T
	descriptions  map[T]string
	stringToValue map[string]T
}

// Is checks if the given value is one of the target values of the given Enum type parameter.
func Is[E Enum](enum E, target E, or ...E) bool {
	if enum == target {
		return true
	}
	for i := range or {
		if enum == or[i] {
			return true
		}
	}
	return false
}

// Validate checks if the given value is one of the registered values of the given Enum type.
func Validate[E Enum](e E) error {
	found, ok := enumSet[e.EnumUid()]
	if !ok {
		return fmt.Errorf("[Enum] %q %w", typeName(e, true), ErrNotRegisteredYet)
	}
	registered := found.(*enum[E])
	for i := range registered.values {
		if e == registered.values[i] {
			return nil
		}
	}
	return fmt.Errorf("%w for %s: must be one of %v, got %v", ErrInvalidValue, registered.shortName, registered.values, e)
}

func New[E Enum](desc string) (*E, error) {
	var e E
	x, registered := enumSet[e.EnumUid()]
	if !registered {
		panic(fmt.Sprintf("[Enum] %q %v", typeName(e, true), ErrNotRegisteredYet))
	}
	member := x.(*enum[E])
	if v, found := member.stringToValue[desc]; found {
		return &v, nil
	}
	oneof := make([]string, 0, len(member.values))
	for str := range member.stringToValue {
		oneof = append(oneof, str)
	}
	return nil, fmt.Errorf("%w for %s: must be one of %v, got %s", ErrInvalidValue, member.shortName, strings.Join(oneof, ", "), desc)
}
