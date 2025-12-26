package convert

import (
	"context"
	"io"
	"time"
)

func Convert(ctx context.Context, pdf io.Reader, out io.Writer) error {
	select {
	case <-time.After(5 * time.Second):
		_, err := out.Write([]byte("DOCX DATA"))
		return err
	case <-ctx.Done():
		return ctx.Err()
	}
}
