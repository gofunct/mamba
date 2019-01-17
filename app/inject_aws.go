//+build wireinject

package app

import (
	"context"
	awsclient "github.com/aws/aws-sdk-go/aws/client"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"gocloud.dev/aws/awscloud"
	"gocloud.dev/blob"
	"gocloud.dev/blob/s3blob"
	"gocloud.dev/mysql/rdsmysql"
)

func Aws(ctx context.Context, name string) (*Application, func(), error) {
	wire.Build(
		awscloud.AWS,
		rdsmysql.Open,
		ApplicationSet,
		awsBucket,
		awsSQLParams,
	)
	return nil, nil, nil
}

func awsBucket(ctx context.Context, cp awsclient.ConfigProvider) (*blob.Bucket, error) {
	return s3blob.OpenBucket(ctx, cp, viper.GetString("aws.bucket"), nil)
}

// awsSQLParams is a Wire provider function that returns the RDS SQL connection
// parameters based on the command-line c. Other providers inside
// awscloud.AWS use the parameters to construct a *sql.DB.
func awsSQLParams() *rdsmysql.Params {
	return &rdsmysql.Params{
		Endpoint: viper.GetString("aws.sql.endpoint"),
		Database: viper.GetString("aws.sql.database"),
		User:     viper.GetString("aws.sql.user"),
		Password: viper.GetString("aws.sql.password"),
	}
}
