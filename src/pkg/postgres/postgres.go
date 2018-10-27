package postgres

import (
	"../tools"
	"database/sql"
	"fmt"
	"github.com/lib/pq"
	"time"
)

// Populate database connection string structure

type User struct {
	Id    int
	Name  string
	Email string
}

type ReportingStore interface {
	AllUsers() []*User
	UpdateEscalationPolicies(tools.EscalationsPolicy)
	UpdateEscalationRules(tools.EscalationsRule)
	UpdateEscalationRuleUsers(tools.EscalationsRuleUser)
	UpdateEscalationRuleSchedules(tools.EscalationsRuleSchedule)
	UpdateUsers(tools.User)
	UpdateServices(tools.Service)
	UpdateSchedules(tools.Schedule)
	UpdateUserSchedules(tools.UserSchedule)
	UpdateIncidents(tools.Incident)
	UpdateLogEntries(tools.LogEntry)
	CalcLastIncidentRecordDate() time.Time
	CalcLastLogEntryRecordDate() time.Time
	TruncateTable(string)
}

// Implements a custome DB type, gives us an option to mock DB connections
type DB struct {
	*sql.DB
}

// Open DB connection
func DatabaseConnect(DataSourceName string) *DB {

	db, err := sql.Open("postgres", DataSourceName)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return &DB{db}

}

// Get all users from reporting database

func (db *DB) UpdateEscalationPolicies(input tools.EscalationsPolicy) {

	sqlStatement := `
	INSERT INTO escalation_policies (Id, name, num_loops)
	VALUES ($1, $2, $3)`

	_, err := db.Exec(sqlStatement, input.APIObject.ID, input.Name, input.NumLoops)
	if err != nil {
		panic(err)
	}
}

func (db *DB) UpdateEscalationRules(input tools.EscalationsRule) {

	sqlStatement := `
	INSERT INTO escalation_rules (Id, escalation_policy_id, escalation_delay_in_minutes, level_index)
	VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(sqlStatement, input.ID, input.PolicyID, input.Delay, input.LevelIndex)
	if err != nil {
		panic(err)
	}
}

func (db *DB) UpdateEscalationRuleUsers(input tools.EscalationsRuleUser) {

	sqlStatement := `
	INSERT INTO escalation_rule_users (Id, escalation_rule_id, user_id)
	VALUES ($1, $2, $3)`

	_, err := db.Exec(sqlStatement, input.ID, input.RuleID, input.UserID)
	if err != nil {
		panic(err)
	}
}

func (db *DB) UpdateEscalationRuleSchedules(input tools.EscalationsRuleSchedule) {

	sqlStatement := `
	INSERT INTO escalation_rule_schedules (Id, escalation_rule_id, schedule_id)
	VALUES ($1, $2, $3)`

	_, err := db.Exec(sqlStatement, input.ID, input.RuleID, input.ScheduleID)
	if err != nil {
		panic(err)
	}
}

func (db *DB) UpdateSchedules(input tools.Schedule) {

	sqlStatement := `
	INSERT INTO schedules (Id, name)
	VALUES ($1, $2)`

	_, err := db.Exec(sqlStatement, input.APIObject.ID, input.Name)
	if err != nil {
		panic(err)
	}
}

func (db *DB) UpdateUserSchedules(input tools.UserSchedule) {

	sqlStatement := `
	INSERT INTO user_schedule (Id, user_id, schedule_id)
	VALUES ($1, $2, $3)`

	_, err := db.Exec(sqlStatement, input.ID, input.UserID, input.ScheduleID)
	if err != nil {
		panic(err)
	}
}

func (db *DB) UpdateServices(input tools.Service) {

	sqlStatement := `
	INSERT INTO services (Id, name, status, type)
	VALUES ($1, $2, $3, $4)`

	_, err := db.Exec(sqlStatement, input.APIObject.ID, input.Name, input.Status, input.APIObject.Type)
	if err != nil {
		panic(err)
	}
}

func (db *DB) UpdateUsers(input tools.User) {

	sqlStatement := `
	INSERT INTO users (Id, name, email)
	VALUES ($1, $2, $3)`

	_, err := db.Exec(sqlStatement, input.APIObject.ID, input.Name, input.Email)
	if err != nil {
		panic(err)
	}
}

func (db *DB) UpdateIncidents(input tools.Incident) {

	sqlStatement := `
	INSERT INTO incidents (Id, incident_number, created_at, html_url, incident_key, service_id,
		escalation_policy_id, trigger_summary_subject, trigger_summary_description, trigger_type)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := db.Exec(sqlStatement, input.APIObject.ID, input.IncidentNumber, input.CreatedAt, input.APIObject.HTMLURL,
		input.IncidentKey, input.Service.ID, input.EscalationPolicy.ID, input.FirstTriggerLogEntry.Summary,
		input.FirstTriggerLogEntry.Self, input.FirstTriggerLogEntry.Type)

	if err != nil {
		fmt.Println(err)
		// panic(err)
	}
}

func (db *DB) UpdateLogEntries(input tools.LogEntry) {

	sqlStatement := `
	INSERT INTO log_entries (Id, type, created_at, incident_id, agent_type, agent_id,
		channel_type, user_id, notification_type, assigned_user_id)
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	// handle a case when log entry has no team assigned
	var assigned_user_id string
	var user_id string

	if len(input.Teams) == 0 {
		assigned_user_id, user_id = "", ""
	} else {
		assigned_user_id, user_id = input.Teams[0].ID, input.Teams[0].ID
	}

	_, err := db.Exec(sqlStatement, input.APIObject.ID, input.APIObject.Type, input.CreatedAt, input.Incident.ID,
		input.Agent.Type, input.Agent.ID, input.Channel.Type, user_id,
		input.APIObject.Type, assigned_user_id)

	if err != nil {
		fmt.Println(err)
	}
}

func (db *DB) TruncateTable(TableName string) {

	TableName = pq.QuoteIdentifier(TableName)
	sqlStatement := fmt.Sprintf("TRUNCATE %v;", TableName)
	res, err := db.Exec(sqlStatement)
	if err != nil {
		panic(err)
	}
	count, err := res.RowsAffected()
	if err != nil {
		panic(err)
	}
	fmt.Println("Truncate complete: ", count)

}

func (db *DB) AllUsers() []*User {
	rows, err := db.Query("SELECT * from users")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	sliceOfUsers := make([]*User, 0)

	for rows.Next() {
		us := new(User)

		// Validate that fetched row contains valid user entry
		err := rows.Scan(&us.Id, &us.Name, &us.Email)
		if err != nil {
			panic(err)
		}

		fmt.Println(us)
		sliceOfUsers = append(sliceOfUsers, us)
	}

	return sliceOfUsers

}

func (db *DB) CalcLastIncidentRecordDate() time.Time {
	/*
		Calculate the point from which we should resume incremental
		updates. Allow a bit of overlap to ensure we don't miss anything.
	*/
	var date string
	var LastRecordedIncidentDate time.Time

	sqlStatement := `SELECT created_at FROM public.incidents ORDER BY 1 DESC LIMIT 1`
	row := db.QueryRow(sqlStatement)
	switch err := row.Scan(&date); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned, falling back to default start date:", tools.EnvironmentVariables.PagerDutyEpoch)
		LastRecordedIncidentDate = tools.EnvironmentVariables.PagerDutyEpoch
	case nil:
		fmt.Println(date)
		LastRecordedIncidentDate, err = time.Parse(time.RFC3339, date)
	default:
		panic(err)
	}

	LastRecordedIncidentDate = LastRecordedIncidentDate.Add(time.Duration(-tools.EnvironmentVariables.IncrementalBuffer) * time.Second)

	return LastRecordedIncidentDate

}

func (db *DB) CalcLastLogEntryRecordDate() time.Time {
	/*
		Calculate the point from which we should resume incremental
		updates. Allow a bit of overlap to ensure we don't miss anything.
	*/
	var date string
	var LastRecordedLogEntryDate time.Time

	sqlStatement := `SELECT created_at FROM public.log_entries ORDER BY 1 DESC LIMIT 1`
	row := db.QueryRow(sqlStatement)
	switch err := row.Scan(&date); err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned, falling back to default start date:", tools.EnvironmentVariables.PagerDutyEpoch)
		LastRecordedLogEntryDate = tools.EnvironmentVariables.PagerDutyEpoch
	case nil:
		fmt.Println(date)
		LastRecordedLogEntryDate, err = time.Parse(time.RFC3339, date)
	default:
		panic(err)
	}

	LastRecordedLogEntryDate = LastRecordedLogEntryDate.Add(time.Duration(-tools.EnvironmentVariables.IncrementalBuffer) * time.Second)

	return LastRecordedLogEntryDate

}
