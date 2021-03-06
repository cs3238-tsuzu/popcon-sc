package main

import (
	"fmt"
	"net/http"
	"strconv"
	"syscall"
	"time"

	"github.com/cs3238-tsuzu/popcon-sc/lib/utility"

	"math/rand"

	"io"

	"golang.org/x/crypto/ssh/terminal"
)

func ArgumentsToArray(vals ...interface{}) []interface{} {
	return vals
}

func createWrapFormInt64(req *http.Request) func(str string) int64 {
	f := createWrapFormStr(req)

	return func(str string) int64 {
		val, err := strconv.ParseInt(f(str), 10, 64)

		if err != nil {
			return -1
		}

		return val
	}
}

func createWrapFormInt(req *http.Request) func(str string) int {
	f := createWrapFormInt64(req)

	return func(str string) int {
		return int(f(str))
	}
}

func createWrapFormStr(req *http.Request) func(str string) string {
	return func(str string) string {
		arr, has := req.Form[str]
		if has && len(arr) != 0 {
			return arr[0]
		}
		return ""
	}
}
func CreateDefaultAdminUser() bool {
	fmt.Println("No user found in the DB")
	fmt.Println("You need to create the default admin")
	var id, name, email, pass, pass2 string

	fmt.Print("User ID: ")
	_, err := fmt.Scan(&id)

	if len(id) == 0 || err != nil {
		return false
	}

	fmt.Print("User Name: ")
	_, err = fmt.Scan(&name)

	if len(name) == 0 || err != nil {
		return false
	}

	fmt.Print("Email(or \"null\"): ")
	_, err = fmt.Scan(&email)

	if email == "null" {
		email = ""
	}

	if len(email) == 0 || err != nil {
		return false
	}

	fmt.Print("Password (hidden): ")
	passArr, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		return false
	}
	fmt.Println()

	fmt.Print("Password (confirmation): ")
	passArr2, err := terminal.ReadPassword(int(syscall.Stdin))

	if err != nil {
		return false
	}
	fmt.Println()

	pass = string(passArr)
	pass2 = string(passArr2)

	if pass != pass2 {
		fmt.Println("Different password")

		return false
	}

	_, err = mainDB.UserAdd(id, name, pass, utility.NullStringCreate(email), 0, true)

	if err != nil {
		fmt.Println("Failed to create user. (", err.Error(), ")")

		return false
	}

	return true
}

func generateRandomString(strlen int) string {
	rand.Seed(time.Now().UTC().UnixNano())
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, strlen)
	for i := 0; i < strlen; i++ {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

func CreateAdminUserAutomatically() bool {
	pass := generateRandomString(8)
	_, err := mainDB.UserAdd("admin", "admin", pass, NullStringInvalid, 0, true)

	if err == nil {
		DBLog().Infof("Default User(ID: admin, Pass: %s) You should change this", pass)

		return true
	}

	return false
}

type TrimNewlineReader struct {
	reader      io.Reader
	prevNewline bool
}

func NewTrimNewlineReader(reader io.Reader) io.Reader {
	return &TrimNewlineReader{reader, false}
}

func (tnlr *TrimNewlineReader) Read(ret []byte) (int, error) {
	p := make([]byte, len(ret))
	n, err := tnlr.reader.Read(p)

	if err != nil {
		return n, err
	}

	len := 0
	for i := 0; i < n; i++ {
		if p[i] == '\r' {
			tnlr.prevNewline = true
			ret[len] = '\n'
		} else if p[i] == '\n' {
			if tnlr.prevNewline {
				tnlr.prevNewline = false
				continue
			} else {
				ret[len] = '\n'
			}
		} else {
			ret[len] = p[i]
		}
		len++
	}

	return len, nil
}

func SetSession(rw http.ResponseWriter, session string) {
	exp, err := mainRM.SessionExpiration()

	if err != nil {
		DBLog().WithError(err).Error("SessionExpiration() error")
		exp = 60
	}

	cookie := http.Cookie{
		Name:     HTTPCookieSession,
		Value:    session,
		MaxAge:   exp,
		HttpOnly: true,
		Path:     "/",
	}

	http.SetCookie(rw, &cookie)
}

func ProcessUntilError(functions ...func() error) error {
	var err error

	for i := range functions {
		err = functions[i]()

		if err != nil {
			return err
		}
	}

	return nil
}
