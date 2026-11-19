package jsonl2zseek

import (
	"bufio"
	"context"
	"io"
	"iter"

	zs "github.com/SaveTheRbtz/zstd-seekable-format-go/pkg"
)

type Frames iter.Seq[[]byte]

type ZstdEncoder func(src, dst []byte) []byte

func (e ZstdEncoder) EncodeAll(src, dst []byte) []byte {
	return e(src, dst)
}

func (e ZstdEncoder) AsZstdEncoder() zs.ZSTDEncoder { return e }

func (e ZstdEncoder) FramesToWriter(
	ctx context.Context,
	frames Frames,
	iwtr io.Writer,
	onCloseError func(error),
) error {
	swtr, err := zs.NewWriter(
		iwtr,
		e.AsZstdEncoder(),
	)
	if nil != err {
		return err
	}
	defer func() {
		cerr := swtr.Close()
		if nil != cerr {
			onCloseError(cerr)
		}
	}()

	for frame := range frames {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		_, werr := swtr.Write(frame)
		if nil != werr {
			return werr
		}
	}

	return nil
}

const NewLine byte = 0x0a

func ReaderToFrames(rdr io.Reader) Frames {
	return func(yield func([]byte) bool) {
		var scanner *bufio.Scanner = bufio.NewScanner(rdr)
		var buf []byte
		for scanner.Scan() {
			buf = buf[:0]
			var line []byte = scanner.Bytes()
			buf = append(buf, line...)
			buf = append(buf, NewLine)
			if !yield(buf) {
				return
			}
		}
	}
}
