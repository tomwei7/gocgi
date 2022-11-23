package gocgi

import (
	"strings"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBasicAuthHandler(t *testing.T) {
	handler := WithBasicAuth([]string{"test:password","longlonglonglongusername:password"}, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){ w.Write([]byte("pong")) }))

	t.Run("OK1", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.SetBasicAuth("test", "password")
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, []byte("pong"), w.Body.Bytes())
	})

	t.Run("OK2", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.SetBasicAuth("longlonglonglongusername", "password")
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, []byte("pong"), w.Body.Bytes())
	})

	t.Run("Unauthorized", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		handler.ServeHTTP(w, r)

		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
		assert.Equal(t, []byte("401 Unauthorized\n"), w.Body.Bytes())
	})

	t.Run("CookieAuth1", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.SetBasicAuth("test", "password")
		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		cookie := w.HeaderMap.Get("Set-Cookie")
		if !assert.NotEmpty(t, cookie) {
			return
		}
		//gocgi-session=v7oUWzWP8wIdwR5pJVkKGA==; Path=/
		firstCookieKV := strings.Split(cookie, ";")[0]
		ss := strings.SplitN(firstCookieKV, "=", 2)

		w = httptest.NewRecorder()
		r = httptest.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: ss[0], Value: ss[1]})
		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, []byte("pong"), w.Body.Bytes())
	})

	t.Run("CookieAuth2", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.SetBasicAuth("longlonglonglongusername", "password")
		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		cookie := w.HeaderMap.Get("Set-Cookie")
		if !assert.NotEmpty(t, cookie) {
			return
		}
		//gocgi-session=v7oUWzWP8wIdwR5pJVkKGA==; Path=/
		firstCookieKV := strings.Split(cookie, ";")[0]
		ss := strings.SplitN(firstCookieKV, "=", 2)

		w = httptest.NewRecorder()
		r = httptest.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: ss[0], Value: ss[1]})
		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
		assert.Equal(t, []byte("pong"), w.Body.Bytes())
	})

	t.Run("CookieAuthError1", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: SessionCookieKey, Value: "v7oUWzWP8wIdwR5pJVkKGA=="})
		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})

	t.Run("CookieAuthError2", func(t *testing.T) {
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/", nil)
		r.AddCookie(&http.Cookie{Name: SessionCookieKey, Value: "aGVsbG8sd29ybGQK"})
		handler.ServeHTTP(w, r)
		assert.Equal(t, http.StatusUnauthorized, w.Result().StatusCode)
	})
}
