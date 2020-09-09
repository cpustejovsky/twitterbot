package main

import "testing"

func TestEmptyError(t *testing.T) {
	want := "no users to send email to."
	got := (*EmptyError).Error(&EmptyError{})

	assertEqualStrings(t, got, want)
}
