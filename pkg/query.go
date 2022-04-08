package pkg

import (
	"fmt"
	ctxLogger "github.com/luizsuper/ctxLoggers"
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"sakura/model"
	"strings"
)

const (
	Normal    = "Normal"
	JsonArray = "JsonArray"
	Array     = "Array"
)
const (
	script = "script"
)

type (
	entityName string
	entityKey  string
	keyTag     string

	entityAttrMap map[entityName]nameAttr
	nameAttr      map[entityKey]keyTag
)

var (
	m1   entityAttrMap
	attr nameAttr
)

func init() {
	m1 = make(entityAttrMap)
	m2 := make(map[entityName]interface{})
	m2[script] = new(model.Scripts)
	attr = make(nameAttr)
	getAttr(m2)
}

type QueryCondition struct {
	Page  int64
	Size  int64
	Query string
}

type (
	RuleType struct {
		Rule       string
		Comparator string
		Value      []string
	}
	QueryMap map[string]RuleType
	Kv       map[string]string
)

func GenerateKv(str string) (QueryMap, error) {
	kv := make(QueryMap)
	s := []byte(str)
	//spilt := make([]string, 0)
	//对于进入的查询条件做检查
	if str == "" {
		return kv, nil
	}

	queryReg := regexp.MustCompile(`([\w]+)(=|>=|<=|>|<)([^,]*)`)
	submatch := queryReg.FindAllStringSubmatch(string(s[1:len(s)-1]), -1)

	if submatch == nil {
		return nil, errors.New("非法query")
	}

	for _, v := range submatch {
		if len(v) != 4 {
			return nil, errors.New("非法query")
		}

		s := keyTag("")
		ok := true
		key := v[1]
		value := v[3]
		comparator := v[2]

		if s, ok = attr[entityKey(key)]; !ok {
			return nil, errors.New(fmt.Sprintf("无效的key:%v", key))
		}

		switch s {
		case JsonArray:
			jsonArr := ""
			isArr := true

			if jsonArr, isArr = processJsonArr(key, value); !isArr {
				return nil, errors.New(fmt.Sprintf("无效的jsonArr:%v", key))
			}

			kv[key] = RuleType{
				Rule: JsonArray,
				Value: []string{
					jsonArr,
				},
			}

		case Array:
			kv[key] = RuleType{
				Rule:  Array,
				Value: processArr(value),
			}

		case Normal:
			kv[key] = RuleType{
				Rule: Normal,
				Value: []string{
					value,
				},
				Comparator: comparator,
			}
		}
	}
	return kv, nil
}

func GetParam(queryMap *QueryCondition) (page, limit int, query string) {
	if queryMap == nil {
		return
	}
	query = queryMap.Query
	page = int(queryMap.Page)
	limit = int(queryMap.Size)
	return
}

func getAttr(entity map[entityName]interface{}) {
	for name, body := range entity {
		m1[name] = attr

		typ := reflect.TypeOf(body)
		val := reflect.ValueOf(body)

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
			attr[entityKey(json)] = keyTag(tag)
			if tag == "" {
				attr[entityKey(json)] = Normal
			}
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
