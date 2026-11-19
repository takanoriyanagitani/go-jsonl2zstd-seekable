package main

import (
	"bufio"
	"context"
	"errors"
	"log/slog"
	"os"

	kz "github.com/klauspost/compress/zstd"
	js "github.com/takanoriyanagitani/go-jsonl2zstd-seekable"
)

const EncoderLevelDefault kz.EncoderLevel = kz.SpeedFastest

//nolint:gochecknoglobals
var zeoptsDefault []kz.EOption = []kz.EOption{
	kz.WithEncoderLevel(EncoderLevelDefault),
}

//nolint:gochecknoglobals
var zeopts []kz.EOption = zeoptsDefault

func sub(ctx context.Context) error {
	zenc, err := kz.NewWriter(nil, zeopts...)
	if nil != err {
		return err
	}
	defer func() {
		zcerr := zenc.Close()
		if nil != zcerr {
			slog.Error(
				"unexpected error",
				"on", "closing zstd encoder",
				"detail", zcerr.Error(),
			)
		}
	}()

	var ze js.ZstdEncoder = zenc.EncodeAll
	var frames js.Frames = js.ReaderToFrames(os.Stdin)
	var bwtr *bufio.Writer = bufio.NewWriter(os.Stdout)

	ferr := ze.FramesToWriter(
		ctx,
		frames,
		bwtr,
		func(cerr error) {
			slog.Warn(
				"unexpected error",
				"on", "closing seekable zstd writer",
				"detail", cerr.Error(),
			)
		},
	)
	return errors.Join(ferr, bwtr.Flush())
}

func main() {
	err := sub(context.Background())
	if nil != err {
		slog.Error(
			"unexpected error",
			"detail", err.Error(),
		)
	}
}
