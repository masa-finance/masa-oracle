package tee

/*

A way to seal the response body before writing it.

Note: This is a simple approach.

Another approach could have been by sealing every
answer at protocol level by sealing directly the
results of the worker.

*/

import (
	"github.com/gin-gonic/gin"
)

// GinTEESealer is a wrapper around the gin.ResponseWriter that seals the response body before writing it.
type GinTEESealer struct {
	gin.ResponseWriter
}

func (w GinTEESealer) Write(b []byte) (int, error) {
	sealedData, err := Seal(b)
	if err != nil {
		return 502, err
	}
	return w.ResponseWriter.Write(sealedData)
}

func RegisterGIN(c *gin.Context) {
	blw := &GinTEESealer{ResponseWriter: c.Writer}
	c.Writer = blw
	c.Next()
}
