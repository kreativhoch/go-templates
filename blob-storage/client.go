package azureClient

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"go.uber.org/zap"
)

type ClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

type AzureClient struct {
	client *azblob.Client
	logger *zap.SugaredLogger
}

func NewClient(logger *zap.SugaredLogger) *AzureClient {
	accountName := os.Getenv("BLOB_STORAGE_ACCOUNT_NAME")
	if accountName == "" {
		logger.Fatalw("error on loading variable BLOB_STORAGE_ACCOUNT_NAME")
		return nil
	}
	credential, err := azidentity.NewDefaultAzureCredential(nil)
	if err != nil {
		logger.Fatalw("credentials could not be build with given data", err)
		return nil
	}
	client, err := azblob.NewClient(fmt.Sprintf("https://%s.blob.core.windows.net/", accountName), credential, nil)
	if err != nil {
		logger.Fatalw("client could not be created", err)
		return nil
	}

	return &AzureClient{
		client: client,
		logger: logger,
	}
}

func (s *AzureClient) PutObject(data []byte, blobStorageKey string) error {
	containerName := os.Getenv("BLOB_STORAGE_CONTAINER_NAME")
	if containerName == "" {
		s.logger.Fatalw("error on loading variable BLOB_STORAGE_CONTAINER_NAME")
	}

	_, err := s.client.UploadBuffer(context.Background(), containerName, blobStorageKey, data, &azblob.UploadBufferOptions{})
	if err != nil {
		return err
	}

	return nil
}
