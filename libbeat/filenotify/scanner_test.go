package filenotify

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGlobbing(t *testing.T) {
	c := Config{
		Paths:         []string{"/logs/**/*.log"},
		RecursiveGlob: true,
	}

	scanner, err := New(c, nil, nil)
	if err != nil {
		t.Fatal(err)
	}

	expectedPaths := []string{
		"/logs/*.log",
		"/logs/*/*.log",
		"/logs/*/*/*.log",
		"/logs/*/*/*/*.log",
		"/logs/*/*/*/*/*.log",
		"/logs/*/*/*/*/*/*.log",
		"/logs/*/*/*/*/*/*/*.log",
		"/logs/*/*/*/*/*/*/*/*.log",
		"/logs/*/*/*/*/*/*/*/*/*.log",
	}
	assert.True(t, reflect.DeepEqual(scanner.paths, expectedPaths))
}
