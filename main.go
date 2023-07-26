package main

import (
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/iam"
	"os/exec"
	"strings"
)

func main() {
	ctx := context.Background()
	fmt.Println("loading aws client from local file")

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		panic(err)
	}

	iamC := iam.NewFromConfig(awsConfig)

	fmt.Println("Getting current keys")
	userKey, err := iamC.ListAccessKeys(ctx, &iam.ListAccessKeysInput{})
	if err != nil {
		panic(err.Error())
	}

	if len(userKey.AccessKeyMetadata) > 1 {
		panic("You cannot have more than one key to run this script")
	}

	oldKey := userKey.AccessKeyMetadata[0].AccessKeyId

	fmt.Println("Creating new key set")
	createdKeys, err := iamC.CreateAccessKey(ctx, &iam.CreateAccessKeyInput{})
	if err != nil {
		panic(err.Error())
	}

	cAccessKey := *createdKeys.AccessKey.AccessKeyId
	cSectetKey := *createdKeys.AccessKey.SecretAccessKey

	fmt.Println(cAccessKey)
	fmt.Println(cSectetKey)

	keySet := strings.Fields(fmt.Sprintf("configure set aws_access_key_id %s --profile default", cAccessKey))
	secretSet := strings.Fields(fmt.Sprintf("configure set aws_secret_access_key %s --profile default", cSectetKey))
	fmt.Println(keySet)
	fmt.Println(secretSet)

	setKeys(keySet)
	setKeys(secretSet)

	fmt.Println("Deleting old key set")
	_, err = iamC.DeleteAccessKey(ctx, &iam.DeleteAccessKeyInput{AccessKeyId: oldKey})
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("successfully updated keys")
}

func setKeys(k []string) {
	_, err := exec.Command("aws", k[0], k[1], k[2], k[3], k[4], k[5]).Output()
	if err != nil {
		panic(err.Error())
	}
	return
}
