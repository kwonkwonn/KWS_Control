package WorkerCont

// import (
// 	"fmt"
// 	"os"
// 	"path/filepath"
// )

// func SshPrivateStore(UUID string, PrivateKey string) {
// 	workingDir, err := os.Getwd()
// 	if err != nil {
// 		fmt.Println("Error getting working directory:", err)
// 		return
// 	}

// 	sshPrivateDir := filepath.Join(workingDir, "sshprivate")

// 	// 폴더가 없으면 생성
// 	if _, err := os.Stat(sshPrivateDir); os.IsNotExist(err) {
// 		err = os.Mkdir(sshPrivateDir, 0700) // 0700: 소유자만 접근 가능
// 		if err != nil {
// 			fmt.Println("Error creating sshprivate directory:", err)
// 			return
// 		}
// 		fmt.Println("Created directory:", sshPrivateDir)
// 	}
// 	privateKeyPath := filepath.Join(sshPrivateDir, "ssh_"+UUID)

// 	err = os.WriteFile(privateKeyPath, []byte(PrivateKey), 0600) // 0600: 소유자만 읽고 쓸 수 있음
// 	if err != nil {
// 		fmt.Println("Error writing private key file:", err)
// 		return
// 	}

// 	fmt.Println("SSH private key saved at:", privateKeyPath)

// }
