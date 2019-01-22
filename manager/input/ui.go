package input

import (
	"bufio"
	"fmt"
	"github.com/gofunct/mamba/logging"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io"
	"sync"
)

type UI struct {
	Queries   []*Query
	ChainFunc ChainFunc
	// Writer is where output is written. For example a query
	// to the user will be written here. By default, it's os.Stdout.
	Writer io.Writer

	// Reader is source of input. By default, it's os.Stdin.
	Reader io.Reader

	// mask is option for read function
	mask    bool
	maskVal string

	bReader *bufio.Reader

	once    sync.Once
	answers map[string]string
}


type Query struct {
	Q    string
	Tag  string
	Opts *Options
}

func NewQuery(q string, tag string, opts *Options) *Query {
	return &Query{Q: q, Tag: tag, Opts: opts}
}

func (u *UI) Query() {
	if u.answers == nil {
		u.answers = make(map[string]string)
	}
	for _, v := range u.Queries {

		if u.ChainFunc != nil {
			v.Opts.ValidateFunc = u.ChainFunc(v.Opts.ValidateFunc)
		}

		ans, err := u.Ask(fmt.Sprintf("%s", v.Q), v.Opts)

		if err != nil {
			logging.L.Fatal("Failed to ask for input", ans, errors.WithStack(err))
		}
		u.answers[v.Tag] = ans
	}
	viper.SetDefault("query.answers", u.answers)
}

func (u *UI) AddQueries(q ...*Query) {
	for _, v := range q {
		u.Queries = append(u.Queries, v)
	}
}

func (u *UI) Enquire(q, tag string) string {
	ans, err := u.Ask(fmt.Sprintf("%s", q), &Options{
		Name: tag,
		ValidateFunc: u.notEmpty(),
		Required: true,
		Loop: true,
	})
	if err != nil {
		logging.L.Fatalln(err)
	}
	return ans

}
