package utils

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
)

// PKCEGenerator 默认PKCE生成器实现
type PKCEGenerator struct{}

// NewPKCEGenerator 创建一个新的PKCE生成器实例
func NewPKCEGenerator() *PKCEGenerator {
	return &PKCEGenerator{}
}

// GenerateCodeVerifier 生成一个随机的code_verifier
func (g *PKCEGenerator) GenerateCodeVerifier() (string, error) {
	// 生成32字节的随机数据
	data := make([]byte, 32)
	_, err := rand.Read(data)
	if err != nil {
		return "", err
	}

	// 使用URL安全的base64编码
	verifier := base64.RawURLEncoding.EncodeToString(data)
	return verifier, nil
}

// GenerateCodeChallenge 根据code_verifier生成code_challenge
// 使用SHA256方法计算挑战值
func (g *PKCEGenerator) GenerateCodeChallenge(codeVerifier string) string {
	hash := sha256.Sum256([]byte(codeVerifier))
	// 使用URL安全的base64编码，不包含填充字符
	challenge := base64.RawURLEncoding.EncodeToString(hash[:])
	return challenge
}

// GetCodeChallengeMethod 返回使用的挑战方法
func (g *PKCEGenerator) GetCodeChallengeMethod() string {
	return "S256"
}
