package input

import (
	"bufio"
	"github.com/gofunct/mamba/runtime/logging"
	"github.com/pkg/errors"
	"io"
	"os"
	"sync"
)

var (
	// defaultWriter and defaultReader is default val for UI.Writer
	// and UI.Reader.
	defaultWriter = os.Stdout
	defaultReader = os.Stdin
)

var (
	// Errs are error returned by input functions.
	// It's useful for handling error from outside of input functions.
	ErrEmpty       = errors.New("default value is not provided but input is empty")
	ErrNotNumber   = errors.New("input must be number")
	ErrOutOfRange  = errors.New("input is out of range")
	ErrInterrupted = errors.New("interrupted")
)

type UI struct {
	Queries []*Query
	Writer  io.Writer
	Reader  io.Reader
	bReader *bufio.Reader
	once    sync.Once
}

// DefaultUI returns default UI. It outputs to stdout and intputs from stdin.
func DefaultUI() *UI {
	return &UI{
		Writer: os.Stdout,
		Reader: os.Stdin,
	}
}

// setDefault sets the default value for UI struct.
func (i *UI) setDefault() {
	// Set the default writer & reader if not provided
	if i.Writer == nil {
		i.Writer = defaultWriter
	}

	if i.Reader == nil {
		i.Reader = defaultReader
	}

	if i.bReader == nil {
		i.bReader = bufio.NewReader(i.Reader)
	}
}

func (u *UI) Enquire() {
	var err error
	var ans string
	for _, v := range u.Queries {
		if v.Q == "" || v.Opts.Key == "" {
			logging.L.Fatalf("%s", "failed to query- query question and key are required")
		}
		if v.Opts.Options != nil {
			ans, err = u.Select(v)
		} else {
			ans, err = u.Ask(v)
		}

		if err != nil {
			logging.L.Fatal("Failed to ask for input", ans, errors.WithStack(err))
		}
	}
}

func (u *UI) AddQueries(q ...*Query) {
	for _, v := range q {
		u.Queries = append(u.Queries, v)
	}
}

func (u *UI) SingleQuery(q string, opts *Options) *Query {
	return &Query{Q: q, Opts: opts}
}
