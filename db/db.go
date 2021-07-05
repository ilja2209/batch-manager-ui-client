package db

import (
	"context"
	"fmt"

	"github.com/go-pg/pg/v10"
)

type Process struct {
	tableName struct{} `pg:"batch_process"`

	Id       int64  `pg:",pk" json:"id"`
	Data     string `json:"data"`
	State    string `json:"state"`
	Started  string `json:"started"`
	Finished string `json:"finished"`
	Context  string `json:"context"`
	Tasks    []Task `pg:"rel:has-many,join_fk:batch_process_id" json:"tasks"`
}

type Task struct {
	tableName struct{} `pg:"task"`

	Id             int64  `pg:",pk" json:"id"`
	Data           string `json:"data"`
	BatchProcessId int64  `json:"process_id"`
	ApplicationId  string `json:"application_id"`
	Attempt        int32  `json:"attempt"`
	MaxAttempt     int32  `json:"max_attempt"`
	Started        string `json:"started"`
	Finished       string `json:"finished"`
	Result         string `json:"result"`
	StageNum       int32  `json:"stage_num"`
	State          string `json:"state"`
	Context        string `json:"context"`
}

type dbLogger struct{}

func (d dbLogger) BeforeQuery(c context.Context, q *pg.QueryEvent) (context.Context, error) {
	return c, nil
}

func (d dbLogger) AfterQuery(c context.Context, q *pg.QueryEvent) error {
	b, _ := q.FormattedQuery()
	fmt.Println(string(b))
	return nil
}

type ProcessRepository struct {
	dbCon *pg.DB
}

const defaultLimit = 50

func NewDatabase(addr string, user string, psswd string, dbName string) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     addr,
		User:     user,
		Password: psswd,
		Database: dbName,
	})

	db.AddQueryHook(dbLogger{})

	ctx := context.Background()

	_, err := db.ExecContext(ctx, "SELECT 1")
	return db, err
}

func NewProcessRepository(dbCon *pg.DB) *ProcessRepository {
	return &ProcessRepository{
		dbCon: dbCon,
	}
}

func (processRepository *ProcessRepository) GetProcesses() ([]Process, error) {
	var processes []Process
	err := processRepository.dbCon.Model(&processes).Order("id DESC").Limit(defaultLimit).Select()
	return processes, err
}

func (processRepository *ProcessRepository) GetProcessById(id int64) (*Process, error) {
	process := new(Process)

	err := processRepository.dbCon.Model(process).
		ColumnExpr("process.id, process.data, process.state, process.started, process.finished, process.context").
		Relation("Tasks"). //todo: Fiusing join. Here 2 queries are generated
		//Join("LEFT JOIN task AS t ON t.batch_process_id = process.id").
		Where("process.id = ?", id).
		First()

	return process, err
}
