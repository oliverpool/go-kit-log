package level_test

import (
	"bytes"
	"errors"
	"strings"
	"testing"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/experimental_level"
)

func TestVariousLevels(t *testing.T) {
	for _, testcase := range []struct {
		allowed []string
		want    string
	}{
		{
			level.AllowAll,
			strings.Join([]string{
				`{"level":"debug","this is":"debug log"}`,
				`{"level":"info","this is":"info log"}`,
				`{"level":"warn","this is":"warn log"}`,
				`{"level":"error","this is":"error log"}`,
			}, "\n"),
		},
		{
			level.AllowDebugAndAbove,
			strings.Join([]string{
				`{"level":"debug","this is":"debug log"}`,
				`{"level":"info","this is":"info log"}`,
				`{"level":"warn","this is":"warn log"}`,
				`{"level":"error","this is":"error log"}`,
			}, "\n"),
		},
		{
			level.AllowInfoAndAbove,
			strings.Join([]string{
				`{"level":"info","this is":"info log"}`,
				`{"level":"warn","this is":"warn log"}`,
				`{"level":"error","this is":"error log"}`,
			}, "\n"),
		},
		{
			level.AllowWarnAndAbove,
			strings.Join([]string{
				`{"level":"warn","this is":"warn log"}`,
				`{"level":"error","this is":"error log"}`,
			}, "\n"),
		},
		{
			level.AllowErrorOnly,
			strings.Join([]string{
				`{"level":"error","this is":"error log"}`,
			}, "\n"),
		},
		{
			level.AllowNone,
			``,
		},
	} {
		var buf bytes.Buffer
		logger := level.New(log.NewJSONLogger(&buf), level.Config{Allowed: testcase.allowed})

		level.Debug(logger).Log("this is", "debug log")
		level.Info(logger).Log("this is", "info log")
		level.Warn(logger).Log("this is", "warn log")
		level.Error(logger).Log("this is", "error log")

		if want, have := testcase.want, strings.TrimSpace(buf.String()); want != have {
			t.Errorf("given Allowed=%v: want\n%s\nhave\n%s", testcase.allowed, want, have)
		}
	}
}

func TestErrSquelch(t *testing.T) {
	myError := errors.New("squelched!")
	logger := level.New(log.NewNopLogger(), level.Config{
		Allowed:      level.AllowWarnAndAbove,
		ErrSquelched: myError,
	})

	if want, have := myError, level.Info(logger).Log("foo", "bar"); want != have {
		t.Errorf("want %#+v, have %#+v", want, have)
	}

	if want, have := error(nil), level.Warn(logger).Log("foo", "bar"); want != have {
		t.Errorf("want %#+v, have %#+v", want, have)
	}
}

func TestErrNoLevel(t *testing.T) {
	myError := errors.New("no level specified")

	var buf bytes.Buffer
	logger := level.New(log.NewJSONLogger(&buf), level.Config{
		AllowNoLevel: false,
		ErrNoLevel:   myError,
	})

	if want, have := myError, logger.Log("foo", "bar"); want != have {
		t.Errorf("want %v, have %v", want, have)
	}
	if want, have := ``, strings.TrimSpace(buf.String()); want != have {
		t.Errorf("want %q, have %q", want, have)
	}
}

func TestAllowNoLevel(t *testing.T) {
	var buf bytes.Buffer
	logger := level.New(log.NewJSONLogger(&buf), level.Config{
		AllowNoLevel: true,
		ErrNoLevel:   errors.New("I should never be returned!"),
	})

	if want, have := error(nil), logger.Log("foo", "bar"); want != have {
		t.Errorf("want %v, have %v", want, have)
	}
	if want, have := `{"foo":"bar"}`, strings.TrimSpace(buf.String()); want != have {
		t.Errorf("want %q, have %q", want, have)
	}
}
