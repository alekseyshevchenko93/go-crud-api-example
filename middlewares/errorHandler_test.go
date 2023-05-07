package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestErrorHandler(t *testing.T) {
	e := echo.New()

	tt := []struct {
		Code    int
		Message string
	}{
		{Code: http.StatusBadRequest, Message: "Validation failed"},
		{Code: http.StatusConflict, Message: "Something already exists"},
		{Code: http.StatusForbidden, Message: "You dont have permission"},
		{Code: http.StatusNotFound, Message: "Something not found"},
		{Code: http.StatusUnauthorized, Message: "Unauthorized"},
		{Code: http.StatusInternalServerError, Message: "Internal Server Error"},
	}

	for _, testcase := range tt {
		req := httptest.NewRequest(http.MethodGet, "/any-route", nil)
		rec := httptest.NewRecorder()
		ctx := e.NewContext(req, rec)

		err := echo.NewHTTPError(testcase.Code, testcase.Message)

		ErrorHandler(err, ctx)

		assert.Equal(t, rec.Code, testcase.Code)
		assert.Contains(t, rec.Body.String(), testcase.Message)
	}
}
