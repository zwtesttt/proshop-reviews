package auth

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// 获取token
func AuthenticateUser(host, email, password string) (string, error) {
	url := host + "/api/users/auth"

	// 构造请求体
	data := map[string]string{
		"email":    email,
		"password": password,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	// 创建 HTTP 请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// 获取 Set-Cookie 字段
	setCookie := resp.Header.Get("Set-Cookie")
	if setCookie == "" {
		return "", fmt.Errorf("Set-Cookie not found in response")
	}
	jwt, err := ExtractJWTFromCookie(setCookie)
	if err != nil {
		return "", err
	}

	return jwt, nil
}

func ExtractJWTFromCookie(setCookie string) (string, error) {
	// 在 Set-Cookie 字段中查找 "jwt=" 的起始位置
	start := strings.Index(setCookie, "jwt=")
	if start == -1 {
		return "", fmt.Errorf("JWT token not found in Set-Cookie")
	}

	// 从 "jwt=" 的起始位置开始查找分号 ";" 的位置
	end := strings.Index(setCookie[start:], ";")
	if end == -1 {
		// 如果没有找到分号，则直接从 "jwt=" 的位置截取到字符串末尾
		return setCookie[start+len("jwt="):], nil
	}

	// 找到分号后，从 "jwt=" 的位置截取到分号的位置
	return setCookie[start+len("jwt=") : start+end], nil
}
