package main

import (
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsrds"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
)

type RDSStackProps struct {
	awscdk.StackProps
}

func NewRDSStack(scope constructs.Construct, id string, props *RDSStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)
	vpc := awsec2.NewVpc(stack, jsii.String("MyVpc"), &awsec2.VpcProps{
		MaxAzs: jsii.Number(2),
	})

	securityGroup := awsec2.NewSecurityGroup(stack, jsii.String("RDS-SecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc:         vpc,
		Description: jsii.String("Allow connection to RDS"),
	})

	securityGroup.AddIngressRule(
		awsec2.Peer_AnyIpv4(),
		awsec2.Port_Tcp(jsii.Number(3306)),
		jsii.String("Allow MySQL connection from the internet"),
		nil,
	)

	// RDSインスタンスの作成
	dbInstance := awsrds.NewDatabaseInstance(stack, jsii.String("MyRDS"), &awsrds.DatabaseInstanceProps{
		Engine: awsrds.DatabaseInstanceEngine_Mysql(&awsrds.MySqlInstanceEngineProps{
			Version: awsrds.MysqlEngineVersion_VER_8_0(),
		}),
		InstanceType:           awsec2.NewInstanceType(jsii.String("t3.micro")),
		Vpc:                    vpc,
		SecurityGroups:         &[]awsec2.ISecurityGroup{securityGroup},
		RemovalPolicy:          awscdk.RemovalPolicy_DESTROY,
		DeleteAutomatedBackups: jsii.Bool(true),
	})

	awscdk.NewCfnOutput(stack, jsii.String("DBInstanceEndpoint"), &awscdk.CfnOutputProps{
		Value: dbInstance.DbInstanceEndpointAddress(),
	})

	return stack
}

func main() {
	app := awscdk.NewApp(nil)

	NewRDSStack(app, "MyRDSStack", &RDSStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}

func env() *awscdk.Environment {
	return &awscdk.Environment{
		Region:  jsii.String("ap-northeast-1"),
		Account: jsii.String(os.Getenv("ACCOUNT")),
	}
}
