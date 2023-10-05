package testdecodecasetostix_test

import (
	"fmt"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"ISEMS-MRSICT/mytest/testdecodecasetostix/datatypes"
	"ISEMS-MRSICT/mytest/testdecodecasetostix/handlerdecoder"
	"ISEMS-MRSICT/mytest/testdecodecasetostix/rules"
	"ISEMS-MRSICT/mytest/testdecodecasetostix/supportfunc"
)

var _ = Describe("Decodecase", Ordered, func() {
	var (
		logging               chan datatypes.MessageLogging
		cmispfDone            chan bool
		chanCreateSTIXFormat  chan datatypes.ChanInputCreateSTIXFormat
		listRule              rules.ListRulesProcessingMsgMISP
		exampleByte           []byte
		errReadFile, errRules error
	)

	BeforeAll(func() {
		logging = make(chan datatypes.MessageLogging)
		cmispfDone = make(chan bool)
		chanCreateSTIXFormat = make(chan datatypes.ChanInputCreateSTIXFormat)

		go func() {
			fmt.Println("___ Logging START")
			defer fmt.Println("___ Logging STOP")

			for log := range logging {
				fmt.Println("----", log, "----")
			}
		}()

		go func() {
			for d := range chanCreateSTIXFormat {
				fmt.Println("--- FieldBranch:", d.FieldBranch)
				fmt.Println("--- Value:", d.Value)
				fmt.Println("")
			}
		}()

		//читаем тестовый файл
		exampleByte, errReadFile = supportfunc.ReadFileJson("mytest/testdecodecasetostix", "example_case_1.json")

		// инициализируем модуль чтения правил обработки MISP сообщений
		listRule, _, errRules = rules.GetRuleProcessingMsgForMISP("mytest/testdecodecasetostix/rules", "rule_1.yaml")

		handlerdecoder.NewHandlerDecoder(exampleByte, listRule, chanCreateSTIXFormat, cmispfDone, logging)
	})

	Context("Тест 1. Проверка чтения тестового файла json", func() {
		It("При чтении тестового файла ошибки быть не должно", func() {

			Expect(errReadFile).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 2. Проверка декодирования тестового json файла", func() {
		It("При декодировании тестового файла ошибки быть не должно", func() {

			fmt.Println("RULES:", listRule)

			Expect(errRules).ShouldNot(HaveOccurred())
		})
	})

	Context("Тест 3. Проверка отработки декодера json сообщения", func() {
		It("Должен быть получен сигнал останов через канал cmispfDone", func() {
			isAllowed := <-cmispfDone

			fmt.Println("isAllowed: ", isAllowed)

			Expect(isAllowed).Should(BeTrue())
		})
	})
})
