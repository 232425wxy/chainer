package cenc

import (
	"bytes"
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

func TestParseFormat(t *testing.T) {
	timeLayout := time.RFC850
	var tests = []struct {
		desc       string
		spec       string
		formatters []Formatter
		err        string
	}{
		{
			desc:       "empty",
			spec:       "",
			formatters: nil,
			err:        "",
		},
		{
			desc: "time color",
			spec: fmt.Sprintf("%%{time:%s}%%{color:bold}", timeLayout),
			formatters: []Formatter{
				TimeFormatter{Layout: timeLayout},
				ColorFormatter{Bold: true},
			},
			err: "",
		},
		{
			desc:       "unknown color",
			spec:       "%{color:unknown}",
			formatters: nil,
			err:        "invalid color option: unknown",
		},
		{
			desc: "message string module string",
			spec: "%{message:4s},xxx%{module:6s}:haha",
			formatters: []Formatter{
				MessageFormatter{FormatVerb: "%4s"},
				StringFormatter{Value: ",xxx"},
				ModuleFormatter{FormatVerb: "%6s"},
				StringFormatter{Value: ":haha"},
			},
			err: "",
		},
		{
			desc: "level id shortfunc",
			spec: "%{level:s}%{id:.4f}%{shortfunc:v}",
			formatters: []Formatter{
				LevelFormatter{FormatVerb: "%s"},
				SequenceFormatter{FormatVerb: "%.4f"},
				ShortFuncFormatter{FormatVerb: "%v"},
			},
			err: "",
		},
	}

	for _, test := range tests {
		t.Run(test.desc, func(t *testing.T) {
			formatters, err := ParseFormat(test.spec)
			if test.err == "" {
				require.NoError(t, err)
				for i, formatter := range formatters {
					require.Equal(t, formatter, test.formatters[i])
				}
			} else {
				require.Equal(t, test.err, err.Error())
			}
		})
	}
}

func TestColorFormatter(t *testing.T) {
	entry := zapcore.Entry{
		Level:   zapcore.ErrorLevel,
		Message: "pos or pbft?",
	}

	spec := "%{color:bold}%{color:reset}%{color}"
	formatters, err := ParseFormat(spec)
	require.NoError(t, err)

	pool := buffer.NewPool()

	for _, formatter := range formatters {
		buf := pool.Get()
		formatter.Format(buf, entry, nil)
		fmt.Fprint(buf, entry.Message)
		fmt.Println(buf.String())
		buf.Free()
	}
	// Output:
	// 红色加粗 (pos or pbft?)
	// 黑色 (pos or pbft?)
	// 红色 (pos or pbft?)
}

func TestSequenceFormatter(t *testing.T) {
	mutex := &sync.Mutex{}
	results := map[string]struct{}{}

	ready := &sync.WaitGroup{}
	ready.Add(100)

	finished := &sync.WaitGroup{}
	finished.Add(100)

	SetSequence(0)
	for i := 1; i <= 100; i++ {
		go func(i int) {
			buf := &bytes.Buffer{}
			entry := zapcore.Entry{Level: zapcore.DebugLevel}
			f := SequenceFormatter{FormatVerb: "%d"}
			ready.Done() // setup complete
			ready.Wait() // wait for all go routines to be ready

			f.Format(buf, entry, nil) // format concurrently

			mutex.Lock()
			results[buf.String()] = struct{}{}
			mutex.Unlock()

			finished.Done()
		}(i)
	}

	finished.Wait()
	for i := 1; i <= 100; i++ {
		require.Contains(t, results, strconv.Itoa(i))
	}

	t.Log(results)
}
