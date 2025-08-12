package ucry

import (
	"crypto/rsa"
	"errors"
	"sync"
)

type RSAHolder struct {
	rSA
	pubKeyFmt string
	priKeyFmt string
	pubKey    *rsa.PublicKey
	priKey    *rsa.PrivateKey
	mu        sync.Mutex
}

func NewRSAHolder(pubKey, priKey string) (*RSAHolder, error) {
	rsaHolder := &RSAHolder{}
	rsaHolder.mu.Lock()
	defer rsaHolder.mu.Unlock()

	err := rsaHolder.parsePubKey([]byte(pubKey))
	if err != nil {
		return nil, err
	}
	err = rsaHolder.parsePriKey([]byte(priKey))

	if err != nil {
		return nil, err
	}

	return rsaHolder, nil
}

func (r *RSAHolder) parsePubKey(pubKeys []byte) error {
	r.pubKeyFmt = r.fmtKey(string(pubKeys), true)
	pubKey, err := r.rSA.parsePubKey([]byte(r.pubKeyFmt))
	if err != nil {
		return err
	}

	r.pubKey = pubKey

	return nil
}

func (r *RSAHolder) parsePriKey(priKeys []byte) error {
	r.priKeyFmt = r.fmtKey(string(priKeys), false)
	priKey, err := r.rSA.parsePriKey([]byte(r.priKeyFmt))
	if err != nil {
		return err
	}

	r.priKey = priKey
	return nil
}

func (r *RSAHolder) Encrypt(plains []byte) ([]byte, error) {
	return r._encrypt(r.pubKey, plains)
}

func (r *RSAHolder) EncryptPlus(plain string) (string, error) {
	//pubKey = r.fmtKey(pubKey, true)

	if ciphers, err := r.Encrypt([]byte(plain)); err != nil {
		return "", err
	} else {
		return string(B64.Encrypt(ciphers)), nil
	}
}

func (r *RSAHolder) Decrypt(ciphers []byte) ([]byte, error) {
	return r._decrypt(r.priKey, ciphers)
}

func (r *RSAHolder) DecryptPlus(plain string) (string, error) {
	//priKey = r.fmtKey(priKey, false)

	b64ciphers, err := B64.Decrypt([]byte(plain))
	if err != nil {
		return "", err
	}

	if ciphers, err := r.Decrypt(b64ciphers); err != nil {
		return "", err
	} else {
		return string(ciphers), nil
	}
}

type rsaCache struct {
	m sync.Map
}

func (r *rsaCache) Store(key, pubKey, priKey string) error {
	//r.m.Delete(key)
	rsaHolder, err := NewRSAHolder(pubKey, priKey)
	if err != nil {
		return err
	}
	r.m.Store(key, rsaHolder)
	return nil
}

func (r *rsaCache) Get(key string) (*RSAHolder, error) {
	if val, ok := r.m.Load(key); !ok {
		return nil, errors.New("can't found rsa Holder with key:" + key)
	} else {
		return val.(*RSAHolder), nil
	}
}

var RsaCache = &rsaCache{}
