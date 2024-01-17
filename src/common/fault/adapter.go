package fault

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	appvalidator "go-kit/src/common/validator"
	"go.mongodb.org/mongo-driver/mongo"
	"strings"
)

func DBWrapf(err error, format string, v ...interface{}) *AppError {
	msg := fmt.Sprintf(format, v...)
	wErr := fmt.Errorf("%s: %w", msg, err)
	tag := getTagFromMongoErr(err)

	return &AppError{tag: tag, cause: wErr}
}

func getTagFromMongoErr(err error) Tag {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return TagNotFound
	}
	if errors.Is(err, redis.Nil) {
		return TagNotFound
	}
	if errors.Is(err, context.Canceled) {
		return TagCancelled
	}
	if errors.Is(err, context.DeadlineExceeded) {
		return TagDeadlineExceeded
	}
	var mErr mongo.WriteException
	if errors.As(err, &mErr) {
		for _, we := range mErr.WriteErrors {
			if we.Code == 11000 {
				return TagAlreadyExists
			}
		}
	}
	return TagInternal
}

func ConvertValidatorErr(err error) *AppError {
	var errs validator.ValidationErrors
	ok := errors.As(err, &errs)

	if !ok {
		return Wrap(err).Tag(TagInvalidArgument)
	}

	trans := appvalidator.GetTranslator()

	messages := make([]string, 0, len(errs))

	for _, e := range errs {
		messages = append(messages, e.Translate(trans))
	}
	msg := strings.Join(messages, "; ")

	return Wrap(fmt.Errorf(msg)).Tag(TagInvalidArgument).Message(msg)
}
