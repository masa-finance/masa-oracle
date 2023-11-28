package cicd_helpers

import (
    "fmt"
    "os"
    "os/user"
    "strings"
    "path/filepath"
    "github.com/sirupsen/logrus"
    "log"
)

func SetEnvVariablesForPipeline (multiAddr string) {
    usr, err := user.Current()
	if err != nil {
		log.Fatal("could not find user.home directory")
	}
    outputFilePath := filepath.Join(usr.HomeDir, ".masa", "masa_oracle_node_output.env")
    builder := strings.Builder{}
	builder.WriteString(fmt.Sprintf("export MASA_NODE_MULTIADDRESS='%s\n'", multiAddr))
	err = os.WriteFile(outputFilePath, []byte(builder.String()), 0755)
	if err != nil {
		logrus.Fatal("could not write to masa_oracle_node_output.env file:", err)
	}
}

