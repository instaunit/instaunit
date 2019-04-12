package crypto

import (
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

// Create a certificate pool containing all the certificates found in the specified directory.
func NewCertPoolFromCertificatesInDirectory(dir string) (*x509.CertPool, error) {

	file, err := os.Open(dir)
	if err != nil {
		return nil, err
	}

	names, err := file.Readdirnames(1024)
	if err != nil {
		return nil, err
	}

	pool := x509.NewCertPool()

	for _, e := range names {

		switch strings.ToLower(path.Ext(e)) {
		case ".crt", ".cert", ".pem": // ok
		default:
			continue // skip
		}

		file, err = os.Open(path.Join(dir, e))
		if err != nil {
			return nil, err
		}

		data, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}

		if !pool.AppendCertsFromPEM(data) {
			return nil, fmt.Errorf("No certificates found")
		}

	}

	return pool, nil
}
