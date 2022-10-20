package models

import (
	"database/sql"
	"errors"
)

type Model struct {
	id             int
	belongs_to     int
	ratio_of_train float32
	look_back      int
	forecast_days  int
	crypto         string
	first_layer    int
	second_layer   sql.NullInt16
	third_layer    sql.NullInt16
	first_index    sql.NullString
	second_index   sql.NullString
	third_index    sql.NullString
	learning_rate  float32
	epoch          int
	batch_size     int
	modelerr       sql.NullFloat64
}

type ModelModel struct {
	DB *sql.DB
}

// This will insert a new snippet into the database.

func (m *ModelModel) Insert(belongs_to int, ratio_of_train float32, look_back int, forecast_days int, crypto string, first_layer int, second_layer int, third_layer int, first_index string, second_index string, third_index string, learning_rate float32, epoch int, batch_size int, modelerr float32) error {
	stmt := "INSERT INTO model (belongs_to,ratio_of_train,look_back,forecast_days,crypto,first_layer,second_layer,third_layer,first_index,second_index,third_index,learning_rate,epoch,batch_size,modelerr) values ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15);"

	_, err := m.DB.Exec(stmt, belongs_to, ratio_of_train, look_back, forecast_days, crypto, first_layer, second_layer, third_layer, first_index, second_index, third_index, learning_rate, epoch, batch_size, modelerr)
	if err != nil {
		return err
	}

	return err
}

//This will return a specific snippet based on its id.

func (m *ModelModel) Get(id int) (*Model, error) {
	stmt := "SELECT id,belongs_to,ratio_of_train,look_back,forecast_days,crypto,first_layer,second_layer,third_layer,first_index,second_index,third_index,learning_rate,epoch,batch_size,modelerr FROM model WHERE id=$1"

	row := m.DB.QueryRow(stmt, id)

	model := &Model{}

	err := row.Scan(&model.id, &model.belongs_to, &model.ratio_of_train, &model.look_back, &model.forecast_days, &model.crypto, &model.first_layer, &model.second_layer, &model.third_layer, &model.first_index, &model.second_index, &model.third_index, &model.learning_rate, &model.epoch, &model.batch_size, &model.modelerr)

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
	stmt := "SELECT id,belongs_to,ratio_of_train,look_back,forecast_days,crypto,first_layer,second_layer,third_layer,first_index,second_index,third_index,learning_rate,epoch,batch_size,modelerr FROM model ORDER BY modelerr ASC LIMIT 10"

	rows, err := m.DB.Query(stmt)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	models := []*Model{}

	for rows.Next() {
		model := &Model{}

		err := rows.Scan(&model.id, &model.belongs_to, &model.ratio_of_train, &model.look_back, &model.forecast_days, &model.crypto, &model.first_layer, &model.second_layer, &model.third_layer, &model.first_index, &model.second_index, &model.third_index, &model.learning_rate, &model.epoch, &model.batch_size, &model.modelerr)
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
