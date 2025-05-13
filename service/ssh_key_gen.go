package service

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
	"os/exec"
	"strings"
)

func SshKeygen() (privateKeyPEM, publicKeyOpenSSH string) {
	// RSA 2048 비트 개인 키 생성
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", ""
	}

	// PEM 포맷으로 변환 (PKCS#1)
	privPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// 임시 파일에 저장 (OpenSSH 변환을 위해)
	tmpPrivFile := "temp_rsa_key.pem"
	err = os.WriteFile(tmpPrivFile, privPem, 0600)
	if err != nil {
		return "", ""
	}

	// OpenSSH 포맷으로 변환 ('-m PEM' 제거)
	cmd := exec.Command("ssh-keygen", "-p", "-f", tmpPrivFile, "-N", "")
	err = cmd.Run()
	if err != nil {
		return "", ""
	}

	// 변환된 OpenSSH Private Key 읽기
	privKeyBytes, err := os.ReadFile(tmpPrivFile)
	if err != nil {
		return "", ""
	}
	privateKeyPEM = string(privKeyBytes)

	// OpenSSH Public Key 생성
	cmd = exec.Command("ssh-keygen", "-y", "-f", tmpPrivFile)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", ""
	}
	publicKeyOpenSSH = strings.TrimSpace(out.String())

	// 임시 파일 삭제
	_ = os.Remove(tmpPrivFile)

	return privateKeyPEM, publicKeyOpenSSH
}
