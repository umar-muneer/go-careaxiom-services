package filetransfer

import (
	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

/*S3IO a struct implement the io.Writer interface*/
type S3IO struct {
	Bucket string
	Key    string
}

func (s S3IO) Write(p []byte) (int, error) {
	session, sessionErr := session.NewSession()
	if sessionErr != nil {
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

func (s S3IO) Read(p []byte) (int, error) {
	session, sessionErr := session.NewSession()
	if sessionErr != nil {
		return 0, sessionErr
	}
	svc := s3.New(session)
	output, copyError := svc.GetObject(&s3.GetObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(s.Key),
	})
	if copyError != nil {
		return 0, copyError
	}
	noOfBytes, readError := output.Body.Read(p)
	if readError != nil {
		return noOfBytes, readError
	}
	return len(p), nil
}
