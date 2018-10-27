package tools

import (
	"github.com/PagerDuty/go-pagerduty"
	"reflect"
	"testing"
)

func TestGetMappedEscalationPolicies(t *testing.T) {
	var EscalationPolicies []pagerduty.EscalationPolicy

	testPolicy := pagerduty.EscalationPolicy{}
	testPolicy.Name = "Test Policy"
	testPolicy.NumLoops = 2
	testPolicy.ID = "Test1"

	EscalationPolicies = append(EscalationPolicies, testPolicy)
	var result = GetMappedEscalationPolicies(EscalationPolicies)

	assertEqual(t, result[0].Name, testPolicy.Name)
}

func TestGetMappedEscalationPoliciesNested(t *testing.T) {

	var testRules []pagerduty.EscalationRule
	testRule1 := pagerduty.EscalationRule{}
	testRule1.ID = "rule1"
	testRule1.Delay = 2
	testRules = append(testRules, testRule1)

	var EscalationPolicies []pagerduty.EscalationPolicy
	testPolicy := pagerduty.EscalationPolicy{}
	testPolicy.Name = "Test Policy"
	testPolicy.NumLoops = 2
	testPolicy.ID = "Test1"
	testPolicy.EscalationRules = testRules
	EscalationPolicies = append(EscalationPolicies, testPolicy)

	var result = GetMappedEscalationPolicies(EscalationPolicies)

	assertEqual(t, result[0].Name, testPolicy.Name)
}

func TestGetMappedEscalationRules(t *testing.T) {

	var testRules []pagerduty.EscalationRule
	testRule1 := pagerduty.EscalationRule{}
	testRule1.ID = "rule1"
	testRule1.Delay = 2
	testRules = append(testRules, testRule1)

	var EscalationPolicies []pagerduty.EscalationPolicy
	testPolicy := pagerduty.EscalationPolicy{}
	testPolicy.Name = "Test Policy"
	testPolicy.NumLoops = 2
	testPolicy.ID = "Test1"
	testPolicy.EscalationRules = testRules
	EscalationPolicies = append(EscalationPolicies, testPolicy)

	var result = GetMappedEscalationRules(testRules, testPolicy.ID)

	assertEqual(t, result[0].ID, testRule1.ID)
	assertEqual(t, result[0].PolicyID, testPolicy.ID)
	assertEqual(t, 0, result[0].LevelIndex)
}

func TestGetMappedSchedules(t *testing.T) {
	var PagerDutySchedules []pagerduty.Schedule

	testSchedule := pagerduty.Schedule{
		Name: "Infra",
		APIObject: pagerduty.APIObject{
			ID: "schedule1",
		},
	}

	PagerDutySchedules = append(PagerDutySchedules, testSchedule)

	var result = GetMappedSchedules(PagerDutySchedules)

	assertEqual(t, result[0].Name, testSchedule.Name)
	assertEqual(t, result[0].APIObject.ID, testSchedule.APIObject.ID)
}

func TestGetMappedServices(t *testing.T) {

	var PagerDutyServices []pagerduty.Service

	testService := pagerduty.Service{
		Name:        "PD Service",
		Description: "grumble grumble",
		APIObject: pagerduty.APIObject{
			ID:   "Service1",
			Type: "Type1",
		},
	}

	PagerDutyServices = append(PagerDutyServices, testService)

	var result = GetMappedServices(PagerDutyServices)

	assertEqual(t, result[0].Name, testService.Name)
	assertEqual(t, result[0].APIObject.ID, testService.APIObject.ID)
}

func TestGetMappedIncidents(t *testing.T) {

	var PagerDutyIncidents []pagerduty.Incident

	testIncident := pagerduty.Incident{
		IncidentNumber: 1,
		IncidentKey:    "grumble grumble",
		APIObject: pagerduty.APIObject{
			ID:   "Incident11",
			Type: "Disaster",
		},
		FirstTriggerLogEntry: pagerduty.APIObject{
			ID:   "trig1",
			Type: "info",
		},
	}

	PagerDutyIncidents = append(PagerDutyIncidents, testIncident)

	var result = GetMappedIncidents(PagerDutyIncidents)

	assertEqual(t, result[0].IncidentNumber, testIncident.IncidentNumber)
	assertEqual(t, result[0].APIObject.ID, testIncident.APIObject.ID)
	assertEqual(t, result[0].FirstTriggerLogEntry.ID, testIncident.FirstTriggerLogEntry.ID)
}

func TestGetMappedLogEntries(t *testing.T) {

	var PagerDutyLogEntries []pagerduty.LogEntry

	testIncident := pagerduty.Incident{
		IncidentNumber: 1,
		IncidentKey:    "grumble grumble",
		APIObject: pagerduty.APIObject{
			ID:   "Incident11",
			Type: "Disaster",
		},
		FirstTriggerLogEntry: pagerduty.APIObject{
			ID:   "trig1",
			Type: "info",
		},
	}

	testLogEntry := pagerduty.LogEntry{
		CreatedAt: "01-Jun-2020",
		APIObject: pagerduty.APIObject{
			ID:   "log1",
			Type: "trigger_log_entry",
		},
		Agent: pagerduty.Agent{
			ID:   "agent001",
			Type: "super_user",
		},
		Channel:  pagerduty.Channel{Type: "incident_channel"},
		Incident: testIncident,
	}

	PagerDutyLogEntries = append(PagerDutyLogEntries, testLogEntry)

	var result = GetMappedLogEntries(PagerDutyLogEntries)

	assertEqual(t, result[0].CreatedAt, testLogEntry.CreatedAt)
	assertEqual(t, result[0].APIObject.ID, testLogEntry.APIObject.ID)
	assertEqual(t, result[0].APIObject.Type, testLogEntry.APIObject.Type)
}

func assertEqual(t *testing.T, e, g interface{}) (r bool) {
	r = compare(e, g)
	if !r {
		t.Errorf("Expected [%v], got [%v]", e, g)
	}

	return
}

func compare(e, g interface{}) (r bool) {
	ev := reflect.ValueOf(e)
	gv := reflect.ValueOf(g)

	if ev.Kind() != gv.Kind() {
		return
	}

	switch ev.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		r = (ev.Int() == gv.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		r = (ev.Uint() == gv.Uint())
	case reflect.Float32, reflect.Float64:
		r = (ev.Float() == gv.Float())
	case reflect.String:
		r = (ev.String() == gv.String())
	case reflect.Bool:
		r = (ev.Bool() == gv.Bool())
	case reflect.Slice, reflect.Map:
		r = reflect.DeepEqual(e, g)
	}

	return
}
