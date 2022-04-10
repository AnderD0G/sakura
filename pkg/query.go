package pkg

import (
	"fmt"
	ctxLogger "github.com/luizsuper/ctxLoggers"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"reflect"
	"regexp"
	"strings"
)

const (
	Normal    = "Normal"
	JsonArray = "JsonArray"
	Array     = "Array"
)

type (
	varName string
	varType string
	TypeMap map[varName]varType
	RuleMap map[string]RuleType
)

type (
	RuleType struct {
		Rule       string
		Comparator string
		Value      []string
	}
	QueryCondition struct {
		Page  int
		Size  int
		Query string
	}
)

type Inquirer[T any] struct {
	M          T
	N          TypeMap
	Db         *gorm.DB
	QueryMap   map[string]RuleType
	condition  string
	page, size int
}

func (k *Inquirer[T]) ParseRule() error {
	c := k.condition
	q := k.QueryMap
	s := []byte(c)

	//对于进入的查询条件做检查
	if c == "" {
		return nil
	}

	queryReg := regexp.MustCompile(`([\w]+)(=|>=|<=|>|<)([^,]*)`)
	sub := queryReg.FindAllStringSubmatch(string(s[1:len(s)-1]), -1)

	if sub == nil {
		return errors.New("非法query")
	}

	for _, v := range sub {
		if len(v) != 4 {
			return errors.New("非法query")
		}

		s := varType("")
		ok := true
		key := v[1]
		value := v[3]
		comparator := v[2]

		if s, ok = k.N[varName(key)]; !ok {
			return errors.New(fmt.Sprintf("无效的key:%v", key))
		}

		switch s {
		case JsonArray:
			jsonArr := ""
			isArr := true

			if jsonArr, isArr = processJsonArr(key, value); !isArr {
				return errors.New(fmt.Sprintf("无效的jsonArr:%v", key))
			}

			q[key] = RuleType{
				Rule: JsonArray,
				Value: []string{
					jsonArr,
				},
			}

		case Array:
			q[key] = RuleType{
				Rule:  Array,
				Value: processArr(value),
			}

		case Normal:
			q[key] = RuleType{
				Rule: Normal,
				Value: []string{
					value,
				},
				Comparator: comparator,
			}
		}
	}
	return nil
}

func (s *Inquirer[T]) GetParam(queryMap *QueryCondition) {
	s.page = queryMap.Page
	s.condition = queryMap.Query
	s.size = queryMap.Size
}

func (s *Inquirer[T]) ParseStruct() {

	typ := reflect.TypeOf(s.M)
	val := reflect.ValueOf(s.M)

	if val.Kind().String() != reflect.Ptr.String() {
		ctxLogger.Error(nil, "is not ptr")
		panic(errors.New("is not ptr"))
	}
	if val.IsNil() {
		ctxLogger.Error(nil, "nil ptr")
		panic(errors.New("nil ptr"))
	}

	num := val.Elem().NumField()
	for i := 0; i < num; i++ {
		field := typ.Elem().Field(i)
		tag := field.Tag.Get("type")
		json := field.Tag.Get("json")
		s.N[varName(json)] = varType(tag)
		if tag == "" {
			s.N[varName(json)] = Normal
		}
	}

}

//处理json数组
func processJsonArr(pair ...string) (string, bool) {
	value := fmt.Sprintf("JSON_CONTAINS(%v,JSON_ARRAY(", pair[0])
	s := []byte(pair[1])

	//处理括号，遍历元素
	arr := strings.Split(string(s[1:len(s)-1]), "|")
	for _, v := range arr {
		//如果大括号里面没有元素
		if v == "" {
			return "", false
		}
		value = fmt.Sprintf("%v%v,", value, v)
	}

	s1 := []byte(value)
	value = string(s1[0:len(s1)-1]) + "))"
	return value, true
}

func processArr(str string) []string {
	s := []byte(str)

	arr := strings.Split(string(s[1:len(s)-1]), "|")
	i := make([]string, len(arr))

	for k, v := range i {
		i[k] = v
	}

	return i
}

func (s *Inquirer[T]) Query(t interface{}) {
	db := s.Db
	limit := s.size
	page := s.page

	for k, v := range s.QueryMap {
		if v.Rule == Normal {
			db = db.Where(fmt.Sprintf("%v %v ?", k, v.Comparator), v.Value[0])
		}
		if v.Rule == JsonArray {
			db = db.Where(v.Value[0])
		}
		if v.Rule == Array {
			db = db.Where(fmt.Sprintf("%v IN ?", k), v.Value)
		}
	}

	if limit > 0 {
		db = db.Limit(limit).Offset((page - 1) * limit)
	}

	db.Find(t)
}
