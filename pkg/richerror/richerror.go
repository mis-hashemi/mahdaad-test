package richerror

import "errors"

type Kind int

const (
	KindInvalid Kind = iota + 1
	KindForbidden
	KindNotFound
	KindUnexpected
	KindRateLimit
	KindBadRequest
)

type Op string

type RichError struct {
	operation    Op
	wrappedError error
	message      string
	kind         Kind
	meta         map[string]any
}

func New(op Op) *RichError {
	return &RichError{operation: op}
}

func (r *RichError) WithOp(op Op) *RichError {
	r.operation = op
	return r
}

func (r *RichError) WithErr(err error) *RichError {
	r.wrappedError = err
	return r
}

func (r *RichError) WithMessage(message string) *RichError {
	r.message = message
	return r
}

func (r *RichError) WithKind(kind Kind) *RichError {
	r.kind = kind
	return r
}

func (r *RichError) WithMeta(meta map[string]any) *RichError {
	r.meta = meta
	return r
}

func (r *RichError) Meta() map[string]any {
	if len(r.meta) != 0 {
		return r.meta
	}
	if re, ok := r.wrappedError.(*RichError); ok {
		return re.Meta()
	}
	return map[string]any{}
}

func (r *RichError) Error() string {
	if r.message != "" {
		return r.message
	}
	if re, ok := r.wrappedError.(*RichError); ok {
		return re.Error()
	}
	if r.wrappedError != nil {
		return r.wrappedError.Error()
	}
	return "unknown error"
}

func (r *RichError) Kind() Kind {
	if r.kind != 0 {
		return r.kind
	}
	var re *RichError
	if errors.As(r.wrappedError, &re) {
		return re.Kind()
	}
	return KindUnexpected
}

func (r *RichError) Unwrap() error {
	return r.wrappedError
}
