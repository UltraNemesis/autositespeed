// TestExecutor.go
package autositespeed

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/mholt/archiver"
)

type testExecutorConf struct {
	ConfigFile           string `json:"configFile"`
	TestFileLocation     string `json:"testFileLocation"`
	ReportFileLocation   string `json:"reportFileLocation"`
	ReportBackupDuration string `json:"reportBackupDuration"`
	SiteSpeedVersion     string `json:"sitespeedVersion"`
	SiteSpeedCommand     string `json:"sitespeedCommand"`
}

type testExecutor struct {
	conf                *testExecutorConf
	currentRunTimestamp time.Time
	running             bool
}

func newTestExecutor(conf *testExecutorConf) *testExecutor {
	return &testExecutor{
		conf: conf,
	}
}

func (te *testExecutor) Run() {
	if !te.running {
		te.running = true
		te.removeOldReports()
		pattern := te.conf.TestFileLocation + "/*.json"
		log.Println("Enumerating TestSuite Files with pattern : ", pattern)

		files, _ := filepath.Glob(pattern)

		log.Println("TestSuite Files found : ", len(files))

		te.currentRunTimestamp = time.Now().UTC()

		for _, file := range files {
			log.Println("Trying to load TestSuite File : ", file)
			testSuite, err := ReadTestSuiteFromFile(file)

			if testSuite != nil {
				log.Println("Loading Successful")
				te.executeTestSuite(testSuite)
			} else {
				log.Println("Loading Failed with Error : ", err)
			}
		}
		te.running = false
	} else {
		log.Println("Last run not finished yet")
	}
}

func (te *testExecutor) executeTestSuite(ts *TestSuite) {
	log.Println("Trying to execute test suite : ", ts.SuiteName)

	for _, test := range ts.Tests {
		log.Println("Trying to Execute Test [Name = " + test.Name + ", Url = " + test.Url + "]")
		te.executeTest(ts.SuiteName, &test)

		time.Sleep(time.Second * 1)
	}
}

func (te *testExecutor) executeTest(tsName string, test *UrlTest) error {
	resultPath := te.GenerateResultPath(tsName, test.Name)

	log.Println("Trying to create directory : ", resultPath)

	err := os.MkdirAll(resultPath, os.ModePerm)

	if err != nil {
		log.Println(err)
	}

	log.Println("Result Path = ", resultPath)

	parts := strings.Fields(te.conf.SiteSpeedCommand)

	log.Println("parts :", parts)

	args := make([]string, 0)

	pwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))

	args = append(args, "run", "--privileged", "--shm-size=1g", "--rm", "-v", pwd+":/sitespeed.io", te.conf.SiteSpeedVersion)

	args = append(args, "--config", te.conf.ConfigFile)

	if len(test.WptScript) > 0 {
		args = append(args, "--webpagetest.file="+test.WptScript)
	}

	args = append(args, "--webpagetest.label="+tsName+"-"+test.Name)

	if len(test.PreScript) > 0 {
		args = append(args, "--preScript="+test.PreScript)
	}

	args = append(args, "--outputFolder", resultPath)

	args = append(args, test.Url)

	//cmd := exec.Command("docker", args...)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	cmd := exec.CommandContext(ctx, "docker", args...)

	log.Println("Command Line : ", cmd.Args)

	writer := &logWriter{}

	cmd.Stdout = writer
	cmd.Stderr = writer

	result := cmd.Run()

	if result != nil {
		log.Println(result)
	}

	log.Println("Test Successful")

	log.Println("Archiving Test Report")

	zipErr := te.Archive(resultPath)

	if zipErr != nil {
		log.Println(zipErr)
	}

	log.Println("Archiving Successful")

	remErr := os.RemoveAll(resultPath)

	if remErr != nil {
		log.Println(remErr)
	}

	return result
}

func (te *testExecutor) Archive(resultPath string) error {
	files, _ := filepath.Glob(resultPath + "/*")
	zipErr := archiver.Zip.Make(resultPath+".zip", files)

	return zipErr
}

func (te *testExecutor) GenerateResultPath(tsName, testName string) string {
	timestamp := te.currentRunTimestamp.Format("2006-01-02-15-04-05")
	dateStr := te.currentRunTimestamp.Format("2006-01-02")
	timeStr := te.currentRunTimestamp.Format("15-04-05")

	return te.conf.ReportFileLocation + "/" + dateStr + "/" + timeStr + "/" + tsName + "/" + testName + "/" + timestamp + "-" + tsName + "-" + testName
}

func (te *testExecutor) removeOldReports() {
	backupDuration, parseErr := time.ParseDuration(te.conf.ReportBackupDuration)

	if parseErr != nil {
		log.Println("Could not parse backup duration : ", parseErr)
	}

	log.Printf("Remove reports older than %s", backupDuration)
	files, err := filepath.Glob(te.conf.ReportFileLocation + "/*")

	if err != nil {
		log.Println(err)
	} else {
		log.Println("Reports Found : ", files)
	}

	timeNow := time.Now().UTC()

	for _, file := range files {
		_, dir := filepath.Split(file)
		creationTime, parseErr := time.Parse("2006-01-02", dir)

		if parseErr != nil {
			log.Println(parseErr)
		} else {
			diff := timeNow.Sub(creationTime)

			if diff >= backupDuration {
				log.Printf("Report @ %s aged %s ==> Remove", dir, diff)
				remErr := os.RemoveAll(file)

				if remErr != nil {
					log.Println(remErr)
				}
			} else {
				log.Printf("Report @ %s aged %s ==> Keep", dir, diff)
			}
		}
	}
}
