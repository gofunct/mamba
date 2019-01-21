package db

import (
	"bytes"
	"github.com/gofunct/mamba/function"
	"github.com/gofunct/mamba/logging"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
 	_ "github.com/go-sql-driver/mysql"
	"github.com/pkg/errors"
	"os"
)

type Model struct {
	gorm.Model
}

func OpenSqlite3(args ...string) *gorm.DB {
	db, err := gorm.Open("sqlite3", args)
	if err != nil {
		panic("failed to connect sqlite3 database")
	}
	return db
}
func OpenMySql(args ...string) *gorm.DB {
	db, err := gorm.Open("mysql", args)
	if err != nil {
		panic("failed to connect mysql database")
	}
	return db
}

func StartDGraph() error {
	var b []byte
	logging.L.Debug("pulling dgraph docker image...")
	pull, err := function.RunBytes("docker", "pull", "dgraph/dgraph")
	b = append(b, pull...)
	if err != nil {
		return errors.Wrap(err, "failed to execute:\n"+ string(b))
	}
	logging.L.Debug("making directory for data at ~/dgraph...")
	mkdir, err := function.RunBytes("mkdir", "-p", "~/dgraph")
	b = append(b, mkdir...)
	if err != nil {
		return errors.Wrap(err, "failed to execute:\n"+ string(b))
	}
	logging.L.Debug("running dgraph docker container...")
	run, err := function.RunBytes("docker", "run", "-it", "-p", "5080:5080", "-p", "6080:6080", "-p", "8080:8080", "-p", "9080:9080", "-p", "8000:8000", "-v", "~/dgraph:/dgraph", "--name", "dgraph dgraph/dgraph dgraph zero")
	b = append(b, run...)
	if err != nil {
		return errors.Wrap(err, "failed to execute:\n"+ string(b))
	}
	logging.L.Debug("executing dgraph docker container...")
	exec, err := function.RunBytes("docker", "exec", "-it", "dgraph dgraph alpha", "--lru_mb", "2048", "--zero", "localhost:5080")
	b = append(b, exec...)
	if err != nil {
		return errors.Wrap(err, "failed to execute:\n"+ string(b))
	}
	logging.L.Debug("executing dgraph ratel container...")
	ratel, err := function.RunBytes("docker", "exec", "-it", "dgraph dgraph-ratel")
	b = append(b, ratel...)
	if err != nil {
		return errors.Wrap(err, "failed to execute:\n"+ string(b))
	}

	buf := bytes.NewBuffer(b)
	_, err = buf.WriteTo(os.Stdout)
	if err != nil {
		return err
	}
	return nil


}
