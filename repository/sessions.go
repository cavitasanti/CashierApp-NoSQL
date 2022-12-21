package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"errors"
	"time"
)

type SessionsRepository struct {
	db db.DB
}

func NewSessionsRepository(db db.DB) SessionsRepository {
	return SessionsRepository{db}
}

func (u *SessionsRepository) ReadSessions() ([]model.Session, error) {
	records, err := u.db.Load("sessions")
	if err != nil {
		return nil, err
	}

	var listSessions []model.Session
	err = json.Unmarshal([]byte(records), &listSessions)
	if err != nil {
		return nil, err
	}

	return listSessions, nil
}

func (u *SessionsRepository) DeleteSessions(tokenTarget string) error {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return err
	}
	// Select target token and delete from listSessions
	m := []model.Session{}
	for i := 0; i < len(listSessions); i++ {
		if listSessions[i].Token == tokenTarget {
			listSessions = append(listSessions[:i], listSessions[i+1:]...)
		}
	}
	// convert listUser ke json
	jsonData, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = u.db.Save("sessions", []byte(jsonData))
	if err != nil {
		return err
	}
	return nil
}

func (u *SessionsRepository) AddSessions(session model.Session) error {
	// return nil // TODO: replace this
	listSessions, err := u.ReadSessions()
	if err != nil {
		return err
	}
	// masukan creds.Username ke data user.json
	listSessions = append(listSessions, session)
	// convert listUser ke json
	jsonData, err := json.Marshal(listSessions)
	if err != nil {
		return err
	} else {
		err = u.db.Save("sessions", []byte(jsonData))
		if err != nil {
			return err
		}
	}
	return nil
	// return errors.New("Token Not found!")
}

func (u *SessionsRepository) CheckExpireToken(token string) (model.Session, error) {
	// return model.Session{}, nil // TODO: replace this
	s, err := u.TokenExist(token)
	if err != nil {
		return model.Session{}, err
	}
	if u.TokenExpired(s) {
		return model.Session{}, errors.New("Token is Expired!")
	}
	return s, nil
}

func (u *SessionsRepository) ResetSessions() error {
	err := u.db.Reset("sessions", []byte("[]"))
	if err != nil {
		return err
	}

	return nil
}

func (u *SessionsRepository) TokenExist(req string) (model.Session, error) {
	listSessions, err := u.ReadSessions()
	if err != nil {
		return model.Session{}, err
	}
	for _, element := range listSessions {
		if element.Token == req {
			return element, nil
		}
	}
	return model.Session{}, errors.New("Token Not Found!")
	// return model.Session{}, fmt.Errorf("Token Not Found!")
}

func (u *SessionsRepository) TokenExpired(s model.Session) bool {
	return s.Expiry.Before(time.Now())
}
