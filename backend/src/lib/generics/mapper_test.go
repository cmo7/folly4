package generics

import "testing"

type Input struct {
	Name  string
	Age   int
	Email string
}

type Output struct {
	Name string
	Age  int
}

type ExtendedOutput struct {
	Name    string
	Age     int
	Address string
}

type LargeInput struct {
	Field1 string
	Field2 int
	Field3 float64
	Field4 bool
	Field5 string
}

type LargeOutput struct {
	Field1 string
	Field2 int
	Field3 float64
	Field4 bool
}

func TestGenericMapperDefault(t *testing.T) {
	mapper := NewGenericMapperDefault[Input, Output]()

	input := Input{Name: "John Doe", Age: 30, Email: "john@example.com"}
	output := mapper.Map(input)

	if output.Name != input.Name || output.Age != input.Age {
		t.Errorf("Expected output.Name: %s, output.Age: %d, got Name: %s, Age: %d", input.Name, input.Age, output.Name, output.Age)
	}
}

func TestGenericMapperExcludingFields(t *testing.T) {
	mapper := NewGenericMapperExcluding[Input, Output]([]string{"Name"})

	input := Input{Name: "John Doe", Age: 30, Email: "john@example.com"}
	output := mapper.Map(input)

	if output.Name != "" || output.Age != input.Age {
		t.Errorf("Expected output.Name to be empty, got: %s. Expected Age: %d, got: %d", output.Name, input.Age, output.Age)
	}
}

func TestGenericMapperIncludingFields(t *testing.T) {
	mapper := NewGenericMapperIncluding[Input, Output]([]string{"Age"})

	input := Input{Name: "John Doe", Age: 30, Email: "john@example.com"}
	output := mapper.Map(input)

	if output.Name != "" || output.Age != input.Age {
		t.Errorf("Expected output.Name to be empty, got: %s. Expected Age: %d, got: %d", output.Name, input.Age, output.Age)
	}
}

func TestGenericMapperInvalidFields(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected a panic for invalid field name, but code did not panic")
		}
	}()

	_ = NewGenericMapperIncluding[Input, Output]([]string{"InvalidField"})
}

func TestGenericMapperExcludingInvalidFields(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected a panic for invalid field name, but code did not panic")
		}
	}()

	_ = NewGenericMapperExcluding[Input, Output]([]string{"InvalidField"})
}

func TestEmptyInputStruct(t *testing.T) {
	mapper := NewGenericMapperDefault[Input, Output]()

	input := Input{}
	output := mapper.Map(input)

	if output.Name != "" || output.Age != 0 {
		t.Errorf("Expected default output values, got Name: %s, Age: %d", output.Name, output.Age)
	}
}

func TestPartialInputFields(t *testing.T) {
	mapper := NewGenericMapperDefault[Input, Output]()

	input := Input{Name: "Partial Name"}
	output := mapper.Map(input)

	if output.Name != input.Name || output.Age != 0 {
		t.Errorf("Expected output.Name: %s, output.Age: 0, got Name: %s, Age: %d", input.Name, output.Name, output.Age)
	}
}

func TestFieldCaseSensitivity(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected a panic for case-sensitive field mismatch, but code did not panic")
		}
	}()

	_ = NewGenericMapperIncluding[Input, Output]([]string{"name"}) // lowercase "name" does not match "Name"
}

func TestMappingToStructWithExtraFields(t *testing.T) {
	mapper := NewGenericMapperDefault[Input, ExtendedOutput]()

	input := Input{Name: "John Doe", Age: 30, Email: "john@example.com"}
	output := mapper.Map(input)

	if output.Name != input.Name || output.Age != input.Age || output.Address != "" {
		t.Errorf("Expected output.Name: %s, output.Age: %d, output.Address to be empty, got Name: %s, Age: %d, Address: %s",
			input.Name, input.Age, output.Name, output.Age, output.Address)
	}
}

func TestMappingFromStructWithExtraFields(t *testing.T) {
	mapper := NewGenericMapperDefault[Input, Output]()

	input := Input{Name: "John Doe", Age: 30, Email: "john@example.com"}
	output := mapper.Map(input)

	if output.Name != input.Name || output.Age != input.Age {
		t.Errorf("Expected output.Name: %s, output.Age: %d, got Name: %s, Age: %d", input.Name, input.Age, output.Name, output.Age)
	}
}

func TestExcludingAllFields(t *testing.T) {
	mapper := NewGenericMapperExcluding[Input, Output]([]string{"Name", "Age"})

	input := Input{Name: "John Doe", Age: 30, Email: "john@example.com"}
	output := mapper.Map(input)

	if output.Name != "" || output.Age != 0 {
		t.Errorf("Expected all fields excluded, got Name: %s, Age: %d", output.Name, output.Age)
	}
}

func TestIncludingNonOverlappingFields(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected a panic for non-overlapping included fields, but code did not panic")
		}
	}()

	_ = NewGenericMapperIncluding[Input, Output]([]string{"Email"}) // "Email" is not in Output
}

func TestLargeStruct(t *testing.T) {
	mapper := NewGenericMapperDefault[LargeInput, LargeOutput]()

	input := LargeInput{
		Field1: "Test",
		Field2: 123,
		Field3: 45.67,
		Field4: true,
		Field5: "Ignored",
	}
	output := mapper.Map(input)

	if output.Field1 != input.Field1 || output.Field2 != input.Field2 || output.Field3 != input.Field3 || output.Field4 != input.Field4 {
		t.Errorf("Mapping failed for large struct. Expected Field1: %s, Field2: %d, Field3: %f, Field4: %t, got Field1: %s, Field2: %d, Field3: %f, Field4: %t",
			input.Field1, input.Field2, input.Field3, input.Field4, output.Field1, output.Field2, output.Field3, output.Field4)
	}
}

func TestNilInput(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Expected a panic for nil input, but code did not panic")
		}
	}()

	var input *Input
	mapper := NewGenericMapperDefault[*Input, Output]()
	_ = mapper.Map(input)
}
