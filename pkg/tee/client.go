package tee

import (
	"os"

	worker "github.com/masa-finance/tee-worker/pkg/client"
)

var teeWorkerURL = os.Getenv("TEE_WORKER_URL")

func NewClient() *worker.Client {
	return worker.NewClient(teeWorkerURL)
}
