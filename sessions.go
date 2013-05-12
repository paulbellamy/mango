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

type SessionWrapper struct {
	r             *http.Request
	changed       bool
	key           string
	secret        string
	cookieOptions *CookieOptions
	values        map[string]interface{}
}

func (s *SessionWrapper) Deserialize() {
	value := sessionCookieValue(s.r, s.key)
	if value == "" {
		s.values = make(map[string]interface{})
		return
	}
	s.values = decodeCookie(value, s.secret)
}

func (s *SessionWrapper) Get(name string) interface{} {
	return s.values[name]
}

func (s *SessionWrapper) Set(name string, value interface{}) {
	if s.values[name] != value {
		s.changed = true
		s.values[name] = value
	}
}

func (s *SessionWrapper) Serialize() string {
	return encodeCookie(s.values, s.secret)
}

func (s *SessionWrapper) Write(w http.ResponseWriter) {
	if s.changed {
		newValue := s.Serialize()
		if newValue == "" {
			return
		}
		cookie := new(http.Cookie)
		cookie.Name = s.key
		cookie.Value = newValue
		cookie.Path = s.cookieOptions.Path
		cookie.Domain = s.cookieOptions.Domain
		cookie.MaxAge = s.cookieOptions.MaxAge
		cookie.Secure = s.cookieOptions.Secure
		cookie.HttpOnly = s.cookieOptions.HttpOnly
		w.Header().Add("Set-Cookie", cookie.String())
	}
}

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

func decodeCookie(value, secret string) map[string]interface{} {
	cookie := make(map[string]interface{})

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

func encodeCookie(values map[string]interface{}, secret string) string {
	data := encodeGob(values)

	return fmt.Sprintf("%s/%s", encode64(data), encode64(hashCookie(data, secret)))
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

func sessionCookieValue(r *http.Request, key string) string {
	if cookie, err := r.Cookie(key); err == nil {
		return cookie.Value
	} else {
		return ""
	}
}

type CookieOptions struct {
	Domain   string
	Path     string
	MaxAge   int
	Secure   bool
	HttpOnly bool
}

func Session(r *http.Request, key, secret string, options *CookieOptions) *SessionWrapper {
	wrapper := &SessionWrapper{
		r:             r,
		key:           key,
		secret:        secret,
		cookieOptions: options,
	}
	wrapper.Deserialize()
	return wrapper
}
