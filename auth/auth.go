package auth

import (
    "fmt"
    "context"
    "encoding/json"
    "net/url"
    "time"
    "errors"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/iam"
    "github.com/aws/aws-sdk-go-v2/service/sts"
    "github.com/spf13/viper"
)

type PolicyStatement struct {
    Effect string   `json:"Effect"`
    Action any `json:"Action"`
    Resource any `json:"Resource"`
}

type PolicyDocument struct {
    Version string  `json:"Version"`
    Statement []PolicyStatement `json:"Statement"`
}

func CheckCredentials() error {
    authenticated := viper.GetBool("auth.authenticated")
    ttl := viper.GetInt64("auth.ttl")
    timeSinceCheck := time.Now().Unix() - viper.GetInt64("auth.lastChecked")

    
    if !authenticated || timeSinceCheck > ttl {
        fmt.Println("Checking AWS permissions")
        cfg, err := config.LoadDefaultConfig(context.TODO())
        if err != nil {
           return err
        }
        err = validatePermissions(&cfg)
        if err != nil {
            return err
        }
        lastChecked := time.Now().Unix()
        viper.Set("auth.authenticated", true)
        viper.Set("auth.lastChecked", lastChecked)
    }
    return nil
}

func validatePermissions(cfg *aws.Config) error {
    
    policy_arns := []string {
        "arn:aws:iam::aws:policy/AmazonS3FullAccess",
        "arn:aws:iam::aws:policy/AWSLambda_FullAccess",
        "arn:aws:iam::aws:policy/AmazonDynamoDBFullAccess",
        "arn:aws:iam::aws:policy/AwsCloudFormationFullAccess",
        "arn:aws:iam::aws:policy/IAMFullAccess",
    }
    policies, err := getRequiredActions(cfg, policy_arns)
    stsClient := sts.NewFromConfig(*cfg)
    ident, _ := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})

    iamClient := iam.NewFromConfig(*cfg)
    out, err := iamClient.SimulatePrincipalPolicy(context.TODO(), &iam.SimulatePrincipalPolicyInput{
        PolicySourceArn: ident.Arn,
        ActionNames: policies,
    })

    if err != nil {
        return err
    }
    notAllowed := make([]string, 0, len(out.EvaluationResults))
    for _, result := range out.EvaluationResults {
        if result.EvalDecision != "allowed" {
            notAllowed = append(notAllowed, *result.EvalActionName)
        }
    }
    if len(notAllowed) > 0 {
        return errors.New(fmt.Sprintf("Your account is missing required permissions: %s", notAllowed))
    }
    return nil



}

type GetPolicyActionsOutput struct {
    Actions []string
    Error error
}

func getPolicyActions(cfg *aws.Config, policyArn string, c chan <- GetPolicyActionsOutput) {
    client := iam.NewFromConfig(*cfg)
    actions := []string {}
    policy, err := client.GetPolicy(context.TODO(), &iam.GetPolicyInput{
            PolicyArn: aws.String(policyArn),
        })

        if err != nil {
            c <- GetPolicyActionsOutput{
                nil,
                err,
            }
            return
        }

        version, err := client.GetPolicyVersion(context.TODO(), &iam.GetPolicyVersionInput{
            PolicyArn: aws.String(policyArn),
            VersionId: aws.String(*policy.Policy.DefaultVersionId),
        })

        policyDoc, err := url.QueryUnescape(*version.PolicyVersion.Document)
        var doc PolicyDocument
        err = json.Unmarshal([]byte(policyDoc), &doc)
        if err != nil {
            c <- GetPolicyActionsOutput{
                nil,
                err,
            }
        }
        for _, stmt := range doc.Statement {
            switch v := stmt.Action.(type) {
                case string: {
                        actions = append(actions, v)
                    }
                case []any: {
                    acts := make([]string, len(v), len(v))
                    for i, act := range v {
                        acts[i] = act.(string)
                    }
                    actions = append(actions, acts...)
                    
                }
            }
        }
    c <- GetPolicyActionsOutput{
        actions,
        nil,
    }


}


func getRequiredActions(cfg *aws.Config, policyArns []string) ([]string, error) {
    var channels []chan GetPolicyActionsOutput
    for _, arn := range policyArns {
        c := make(chan GetPolicyActionsOutput)
        go getPolicyActions(cfg, arn, c)
        channels = append(channels, c)
    }
    actions := []string {}
    for _, c := range channels {
        result := <-c
        if result.Error != nil {
            return nil, result.Error
        }
        actions = append(actions, result.Actions...)
    }
    return actions, nil


}

