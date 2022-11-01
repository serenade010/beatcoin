package models

import (
	"database/sql"
	"errors"
	"fmt"
)

type Model struct {
	Id             int
	Name           string
	Belongs_to     int
	Ratio_of_train float32
	Look_back      int
	Forecast_days  int
	Crypto         string
	First_layer    int
	Second_layer   sql.NullInt16
	Third_layer    sql.NullInt16
	First_index    sql.NullString
	Second_index   sql.NullString
	Third_index    sql.NullString
	Learning_rate  float32
	Epoch          int
	Batch_size     int
	Modelerr       sql.NullFloat64
}

type ModelModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.

func (m *ModelModel) Insert(name string, belongs_to int, ratio_of_train float32, look_back int, forecast_days int, crypto string, first_layer int, second_layer int, third_layer int, first_index string, second_index string, third_index string, learning_rate float32, epoch int, batch_size int, modelerr float32) (int, error) {
	stmt := "INSERT INTO model (modelname,belongs_to,ratio_of_train,look_back,forecast_days,crypto,first_layer,second_layer,third_layer,first_index,second_index,third_index,learning_rate,epoch,batch_size,modelerr) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15,$16) RETURNING id"

	var lastInsertId int
	err := m.DB.QueryRow(stmt, name, belongs_to, ratio_of_train, look_back, forecast_days, crypto, first_layer, second_layer, third_layer, first_index, second_index, third_index, learning_rate, epoch, batch_size, modelerr).Scan(&lastInsertId)
	if err != nil {
		return -1, err
	}

	return lastInsertId, err
}

//This will return a specific snippet based on its id.

func (m *ModelModel) Get(id int) (*Model, error) {
	stmt := "SELECT id,modelname,belongs_to,ratio_of_train,look_back,forecast_days,crypto,first_layer,second_layer,third_layer,first_index,second_index,third_index,learning_rate,epoch,batch_size,modelerr FROM model WHERE id=$1"

	row := m.DB.QueryRow(stmt, id)

	model := &Model{}

	err := row.Scan(&model.Id, &model.Name, &model.Belongs_to, &model.Ratio_of_train, &model.Look_back, &model.Forecast_days, &model.Crypto, &model.First_layer, &model.Second_layer, &model.Third_layer, &model.First_index, &model.Second_index, &model.Third_index, &model.Learning_rate, &model.Epoch, &model.Batch_size, &model.Modelerr)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNoRecord
		} else {
			return nil, err
		}
	}

	return model, nil
}

func (m *ModelModel) Best() ([]*Model, error) {
	stmt := "SELECT id,modelname,belongs_to,ratio_of_train,look_back,forecast_days,crypto,first_layer,second_layer,third_layer,first_index,second_index,third_index,learning_rate,epoch,batch_size,modelerr FROM model ORDER BY modelerr ASC LIMIT 20"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	models := []*Model{}

	for rows.Next() {
		model := &Model{}

		err := rows.Scan(&model.Id, &model.Name, &model.Belongs_to, &model.Ratio_of_train, &model.Look_back, &model.Forecast_days, &model.Crypto, &model.First_layer, &model.Second_layer, &model.Third_layer, &model.First_index, &model.Second_index, &model.Third_index, &model.Learning_rate, &model.Epoch, &model.Batch_size, &model.Modelerr)

		if err != nil {
			return nil, err
		}

		models = append(models, model)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}

func (m *ModelModel) MyModels(id int) ([]*Model, error) {
	stmt := "SELECT id,modelname,belongs_to,ratio_of_train,look_back,forecast_days,crypto,first_layer,second_layer,third_layer,first_index,second_index,third_index,learning_rate,epoch,batch_size,modelerr FROM model WHERE belongs_to=$1"

	rows, err := m.DB.Query(stmt, id)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	models := []*Model{}

	for rows.Next() {
		model := &Model{}

		err := rows.Scan(&model.Id, &model.Name, &model.Belongs_to, &model.Ratio_of_train, &model.Look_back, &model.Forecast_days, &model.Crypto, &model.First_layer, &model.Second_layer, &model.Third_layer, &model.First_index, &model.Second_index, &model.Third_index, &model.Learning_rate, &model.Epoch, &model.Batch_size, &model.Modelerr)

		if err != nil {
			return nil, err
		}

		models = append(models, model)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return models, nil
}

func (m *ModelModel) Belong(modelId int, userId int) bool {
	stmt := "SELECT belongs_to FROM model WHERE id=$1"
	row := m.DB.QueryRow(stmt, modelId)

	var belongsID int
	err := row.Scan(&belongsID)
	if err != nil || belongsID != userId {
		return false
	} else {
		return true
	}
}

func (m *ModelModel) UpdateMAPE(modelId int, mape float64) error {
	stmt := "SELECT modelerr FROM model WHERE id=$1"
	row := m.DB.QueryRow(stmt, modelId)

	var modelerr float64
	err := row.Scan(&modelerr)
	if err != nil {
		return err
	}
	if mape <= modelerr {
		fmt.Println(mape, modelerr)
		stmt := "UPDATE model SET modelerr=$1 WHERE id=$2"
		_, err := m.DB.Exec(stmt, mape, modelId)
		if err != nil {
			return err
		}
	}
	return nil
}
