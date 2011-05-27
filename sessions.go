package mango

import (
	"bytes"
	"hash"
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"gob"
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

func decode64(value string) (result string) {
	buffer := bytes.NewBufferString(value)
	encoder := base64.NewDecoder(base64.StdEncoding, buffer)
	decoded, _ := ioutil.ReadAll(encoder)
	return string(decoded)
}

func decodeCookie(value, secret string) (cookie map[string]interface{}) {
	cookie = make(map[string]interface{})

	split := strings.Split(string(value), "--", 2)

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
	buffer := &bytes.Buffer{}
	encoder := gob.NewEncoder(buffer)
	encoder.Encode(value)
	return buffer.String()
}

func encode64(value string) (result string) {
	buffer := &bytes.Buffer{}
	encoder := base64.NewEncoder(base64.StdEncoding, buffer)
	encoder.Write([]byte(value))
	encoder.Close()
	return buffer.String()
}

func encodeCookie(value map[string]interface{}, secret string) (cookie string) {
	data := encodeGob(value)

	return fmt.Sprintf("%s--%s", encode64(data), encode64(hashCookie(data, secret)))
}

func prepareSession(env Env, key, secret string) {
	for _, value := range env.Request().Cookie {
		if value.Name == key {
			env["mango.session"] = decodeCookie(value.Value, secret)
			return
		}
	}

	// Didn't find a session to decode
	env["mango.session"] = make(map[string]interface{})
}

func commitSession(headers Headers, env Env, key, secret, domain string) {
	headers.Add("Set-Cookie", fmt.Sprintf("%s=%s; Domain=%s;", key, encodeCookie(env["mango.session"].(map[string]interface{}), secret), domain))
}

func Sessions(secret, key, domain string) Middleware {
	return func(env Env, app App) (status Status, headers Headers, body Body) {
		prepareSession(env, key, secret)
		status, headers, body = app(env)
		commitSession(headers, env, key, secret, domain)
		return status, headers, body
	}
}
