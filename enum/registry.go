package enum

import (
	"log"
)

var enumSet = make(map[string]any, 0)

// Register registers the given values as the possible values of the given enum type
// or panics if no values are provided.
//
// Important:
//
// - Do not use this function synchronously/in runtime due to race conditions.
//
// - Use it once in the init() function of the package where the enum type is defined.
func Register[E Enum](descriptions map[E]string) {
	if len(descriptions) == 0 {
		log.Panic("[Enum] no values provided to register.")
	}
	var e E
	uid := e.EnumUid()
	// check if the enum type is already registered
	if _, alreadyRegistered := enumSet[uid]; alreadyRegistered {
		// TODO: check if we want to panic here
		log.Printf("[Enum] %q is already registered with uid %q.", typeName(e, true), uid)
		return
	}
	values := make([]E, 0, len(descriptions))
	stringToValue := make(map[string]E, len(descriptions))
	for key, value := range descriptions {
		values = append(values, key)
		stringToValue[value] = key
	}
	// register the enum type as a member of the enumSet
	member := &enum[E]{
		uid:           uid,
		shortName:     typeName(e, false),
		fullName:      typeName(e, true),
		values:        values,
		descriptions:  descriptions,
		stringToValue: stringToValue,
	}
	enumSet[uid] = member
	log.Printf("[Enum] Successfully registered %q as %v.", member.fullName, member.descriptions)
}

// List returns all the registered values of the given enum type.
// It returns nil if the enum type is not registered.
func List[E Enum]() []E {
	var e E
	values, ok := enumSet[e.EnumUid()]
	if !ok {
		return nil
	}
	return values.(*enum[E]).values
}

func Strings[E Enum]() []string {
	var e E
	x, registered := enumSet[e.EnumUid()]
	if !registered {
		// TODO: check if we want to panic here or return an error
		return nil
	}
	member := x.(*enum[E])
	result := make([]string, len(member.values))
	for i := range member.values {
		result[i] = member.descriptions[member.values[i]]
	}
	return result
}
