package headers

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHeadersParse(t *testing.T) {
	// Test: Valid single header
	headers := NewHeaders()
	data := []byte("Host: localhost:42069\r\n\r\n")
	n, done, err := headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 23, n)
	assert.False(t, done)

	//Test: Valid single header with extra spaces
	headers = NewHeaders()
	data = []byte("Host:    localhost:42069 \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "localhost:42069", headers["host"])
	assert.Equal(t, 27, n)
	assert.False(t, done)

	//Test: Valid 2 header with existing headers
	data = []byte("This: Roland \r\nThat: Angela \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "Roland", headers["this"])
	m, done, err := headers.Parse(data[n:])
	assert.Equal(t, "Angela", headers["that"])
	assert.Equal(t, err, nil)

	assert.False(t, done)
	m, done, err = headers.Parse(data[n+m:])
	assert.True(t, done)

	// Test: Invalid spacing header
	headers = NewHeaders()
	data = []byte("       Host: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	//Test: Invalid field-name characters
	headers = NewHeaders()
	data = []byte("H©st: localhost:42069\r\n\r\n")
	n, done, err = headers.Parse(data)
	require.Error(t, err)
	assert.Equal(t, 0, n)
	assert.False(t, done)

	//Test: Multiple entries for the same header-field
	headers = NewHeaders()
	headers["librarians"] = "Malkuth"
	data = []byte("Librarians: Binah \r\nLibrarians: Yesod \r\n\r\n")
	n, done, err = headers.Parse(data)
	require.NoError(t, err)
	require.NotNil(t, headers)
	assert.Equal(t, "Malkuth, Binah", headers["librarians"])
	m, done, err = headers.Parse(data[n:])
	assert.Equal(t, "Malkuth, Binah, Yesod", headers["librarians"])
	assert.Equal(t, err, nil)

}
