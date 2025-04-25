package stack

import (
    "github.com/aws/jsii-runtime-go"
    "github.com/aws/constructs-go/constructs/v10"
    "github.com/aws/aws-cdk-go/awscdk/v2"
    "github.com/aws/aws-cdk-go/awscdk/v2/awss3"
    "github.com/aws/aws-cdk-go/awscdk/v2/awslambda"
    "fmt"
)

type Buckets struct {
    inputBucket awss3.Bucket
    outputBucket awss3.Bucket
}

func SynthDataHandlerStack() {
    appProps := awscdk.AppProps{
	Outdir: jsii.String("synth-output"),
    }
    app := awscdk.NewApp(&appProps)
    var sprops awscdk.StackProps


    stack_ := awscdk.NewStack(app, jsii.String("gbtest-processing-stack"), &sprops)
    projectName := awscdk.NewCfnParameter(stack_, jsii.String("ProjectName"), &awscdk.CfnParameterProps{
	Type: jsii.String("String"),
    })
    buckets := makeBuckets(stack_, projectName)
    makeLambda(stack_, projectName, &buckets)
    ast := app.Synth(nil)
    fmt.Printf("%s", ast)
}

func makeLambda(scope constructs.Construct, projectName awscdk.CfnParameter, buckets *Buckets) {
    lambda := awslambda.NewFunction(scope, jsii.String("MoveDataLambda"), &awslambda.FunctionProps{
	Runtime: awslambda.Runtime_PYTHON_3_12(),
	Handler: jsii.String("lambda.handler"),
	Code: awslambda.NewAssetCode(jsii.String("stack/lambda"), nil),
	Environment: &map[string]*string {
	    "INPUT_BUCKET_NAME": buckets.inputBucket.BucketName(),
	    "OUTPUT_BUCKET_NAME": buckets.outputBucket.BucketName(),
	},
	FunctionName: jsii.String(*projectName.ValueAsString() + "-transfer-function"),
    })
    buckets.inputBucket.GrantReadWrite(lambda, nil)
    buckets.outputBucket.GrantReadWrite(lambda, nil)
}



func makeBuckets(scope constructs.Construct, projectName awscdk.CfnParameter) Buckets {
    inputBucket := awss3.NewBucket(scope, jsii.String("InputBucket"), &awss3.BucketProps{
	BucketName: jsii.String(*projectName.ValueAsString() + "-input-bucket"),
	RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	AutoDeleteObjects: jsii.Bool(true),
    })
    outputBucket := awss3.NewBucket(scope, jsii.String("OutputBucket"), &awss3.BucketProps{
	    BucketName: jsii.String(*projectName.ValueAsString() + "-output-bucket"),
	    RemovalPolicy: awscdk.RemovalPolicy_DESTROY,
	    AutoDeleteObjects: jsii.Bool(true),
	})
    return Buckets {
	inputBucket: inputBucket,
	outputBucket: outputBucket,
    }

}


