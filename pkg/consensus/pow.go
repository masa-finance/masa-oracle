package consensus

import (
	"crypto/sha256"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
)

// @@note WIP for future use case
// PoW Usage
// apiKey := consensus.GeneratePoW(node.Host.ID().String())
// logrus.Infof("api key: %s", apiKey)
// PoW

func GeneratePoW(peerId string) string {
	difficulty := 1
	// difficulty := time.Now().Unix()

	targetTime := 5 * time.Second
	start := time.Now()
	proofOfWork, err := ComputeProofOfWork(peerId, difficulty)
	if err != nil {
		logrus.Error(err)
	}
	elapsed := time.Since(start)

	difficulty = adjustDifficulty(difficulty, targetTime, elapsed)

	logrus.Infof("Peer ID: %s\n", peerId)
	logrus.Infof("Difficulty: %d\n", difficulty)
	logrus.Infof("Proof of Work: %s\n", proofOfWork)
	logrus.Infof("Elapsed: %s\n", elapsed)

	return proofOfWork
}

func ComputeProofOfWork(peerId string, difficulty int) (string, error) {
	mhHash, err := mh.Sum([]byte(peerId), mh.SHA2_256, -1)
	if err != nil {
		return "", err
	}
	cidKey := cid.NewCidV1(cid.Raw, mhHash).String()
	nonce := 0
	for {
		data := fmt.Sprintf("%s%d", cidKey, nonce)
		hash := sha256.Sum256([]byte(data))
		if hasLeadingZeroes(hash[:], difficulty) {
			return fmt.Sprintf("%x", hash), nil
		}
		nonce++
	}
}

func adjustDifficulty(difficulty int, targetTime, actualTime time.Duration) int {
	if actualTime < targetTime {
		difficulty++
	} else if actualTime > targetTime {
		difficulty--
	}
	return difficulty
}

func hasLeadingZeroes(hash []byte, difficulty int) bool {
	for i := 0; i < difficulty; i++ {
		byteIndex := i / 8
		bitIndex := i % 8
		if byteIndex >= len(hash) || hash[byteIndex]&(1<<(7-bitIndex)) != 0 {
			return false
		}
	}
	return true
}
