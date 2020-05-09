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
	ErrInvalidCategory   = errors.New("invalid category param")
	ErrRequestFailed     = errors.New("failed to perform http request")
	ErrInvalidWidgetType = errors.New("invalid widget type")
)
