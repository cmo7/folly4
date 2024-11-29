package generics

import (
	"reflect"
	"slices"
)

type MapperFunc[I, O interface{}] func(input I) O

// Mapper is an interface that defines a method to map an input object to an output object.
// The fields of the input object are copied to the output object.
type Mapper[I, O interface{}] interface {
	Map(input I) O
}

// GenericMapperImpl is a struct that implements the Mapper interface.
// It has two fields: excludedFields and includedFields.
// excludedFields is a list of fields that should not be copied from the input object to the output object.
// includedFields is a list of fields that should be copied from the input object to the output object.
// If both excludedFields and includedFields are empty, all fields are copied.
type GenericMapperImpl[I, O interface{}] struct {
	excludedFields []string
	includedFields []string
}

func NewGenericMapper[I, O interface{}](excludedFields []string, includedFields []string) Mapper[I, O] {
	return &GenericMapperImpl[I, O]{excludedFields: excludedFields, includedFields: includedFields}
}

func NewGenericMapperExcluding[I, O interface{}](excludedFields []string) Mapper[I, O] {
	// Check if excludedFields are properties of the input type and output type

	var helperInputInstance I
	var helperOutputInstance O

	inputFields := reflect.VisibleFields(reflect.TypeOf(helperInputInstance))
	outputFields := reflect.VisibleFields(reflect.TypeOf(helperOutputInstance))

	inputFieldsNames := make([]string, len(inputFields))
	for i, field := range inputFields {
		inputFieldsNames[i] = field.Name
	}

	outputFieldsNames := make([]string, len(outputFields))
	for i, field := range outputFields {
		outputFieldsNames[i] = field.Name
	}

	for _, excludedField := range excludedFields {
		if !slices.Contains(inputFieldsNames, excludedField) {
			panic("Excluded field " + excludedField + " is not a property of the input type")
		}

		if !slices.Contains(outputFieldsNames, excludedField) {
			panic("Excluded field " + excludedField + " is not a property of the output type")
		}
	}

	// Not excluded fields must be of the same type
	for _, inputField := range inputFields {
		if !slices.Contains(excludedFields, inputField.Name) {
			for _, outputField := range outputFields {
				if inputField.Name == outputField.Name && inputField.Type != outputField.Type {
					panic("Field " + inputField.Name + " is not of the same type in the input and output types")
				}
			}
		}
	}

	return &GenericMapperImpl[I, O]{excludedFields: excludedFields}
}

func NewGenericMapperIncluding[I, O interface{}](includedFields []string) Mapper[I, O] {
	// Check if includedFields are properties of the input type and output type

	var helperInputInstance I
	var helperOutputInstance O

	inputFields := reflect.VisibleFields(reflect.TypeOf(helperInputInstance))
	outputFields := reflect.VisibleFields(reflect.TypeOf(helperOutputInstance))

	inputFieldsNames := getFieldNames(inputFields)

	outputFieldsNames := getFieldNames(outputFields)

	for _, includedField := range includedFields {
		if !slices.Contains(inputFieldsNames, includedField) {
			panic("Included field " + includedField + " is not a property of the input type")
		}

		if !slices.Contains(outputFieldsNames, includedField) {
			panic("Included field " + includedField + " is not a property of the output type")
		}
	}

	// Included fields must be of the same type
	for _, includedField := range includedFields {
		for _, inputField := range inputFields {
			if includedField == inputField.Name {
				for _, outputField := range outputFields {
					if includedField == outputField.Name && inputField.Type != outputField.Type {
						panic("Field " + includedField + " is not of the same type in the input and output types")
					}
				}
			}
		}
	}

	return &GenericMapperImpl[I, O]{includedFields: includedFields}
}

func NewGenericMapperDefault[I, O interface{}]() Mapper[I, O] {

	// Every field pressent in the output type must be present in the input type and must be of the same type
	var helperInputInstance I
	var helperOutputInstance O

	inputFields := reflect.VisibleFields(reflect.TypeOf(helperInputInstance))
	outputFields := reflect.VisibleFields(reflect.TypeOf(helperOutputInstance))

	// All fields shared between input and output types must be of the same type
	for _, inputField := range inputFields {
		for _, outputField := range outputFields {
			if inputField.Name == outputField.Name && inputField.Type != outputField.Type {
				panic("Field " + inputField.Name + " is not of the same type in the input and output types")
			}
		}
	}

	return &GenericMapperImpl[I, O]{
		excludedFields: []string{},
		includedFields: []string{},
	}
}

func (m GenericMapperImpl[I, O]) Map(input I) O {
	var output O
	inputInstance := reflect.ValueOf(input)
	outputInstance := reflect.ValueOf(&output).Elem()

	inputFields := reflect.VisibleFields(inputInstance.Type())
	outputFields := reflect.VisibleFields(outputInstance.Type())

	outputFieldsNames := getFieldNames(outputFields)

	for i, inputField := range inputFields {
		// If there are excluded fields and the current field is in the excluded fields, skip it
		if m.excludedFields != nil && len(m.excludedFields) > 0 && slices.Contains(m.excludedFields, inputField.Name) {
			continue
		}
		// If there are included fields and the current field is not in the included fields, skip it
		if m.includedFields != nil && len(m.includedFields) > 0 && !slices.Contains(m.includedFields, inputField.Name) {
			continue
		}
		// If the current field is not in the output fields, skip it
		if !slices.Contains(outputFieldsNames, inputField.Name) {
			continue
		}
		// Copy the field from the input object to the output object
		outputField := outputInstance.FieldByName(inputField.Name)
		outputField.Set(inputInstance.Field(i))
	}

	return output
}

func getFieldNames(fields []reflect.StructField) []string {
	fieldsNames := make([]string, len(fields))
	for i, field := range fields {
		fieldsNames[i] = field.Name
	}
	return fieldsNames
}
