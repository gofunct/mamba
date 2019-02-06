package input

import (
	"bytes"
	"fmt"
	"strconv"
)

func (i *UI) Select(q *Query) (string, error) {
	// Set default val
	i.once.Do(i.setDefault)
	q.Opts.Required = true

	// Find default index which opts.Default indicates
	defaultIndex := -1
	defaultVal := q.Opts.Default
	if defaultVal != "" {
		for i, item := range q.Opts.Options {
			if item == defaultVal {
				defaultIndex = i
			}
		}

		// DefaultVal is set but doesn't exist in q.Opts.Options
		if defaultIndex == -1 {
			// This error message is not for user
			// Should be found while development
			return "", fmt.Errorf("opt.Default is specified but item does not exist in q.Opts.Options")
		}
	}

	// Construct the query & display it to user
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("%s\n\n", q.Q))
	for i, item := range q.Opts.Options {
		buf.WriteString(fmt.Sprintf("%d. %s\n", i+1, item))
	}

	buf.WriteString("\n")
	fmt.Fprintf(i.Writer, buf.String())

	// resultStr and resultErr are return val of this function
	var resultStr string
	var resultErr error
	for {

		// Construct the asking line to input
		var buf bytes.Buffer
		buf.WriteString("Enter a number")

		if defaultIndex >= 0 {
			buf.WriteString(fmt.Sprintf(" (Default is %d)", defaultIndex+1))
		}

		buf.WriteString(": ")
		fmt.Fprintf(i.Writer, buf.String())

		// Read user input from reader.
		line, err := i.read()
		if err != nil {
			resultErr = err
			break
		}

		// line is empty but default is provided returns it
		if line == "" && defaultIndex >= 0 {
			resultStr = q.Opts.Options[defaultIndex]
			break
		}

		if line == "" && q.Opts.Required {
			if !q.Opts.Loop {
				resultErr = ErrEmpty
				break
			}

			fmt.Fprintf(i.Writer, "Input must not be empty. Answer by a number.\n\n")
			continue
		}

		// Convert user input string to int val
		n, err := strconv.Atoi(line)
		if err != nil {
			if !q.Opts.Loop {
				resultErr = ErrNotNumber
				break
			}

			fmt.Fprintf(i.Writer,
				"%q is not a valid input. Answer by a number.\n\n", line)
			continue
		}

		// Check answer is in range of q.Opts.Options
		if n < 1 || len(q.Opts.Options) < n {
			if !q.Opts.Loop {
				resultErr = ErrOutOfRange
				break
			}

			fmt.Fprintf(i.Writer,
				"%q is not a valid choice. Choose a number from 1 to %d.\n\n",
				line, len(q.Opts.Options))
			continue
		}

		// validate input by custom function
		validate := bind(q)
		if err := validate(line); err != nil {
			if !q.Opts.Loop {
				resultErr = err
				break
			}

			fmt.Fprintf(i.Writer, "Failed to bind tag to provided input string: %s\n\n", err)
			continue
		}

		// Reach here means it gets ideal input.
		resultStr = q.Opts.Options[n-1]
		break
	}

	// Insert the new line for next output
	fmt.Fprintf(i.Writer, "\n")

	return resultStr, resultErr
}
