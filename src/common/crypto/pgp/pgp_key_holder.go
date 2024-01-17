package pgpprovider

import (
	"go-kit/src/common/configs"
	"golang.org/x/crypto/openpgp"
	"strings"
)

type PgpKeyHolder struct {
	pgpMapKey map[string]*openpgp.Entity
}

func (h *PgpKeyHolder) GetKey(code string) (*openpgp.Entity, bool) {
	key, exists := h.pgpMapKey[strings.ToLower(code)]
	return key, exists
}

func NewPgpKeyHolder(cf *configs.Config) (*PgpKeyHolder, error) {
	mapPgpKeys := make(map[string]*openpgp.Entity)

	for _, key := range cf.PgpKeys {
		pgpEntity, err := ParsePgpKey(key.Path, key.Passphrase)
		if err == nil {
			mapPgpKeys[strings.ToLower(key.Name)] = pgpEntity
		}
	}

	return &PgpKeyHolder{mapPgpKeys}, nil
}
