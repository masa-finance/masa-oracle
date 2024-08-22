package masa_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/masa-finance/masa-oracle/pkg/masacrypto"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestOracle(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Oracle integration test suite")
}

func generateNodeKeys(configPath, keyFiles string) {
	dir := filepath.Dir(configPath)
	key := filepath.Join(dir, keyFiles+".key")
	err := masacrypto.GenerateSelfSignedCert(filepath.Join(dir, keyFiles+".cert"), key)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
	// This is not optimal, but for now this is the only way to initialize a oracle node with a config
	os.Setenv("FILE_PATH", dir)
	err = os.WriteFile(configPath, []byte(generateconfig(key)), 0644)
	ExpectWithOffset(1, err).ToNot(HaveOccurred())
}

func generateconfig(path string) string {
	return `
privateKeyFile: "` + path + `"
`
}
