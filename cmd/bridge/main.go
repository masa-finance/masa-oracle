package main

import (
	"io"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	masa "github.com/masa-finance/masa-oracle/pkg"
)

func init() {
	f, err := os.OpenFile("masa_bridge.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0755)
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
	envFilePath := filepath.Join(usr.HomeDir, ".masa", "masa_bridge.env")
	//certFilePath := filepath.Join(usr.HomeDir, ".masa", "webhook-selfsigned-cert.pem")
	//certKeyFilePath := filepath.Join(usr.HomeDir, ".masa", "webhook - selfsigned - key.pem")

	// Create the directories if they don't already exist
	if _, err := os.Stat(filepath.Dir(envFilePath)); os.IsNotExist(err) {
		err = os.MkdirAll(filepath.Dir(envFilePath), 0755)
		if err != nil {
			logrus.Fatal("could not create directory:", err)
		}
	}
	// Check if the .env file exists
	if _, err := os.Stat(envFilePath); os.IsNotExist(err) {
		// If not, create it with default values
		builder := strings.Builder{}
		//builder.WriteString(fmt.Sprintf("%s=%s\n", masa.Cert, certFilePath))
		//builder.WriteString(fmt.Sprintf("%s=%s\n", masa.CertPem, certKeyFilePath))
		err = os.WriteFile(envFilePath, []byte(builder.String()), 0644)
		if err != nil {
			logrus.Fatal("could not write to .env file:", err)
		}
	}
	err = godotenv.Load(envFilePath)
	if err != nil {
		logrus.Error("Error loading .env file")
	}
}

func main() {
	//certFile := os.Getenv(masa.Cert)
	//keyFile := os.Getenv(masa.CertPem)

	// Check if the certificate and key files already exist
	//if _, err := os.Stat(certFile); os.IsNotExist(err) {
	//	if _, err := os.Stat(keyFile); os.IsNotExist(err) {
	//		// If not, generate them
	//		if err := crypto.GenerateSelfSignedCert(certFile, keyFile); err != nil {
	//			logrus.Fatal("Failed to generate self-signed certificate:", err)
	//		}
	//	}
	//}

	err := masa.NewBridge()
	if err != nil {

	}
}
