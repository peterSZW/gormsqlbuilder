# GormSQLBuilder

```
type PersonDemo struct {
	Id          string    `json:"id" gorm:"column:id;primary_key"`                         
	User_id     string    `json:"guser_id" gorm:"column:guser_id"  gensql:"notnull"`     
	Age         int       `json:"age" gorm:"column:age"  gensql:"notnull"` 
	Weight      float32   `json:"weight" gorm:"column:weight"  gensql:"notnull"` 
	Create_time time.Time `json:"create_time" gorm:"column:create_time" gensql:"notnull"` 
}
```
1.Add tags on structure 
```
gensql:"notnull"
```
2.Call CreateSQL to get the insert sql
```
var person PersonDemo
sql,err:=CreateSQL(person)
fmt.Println(sql)
```
