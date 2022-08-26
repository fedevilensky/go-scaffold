package templates

import (
	"embed"
)

//go:embed "embedded/db_common"
var common embed.FS

//go:embed "embedded/gin"
var gin embed.FS

//go:embed "embedded/gin/pkg/httphelpers" "embedded/gin/pkg/taskutils"
var ginLean embed.FS

//go:embed "embedded/fiber"
var fiber embed.FS

//go:embed "embedded/fiber/pkg/httphelpers" "embedded/fiber/pkg/taskutils"
var fiberLean embed.FS

//go:embed "embedded/http" "embedded/httpcommon"
var http embed.FS

//go:embed "embedded/httpcommon"
var httpLean embed.FS

//go:embed "embedded/gorillamux" "embedded/httpcommon"
var gorillamux embed.FS

//go:embed "embedded/httpcommon"
var gorillamuxLean embed.FS

//go:embed "embedded/sql"
var sql embed.FS

//go:embed "embedded/sqlx"
var sqlx embed.FS

//go:embed "embedded/gorm"
var gorm embed.FS

//go:embed "embedded/mysql.tmpl"
var mysql embed.FS

//go:embed "embedded/postgresql.tmpl"
var postgresql embed.FS

//go:embed "embedded/gorm_mysql.tmpl"
var gormMysql embed.FS

//go:embed "embedded/gorm_postgresql.tmpl"
var gormPostgresql embed.FS

//go:embed "embedded/db_common/cmd/init_example_db" "embedded/db_common/helloworld" "embedded/db_common/internal"
var commonLean embed.FS
