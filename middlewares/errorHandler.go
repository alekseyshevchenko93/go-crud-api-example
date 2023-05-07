package middlewares

import (
	"net/http"

	echo "github.com/labstack/echo/v4"
)

func ErrorHandler(err error, c echo.Context) {
	code := http.StatusInternalServerError
	message := "Internal Server Error"

	if he, ok := err.(*echo.HTTPError); ok {
		code = he.Code

		if code != http.StatusInternalServerError {
			message = he.Message.(string)
		}
	}

	c.JSON(code, map[string]string{"message": message})
}
