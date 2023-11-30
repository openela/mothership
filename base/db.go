package base

import (
	"context"
	"errors"
	"fmt"
	orderedmap "github.com/wk8/go-ordered-map/v2"
	"go.ciq.dev/pika"
	"math/rand"
)

type Pika[T any] interface {
	pika.QuerySet[T]

	U(x any) error
	F(keyval ...any) Pika[T]
	D(x any) error
	Transaction(ctx context.Context) (Pika[T], error)
	Commit() error
}

type DB struct {
	*pika.PostgreSQL
}

type innerDB[T any] struct {
	pika.QuerySet[T]
	*DB
}

type nameInterface interface {
	GetID() string
}

func NewDB(databaseURL string) (*DB, error) {
	db, err := pika.NewPostgreSQL(databaseURL)
	if err != nil {
		return nil, err
	}

	return &DB{db}, nil
}

func NewDBArgs(keyval ...any) *orderedmap.OrderedMap[string, any] {
	args := pika.NewArgs()
	for i := 0; i < len(keyval); i += 2 {
		args.Set(keyval[i].(string), keyval[i+1])
	}

	return args
}

func NameGen(prefix string) string {
	// last part is a random int64
	// generate a length=18 random int
	minRan := 100000000000000000
	maxRan := 999999999999999999
	ran := minRan + rand.Intn(maxRan-minRan)

	return fmt.Sprintf("%s/%d", prefix, ran)
}

//goland:noinspection GoExportedFuncWithUnexportedType
func Q[T any](db *DB) Pika[T] {
	return &innerDB[T]{pika.Q[T](db.PostgreSQL), db}
}

func (inner *innerDB[T]) F(keyval ...any) Pika[T] {
	var qs pika.QuerySet[T] = inner
	args := pika.NewArgs()
	for i := 0; i < len(keyval); i += 2 {
		args.Set(keyval[i].(string), keyval[i+1])
		qs = qs.Filter(fmt.Sprintf("%s=:%s", keyval[i].(string), keyval[i].(string)))
	}

	inner.QuerySet = qs.Args(args)
	return inner
}

func (inner *innerDB[T]) D(x any) error {
	// Check if x has GetID() method
	var id any

	name, ok := x.(nameInterface)
	if ok {
		stringID := name.GetID()
		if stringID == "" {
			return fmt.Errorf("id is empty")
		}
		id = stringID
	}

	if id == nil {
		return errors.New("id is nil")
	}

	qs := inner.F("name", id)
	return qs.Delete()
}

func (inner *innerDB[T]) Transaction(ctx context.Context) (Pika[T], error) {
	innerDeref := *inner
	ts := &innerDeref
	err := ts.Begin(ctx)
	if err != nil {
		return nil, err
	}

	return &innerDB[T]{pika.Q[T](ts), inner.DB}, nil
}

func (inner *innerDB[T]) Commit() error {
	return inner.DB.Commit()
}

func (inner *innerDB[T]) U(x any) error {
	y := x.(*T)

	// Check if x has GetID() method
	var id any

	name, ok := x.(nameInterface)
	if ok {
		stringID := name.GetID()
		if stringID == "" {
			return fmt.Errorf("id is empty")
		}
		id = stringID
	}

	if id == nil {
		return fmt.Errorf("id is nil")
	}

	qs := inner.F("name", id)
	return qs.Update(y)
}

func (inner *innerDB[T]) Create(x *T) error {
	ctx := context.TODO()
	ts := pika.NewPostgreSQLFromDB(inner.DB.DB())
	err := ts.Begin(ctx)
	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	var y any = x
	err = pika.Q[T](ts).Create(x)
	if err != nil {
		return err
	}

	var nameID string
	name, ok := y.(nameInterface)
	if ok {
		stringID := name.GetID()
		if stringID == "" {
			return fmt.Errorf("id is empty")
		}
		nameID = stringID
	}

	if nameID == "" {
		return errors.New("id is empty")
	}

	args := pika.NewArgs()
	args.Set("name", nameID)
	newX, err := pika.Q[T](ts).Filter("name=:name").Args(args).Get()
	if err != nil {
		return err
	}

	err = ts.Commit()
	if err != nil {
		return err
	}

	*x = *newX

	return nil
}
