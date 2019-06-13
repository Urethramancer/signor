package main

import "time"

type Nothing struct {
	ID     int64     `db:"smallint,unique,index,primary"`
	Name   string    `db:"string,unique"`
	When   time.Time `db:"timestamp"`
	Active bool      `db:"boolean"`
	Pi     float64   `db:"real"`
	Data   []byte    `db:"bytea"`
}
