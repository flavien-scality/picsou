package secrets

import (
  "os"
  "io/ioutil"
  "strings"
  "fmt"

  "github.com/aws/aws-sdk-go/aws"
  "github.com/aws/aws-sdk-go/aws/session"
  "github.com/aws/aws-sdk-go/service/kms"
  "github.com/aws/aws-sdk-go/service/s3"
  "github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func encryptSecret(name string, value string) (string, error) {

  kmsKeyARN := "arn:aws:kms:us-east-1:012345678910:key/0000000-0000-0000-0000-000000000000"
  kmsClient := kms.New(session.New(&aws.Config{
    Region: aws.String("us-east-1"),
  }))

  params := &kms.EncryptInput{
    KeyId:     aws.String(kmsKeyARN),
    Plaintext: []byte(value),
  }

  resp, err := kmsClient.Encrypt(params)
  if err != nil {
    return "", err
  }

  err = ioutil.WriteFile(name, resp.CiphertextBlob, 0644)
  if err != nil {
    return "", err
  }

  return name, nil
}

func uploadSecret(secretFileName string) error {

  s3Uploader := s3manager.NewUploader(session.New(&aws.Config{
    Region: aws.String("us-east-1"),
  }))

  reader, err := os.Open(secretFileName)
  if err != nil {
    return err
  }
  defer reader.Close()

  input := &s3manager.UploadInput{
    Bucket:           aws.String("my-bucket"),
    Key:              aws.String(secretFileName),
    Body:             reader,
  }
  _, err = s3Uploader.Upload(input)
  if err != nil {
    return err
  }

  err = os.Remove(secretFileName)
  if err != nil {
    return err
  }

  return nil
}

func downloadSecret(secretFileName string) (string, error) {

  s3Downloader := s3manager.NewDownloader(session.New(&aws.Config{
    Region: aws.String("us-east-1"),
  }))

  f, err := os.Create(secretName)
  if err != nil {
    return "", err
  }

  input := &s3.GetObjectInput{
    Bucket: aws.String("my-bucket"),
    Key:    aws.String(secretFileName),
  }
  _, err = s3Downloader.Download(f, input)
  if err != nil {
    return "", err
  }

  return f.Name(), nil
}

func decryptSecretFile(secretFile string) (string, error) {

  secretBytes, err := ioutil.ReadFile(secretFile)
  if err != nil {
    return "", err
  }

  kmsClient := kms.New(session.New(&aws.Config{
    Region: aws.String("us-east-1"),
  }))

  params := &kms.DecryptInput{
    CiphertextBlob: secretBytes,
  }

  resp, err := kmsClient.Decrypt(params)
  if err != nil {
    return "", err
  }

  err = os.Remove(secretFile)
  if err != nil {
    return "", err
  }

  return string(resp.Plaintext), nil
}
