// Server.go
package autositespeed

type Configuration struct {
	TestScheduler testSchedulerConf `json:"testScheduler"`
	TestExecutor  testExecutorConf  `json:"testExecutor"`
	Logging       LogOptions        `json:logging`
}

type AutoSiteSpeedServer struct {
	conf          *Configuration
	testScheduler *testScheduler
	testExecutor  *testExecutor
}

func NewServer(conf Configuration) *AutoSiteSpeedServer {
	return &AutoSiteSpeedServer{
		conf:          &conf,
		testScheduler: newTestScheduler(&conf.TestScheduler),
		testExecutor:  newTestExecutor(&conf.TestExecutor),
	}
}

func (ass *AutoSiteSpeedServer) Start() {
	initLogging(&ass.conf.Logging)
	ass.testScheduler.Start(ass.testExecutor)
}
