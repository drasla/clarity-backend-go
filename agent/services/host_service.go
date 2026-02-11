package services

import (
	"agent/utils"
	"fmt"
	"os"
)

func CreateHosting(ftpId, ftpPw string, quotaMB int64) error {
	baseDir := fmt.Sprintf("/home/%s", ftpId)

	if err := os.MkdirAll(baseDir, 0755); err != nil {
		return fmt.Errorf("디렉토리 생성 실패: %v", err)
	}

	if err := utils.RunCommand("chown", "-R", "33:33", baseDir); err != nil {
		return fmt.Errorf("권한 설정 실패: %v", err)
	}

	if err := CreateFtpUser(ftpId, ftpPw, quotaMB); err != nil {
		return err
	}

	return nil
}
