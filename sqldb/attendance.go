package sqldb

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/boseungjeong/wedding-invitation-server/types"
)

func initializeAttendanceTable() error {
	_, err := sqlDb.Exec(`
		CREATE TABLE IF NOT EXISTS attendance (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			side VARCHAR(10),
			name VARCHER(20),
			meal VARCHAR(20),
			count INTEGER,
			timestamp INTEGER
		)
	`)
	return err
}

func CreateAttendance(side, name, meal string, count int) error {
	_, err := sqlDb.Exec(`
		INSERT INTO attendance (side, name, meal, count, timestamp)
		VALUES (?, ?, ?, ?, ?)
	`, side, name, meal, count, time.Now().Unix())
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

func GetAttendance() (*types.AttendanceGetResponse, error) {
	rows, err := sqlDb.Query(`
		SELECT id, side, name, meal, count, timestamp
		FROM attendance
		ORDER BY timestamp DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	response := &types.AttendanceGetResponse{
		Attendances: []types.AttendanceForGet{},
	}

	for rows.Next() {
		attendance := types.AttendanceForGet{}
		err := rows.Scan(&attendance.Id, &attendance.Side, &attendance.Name, &attendance.Meal, &attendance.Count, &attendance.Timestamp)
		if err != nil {
			return nil, err
		}
		response.Attendances = append(response.Attendances, attendance)
	}

	var totalCount sql.NullInt64
	err = sqlDb.QueryRow(`
		SELECT COALESCE(SUM(count), 0)
		FROM attendance
	`).Scan(&totalCount)
	if err != nil {
		return nil, err
	}

	response.Total = len(response.Attendances)
	response.TotalCount = int(totalCount.Int64)

	return response, nil
}
