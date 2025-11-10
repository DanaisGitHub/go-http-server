package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRequestLineParse1(t *testing.T) {

	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("HOst: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers.Get("hOsT"))
	assert.Equal(t, 23, n)
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host : localhost:42069       \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)
}
func TestErronous(t *testing.T) {

	// Test: Multi-field values
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n Foo: Bar \r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["Host"])
	assert.Equal(t, "Bar", headers["Foo"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	// // Test: Empty
	// headers = NewHeaders()
	// data = []byte("\r\n\r\n       Host: localhost:42069       \r\n\r\n")
	// n, done, err = headers.Parse(data)
	// require.NoError(t, err)
	// assert.NotEqual(t, 0, n)
	// assert.False(t, done)

	// // Test: Malformed
	// headers = NewHeaders()
	// data = []byte("Host     \r\n\r\n")
	// n, done, err = headers.Parse(data)
	// require.Error(t, err)
	// assert.Equal(t, 0, n)
	// assert.False(t, done)

}
