package tools

import (
	"testing"
)

func TestIsValidStruct(t *testing.T) {

	var tests = []struct {
		input    string
		expected bool
	}{
		{"incident", true},
		{"service", true},
		{"escalation_policy", true},
		{"escalation_rule", true},
		{"user", true},
		{"notify_log_entry", true},
		{"trigger_log_entry", true},
		{"acknowledge_log_entry", true},
		{"acknowledge_log_entry_reference", true},
		{"annotate_log_entry", true},
		{"annotate_log_entry_reference", true},
		{"assign_log_entry", true},
		{"assign_log_entry_reference", true},
		{"escalate_log_...", false},
		{"escalate_log_entry_referce", false},
		{"exhaust_escalat", false},
		{"ex", false},
	}

	for _, test := range tests {
		if output := IsValidStruct(test.input); output != test.expected {
			t.Error("Test Failed: {} inputted, {} expected, recieved: {}", test.input, test.expected, output)
		}
	}
}
