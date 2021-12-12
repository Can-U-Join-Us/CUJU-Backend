package errors

import (
	"github.com/Can-U-Join-Us/CUJU-Backend/modules/logging"
)

func Check(err error) error {
	if err != nil {
		logging.Warn(err)
		return err
	}
	return nil
}
