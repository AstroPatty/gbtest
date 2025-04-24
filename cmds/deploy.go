package cmds


import (
    "fmt"
    "context"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/sts" 
    "github.com/aws/aws-sdk-go-v2/service/iam")

func validatePermisions() {
    cfg, err := config.LoadDefaultConfig(context.TODO())
    if err != nil {
        panic(fmt.Sprintf("Failed to load config. Have you set up the CLI?"))
    }

    stsClient := sts.NewFromConfig(cfg)
    identity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
    if err != nil {
        panic(fmt.Sprintf("Failed to get caller identity"))
    }
    
    iamClient := iam.NewFromConfig(cfg)
    actions := []string{
        "s3:CreateBucket",
    }
    result, err := iamClient.SimulatePrincipalPolicy(context.TODO(), &iam.SimulatePrincipalPolicyInput {
        PolicySourceArn: aws.String(*identity.Arn),
        ActionNames: actions,
    })
    fmt.Println(*result.EvaluationResults[0].EvalActionName)
    fmt.Println(result.EvaluationResults[0].EvalDecision)

    fmt.Println(*result.EvaluationResults[0].EvalActionName)


}

func Deploy() {
    validatePermisions()
}
