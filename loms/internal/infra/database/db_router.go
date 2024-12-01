package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/spaolacci/murmur3"
	"log"
	"strconv"
)

type FallbackConnection struct {
	currentConn  *pgxpool.Conn
	fallBackConn *pgxpool.Conn
}

func (fc *FallbackConnection) Begin(ctx context.Context) (pgx.Tx, error) {
	return fc.currentConn.Begin(ctx)
}

func (fc *FallbackConnection) Release() {
	if fc.currentConn != nil {
		fc.currentConn.Release()
	}
	if fc.fallBackConn != nil {
		fc.fallBackConn.Release()
	}
}

/*
также хочу реализовать примитивно fallback поведение
при ошибке на реплике, попробуем выполнить запрос на мастере
проблема в том что любая ошибка будет расценена как ошибка на реплике и запроску к master
может быть стоить различать ошибки бд и ошибки соединения, но пока хочу сделать по-простому
*/

func (fc *FallbackConnection) Exec(ctx context.Context, query string, args ...interface{}) (pgconn.CommandTag, error) {
	tag, err := fc.currentConn.Exec(ctx, query, args...)
	if err != nil {
		fmt.Println("[infra] Detected connection error")
		// Выполнить фоллбэк, если ошибка подключения
		if fc.fallBackConn != nil {
			return fc.fallBackConn.Exec(ctx, query, args...)
		}
	}
	return tag, err
}

func (fc *FallbackConnection) Query(ctx context.Context, query string, args ...interface{}) (pgx.Rows, error) {
	rows, err := fc.currentConn.Query(ctx, query, args...)
	if err != nil {
		fmt.Println("[infra] Detected connection error")
		if fc.fallBackConn != nil {
			return fc.fallBackConn.Query(ctx, query, args...)
		}
	}
	return rows, err
}

func (fc *FallbackConnection) QueryRow(ctx context.Context, query string, args ...interface{}) pgx.Row {
	return &fallbackRow{
		currentRow:   fc.currentConn.QueryRow(ctx, query, args...),
		fallBackConn: fc.fallBackConn,
		ctx:          ctx,
		query:        query,
		args:         args,
	}
}

// fallbackRow структура обертки для Row с fallback
type fallbackRow struct {
	currentRow   pgx.Row
	fallBackConn *pgxpool.Conn
	ctx          context.Context
	query        string
	args         []interface{}
}

// Scan метод для сканирования данных с fallback
func (f *fallbackRow) Scan(dest ...interface{}) error {
	// Выполняем первичное сканирование
	err := f.currentRow.Scan(dest...)
	if err != nil {
		if f.fallBackConn != nil {
			fmt.Println("[infra] Fallback to master database")
			masterRow := f.fallBackConn.QueryRow(f.ctx, f.query, f.args...)
			return masterRow.Scan(dest...)
		}
	}
	return err
}

/*
Структура которая может поддерживать работу с двумя базами данных, реализует метод для извлечения пула а так же
методы для выполнения запросов к базе данных, то что нужно моему приложению в будущем можно расширить либо переходит на
внешнюю балансировку
*/
type DBRouter struct {
	shards     []*MasterAndReplica
	countShard int
}

type MasterAndReplica [2]*pgxpool.Pool

// NewDBRouter masterDB mustn't nil, and replicaDB can be nil
func NewDBRouter(shards []*MasterAndReplica) *DBRouter {
	return &DBRouter{shards: shards, countShard: len(shards)}
}

func (db *DBRouter) PickConnFromUserId(ctx context.Context, userId int64, readOnlyOperation bool) (*FallbackConnection, error) {
	shardKey := strconv.FormatInt(userId, 10)
	hash := hashCode(shardKey)
	shardIndex := int(hash) % db.countShard
	return db.pickConnectionFromShards(ctx, shardIndex, readOnlyOperation)
}

func (db *DBRouter) PickConnFromOrderId(ctx context.Context, orderID int64, readOnlyOperation bool) (*FallbackConnection, error) {
	shardIndex := orderID % 1000
	return db.pickConnectionFromShards(ctx, int(shardIndex), readOnlyOperation)
}

func (db *DBRouter) PickDefaultShard(ctx context.Context, readOnlyOperation bool) (*FallbackConnection, error) {
	return db.pickConnectionFromShards(ctx, 0, readOnlyOperation)
}

func (db *DBRouter) PickAllShards(ctx context.Context, readOnlyOperation bool) ([]*FallbackConnection, error) {
	var connections []*FallbackConnection
	for i := 0; i < db.countShard; i++ {
		conn, err := db.pickConnectionFromShards(ctx, i, readOnlyOperation)
		if err != nil {
			return connections, err
		}
		connections = append(connections, conn)
	}
	return connections, nil
}

func (db *DBRouter) pickConnectionFromShards(ctx context.Context, shardIndex int, readOnlyOperation bool) (*FallbackConnection, error) {
	masterConn, err := db.shards[shardIndex][0].Acquire(ctx)
	if err != nil {
		return &FallbackConnection{}, err
	}
	var replicaConn *pgxpool.Conn
	if readOnlyOperation && db.shards[shardIndex][1] != nil {
		replicaConn, err = db.shards[shardIndex][1].Acquire(ctx)
		if err != nil {
			log.Println(err)
			return &FallbackConnection{currentConn: masterConn}, nil
		}
	} else {
		return &FallbackConnection{currentConn: masterConn}, nil
	}

	return &FallbackConnection{currentConn: replicaConn, fallBackConn: masterConn}, nil
}

func hashCode(key string) uint32 {
	var hasher = murmur3.New32()

	_, _ = hasher.Write([]byte(key))

	return hasher.Sum32()
}
