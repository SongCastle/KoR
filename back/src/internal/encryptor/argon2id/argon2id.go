package argon2id

import (
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"regexp"
	"strconv"

	"github.com/SongCastle/KoR/internal/log"
	"github.com/SongCastle/KoR/internal/random"
	"golang.org/x/crypto/argon2"
)

const (
	Version     = argon2.Version
	Memory      = uint32(32 << 10) // 32 MiB
	TimeCost    = uint32(3)
	Parallelism = uint8(4)
	KeyLen      = uint32(32)
	SaltLen     = uint32(32)
	HashFormat  = "$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s"
)

var (
	Encoder   *base64.Encoding = base64.RawStdEncoding
	Decoder   *base64.Encoding = Encoder.Strict()
	Regex     *regexp.Regexp   = regexp.MustCompile(`^\$argon2id`)
	HashRegex *regexp.Regexp   = regexp.MustCompile(`^\$argon2id\$v=(\d+)\$m=(\d+),t=(\d+),p=(\d+)\$([[:alnum:]+/]+)\$([[:alnum:]+/]+)$`)
	InvalidHashFormat error    = errors.New("invalid_hash_format")
	InvalidComplexity error    = errors.New("invalid_complexity")
)

type rawHash struct {
	salt, hash []byte
}

func Digest(password string) (string, error) {
	salt := []byte(random.Generate(int(SaltLen)))
	hash := argon2.IDKey(
		[]byte(password), salt, TimeCost, Memory, Parallelism, KeyLen,
	)
	return fmt.Sprintf(
		HashFormat,
		Version, Memory, TimeCost, Parallelism,
		Encoder.EncodeToString(salt), Encoder.EncodeToString(hash),
	), nil
}

func Compare(hash, password string) bool {
	raw, err := decodeHash(hash)
	if err != nil {
		log.Warnf("argon2id Compare Error: %v", err)
		return false
	}
	calcHash := argon2.IDKey(
		[]byte(password), raw.salt, TimeCost, Memory, Parallelism, KeyLen,
	)
	return subtle.ConstantTimeCompare(raw.hash, calcHash) == 1
}

func EncryptedByArgon2id(hash string) bool {
	return Regex.MatchString(hash)
}

func decodeHash(hash string) (*rawHash, error) {
	keys := HashRegex.FindStringSubmatch(hash)
	if len(keys) != 7 {
		return nil, InvalidHashFormat
	}
	if v, err := strconv.Atoi(keys[1]); v != Version || err != nil {
		return nil, InvalidComplexity
	}
	if m, err := strconv.Atoi(keys[2]); uint32(m) != Memory || err != nil {
		return nil, InvalidComplexity
	}
	if t, err := strconv.Atoi(keys[3]); uint32(t) != TimeCost || err != nil {
		return nil, InvalidComplexity
	}
	if p, err := strconv.Atoi(keys[4]); uint8(p) != Parallelism || err != nil {
		return nil, InvalidComplexity
	}
	var err error
	raw := &rawHash{}
	raw.salt, err = Decoder.DecodeString(keys[5])
	if err != nil {
		return nil, err
	}
	raw.hash, err = Decoder.DecodeString(keys[6])
	if err != nil {
		return nil, err
	}
	if uint32(len(raw.hash)) != KeyLen {
		return nil, InvalidComplexity
	}
	return raw, nil
}
