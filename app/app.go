package app

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"time"
)

const (
	codeOK = iota
	codeInitAppErr
)

func Main(args []string) int {
	a, err := newApp(args, os.Stdout, os.Stderr)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
		return codeInitAppErr
	}

	return a.run()
}

type App struct {
	content  string
	reverse  bool
	frames   int
	interval time.Duration

	outStream io.Writer
	errStream io.Writer
}

func newApp(args []string, outStream, errStream io.Writer) (*App, error) {
	flags := flag.NewFlagSet("gogo", flag.ContinueOnError)
	flags.SetOutput(errStream)
	flags.Usage = func() {
		_, _ = fmt.Fprintf(errStream, helpMessage)
		flags.PrintDefaults()
	}

	var (
		file    string
		reverse bool
		frames  int
		fps     int
	)
	flags.StringVar(&file, "file", "", "filepath whose content will run. if not set, the content is got from standard input")
	flags.BoolVar(&reverse, "reverse", false, "if set, the animation run from right to left")
	flags.IntVar(&frames, "frames", 50, "number of frames of the animation (default: 50, min: 1, max: 200)")
	flags.IntVar(&fps, "fps", 10, "fps of the animation (default: 10, min: 1, max: 60)")
	if err := flags.Parse(args[1:]); err != nil {
		return nil, fmt.Errorf("failed to parse command line options: %s", strings.Join(args[1:], " "))
	}

	content, err := getContent(file)
	if err != nil {
		return nil, err
	}

	if frames < 1 {
		return nil, fmt.Errorf("min value of frames is 1. got=%d", frames)
	} else if frames > 200 {
		return nil, fmt.Errorf("max value of frames is 200. got=%d", frames)
	}

	if fps < 1 {
		return nil, fmt.Errorf("min value of fps is 1. got=%d", fps)
	} else if fps > 60 {
		return nil, fmt.Errorf("max value of fps is 60. got=%d", fps)
	}
	interval := time.Duration(1000/fps) * time.Millisecond

	return &App{
		content:  content,
		reverse:  reverse,
		frames:   frames,
		interval: interval,

		outStream: outStream,
		errStream: errStream,
	}, nil
}

func (a *App) run() int {
	for t := 0; t < a.frames; t++ {
		a.clear()

		var spaces string
		if a.reverse {
			spaces = strings.Repeat(" ", a.frames-1-t)
		} else {
			spaces = strings.Repeat(" ", t)
		}

		rdr := strings.NewReader(a.content)
		sc := bufio.NewScanner(rdr)
		for sc.Scan() {
			_, _ = fmt.Fprintln(a.outStream, spaces+sc.Text())
		}

		time.Sleep(a.interval)
	}

	return codeOK
}

func (a *App) clear() {
	_, _ = fmt.Fprint(a.outStream, "\033[H\033[2J")
}

func getContent(file string) (string, error) {
	var (
		in  io.Reader
		err error
	)
	if file == "" {
		in = os.Stdin
	} else {
		in, err = os.Open(file)
		if err != nil {
			return "", err
		}
	}

	b, err := ioutil.ReadAll(in)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

const helpMessage = `
gogo is a command line tool to run some input in console.

EXAMPLE: 
$ gogo -contest ABC051 -problem C -command 'python c.py'

OPTION:
`
