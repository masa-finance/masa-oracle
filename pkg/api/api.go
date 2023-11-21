package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	masa "github.com/masa-finance/masa-oracle/pkg"
)

type API struct {
	Node *masa.OracleNode
}

func NewAPI(node *masa.OracleNode) *API {
	return &API{Node: node}
}

func (api *API) GetPeersHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		routingTable := api.Node.DHT.RoutingTable()
		peers := routingTable.ListPeers()
		c.JSON(http.StatusOK, gin.H{"peers": peers})
	}
}

func (api *API) GetPeerAddresses() gin.HandlerFunc {
	return func(c *gin.Context) {
		peers := api.Node.Host.Network().Peers()
		peerAddresses := make(map[string][]string)
		for _, peer := range peers {
			conns := api.Node.Host.Network().ConnsToPeer(peer)
			for _, conn := range conns {
				addr := conn.RemoteMultiaddr()
				peerAddresses[peer.String()] = append(peerAddresses[peer.String()], addr.String())
			}
		}
		c.JSON(http.StatusOK, gin.H{"peerAddresses": peerAddresses})
	}
}
