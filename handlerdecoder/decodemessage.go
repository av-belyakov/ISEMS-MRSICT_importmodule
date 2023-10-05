package handlerdecoder

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"

	"ISEMS-MRSICT/mytest/testdecodecasetostix/datatypes"
	"ISEMS-MRSICT/mytest/testdecodecasetostix/rules"
)

func ProcessingReflectAnySimpleType(
	logging chan<- datatypes.MessageLogging,
	chanCreateSTIXFormat chan datatypes.ChanInputCreateSTIXFormat,
	name interface{},
	anyType interface{},
	listRule *rules.ListRulesProcessingMsgMISP,
	num int,
	fieldBranch string) interface{} {

	var nameStr string
	r := reflect.TypeOf(anyType)

	if n, ok := name.(int); ok {
		nameStr = fmt.Sprintln(n)
	} else if n, ok := name.(string); ok {
		nameStr = n
	}

	if r == nil {
		return anyType
	}

	switch r.Kind() {
	case reflect.String:
		result := reflect.ValueOf(anyType).String()

		//
		//
		// а правила мне нужны обязательно, хотябы для исключения ненужных событий
		//
		//

		ncv, num, err := ReplacementRuleHandler(listRule, "string", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datatypes.MessageLogging{
				MsgData: fmt.Sprintf("'search value \"%s\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanCreateSTIXFormat <- datatypes.ChanInputCreateSTIXFormat{
			FieldName:   nameStr,
			ValueType:   "string",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv

	case reflect.Int, reflect.Int16, reflect.Int32, reflect.Int64:
		result := reflect.ValueOf(anyType).Int()

		ncv, num, err := ReplacementRuleHandler(listRule, "int", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datatypes.MessageLogging{
				MsgData: fmt.Sprintf("'search value \"%d\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanCreateSTIXFormat <- datatypes.ChanInputCreateSTIXFormat{
			FieldName:   nameStr,
			ValueType:   "int",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv

	case reflect.Uint, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		result := reflect.ValueOf(anyType).Uint()

		ncv, num, err := ReplacementRuleHandler(listRule, "uint", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datatypes.MessageLogging{
				MsgData: fmt.Sprintf("'search value \"%d\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanCreateSTIXFormat <- datatypes.ChanInputCreateSTIXFormat{
			FieldName:   nameStr,
			ValueType:   "uint",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv

	case reflect.Float32, reflect.Float64:
		result := reflect.ValueOf(anyType).Float()

		ncv, num, err := ReplacementRuleHandler(listRule, "float", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datatypes.MessageLogging{
				MsgData: fmt.Sprintf("'search value \"%v\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanCreateSTIXFormat <- datatypes.ChanInputCreateSTIXFormat{
			FieldName:   nameStr,
			ValueType:   "float",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv

	case reflect.Bool:
		result := reflect.ValueOf(anyType).Bool()

		ncv, num, err := ReplacementRuleHandler(listRule, "bool", nameStr, result)
		if err != nil {
			_, f, l, _ := runtime.Caller(0)

			logging <- datatypes.MessageLogging{
				MsgData: fmt.Sprintf("'search value \"%v\" from rule number \"%d\" of section \"REPLACE\" is not fulfilled' %s:%d", result, num, f, l-1),
				MsgType: "warning",
			}
		}

		PassRuleHandler(listRule.Rules.Pass, nameStr, ncv)

		chanCreateSTIXFormat <- datatypes.ChanInputCreateSTIXFormat{
			FieldName:   nameStr,
			ValueType:   "bool",
			Value:       ncv,
			FieldBranch: fieldBranch,
		}

		return ncv
	}

	return anyType
}

func ProcessingReflectMap(
	logging chan<- datatypes.MessageLogging,
	chanCreateSTIXFormat chan datatypes.ChanInputCreateSTIXFormat,
	l map[string]interface{},
	lr *rules.ListRulesProcessingMsgMISP,
	num int,
	fieldBranch string) map[string]interface{} {

	var (
		newMap  map[string]interface{}
		newList []interface{}
	)
	nl := map[string]interface{}{}

	for k, v := range l {
		var fbTmp string
		r := reflect.TypeOf(v)

		if r == nil {
			return nl
		}

		fbTmp = fieldBranch
		if fbTmp == "" {
			fbTmp += k
		} else {
			fbTmp += "." + k
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				newMap = ProcessingReflectMap(logging, chanCreateSTIXFormat, v, lr, num+1, fbTmp)
				nl[k] = newMap
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = ProcessingReflectSlice(logging, chanCreateSTIXFormat, v, lr, num+1, fbTmp)
				nl[k] = newList
			}

		default:
			nl[k] = ProcessingReflectAnySimpleType(logging, chanCreateSTIXFormat, k, v, lr, num, fbTmp)
		}
	}

	return nl
}

func ProcessingReflectSlice(
	logging chan<- datatypes.MessageLogging,
	chanCreateSTIXFormat chan datatypes.ChanInputCreateSTIXFormat,
	l []interface{},
	lr *rules.ListRulesProcessingMsgMISP,
	num int,
	fieldBranch string) []interface{} {

	var (
		newMap  map[string]interface{}
		newList []interface{}
	)
	nl := make([]interface{}, 0, len(l))

	for k, v := range l {
		r := reflect.TypeOf(v)

		if r == nil {
			return nl
		}

		switch r.Kind() {
		case reflect.Map:
			if v, ok := v.(map[string]interface{}); ok {
				newMap = ProcessingReflectMap(logging, chanCreateSTIXFormat, v, lr, num+1, fieldBranch)

				nl = append(nl, newMap)
			}

		case reflect.Slice:
			if v, ok := v.([]interface{}); ok {
				newList = ProcessingReflectSlice(logging, chanCreateSTIXFormat, v, lr, num+1, fieldBranch)

				nl = append(nl, newList...)
			}

		default:
			nl = append(nl, ProcessingReflectAnySimpleType(logging, chanCreateSTIXFormat, k, v, lr, num, fieldBranch))
		}
	}

	return nl
}

func PassRuleHandler(rulePass []rules.PassListAnd, fn string, cv interface{}) {
	cvstr := fmt.Sprint(cv)

	for key, value := range rulePass {
		for k, v := range value.ListAnd {
			if fn != v.SearchField {
				continue
			}

			if cvstr != v.SearchValue {
				continue
			}

			rulePass[key].ListAnd[k].StatementExpression = true
		}
	}
}

func ReplacementRuleHandler(lr *rules.ListRulesProcessingMsgMISP, svt, fn string, cv interface{}) (interface{}, int, error) {
	getReplaceValue := func(svt, rv string) (interface{}, error) {
		switch svt {
		case "string":
			return rv, nil

		case "int":
			return strconv.ParseInt(rv, 10, 64)

		case "uint":
			return strconv.ParseUint(rv, 10, 64)

		case "float":
			return strconv.ParseFloat(rv, 64)

		case "bool":
			return strconv.ParseBool(rv)
		}

		return rv, nil
	}

	for k, v := range lr.Rules.Replace {
		if v.SearchValue != fmt.Sprint(cv) {
			continue
		}

		if v.SearchField != "" {
			if v.SearchField == fn {
				nv, err := getReplaceValue(svt, v.ReplaceValue)

				return nv, k, err
			}

			continue
		}

		nv, err := getReplaceValue(svt, v.ReplaceValue)

		return nv, k, err
	}

	return cv, 0, nil
}
