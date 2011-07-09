package mango

import (
	"bytes"
	"hash"
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"gob"
	"http"
	"strings"
)

func hashCookie(data, secret string) (sum string) {
	var h hash.Hash = hmac.NewSHA1([]byte(secret))
	h.Write([]byte(data))
	return string(h.Sum())
}

func verifyCookie(data, secret, sum string) bool {
	return hashCookie(data, secret) == sum
}

func decodeGob(value string) (result map[string]interface{}) {
	buffer := bytes.NewBufferString(value)
	decoder := gob.NewDecoder(buffer)
	result = make(map[string]interface{})
	decoder.Decode(&result)
	return result
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

	split := strings.Split(string(value), "--")

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

func encodeGob(value interface{}) (result string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	encoder.Encode(value)
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

	return fmt.Sprintf("%s--%s", encode64(data), encode64(hashCookie(data, secret)))
}

func prepareSession(env Env, key, secret string) {
	for _, cookie := range env.Request().Cookies() {
		if cookie.Name == key {
			env["mango.session"] = decodeCookie(cookie.Value, secret)
			return
		}
	}

	// Didn't find a session to decode
	env["mango.session"] = make(map[string]interface{})
}

func commitSession(headers Headers, env Env, key, secret, domain string) {
	cookie := new(http.Cookie)
	cookie.Name = key
	cookie.Value = encodeCookie(env["mango.session"].(map[string]interface{}), secret)
	cookie.Domain = domain
	headers.Add("Set-Cookie", cookie.String())
}

func Sessions(secret, key, domain string) Middleware {
	return func(env Env, app App) (status Status, headers Headers, body Body) {
		prepareSession(env, key, secret)
		status, headers, body = app(env)
		commitSession(headers, env, key, secret, domain)
		return status, headers, body
	}
}
