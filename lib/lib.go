package lib

import (
	"context"
	"net/http"

	"github.com/digitalcircle-com-br/dbapi/lib/types"
	"github.com/digitalcircle-com-br/service"
	"golang.org/x/crypto/bcrypt"
)

func Run() error {

	service.Init("dbapi")

	service.HttpHandle("/tenants", http.MethodGet, "dbapi.admin", func(ctx context.Context, in *service.EMPTY_TYPE) (out []string, err error) {
		ts, err := service.DataHGetAll("dsn")
		if err != nil {
			return
		}
		out = make([]string, len(ts))
		for k := range ts {
			out = append(out, k)
		}
		return
	})

	service.HttpHandle("/tenant/{id}", http.MethodGet, "dbapi.admin", func(ctx context.Context, in *service.EMPTY_TYPE) (out string, err error) {
		dsn := service.CtxVars(ctx)["id"]
		out, err = service.DataHGet("dsn", dsn)
		return
	})

	service.HttpHandle("/tenant/{id}", http.MethodPost, "dbapi.admin", func(ctx context.Context, in string) (out int64, err error) {
		dsn := service.CtxVars(ctx)["id"]
		out, err = service.DataHSet("dsn", dsn, in)
		return
	})

	service.HttpHandle("/tenant/{id}", http.MethodDelete, "dbapi.admin", func(ctx context.Context, in string) (out int64, err error) {
		dsn := service.CtxVars(ctx)["id"]
		out, err = service.DataHDel("dsn", dsn)
		return
	})

	service.HttpHandle("/admin", http.MethodPost, "dbapi.admin", func(ctx context.Context, in *types.DBIn) (out *types.DBOut, err error) {
		if in.T == "" {
			in.T = "default"
		}

		db, err := service.DBN(in.T)

		if err != nil {
			return
		}

		out = &types.DBOut{}
		ret := make([]map[string]interface{}, 0)
		ierr := db.Raw(in.Q, in.Params...).Scan(&ret).Error
		out.Data = ret
		if ierr != nil {
			out.Err = ierr.Error()
		}

		return
	})

	service.HttpHandle("/dbs", http.MethodGet, "dbapi.admin", func(ctx context.Context, in *service.EMPTY_TYPE) (out *types.DBOut, err error) {

		db, err := service.DBN("default")

		if err != nil {
			return
		}

		out = &types.DBOut{}
		ret := make([]map[string]interface{}, 0)
		ierr := db.Raw("SELECT * FROM pg_database;").Scan(&ret).Error
		out.Data = ret
		if ierr != nil {
			out.Err = ierr.Error()
		}

		return
	})

	service.HttpHandle("/db/{n}/tables", http.MethodGet, "dbapi.admin", func(ctx context.Context, in *service.EMPTY_TYPE) (out *types.DBOut, err error) {

		t := service.CtxVars(ctx)["n"]

		db, err := service.DBN(t)

		if err != nil {
			return
		}

		out = &types.DBOut{}
		ret := make([]map[string]interface{}, 0)
		ierr := db.Raw(`SELECT *
		FROM pg_catalog.pg_tables
		WHERE schemaname != 'pg_catalog' AND 
			schemaname != 'information_schema';`).Scan(&ret).Error
		out.Data = ret
		if ierr != nil {
			out.Err = ierr.Error()
		}

		return
	})

	service.HttpHandle("/db/{n}/init", http.MethodGet, "dbapi.admin", func(ctx context.Context, in *service.EMPTY_TYPE) (out *types.DBOut, err error) {

		t := service.CtxVars(ctx)["n"]

		db, err := service.DBN(t)

		if err != nil {
			return
		}

		out = &types.DBOut{}
		db.AutoMigrate(&service.SecUser{})
		user := &service.SecUser{Username: "root"}
		passbs, err := bcrypt.GenerateFromPassword([]byte("Aa1234"), 0)
		if err != nil {
			return
		}
		ptrTrue := true
		err = db.Find(user).First(user).Error
		user.Hash = string(passbs)
		user.Enabled = &ptrTrue
		if err != nil {
			err = db.Create(user).Error
		} else {
			err = db.Save(user).Error
		}
		out.Data = "ok"
		return
	})

	service.HttpRun("")
	return nil

}
