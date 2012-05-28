package mango

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"hash"
	"io/ioutil"
	"net/http"
	"sort"
	"strings"
)

type sessionItem struct {
	Key   string
	Value interface{}
}

type sessionItems []sessionItem

func (sis sessionItems) Len() int {
	return len(sis)
}

func (sis sessionItems) Less(i, j int) bool {
	return sis[i].Key < sis[j].Key
}

func (sis sessionItems) Swap(i, j int) {
	sis[i], sis[j] = sis[j], sis[i]
}

func (sis sessionItems) ToMap() (m map[string]interface{}) {
	m = make(map[string]interface{})
	for _, item := range sis {
		m[item.Key] = item.Value
	}
	return
}

func sessionItemsFromMap(m map[string]interface{}) (sis sessionItems) {
	for k, v := range m {
		sis = append(sis, sessionItem{
			Key:   k,
			Value: v,
		})
	}
	sort.Sort(sis)
	return
}

func hashCookie(data, secret string) (sum string) {
	var h hash.Hash = hmac.New(sha1.New, []byte(secret))
	h.Write([]byte(data))
	return string(h.Sum(nil))
}

func verifyCookie(data, secret, sum string) bool {
	return hashCookie(data, secret) == sum
}

func decodeGob(value string) (result map[string]interface{}) {
	buffer := bytes.NewBufferString(value)

	decoder := gob.NewDecoder(buffer)
	sis := sessionItems{}
	decoder.Decode(&sis)

	return sis.ToMap()
}

// Due to a bug in golang where when using
// base64.URLEncoding padding is still added
// (it shouldn't be), we have to strip and add
// it ourselves.
func pad64(value string) (result string) {
	padding := strings.Repeat("=", len(value)%4)
	return strings.Join([]string{value, padding}, "")
}

func decode64(value string) (result string) {
	buffer := bytes.NewBufferString(pad64(value))
	encoder := base64.NewDecoder(base64.URLEncoding, buffer)
	decoded, _ := ioutil.ReadAll(encoder)
	return string(decoded)
}

func decodeCookie(value, secret string) (cookie map[string]interface{}) {
	cookie = make(map[string]interface{})

	split := strings.Split(string(value), "/")

	if len(split) < 2 {
		return cookie
	}

	data := decode64(split[0])
	sum := decode64(split[1])
	if verifyCookie(data, secret, sum) {
		cookie = decodeGob(data)
	}

	return cookie
}

func encodeGob(value map[string]interface{}) (result string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	sis := sessionItemsFromMap(value)
	encoder.Encode(sis)
	return buffer.String()
}

// Due to a bug in golang where when using
// base64.URLEncoding padding is still added
// (it shouldn't be), we have to strip and add
// it ourselves.
func dePad64(value string) (result string) {
	return strings.TrimRight(value, "=")
}

func encode64(value string) (result string) {
	buffer := new(bytes.Buffer)
	encoder := base64.NewEncoder(base64.URLEncoding, buffer)
	encoder.Write([]byte(value))
	encoder.Close()
	return dePad64(buffer.String())
}

func encodeCookie(value map[string]interface{}, secret string) (cookie string) {
	data := encodeGob(value)

	return fmt.Sprintf("%s/%s", encode64(data), encode64(hashCookie(data, secret)))
}

func prepareSession(env Env, key, secret string) {
	value := sessionCookieValue(env, key)
	if value == "" {
		// Didn't find a session to decode
		env["mango.session"] = make(map[string]interface{})
		return
	}
	env["mango.session"] = decodeCookie(value, secret)
}

func commitSession(headers Headers, env Env, key, secret string, newValue string, options *CookieOptions) {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = newValue
	cookie.Path = options.Path
	cookie.Domain = options.Domain
	cookie.MaxAge = options.MaxAge
	cookie.Secure = options.Secure
	cookie.HttpOnly = options.HttpOnly
	headers.Add("Set-Cookie", cookie.String())
}

func sessionCookieValue(env Env, key string) (value string) {
	for _, cookie := range env.Request().Cookies() {
		if cookie.Name == key {
			value = cookie.Value
			return
		}
	}
	return

}

func cookieChanged(env Env, key, secret string) string {
	oldCookieValue := sessionCookieValue(env, key)
	value := env["mango.session"].(map[string]interface{})

	// old and new both are empty
	if oldCookieValue == "" && len(value) == 0 {
		return ""
	}

	newCookieValue := encodeCookie(value, secret)
	if oldCookieValue == newCookieValue {
		return ""
	}
	return newCookieValue
}

type CookieOptions struct {
	Domain   string
	Path     string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

func Sessions(secret, key string, options *CookieOptions) Middleware {
	return func(env Env, app App) (status Status, headers Headers, body Body) {
		prepareSession(env, key, secret)
		status, headers, body = app(env)
		newValue := cookieChanged(env, key, secret)
		if newValue == "" {
			return
		}
		commitSession(headers, env, key, secret, newValue, options)
		return
	}
}
