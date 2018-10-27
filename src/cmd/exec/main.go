package main

import (
	"../../pkg/pagerdutysvc"
	"../../pkg/postgres"
	"../../pkg/tools"
	"context"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"github.com/aws/aws-lambda-go/lambda"
	"time"
)

type MyEvent struct {
	Name string `json:"name"`
}

// links reportingstore interface
type Env struct {
	db postgres.ReportingStore
}

func HandleRequest(ctx context.Context, name MyEvent) (string, error) {
	return fmt.Sprintf("Data transfer completed successfully"), nil
}

func TransferEscalationPolicies(env *Env) {

	EscalationsPolicies := pagerdutysvc.GetPagerDutyEscalationPolicies()
	MappedEscalationPolicies := tools.GetMappedEscalationPolicies(EscalationsPolicies)

	env.db.TruncateTable("escalation_policies")

	for i := range MappedEscalationPolicies {
		env.db.UpdateEscalationPolicies(MappedEscalationPolicies[i])
	}
}

func TransferSchedules(env *Env) {
	Schedules := pagerdutysvc.GetPagerDutySchedules()
	MappedSchedules := tools.GetMappedSchedules(Schedules)

	MappedUserSchedules := []tools.UserSchedule{}

	env.db.TruncateTable("schedules")
	env.db.TruncateTable("user_schedule")

	for i := range MappedSchedules {
		env.db.UpdateSchedules(MappedSchedules[i])
	}

	// Loop over mapped schedules again extract user IDs and build UserSchedule mapping

	for i := range Schedules {

		CurrentcheduleID := Schedules[i].APIObject.ID

		// Loop over user struct in schedule  and extract User ID
		for y := range Schedules[i].Users {

			CurrentUserSchedule := tools.UserSchedule{}

			CurrentUserSchedule.ScheduleID = CurrentcheduleID
			CurrentUserSchedule.UserID = Schedules[i].Users[y].ID
			CurrentUserSchedule.ID = CurrentUserSchedule.UserID + CurrentUserSchedule.ScheduleID

			MappedUserSchedules = append(MappedUserSchedules, CurrentUserSchedule)

		}

	}

	for i := range MappedUserSchedules {
		env.db.UpdateUserSchedules(MappedUserSchedules[i])
	}

}

func TransferEscalationRules(env *Env) {

	env.db.TruncateTable("escalation_rules")
	env.db.TruncateTable("escalation_rule_schedules")
	env.db.TruncateTable("escalation_rule_users")

	// Retrieve escalation policies
	EscalationsPolicies := pagerdutysvc.GetPagerDutyEscalationPolicies()
	EscalationsRulesSlice := []pagerduty.EscalationRule{}
	var MappedEscalationRules = []tools.EscalationsRule{}

	// Map Escalation Rules to Escalation Policy
	for i := range EscalationsPolicies {
		EscalationsPolicyID := EscalationsPolicies[i].APIObject.ID

		EscalationsRules := pagerdutysvc.GetPagerDutyEscalationRule(EscalationsPolicyID)

		// Append API response to slice for future use
		EscalationsRulesSlice = append(EscalationsRulesSlice, EscalationsRules...)

		MappedEscalationRules = append(tools.GetMappedEscalationRules(EscalationsRules, EscalationsPolicyID), MappedEscalationRules...)

	}

	for i := range MappedEscalationRules {
		env.db.UpdateEscalationRules(MappedEscalationRules[i])
	}

	// Map Escalation Rules to User IDs and Schedule IDs
	EscalationRuleUserStruct := ExtractEscalationRulesUser(EscalationsRulesSlice)
	TransferEscalationRulesUser(env, EscalationRuleUserStruct)

	EscalationRuleScheduleStruct := ExtractEscalationRulesSchedule(EscalationsRulesSlice)
	TransferEscalationRulesSchedule(env, EscalationRuleScheduleStruct)

}

func ExtractEscalationRulesUser(EscalationRules []pagerduty.EscalationRule) []tools.EscalationsRuleUser {

	EscalationRuleUsers := []tools.EscalationsRuleUser{}

	for i := range EscalationRules {
		currentRuleUser := tools.EscalationsRuleUser{}
		currentRuleUser.RuleID = EscalationRules[i].ID

		// Loop over targets in Escalation Rule and extract User ID
		for y := range EscalationRules[i].Targets {

			if EscalationRules[i].Targets[y].Type == "user_reference" {

				currentRuleUser.UserID = EscalationRules[i].Targets[y].ID
				currentRuleUser.ID = currentRuleUser.RuleID + currentRuleUser.UserID

				EscalationRuleUsers = append(EscalationRuleUsers, currentRuleUser)

			}

		}
	}

	return EscalationRuleUsers

}

func TransferEscalationRulesUser(env *Env, EscalationRuleUsers []tools.EscalationsRuleUser) {

	for i := range EscalationRuleUsers {
		// fmt.Printf("%+v\n", EscalationRuleUsers[i])
		env.db.UpdateEscalationRuleUsers(EscalationRuleUsers[i])
	}

}

func ExtractEscalationRulesSchedule(EscalationRules []pagerduty.EscalationRule) []tools.EscalationsRuleSchedule {

	EscalationRuleSchedules := []tools.EscalationsRuleSchedule{}

	for i := range EscalationRules {
		currentRuleSchedule := tools.EscalationsRuleSchedule{}
		currentRuleSchedule.RuleID = EscalationRules[i].ID

		// Loop over targets in Escalation Rule and extract Schedule ID
		for y := range EscalationRules[i].Targets {

			if EscalationRules[i].Targets[y].Type == "schedule_reference" {

				currentRuleSchedule.ScheduleID = EscalationRules[i].Targets[y].ID
				currentRuleSchedule.ID = currentRuleSchedule.RuleID + currentRuleSchedule.ScheduleID

				EscalationRuleSchedules = append(EscalationRuleSchedules, currentRuleSchedule)

			}

		}
	}

	return EscalationRuleSchedules

}

func TransferEscalationRulesSchedule(env *Env, EscalationRuleSchedules []tools.EscalationsRuleSchedule) {

	for i := range EscalationRuleSchedules {
		// fmt.Printf("%+v\n", EscalationRuleSchedules[i])
		env.db.UpdateEscalationRuleSchedules(EscalationRuleSchedules[i])
	}

}

func TransferUsers(env *Env) {
	Users := pagerdutysvc.GetPagerDutyUsers()
	MappedUsers := tools.GetMappedUsers(Users)
	env.db.TruncateTable("users")

	for i := range MappedUsers {
		env.db.UpdateUsers(MappedUsers[i])
	}
}

func TransferServices(env *Env) {
	Services := pagerdutysvc.GetPagerDutyServices()
	MappedServices := tools.GetMappedServices(Services)
	env.db.TruncateTable("services")

	for i := range MappedServices {
		env.db.UpdateServices(MappedServices[i])
	}
}

func TransferIncidents(env *Env) {

	/*
		Update data in windowed time chunks. This will give us manageable
		amounts of data request from the API coherently.
		while latest < Time.now
		through = latest + INCREMENTAL_WINDOW
		log("refresh_incremental.window", collection: collection, since: since.iso8601, through: through.iso8601)
	*/

	dateFrom := env.db.CalcLastIncidentRecordDate()
	dateTo := dateFrom.Add(time.Duration(tools.EnvironmentVariables.IncrementalWindow) * time.Second)

	for time.Now().After(dateFrom) {
		Incidents := pagerdutysvc.GetPagerDutyIncidents(dateFrom, dateTo)
		MappedIncidents := tools.GetMappedIncidents(Incidents)

		for i := range MappedIncidents {
			env.db.UpdateIncidents(MappedIncidents[i])
		}

		dateFrom = dateTo
		dateTo = dateFrom.Add(time.Duration(tools.EnvironmentVariables.IncrementalWindow) * time.Second)
	}
}

func TransferLogEntries(env *Env) {

	/*
		Update data in windowed time chunks. This will give us manageable
		amounts of data request from the API coherently.
		while latest < Time.now
		through = latest + INCREMENTAL_WINDOW
		log("refresh_incremental.window", collection: collection, since: since.iso8601, through: through.iso8601)
	*/

	dateFrom := env.db.CalcLastLogEntryRecordDate()
	dateTo := dateFrom.Add(time.Duration(tools.EnvironmentVariables.IncrementalWindow) * time.Second)

	for time.Now().After(dateFrom) {
		LogEntries := pagerdutysvc.GetPagerDutyLogEntries(dateFrom, dateTo)
		MappedLogEntries := tools.GetMappedLogEntries(LogEntries)

		for i := range MappedLogEntries {
			env.db.UpdateLogEntries(MappedLogEntries[i])
		}

		dateFrom = dateTo
		dateTo = dateFrom.Add(time.Duration(tools.EnvironmentVariables.IncrementalWindow) * time.Second)
	}
}

func main() {

	// Retreive environment variables

	tools.PopulateEnvVariables()

	// Make the handler available for Remote Procedure Call by AWS Lambda
	// Get variables for database connection

	var Host = tools.EnvironmentVariables.DatabaseEndpoint
	var Dbname = tools.EnvironmentVariables.DatabaseName
	var User = tools.EnvironmentVariables.DatabaseUserName
	var Password = tools.EnvironmentVariables.DatabasePassword

	// // Build connection string
	ConnectionString := fmt.Sprintf("host=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		Host, User, Password, Dbname)

	db := postgres.DatabaseConnect(ConnectionString)

	// Instantiate env struct with pointer to db connections, pass DB connection as parameter
	env := &Env{db}

	TransferEscalationPolicies(env)
	TransferUsers(env)
	TransferSchedules(env)
	TransferServices(env)
	TransferEscalationRules(env)
	TransferLogEntries(env)
	TransferIncidents(env)

	lambda.Start(HandleRequest)
}
