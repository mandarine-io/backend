package logging

import (
	"github.com/mdobak/go-xerrors"
	"github.com/rs/zerolog"
	"path/filepath"
)

type stackFrame struct {
	Func   string `json:"func"`
	Source string `json:"source"`
	Line   int    `json:"line"`
}

type stackFrames []stackFrame

func marshalStack(err error) interface{} {
	err = xerrors.WithStackTrace(err, 3)
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

	return (*stackFrames)(&s)
}

func (s *stackFrames) MarshalZerologObject(e *zerolog.Event) {
	e.Array("stack", s)
}

func (s *stackFrames) MarshalZerologArray(a *zerolog.Array) {
	for _, v := range *s {
		a.Object(v)
	}
}

func (s stackFrame) MarshalZerologObject(e *zerolog.Event) {
	e.Str("func", s.Func).
		Str("source", s.Source).
		Int("line", s.Line)
}
