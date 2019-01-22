package input

import (
	"fmt"
	"github.com/gofunct/mamba/manager/logging"
	"github.com/pkg/errors"
)

func validateTF(q string) ValidateFunc {
	return func(ans string) error {
		logging.L.Debug("received response")
		logging.L.Log("question", q, "answer", ans)
		if ans != "true" && ans != "false" {
			return fmt.Errorf("input must be true or false")
		}
		return nil
	}
}

func validateYN(q string) ValidateFunc {
	return func(ans string) error {
		logging.L.Debug("received response")
		logging.L.Log("question", q, "answer", ans)

		if ans != "y" && ans != "n" {
			return fmt.Errorf("input must be y or n")
		}
		return nil
	}
}

func (i *UI) isYesOrNo(q string) (bool, error) {
	ans, err := i.Ask(fmt.Sprintf("%s [y/n] ", q), &Options{
		HideOrder:    true,
		Loop:         true,
		ValidateFunc: validateYN(q),
	})
	if err != nil {
		return false, errors.WithStack(err)
	}
	return ans == "y", nil
}

func (u *UI) trueFalseBool(q string) (bool, error) {
	ans, err := u.Ask(fmt.Sprintf("%s [true/fase ]", q), &Options{
		HideOrder:    true,
		Loop:         true,
		ValidateFunc: validateTF(q),
	})
	if err != nil {
		return false, errors.WithStack(err)
	}
	return ans == "true", nil
}

func (u *UI) notEmpty() ValidateFunc {
	return func(ans string) error {
		if ans == "" {
			return ErrEmpty
		}
		return nil
	}

}

type ChainFunc func(func(string) error) func(string) error
