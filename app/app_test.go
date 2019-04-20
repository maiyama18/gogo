package app

import (
	"bytes"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name string

		inputArgs []string

		expectedContent  string
		expectedReverse  bool
		expectedFrames   int
		expectedInterval time.Duration

		expectedErrMsg       string
		expectedErrStreamMsg string
	}{
		{
			name:      "success-no_option",
			inputArgs: strings.Fields("gogo -file testdata/hello.txt"),

			expectedContent:  "hello",
			expectedReverse:  false,
			expectedFrames:   defaultFrames,
			expectedInterval: time.Duration(1000/defaultFps) * time.Millisecond,
		},
		{
			name:      "success-full_options",
			inputArgs: strings.Fields("gogo -file testdata/hello.txt -fps 20 -frames 100 -reverse"),

			expectedContent:  "hello",
			expectedReverse:  true,
			expectedFrames:   100,
			expectedInterval: time.Duration(1000/20) * time.Millisecond,
		},
		{
			name:      "failure-undefined_option",
			inputArgs: strings.Fields("gogo -undefined undefined"),

			expectedErrMsg:       "failed to parse",
			expectedErrStreamMsg: "gogo is a command line tool",
		},
		{
			name:      "failure-file_not_exist",
			inputArgs: strings.Fields("gogo -file nonexistent.txt"),

			expectedErrMsg: "no such file or directory",
		},
		{
			name:      "failure-fps_too_small",
			inputArgs: strings.Fields("gogo -fps 0 -file testdata/hello.txt"),

			expectedErrMsg:       "min value of fps",
			expectedErrStreamMsg: "gogo is a command line tool",
		},
		{
			name:      "failure-fps_too_large",
			inputArgs: strings.Fields("gogo -fps 1000 -file testdata/hello.txt"),

			expectedErrMsg:       "max value of fps",
			expectedErrStreamMsg: "gogo is a command line tool",
		},
		{
			name:      "failure-frames_too_small",
			inputArgs: strings.Fields("gogo -frames 0 -file testdata/hello.txt"),

			expectedErrMsg:       "min value of frames",
			expectedErrStreamMsg: "gogo is a command line tool",
		},
		{
			name:      "failure-frames_too_large",
			inputArgs: strings.Fields("gogo -frames 1000 -file testdata/hello.txt"),

			expectedErrMsg:       "max value of frames",
			expectedErrStreamMsg: "gogo is a command line tool",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var outStream, errStream bytes.Buffer

			a, err := New(test.inputArgs, &outStream, &errStream)
			if test.expectedErrMsg == "" {
				if err != nil {
					t.Fatalf("got error: %s", err.Error())
				}

				if a.content != test.expectedContent {
					t.Fatalf("content wrong.\nwant\n%s\ngot\n%s", test.expectedContent, a.content)
				}
				if a.reverse != test.expectedReverse {
					t.Fatalf("reverse wrong. want=%t, got=%t", test.expectedReverse, a.reverse)
				}
				if a.frames != test.expectedFrames {
					t.Fatalf("frames wrong. want=%d, got=%d", test.expectedFrames, a.frames)
				}
				if a.interval != test.expectedInterval {
					t.Fatalf("interval wrong. want=%d, got=%d", test.expectedInterval, a.interval)
				}
			} else {
				if !strings.Contains(err.Error(), test.expectedErrMsg) {
					t.Fatalf("error message wrong. expect %q to contain %q", test.expectedErrMsg, err.Error())
				}
				if !strings.Contains(errStream.String(), test.expectedErrStreamMsg) {
					t.Fatalf("error stream message wrong. expect %q to contain %q", test.expectedErrStreamMsg, errStream.String())
				}
			}
		})
	}
}
