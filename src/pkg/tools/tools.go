package tools

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"log"
	"os"
	"strconv"
	"time"
)

type AWS struct {
	session     *session.Session
	Environment string
	acreds      *credentials.Credentials
	assumedRole bool
}

type EnvVariables struct {
	PagerDutySubdomain        string
	PagerDutyApiKey           string
	DatabaseEndpoint          string
	DatabaseName              string
	DatabaseUserName          string
	DatabasePasswordParameter string
	DatabasePassword          string
	PaginationLimit           uint
	IncrementalBuffer         int
	IncrementalWindow         int
	PagerDutyEpoch            time.Time
}

type EscalationsPolicy struct {
	APIObject pagerduty.APIObject
	Name      string `API:"Name" DB:"name"`
	NumLoops  uint   `API:"NumLoops" DB:"num_loops"`
}

// Integration in PD API
type Service struct {
	APIObject pagerduty.APIObject
	Name      string `API:"Name" DB:"name"`
	Status    string `API:"Status" DB:"status"`
}

type EscalationsRule struct {
	ID         string `API:"ID" DB:"id"`
	PolicyID   string `API:"N/A" DB:"escalation_policy_id" `
	Delay      uint   `API:"Name" DB:"escalation_delay_in_minutes" `
	LevelIndex int    `API:"LevelIndex" DB:"level_index"`
}

type EscalationsRuleUser struct {
	ID     string `API:"ID" DB:"id"`
	RuleID string `API:"N/A" DB:"escalation_rule_id" `
	UserID string `API:"N/A" DB:"user_id"`
}

type EscalationsRuleSchedule struct {
	ID         string `API:"ID" DB:"id"`
	RuleID     string `API:"N/A" DB:"escalation_rule_id" `
	ScheduleID string `API:"N/A" DB:"schedule_id"`
}

type User struct {
	APIObject pagerduty.APIObject
	Name      string `API:"Name" DB:"name"`
	Email     string `API:"Email" DB:"email"`
}

type Incident struct {
	APIObject            pagerduty.APIObject
	IncidentNumber       uint   `API:"IncidentNumber" DB:"incident_number"`
	CreatedAt            string `API:"CreatedAt" DB:"created_at"`
	IncidentKey          string `API:"IncidentKey" DB:"incident_key"`
	Service              pagerduty.APIObject
	EscalationPolicy     pagerduty.APIObject
	FirstTriggerLogEntry pagerduty.APIObject
}

type LogEntry struct {
	APIObject pagerduty.APIObject
	CreatedAt string `API:"CreatedAt" DB:"created_at"`
	Incident  pagerduty.Incident
	Agent     pagerduty.Agent
	Channel   pagerduty.Channel
	Teams     []pagerduty.Team
}

type Schedule struct {
	APIObject pagerduty.APIObject
	Name      string `API:"Name" DB:"name"`
}

type UserSchedule struct {
	ID         string `API:"ID" DB:"id"`
	UserID     string `API:"N/A" DB:"user_id" `
	ScheduleID string `API:"N/A" DB:"schedule_id"`
}

// Initialize a new struct

var EnvironmentVariables = new(EnvVariables)

func PopulateEnvVariables() {

	// Populate struct content
	var err error

	EnvironmentVariables.PagerDutySubdomain = os.Getenv("PAGERDUTY_SUBDOMAIN")
	EnvironmentVariables.PagerDutyApiKey = os.Getenv("PAGERDUTY_API_KEY")
	EnvironmentVariables.DatabaseEndpoint = os.Getenv("DATABASE_URL")
	EnvironmentVariables.DatabaseName = os.Getenv("DATABASE_NAME")
	EnvironmentVariables.DatabaseUserName = os.Getenv("DATABASE_USER_NAME")
	EnvironmentVariables.DatabasePasswordParameter = os.Getenv("DATABASE_PASSWORD_PARAMETER")
	PaginationLimitInt, err := strconv.Atoi(os.Getenv("PAGINATION_LIMIT"))
	if err != nil {
		fmt.Println(err)
	}
	EnvironmentVariables.PaginationLimit = uint(PaginationLimitInt)
	EnvironmentVariables.IncrementalBuffer, err = strconv.Atoi(os.Getenv("INCREMENTAL_BUFFER"))
	if err != nil {
		fmt.Println(err)
	}
	EnvironmentVariables.IncrementalWindow, err = strconv.Atoi(os.Getenv("INCREMENTAL_WINDOW"))
	if err != nil {
		fmt.Println(err)
	}
	EnvironmentVariables.PagerDutyEpoch, err = time.Parse(time.RFC3339, os.Getenv("PAGERDUTY_EPOCH"))
	if err != nil {
		fmt.Println(err)
	}

	// create AWS object and retrieve SSM parameter

	var AWSSession AWS
	AWSSession.Open("region-id")
	defer AWSSession.Close()

	fmt.Println("Retrieving SSM Parameter Store entry")
	response, err := AWSSession.GetParameterValue(EnvironmentVariables.DatabasePasswordParameter)

	if err != nil {
		log.Fatal(err)
	}

	EnvironmentVariables.DatabasePassword = *response

}

//Open starts AWS connections
func (a *AWS) Open(region string) error {
	sess, err := session.NewSessionWithOptions(session.Options{
		Config: aws.Config{Region: aws.String(region)},
	})
	if err != nil {
		log.Fatal(err)
	}

	a.session = sess

	return err
}

//Close destroys all AWS connecitons
func (a *AWS) Close() {
	a.session = nil
}

// Get SSM Parameter Store parameter
func (a *AWS) GetParameters(paramName string) (*ssm.GetParametersOutput, error) {

	svc := ssm.New(a.session)

	params := &ssm.GetParametersInput{
		Names: []*string{
			aws.String(paramName),
		},
		WithDecryption: aws.Bool(true),
	}
	resp, err := svc.GetParameters(params)

	if err != nil {
		return nil, err
	}

	return resp, nil
}

func (a *AWS) GetParameterValue(paramName string) (*string, error) {
	r, err := a.GetParameters(paramName)
	if err != nil {
		return nil, err
	}

	params := r.Parameters[0]

	return params.Value, err
}
