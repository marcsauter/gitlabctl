package auth_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/marcsauter/gitlabctl/internal/auth"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokener(t *testing.T) {
	hostname1 := "localhost"
	authtoken1 := hostname1 + "-secret"

	// temporary file
	f, err := os.CreateTemp("", "auth-test-*")
	require.NoError(t, err)
	require.NoError(t, f.Close())

	testfile := f.Name()

	defer require.NoError(t, os.Remove(testfile))

	t.Run("invalid content", func(t *testing.T) {
		data := []byte("---")
		require.NoError(t, os.WriteFile(testfile, data, 0600))

		_, err = auth.New(testfile)
		require.Error(t, err)
	})

	t.Run("valid content", func(t *testing.T) {
		data := []byte(fmt.Sprintf("{%q:%q}", hostname1, authtoken1))
		require.NoError(t, os.WriteFile(testfile, data, 0600))

		tkner, err := auth.New(testfile)
		require.NoError(t, err)

		token, err := tkner.Get(hostname1)
		assert.NoError(t, err)
		assert.Equal(t, authtoken1, token)
	})

	t.Run("no content", func(t *testing.T) {
		data := []byte("")
		require.NoError(t, os.WriteFile(testfile, data, 0600))

		tkner, err := auth.New(testfile)
		require.NoError(t, err)

		token, err := tkner.Get(hostname1)
		t.Log(token, err)
		assert.Error(t, err)
		assert.Empty(t, token)

		err = tkner.Remove(hostname1)
		assert.NoError(t, err)
	})

	t.Run("CRUD", func(t *testing.T) {
		data := []byte("")
		require.NoError(t, os.WriteFile(testfile, data, 0600))

		tkner, err := auth.New(testfile)
		require.NoError(t, err)

		err = tkner.Set(hostname1, authtoken1)
		assert.NoError(t, err)

		token, err := tkner.Get(hostname1)
		assert.NoError(t, err)
		assert.Equal(t, authtoken1, token)

		// update a token
		authtoken1new := authtoken1 + "-new"

		err = tkner.Set(hostname1, authtoken1new)
		assert.NoError(t, err)

		token, err = tkner.Get(hostname1)
		assert.NoError(t, err)
		assert.Equal(t, authtoken1new, token)

		// add a second token
		hostname2 := "gitlab.example.com"
		authtoken2 := hostname2 + "-secret"

		err = tkner.Set(hostname2, authtoken2)
		assert.NoError(t, err)

		token, err = tkner.Get(hostname1)
		assert.NoError(t, err)
		assert.Equal(t, authtoken1new, token)

		token, err = tkner.Get(hostname2)
		assert.NoError(t, err)
		assert.Equal(t, authtoken2, token)

		// remove a token
		err = tkner.Remove(hostname1)
		assert.NoError(t, err)

		token, err = tkner.Get(hostname1)
		assert.Error(t, err)
		assert.Empty(t, token)

		token, err = tkner.Get(hostname2)
		assert.NoError(t, err)
		assert.Equal(t, authtoken2, token)
	})
}
