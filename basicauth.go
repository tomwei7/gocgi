package gocgi

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/binary"
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	SessionCookieKey = "gocgi-session"
)

func (h *BasicAuthHandler) encryptSession(pt []byte) []byte {
	p := len(pt) % h.sessionCipher.BlockSize()
	if p > 0 {
		b := make([]byte, h.sessionCipher.BlockSize()-p)
		pt = append(pt, b...)
	}

	ct := make([]byte, len(pt))
	bs := h.sessionCipher.BlockSize()
	for i := 0; i < len(pt); i += bs {
		h.sessionCipher.Encrypt(ct[i:], pt[i:])
	}
	return ct
}

func (h *BasicAuthHandler) decryptSession(ct []byte) []byte {
	if len(ct)%h.sessionCipher.BlockSize() != 0 {
		return nil
	}

	pt := make([]byte, len(ct))
	bs := h.sessionCipher.BlockSize()
	for i := 0; i < len(pt); i += bs {
		h.sessionCipher.Decrypt(pt[i:], ct[i:])
	}

	i := bytes.IndexByte(pt, 0)
	if i > 0 {
		pt = pt[:i]
	}
	return pt
}

type BasicAuthHandler struct {
	userInfos     []*url.Userinfo
	handler       http.Handler
	sessionCipher cipher.Block
}

func (h *BasicAuthHandler) CheckCookie(r *http.Request) bool {
	cookie, err := r.Cookie(SessionCookieKey)
	if err != nil {
		return false
	}

	ct, err := base64.StdEncoding.DecodeString(cookie.Value)
	if err != nil {
		return false
	}

	pt := h.decryptSession(ct)
	return h.CheckUsername(string(pt))
}

func (h *BasicAuthHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if h.CheckCookie(r) {
		h.handler.ServeHTTP(w, r)
		return
	}

	username, password, ok := r.BasicAuth()
	if !ok || !h.Check(username, password) {
		w.Header().Set("WWW-Authenticate", fmt.Sprintf(`Basic realm="%s"`, r.URL.Path))
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("401 Unauthorized\n"))
		return
	}
	ct := h.encryptSession([]byte(username))
	sessionValue := base64.StdEncoding.EncodeToString(ct)
	http.SetCookie(w, &http.Cookie{Name: SessionCookieKey, Value: sessionValue, Path: "/"})
	h.handler.ServeHTTP(w, r)
}

func (h *BasicAuthHandler) CheckUsername(username string) bool {
	for _, userInfo := range h.userInfos {
		if userInfo.Username() == username {
			return true
		}
	}
	return false
}

func (h *BasicAuthHandler) Check(username, passwd string) bool {
	for _, userInfo := range h.userInfos {
		if userInfo.Username() != username {
			continue
		}
		if expectPasswd, _ := userInfo.Password(); passwd == expectPasswd {
			return true
		}
	}
	return false
}

func parseUserInfos(users []string) []*url.Userinfo {
	var userInfos []*url.Userinfo
	for _, pair := range users {
		ss := strings.SplitN(pair, ":", 2)
		if len(ss) == 2 {
			userInfos = append(userInfos, url.UserPassword(ss[0], ss[1]))
		} else {
			userInfos = append(userInfos, url.User(ss[0]))
		}
	}
	return userInfos
}

func WithBasicAuth(users []string, handler http.Handler) http.Handler {
	if len(users) == 0 {
		return handler
	}

	userInfos := parseUserInfos(users)

	var key [16]byte
	source := rand.NewSource(time.Now().UnixNano())
	binary.BigEndian.PutUint64(key[:], uint64(source.Int63()))
	binary.BigEndian.PutUint64(key[8:], uint64(source.Int63()))

	sessionCipher, err := aes.NewCipher(key[:])
	if err != nil {
		panic(err)
	}

	return &BasicAuthHandler{userInfos: userInfos, handler: handler, sessionCipher: sessionCipher}
}
