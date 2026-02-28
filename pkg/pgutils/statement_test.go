package pgutils_test

import (
	"reflect"
	"testing"

	"github.com/puriice/httplibs/pkg/pgutils"
)

type user struct {
	Id       *string `db:"id"`
	Username *string `db:"username"`
	Password *string `db:"password"`
}

func TestFromOne(t *testing.T) {
	id := "123"
	username := "Test"
	password := "Test1234"

	user := &user{
		Id:       &id,
		Username: &username,
		Password: &password,
	}

	statement, argv, err := pgutils.CreateSetStatement(*user, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedStatement := "id = $1, username = $2, password = $3"
	if statement != expectedStatement {
		t.Fatalf("unexpected statement.\nexpected: %s\ngot:      %s",
			expectedStatement, statement)
	}

	expectedArgs := []any{"123", "Test", "Test1234"}

	if !reflect.DeepEqual(argv, expectedArgs) {
		t.Fatalf("unexpected args.\nexpected: %#v\ngot:      %#v",
			expectedArgs, argv)
	}
}

func TestFromFive(t *testing.T) {
	id := "123"
	username := "Test"
	password := "Test1234"

	user := &user{
		Id:       &id,
		Username: &username,
		Password: &password,
	}

	statement, argv, err := pgutils.CreateSetStatement(*user, 5)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedStatement := "id = $5, username = $6, password = $7"
	if statement != expectedStatement {
		t.Fatalf("unexpected statement.\nexpected: %s\ngot:      %s",
			expectedStatement, statement)
	}

	expectedArgs := []any{"123", "Test", "Test1234"}

	if !reflect.DeepEqual(argv, expectedArgs) {
		t.Fatalf("unexpected args.\nexpected: %#v\ngot:      %#v",
			expectedArgs, argv)
	}
}

func TestMissingSomeFields(t *testing.T) {
	id := "123"
	username := "Test"

	user := &user{
		Id:       &id,
		Username: &username,
	}

	statement, argv, err := pgutils.CreateSetStatement(*user, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedStatement := "id = $1, username = $2"
	if statement != expectedStatement {
		t.Fatalf("unexpected statement.\nexpected: %s\ngot:      %s",
			expectedStatement, statement)
	}

	expectedArgs := []any{"123", "Test"}

	if !reflect.DeepEqual(argv, expectedArgs) {
		t.Fatalf("unexpected args.\nexpected: %#v\ngot:      %#v",
			expectedArgs, argv)
	}
}

func TestFromPointer(t *testing.T) {
	id := "123"
	username := "Test"

	user := &user{
		Id:       &id,
		Username: &username,
	}

	statement, argv, err := pgutils.CreateSetStatement(user, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expectedStatement := "id = $1, username = $2"
	if statement != expectedStatement {
		t.Fatalf("unexpected statement.\nexpected: %s\ngot:      %s",
			expectedStatement, statement)
	}

	expectedArgs := []any{"123", "Test"}

	if !reflect.DeepEqual(argv, expectedArgs) {
		t.Fatalf("unexpected args.\nexpected: %#v\ngot:      %#v",
			expectedArgs, argv)
	}
}
