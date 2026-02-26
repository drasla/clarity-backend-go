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

	"github.com/99designs/gqlgen/graphql"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Options struct {
	Domain     string
	Endpoint   string
	AccessKey  string
	SecretKey  string
	BucketName string
}

type fileService struct {
	client     *minio.Client
	bucketName string
	domain     string
}

type UploadResult struct {
	URL        string
	StoredName string
}

type FileService interface {
	UploadFile(ctx context.Context, file graphql.Upload, directory string) (*UploadResult, error)
}

func NewS3Service(opts S3Options) (FileService, error) {
	endpoint := strings.TrimPrefix(opts.Endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")

	minioClient, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(opts.AccessKey, opts.SecretKey, ""),
		Secure: true,
	})
	if err != nil {
		return nil, fmt.Errorf("❌ MinIO 클라이언트 생성 실패: %w", err)
	}

	exists, err := minioClient.BucketExists(context.TODO(), opts.BucketName)
	if err != nil || !exists {
		return nil, fmt.Errorf("❌ S3 연결 실패 또는 버킷 없음: %w", err)
	}

	return &fileService{
		client:     minioClient,
		bucketName: opts.BucketName,
		domain:     opts.Domain,
	}, nil
}

func (s *fileService) UploadFile(ctx context.Context, file graphql.Upload, directory string) (*UploadResult, error) {
	if closer, ok := file.File.(io.Closer); ok {
		defer func() {
			if err := closer.Close(); err != nil {
				log.Printf("⚠️ 파일 닫기 실패: %v", err)
			}
		}()
	}

	fileBytes, err := io.ReadAll(file.File)
	if err != nil {
		return nil, fmt.Errorf("파일 스트림 읽기 에러: %w", err)
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
		return nil, fmt.Errorf("스토리지 업로드 에러: %w", err)
	}

	fileURL := fmt.Sprintf("%s/%s/%s", s.domain, s.bucketName, objectKey)
	return &UploadResult{
		URL:        fileURL,
		StoredName: uniqueFilename,
	}, nil
}
