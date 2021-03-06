// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//
// Note:
// 1. It needs manually import: _ "github.com/lib/pq"
// 2. It does not support Save/Replace features.
// 3. It does not support LastInsertId.

package gdb

import (
	"database/sql"
	"fmt"
	"github.com/gogf/gf/internal/intlog"
	"github.com/gogf/gf/text/gstr"
	"strings"

	"github.com/gogf/gf/text/gregex"
)

type dbPgsql struct {
	*dbBase
}

// Open creates and returns a underlying sql.DB object for pgsql.
func (db *dbPgsql) Open(config *ConfigNode) (*sql.DB, error) {
	var source string
	if config.LinkInfo != "" {
		source = config.LinkInfo
	} else {
		source = fmt.Sprintf(
			"user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			config.User, config.Pass, config.Host, config.Port, config.Name,
		)
	}
	intlog.Printf("Open: %s", source)
	if db, err := sql.Open("postgres", source); err == nil {
		return db, nil
	} else {
		return nil, err
	}
}

// getChars returns the security char for this type of database.
func (db *dbPgsql) getChars() (charLeft string, charRight string) {
	return "\"", "\""
}

// handleSqlBeforeExec deals with the sql string before commits it to underlying sql driver.
func (db *dbPgsql) handleSqlBeforeExec(sql string) string {
	var index int
	// Convert place holder char '?' to string "$x".
	sql, _ = gregex.ReplaceStringFunc("\\?", sql, func(s string) string {
		index++
		return fmt.Sprintf("$%d", index)
	})
	sql, _ = gregex.ReplaceString(` LIMIT (\d+),\s*(\d+)`, ` LIMIT $1 OFFSET $2`, sql)
	return sql
}

// Tables retrieves and returns the tables of current schema.
// TODO
func (db *dbPgsql) Tables(schema ...string) (tables []string, err error) {
	return
}

// TableFields retrieves and returns the fields information of specified table of current schema.
func (db *dbPgsql) TableFields(table string, schema ...string) (fields map[string]*TableField, err error) {
	table = gstr.Trim(table)
	if gstr.Contains(table, " ") {
		panic("function TableFields supports only single table operations")
	}
	table, _ = gregex.ReplaceString("\"", "", table)
	checkSchema := db.schema.Val()
	if len(schema) > 0 && schema[0] != "" {
		checkSchema = schema[0]
	}
	v := db.cache.GetOrSetFunc(
		fmt.Sprintf(`pgsql_table_fields_%s_%s`, table, checkSchema), func() interface{} {
			var result Result
			var link *sql.DB
			link, err = db.getSlave(checkSchema)
			if err != nil {
				return nil
			}
			result, err = db.doGetAll(link, fmt.Sprintf(`
			SELECT a.attname AS field, t.typname AS type FROM pg_class c, pg_attribute a 
	        LEFT OUTER JOIN pg_description b ON a.attrelid=b.objoid AND a.attnum = b.objsubid,pg_type t
	        WHERE c.relname = '%s' and a.attnum > 0 and a.attrelid = c.oid and a.atttypid = t.oid 
			ORDER BY a.attnum`, strings.ToLower(table)))
			if err != nil {
				return nil
			}

			fields = make(map[string]*TableField)
			for i, m := range result {
				fields[m["field"].String()] = &TableField{
					Index: i,
					Name:  m["field"].String(),
					Type:  m["type"].String(),
				}
			}
			return fields
		}, 0)
	if err == nil {
		fields = v.(map[string]*TableField)
	}
	return
}
