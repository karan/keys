package saltpack

import (
	"io"

	"github.com/keys-pub/keys"
	"github.com/pkg/errors"

	ksaltpack "github.com/keybase/saltpack"
)

// Encrypt to recipients.
// Sender can be nil, if you want it to be anonymous.
// https://saltpack.org/encryption-format-v2
func (s *Saltpack) Encrypt(b []byte, signer *keys.X25519Key, recipients ...keys.ID) ([]byte, error) {
	recs, err := s.boxPublicKeys(recipients)
	if err != nil {
		return nil, err
	}
	var sbk ksaltpack.BoxSecretKey
	if signer != nil {
		sbk = newBoxKey(signer)
	}
	return ksaltpack.Seal(ksaltpack.Version2(), b, sbk, recs)
}

// EncryptArmored to recipients.
// Sender can be nil, if you want it to be anonymous.
// https://saltpack.org/encryption-format-v2
func (s *Saltpack) EncryptArmored(b []byte, signer *keys.X25519Key, recipients ...keys.ID) (string, error) {
	recs, err := s.boxPublicKeys(recipients)
	if err != nil {
		return "", err
	}
	var sbk ksaltpack.BoxSecretKey
	if signer != nil {
		sbk = newBoxKey(signer)
	}
	return ksaltpack.EncryptArmor62Seal(ksaltpack.Version2(), b, sbk, recs, "")
}

// Decrypt bytes.
// If there was a signer, will return a X25519 key ID.
func (s *Saltpack) Decrypt(b []byte) ([]byte, keys.ID, error) {
	info, out, err := ksaltpack.Open(encryptVersionValidator, b, s)
	if err != nil {
		return nil, "", convertBoxKeyErr(err)
	}
	signer := keys.ID("")
	if !info.SenderIsAnon {
		kid, err := x25519KeyID(info.SenderKey.ToKID())
		if err != nil {
			return nil, "", errors.Wrapf(err, "failed to decrypt")
		}
		signer = kid
	}
	return out, signer, nil
}

// DecryptArmored text.
// If there was a signer, will return a X25519 key ID.
func (s *Saltpack) DecryptArmored(str string) ([]byte, keys.ID, error) {
	// TODO: Casting to string could be a performance issue
	info, out, _, err := ksaltpack.Dearmor62DecryptOpen(encryptVersionValidator, str, s)
	if err != nil {
		return nil, "", convertBoxKeyErr(err)
	}
	signer := keys.ID("")
	if !info.SenderIsAnon {
		kid, err := x25519KeyID(info.SenderKey.ToKID())
		if err != nil {
			return nil, "", errors.Wrapf(err, "failed to decrypt")
		}
		signer = kid
	}
	return out, signer, nil
}

// NewEncryptStream creates an encrypted io.WriteCloser.
// Sender can be nil, if you want it to be anonymous.
func (s *Saltpack) NewEncryptStream(w io.Writer, signer *keys.X25519Key, recipients ...keys.ID) (io.WriteCloser, error) {
	recs, err := s.boxPublicKeys(recipients)
	if err != nil {
		return nil, err
	}
	var sbk ksaltpack.BoxSecretKey
	if signer != nil {
		sbk = newBoxKey(signer)
	}
	return ksaltpack.NewEncryptStream(ksaltpack.Version2(), w, sbk, recs)
}

// NewEncryptArmoredStream creates an encrypted armored io.WriteCloser.
// Sender can be nil, if you want it to be anonymous.
func (s *Saltpack) NewEncryptArmoredStream(w io.Writer, signer *keys.X25519Key, recipients ...keys.ID) (io.WriteCloser, error) {
	recs, err := s.boxPublicKeys(recipients)
	if err != nil {
		return nil, err
	}
	var sbk ksaltpack.BoxSecretKey
	if signer != nil {
		sbk = newBoxKey(signer)
	}
	return ksaltpack.NewEncryptArmor62Stream(ksaltpack.Version2(), w, sbk, recs, "")
}

// NewDecryptStream create decryption stream.
// If there was a signer, will return a X25519 key ID.
func (s *Saltpack) NewDecryptStream(r io.Reader) (io.Reader, keys.ID, error) {
	info, stream, err := ksaltpack.NewDecryptStream(encryptVersionValidator, r, s)
	if err != nil {
		return nil, "", convertBoxKeyErr(err)
	}
	signer := keys.ID("")
	if !info.SenderIsAnon {
		kid, err := x25519KeyID(info.SenderKey.ToKID())
		if err != nil {
			return nil, "", errors.Wrapf(err, "failed to decrypt")
		}
		signer = kid
	}
	return stream, signer, nil
}

// NewDecryptArmoredStream creates decryption stream.
// If there was a signer, will return a X25519 key ID.
func (s *Saltpack) NewDecryptArmoredStream(r io.Reader) (io.Reader, keys.ID, error) {
	// TODO: Specifying nil for resolver will panic if box keys not found
	info, stream, _, err := ksaltpack.NewDearmor62DecryptStream(encryptVersionValidator, r, s)
	if err != nil {
		return nil, "", convertBoxKeyErr(err)
	}
	signer := keys.ID("")
	if !info.SenderIsAnon {
		kid, err := x25519KeyID(info.SenderKey.ToKID())
		if err != nil {
			return nil, "", errors.Wrapf(err, "failed to decrypt")
		}
		signer = kid
	}
	return stream, signer, nil
}
