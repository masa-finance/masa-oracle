package ethereum

import (
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
	"github.com/masa-finance/masa-oracle/pkg/crypto"
)

func init() {
	f, err := os.OpenFile("identity_test.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
	if err != nil {
		log.Fatal(err)
	}
	mw := io.MultiWriter(os.Stdout, f)
	logrus.SetOutput(mw)
	logrus.SetLevel(logrus.DebugLevel)

	usr, err := user.Current()
	if err != nil {
		log.Fatal("could not find user.home directory")
	}
	envFilePath := filepath.Join(usr.HomeDir, ".masa", "masa_oracle_node.env")
	//keyFilePath := filepath.Join(usr.HomeDir, ".masa", "masa_oracle_key")

	err = godotenv.Load(envFilePath)
	if err != nil {
		logrus.Error("Error loading .env file")
	}
}

func TestMint(t *testing.T) {
	toAddress := "0x52f823a4dbe2Dc2934d5F5a854dCb8B407FEa24A"
	privateKey, _, err := crypto.GetOrCreatePrivateKey(os.Getenv(masa.KeyFileKey))
	if err != nil {
		t.Fatal(err)
	}
	err = Mint(privateKey, toAddress)
	if err != nil {
		t.Error(err)
	}
}
