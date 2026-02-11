package services

import (
	"fmt"
	"io"
	"os/exec"
)

func CreateFtpUser(ftpId, ftpPw string, quotaMB int64) error {
	args := []string{
		"exec", "-i", "ftp",
		"pure-pw", "useradd", ftpId,
		"-f", "/etc/pure-ftpd/passwd",
		"-m",
		"-u", "33", "-g", "33",
		"-d", fmt.Sprintf("/home/%s", ftpId),
		"-N", fmt.Sprintf("%d", quotaMB),
	}
	cmd := exec.Command("docker", args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("stdin pipe error: %v", err)
	}

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("command start error: %v", err)
	}

	_, _ = io.WriteString(stdin, ftpPw+"\n")
	_, _ = io.WriteString(stdin, ftpPw+"\n")
	_ = stdin.Close()

	if err := cmd.Wait(); err != nil {
		return fmt.Errorf("FTP 유저 생성 실패 (명령어 에러): %v", err)
	}

	return nil
}
