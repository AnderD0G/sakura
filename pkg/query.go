package pkg

import (
	"fmt"
	ctxLogger "github.com/luizsuper/ctxLoggers"
	"github.com/pkg/errors"
	"go.uber.org/zap"
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
	Query struct {
		Page      int
		Size      int
		Condition string
	}
)

type Inquirer[T any] struct {
	M          interface{}
	N          TypeMap
	Db         *gorm.DB
	QueryMap   map[string]RuleType
	condition  string
	Page, Size int
}

func (k *Inquirer[T]) ParseQuery() error {
	c := k.condition
	m := make(map[string]RuleType)
	k.QueryMap = m
	s := []byte(c)

	//对于进入的查询条件做检查
	if c == "" {
		return nil
	}

	queryReg := regexp.MustCompile(`([\w]+)(=|>=|<=|>|<|:)([^,]*)`)
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

			m[key] = RuleType{
				Rule: JsonArray,
				Value: []string{
					jsonArr,
				},
			}

		case Normal:
			m[key] = RuleType{
				Rule: Normal,
				Value: []string{
					value,
				},
				Comparator: comparator,
			}
			if comparator == ":" {
				m[key] = RuleType{
					Rule:       Array,
					Value:      strings.Split(strings.TrimSuffix(strings.TrimPrefix(value, "("), ")"), "|"),
					Comparator: comparator,
				}
			}
		}
	}
	return nil
}

func (s *Inquirer[T]) InjectParam(queryMap *Query) {
	s.Page = queryMap.Page
	s.condition = queryMap.Condition
	s.Size = queryMap.Size
}

func (s *Inquirer[T]) parseStruct() {

	if s.N == nil {
		s.N = make(TypeMap)
	}

	typ := reflect.TypeOf(s.M)
	val := reflect.ValueOf(s.M)

	if val.Kind().String() != reflect.Ptr.String() {
		ctxLogger.FError(nil, "is not ptr", zap.String("finally get", val.Kind().String()))
		panic(errors.New("is not ptr"))
	}
	if val.IsNil() {
		ctxLogger.Error(nil, "nil ptr")
		panic(errors.New("nil ptr"))
	}

	num := val.Elem().NumField()
	for i := 0; i < num; i++ {
		field := typ.Elem().Field(i)
		//fmt.Printf("name:%v,kind:%v", field.Name, val.Elem().Field(i).Kind())
		//递归解析结构体
		if v := val.Elem().Field(i); v.Kind() == reflect.Struct {
			s.M = v.Addr().Interface()
			s.parseStruct()
		}

		tag := field.Tag.Get("type")
		json := field.Tag.Get("json")

		s.N[varName(json)] = varType(tag)
		if tag == "" {
			s.N[varName(json)] = Normal
		}
	}

	return
}
func (s *Inquirer[T]) ParseStruct() {
	if s.N == nil {
		ctxLogger.FInfo(nil, "first", zap.String("which", fmt.Sprintf("%T", s)))
		s.parseStruct()
	}
	return

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

func (s *Inquirer[T]) Query(table string, t interface{}, f ...func(d *gorm.DB)) {
	db := s.Db.Table(table)
	limit := s.Size
	page := s.Page

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

	if len(f) == 1 {
		f[0](db)
		return
	}

	db.Debug().Find(t)
}
