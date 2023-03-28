package cenc

import (
	"runtime"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRuntimeFuncForPC(t *testing.T) {
	pc, file, _, _ := runtime.Caller(0)

	testFileIdx := strings.LastIndex(file, "/")
	testFile := file[testFileIdx+1:]

	f := runtime.FuncForPC(pc)
	fNameIdx := strings.LastIndex(f.Name(), ".")
	fName := f.Name()[fNameIdx+1:]

	require.Equal(t, "runtime_test.go", testFile)
	require.Equal(t, "TestRuntimeFuncForPC", fName)
}
