package log_test

import (
	"bytes"
	"flag"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zapcore"
	"kdex.dev/crds/log"
	crzap "sigs.k8s.io/controller-runtime/pkg/log/zap"
)

func TestNew(t *testing.T) {
	type message struct {
		level         int
		msg           string
		keysAndValues []interface{}
	}
	tests := []struct {
		name         string
		category     string
		defaultLevel string
		levelMap     map[string]string
		message      message
		want         string
		wantErr      bool
	}{
		{
			name:         "info nil 0 foo",
			category:     "foo",
			defaultLevel: "info",
			levelMap:     map[string]string{},
			message: message{
				level:         0,
				msg:           "test",
				keysAndValues: nil,
			},
			want: "INFO\tfoo\ttest\n",
		},
		{
			name:         "info info 0 foo",
			category:     "foo",
			defaultLevel: "info",
			levelMap: map[string]string{
				"foo": "info",
			},
			message: message{
				level:         0,
				msg:           "test",
				keysAndValues: nil,
			},
			want: "INFO\tfoo\ttest\n",
		},
		{
			name:         "error debug 1 foo",
			category:     "foo",
			defaultLevel: "error",
			levelMap: map[string]string{
				"foo": "debug",
			},
			message: message{
				level:         1,
				msg:           "test",
				keysAndValues: nil,
			},
			want: "DEBUG\tfoo\ttest\n",
		},
		{
			name:         "error nil 0 foo",
			category:     "foo",
			defaultLevel: "error",
			levelMap:     map[string]string{},
			message: message{
				level:         0,
				msg:           "test",
				keysAndValues: nil,
			},
			want: "",
		},
		{
			name:         "error 2 1 foo",
			category:     "foo",
			defaultLevel: "error",
			levelMap: map[string]string{
				"foo": "2",
			},
			message: message{
				level:         1,
				msg:           "test",
				keysAndValues: nil,
			},
			want: "DEBUG\tfoo\ttest\n",
		},
		{
			name:         "error 2 0 foo",
			category:     "foo",
			defaultLevel: "error",
			levelMap: map[string]string{
				"foo": "2",
			},
			message: message{
				level:         0,
				msg:           "test",
				keysAndValues: nil,
			},
			want: "INFO\tfoo\ttest\n",
		},
		{
			name:         "error 2 0 other",
			category:     "other",
			defaultLevel: "error",
			levelMap: map[string]string{
				"foo": "2",
			},
			message: message{
				level:         0,
				msg:           "test",
				keysAndValues: nil,
			},
			want: "",
		},
		{
			name:         "invalid default level int",
			defaultLevel: "0",
			levelMap:     map[string]string{},
			wantErr:      true,
		},
		{
			name:         "invalid default level string",
			defaultLevel: "foo",
			levelMap:     map[string]string{},
			wantErr:      true,
		},
		{
			name:         "invalid mapped level int",
			defaultLevel: "error",
			levelMap: map[string]string{
				"foo": "0",
			},
			wantErr: true,
		},
		{
			name:         "invalid mapped level string",
			defaultLevel: "error",
			levelMap: map[string]string{
				"foo": "foo",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			buffer := bytes.Buffer{}

			flags := flag.NewFlagSet("dummy-flags", flag.ContinueOnError)
			opts := crzap.Options{
				DestWriter:  &buffer,
				TimeEncoder: func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {},
			}
			opts.BindFlags(flags)
			simulatedArgs := []string{
				fmt.Sprintf("--zap-log-level=%s", tt.defaultLevel),
				"--zap-encoder=console",
				"--zap-stacktrace-level=error",
			}
			err := flags.Parse(simulatedArgs)
			if err != nil {
				if !tt.wantErr {
					t.Fatal(err)
				}
				return
			}

			logger, err := log.New(&opts, tt.levelMap)
			if err != nil {
				if !tt.wantErr {
					t.Fatal(err)
				}
				return
			}

			logger = logger.WithName(tt.category)

			if tt.message.keysAndValues == nil {
				logger.V(tt.message.level).Info(tt.message.msg)
			} else {
				logger.V(tt.message.level).Info(tt.message.msg, tt.message.keysAndValues)
			}

			assert.Equal(t, tt.want, buffer.String())
		})
	}
}
