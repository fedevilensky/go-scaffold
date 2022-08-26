package project

const (
	DBProviderGormMysql    = "gorm.io/driver/mysql"
	DBProviderGormPostgres = "gorm.io/driver/postgres"
	DBProviderMysql        = "github.com/go-sql-driver/mysql"
	DBProviderPostgres     = "github.com/lib/pq"
	DBProviderNone         = ""
)

const (
	DBLibraryGorm = "gorm.io/gorm"
	DBLibrarySql  = "database/sql"
	DBLibrarySqlx = "github.com/jmoiron/sqlx"
	DBLibraryNone = ""
)

const (
	WebLibraryGin        = "github.com/gin-gonic/gin"
	WebLibraryFiber      = "github.com/gofiber/fiber/v2"
	WebLibraryGorillamux = "github.com/gorilla/mux"
	WebLibraryHttp       = "net/http"
	WebLibraryNone       = ""
)
