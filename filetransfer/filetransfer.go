package filetransfer

import (
	"bytes"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

/*S3IOWriter a struct implement the io.Writer interface*/
type S3IOWriter struct {
	Bucket string
	Key    string
}

func (s S3IOWriter) Write(p []byte) (int, error) {
	session, sessionErr := session.NewSession()
	if sessionErr != nil {
		fmt.Println("error while creating aws session")
		return 0, sessionErr
	}

	svc := s3.New(session)
	_, uploadError := svc.PutObject(&s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.Key),
		Body:   bytes.NewReader(p),
	})
	if uploadError != nil {
		return 0, uploadError
	}
	return len(p), nil
}
