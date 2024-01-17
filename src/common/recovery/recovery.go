package recovery

import (
	"context"
	"github.com/go-errors/errors"
	"go-kit/src/common/fault"
	"go-kit/src/common/logger"
)

func HandleRoutine() {
	if err := recover(); err != nil {
		goErr := errors.Wrap(err, 2)
		wErr := fault.Wrapf(goErr, "[Recovery] go routine")
		logger.Fault(context.Background(), wErr, "[Recovery] stack %s", goErr.Stack())
	}
}
