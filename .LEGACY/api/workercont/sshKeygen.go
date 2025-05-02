package WorkerCont

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func SshKeygen() (privateKeyPEM, publicKeyOpenSSH string, err error) {
	// 1️⃣ RSA 2048 비트 개인 키 생성
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", "", fmt.Errorf("🔴 키 생성 실패: %v", err)
	}

	// 2️⃣ PEM 포맷으로 변환 (PKCS#1)
	privPem := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	})

	// 3️⃣ 임시 파일에 저장 (OpenSSH 변환을 위해)
	tmpPrivFile := "temp_rsa_key.pem"
	err = os.WriteFile(tmpPrivFile, privPem, 0600)
	if err != nil {
		return "", "", fmt.Errorf("🔴 임시 파일 저장 실패: %v", err)
	}

	// 🚀 4️⃣ OpenSSH 포맷으로 변환 (💡 '-m PEM' 제거)
	cmd := exec.Command("ssh-keygen", "-p", "-f", tmpPrivFile, "-N", "")
	err = cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("🔴 OpenSSH 변환 실패: %v", err)
	}

	// 5️⃣ 변환된 OpenSSH Private Key 읽기
	privKeyBytes, err := os.ReadFile(tmpPrivFile)
	if err != nil {
		return "", "", fmt.Errorf("🔴 변환된 키 읽기 실패: %v", err)
	}
	privateKeyPEM = string(privKeyBytes)

	// 6️⃣ OpenSSH Public Key 생성
	cmd = exec.Command("ssh-keygen", "-y", "-f", tmpPrivFile)
	var out bytes.Buffer
	cmd.Stdout = &out
	err = cmd.Run()
	if err != nil {
		return "", "", fmt.Errorf("🔴 공개 키 변환 실패: %v", err)
	}
	publicKeyOpenSSH = strings.TrimSpace(out.String())

	// 7️⃣ 임시 파일 삭제
	_ = os.Remove(tmpPrivFile)

	fmt.Println("✅ OpenSSH Private Key:\n", privateKeyPEM)
	fmt.Println("✅ OpenSSH Public Key:\n", publicKeyOpenSSH)

	return privateKeyPEM, publicKeyOpenSSH, nil
}
