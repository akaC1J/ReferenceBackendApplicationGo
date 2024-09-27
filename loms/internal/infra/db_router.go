package infra

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

const ReadOnlyKey = "READ_ONLY"

/*
Структура которая может поддерживать работу с двумя базами данных, реализует метод для извлечения пула а так же
методы для выполнения запросов к базе данных, то что нужно моему приложению в будущем можно расширить либо переходит на
внешнюю балансировку
*/
type DBRouter struct {
	masterDB  *pgxpool.Pool
	replicaDB *pgxpool.Pool
}

func (db *DBRouter) Acquire(ctx context.Context) (c *pgxpool.Conn, err error) {
	dbInstance := db.chooseDBTXInstanceByContext(ctx)
	return dbInstance.Acquire(ctx)
}

// NewDBRouter masterDB mustn't nil, and replicaDB can be nil
func NewDBRouter(masterDB *pgxpool.Pool, replicaDB *pgxpool.Pool) *DBRouter {
	return &DBRouter{masterDB: masterDB, replicaDB: replicaDB}
}

/*
также хочу реализовать примитивно fallback поведение
при ошибке на реплике, попробуем выполнить запрос на мастере
проблема в том что любая ошибка будет расценена как ошибка на реплике и запроску к master
может быть стоить различать ошибки бд и ошибки соединения, но пока хочу сделать по-простому
*/

func (db *DBRouter) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	dbInstance := db.chooseDBTXInstanceByContext(ctx)
	tag, err := dbInstance.Exec(ctx, query, args...)
	if err != nil {
		fmt.Println("[infra] Detected connection error")
		// Выполнить фоллбэк, если ошибка подключения
		if dbInstance != db.masterDB {
			return db.masterDB.Exec(ctx, query, args...)
		}
	}
	return tag, err
}

func (db *DBRouter) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	dbInstance := db.chooseDBTXInstanceByContext(ctx)
	rows, err := dbInstance.Query(ctx, query, args...)
	if err != nil {
		fmt.Println("[infra] Detected connection error")
		if dbInstance != db.masterDB {
			return db.masterDB.Query(ctx, query, args...)
		}
	}
	return rows, err
}

func (db *DBRouter) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	dbInstance := db.chooseDBTXInstanceByContext(ctx)
	return &fallbackRow{
		primaryRow: dbInstance.QueryRow(ctx, query, args...),
		primaryDB:  dbInstance,
		masterDB:   db.masterDB,
		ctx:        ctx,
		query:      query,
		args:       args,
	}
}

// fallbackRow структура обертки для Row с fallback
type fallbackRow struct {
	primaryRow pgx.Row
	primaryDB  *pgxpool.Pool
	masterDB   *pgxpool.Pool
	ctx        context.Context
	query      string
	args       []interface{}
}

// Scan метод для сканирования данных с fallback
func (f *fallbackRow) Scan(dest ...interface{}) error {
	// Выполняем первичное сканирование
	err := f.primaryRow.Scan(dest...)
	if err != nil {
		if f.primaryDB != f.masterDB {
			fmt.Println("[infra] Fallback to master database")
			masterRow := f.masterDB.QueryRow(f.ctx, f.query, f.args...)
			return masterRow.Scan(dest...)
		}
	}
	return err
}

// выбирать будем на основе ключа контекста(учитываем что replicapool can be nil)
// может быть достаточно примитивная реализация, но свои задачи должна выполнять
func (db *DBRouter) chooseDBTXInstanceByContext(ctx context.Context) *pgxpool.Pool {
	if db.replicaDB == nil {
		return db.masterDB
	}
	if readOnly, ok := ctx.Value(ReadOnlyKey).(bool); ok && readOnly {
		return db.replicaDB
	}
	return db.masterDB
}
