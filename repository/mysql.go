package repository

import (
	"database/sql"
	"douyin_lite/settings"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
)

var Database *sql.DB

func InitSql() error {
	db, err := sql.Open("mysql",
		fmt.Sprintf("%s:%s@(%s:%s)/douyin?parseTime=true",
			settings.SqlUsername,
			settings.SqlPassword,
			settings.SqlIp,
			settings.SqlPort),
	)
	if err != nil {
		return err
	}
	Database = db
	return nil
}

// __UpdateByOptimisticLockUseIncrement__
// @param: queryGet(string), queryUpdate(string), id(int64), increment(int64)
// @return: err(error)
// queryGet format as "select change_element from table_chose where id=?" values(id) -> pre
// queryUpdate format as "update set change_element=? from table_chose where id=? and change_element=?" values(pre+increment,id,pre)
func __UpdateByOptimisticLockUseIncrement__(queryGet string, queryUpdate string, id int64, increment int64) error {
	for cnt := 0; cnt < settings.OptimisticLockMaxTryNumber; cnt++ {
		var pre int64
		err := Database.QueryRow(queryGet, id).Scan(&pre)
		if err != nil {
			return err
		}

		res, err := Database.Exec(queryUpdate, pre+increment, id, pre)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if affected > 0 {
			return nil
		}
	}
	return fmt.Errorf("fail to update after use so much times")
}

// __UpdateByOptimisticLock__
// @param: queryGet(string), queryUpdate(string), id(int64), newNumber(int64)
// @return: err(error)
// queryGet format as "select change_element from table_chose where id=?" values(id) -> pre
// queryUpdate format as "update set change_element=? from table_chose where id=? and change_element=?" values(newNumber,id,pre)
func __UpdateByOptimisticLock__(queryGet string, queryUpdate string, id int64, newNumber int64) error {
	for cnt := 0; cnt < settings.OptimisticLockMaxTryNumber; cnt++ {
		var pre int64
		err := Database.QueryRow(queryGet, id).Scan(&pre)
		if err != nil {
			return err
		}

		if newNumber == pre {
			return nil
		}

		res, err := Database.Exec(queryUpdate, newNumber, id, pre)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if affected > 0 {
			return nil
		}
	}
	return fmt.Errorf("fail to update after use so much times")
}
