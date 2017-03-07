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

type KMS struct {
  KmsARN string
}

type Secret struct {
  BucketName string
  Kms *KMS
  Region string
}

func New(bucketName, region, kmsARN string) *Secret {
	k := &KMS{KmsARN: kmsARN}
	s := &Secret{BucketName: bucketName, Kms: k, Region: region}
	return s
}

func (s *Secret) encryptSecret(name string, value string) (string, error) {

  kmsKeyARN := "arn:aws:kms:eu-west-2:944690102204:key/a90a57ac-bb02-496c-a2d9-21e985f3387c "
  kmsClient := kms.New(session.New(&aws.Config{
    Region: aws.String("eu-west-2"),
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

func (s *Secret) uploadSecret(secretFileName string) error {

  s3Uploader := s3manager.NewUploader(session.New(&aws.Config{
    Region: aws.String("eu-west-2"),
  }))

  reader, err := os.Open(secretFileName)
  if err != nil {
    return err
  }
  defer reader.Close()

  input := &s3manager.UploadInput{
    Bucket:           aws.String("picsou-data"),
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

func (s *Secret) downloadSecret(secretFileName string) (string, error) {

  s3Downloader := s3manager.NewDownloader(session.New(&aws.Config{
    Region: aws.String("eu-west-2"),
  }))

  f, err := os.Create(secretName)
  if err != nil {
    return "", err
  }

  input := &s3.GetObjectInput{
    Bucket: aws.String("picsou-data"),
    Key:    aws.String(secretFileName),
  }
  _, err = s3Downloader.Download(f, input)
  if err != nil {
    return "", err
  }

  return f.Name(), nil
}

func (s *Secret) decryptSecretFile(secretFile string) (string, error) {

  secretBytes, err := ioutil.ReadFile(secretFile)
  if err != nil {
    return "", err
  }

  kmsClient := kms.New(session.New(&aws.Config{
    Region: aws.String("us-west-2"),
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
