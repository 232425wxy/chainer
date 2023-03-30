package clogging

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestLoggerNameRegexp(t *testing.T) {
	var tests = []struct{
		str string
		ok bool
	}{
		{str: "logger", ok: true},
		{str: "_", ok: true},
		{str: "logger.name", ok: true},
		{str: "logger_#-:._", ok: true},
		{str: "_._", ok: true},
		{str: "_.:", ok: true},
		{str: ":", ok: true},
		{str: "._", ok: false},
		{str: "a._.", ok: false},
		{str: "a.:.", ok: false},
	}

	for _, test := range tests {
		t.Run(test.str, func(t *testing.T) {
			require.Equal(t, test.ok, loggerNameRegexp.MatchString(test.str))
		})
	}
}
