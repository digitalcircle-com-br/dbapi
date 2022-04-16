package client

import (
	"fmt"

	"github.com/digitalcircle-com-br/dbapi/lib/types"
	"github.com/digitalcircle-com-br/httpcli"
)

type Client struct {
	cli httpcli.Client
}

func New() *Client {
	ret := Client{}
	ret.cli = *httpcli.New()
	return &ret
}
func (c *Client) SetBasePath(s string) {
	c.cli.BasePath = s
}

func (c Client) Tenants() (out []string, err error) {
	out = make([]string, 0)
	err = c.cli.JsonGet("/tenants", &out)
	return
}

func (c Client) Tenant(t string) (out string, err error) {
	out = ""
	err = c.cli.JsonGet(fmt.Sprintf("/tenant/%s", t), &out)
	return
}

func (c Client) SetTenant(t string, d string) (out int64, err error) {
	err = c.cli.JsonPost(fmt.Sprintf("/tenant/%s", t), d, &out)
	return
}

func (c Client) DelTenant(t string) (out int64, err error) {
	err = c.cli.JsonDelete(fmt.Sprintf("/tenant/%s", t), &out)
	return
}

func (c Client) Admin(t *types.DBIn) (out *types.DBOut, err error) {
	out = &types.DBOut{}
	err = c.cli.JsonPost("/admin", t, out)
	return
}

func (c Client) DBs() (out *types.DBOut, err error) {
	out = &types.DBOut{}
	err = c.cli.JsonGet("/dbs", &out)
	return
}

func (c Client) Tables(n string) (out *types.DBOut, err error) {
	out = &types.DBOut{}
	err = c.cli.JsonGet(fmt.Sprintf("/db/%s/tables", n), out)
	return
}

func (c Client) Init(n string) (out string, err error) {
	lout := &types.DBOut{}
	err = c.cli.JsonGet(fmt.Sprintf("/db/%s/init", n), lout)
	out = lout.Data.(string)
	return
}

func (c Client) CreateDB(n string, dsn string) (out string, err error) {
	_, err = c.Admin(&types.DBIn{Q: "create database " + n})
	if err != nil {
		return
	}
	_, err = c.SetTenant(n, dsn)
	if err != nil {
		return
	}

	out, err = c.Init(n)
	return
}
