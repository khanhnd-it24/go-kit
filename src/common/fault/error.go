package fault

import (
	"context"
	"errors"
	"fmt"
)

//go:generate go run ./cmd/gen.go

type attribute struct {
	Key, Value interface{}
}

type ErrorAttribute string

const Retryable = ErrorAttribute("Retryable")

type AppError struct {
	tag     Tag
	cause   error
	key     string
	message string
	meta    map[string]string
	attr    []attribute
}

func defaultTag(err error) Tag {
	var tag Tag
	if rTag, exists := getRootTag(err); exists {
		tag = rTag
	} else if errors.Is(err, context.DeadlineExceeded) {
		tag = TagDeadlineExceeded
	} else if errors.Is(err, context.Canceled) {
		tag = TagCancelled
	} else {
		tag = TagInternal
	}
	
	return tag
}

func Wrap(err error) *AppError {
	tag := defaultTag(err)
	return &AppError{tag: tag, cause: err}
}

func Wrapf(err error, format string, v ...interface{}) *AppError {
	msg := fmt.Sprintf(format, v...)
	wErr := fmt.Errorf("%s: %w", msg, err)
	tag := defaultTag(err)

	return &AppError{tag: tag, cause: wErr}
}

func (e *AppError) Unwrap() error {
	return e.cause
}

func (e *AppError) Error() string {
	if e.cause != nil {
		return e.cause.Error()
	}
	return "unknown"
}

func (e *AppError) Tag(tag Tag) *AppError {
	e.tag = tag
	return e
}

func (e *AppError) GetTag() Tag {
	return e.tag
}

func (e *AppError) Message(msg string) *AppError {
	e.message = msg
	return e
}

func (e *AppError) GetMessage() string {
	return e.message
}

func (e *AppError) Meta(kv ...string) *AppError {
	if e.meta == nil {
		e.meta = map[string]string{}
	}
	var key interface{}
	for _, item := range kv {
		if key == nil {
			key = item
			continue
		}
		keyStr, ok := key.(string)
		if ok {
			e.meta[keyStr] = item
		}
		key = nil
	}
	return e
}

func (e *AppError) MetaMap(meta map[string]string) *AppError {
	if e.meta == nil {
		e.meta = map[string]string{}
	}

	for k, v := range meta {
		e.meta[k] = v
	}
	return e
}

func (e *AppError) GetMeta() map[string]string {
	return e.meta
}

func (e *AppError) Key(key string) *AppError {
	e.key = key
	return e
}

func (e *AppError) GetKey() string {
	return e.key
}

func (e *AppError) Attr(key, value interface{}) *AppError {
	e.attr = append(e.attr, attribute{Key: key, Value: value})
	return e
}

func (e *AppError) GetAttribute(key interface{}) (interface{}, bool) {
	for _, attr := range e.attr {
		if attr.Key == key {
			return attr.Value, true
		}
	}
	return nil, false
}
