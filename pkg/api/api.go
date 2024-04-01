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

// NewAPI creates a new API instance with the given OracleNode.
func NewAPI(node *masa.OracleNode) *API {
	return &API{Node: node}
}

// GetPathInt converts the path parameter with name to an int.
// It returns the int value and nil error if the path parameter is present and a valid integer.
// It returns 0 and a formatted error if the path parameter is missing or not a valid integer.
func GetPathInt(ctx *gin.Context, name string) (int, error) {
	val, ok := ctx.GetQuery(name)
	if !ok {
		return 0, errors.New(fmt.Sprintf("the value for path parameter %s empty or not specified", name))
	}
	return strconv.Atoi(val)
}
