package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	Models "gocloudcamp/pkg/models"
	"net/http"
	"time"
)

type ConfigService interface {
	SetConfig(ctx context.Context, req interface{}) (*Models.ConfigRequest, error)
	GetConfig(ctx context.Context, req interface{}) (*Models.ConfigRequest, error)
	UpdConfig(ctx context.Context, req interface{}) (*Models.ConfigRequest, error)
	DelConfig(ctx context.Context, req interface{}) (*Models.ConfigRequest, error)
}

func (svc configService) SetConfig(_ context.Context, req interface{}) (*Models.ConfigRequest, error) {
	r := req.(*Models.ConfigRequest)
	json, err := json.Marshal(r.Data)
	if err != nil {
		return nil, Models.ResponseError{ErrorDescr: "Data marshaling failed"}
	}

	ctx := context.Background()
	tx, err := svc.DB.BeginTx(ctx, nil)
	if err != nil {
		return nil, Models.ResponseError{ErrorDescr: err.Error()}
	}

	row := tx.QueryRow("select coalesce(max(version), 0) from configs where service = $1", r.Service)
	var version int
	err = row.Scan(&version)
	if err != nil {
		tx.Rollback()
		return nil, Models.ResponseError{ErrorDescr: err.Error()}
	}

	if version > 0 {
		_, err = tx.ExecContext(ctx, "update configs set used=false where service = $1 and used = true", r.Service)
		if err != nil {
			tx.Rollback()
			return nil, Models.ResponseError{ErrorDescr: err.Error()}
		}
	}
	version++

	_, err = tx.ExecContext(ctx, "insert into configs (service, version, data) values ($1, $2, $3)", r.Service, version, json)
	if err != nil {
		tx.Rollback()
		return nil, Models.ResponseError{ErrorDescr: err.Error()}
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, Models.ResponseError{ErrorDescr: err.Error()}
	}

	r.Version = version
	return r, nil
}

func (svc configService) GetConfig(_ context.Context, req interface{}) (*Models.ConfigRequest, error) {
	r := req.(*Models.ConfigRequest)
	var rows *sql.Rows
	var err error
	if r.Version == 0 {
		rows, err = svc.DB.Query("select coalesce(id, 0), coalesce(version, 0), used, data from configs where service = $1 and used = true limit 1", r.Service)
	} else {
		rows, err = svc.DB.Query("select coalesce(id, 0), coalesce(version, 0), used, data from configs where service = $1 and version = $2 limit 1", r.Service, r.Version)
	}
	if err != nil {
		return nil, Models.ResponseError{ErrorDescr: err.Error()}
	}
	defer rows.Close()

	var id, version int
	var used bool
	var jsonData []byte
	for rows.Next() {
		err := rows.Scan(&id, &version, &used, &jsonData)
		if err != nil {
			return nil, Models.ResponseError{ErrorDescr: err.Error()}
		}
	}
	if id == 0 {
		return nil, Models.ResponseError{ErrorDescr: "No data on request parameters", Status: http.StatusNotFound}
	}

	err = json.Unmarshal(jsonData, &r.Data)
	if err != nil {
		return nil, Models.ResponseError{ErrorDescr: err.Error()}
	}
	r.Version = version
	r.Used = used
	return r, nil
}

func (svc configService) UpdConfig(_ context.Context, req interface{}) (*Models.ConfigRequest, error) {
	r := req.(*Models.ConfigRequest)
	var rows *sql.Rows
	var err error
	if r.Version == 0 {
		rows, err = svc.DB.Query("select coalesce(id, 0), used from configs where service = $1 and used = true limit 1", r.Service)
	} else {
		rows, err = svc.DB.Query("select coalesce(id, 0), used from configs where service = $1 and version = $2 limit 1", r.Service, r.Version)
	}
	if err != nil {
		return nil, Models.ResponseError{ErrorDescr: err.Error()}
	}
	defer rows.Close()

	var id int
	var used bool
	for rows.Next() {
		err := rows.Scan(&id, &used)
		if err != nil {
			return nil, Models.ResponseError{ErrorDescr: err.Error()}
		}
	}
	if id == 0 {
		return nil, Models.ResponseError{ErrorDescr: "No data on request parameters", Status: http.StatusNotFound}
	}

	if r.Used != used {
		if used {

			_, err = svc.DB.Exec("update configs set used=false where id = $1", id)
			if err != nil {
				return nil, Models.ResponseError{ErrorDescr: err.Error()}
			}

		} else {

			ctx := context.Background()
			tx, err := svc.DB.BeginTx(ctx, nil)
			if err != nil {
				return nil, Models.ResponseError{ErrorDescr: err.Error()}
			}

			_, err = tx.ExecContext(ctx, "update configs set used=false where service = $1 and used=true", r.Service)
			if err != nil {
				tx.Rollback()
				return nil, Models.ResponseError{ErrorDescr: err.Error()}
			}

			_, err = tx.ExecContext(ctx, "update configs set used=true where id = $1", id)
			if err != nil {
				tx.Rollback()
				return nil, Models.ResponseError{ErrorDescr: err.Error()}
			}

			err = tx.Commit()
			if err != nil {
				tx.Rollback()
				return nil, Models.ResponseError{ErrorDescr: err.Error()}
			}
		}
	}
	return r, nil
}

func (svc configService) DelConfig(_ context.Context, req interface{}) (*Models.ConfigRequest, error) {
	r := req.(*Models.ConfigRequest)
	if r.Version == 0 {
		return nil, Models.ResponseError{ErrorDescr: "version parameter must be specified", Status: http.StatusBadRequest}
	}

	rows, err := svc.DB.Query("select coalesce(id, 0), used from configs where service = $1 and version = $2 limit 1", r.Service, r.Version)
	if err != nil {
		return nil, Models.ResponseError{ErrorDescr: err.Error()}
	}
	defer rows.Close()

	var id int
	var used bool
	for rows.Next() {
		err := rows.Scan(&id, &used)
		if err != nil {
			return nil, Models.ResponseError{ErrorDescr: err.Error()}
		}
	}
	if id == 0 {
		return nil, Models.ResponseError{ErrorDescr: "No data on request parameters", Status: http.StatusNotFound}
	}
	if used {
		return nil, Models.ResponseError{ErrorDescr: "Specified config is used", Status: http.StatusForbidden}
	}

	_, err = svc.DB.Exec("delete from configs where id = $1", id)
	if err != nil {
		return nil, Models.ResponseError{ErrorDescr: err.Error()}
	}
	return r, nil
}

type configService struct {
	DB *sql.DB
}

func NewConfigService(postgresUri string) (*configService, error) {
	time.Sleep(2 * time.Second)
	db, err := sql.Open("postgres", postgresUri)
	if err != nil {
		return nil, fmt.Errorf("Can't connect to postgresql: %v", err)
	}
	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("Can't test ping to postgresql: %v", err)
	}
	return &configService{DB: db}, nil
}

func Shutdown(s *configService) {
	_ = s.DB.Close()
}
