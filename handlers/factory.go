package handlers

import (
	"github.com/onepointsixtwo/golangredisserver/router"
)

type Factory interface {
	AddHandlersToRouter(router router.Router)
}
