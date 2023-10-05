package handlerdecoder

import (
	"encoding/json"
	"fmt"
	"runtime"

	"ISEMS-MRSICT/mytest/testdecodecasetostix/datatypes"
	"ISEMS-MRSICT/mytest/testdecodecasetostix/rules"
	"ISEMS-MRSICT/mytest/testdecodecasetostix/supportfunc"
)

func NewHandlerDecoder(b []byte,
	listRule rules.ListRulesProcessingMsgMISP,
	cstixf chan datatypes.ChanInputCreateSTIXFormat,
	cmispfDone chan<- bool,
	logging chan<- datatypes.MessageLogging) {

	var isAllowed bool
	listTmp := map[string]interface{}{}

	//для записи событий в файл events
	str, _ := supportfunc.NewReadReflectJSONSprint(b)
	logging <- datatypes.MessageLogging{
		MsgData: fmt.Sprintf("\t---------------\n\tEVENTS:\n%s\n", str),
		MsgType: "events",
	}

	if err := json.Unmarshal(b, &listTmp); err != nil {
		_, f, l, _ := runtime.Caller(0)

		logging <- datatypes.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", err.Error(), f, l-1),
			MsgType: "error",
		}

		return
	}

	if len(listTmp) == 0 {
		_, f, l, _ := runtime.Caller(0)

		logging <- datatypes.MessageLogging{
			MsgData: fmt.Sprintf("'%s' %s:%d", fmt.Errorf("error decoding the json message, it may be empty"), f, l-1),
			MsgType: "error",
		}

		return
	}

	_ = ProcessingReflectMap(logging, cstixf, listTmp, &listRule, 0, "")

	if listRule.Rules.Passany {
		isAllowed = true
	} else {
		//проверяем соответствие сообщения правилам из раздела Pass
		for _, v := range listRule.Rules.Pass {
			skipMsg := true
			for _, value := range v.ListAnd {
				if !value.StatementExpression {
					skipMsg = false

					break
				}
			}

			if skipMsg {
				isAllowed = true

				break
			}
		}
	}

	//останавливаем обработчик формирующий STIX формат
	cmispfDone <- isAllowed
}
