package pagerdutysvc

import (
	"../tools"
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
	"time"
)

func GetPagerDutyEscalationPolicies() []pagerduty.EscalationPolicy {

	var EscalationPolicies []pagerduty.EscalationPolicy
	var APIList pagerduty.APIListObject

	// Override default pagination limit
	APIList.Limit = tools.EnvironmentVariables.PaginationLimit

	opts := pagerduty.ListEscalationPoliciesOptions{APIListObject: APIList}

	client := pagerduty.NewClient(tools.EnvironmentVariables.PagerDutyApiKey)

	for {

		eps, err := client.ListEscalationPolicies(opts)

		if err != nil {
			panic(err)
		}

		EscalationPolicies = append(EscalationPolicies, eps.EscalationPolicies...)
		APIList.Offset += tools.EnvironmentVariables.PaginationLimit
		APIList.Limit = tools.EnvironmentVariables.PaginationLimit
		opts = pagerduty.ListEscalationPoliciesOptions{APIListObject: APIList}

		if eps.APIListObject.More != true {
			fmt.Println("Escalation Policies Extracted")

			return EscalationPolicies

		}

	}
	return EscalationPolicies

}

func GetPagerDutyEscalationRule(escID string) []pagerduty.EscalationRule {

	var EscalationRules []pagerduty.EscalationRule

	client := pagerduty.NewClient(tools.EnvironmentVariables.PagerDutyApiKey)

	ers, err := client.ListEscalationRules(escID)

	if err != nil {
		panic(err)
	}

	EscalationRules = ers.EscalationRules

	return EscalationRules
}

func GetPagerDutyUsers() []pagerduty.User {

	var Users []pagerduty.User
	var APIList pagerduty.APIListObject

	// Override default pagination limit
	APIList.Limit = tools.EnvironmentVariables.PaginationLimit

	opts := pagerduty.ListUsersOptions{APIListObject: APIList}

	client := pagerduty.NewClient(tools.EnvironmentVariables.PagerDutyApiKey)

	for {

		usr, err := client.ListUsers(opts)
		if err != nil {
			panic(err)
		}

		Users = append(Users, usr.Users...)
		APIList.Offset += tools.EnvironmentVariables.PaginationLimit
		APIList.Limit = tools.EnvironmentVariables.PaginationLimit
		opts = pagerduty.ListUsersOptions{APIListObject: APIList}

		if usr.APIListObject.More != true {
			fmt.Println("Users Extracted")

			return Users

		}

	}
	return Users
}

func GetPagerDutySchedules() []pagerduty.Schedule {

	var Schedules []pagerduty.Schedule
	var APIList pagerduty.APIListObject

	// Override default pagination limit
	APIList.Limit = tools.EnvironmentVariables.PaginationLimit

	opts := pagerduty.ListSchedulesOptions{APIListObject: APIList}

	client := pagerduty.NewClient(tools.EnvironmentVariables.PagerDutyApiKey)

	for {

		sch, err := client.ListSchedules(opts)
		if err != nil {
			panic(err)
		}

		Schedules = append(Schedules, sch.Schedules...)
		APIList.Offset += tools.EnvironmentVariables.PaginationLimit
		APIList.Limit = tools.EnvironmentVariables.PaginationLimit
		opts = pagerduty.ListSchedulesOptions{APIListObject: APIList}

		if sch.APIListObject.More != true {
			fmt.Println("Schedules Extracted")

			return Schedules

		}

	}
	return Schedules
}

func GetPagerDutyServices() []pagerduty.Service {

	var Services []pagerduty.Service
	var APIList pagerduty.APIListObject

	// Override default pagination limit
	APIList.Limit = tools.EnvironmentVariables.PaginationLimit

	opts := pagerduty.ListServiceOptions{APIListObject: APIList}

	client := pagerduty.NewClient(tools.EnvironmentVariables.PagerDutyApiKey)

	for {

		ser, err := client.ListServices(opts)
		if err != nil {
			panic(err)
		}

		Services = append(Services, ser.Services...)
		APIList.Offset += tools.EnvironmentVariables.PaginationLimit
		APIList.Limit = tools.EnvironmentVariables.PaginationLimit
		opts = pagerduty.ListServiceOptions{APIListObject: APIList}

		if ser.APIListObject.More != true {
			fmt.Println("Services Extracted")

			return Services

		}

	}
	return Services
}

func GetPagerDutyIncidents(dateFrom time.Time, dateTo time.Time) []pagerduty.Incident {

	fmt.Println("Working with:", dateFrom.String(), dateTo.String())

	opts := pagerduty.ListIncidentsOptions{
		Since: dateFrom.String(),
		Until: dateTo.String(),
	}

	var Incidents []pagerduty.Incident
	var APIList pagerduty.APIListObject

	// Override default pagination limit
	APIList.Limit = tools.EnvironmentVariables.PaginationLimit

	client := pagerduty.NewClient(tools.EnvironmentVariables.PagerDutyApiKey)

	for {

		inc, err := client.ListIncidents(opts)
		if err != nil {
			panic(err)
		}

		// fmt.Printf("%+v", inc)

		Incidents = append(Incidents, inc.Incidents...)
		APIList.Offset += tools.EnvironmentVariables.PaginationLimit
		APIList.Limit = tools.EnvironmentVariables.PaginationLimit
		opts = pagerduty.ListIncidentsOptions{APIListObject: APIList, Since: dateFrom.String(),
			Until: dateTo.String()}

		if inc.APIListObject.More != true {
			fmt.Println("Incidents Extracted")

			return Incidents

		}
	}

	inc, err := client.ListIncidents(opts)
	if err != nil {
		panic(err)
	}

	// fmt.Printf("%+v", inc)

	Incidents = append(Incidents, inc.Incidents...)

	return Incidents
}

func GetPagerDutyLogEntries(dateFrom time.Time, dateTo time.Time) []pagerduty.LogEntry {

	fmt.Println("Working with:", dateFrom.String(), dateTo.String())

	opts := pagerduty.ListLogEntriesOptions{
		Since:    dateFrom.String(),
		Until:    dateTo.String(),
		TimeZone: "UTC",
	}

	var LogEntries []pagerduty.LogEntry
	var APIList pagerduty.APIListObject

	// Override default pagination limit
	APIList.Limit = tools.EnvironmentVariables.PaginationLimit

	client := pagerduty.NewClient(tools.EnvironmentVariables.PagerDutyApiKey)

	for {

		log, err := client.ListLogEntries(opts)
		if err != nil {
			panic(err)
		}

		// fmt.Printf("%+v", log)

		LogEntries = append(LogEntries, log.LogEntries...)
		APIList.Offset += tools.EnvironmentVariables.PaginationLimit
		APIList.Limit = tools.EnvironmentVariables.PaginationLimit
		opts = pagerduty.ListLogEntriesOptions{APIListObject: APIList, Since: dateFrom.String(),
			Until: dateTo.String(), TimeZone: "UTC"}

		if log.APIListObject.More != true {
			fmt.Println("Log Entries Extracted")

			return LogEntries

		}
	}

	log, err := client.ListLogEntries(opts)
	if err != nil {
		panic(err)
	}

	LogEntries = append(LogEntries, log.LogEntries...)

	return LogEntries

}
