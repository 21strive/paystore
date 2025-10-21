package helper

import "reflect"

func FetchColumns(s interface{}) []string {
	var tags []string

	t := reflect.TypeOf(s)

	// If pointer, get the underlying element
	if t.Kind() == reflect.Ptr {
		t = t.Elem()
	}

	// Ensure we're working with a struct
	if t.Kind() != reflect.Struct {
		return tags
	}

	// Iterate through all fields
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		// Get the tag value for the specified tag name
		if tag, ok := field.Tag.Lookup("db"); ok {
			tags = append(tags, tag)
		}

		// Handle nested structs
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		if fieldType.Kind() == reflect.Struct {
			// Recursively get tags from nested struct
			nestedValue := reflect.New(fieldType).Elem().Interface()
			nestedTags := FetchColumns(nestedValue)
			tags = append(tags, nestedTags...)
		}
	}

	return tags
}
