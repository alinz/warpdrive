package data

import (
	"encoding/json"
	"fmt"
	"strings"

	db "upper.io/db.v2"
)

//PermissionType is a type to show permission
type PermissionType int

const (
	//AGENT can delete the project and make other people AGENT
	AGENT PermissionType = iota
	//ADMIN can invite other users and publish a build
	ADMIN
)

var (
	permissionTypeNameToValue = map[string]PermissionType{
		"AGENT": AGENT,
		"ADMIN": ADMIN,
	}

	permissionTypeValueToName = map[PermissionType]string{
		AGENT: "AGENT",
		ADMIN: "ADMIN",
	}
)

//MarshalJSON for type PermissionType
func (p PermissionType) MarshalJSON() ([]byte, error) {
	if s, ok := interface{}(p).(fmt.Stringer); ok {
		return json.Marshal(strings.ToLower(s.String()))
	}
	s, ok := permissionTypeValueToName[p]
	if !ok {
		return nil, fmt.Errorf("invalid Permission Type: %d", p)
	}

	return json.Marshal(strings.ToLower(s))
}

//UnmarshalJSON for type PermissionType
func (p *PermissionType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("Permission Type should be a string, got %s", data)
	}
	v, ok := permissionTypeNameToValue[strings.ToUpper(s)]
	if !ok {
		return fmt.Errorf("invalid Permission Type %q", s)
	}
	*p = v
	return nil
}

//Permission this is reppresentation of Permissions tbale
type Permission struct {
	ID         int64          `db:"id,omitempty,pk" json:"-"`
	UserID     int64          `db:"user_id" json:"-"`
	AppID      int64          `db:"app_id" json:"-"`
	Permission PermissionType `db:"permission" json:"-"`
}

func (p Permission) CollectionName() string {
	return "permissions"
}

func (p Permission) Query(session db.Database, query db.Cond) db.Result {
	return session.Collection(p.CollectionName()).Find(query)
}

func (p *Permission) Find(session db.Database, query db.Cond) error {
	return p.Query(session, query).One(p)
}

func (p *Permission) Save(session db.Database) error {
	collection := session.Collection(p.CollectionName())
	var err error

	if p.ID == 0 {
		var id interface{}
		id, err = collection.Insert(p)
		if err == nil {
			p.ID = id.(int64)
		}
	} else {
		err = collection.
			Find(db.Cond{"id": p.ID}).
			Update(p)
	}

	return err
}

func (p *Permission) Remove(session db.Database) error {
	return p.Query(session, db.Cond{"id": p.ID}).Delete()
}
