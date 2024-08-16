package masacrypto

import (
	"github.com/ipfs/go-cid"
	mh "github.com/multiformats/go-multihash"
	"github.com/sirupsen/logrus"
)

// ComputeSha256Cid calculates the CID (Content Identifier) for a given string.
//
// Parameters:
//   - str: The input string for which to compute the CID.
//
// Returns:
//   - string: The computed CID as a string.
//   - error: An error, if any occurred during the CID computation.
//
// The function uses the multihash package to create a SHA2-256 hash of the input string.
// It then creates a CID (version 1) from the multihash and returns the CID as a string.
// If an error occurs during the multihash computation or CID creation, it is returned.
func ComputeSha256Cid(str string) (string, error) {
	logrus.Infof("Computing CID for string: %s", str)
	// Create a multihash from the string
	mhHash, err := mh.Sum([]byte(str), mh.SHA2_256, -1)
	if err != nil {
		logrus.Errorf("Error computing multihash for string: %s, error: %v", str, err)
		return "", err
	}
	// Create a CID from the multihash
	cidKey := cid.NewCidV1(cid.Raw, mhHash).String()
	logrus.Infof("Computed CID: %s", cidKey)
	return cidKey, nil
}
