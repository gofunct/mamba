//+build wireinject

package app

import (
	"context"
	"database/sql"
	"github.com/go-sql-driver/mysql"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"go.opencensus.io/trace"
	"gocloud.dev/blob"
	"gocloud.dev/blob/fileblob"
	"gocloud.dev/requestlog"
	"gocloud.dev/server"
)

func Local(ctx context.Context, name string) (*Application, func(), error) {
	// This will be filled in by Wire with providers from the provider sets in
	// wire.Build.
	wire.Build(
		wire.InterfaceValue(new(requestlog.Logger), requestlog.Logger(nil)),
		wire.InterfaceValue(new(trace.Exporter), trace.Exporter(nil)),
		server.Set,
		ApplicationSet,
		dialLocalSQL,
		localBucket,
	)
	return nil, nil, nil
}

// localBucket is a Wire provider function that returns a directory-based bucket
// based on the command-line c.
func localBucket() (*blob.Bucket, error) {
	return fileblob.OpenBucket(viper.GetString("local.bucket"), nil)
}

// dialLocalSQL is a Wire provider function that connects to a MySQL database
// (usually on localhost).
func dialLocalSQL() (*sql.DB, error) {
	cfg := &mysql.Config{
		Net:                  "tcp",
		Addr:                 viper.GetString("local.sql.region"),
		DBName:               viper.GetString("local.sql.name"),
		User:                 viper.GetString("local.sql.user"),
		Passwd:               viper.GetString("local.sql.password"),
		AllowNativePasswords: true,
	}
	return sql.Open("mysql", cfg.FormatDSN())
}
