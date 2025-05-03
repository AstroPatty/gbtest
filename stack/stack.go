package stack

import (
	"os"
	"errors"
	"os/exec"
	"path/filepath"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsdynamodb"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3"
	"github.com/aws/aws-cdk-go/awscdk/v2/awss3notifications"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/spf13/viper"
)

type Buckets struct {
    inputBucket awss3.Bucket
    outputBucket awss3.Bucket
}

func DestroyStack() error {
    configDir, _ := os.UserConfigDir()
    appConfigDir := filepath.Join(configDir, "gbtest")
    stackOutputDir := filepath.Join(appConfigDir, "stack")
    if fData, err := os.Stat(stackOutputDir); err != nil {
	return errors.New("Unable to destroy stack. Are you sure it has been deployed?")
    } else if !fData.IsDir() {
	return errors.New("Unable to destroy stack. Did you update any files manually?")
    }

	
    cmd := exec.Command("cdk", "destroy", "--app", stackOutputDir)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    cmd.Run()
    os.RemoveAll(stackOutputDir)
    return nil
}

func DeployStack() error {
    configDir, _ := os.UserConfigDir()
    appConfigDir := filepath.Join(configDir, "gbtest")
    stackOutputDir := filepath.Join(appConfigDir, "stack")
    err := os.MkdirAll(stackOutputDir, 0750)
    if err != nil {
	return err
	
    }
    synthDataHandlerStack(stackOutputDir)
    cmd := exec.Command("cdk", "deploy", "--app", stackOutputDir)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin
    cmd.Run()
    return nil
}

func synthDataHandlerStack(stackOutputDir string) {
    stackName := jsii.String(viper.GetString("stacks.processing.stackName"))

    appProps := awscdk.AppProps{
	Outdir: jsii.String(stackOutputDir),
    }
    app := awscdk.NewApp(&appProps)
    var sprops awscdk.StackProps

    stack_ := awscdk.NewStack(app, stackName, &sprops)
    buckets := makeBuckets(stack_)
    db := makeDb(stack_)
    makeTransferLambda(stack_, &buckets, db)
    makeInitializeLambda(stack_, db)
    makeDeleteLambda(stack_, db, buckets.outputBucket)
    app.Synth(nil)

}

func makeUsers(scope constructs.Construct, bucket awss3.Bucket) {
    uploadUser := awsiam.NewUser(scope, jsii.String("gbtest-upload-user"), &awsiam.UserProps{})
    bucket.GrantPut(uploadUser, nil)
    uploadApiKey := awsiam.NewAccessKey(scope, jsii.String("upload-user-access-key"), &awsiam.AccessKeyProps{
	User: uploadUser,
    })
    awscdk.NewCfnOutput(scope, jsii.String("upload-user-access-key"), &awscdk.CfnOutputProps {
	Value: uploadApiKey.AccessKeyId(),
    })
    awscdk.NewCfnOutput(scope, jsii.String("upload-user-access-secret"), &awscdk.CfnOutputProps {
	Value: uploadApiKey.SecretAccessKey().ToString(),
    })
}


func makeTransferLambda(scope constructs.Construct, buckets *Buckets, db awsdynamodb.TableV2) {
    transferLambda := awslambda.NewFunction(scope, jsii.String("MoveDataLambda"), &awslambda.FunctionProps{
	Runtime: awslambda.Runtime_PYTHON_3_12(),
	Handler: jsii.String("lambda.handler"),
	Code: awslambda.NewAssetCode(jsii.String("stack/transfer"), nil),
	Environment: &map[string]*string {
	    "INPUT_BUCKET_NAME": buckets.inputBucket.BucketName(),
	    "OUTPUT_BUCKET_NAME": buckets.outputBucket.BucketName(),
	    "DB_NAME": db.TableName(),
	},
	FunctionName: jsii.String("transfer-function"),
    })
    buckets.inputBucket.GrantReadWrite(transferLambda, nil)
    buckets.outputBucket.GrantReadWrite(transferLambda, nil)
    db.GrantReadWriteData(transferLambda)
}

func makeUploadNotification(scope constructs.Construct, bucket awss3.Bucket, lambda awslambda.Function) error {
    bucket.AddEventNotification(awss3.EventType_OBJECT_CREATED, awss3notifications.NewLambdaDestination(lambda))
    return nil
}

func makeInitializeLambda(scope constructs.Construct, db awsdynamodb.TableV2) {
    initializeLambda := awslambda.NewFunction(scope, jsii.String("initProjectLambda"), &awslambda.FunctionProps{
	Runtime: awslambda.Runtime_PYTHON_3_12(),
	Handler: jsii.String("lambda.handler"),
	Code: awslambda.NewAssetCode(jsii.String("stack/initialize"), nil),
	Environment: &map[string]*string {
	    "DB_NAME": db.TableName(),
	},
	FunctionName: jsii.String("init-project-function"),
    })
    initializeLambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
	Actions: jsii.Strings(
	    "dynamodb:PutItem",
	    "dynamodb:GetItem",
	),
	Resources: jsii.Strings(*db.TableArn()),
    }))
}

func makeDeleteLambda(scope constructs.Construct, db awsdynamodb.TableV2, outputBucket awss3.Bucket) {
    deleteLambda := awslambda.NewFunction(scope, jsii.String("deleteProjectLambda"), &awslambda.FunctionProps{
	Runtime: awslambda.Runtime_PYTHON_3_12(),
	Handler: jsii.String("lambda.handler"),
	Code: awslambda.NewAssetCode(jsii.String("stack/delete"), nil),
	Environment: &map[string]*string {
	    "DB_NAME": db.TableName(),
	    "BUCKET_NAME": outputBucket.BucketName(),
	},
	FunctionName: jsii.String("delete-project-function"),
    })
    deleteLambda.AddToRolePolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
	Actions: jsii.Strings(
	    "dynamodb:GetItem",
	    "dynamodb:DeleteItem",
	),
	Resources: jsii.Strings(*db.TableArn()),
    }))
    outputBucket.GrantRead(deleteLambda, nil)
    outputBucket.GrantDelete(deleteLambda, nil)
}


func makeBuckets(scope constructs.Construct) Buckets {
    inputBucket := awss3.NewBucket(scope, jsii.String("InputBucket"), &awss3.BucketProps{
	BucketName: jsii.String("gbtest-input-bucket"),
	RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	AutoDeleteObjects: jsii.Bool(true),
    })
    outputBucket := awss3.NewBucket(scope, jsii.String("OutputBucket"), &awss3.BucketProps{
	    BucketName: jsii.String("gbtest-output-bucket"),
	    RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	    AutoDeleteObjects: jsii.Bool(true),
	})

    awscdk.NewCfnOutput(scope, jsii.String("upload-bucket-arn"), &awscdk.CfnOutputProps {
	Value: inputBucket.BucketArn(),
    })

    awscdk.NewCfnOutput(scope, jsii.String("downlaod-bucket-arn"), &awscdk.CfnOutputProps {
	Value: outputBucket.BucketArn(),
    })
    return Buckets {
	inputBucket: inputBucket,
	outputBucket: outputBucket,
    }

}

func makeDb(scope constructs.Construct) awsdynamodb.TableV2 {
    db := awsdynamodb.NewTableV2(scope, jsii.String("gbtest-database"), 
	&awsdynamodb.TablePropsV2{
	    PartitionKey: &awsdynamodb.Attribute {
		Name: jsii.String("project-name"),
		Type: awsdynamodb.AttributeType_STRING,
	    },
	})
    return db
}


