package output

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"
	"text/tabwriter"

	"gopkg.in/yaml.v3"
)

type Printer interface {
	Print(w io.Writer) error
}

// Kong for embedding in kong structs
type Kong struct {
	Format string `help:"Output format (${enum})" enum:"json,yaml,csv,table" default:"table"`
}

type output struct {
	data   interface{}
	fields []string
}

func New(format string, data interface{}, fields []string) (Printer, error) {
	if reflect.TypeOf(data).Kind() != reflect.Slice {
		return nil, fmt.Errorf("data must be a slice")
	}

	f := output{
		data:   data,
		fields: fields,
	}

	switch format {
	case "json":
		return &jsonOutput{f}, nil
	case "yaml":
		return &yamlOutput{f}, nil
	case "csv":
		return &csvOutput{f}, nil
	case "table":
		return &tableOutput{f}, nil
	default:
		return nil, fmt.Errorf("unknown format")
	}
}

type jsonOutput struct {
	output
}

func (j *jsonOutput) Print(w io.Writer) error {
	s := reflect.ValueOf(j.data) // data is a slice

	data := []map[string]interface{}{}
	for i := 0; i < s.Len(); i++ {
		row := make(map[string]interface{})
		for _, f := range j.fields {
			v, err := getFieldValue(s.Index(i).Interface(), f)
			if err != nil {
				return err
			}
			row[f] = v
		}
		data = append(data, row)
	}

	return json.NewEncoder(w).Encode(data)
}

type yamlOutput struct {
	output
}

func (y *yamlOutput) Print(w io.Writer) error {
	s := reflect.ValueOf(y.data) // data is a slice

	data := []map[string]interface{}{}
	for i := 0; i < s.Len(); i++ {
		row := make(map[string]interface{})
		for _, f := range y.fields {
			v, err := getFieldValue(s.Index(i).Interface(), f)
			if err != nil {
				return err
			}
			row[f] = v
		}
		data = append(data, row)
	}

	return yaml.NewEncoder(w).Encode(data)
}

type csvOutput struct {
	output
}

func (c *csvOutput) Print(out io.Writer) error {
	s := reflect.ValueOf(c.data) // data is a slice

	w := csv.NewWriter(out)
	if s.Len() > 0 {
		if err := w.Write(c.fields); err != nil {
			return err
		}
	}

	for i := 0; i < s.Len(); i++ {
		row := make([]string, len(c.fields))
		for j, f := range c.fields {
			v, err := getFieldValue(s.Index(i).Interface(), f)
			if err != nil {
				return err
			}
			row[j] = fmt.Sprintf("%v", v)
		}
		if err := w.Write(row); err != nil {
			return err
		}
	}

	w.Flush()

	return nil
}

type tableOutput struct {
	output
}

func (t *tableOutput) Print(w io.Writer) error {
	s := reflect.ValueOf(t.data) // data is a slice

	tw := tabwriter.NewWriter(w, 0, 8, 1, ' ', tabwriter.Debug)

	if s.Len() > 0 {
		fmt.Fprintf(tw, "%s\n", strings.Join(t.fields, "\t"))
	}

	format := fmt.Sprintf("%s\n", strings.Repeat("%s\t", len(t.fields)))

	for i := 0; i < s.Len(); i++ {
		row := make([]interface{}, len(t.fields))
		for j, f := range t.fields {
			v, err := getFieldValue(s.Index(i).Interface(), f)
			if err != nil {
				return err
			}
			row[j] = fmt.Sprintf("%v", v)
		}
		fmt.Fprintf(tw, format, row...)
	}

	return tw.Flush()
}

func getFieldValue(s interface{}, f string) (interface{}, error) {
	v := reflect.ValueOf(s)

	// if pointer get the underlying element
	for v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		panic("not struct")
	}

	fp := strings.Split(f, ".")

	fv := v.FieldByName(fp[0])

	// if pointer get the underlying element
	if fv.Kind() == reflect.Ptr {
		fv = fv.Elem()
	}

	switch fv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fv.Int(), nil
	case reflect.Bool:
		return fv.Bool(), nil
	case reflect.String:
		return fv.String(), nil
	case reflect.Struct:
		if len(fp) > 1 {
			return getFieldValue(fv.Interface(), strings.Join(fp[1:], "."))
		}

		return "", nil
	default:
		return "", nil
	}
}
