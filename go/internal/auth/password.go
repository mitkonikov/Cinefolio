package auth

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"

	"cinefolio/internal/store"
)

const (
	memoryKiB      = 64 * 1024
	iterations     = 3
	parallelism    = 2
	hashBytes      = 32
	saltBytes      = 16
	passwordParams = "m=65536,t=3,p=2"
)

var ErrHashBusy = errors.New("password service is busy; try again shortly")

type hashAdmission struct {
	active chan struct{}
	total  chan struct{}
}

var hashes = hashAdmission{active: make(chan struct{}, 2), total: make(chan struct{}, 10)}

func (admission hashAdmission) acquire(ctx context.Context) (func(), error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	select {
	case admission.total <- struct{}{}:
	default:
		return nil, ErrHashBusy
	}
	select {
	case admission.active <- struct{}{}:
		release := func() {
			<-admission.active
			<-admission.total
		}
		if err := ctx.Err(); err != nil {
			release()
			return nil, err
		}
		return release, nil
	case <-ctx.Done():
		<-admission.total
		return nil, ctx.Err()
	}
}

func validatePassword(password string) error {
	if len(password) < 8 || len(password) > 1024 {
		return store.ValidationError("password must be between 8 and 1024 bytes")
	}
	return nil
}

func derive(ctx context.Context, password string, salt []byte) ([]byte, error) {
	release, err := hashes.acquire(ctx)
	if err != nil {
		return nil, err
	}
	defer release()
	digest := argon2.IDKey([]byte(password), salt, iterations, memoryKiB, parallelism, hashBytes)
	return digest, ctx.Err()
}

func hashPassword(ctx context.Context, password string) (string, error) {
	salt := make([]byte, saltBytes)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}
	digest, err := derive(ctx, password, salt)
	if err != nil {
		return "", err
	}
	b64 := base64.RawStdEncoding
	return fmt.Sprintf("$argon2id$v=%d$%s$%s$%s", argon2.Version, passwordParams,
		b64.EncodeToString(salt), b64.EncodeToString(digest)), nil
}

func verifyPassword(ctx context.Context, encoded, password string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 6 || parts[1] != "argon2id" || parts[2] != "v=19" || parts[3] != passwordParams {
		return false, nil
	}
	b64 := base64.RawStdEncoding
	salt, saltErr := b64.DecodeString(parts[4])
	digest, digestErr := b64.DecodeString(parts[5])
	if saltErr != nil || digestErr != nil || len(salt) != saltBytes || len(digest) != hashBytes {
		return false, nil
	}
	candidate, err := derive(ctx, password, salt)
	return subtle.ConstantTimeCompare(digest, candidate) == 1, err
}
