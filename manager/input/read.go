package input

import (
	"fmt"
	"io"
	"os"
	"os/signal"
	"strings"
)

// read reads input from UI.Reader
func (i *UI) read() (string, error) {
	i.once.Do(i.setDefault)

	// sigCh is channel which is watch Interruptted signal (SIGINT)
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt)
	defer signal.Stop(sigCh)

	var resultStr string
	var resultErr error
	doneCh := make(chan struct{})

	go func() {
		defer close(doneCh)

		line, err := i.bReader.ReadString('\n')
		if err != nil && err != io.EOF {
			resultErr = fmt.Errorf("failed to read the input: %s", err)
		}

		resultStr = strings.TrimSuffix(line, LineSep)
		// brute force for the moment
		resultStr = strings.TrimSuffix(line, "\n")

	}()

	select {
	case <-sigCh:
		return "", ErrInterrupted
	case <-doneCh:
		return resultStr, resultErr
	}
}

func (i *UI) rawReadline(f *os.File) (string, error) {
	var resultBuf []byte
	for {
		var buf [1]byte
		n, err := f.Read(buf[:])
		if err != nil && err != io.EOF {
			return "", err
		}

		if n == 0 || buf[0] == '\n' || buf[0] == '\r' {
			break
		}

		if buf[0] == 3 {
			return "", ErrInterrupted
		}

		resultBuf = append(resultBuf, buf[0])
	}

	fmt.Fprintf(i.Writer, "\n")
	return string(resultBuf), nil
}
