package mapper

import (
	"net/http"
)

// ErrorMapper maps error on ErrorInfo
type ErrorMapper struct {
	mapper ErrorMap
}

// NewErrorMapper init new ErrorMapper
func NewErrorMapper(m ErrorMap) ErrorMapper {
	mapper := ErrorMapper{mapper: m}
	return mapper
}

// ErrorInfo holds status code(http statuses) and additional info
type ErrorInfo struct {
	StatusCode int
	Msg        string
}

type ErrorMap map[error]ErrorInfo

// MapError for provided error returns from ErrorMapper instance ErrorInfo,
// if error not found returns 500,"Internal server error"
func (m ErrorMapper) MapError(err error) ErrorInfo {
	if v, ok := m.mapper[err]; ok {
		return v
	}

	inf := ErrorInfo{
		StatusCode: http.StatusInternalServerError,
		Msg:        "Internal server error",
	}
	return inf
}
