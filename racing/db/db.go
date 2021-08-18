package db

import (
	"database/sql"
	"time"

	"github.com/ashleyjlive/entain/racing/proto/racing"
	"github.com/golang/protobuf/ptypes"
	"syreclabs.com/go/faker"
)

func (r *racesRepo) seed() error {
	statement, err := r.db.Prepare(`CREATE TABLE IF NOT EXISTS races (id INTEGER PRIMARY KEY, meeting_id INTEGER, name TEXT, number INTEGER, visible INTEGER, advertised_start_time DATETIME)`)
	if err == nil {
		_, err = statement.Exec()
	}

	for i := 1; i <= 100; i++ {
		statement, err = r.db.Prepare(`INSERT OR IGNORE INTO races(id, meeting_id, name, number, visible, advertised_start_time) VALUES (?,?,?,?,?,?)`)
		if err == nil {
			_, err = statement.Exec(
				i,
				faker.Number().Between(1, 10),
				faker.Team().Name(),
				faker.Number().Between(1, 12),
				faker.Number().Between(0, 1),
				faker.Time().Between(time.Now().AddDate(0, 0, -1), time.Now().AddDate(0, 0, 2)).Format(time.RFC3339),
			)
		}
	}

	return err
}

func (r *racesRepo) insert(race *racing.Race) error {
	var statement *sql.Stmt
	ts, err := ptypes.Timestamp(race.AdvertisedStartTime)
	if err != nil {
		return err
	}
	statement, err = r.db.Prepare(`INSERT INTO races(id, meeting_id, name, number, visible, advertised_start_time) VALUES (?,?,?,?,?,?)`)
	if err == nil {
		_, err = statement.Exec(
			&race.Id,
			&race.MeetingId,
			&race.Name,
			&race.Number,
			&race.Visible,
			ts,
		)
	}
	return err
}

func (r *racesRepo) clear() error {
	_, err := r.db.Exec("DELETE FROM races")
	return err
}
