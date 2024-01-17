package fault

import (
	"google.golang.org/grpc/codes"
	"net/http"
)

type Tag string

const (
	TagNone               Tag = ""
	TagInvalidArgument    Tag = "InvalidArgument"
	TagAlreadyExists      Tag = "AlreadyExists"
	TagNotFound           Tag = "NotFound"
	TagInternal           Tag = "Internal"
	TagFailedPrecondition Tag = "FailedPrecondition"
	TagUnauthenticated    Tag = "Unauthenticated"
	TagPermissionDenied   Tag = "PermissionDenied"
	TagDeadlineExceeded   Tag = "TagDeadlineExceeded"
	TagCancelled          Tag = "TagCancelled"
	TagUnavailable        Tag = "Unavailable"
)

var (
	tagFactoryMessage = map[Tag]string{
		TagNone:               "something went wrong",
		TagInvalidArgument:    "invalid argument",
		TagAlreadyExists:      "present already exists",
		TagNotFound:           "present not found",
		TagInternal:           "internal error",
		TagFailedPrecondition: "failed precondition",
		TagUnauthenticated:    "unauthenticated",
		TagPermissionDenied:   "permission denied",
		TagDeadlineExceeded:   "timeout operation",
		TagCancelled:          "tag canceled",
		TagUnavailable:        "service is unavailable",
	}
)

var (
	tagFactoryHttpCode = map[Tag]int{
		TagNone:               http.StatusBadRequest,
		TagInvalidArgument:    http.StatusBadRequest,
		TagAlreadyExists:      http.StatusBadRequest,
		TagNotFound:           http.StatusBadRequest,
		TagInternal:           http.StatusInternalServerError,
		TagFailedPrecondition: http.StatusBadRequest,
		TagUnauthenticated:    http.StatusUnauthorized,
		TagPermissionDenied:   http.StatusForbidden,
		TagCancelled:          http.StatusBadRequest,
		TagDeadlineExceeded:   http.StatusServiceUnavailable,
		TagUnavailable:        http.StatusServiceUnavailable,
	}
)

var (
	tagFactoryGrpcCode = map[Tag]codes.Code{
		TagNone:               codes.Unknown,
		TagInvalidArgument:    codes.InvalidArgument,
		TagAlreadyExists:      codes.AlreadyExists,
		TagNotFound:           codes.NotFound,
		TagInternal:           codes.Internal,
		TagFailedPrecondition: codes.FailedPrecondition,
		TagUnauthenticated:    codes.Unauthenticated,
		TagPermissionDenied:   codes.PermissionDenied,
		TagCancelled:          codes.Canceled,
		TagDeadlineExceeded:   codes.DeadlineExceeded,
		TagUnavailable:        codes.Unavailable,
	}
)

func getDescriptionFromTag(tag Tag) string {
	if des, exists := tagFactoryMessage[tag]; exists {
		return des
	}
	return "unknown"
}

func getHttpStatusCodeFromTag(tag Tag) int {
	if des, exists := tagFactoryHttpCode[tag]; exists {
		return des
	}
	return http.StatusInternalServerError
}

func getCodeFromTag(tag Tag) int64 {
	code := getGrpcStatusCodeFromTag(tag)

	return int64(code)
}

func getGrpcStatusCodeFromTag(tag Tag) codes.Code {
	if des, exists := tagFactoryGrpcCode[tag]; exists {
		return des
	}
	return codes.Unknown
}
