package fault

import (
	"errors"
)

type metaHolder interface {
	GetMeta() map[string]string
}

type attrHolder interface {
	GetAttribute(key interface{}) (interface{}, bool)
}

type tagHolder interface {
	GetTag() Tag
}

func IsServerErr(err error) bool {
	appErr := ToAppErr(err)

	if OneOfTags(appErr, TagInternal, TagDeadlineExceeded, TagUnavailable, TagNone) {
		return true
	}
	return false
}

func getRootTag(err error) (Tag, bool) {
	if err == nil {
		return TagNone, false
	}

	tagHolder, ok := err.(tagHolder)
	if ok {
		return tagHolder.GetTag(), true
	}

	return getRootTag(errors.Unwrap(err))
}

func IsTag(err error, tag Tag) bool {
	if err == nil {
		return false
	}

	tagHolder, ok := err.(tagHolder)
	if ok {
		return tagHolder.GetTag() == tag
	}

	return IsTag(errors.Unwrap(err), tag)
}

func OneOfTags(err error, tags ...Tag) bool {
	if err == nil {
		return false
	}

	tagHolder, ok := err.(tagHolder)
	if ok {
		tag := tagHolder.GetTag()
		for _, v := range tags {
			if v == tag {
				return true
			}
		}
		return false
	}

	return OneOfTags(errors.Unwrap(err), tags...)
}

func ExtractMeta(err error) map[string]string {
	if err == nil {
		return nil
	}
	meta := map[string]string{}

	e := err
	for e != nil {
		if auxHolder, ok := e.(metaHolder); ok {
			for k, v := range auxHolder.GetMeta() {
				if _, exist := meta[k]; !exist {
					meta[k] = v
				}
			}
		}

		e = errors.Unwrap(e)
	}
	return meta
}

func GetAttr(err error, key interface{}) (interface{}, bool) {
	if err == nil {
		return nil, false
	}

	if auxHolder, ok := err.(attrHolder); ok {
		return auxHolder.GetAttribute(key)
	}
	return nil, false
}

func ToAppErr(err error) *AppError {
	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}
	return Wrap(err).Tag(TagInternal)
}

func GetDescription(err error) string {
	appErr := ToAppErr(err)

	if msg := appErr.GetMessage(); msg != "" {
		return msg
	}

	if msg := getDescriptionFromKey(appErr.GetKey()); msg != "" {
		return msg
	}

	return getDescriptionFromTag(appErr.GetTag())
}

func GetKey(err error) string {
	appErr := ToAppErr(err)

	if code := appErr.GetKey(); code != "" {
		return code
	}

	return string(appErr.GetTag())
}

func GetCode(err error) int64 {
	appErr := ToAppErr(err)

	if code := getCodeFromKey(appErr.GetKey()); code != CodeUnknown {
		return code
	}

	return getCodeFromTag(appErr.GetTag())
}

func GetHttpStatusCode(err error) int {
	appErr := ToAppErr(err)

	return getHttpStatusCodeFromTag(appErr.GetTag())
}

func IsRetryable(err error) bool {
	if !IsServerErr(err) {
		return false
	}
	attrRetryable, _ := GetAttr(err, Retryable)
	if isRetryable, ok := attrRetryable.(bool); ok {
		return isRetryable
	}
	return false
}
