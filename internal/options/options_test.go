package options_test

import (
	"testing"

	"github.com/marcsauter/gitlabctl/internal/options"
	"github.com/stretchr/testify/assert"
)

func TestNew(t *testing.T) {
	type source struct {
		Int     int
		IntP    int
		Bool    bool
		BoolP   bool
		String  string
		StringP string

		NotSet string

		WrongType int

		NotInTarget bool
	}
	src := source{}

	type target struct {
		String  string
		StringP *string
		Bool    bool
		BoolP   *bool
		Int     int
		IntP    *int

		NotSet *string

		WrongType string
	}
	tgt := target{}

	t.Run("not pointer", func(t *testing.T) {
		assert.Error(t, options.Transfer(src, &tgt))
		assert.Error(t, options.Transfer(&src, tgt))
	})

	t.Run("not struct", func(t *testing.T) {
		s1 := ""
		s2 := ""
		assert.Error(t, options.Transfer(&s1, &tgt))
		assert.Error(t, options.Transfer(&src, &s2))
	})

	t.Run("set", func(t *testing.T) {
		src := &source{
			Bool:        true,
			BoolP:       true,
			String:      "123",
			StringP:     "456",
			Int:         123,
			IntP:        456,
			WrongType:   789,
			NotInTarget: true,
		}

		tgt := &target{}

		err := options.Transfer(src, tgt)
		assert.NoError(t, err)
		assert.Equal(t, src.Bool, tgt.Bool)
		assert.Equal(t, src.BoolP, *tgt.BoolP)
		assert.Equal(t, src.String, tgt.String)
		assert.Equal(t, src.StringP, *tgt.StringP)
		assert.Equal(t, src.Int, tgt.Int)
		assert.Equal(t, src.IntP, *tgt.IntP)
		assert.Nil(t, tgt.NotSet)
		assert.NotEqual(t, src.WrongType, tgt.WrongType)
	})
}
