package input

import (
	"bytes"
	"fmt"
)

func (i *UI) Ask(q *Query) (string, error) {
	i.once.Do(i.setDefault)

	// Display the query to the user.
	fmt.Fprintf(i.Writer, "%s", q.Q)

	// resultStr and resultErr are return val of this function
	var resultStr string
	var resultErr error

	loopCount := 0
	for {
		loopCount++

		// Construct the instruction to user.
		var buf bytes.Buffer

		buf.WriteString("\nEnter a value")
		defaultVal := q.Opts.Default
		buf.WriteString(fmt.Sprintf(" (Default is %s)", defaultVal))
		buf.WriteString(": ")
		fmt.Fprintf(i.Writer, buf.String())

		// Read user input from UI.Reader.
		line, err := i.read()
		if err != nil {
			resultErr = err
			break
		}

		// line is empty but default is provided returns it
		if line == "" && q.Opts.Default != "" {
			resultStr = q.Opts.Default
			break
		}

		if line == "" && q.Opts.Required {
			if !q.Opts.Loop {
				resultErr = ErrEmpty
				break
			}

			fmt.Fprintf(i.Writer, "Input must not be empty.\n\n")
			continue
		}

		// validate input by custom fuction
		validate := bind(q)
		if err := validate(line); err != nil {
			if !q.Opts.Loop {
				resultErr = err
				break
			}

			fmt.Fprintf(i.Writer, "Failed to validate input string: %s\n\n", err)
			continue
		}

		// Reach here means it gets ideal input.
		resultStr = line
		break
	}

	// Insert the new line for next output
	fmt.Fprintf(i.Writer, "\n")

	return resultStr, resultErr
}

func (i *UI) AskStringVar(q *Query, v *string) error {
	i.once.Do(i.setDefault)

	// Display the query to the user.
	fmt.Fprintf(i.Writer, "%s", q.Q)

	// resultStr and resultErr are return val of this function
	var resultStr string
	var resultErr error

	loopCount := 0
	for {
		loopCount++

		// Construct the instruction to user.
		var buf bytes.Buffer

		buf.WriteString("\nEnter a value")
		defaultVal := q.Opts.Default
		buf.WriteString(fmt.Sprintf(" (Default is %s)", defaultVal))
		buf.WriteString(": ")
		fmt.Fprintf(i.Writer, buf.String())

		// Read user input from UI.Reader.
		line, err := i.read()
		if err != nil {
			resultErr = err
			break
		}

		// line is empty but default is provided returns it
		if line == "" && q.Opts.Default != "" {
			resultStr = q.Opts.Default
			break
		}

		if line == "" && q.Opts.Required {
			if !q.Opts.Loop {
				resultErr = ErrEmpty
				break
			}

			fmt.Fprintf(i.Writer, "Input must not be empty.\n\n")
			continue
		}

		// validate input by custom fuction
		validate := bind(q)
		if err := validate(line); err != nil {
			if !q.Opts.Loop {
				resultErr = err
				break
			}

			fmt.Fprintf(i.Writer, "Failed to validate input string: %s\n\n", err)
			continue
		}

		// Reach here means it gets ideal input.
		resultStr = line
		break
	}

	fmt.Fprintf(i.Writer, "\n")

	v = &resultStr

	return resultErr
}
