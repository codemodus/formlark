package main

import (
	"encoding/json"
	"errors"
)

func (n *node) newUser() *user {
	return &user{
		dtUser: &dtUser{},
		boltItem: &boltItem{
			DS: n.su.ds,
		},
	}
}

type user struct {
	*boltItem
	*dtUser
}

func (u *user) getID() (string, error) {
	if u.ID != "" {
		return u.ID, nil
	}
	if u.PublicID != "" {
		k, err := u.DS.dcbMrks.getBytes(u.PublicID)
		if err != nil {
			return "", err
		}
		if len(k) > 0 {
			return string(k), nil
		}
	}
	if u.Email != "" {
		k, err := u.DS.dcbMrks.getBytes(u.Email)
		if err != nil {
			return "", err
		}
		if len(k) > 0 {
			return string(k), nil
		}
	}
	return "", errors.New("no id")
}

func (u *user) get() error {
	id, err := u.getID()
	if err != nil {
		return err
	}
	v, err := u.DS.dcbAsts.getBytes(id)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(v, u); err != nil {
		return err
	}
	return nil
}

func (u *user) set() error {
	v, err := json.Marshal(u)
	if err != nil {
		return err
	}
	id, err := u.getID()
	if err != nil {
		return err
	}
	var pidErr, eErr error
	if u.PublicID != "" {
		pidErr = u.DS.dcbMrks.setBytes(u.PublicID, []byte(id))
	}
	if u.Email != "" {
		eErr = u.DS.dcbMrks.setBytes(u.Email, []byte(id))
	}
	if pidErr != nil && eErr != nil {
		return errors.New("cannot save markers")
	}
	if err = u.DS.dcbAsts.setBytes(id, v); err != nil {
		return err
	}
	return nil
}
