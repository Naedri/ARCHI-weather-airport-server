package main

import (
	"meteo_des_aeroports/internal/handlers"
	"net/http"

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
)

func main() {

	router := echo.New()

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	router.GET("/", func(c echo.Context) error {
		// handlers.GetValueOfDataTypeWithRange("NYC", "-inf", "+inf", "temperature")
		return c.String(http.StatusOK, "Welcome!")
	})

	router.GET("/iata/:IATA/probes", func(c echo.Context) error {
		iata := c.Param("IATA")
		start := c.QueryParam("start")

		if start == "" {
			start = "-inf"
		}

		end := c.QueryParam("end")

		if end == "" {
			end = "+inf"
		}

		dataType := c.QueryParam("dataType")

		if dataType == "" {
			return c.String(http.StatusBadRequest, "dataType is required")
		}

		result, err := handlers.GetValueOfDataTypeWithRange(iata, start, end, dataType)

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.String(http.StatusOK, result)
	})

	router.GET("/iata/:IATA/probes/average", func(c echo.Context) error {
		iata := c.Param("IATA")

		result, err := handlers.GetAverageValueOfTheDay(iata)

		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}

		return c.String(http.StatusOK, result)
	})

	router.Logger.Fatal(router.Start(":8080"))
}
