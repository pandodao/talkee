package attachment

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"talkee/core"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
)

type Config struct {
	AwsBucket string
	AwsRegion string
	AwsKey    string
	AwsSecret string
}

type AttachmentService struct {
	cfg         Config
	awsSession  *session.Session
	attachments core.AttachmentStore
}

func New(cfg Config, attachments core.AttachmentStore) *AttachmentService {
	awsSession, err := session.NewSession(&aws.Config{
		Region: aws.String(cfg.AwsRegion),
		Credentials: credentials.NewStaticCredentials(
			cfg.AwsKey,    // id
			cfg.AwsSecret, // secret
			""),           // token can be left blank for now
	})
	if err != nil {
		panic(err)
	}

	return &AttachmentService{
		cfg:         cfg,
		awsSession:  awsSession,
		attachments: attachments,
	}
}

func (s *AttachmentService) uploadFileToS3(url string, s3FileName string, sizeOutput *int64) error {
	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Close = true
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	*sizeOutput = resp.ContentLength
	buffer, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// upload to s3
	if _, err := s3.New(s.awsSession).PutObject(&s3.PutObjectInput{
		Bucket:      aws.String(s.cfg.AwsBucket),
		Key:         aws.String(s3FileName),
		ACL:         aws.String("public-read"), // could be private if you want it to be access by only authorized users
		Body:        bytes.NewReader(buffer),
		ContentType: aws.String("audio/mpeg"),
		// ContentLength:        aws.Int64(resp.ContentLength),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
	}); err != nil {
		fmt.Printf("uploader.Upload err: %v\n", err)
		return err
	}

	return err
}
