package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

type sometype struct {
	A string
	B int
	C bool
	D *sometype
}

func TestNew(t *testing.T) {

	data := []*sometype{
		{
			A: "a",
			B: 98,
			C: false,
		},
		{
			A: "A",
			B: 66,
			C: true,
		},
	}

	t.Run("no slice", func(t *testing.T) {
		p, err := New("json", data[0], []string{})
		assert.Error(t, err)
		assert.Nil(t, p)
	})

	t.Run("unknown", func(t *testing.T) {
		p, err := New("unknown", data, []string{})
		assert.Error(t, err)
		assert.Nil(t, p)
	})
}

func TestGetValue(t *testing.T) {
	t.Run("no struct", func(t *testing.T) {
		defer func() { recover() }() //nolint:errcheck

		getFieldValue(1, "") //nolint:errcheck

		t.Error("did not panic")
	})
}

func TestJSON(t *testing.T) {
	data := []*sometype{
		{
			A: "a",
			B: 98,
			C: false,
			D: &sometype{
				A: "A",
				B: 66,
				C: true,
			},
		},
	}

	expected := "[{\"A\":\"a\",\"B\":98,\"C\":false,\"D.A\":\"A\"}]\n"

	p, err := New("json", data, []string{"A", "B", "C", "D.A"})
	assert.NoError(t, err)
	assert.Implements(t, (*Printer)(nil), p)
	assert.IsType(t, &jsonOutput{}, p)

	var buf bytes.Buffer

	assert.NoError(t, p.Print(&buf))
	assert.Equal(t, expected, buf.String())
}

func TestYAML(t *testing.T) {
	data := []*sometype{
		{
			A: "a",
			B: 98,
			C: false,
			D: &sometype{
				A: "A",
				B: 66,
				C: true,
			},
		},
	}

	expected := "- A: a\n  B: 98\n  C: false\n  D.A: A"

	p, err := New("yaml", data, []string{"A", "B", "C", "D.A"})
	assert.NoError(t, err)
	assert.Implements(t, (*Printer)(nil), p)
	assert.IsType(t, &yamlOutput{}, p)

	var buf bytes.Buffer

	assert.NoError(t, p.Print(&buf))
	assert.Equal(t, expected, strings.TrimSpace(buf.String()))
}
func TestCSV(t *testing.T) {
	data := []*sometype{
		{
			A: "a",
			B: 98,
			C: false,
			D: &sometype{
				A: "A",
				B: 66,
				C: true,
			},
		},
	}

	expected := "A,B,C,D.A\na,98,false,A"

	p, err := New("csv", data, []string{"A", "B", "C", "D.A"})
	assert.NoError(t, err)
	assert.Implements(t, (*Printer)(nil), p)
	assert.IsType(t, &csvOutput{}, p)

	var buf bytes.Buffer

	assert.NoError(t, p.Print(&buf))
	assert.Equal(t, expected, strings.TrimSpace(buf.String()))
}

func TestTABLE(t *testing.T) {
	data := []*sometype{
		{
			A: "a",
			B: 98,
			C: false,
			D: &sometype{
				A: "A",
				B: 66,
				C: true,
			},
		},
	}

	expected := "A |B  |C     |D.A\na |98 |false |A |"

	p, err := New("table", data, []string{"A", "B", "C", "D.A"})
	assert.NoError(t, err)
	assert.Implements(t, (*Printer)(nil), p)
	assert.IsType(t, &tableOutput{}, p)

	var buf bytes.Buffer

	assert.NoError(t, p.Print(&buf))
	assert.Equal(t, expected, strings.TrimSpace(buf.String()))
}
