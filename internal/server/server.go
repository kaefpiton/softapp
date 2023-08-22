package server

import (
	"context"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"strconv"
	"strings"
)

type HTTPServer interface {
	Start()
	Stop(ctx context.Context)
}

type EchoHTTPServer struct {
	echo       *echo.Echo
	serverPort string
}

func NewEchoHTTPServer(
	ServerPort string,
) *EchoHTTPServer {
	server := &EchoHTTPServer{
		echo:       echo.New(),
		serverPort: ServerPort,
	}

	return server
}

func (s *EchoHTTPServer) Start() {

	s.echo.POST("/", s.calculate)

	func() {
		port := fmt.Sprintf(":%v", s.serverPort)
		if err := s.echo.Start(port); err != nil {
			fmt.Println("Echo error:", err)
		}
	}()
}

func (s *EchoHTTPServer) Stop(ctx context.Context) {
	err := s.echo.Shutdown(ctx)
	if err != nil {
		fmt.Println("echo error")
	}
}

const userAccessValidHeader = "superuser"

type calculateRequest struct {
	Data string `json:"data"`
}

func (s *EchoHTTPServer) calculate(ctx echo.Context) error {
	if ctx.Request().Header.Get("User-Access") != userAccessValidHeader {
		fmt.Println("invalid header")
		return ctx.JSON(http.StatusForbidden, "")
	}

	request := new(calculateRequest)
	err := ctx.Bind(&request)
	if err != nil {
		fmt.Println(err)
	}

	arr := strings.Split(request.Data, "+")
	var result int
	for _, val := range arr {
		var num int
		if rune(val[0]) != '-' {
			num, err = strconv.Atoi(val)
			if err != nil {
				return ctx.JSON(http.StatusBadRequest, "invalid data")
			}
			result = result + num
		} else {
			numpart := val[1:]
			num, err = strconv.Atoi(numpart)
			if err != nil {
				return ctx.JSON(http.StatusBadRequest, "invalid data")
			}
			result = result + num*(-1)
		}
	}

	return ctx.JSON(http.StatusOK, result)
}
