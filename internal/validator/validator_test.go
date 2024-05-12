package validator

import (
	"fmt"
	"regexp"
	"testing"
)

// Unit Tests
func TestValidator_Valid(t *testing.T) {
	v := New()
	if !v.Valid() {
		t.Error("Expected no errors, got errors")
	}

	v.AddError("key", "message")
	if v.Valid() {
		t.Error("Expected errors, got none")
	}
}

func TestValidator_AddError(t *testing.T) {
	v := New()
	v.AddError("key", "message")
	if len(v.Errors) != 1 || v.Errors["key"] != "message" {
		t.Errorf("Expected errors map to contain one entry with key 'key', got: %v", v.Errors)
	}
}

func TestValidator_Check(t *testing.T) {
	v := New()
	v.Check(false, "key", "message")
	if len(v.Errors) != 1 || v.Errors["key"] != "message" {
		t.Errorf("Expected errors map to contain one entry with key 'key', got: %v", v.Errors)
	}

	v.Check(true, "key2", "message2")
	if len(v.Errors) != 1 {
		t.Errorf("Expected errors map to contain one entry, got: %v", v.Errors)
	}
}

func TestPermittedValue(t *testing.T) {
	ok := PermittedValue(3, 1, 2, 3)
	if !ok {
		t.Error("Expected 3 to be in the list")
	}

	ok = PermittedValue(4, 1, 2, 3)
	if ok {
		t.Error("Expected 4 not to be in the list")
	}
}

func TestMatches(t *testing.T) {
	rx := regexp.MustCompile("^\\d{3}$")
	ok := Matches("123", rx)
	if !ok {
		t.Error("Expected 123 to match pattern")
	}

	ok = Matches("abc", rx)
	if ok {
		t.Error("Expected abc not to match pattern")
	}
}

func TestUnique(t *testing.T) {
	ok := Unique([]int{1, 2, 3})
	if !ok {
		t.Error("Expected unique values")
	}

	ok = Unique([]int{1, 2, 3, 1})
	if ok {
		t.Error("Expected non-unique values")
	}
}

// Integration Tests
func TestEmailValidationIntegration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		email    string
		expected bool
	}{
		{"john@example.com", true},
		{"invalidemail@", false},
		{"missingat.com", false},
		{"@missinglocalpart.com", false},
		{"missingdomain@", false},
	}

	for _, test := range tests {
		t.Run(test.email, func(t *testing.T) {
			valid := Matches(test.email, EmailRX)
			if valid != test.expected {
				t.Errorf("Expected %v for %s, got %v", test.expected, test.email, valid)
			}
		})
	}
}

func TestUniqueIntegration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		values   []int
		expected bool
	}{
		{[]int{1, 2, 3}, true},
		{[]int{1, 2, 3, 1}, false},
		{[]int{1, 1, 1, 1}, false},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%v", test.values), func(t *testing.T) {
			valid := Unique(test.values)
			if valid != test.expected {
				t.Errorf("Expected %v for %v, got %v", test.expected, test.values, valid)
			}
		})
	}
}

func TestPermittedValueIntegration(t *testing.T) {
	t.Parallel()
	tests := []struct {
		value    int
		list     []int
		expected bool
	}{
		{3, []int{1, 2, 3}, true},
		{4, []int{1, 2, 3}, false},
		{5, []int{5, 6, 7}, true},
	}

	for _, test := range tests {
		t.Run(fmt.Sprintf("%d in %v", test.value, test.list), func(t *testing.T) {
			valid := PermittedValue(test.value, test.list...)
			if valid != test.expected {
				t.Errorf("Expected %v for %d in %v, got %v", test.expected, test.value, test.list, valid)
			}
		})
	}
}

func TestMultipleRulesIntegration(t *testing.T) {
	t.Parallel()
	v := New()

	// Checking multiple conditions
	email := "test@example.com"
	v.Check(Matches(email, EmailRX), "email", "Invalid email format")
	v.Check(Unique([]string{"a", "b", "c"}), "unique", "Values must be unique")

	// Validating the results
	if !v.Valid() {
		t.Errorf("Expected validation to pass, but got errors: %v", v.Errors)
	}

	// Adding a duplicate email to fail the test
	v.AddError("email", "Email already exists")

	if v.Valid() {
		t.Errorf("Expected validation to fail, but got no errors")
	}
}
