package service

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/crypto/argon2"

	"github.com/arasan1289/hexagonal-demo/internal/adapters/config"
	"github.com/arasan1289/hexagonal-demo/internal/adapters/logger"
	"github.com/arasan1289/hexagonal-demo/internal/core/domain"
	"github.com/arasan1289/hexagonal-demo/internal/core/port"
	"github.com/arasan1289/hexagonal-demo/internal/core/util"
)

// UserService struct represents the user service with its dependencies
type UserService struct {
	repo port.IUserRepository // user repository interface
	log  *logger.Logger       // logger instance
}

// NewUserService constructor function
func NewUserService(repo port.IUserRepository, log *logger.Logger) port.IUserService {
	return &UserService{
		repo: repo,
		log:  log,
	}
}

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// Register function: encrypt phone number, generate ID if needed, upsert user
func (us *UserService) Register(ctx context.Context, user *domain.User, config *config.App) (*domain.User, error) {
	phoneNumberEnc, err := util.EncryptString(user.PhoneNumber, config.SecretKey)
	if err != nil {
		return nil, err
	}
	if user.ID == "" {
		user.ID = util.GenerateULID()
	}

	if user.Email != nil {
		emailEnc, err := util.EncryptString(*user.Email, config.SecretKey)
		if err != nil {
			return nil, err
		}
		emailHash := util.HashString(*user.Email)
		user.EmailEncrypted = &emailEnc
		user.EmailHash = &emailHash
	}

	if user.Password != nil {
		passwordHash, err := generatePasswordHash(*user.Password)
		if err != nil {
			return nil, err
		}
		user.Password = &passwordHash
	}
	user.PhoneNumberHash = util.HashString(user.PhoneNumber)
	user.PhoneNumberEncrypted = phoneNumberEnc

	usr, err := us.repo.UpsertUser(ctx, user)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

// GetUser function: retrieve user by ID, decrypt phone number
func (us *UserService) GetUser(ctx context.Context, id string, config *config.App) (*domain.User, error) {
	usr, err := us.repo.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

// GetUserByPhoneNumberOrEmail function: retrieve user by phone number hash or email hash
func (us *UserService) GetUserByPhoneNumberOrEmail(ctx context.Context, str string) (*domain.User, error) {
	hash := util.HashString(str)
	usr, err := us.repo.GetUserByPhoneNumberOrEmail(ctx, hash)
	if err != nil {
		return nil, err
	}
	return usr, nil
}

// GetUserAndComparePassword function: retrieve user by phone number hash or email hash and compare the password
func (us *UserService) GetUserAndComparePassword(ctx context.Context, email, password string) (*domain.User, bool, error) {
	user, err := us.GetUserByPhoneNumberOrEmail(ctx, email)
	if err != nil {
		return nil, false, err
	}
	if user.Password == nil {
		return nil, false, domain.ErrPasswordNotSet
	}

	match, err := comparePasswordAndHash(password, *user.Password)
	if err != nil {
		return nil, false, err
	}

	return user, match, nil
}

func generatePasswordHash(password string) (string, error) {
	// Establish the parameters to use for Argon2.
	p := &params{
		memory:      64 * 1024,
		iterations:  3,
		parallelism: 2,
		saltLength:  16,
		keyLength:   32,
	}

	salt, err := generateRandomBytes(p.saltLength)
	if err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Base64 encode the salt and hashed password.
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Return a string using the standard encoded hash representation.
	encodedHash := fmt.Sprintf("$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s", argon2.Version, p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encodedHash, nil
}

func generateRandomBytes(n uint32) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}

	return b, nil
}

func comparePasswordAndHash(password, encodedHash string) (bool, error) {
	// Extract the parameters, salt and derived key from the encoded password
	// hash.
	p, salt, hash, err := decodeHash(encodedHash)
	if err != nil {
		return false, err
	}

	// Derive the key from the other password using the same parameters.
	otherHash := argon2.IDKey([]byte(password), salt, p.iterations, p.memory, p.parallelism, p.keyLength)

	// Check that the contents of the hashed passwords are identical. Note
	// that we are using the subtle.ConstantTimeCompare() function for this
	// to help prevent timing attacks.
	if subtle.ConstantTimeCompare(hash, otherHash) == 1 {
		return true, nil
	}
	return false, nil
}

func decodeHash(encodedHash string) (p *params, salt, hash []byte, err error) {
	vals := strings.Split(encodedHash, "$")
	if len(vals) != 6 {
		return nil, nil, nil, domain.ErrInvalidHash
	}

	var version int
	_, err = fmt.Sscanf(vals[2], "v=%d", &version)
	if err != nil {
		return nil, nil, nil, err
	}
	if version != argon2.Version {
		return nil, nil, nil, domain.ErrIncompatibleVersion
	}

	p = &params{}
	_, err = fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return nil, nil, nil, err
	}

	salt, err = base64.RawStdEncoding.Strict().DecodeString(vals[4])
	if err != nil {
		return nil, nil, nil, err
	}
	p.saltLength = uint32(len(salt))

	hash, err = base64.RawStdEncoding.Strict().DecodeString(vals[5])
	if err != nil {
		return nil, nil, nil, err
	}
	p.keyLength = uint32(len(hash))

	return p, salt, hash, nil
}
