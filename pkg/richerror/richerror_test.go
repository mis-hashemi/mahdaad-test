package richerror_test

import (
	"errors"
	richerror2 "github.com/mis-hashemi/mahdaad-test/pkg/richerror"
	"net/http"
	"testing"
)

func TestRichError_Chaining(t *testing.T) {
	base := richerror2.New("Service.SendSMS").
		WithKind(richerror2.KindInvalid).
		WithMessage("invalid phone number").
		WithMeta(map[string]any{"phone": "123"}).
		WithErr(richerror2.New("underlying parse error"))

	if base.Kind() != richerror2.KindInvalid {
		t.Errorf("expected KindInvalid, got %v", base.Kind())
	}

	if base.Error() != "invalid phone number" {
		t.Errorf("expected message 'invalid phone number', got %s", base.Error())
	}

	if base.Meta()["phone"] != "123" {
		t.Errorf("expected meta phone=123, got %v", base.Meta()["phone"])
	}
}

func TestRichError_Unwrap(t *testing.T) {
	inner := richerror2.New("DB.Query").WithMessage("sql error")
	outer := richerror2.New("Repo.FindUser").WithErr(inner)

	if !errors.Is(outer, inner) {
		t.Errorf("outer should unwrap to inner")
	}
}

func TestHTTPMapping_ClientError(t *testing.T) {
	re := richerror2.New("Handler.Create").
		WithKind(richerror2.KindInvalid).
		WithMessage("invalid input")

	resp, status := richerror2.ToHTTP(re)

	if status != http.StatusUnprocessableEntity {
		t.Errorf("expected 422, got %d", status)
	}

	if resp.Message != "invalid input" {
		t.Errorf("expected message 'invalid input', got %s", resp.Message)
	}
}

func TestHTTPMapping_ServerError(t *testing.T) {
	re := richerror2.New("Handler.Save").
		WithKind(richerror2.KindUnexpected).
		WithMessage("db timeout")

	resp, status := richerror2.ToHTTP(re)

	if status != http.StatusInternalServerError {
		t.Errorf("expected 500, got %d", status)
	}

	if resp.Message == "db timeout" {
		t.Errorf("should not leak server error, got %s", resp.Message)
	}
}
