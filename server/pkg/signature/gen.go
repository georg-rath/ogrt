package signature

import (
	"io/ioutil"

	"github.com/gofrs/uuid"
)

func GenerateRandomSignatureObject(outputFile string) {
	uuid := uuid.Must(uuid.NewV4())
	copy(SignatureTemplate[85:], uuid.String())
	ioutil.WriteFile(outputFile, SignatureTemplate, 0644)
}
