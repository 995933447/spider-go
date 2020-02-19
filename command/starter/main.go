package main

import (
	"spider-go/config"
	"spider-go/engine"
	recycler "spider-go/failure/recycle"
	//"spider-go/process/app/javdove/itemworker"
	"spider-go/fetcher"
	"spider-go/logger"
	//airavConfig "spider-go/process/app/airav/config"
	"spider-go/process/app/airav/parser"
	"spider-go/process/app/avhd101/persist"
	"spider-go/process/app/javdove/itemworker"
	//javdoveConfig "spider-go/process/app/javdove/config"
	//javdoveParse "spider-go/process/app/javdove/parser"
	tangccConfig "spider-go/process/app/tangcc/config"
	tangccParser "spider-go/process/app/tangcc/parser"
	tangccPersist "spider-go/process/app/tangcc/persist"
	"spider-go/schedule"
)

func main() {
	registerErrorHandle()
	runEngine(engine.Engine{
			Schedule: &schedule.QueueSchedule{},
			WorkerNum: 3,
			FailedRequestRecycler: recycler.RecycleFailedRequest(),
	})
}

func runEngine(engine config.Engine)  {
	engine.Run(
		//config.Request{
		//	Url: avhd101Config.Host + "/category",
		//	Fetch: fetcher.Html,
		//	Parse: avhd101Parser.CategoryWithSubjects,
		//	ItemChan: persist.CategoryWithSubjects(),
		//	Distinctor: config.NilDistinctor{},
		//},
		//
		config.Request{
			Url:        javdoveConfig.Host + "/categories",
			Fetch:      fetcher.Html,
			Parse:      javdoveParse.Categories,
			ItemChan:   itemworker.Categories(),
			Distinctor: config.NilDistinctor{},
		},
		//
		//config.Request{
		//	Url: julypornConfig.Host + "/videos/amateur/free",
		//	Fetch: julypornFetcher.Categories,
		//	Parse: parser.Categories,
		//	ItemChan: julypornPersist.Categories(),
		//	Distinctor: config.NilDistinctor{},
		//},
		//
		//config.Request{
		//	Url:        airavConfig.Host,
		//	Fetch:      fetcher.Html,
		//	Parse:      parser.Categories,
		//	ItemChan:   config.NilItemWork(),
		//	Distinctor: config.NilDistinctor{},
		//},

		//config.Request{
		//	Url: aiChangeConfig.Host,
		//	Fetch: config.NilFetch,
		//	Parse: aiChangeParse.Categories,
		//	ItemChan: aiChangePersist.Categories(),
		//	Distinctor: config.NilDistinctor{},
		//},

		config.Request{
			Url: tangccConfig.Host,
			Fetch: fetcher.Html,
			Parse: tangccParser.Topic,
			ItemChan: tangccPersist.Subjects(),
			Distinctor: config.NilDistinctor{},
		},
	)
}

func registerErrorHandle() {
	defer func() {
		if err := recover(); err != nil {
			logger.DefaultLogger.Error(err, nil)
		}
		return
	}()
}