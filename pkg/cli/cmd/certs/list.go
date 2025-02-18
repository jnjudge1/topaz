package certs

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/aserto-dev/topaz/pkg/cli/cc"
)

type ListCertsCmd struct {
	CertsDir string `flag:"" required:"false" default:"${topaz_certs_dir}" help:"path to dev certs folder" `
}

func (cmd *ListCertsCmd) Run(c *cc.CommonCtx) error {
	certsDir := cmd.CertsDir
	if fi, err := os.Stat(certsDir); os.IsNotExist(err) || !fi.IsDir() {
		return fmt.Errorf("directory %s not found", certsDir)
	}

	c.UI.Normal().Msgf("certs directory: %s", certsDir)

	certDetails := make(map[string]*x509.Certificate)

	for _, fqn := range getFileList(certsDir, withCerts()) {
		if fi, err := os.Stat(fqn); os.IsNotExist(err) || fi.IsDir() {
			continue
		}

		content, err := os.ReadFile(fqn)
		if err != nil {
			return err
		}

		block, _ := pem.Decode(content)

		cert, err := x509.ParseCertificate(block.Bytes)
		if err != nil {
			return err
		}

		_, fn := filepath.Split(fqn)
		certDetails[fn] = cert
	}

	table := c.UI.Normal().WithTable("File", "Not Before", "Not After", "Valid", "CN", "DNS names")

	fileNames := make([]string, 0, len(certDetails))
	for k := range certDetails {
		fileNames = append(fileNames, k)
	}

	sort.Strings(fileNames)

	table.WithTableNoAutoWrapText()
	for _, k := range fileNames {
		isValid := true
		if time.Until(certDetails[k].NotAfter) < 0 {
			isValid = false
		}

		table.WithTableRow(k,
			certDetails[k].NotBefore.Format(time.RFC3339),
			certDetails[k].NotAfter.Format(time.RFC3339),
			fmt.Sprintf("%t", isValid),
			certDetails[k].Issuer.CommonName,
			strings.Join(certDetails[k].DNSNames, ","),
		)
	}
	table.Do()

	return nil
}
