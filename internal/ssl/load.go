package ssl

import (
	"app/internal/modules/log"
	"os"
)

// CheckSSLPath checks if the SSL private key and certificate files exist
func CheckSSLPath(pk, cert string) (pkOut, certOut string, okPK, okCert bool) {
	if _, err := os.Stat(pk); os.IsNotExist(err) {
		log.With(log.ErrorString(err)).Error("SSL private key file not found: %s", pk)
		okPK = false
	} else {
		okPK = true
	}
	if _, err := os.Stat(cert); os.IsNotExist(err) {
		log.With(log.ErrorString(err)).Error("SSL certificate file not found: %s", cert)
		okCert = false
	} else {
		okCert = true
	}
	return pk, cert, okPK, okCert
}
