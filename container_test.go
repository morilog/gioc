package gioc_test

import (
	"testing"

	"github.com/morilog/gioc"
	"github.com/stretchr/testify/require"
)

type exampleType string

func Test_Bind(t *testing.T) {
	gioc.Bind(func() (exampleType, error) {
		return exampleType("hello"), nil
	}, false)

	var x exampleType
	err := gioc.Make(&x)
	require.Nil(t, err)
	require.Equal(t, x, exampleType("hello"))

	var y string
	err = gioc.Make(&y)
	require.NotNil(t, err)
}

type someCounter int

func Test_Singleton(t *testing.T) {
	resolver := func() (someCounter, error) {
		return someCounter(1) + 1, nil
	}

	gioc.Singleton(resolver)

	var x someCounter
	err := gioc.Make(&x)
	require.Nil(t, err)
	require.Equal(t, 2, int(x))

	var y someCounter
	err = gioc.Make(&y)
	require.Nil(t, err)
	require.Equal(t, 2, int(y))
}
