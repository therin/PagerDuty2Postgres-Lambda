package main

import (
	"github.com/PagerDuty/go-pagerduty"
	"reflect"
	"testing"
)

func TestExtractEscalationRulesUser(t *testing.T) {

	var testTargets []pagerduty.APIObject
	testTarget1 := pagerduty.APIObject{}
	testTarget1.ID = "TestUserID"
	testTarget1.Type = "user_reference"

	var testRules []pagerduty.EscalationRule
	testRule1 := pagerduty.EscalationRule{}
	testRule1.ID = "TestRuleID"
	testRule1.Delay = 2
	testRule1.Targets = append(testTargets, testTarget1)
	testRules = append(testRules, testRule1)

	var result = ExtractEscalationRulesUser(testRules)

	assertEqual(t, result[0].ID, "TestRuleIDTestUserID")
	assertEqual(t, result[0].RuleID, "TestRuleID")
	assertEqual(t, result[0].UserID, "TestUserID")
}

func TestExtractEscalationRulesSchedule(t *testing.T) {

	var testTargets []pagerduty.APIObject
	testTarget1 := pagerduty.APIObject{}
	testTarget1.ID = "TestScheduleID"
	testTarget1.Type = "schedule_reference"

	var testRules []pagerduty.EscalationRule
	testRule1 := pagerduty.EscalationRule{}
	testRule1.ID = "TestRuleID"
	testRule1.Delay = 2
	testRule1.Targets = append(testTargets, testTarget1)
	testRules = append(testRules, testRule1)

	var result = ExtractEscalationRulesSchedule(testRules)

	assertEqual(t, result[0].ID, "TestRuleIDTestScheduleID")
	assertEqual(t, result[0].RuleID, "TestRuleID")
	assertEqual(t, result[0].ScheduleID, "TestScheduleID")
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
