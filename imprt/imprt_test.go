package imprt_test

import (
	"testing"

	"github.com/clear-street/gogen-avro/imprt"
	"github.com/stretchr/testify/require"
)

func TestImprt(t *testing.T) {
	require.Equal(t, "x", imprt.Path("x", "a"))
	require.Equal(t, "x/b", imprt.Path("x", "a.b"))
	require.Equal(t, "x/b/c", imprt.Path("x", "a.b.c"))

	require.Equal(t, "x", imprt.Pkg("x", "a"))
	require.Equal(t, "b", imprt.Pkg("x", "a.b"))
	require.Equal(t, "c", imprt.Pkg("x", "a.b.c"))
	require.Equal(t, "e", imprt.Pkg("x", "a.b.c.c.e"))
}
