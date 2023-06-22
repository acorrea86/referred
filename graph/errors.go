package graph

import (
	"context"
	"fmt"
	"github.com/99designs/gqlgen/graphql"
	"github.com/vektah/gqlparser/v2/gqlerror"
)

type AppError string

const (
	NotFound               AppError = "Could not find resource"
	UserNotFound           AppError = "could not find user"
	EmailNotFound          AppError = "Could not find email"
	PhoneNotFound          AppError = "Could not find phone"
	Unauthorized           AppError = "Unauthorized"
	Forbidden              AppError = "Forbidden"
	ErrorWithoutExtensions AppError = "No Extensions"
	AnyHow                 AppError = "Transparent"
	ServerError            AppError = "ServerError"
	MaxFileSizeError       AppError = "File size exceeds the maximum limit %s"
	ContentTypeError       AppError = "Content Type not allowed %s"
	DataSourceError        AppError = "Could not get info from datasource"
	ValidationError        AppError = "ValidationError: %s %s"
	UnauthorizedReason              = "UNAUTHORIZED"
	ForbiddenReason                 = "FORBIDDEN"
)

type AppErrorRetry string

const (
	None         AppErrorRetry = "NONE"
	Retry        AppErrorRetry = "RETRY"
	WaitAndRetry AppErrorRetry = "WAIT_AND_RETRY"
	Cancel       AppErrorRetry = "CANCEL"
)

var decisionMapLevel = map[AppError]AppErrorRetry{
	NotFound:         None,
	UserNotFound:     None,
	EmailNotFound:    None,
	PhoneNotFound:    None,
	ServerError:      Cancel,
	DataSourceError:  WaitAndRetry,
	ValidationError:  None,
	MaxFileSizeError: Cancel,
	ContentTypeError: Cancel,
	AnyHow:           Cancel,
	Unauthorized:     Cancel,
	Forbidden:        Cancel,
}

type ErrorExtensionValues struct {
	Reason string
	Code   string
	Level  string
}

func createExtensions(reason, code string, level AppErrorRetry) *ErrorExtensionValues {
	return &ErrorExtensionValues{
		Reason: reason,
		Code:   code,
		Level:  string(level),
	}
}

type ErrorExtensionParams struct {
	Reason   string
	Code     string
	AppError AppError
}

func createExtensionForAppError(params ErrorExtensionParams) (*ErrorExtensionValues, string) {
	code := ""
	retry := decisionMapLevel[params.AppError]

	if params.AppError == ErrorWithoutExtensions {
		return nil, ""
	}

	decisionMapCode := map[AppError]string{
		NotFound:         "NOT_FOUND",
		ServerError:      "SERVER_ERROR",
		DataSourceError:  "DATA_SOURCE_ERROR",
		ValidationError:  "VALIDATION_ERROR",
		MaxFileSizeError: "MAX_FILE_SIZE_ERROR",
		ContentTypeError: "CONTENT_TYPE_ERROR",
		AnyHow:           "SERVER_ERROR",
		Unauthorized:     UnauthorizedReason,
		Forbidden:        ForbiddenReason,
	}

	for key, decision := range decisionMapCode {
		if key == params.AppError {
			code = decision
			break
		}
	}

	for key, appErrorRetry := range decisionMapLevel {
		if key == params.AppError {
			retry = appErrorRetry
			break
		}
	}

	message := string(params.AppError)
	switch params.AppError {
	case ValidationError:
		message = fmt.Sprintf(string(ValidationError), params.Reason, params.Code)
		break
	case MaxFileSizeError:
		message = fmt.Sprintf(string(MaxFileSizeError), params.Reason)
		break
	case ContentTypeError:
		message = fmt.Sprintf(string(ContentTypeError), params.Reason)
		break
	}

	return createExtensions(params.Reason, code, retry), message
}

func PresentTypedError(ctx context.Context, errExtensionParam ErrorExtensionParams) *gqlerror.Error {
	errorExtensionsValues, message := createExtensionForAppError(errExtensionParam)
	presentedError := graphql.DefaultErrorPresenter(ctx, fmt.Errorf("%q", message))
	if errorExtensionsValues != nil {
		presentedError.Extensions = map[string]interface{}{
			"reason": errorExtensionsValues.Reason,
			"code":   errorExtensionsValues.Code,
			"level":  errorExtensionsValues.Level,
		}
	}
	return presentedError
}
