//+build wireinject

package app

import (
	"context"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"gocloud.dev/blob"
	"gocloud.dev/blob/gcsblob"
	"gocloud.dev/gcp"
	"gocloud.dev/gcp/gcpcloud"
	"gocloud.dev/mysql/cloudmysql"
)

func Gcp(ctx context.Context, name string) (*Application, func(), error) {
	// This will be filled in by Wire with providers from the provider sets in
	// wire.Build.
	wire.Build(
		gcpcloud.GCP,
		cloudmysql.Open,
		ApplicationSet,
		gcpBucket,
		gcpSQLParams,
	)
	return nil, nil, nil
}

func gcpBucket(ctx context.Context, client *gcp.HTTPClient) (*blob.Bucket, error) {
	return gcsblob.OpenBucket(ctx, client, viper.GetString("gcs.bucket"), nil)
}

// gcpSQLParams is a Wire provider function that returns the Cloud SQL
// connection parameters based on the command-line c. Other providers inside
// gcpcloud.GCP use the parameters to construct a *sql.DB.
func gcpSQLParams(id gcp.ProjectID) *cloudmysql.Params {
	return &cloudmysql.Params{
		ProjectID: string(id),
		Region:    viper.GetString("gcs.sql.region"),
		Instance:  viper.GetString("gcs.sql.instance"),
		Database:  viper.GetString("gcs.sql.database"),
		User:      viper.GetString("gcs.sql.user"),
		Password:  viper.GetString("gcs.sql.password"),
	}
}
