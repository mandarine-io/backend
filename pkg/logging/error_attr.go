package logging

import (
	"github.com/mdobak/go-xerrors"
	"log/slog"
	"path/filepath"
)

type stackFrame struct {
	Func   string `json:"func"`
	Source string `json:"source"`
	Line   int    `json:"line"`
}

func ErrorAttr(err error) slog.Attr {
	return slog.Any(
		"error", slog.GroupValue(
			slog.String("msg", err.Error()),
			slog.Any("stack", marshalStack(err)),
		),
	)
}

func ErrorStringAttr(err string) slog.Attr {
	return slog.Any(
		"error", slog.GroupValue(
			slog.String("msg", err),
		),
	)
}

func marshalStack(err error) []stackFrame {
	err = xerrors.WithStackTrace(err, 2)
	trace := xerrors.StackTrace(err)

	if len(trace) == 0 {
		return nil
	}

	frames := trace.Frames()

	s := make([]stackFrame, len(frames))

	for i, v := range frames {
		f := stackFrame{
			Source: filepath.Join(
				filepath.Base(filepath.Dir(v.File)),
				filepath.Base(v.File),
			),
			Func: filepath.Base(v.Function),
			Line: v.Line,
		}

		s[i] = f
	}

	return s
}
