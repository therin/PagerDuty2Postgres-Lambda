package tools

import (
	"github.com/PagerDuty/go-pagerduty"
	"github.com/jeevatkm/go-model"
)

// GetMappedEscalationPolicies - model.Copy returns a slice of errors if any occur in case we want to do something with them.
func GetMappedEscalationPolicies(policies []pagerduty.EscalationPolicy) []EscalationsPolicy {
	escalationPoliciesToPersist := []EscalationsPolicy{}
	for i := range policies {
		nextPolicy := EscalationsPolicy{}
		model.Copy(&nextPolicy, policies[i])
		escalationPoliciesToPersist = append(escalationPoliciesToPersist, nextPolicy)
	}
	return escalationPoliciesToPersist
}

func GetMappedEscalationRules(rules []pagerduty.EscalationRule, policyID string) []EscalationsRule {
	escalationRulesToPersist := []EscalationsRule{}
	for i := range rules {
		nextRule := EscalationsRule{}
		model.Copy(&nextRule, rules[i])
		nextRule.PolicyID = policyID
		nextRule.LevelIndex = i
		escalationRulesToPersist = append(escalationRulesToPersist, nextRule)
	}
	return escalationRulesToPersist
}

func GetMappedUsers(users []pagerduty.User) []User {
	usersToPersist := []User{}
	for i := range users {
		nextUser := User{}
		model.Copy(&nextUser, users[i])
		usersToPersist = append(usersToPersist, nextUser)
	}
	return usersToPersist
}

func GetMappedSchedules(schedules []pagerduty.Schedule) []Schedule {
	schedulesToPersist := []Schedule{}
	for i := range schedules {
		nextSchedule := Schedule{}
		model.Copy(&nextSchedule, schedules[i])
		schedulesToPersist = append(schedulesToPersist, nextSchedule)
	}
	return schedulesToPersist
}

func GetMappedServices(services []pagerduty.Service) []Service {
	servicesToPersist := []Service{}
	for i := range services {
		nextService := Service{}
		model.Copy(&nextService, services[i])
		servicesToPersist = append(servicesToPersist, nextService)
	}
	return servicesToPersist
}

func GetMappedIncidents(incidents []pagerduty.Incident) []Incident {
	incidentsToPersist := []Incident{}
	for i := range incidents {
		nextIncident := Incident{}
		model.Copy(&nextIncident, incidents[i])
		incidentsToPersist = append(incidentsToPersist, nextIncident)
	}
	return incidentsToPersist
}

func GetMappedLogEntries(logEntries []pagerduty.LogEntry) []LogEntry {
	logEntriesToPersist := []LogEntry{}
	for i := range logEntries {
		nextLogEntry := LogEntry{}
		model.Copy(&nextLogEntry, logEntries[i])
		logEntriesToPersist = append(logEntriesToPersist, nextLogEntry)
	}
	return logEntriesToPersist
}
