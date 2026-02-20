package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"path/filepath"
	"strings"
	"time"
	"tower/pkg/fnEnv"

	"github.com/99designs/gqlgen/graphql"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type FileService interface {
	UploadFile(ctx context.Context, file graphql.Upload, directory string) (string, error)
}

type fileService struct {
	client     *minio.Client
	bucketName string
	domain     string
}

func NewS3Service() (FileService, error) {
	endpoint := fnEnv.GetString("S3_ENDPOINT", "")
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	accessKey := fnEnv.GetString("S3_ACCESS_KEY", "")
	secretKey := fnEnv.GetString("S3_SECRET_KEY", "")
	bucketName := fnEnv.GetString("S3_BUCKET_NAME", "")
	domain := fnEnv.GetString("S3_DOMAIN", "")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: true, // 여기서 https 사용 여부를 결정합니다.
	})
	if err != nil {
		return nil, fmt.Errorf("❌ MinIO 클라이언트 생성 실패: %w", err)
	}

	// 핑 테스트 (버킷 존재 확인)
	exists, err := minioClient.BucketExists(context.TODO(), bucketName)
	if err != nil || !exists {
		return nil, fmt.Errorf("❌ S3(iwinv) 연결 실패 또는 버킷 없음: %w", err)
	}

	return &fileService{
		client:     minioClient,
		bucketName: bucketName,
		domain:     domain,
	}, nil
}

func (s *fileService) UploadFile(ctx context.Context, file graphql.Upload, directory string) (string, error) {
	if closer, ok := file.File.(io.Closer); ok {
		defer func() {
			if err := closer.Close(); err != nil {
				log.Printf("⚠️ 파일 닫기 실패: %v", err)
			}
		}()
	}

	fileBytes, err := io.ReadAll(file.File)
	if err != nil {
		return "", fmt.Errorf("파일 스트림 읽기 에러: %w", err)
	}

	ext := filepath.Ext(file.Filename)
	uniqueFilename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)

	dir := "general"
	if directory != "" {
		dir = directory
	}
	objectKey := fmt.Sprintf("%s/%s", dir, uniqueFilename)

	reader := bytes.NewReader(fileBytes)
	_, err = s.client.PutObject(ctx, s.bucketName, objectKey, reader, int64(len(fileBytes)), minio.PutObjectOptions{
		ContentType: file.ContentType,
	})
	if err != nil {
		return "", fmt.Errorf("스토리지 업로드 에러: %w", err)
	}

	fileURL := fmt.Sprintf("%s/%s", s.domain, objectKey)
	return fileURL, nil
}
