package controllers

import (
	"strconv"
	"time"

	"auth-api/models"
	"auth-api/database"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
)

type Claims struct {
	jwt.StandardClaims
}

func Register(c *fiber.Ctx) error {
	var data map[string]string
 
	// リクエストデータをパースする
	if err := c.BodyParser(&data); err != nil {
		return err
	}

	// パスワードチェック
	if data["password"] != data["password_confirm"] {
		c.Status(400)
		return c.JSON(fiber.Map{
			"message": "Passwords do not match!",
		})
	}

	// パスワードをエンコード
	password, _ := bcrypt.GenerateFromPassword([]byte(data["password"]), 14)
 
	user := models.User{
		FirstName: data["first_name"],
		LastName:  data["last_name"],
		Email:     data["email"],
		Password:  password,
	}

	database.DB.Create(&user)
 
	return c.JSON(user)
}

// ログイン
func Login(c *fiber.Ctx) error {
	var data map[string]string
 
	if err := c.BodyParser(&data); err != nil {
		return err
	}
 
	var user models.User
	// emailに紐づくユーザーを取得
	// &userを指定することでDBから取得したデータを直接格納できる
	database.DB.Where("email = ?", data["email"]).First(&user)

	if user.ID == 0 {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "User not found",
		})
	}

	// パスワードのチェック
	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(data["password"])); err != nil {
		c.Status(404)
		return c.JSON(fiber.Map{
			"message": "Invalid password",
		})
	}
 
	// JWT
	claims := jwt.StandardClaims{
		Issuer:        strconv.Itoa(int(user.ID)),			// stringに型変換
		ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),	// 有効期限
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := jwtToken.SignedString([]byte("secret"))
	if err != nil {
		return c.SendStatus(fiber.StatusInternalServerError)
	}

	// Cookie
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    token,
		Expires:  time.Now().Add(time.Hour * 24),
		HTTPOnly: true,
	}
	c.Cookie(&cookie)
 
	return c.JSON(fiber.Map{
		"jwt": token,
	})
}

func User(c *fiber.Ctx) error {
	// CookieからJWTを取得
	cookie := c.Cookies("jwt") // Loginで保存したもの
	// token取得
	token, err := jwt.ParseWithClaims(cookie, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})
	if err != nil || !token.Valid {
		c.Status(fiber.StatusUnauthorized)
		return c.JSON(fiber.Map{
			"message": "unauthenticated",
		})
	}
 
	claims := token.Claims.(*Claims)
	// User IDを取得
	id := claims.Issuer
 
	var user models.User
	database.DB.Where("id = ?", id).First(&user)
 
	return c.JSON(user)
}

func Logout(c *fiber.Ctx) error {
	cookie := fiber.Cookie{
		Name:     "jwt",
		Value:    "",                         // tokenを空にする
		Expires:  time.Now().Add(-time.Hour), // マイナス値を入れて期限切れ
		HTTPOnly: true,
	}
 
	c.Cookie(&cookie)
	return c.JSON(fiber.Map{
		"message": "success",
	})
}