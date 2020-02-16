package main

import (
	"fmt"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nautilus/gateway"
	"github.com/nautilus/graphql"
	"net/http"
)

func main() {
	fmt.Println("Starting introspect .... ")
	schemas, err := graphql.IntrospectRemoteSchemas(
		"http://localhost:8081/graphql",
		//"http://localhost:8083",
	)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Creating gateway .... ")
	gw, err := gateway.New(schemas)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Booting server .... ")
	e := echo.New()
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins:     []string{"http://localhost:4200"},
		AllowMethods:     []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowCredentials: true,
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.Any("/graphql-playground", func(context echo.Context) error {
		gw.PlaygroundHandler(context.Response(), context.Request())
		return nil
	})
	e.Logger.Fatal(e.Start(":1323"))
}
