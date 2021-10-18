package main

import (
	"meteo_des_aeroports/internal/handlers"
	"net/http"

	_ "meteo_des_aeroports/cmd/httpServer/docs" // docs is generated by Swag CLI, you have to import it.

	echo "github.com/labstack/echo/v4"
	middleware "github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger" // echo-swagger middleware
)

// @title Meteo des aéroports API
// @version 1.0
// @description Une API pour récupérer les données des sondes de différents aéroports.
// @termsOfService http://swagger.io/terms/

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /
// @schemes http
func main() {

	router := echo.New()
	router.GET("/swagger/*", echoSwagger.WrapHandler)

	router.Use(middleware.Logger())
	router.Use(middleware.Recover())

	router.GET("/swagger/*", echoSwagger.WrapHandler)
	router.GET("/", HealthCheck)

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

// HealthCheck godoc
// @Summary Show the status of server.
// @Description get the status of server.
// @Tags root
// @Accept */*
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router / [get]
func HealthCheck(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": "Server is up and running",
	})
}