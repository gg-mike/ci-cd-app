package vault

import (
	"context"
	"time"

	"github.com/hashicorp/vault-client-go"
	"github.com/hashicorp/vault-client-go/schema"
)

var client *vault.Client

func Init(address, token string) error {
	var err error
	client, err = vault.New(
		vault.WithAddress(address),
		vault.WithRequestTimeout(10*time.Second),
	)
	if err != nil {
		return err
	}

	if err = client.SetToken(token); err != nil {
		return err
	}
	return nil
}

func Set(key string, data map[string]any) error {
	ctx := context.Background()
	_, err := client.Secrets.KvV2Write(ctx, key, schema.KvV2WriteRequest{ Data: data }, vault.WithMountPath("secret"))
	return err
}

func Get(key string) (map[string]any, error) {
	ctx := context.Background()
	secret, err := client.Secrets.KvV2Read(ctx, key, vault.WithMountPath("secret"))
	if err != nil {
		return map[string]any{}, err
	}
	return secret.Data.Data, nil
}

func Del(key string) error {
	ctx := context.Background()
	_, err := client.Secrets.KvV2Delete(ctx, key, vault.WithMountPath("secret"))
	if err != nil {
		return err
	}
	_, err = client.Secrets.KvV2DeleteMetadataAndAllVersions(ctx, key, vault.WithMountPath("secret"))
	return err
}
