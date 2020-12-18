package GormSQLBuilder

//package main

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"
)

type PersonDemo struct {
	Id          string    `json:"id" gorm:"column:id;primary_key"`                        //uuid
	User_id     string    `json:"guser_id" gorm:"column:guser_id"  gensql:"notnull"`      //uuid
	Age         int       `json:"age" gorm:"column:age"  gensql:"notnull"`                //int
	Weight      float32   `json:"weight" gorm:"column:weight"  gensql:"notnull"`          //float32
	Create_time time.Time `json:"create_time" gorm:"column:create_time" gensql:"notnull"` //timestamp without time zone
}

func (PersonDemo) TableName() string {
	return "Info.PersonDemo"
}

func test_main() {
	var person PersonDemo
	sql, _ := CreateSQL(person)
	fmt.Println(sql)
}

func FindColumn(s string) string {
	ss := strings.Split(s, ";")
	for _, v := range ss {
		kv := strings.Split(v, ":")
		if len(kv) >= 2 {
			if kv[0] == "column" {
				return kv[1]
			}
		}
	}
	return ""
}

func GetTableName(data interface{}) string {

	object := reflect.ValueOf(data)
	v := object.MethodByName("TableName")
	if v.Kind() == reflect.Func {

		return fmt.Sprintf("%v", v.Call([]reflect.Value{})[0])

	} else {
		s := fmt.Sprintf("%T", data)
		ss := strings.Split(s, ".")
		return ss[len(ss)-1]
	}

}

func CreateSQL(data interface{}) (string, error) {
	dataValue := reflect.ValueOf(data)

	//fmt.Println("dataValue:", dataValue)

	t := reflect.TypeOf(data)

	if dataValue.Kind() != reflect.Struct {

		return "", errors.New("Input is not a Struct.")
	}

	columns := ""
	values := ""

	for i := 0; i < t.NumField(); i++ {

		k := dataValue.Type().Field(i)

		value := dataValue.Field(i).Interface()

		var valueStr string

		switch value.(type) {
		case *string:
			if value.(*string) == nil {
			    valueStr = `null`
			} else {
			    valueStr = fmt.Sprintf(`'%s'`, *value.(*string))
			}
		case string:
			valueStr = fmt.Sprintf(`'%v'`, value)
		case time.Time:
			valueStr = fmt.Sprintf(`'%s'`, value.(time.Time).Format("2006-01-02 15:04:05"))
		default:
			valueStr = fmt.Sprintf("%v", value)
		}

		column := FindColumn(k.Tag.Get("gorm"))
		if column !="" {

			flag := k.Tag.Get("gensql")

			//fmt.Println(column, flag)

			if flag == "notnull" {
				if valueStr == `''` || valueStr == "0" || valueStr == `'0001-01-01 00:00:00 +0000 UTC'` || valueStr == `'0001-01-01 00:00:00'` {

				} else {
					values = values + valueStr + ","
					columns = columns + `"` + column + `",`
				}

			} else {
				values = values + valueStr + ","
				columns = columns + `"` + column + `",`
			}
		}

	}

	tablename := GetTableName(data)

	sql := fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`, tablename, columns[:len(columns)-1], values[:len(values)-1])
	//fmt.Println(sql)
	return sql, nil

}
