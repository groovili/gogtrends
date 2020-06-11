package gogtrends

import "github.com/pkg/errors"

const (
	errParsing        = "failed to parse json"
	errReqDataF       = "request data: code = %d, status = %s"
	errInvalidRequest = "invalid request param"
	errCreateRequest  = "failed to create request"
	errDoRequest      = "failed to perform request"
)

var (
	// ErrInvalidCategory - user input is not in trendsCategories list (binding to available options in Google Trends)
	ErrInvalidCategory = errors.New("invalid category param")
	// ErrRequestFailed - response status != 200
	ErrRequestFailed = errors.New("failed to perform http request")
	// ErrInvalidWidgetType - provided widget is invalid or is used for another method
	ErrInvalidWidgetType = errors.New("invalid widget type")
)
