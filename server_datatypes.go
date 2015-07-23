package main

import "encoding/json"

func (n *node) newUser() *user {
	return &user{
		dtUser: &dtUser{},
		bolterItem: &bolterItem{
			DS: n.su.ds.dcbAsts,
		},
	}
}

type user struct {
	*bolterItem
	*dtUser
}

func (u *user) get() error {
	v, err := u.DS.getBytes(u.getID())
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
	if err = u.DS.setBytes(u.getID(), v); err != nil {
		return err
	}
	return nil
}
