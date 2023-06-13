package stores

import (
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"xorm.io/xorm"
)

// connectDB는 데이터베이스에 연결합니다.
func ConnectDB(dbPath string) (*xorm.Engine, error) {
	engine, err := xorm.NewEngine("sqlite3", fmt.Sprintf(dbPath))
	if err != nil {
		return nil, err
	}
	return engine, nil
}
