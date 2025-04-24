package auth

import (
    "fmt"
    "context"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/iam"
    "github.com/aws/aws-sdk-go-v2/service/sts"
)

func GetCredentials() aws.Config {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        panic(fmt.Sprintf("Failed to load config. Have you set up the CLI?"))
    }
    validatePermissions(&cfg)
    return cfg

}

func validatePermissions(cfg *aws.Config) {
    stsClient := sts.NewFromConfig(*cfg)
    identity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
    if err != nil {
        panic(fmt.Sprintf("Unable to get caller identity"))
    }

    client := iam.NewFromConfig(*cfg)
    policy_arns := []string {
        "arn:aws:iam::aws:policy/AllowS3FullAccess",
    }

    for _, arn := range policy_arns {
        policy, err := client.GetPolicy(context.TODO(), &iam.GetPolicyInput {
            PolicyArn: aws.String(arn),
        })


    }




}

