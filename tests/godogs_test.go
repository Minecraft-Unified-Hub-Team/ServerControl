package tests

import (
	"context"
	"testing"

	"github.com/cucumber/godog"
	"github.com/sirupsen/logrus"
)

type logWriterType struct{ t *testing.T }

func (l logWriterType) Write(p []byte) (n int, err error) {
	l.t.Logf(string(p))
	return len(p), nil
}

func TestFeatures(t *testing.T) {
	logrus.SetOutput(logWriterType{t: t})
	fm, err := NewFeatureManager(context.Background())
	if err != nil {
		logrus.Panic(err)
	}

	suite := godog.TestSuite{
		ScenarioInitializer: func(sc *godog.ScenarioContext) {
			InitializeScenario(fm, sc)
		},
		Options: &godog.Options{
			Format:        "pretty",
			Paths:         []string{"features"},
			TestingT:      t,
			StopOnFailure: true,
		},
	}
	t.Cleanup(func() {
		fm.StepCleanup(context.Background())
	})
	if suite.Run() != 0 {
		logrus.Info("non-zero status returned, failed to run feature tests")
	}
}

func InitializeScenario(fm *FeatureManager, sc *godog.ScenarioContext) {
	sc.Step(`^ServerControl is up$`, fm.serverControlIsUp)

	sc.Step(`^I connect to service control$`, fm.iConnectToServiceControl)
	sc.Step(`^I install "([^"]*)" server version$`, fm.iInstallServer)

	sc.Step(`^I start server$`, fm.iStartServer)

	sc.Step(`^I ping to the server$`, fm.iPingToTheServer)

	sc.Step(`^I have the error "([^"]*)"$`, fm.iHaveAnError)
	sc.Step(`^I have no errors$`, fm.iHaveNoErrors)

	sc.Step(`^I set the option "([^"]*)" with value "([^"]*)"$`, fm.iSetOption)
	sc.Step(`^option "([^"]*)" equal to "([^"]*)"$`, fm.optionEqualTo)
}
