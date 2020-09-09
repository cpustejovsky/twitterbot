package bot

import (
	"testing"
)

func TestEmptyError(t *testing.T) {
	t.Run("Ensures EmptyError is of error type", func(t *testing.T) {

		var _ error = &EmptyError{}

	})

	t.Run("returns same string value when implemented", func(t *testing.T) {

		want := "no users to send email to."
		got := (*EmptyError).Error(&EmptyError{})

		assertEqualStrings(t, got, want)
	})
}
