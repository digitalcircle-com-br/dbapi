package client_test

import (
	"log"
	"testing"

	"github.com/digitalcircle-com-br/dbapi/lib/client"
	"github.com/digitalcircle-com-br/dbapi/lib/types"
	"github.com/stretchr/testify/assert"
)

func TestAPI(t *testing.T) {

	cli := client.New()
	cli.SetBasePath("http://localhost:8080")

	t.Run("tenant", func(t *testing.T) {
		t.Run("Get", func(t *testing.T) {
			o, err := cli.Tenants()
			assert.NoError(t, err)
			log.Printf("Got: %v", o)
		})
		t.Run("Add", func(t *testing.T) {
			o, err := cli.SetTenant("some", "asd123")
			assert.NoError(t, err)
			log.Printf("Got: %v", o)
		})
		t.Run("Get", func(t *testing.T) {
			o, err := cli.Tenant("some")
			assert.NoError(t, err)
			log.Printf("Got: %v", o)
		})
		t.Run("Get", func(t *testing.T) {

			o, err := cli.Tenants()
			assert.NoError(t, err)
			log.Printf("Got: %v", o)
		})
		t.Run("Del", func(t *testing.T) {
			o, err := cli.DelTenant("some")
			assert.NoError(t, err)
			log.Printf("Got: %v", o)
		})
		t.Run("Get", func(t *testing.T) {

			o, err := cli.Tenants()
			assert.NoError(t, err)
			log.Printf("Got: %v", o)
		})

		t.Run("Set Default DSN", func(t *testing.T) {
			o, err := cli.SetTenant("default", "host=localhost user=postgres password=Aa1234 dbname=postgres")
			assert.NoError(t, err)
			log.Printf("Got: %v", o)
		})
	})
}

func TestSetDefaultDSN(t *testing.T) {
	cli := client.New()
	cli.SetBasePath("http://localhost:8080")
	o, err := cli.SetTenant("default", "host=localhost user=postgres password=Aa1234 dbname=postgres")
	assert.NoError(t, err)
	log.Printf("Got: %v", o)
}
func TestAdmin(t *testing.T) {
	cli := client.New()
	cli.SetBasePath("http://localhost:8080")
	out, err := cli.Admin(&types.DBIn{Q: "select 1+1 , 2+2"})
	assert.NoError(t, err)
	log.Printf("got: %#v", out)
}

func TestCreateDBAndDrop(t *testing.T) {
	cli := client.New()
	cli.SetBasePath("http://localhost:8080")

	out, err := cli.Admin(&types.DBIn{Q: "create database db01"})
	assert.NoError(t, err)
	log.Printf("got: %#v", out)

	o, err := cli.SetTenant("db01", "host=localhost user=postgres password=Aa1234 dbname=db01")
	assert.NoError(t, err)
	log.Printf("got: %#v", o)

	out, err = cli.Admin(&types.DBIn{T: "db01", Q: "create table tb01(a text,b text)"})
	assert.NoError(t, err)
	log.Printf("got: %#v", out)

	out, err = cli.Admin(&types.DBIn{T: "db01", Q: "Select * from tb01"})
	assert.NoError(t, err)
	log.Printf("got: %#v", out)

	out, err = cli.Admin(&types.DBIn{Q: "drop database db01"})
	assert.NoError(t, err)
	log.Printf("got: %#v", out)

	o, err = cli.DelTenant("db01")
	assert.NoError(t, err)
	log.Printf("got: %#v", o)
}

func TestListDBs(t *testing.T) {
	cli := client.New()
	cli.SetBasePath("http://localhost:8080")
	o, err := cli.DBs()
	assert.NoError(t, err)
	log.Printf("%#v", o)

}

func TestListTables(t *testing.T) {
	cli := client.New()
	cli.SetBasePath("http://localhost:8080")
	o, err := cli.Tables("default")
	assert.NoError(t, err)
	log.Printf("%#v", o)

}

func TestCreateDB(t *testing.T) {
	cli := client.New()
	cli.SetBasePath("http://localhost:8080")
	o, err := cli.CreateDB("auth","host=localhost user=postgres password=Aa1234 dbname=auth")
	assert.NoError(t, err)
	log.Printf("%#v", o)
}
