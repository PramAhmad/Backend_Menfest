package controllers

import (
	"Bemenfest/initializers"
	"Bemenfest/models"
	"errors"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"golang.org/x/crypto/bcrypt"
)

func Register(c *gin.Context) {
	var body struct {
		Uuid     string `json:"id"`
		Email    string `json:"email"`
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Email == "" || body.Username == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email, username and password are required"})
		return
	}

	// uuid
	thisid := uuid.Must(uuid.NewV4())

	// Hash the password before saving it in the database
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash the password"})
		return
	}
	user := models.User{
		Uuid:     thisid.String(),
		Email:    body.Email,
		Username: body.Username,
		Password: string(hashedPassword),
	}

	// query email sudah ada
	mail := initializers.DB.Where("email = ?", body.Email).First(&user)
	if mail.Error == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email sudah ada"})
		return
	}

	result := initializers.DB.Create(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create the user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"result": "silahkan verif emailnya"})
}

func Login(c *gin.Context) {
	var body struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if body.Email == "" || body.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	var user models.User
	result := initializers.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password is invalid"})
		return
	}

	// Compare the stored hashed password, with the hashed version of the password that was received
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password is invalid"})
		return
	}
	// generate token
	token, err := GenerateJWT(1, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}
	// save token
	user.Token = token
	initializers.DB.Save(&user)
	// store token in cookie
	c.SetCookie("token", token, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"result": "Login berhasil"})
}

type JWTClaim struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

var jwtKey = []byte("7EKhqyxC2j-zXmBs0VOTqq-6Kk2lYA3G2bFqoLS3fLTa8zioEyxAP6Xbjv4vyWVVN5pDdRd9QiPkFWk5Lj5WQC")

func GenerateJWT(role int8, username string) (tokenString string, err error) {
	expirationTime := time.Now().Add(1 * time.Hour)
	claims := &JWTClaim{

		Username: username,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err = token.SignedString(jwtKey)
	return
}

func ValidateToken(signedToken string) (err error) {
	token, err := jwt.ParseWithClaims(
		signedToken,
		&JWTClaim{},
		func(token *jwt.Token) (interface{}, error) {
			return []byte(jwtKey), nil
		},
	)
	if err != nil {
		return
	}
	claims, ok := token.Claims.(*JWTClaim)
	if !ok {
		err = errors.New("gak bisa claim")
		return
	}
	if claims.ExpiresAt < time.Now().Local().Unix() {
		err = errors.New("token exp")
		return
	}
	return
}

func Logout(c *gin.Context) {
	c.SetCookie("token", "", -1, "/", "localhost", false, true)
	// delete token in db
	var user models.User
	var body struct {
		Email string `json:"email"`
	}
	if err := c.Bind(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result := initializers.DB.Where("email = ?", body.Email).First(&user)
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email or password is invalid"})
		return
	}
	user.Token = ""
	initializers.DB.Save(&user)
	// saat cookie kosong
	if c.Request.Header.Get("Cookie") == "" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Kamu Sudah logout"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"result": "Logout berhasil"})
}
