package service

import (
	"crypto/sha256"
	"errors"
	"films/internal/models"
	"fmt"
	"github.com/dgrijalva/jwt-go/v4"
	"time"
)

type Email interface {
	SendEmail(email string, body []byte) error
}

type AuthClaims struct {
	jwt.StandardClaims
	User *models.User `json:"user"`
}

type AuthService struct {
	emailClient    Email
	userStorage    UserStorage
	hashSalt       string
	signingKey     []byte
	expireDuration time.Duration
}

func NewAuthService(
	emailClient Email,
	userStorage UserStorage,
	hashSalt string,
	signingKey []byte,
	tokenTTLSeconds time.Duration) *AuthService {
	return &AuthService{
		emailClient:    emailClient,
		userStorage:    userStorage,
		hashSalt:       hashSalt,
		signingKey:     signingKey,
		expireDuration: time.Second * tokenTTLSeconds,
	}
}

func (a *AuthService) SignUp(user *models.User) error {
	// hashing password
	pwd := sha256.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))

	// hashing user data
	hash := sha256.New()
	hash.Write([]byte(user.FullName))
	hash.Write([]byte(user.Phone))
	hash.Write([]byte(user.Email))
	hash.Write([]byte(user.Password))

	// saving hash
	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))
	user.Hash = fmt.Sprintf("%x", hash.Sum(nil))

	// inserting data to db
	err := a.userStorage.Insert(user)
	if err != nil {
		return err
	}
	return nil
}

func (a *AuthService) SignIn(user *models.User) (string, error) {
	pwd := sha256.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))

	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	user, err := a.userStorage.GetByCredentials(user.Email, user.Password)
	if err != nil {
		return "", errors.New("user not found")
	}

	claims := AuthClaims{
		User: user,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: jwt.At(time.Now().Add(a.expireDuration)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(a.signingKey)
}

func (a *AuthService) ParseToken(accessToken string) (*models.User, error) {
	token, err := jwt.ParseWithClaims(accessToken, &AuthClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.signingKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AuthClaims); ok && token.Valid {
		return claims.User, nil
	}

	return nil, errors.New("not valid token")
}

func (a *AuthService) SendVerificationEmail(user *models.User) error {
	msgBody := `http://localhost:8000/verify?email=` + user.Email + `&hash=` + user.Hash
	err := a.emailClient.SendEmail(user.Email, []byte(msgBody + "\n"))
	if err != nil {
		return err
	}

	return nil
}

func (a *AuthService) VerifyEmail(email, hash string) error {
	userHash, err := a.userStorage.GetUserHash(email)
	if userHash == "" {
		return errors.New("there is no user with this email")
	}
	if err != nil {
		return err
	}

	if userHash == hash {
		err = a.userStorage.VerifyUser(email)
		if err != nil {
			return errors.New("something went wrong")
		}
	} else {
		return errors.New("hashes is not equal")
	}

	return nil
}

func (a *AuthService) ResetPassword(user *models.User) error {
	oldPwd := sha256.New()
	oldPwd.Write([]byte(user.OldPassword))
	oldPwd.Write([]byte(a.hashSalt))

	userHash, err := a.userStorage.GetPasswordHash(fmt.Sprintf("%x", oldPwd.Sum(nil)))
	if userHash == "" {
		return errors.New("your old password is incorrect")
	}
	if err != nil {
		return err
	}

	pwd := sha256.New()
	pwd.Write([]byte(user.Password))
	pwd.Write([]byte(a.hashSalt))

	user.Password = fmt.Sprintf("%x", pwd.Sum(nil))

	err = a.userStorage.UpdateUserPassword(user)
	if err != nil {
		return err
	}

	return nil
}
