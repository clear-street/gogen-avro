package imprt_test

import (
	"testing"

	"github.com/clear-street/gogen-avro/imprt"
	"github.com/stretchr/testify/require"
)

func TestImprt(t *testing.T) {
	require.Equal(t, "1/2/3/x", imprt.Path("1/2/3/x", "a"))
	require.Equal(t, "x/b", imprt.Path("x", "a.b"))
	require.Equal(t, "x/b/c", imprt.Path("x", "a.b.c"))

	require.Equal(t, "x", imprt.Pkg("1/2/3/x", "a"))
	require.Equal(t, "b", imprt.Pkg("x", "a.b"))
	require.Equal(t, "c", imprt.Pkg("x", "a.b.c"))
	require.Equal(t, "e", imprt.Pkg("1/2/x", "a.b.c.c.e"))

	require.Equal(t, "x.Type", imprt.Type("1/2/x", "a", "Type"))
	require.Equal(t, "*x.Type", imprt.Type("x", "a", "*Type"))

	require.Equal(t, "XFoo", imprt.UniqName("x", "a", "Foo"))
	require.Equal(t, "XFoo", imprt.UniqName("1/2/x", "a", "Foo"))
	require.Equal(t, "CFoo", imprt.UniqName("x", "a.b.c", "Foo"))

	require.True(t, imprt.IsRootPkg("1/2/3/34/x", "a"))
	require.False(t, imprt.IsRootPkg("1/2/3/34/x", "a.b"))
}
