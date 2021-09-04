package twitterbot

import (
	"reflect"
	"testing"
)

func assertEqual(t *testing.T, got, want interface{}) {
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got:\n%v\nwant\n%v\n", got, want)
	}
}
