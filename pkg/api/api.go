package api

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"

	masa "github.com/masa-finance/masa-oracle/pkg"
)

type API struct {
	Node *masa.OracleNode
}

func NewAPI(node *masa.OracleNode) *API {
	return &API{Node: node}
}

func GetPathInt(ctx *gin.Context, name string) (int, error) {
	val, ok := ctx.GetQuery(name)
	if !ok {
		return 0, errors.New(fmt.Sprintf("the value for path parameter %s empty or not specified", name))
	}
	return strconv.Atoi(val)
}
