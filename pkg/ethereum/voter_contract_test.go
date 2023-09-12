package ethereum

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"testing"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/masa-finance/masa-oracle/pkg/crypto"
)

func init() {
	f, err := os.OpenFile("masa_oracle_node.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
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

	err = godotenv.Load(envFilePath)
	if err != nil {
		logrus.Error("Error loading .env file")
	}
}

func TestAddUser(t *testing.T) {
	privKey, err := crypto.GetOrCreatePrivateKey(os.Getenv("private.key"))
	if err != nil {
		logrus.Fatal(err)
	}

	result, err := AddUser(privKey, 31337, "testUser", "100")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(result)
}
