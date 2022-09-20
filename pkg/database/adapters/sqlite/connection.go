package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/aabri-ankorstore/cli-auth/utils"
	"github.com/ankorstore/ankorstore-cli-core/pkg/util"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/sqlitedialect"
	"github.com/uptrace/bun/driver/sqliteshim"
	"github.com/uptrace/bun/extra/bundebug"
	"runtime"
)

var db *bun.DB

// InitDB creates sqlite db instance.
func InitDB(verbose bool) (*SqliteDB, error) {
	if db != nil {
		return nil, errors.New("nil configuration")
	}
	dirs := util.NewDirs()
	plugin := dirs.GetPluginsDir()
	sqlDB, err := sql.Open(
		sqliteshim.DriverName(),
		fmt.Sprintf("%s/%s/sessions.db", plugin, utils.PluginPath),
	)
	if err != nil {
		panic(err)
	}
	// for production use
	maxOpenConns := 4 * runtime.GOMAXPROCS(0)
	sqlDB.SetMaxOpenConns(maxOpenConns)
	sqlDB.SetMaxIdleConns(maxOpenConns)

	// create db
	db := bun.NewDB(sqlDB, sqlitedialect.New())
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("db ping error: %w", err)
	}
	v := bundebug.WithVerbose(verbose)
	hook := bundebug.NewQueryHook(v)
	db.AddQueryHook(hook)

	return &SqliteDB{DB: db}, nil
}
