package utils

import (
	"database/sql"
	"encoding/json"
	"reflect"
	"time"
)

type (
	NullString sql.NullString
	NullInt32  sql.NullInt32
	NullTime   sql.NullTime
)

func NewNullString(s string) NullString {
	if s == "" {
		return NullString{String: "", Valid: false}
	}
	return NullString{String: s, Valid: true}
}

func (ns *NullString) MarshalJSON() ([]byte, error) {
	if !ns.Valid {
		ns.Valid = true
		ns.String = ""
		// return []byte("null"), nil
	}
	return json.Marshal(ns.String)
}

// UnmarshalJSON implements the json.Unmarshaler interface for NullString
func (ns *NullString) UnmarshalJSON(data []byte) error {
	// If the value is "null", set the NullString to be invalid
	if string(data) == "null" {
		ns.Valid = false
		ns.String = ""
		return nil
	}
	// Unmarshal the data into a string
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return err
	}
	ns.String = str
	ns.Valid = true
	return nil
}

// Scan implements the Scanner interface for NullString
func (ns *NullString) Scan(value interface{}) error {
	var s sql.NullString
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullString{s.String, false}
	} else {
		*ns = NullString{s.String, true}
	}

	return nil
}

func NewNullInt32(i int32) NullInt32 {
	return NullInt32{Int32: i, Valid: true}
}

func (ni *NullInt32) MarshalJSON() ([]byte, error) {
	if !ni.Valid {
		ni.Valid = true
		ni.Int32 = 0
		// return []byte("null"), nil
	}
	return json.Marshal(ni.Int32)
}

// UnmarshalJSON implements the json.Unmarshaler interface for NullInt32
func (ni *NullInt32) UnmarshalJSON(data []byte) error {
	// If the value is "null", set the NullInt32 to be invalid
	if string(data) == "null" {
		ni.Valid = false
		ni.Int32 = 0
		return nil
	}
	// Unmarshal the data into an int32
	var i int32
	if err := json.Unmarshal(data, &i); err != nil {
		return err
	}
	ni.Int32 = i
	ni.Valid = true
	return nil
}

// Scan implements the Scanner interface for NullInt32
func (ns *NullInt32) Scan(value interface{}) error {
	var s sql.NullInt32
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullInt32{s.Int32, false}
	} else {
		*ns = NullInt32{s.Int32, true}
	}

	return nil
}

func NewNullTime(t time.Time) NullTime {
	return NullTime{Time: t, Valid: true}
}

func (nt *NullTime) MarshalJSON() ([]byte, error) {
	if !nt.Valid {
		nt.Valid = true
		nt.Time = time.Time{}
		// return []byte("null"), nil
	}
	return json.Marshal(nt.Time)
}

// UnmarshalJSON implements the json.Unmarshaler interface for NullTime
func (nt *NullTime) UnmarshalJSON(data []byte) error {
	// If the value is "null", set the NullTime to be invalid
	if string(data) == "null" {
		nt.Valid = false
		nt.Time = time.Time{}
		return nil
	}
	// Unmarshal the data into a time.Time
	var t time.Time
	if err := json.Unmarshal(data, &t); err != nil {
		return err
	}
	nt.Time = t
	nt.Valid = true
	return nil
}

// Scan implements the Scanner interface for NullInt32
func (ns *NullTime) Scan(value interface{}) error {
	var s sql.NullTime
	if err := s.Scan(value); err != nil {
		return err
	}

	// if nil then make Valid false
	if reflect.TypeOf(value) == nil {
		*ns = NullTime{s.Time, false}
	} else {
		*ns = NullTime{s.Time, true}
	}

	return nil
}
