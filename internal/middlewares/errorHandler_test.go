package middlewares

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandlerHttpError(t *testing.T) {
	e := echo.New()

	tt := []error{
		echo.NewHTTPError(http.StatusBadRequest, "Validation failed"),
		echo.NewHTTPError(http.StatusConflict, "Something already exists"),
		echo.NewHTTPError(http.StatusForbidden, "You dont have permission"),
		echo.NewHTTPError(http.StatusNotFound, "Something not found"),
		echo.NewHTTPError(http.StatusUnauthorized, "Unauthorized"),
		echo.NewHTTPError(http.StatusInternalServerError, "Internal Server Error"),
	}

	for _, testcaseErr := range tt {
		req := httptest.NewRequest(http.MethodGet, "/any-route", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		ErrorHandler(testcaseErr, ctx)

		err, ok := testcaseErr.(*echo.HTTPError)

		assert.True(t, ok)
		assert.Equal(t, rec.Code, err.Code)
		assert.Contains(t, rec.Body.String(), err.Message)
	}
}

func TestErrorHandlerNonHttpError(t *testing.T) {
	e := echo.New()

	tt := []error{
		errors.New("some bullshit happened"),
		errors.New("1"),
	}

	for _, testcaseErr := range tt {
		req := httptest.NewRequest(http.MethodGet, "/any-route", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		ErrorHandler(testcaseErr, ctx)

		assert.Equal(t, rec.Code, http.StatusInternalServerError)
		assert.Contains(t, rec.Body.String(), "Internal Server Error")
	}
}
