package bug_test

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"testing"

	"entgo.io/bug"
	"entgo.io/ent/dialect"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"

	"entgo.io/bug/ent"
	"entgo.io/bug/ent/enttest"
	"entgo.io/bug/ent/user"
)

func TestBugSQLite(t *testing.T) {
	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
	defer client.Close()
	test(t, client)
}

func TestBugMySQL(t *testing.T) {
	for version, port := range map[string]int{"56": 3306, "57": 3307, "8": 3308} {
		addr := net.JoinHostPort("localhost", strconv.Itoa(port))
		t.Run(version, func(t *testing.T) {
			client := enttest.Open(t, dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
			defer client.Close()
			test(t, client)
		})
	}
}

func TestBugPostgres(t *testing.T) {
	for version, port := range map[string]int{"10": 5430, "11": 5431, "12": 5432, "13": 5433, "14": 5434} {
		t.Run(version, func(t *testing.T) {
			client := enttest.Open(t, dialect.Postgres, fmt.Sprintf("host=localhost port=%d user=postgres dbname=test password=pass sslmode=disable", port))
			defer client.Close()
			test(t, client)
		})
	}
}

func TestBugMaria(t *testing.T) {
	for version, port := range map[string]int{"10.5": 4306, "10.2": 4307, "10.3": 4308} {
		t.Run(version, func(t *testing.T) {
			addr := net.JoinHostPort("localhost", strconv.Itoa(port))
			client := enttest.Open(t, dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
			defer client.Close()
			test(t, client)
		})
	}
}

func test(t *testing.T, client *ent.Client) {
	ctx := context.Background()
	client.User.Delete().ExecX(ctx)

	// create custom type with nil
	client.User.Create().SetName("Ariel").SetAge(30).ExecX(ctx)

	// create custom type with empty
	client.User.Create().SetName("Bob").SetAge(30).SetFile(bug.File{File: ""}).ExecX(ctx)

	// create custom type with value
	client.User.Create().SetName("Bob").SetAge(30).SetFile(bug.File{File: "aaaa"}).ExecX(ctx)

	users := client.User.Query().Order(ent.Asc(user.FieldFile)).AllX(ctx)

	require.Equal(t, users[0].File.File, "")
	require.Equal(t, users[1].File.File, "")
	require.Equal(t, users[2].File.File, "aaaaa")

}
