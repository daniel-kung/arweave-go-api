package e

import (
	"fmt"
)

type ResponseBase struct {
	Code    ERRCode `json:"code"`    // error code of this api server
	Message string  `json:"message"` // error message
	Field   string  `json:"field"`   // special field prompt, especially when an error occu
}

// Response data response
type Response struct {
	ResponseBase
	Data interface{} `json:"data"` // data
}

// ListResponse List response
type ListResponse struct {
	ResponseBase
	Offset int64       `json:"offset"` // pagination offset
	Limit  int64       `json:"limit"`  // pagination limit
	Total  int64       `json:"total"`  // total record count
	Data   interface{} `json:"data"`   // list data
}

func NewResponse(code ERRCode) *Response {
	return &Response{
		ResponseBase: ResponseBase{
			Code:    code,
			Message: GetMsg(code),
		},
	}
}

func (r *Response) WithField(f string) *Response {
	r.Field = f
	return r
}

func (r *Response) Error() string {
	return fmt.Sprintf("code: %d, message: %s, field: %s", r.Code, r.Message, r.Field)
}

func NewListResponse(code ERRCode, offset, limit, total int64) *ListResponse {
	return &ListResponse{
		ResponseBase: ResponseBase{
			Code:    code,
			Message: GetMsg(code),
		},
		Offset: offset,
		Limit:  limit,
		Total:  total,
	}
}

func (r *ListResponse) WithField(f string) *ListResponse {
	r.Field = f
	return r
}

func (r *ListResponse) Error() string {
	return fmt.Sprintf("code: %d, message: %s, field: %s", r.Code, r.Message, r.Field)
}
