package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"html/template"
	"jwt_case_1/templates"
	"log"
	"math/big"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey

	verifyKeyPEM []byte // zafiyet yaratan değişken

	// MockUsers
	users = []User{
		{ID: 101, Username: "admin", Password: "admin123", Role: "admin"},
		{ID: 202, Username: "user", Password: "user123", Role: "user"},
	}

	CTF_FLAG = "CTF{G0lang_JWK_RS256_M4st3r}"
)

type User struct {
	ID       int
	Username string
	Password string
	Role     string
}

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Claims struct {
	Username             string `json:"username"`
	AccountRole          string `json:"role"`
	ID                   int    `json:"id"`
	jwt.RegisteredClaims        // iat, exp, iss gibi standart jwt payload alanları
}

type JWK struct {
	Kty string `json:"kty"` //Key Type
	Alg string `json:"alg"` //Algorithm => RS256
	//	Kid string `json:"kid"` //Key ID
	//	Use string `json:"use"` // Public key use (sig)
	N string `json:"n"` // Modulus
	E string `json:"e"` // Exponent
}

type JWKS struct {
	Keys []JWK `json:"keys"`
}

type PageData struct {
	Error string
}

type AdminPageData struct {
	Username string
	Flag     string
}

// contextKey is a custom type for context keys to avoid collisions with built-in types.
type contextKey string

const claimsContextKey contextKey = "claims"

// init() fonksiyonu zaten golang'de main'den önce çalıştırılan özel bir fonksiyon, o yüzden bunu main içinde çağırmana gerek yok.
func init() {
	// Uygulama açılırken RSA anahtarları oluşturulur.
	var err error
	signKey, err = rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		log.Fatal("RSA key generation error: ", err)
	}
	if signKey == nil {
		log.Fatal("Signkey oluşturulamadı...")
	}
	verifyKey = &signKey.PublicKey
	//fmt.Println("RSA anahtarları başarılı bir şekilde oluşturuldu...")

	// Zafiyet oluşturacak kısım.
	pubBytes, err := x509.MarshalPKIXPublicKey(verifyKey)
	if err != nil {
		log.Fatal(err)
	}
	verifyKeyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubBytes,
	})
}

// Login page
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	/*
		var creds Credentials
		if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		} */

	// method get ise direkt olarak login page'i göster.

	cookie, err := r.Cookie("access_token")
	if err == nil {
		// Cookie varsa token'ı doğrula
		token, err := jwt.ParseWithClaims(cookie.Value, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			return verifyKey, nil
		})

		// Token hem hatasız parse edildi hem de süresi dolmamışsa (Valid)
		if err == nil && token.Valid {
			// Kullanıcıyı direkt profil sayfasına yönlendir
			http.Redirect(w, r, "/profile", http.StatusSeeOther)
			return
		}
	}

	if r.Method == http.MethodGet {
		tmpl := template.Must(template.New("login").Parse(templates.LoginPageHTML))
		tmpl.Execute(w, nil)
		return
	}

	if r.Method != http.MethodPost {
		return
	}

	username := r.FormValue("username")
	password := r.FormValue("password")

	// Kullanıcı doğrulama
	var user *User
	for _, u := range users {
		if u.Username == username && u.Password == password {
			user = &u
			break
		}
	}

	// Kullanıcı bulunamadı.
	if user == nil {
		tmpl := template.Must(template.New("login").Parse(templates.LoginPageHTML))
		tmpl.Execute(w, PageData{Error: "Username or Password inccorect"})
		return
	}

	//Kullanıcı bulduktan sonra JWT token oluşturulup user'a atıyoruz.

	claims := Claims{
		Username:    user.Username,
		AccountRole: user.Role,
		ID:          user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)

	delete(token.Header, "kid") // Eğer "kid" alanı varsa silsin.

	nStr := base64.RawURLEncoding.EncodeToString(verifyKey.N.Bytes())
	eStr := base64.RawURLEncoding.EncodeToString(big.NewInt(int64(verifyKey.E)).Bytes())

	//jwk oluşturup Header'a ekliyoruz.
	token.Header["jwk"] = map[string]string{
		"kty": "RSA",
		"alg": "RS256",
		"n":   nStr,
		"e":   eStr,
	}

	// token sign ediyoruz.
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name:     "access_token",
		Value:    tokenString,
		Expires:  time.Now().Add(1 * time.Hour),
		HttpOnly: true,
		Path:     "/",
	})

	http.Redirect(w, r, "/profile", http.StatusSeeOther)
}

// middleware => Her gelen istekte token doğrulama
func AuthMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := ""

		// 1. Önce Cookie'ye bak
		cookie, err := r.Cookie("access_token")
		if err == nil {
			tokenString = cookie.Value
		}

		// 2. Cookie yoksa Header'a bak (Postman/Curl desteği için)
		if tokenString == "" {
			authHeader := r.Header.Get("Authorization")
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				tokenString = authHeader[7:]
			}
		}

		if tokenString == "" {
			// Web tarayıcısı ise login sayfasına at, değilse 401 ver
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		fmt.Println("Tokenstring => ", tokenString)

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
			/*
				if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
					return nil, fmt.Errorf("unexpected method")
				}
			*/

			if _, ok := token.Method.(*jwt.SigningMethodRSA); ok {
				return verifyKey, nil
			}

			// zafiyetli kısım
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); ok {
				fmt.Println("HMAC algoritması")
				return verifyKeyPEM, nil
			}

			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		})
		fmt.Println("token => ", token.Valid)
		if err != nil || !token.Valid {
			// Token geçersizse cookie'yi silip logine at
			http.SetCookie(w, &http.Cookie{Name: "access_token", MaxAge: -1})
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}
		fmt.Println("middleware'deki claims:", claims)
		ctx := context.WithValue(r.Context(), claimsContextKey, claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func ProfileHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsContextKey).(*Claims)
	fmt.Println(r.Context().Value("claims"))
	if !ok {
		http.Redirect(w, r, "/login", http.StatusPermanentRedirect)
		return
	}

	tmplt := template.Must(template.New("profile").Parse(templates.ProfilePageHTML))
	tmplt.Execute(w, claims)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:    "access_token",
		Value:   "",
		Expires: time.Unix(0, 0),
		Path:    "/",
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func AdminHandler(w http.ResponseWriter, r *http.Request) {
	claims, ok := r.Context().Value(claimsContextKey).(*Claims)

	if !ok {
		http.Error(w, "Session error...", http.StatusUnauthorized)
		return
	}

	if claims.AccountRole != "admin" {
		http.Error(w, "Unauthorized...", http.StatusUnauthorized)
		return
	}

	data := AdminPageData{
		Username: claims.Username,
		Flag:     CTF_FLAG,
	}

	tmpl := template.Must(template.New("admin").Parse(templates.AdminPageHTML))
	err := tmpl.Execute(w, data)

	if err != nil {
		http.Error(w, "Page rendering error.", http.StatusBadGateway)
	}
}

func main() {

	http.HandleFunc("/login", LoginHandler)

	http.HandleFunc("/admin", AuthMiddleware(AdminHandler))

	http.HandleFunc("/profile", AuthMiddleware(ProfileHandler))

	http.HandleFunc("/logout", AuthMiddleware(LogoutHandler))

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
	})

	fmt.Println("Server running on => 8080")
	http.ListenAndServe(":8080", nil)
}
