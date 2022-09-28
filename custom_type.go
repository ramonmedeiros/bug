package bug

import (
	"context"
	"database/sql/driver"
	"fmt"
	"io"

	"entgo.io/ent/schema/field"
)

type File struct {
	File string
}

var _ field.ValueScanner = new(File)

func (f *File) Scan(i any) error {
	v, ok := i.(string)
	if !ok {
		return nil
	}

	*f = File{
		File: v,
	}
	return nil
}

func (f File) Value() (driver.Value, error) {
	return f.File, nil
}

func (f *File) UnmarshalGQLContext(ctx context.Context, i any) error {
	return f.Scan(i)
}

func (f File) MarshalGQLContext(ctx context.Context, w io.Writer) error {
	w.Write([]byte(fmt.Sprintf(`"%s"`, f)))
	return nil
}
