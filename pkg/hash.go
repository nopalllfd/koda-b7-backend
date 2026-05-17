package pkg

import (
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"log"
	"strings"

	"golang.org/x/crypto/argon2"
)

type HashConfig struct {
	Memory  uint32
	Time    uint32
	Threds  uint8
	KeyLen  uint32
	SaltLen uint32
}

func NewHashConfig(memory, time uint32, threads uint8, keylen, saltlen uint32) *HashConfig {
	return &HashConfig{
		Memory:  memory,
		Time:    time,
		Threds:  threads,
		KeyLen:  keylen,
		SaltLen: saltlen,
	}
}

func (h *HashConfig) OwaspRecomendedHashConfig() {
	// owasp min recomendation (2023 may)
	h.Memory = 32 * 1024
	h.Time = 2
	h.Threds = 1
	h.KeyLen = 32
	h.SaltLen = 16
}

func (h *HashConfig) genSalt() []byte {
	salt := make([]byte, h.SaltLen)
	rand.Read(salt)
	return salt
}

func (h *HashConfig) Hash(pwd string) string {
	salt := h.genSalt()
	hash := argon2.IDKey([]byte(pwd), salt, h.Time, h.Memory, h.Threds, h.KeyLen)

	// format hash
	// $argon2id$v=$m=,t=,p=$salt$hash
	version := argon2.Version
	encodedSalt := base64.RawStdEncoding.EncodeToString(salt)
	encodedHash := base64.RawStdEncoding.EncodeToString(hash)
	out := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", version, h.Memory, h.Time, h.Threds, encodedSalt, encodedHash)
	return out
}

func (h *HashConfig) Compare(pwd string, hashedPwd string) error {
	// deconstructure hash
	splittedHash := strings.Split(hashedPwd, "$")

	// cek panjang
	log.Println(splittedHash, len(splittedHash))
	log.Println("masuk invalid hash")
	if len(splittedHash) != 6 {
		return errors.New("invalid Hash")
	}

	// cek argon2id
	log.Println("masuk argon2id salah")
	if splittedHash[1] != "argon2id" {
		return errors.New("not an argon2id hash")
	}
	log.Println("wrong ssca")
	var version int
	if _, err := fmt.Sscanf(splittedHash[2], "v=%d", &version); err != nil {
		return errors.New("wrong sscanf syntax")
	}

	log.Println("wrong versi argon")
	if version != argon2.Version {
		return errors.New("wrong argon2id version")
	}

	var memory, time uint32
	var threads uint8
	log.Println("wrong sscan")
	if _, err := fmt.Sscanf(splittedHash[3], "m=%d,t=%d,p=%d", &memory, &time, &threads); err != nil {
		return errors.New("wrong sscanf syntax")
	}

	// ambil salt dan hash
	log.Println("wrong decode salt")
	salt, err := base64.RawStdEncoding.DecodeString(splittedHash[4])
	if err != nil {
		return errors.New("failed decode salt")
	}
	log.Println("wrong decode hash")
	hash, err := base64.RawStdEncoding.DecodeString(splittedHash[5])
	if err != nil {
		return errors.New("failed decode hash")
	}
	// generate hash from incoming password (input user passwd)
	newHash := argon2.IDKey([]byte(pwd), salt, time, memory, threads, uint32(len(hash)))

	log.Println("wrong decode pw")
	if subtle.ConstantTimeCompare(hash, newHash) == 0 {
		return errors.New("wrong password")
	}
	return nil
}
