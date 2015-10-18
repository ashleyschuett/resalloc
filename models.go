package main

import (
	"database/sql"
	"errors"
	"log"

	"github.com/doug-martin/goqu"
	"github.com/jmcvetta/randutil"
	_ "github.com/mattn/go-sqlite3" // sqlite3 adapter
)

var db = setDB()

func setDB() *goqu.Database {
	sqliteDb, err := sql.Open("sqlite3", "./resalloc.db")
	if err != nil {
		log.Print(err.Error())
	}
	db := goqu.New("sqlite3", sqliteDb)
	return db
}

// User represents a row from
// the users table
type User struct {
	Username    string `db:"username"`
	Password    string `db:"password"`
	AccountType string `db:"account_type"`
	Activated   int    `db:"activated"`
	Token       string `db:"token"`
}

// Create inserts one row
// into the database
func (u User) Create() error {
	// Generate hashed password and store to db
	u.Password = generatePassword(u.Password)
	if _, err := db.From("users").
		Insert(goqu.Record{
		"username":     u.Username,
		"password":     u.Password,
		"token":        "",
		"account_type": "USR",
		"activated":    1,
	}).Exec(); err != nil {
		return err
	}
	return nil
}

// Polulate takes a username and populates
// a User struct with information from the db
func (u *User) Polulate(username string) error {
	if _, err := db.From("users").
		Where(goqu.Ex{
		"username": username,
	}).ScanStruct(u); err != nil {
		return err
	}
	return nil
}

// PopulateFromToken populates the user struct
// when only a token is available
func (u *User) PopulateFromToken(token string) error {
	if _, err := db.From("users").
		Where(goqu.Ex{
		"token": token,
	}).ScanStruct(u); err != nil {
		return err
	}
	return nil
}

// GenerateToken creates a sudo random string
// and stores it in the database
func (u *User) GenerateToken() error {
	u.Token, _ = randutil.String(50, "7148fa22885551edb6405895ddf7cd0f")
	if _, err := db.From("users").
		Where(goqu.Ex{
		"username": u.Username,
	}).Update(goqu.Record{
		"token": u.Token,
	}).Exec(); err != nil {
		return err
	}
	return nil
}

// VerifyToken ensures that the passed in
// token is a valid token in the database
func VerifyToken(token string) error {
	var count int
	if _, err := db.From("users").
		Select(goqu.COUNT("*")).
		Where(goqu.Ex{
		"token": token,
	}).ScanVal(&count); err != nil {
		return err
	}

	if count != 1 {
		return errors.New("No valid token found")
	}

	return nil
}

// Resource represents the structure
// of the resources table in the database
type Resource struct {
	Name string `db:"name"`
	File string `db:"file"`
}

// Create inserts a new row into the
// resources table
func (r Resource) Create() error {
	// Generate hashed password and store to db
	if _, err := db.From("resources").
		Insert(goqu.Record{
		"name": r.Name,
		"file": r.File,
	}).Exec(); err != nil {
		return err
	}
	return nil
}

// Fetch grabs one machine for creation
func (r *Resource) Fetch(name string) error {
	if _, err := db.From("resources").
		Where(goqu.Ex{
		"name": name,
	}).ScanStruct(r); err != nil {
		return err
	}
	return nil
}

// FetchAll gets a list of all available resources
// that are in the database this is kind of odd
// since it is actually not requires that you have
// already instantiated a Resource struct to make this
// call... however I find from a readablity standpoint
// it makes more sense if you call this on an empty
// Resource struct. Renaming it FetchResouces would
// be another option.
func (r Resource) FetchAll() ([]Resource, error) {
	var resources []Resource
	if err := db.From("resources").
		ScanStructs(&resources); err != nil {
		return nil, err
	}
	return resources, nil
}

// Machine represents the structure
// of the resources table in the database
type Machine struct {
	Name     string `db:"name"`
	Username string `db:"username"`
	IP       string `db:"ip"`
}

// Create inserts a new row into the
// resources table
func (m Machine) Create() error {
	// Generate hashed password and store to db
	if _, err := db.From("machines").
		Insert(goqu.Record{
		"name":     m.Name,
		"username": m.Username,
		"ip":       m.IP,
	}).Exec(); err != nil {
		return err
	}
	return nil
}

// FetchRand will grab a random machine
// to spin up the requests lease on
// In reality images should likely have tags
// so that you can better optimize how the
// hardware would be shared among a large
// amount of users
func (m *Machine) FetchRand() error {
	// Goqu doesn't support order by random
	// So i'm faking it in order to take
	// advantage of the ScanStructs function still
	var ms []Machine
	if err := db.From("machines").
		ScanStructs(&ms); err != nil {
		return err
	}
	// Get "random" row
	total := len(ms)
	randInt, _ := randutil.IntRange(0, 1000)
	row := randInt % total
	*m = ms[row]
	return nil
}

// Lease represents a lease that has been
// created on a remote machine
type Lease struct {
	Name        string `db:"name"`
	Username    string `db:"username"`
	MachineName string `db:"machine_name"`
}

// Create inserts a new row into the
// resources table
func (l Lease) Create() error {
	// Generate hashed password and store to db
	if _, err := db.From("leases").
		Insert(goqu.Record{
		"name":         l.Name,
		"username":     l.Username,
		"machine_name": l.MachineName,
	}).Exec(); err != nil {
		return err
	}
	return nil
}

// FetchAll returns a list of all active leases
func (l Lease) FetchAll() ([]Lease, error) {
	var leases []Lease
	if err := db.From("leases").
		ScanStructs(&leases); err != nil {
		return nil, err
	}
	return leases, nil
}

// Fetch returns a lease if it exists
func (l *Lease) Fetch(name string) error {
	if _, err := db.From("leases").
		Where(goqu.Ex{
		"name": name,
	}).ScanStruct(l); err != nil {
		return err
	}
	return nil
}

// Delete will remove the current
// value stored in l.Name from the
// database.
func (l *Lease) Delete() error {
	if _, err := db.From("leases").
		Where(goqu.Ex{
		"name": l.Name,
	}).Delete().Exec(); err != nil {
		return err
	}
	return nil
}
