package pgpprovider

import (
	"bytes"
	"golang.org/x/crypto/openpgp"
	"golang.org/x/crypto/openpgp/armor"
	"golang.org/x/crypto/openpgp/packet"
	"io"
	"os"
)

func ParsePgpKey(path, passphrase string) (*openpgp.Entity, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	block, err := armor.Decode(file)
	if err != nil {
		return nil, err
	}

	pgpEntity, err := openpgp.ReadEntity(packet.NewReader(block.Body))
	if err != nil {
		return nil, err
	}

	if passphrase != "" {
		for _, subKey := range pgpEntity.Subkeys {
			if subKey.PrivateKey == nil {
				continue
			}

			if ierr := subKey.PrivateKey.Decrypt([]byte(passphrase)); err != nil {
				return nil, ierr
			}
		}
	}
	return pgpEntity, nil
}

func Encrypt(buffer *bytes.Buffer, entity *openpgp.Entity) (*bytes.Buffer, error) {
	encBuff := &bytes.Buffer{}

	armoredWriter, err := armor.Encode(encBuff, "PGP MESSAGE", nil)
	if err != nil {
		return nil, err
	}

	entityList := openpgp.EntityList{
		entity,
	}

	cipheredWriter, err := openpgp.Encrypt(armoredWriter, entityList, nil, nil, nil)
	if err != nil {
		return nil, err
	}

	_, err = cipheredWriter.Write(buffer.Bytes())
	if err != nil {
		return nil, err
	}

	_ = cipheredWriter.Close()
	_ = armoredWriter.Close()
	return encBuff, nil
}

func Decrypt(buff *bytes.Buffer, entity *openpgp.Entity) (io.Reader, error) {
	entityList := openpgp.EntityList{
		entity,
	}
	md, err := openpgp.ReadMessage(buff, entityList, nil, nil)
	if err != nil {
		return nil, err
	}
	return md.UnverifiedBody, nil
}
