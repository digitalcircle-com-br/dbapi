package lib

import (
	"context"
	"fmt"
	"net/http"

	"github.com/digitalcircle-com-br/dbapi/lib/types"
	"github.com/digitalcircle-com-br/service"
	"golang.org/x/crypto/bcrypt"
)

const (
	PERM_DBADMIN service.PermDef = "dbapi.admin"
)

func Run() error {

	service.Init("dbapi")

	service.HttpHandle("/tenants", http.MethodGet, PERM_DBADMIN, func(ctx context.Context, in *service.EMPTY_TYPE) (out []string, err error) {
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

	service.HttpHandle("/tenant/{id}", http.MethodGet, PERM_DBADMIN, func(ctx context.Context, in *service.EMPTY_TYPE) (out string, err error) {
		dsn := service.CtxVars(ctx)["id"]
		out, err = service.DataHGet("dsn", dsn)
		return
	})

	service.HttpHandle("/tenant/{id}", http.MethodPost, PERM_DBADMIN, func(ctx context.Context, in string) (out int64, err error) {
		dsn := service.CtxVars(ctx)["id"]
		out, err = service.DataHSet("dsn", dsn, in)
		return
	})

	service.HttpHandle("/tenant/{id}", http.MethodDelete, PERM_DBADMIN, func(ctx context.Context, in string) (out int64, err error) {
		dsn := service.CtxVars(ctx)["id"]
		out, err = service.DataHDel("dsn", dsn)
		return
	})

	service.HttpHandle("/admin", http.MethodPost, PERM_DBADMIN, func(ctx context.Context, in *types.DBIn) (out *types.DBOut, err error) {
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

	service.HttpHandle("/dbs", http.MethodGet, PERM_DBADMIN, func(ctx context.Context, in *service.EMPTY_TYPE) (out *types.DBOut, err error) {

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

	service.HttpHandle("/db/{n}/tables", http.MethodGet, PERM_DBADMIN, func(ctx context.Context, in *service.EMPTY_TYPE) (out *types.DBOut, err error) {

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

	service.HttpHandle("/db/{n}/init", http.MethodGet, PERM_DBADMIN, func(ctx context.Context, in *service.EMPTY_TYPE) (out *types.DBOut, err error) {

		t := service.CtxVars(ctx)["n"]

		db, err := service.DBN(t)

		if err != nil {
			return
		}

		out = &types.DBOut{}
		db.AutoMigrate(&service.SecUser{})
		db.AutoMigrate(&service.SecGroup{})
		db.AutoMigrate(&service.SecPerm{})
		user := &service.SecUser{Username: "root"}
		passbs, err := bcrypt.GenerateFromPassword([]byte("Aa1234"), 0)
		if err != nil {
			return
		}
		ptrTrue := true
		err = db.Find(user).First(user).Error
		user.Hash = string(passbs)
		user.Enabled = &ptrTrue

		perm := &service.SecPerm{Name: "*", Val: "*"}
		group := &service.SecGroup{Name: "root", Perms: []*service.SecPerm{perm}}
		user.Groups = []*service.SecGroup{group}

		if err != nil {

			err = db.Create(user).Error
		} else {
			err = db.Save(user).Error
		}
		out.Data = "ok"
		return
	})

	service.HttpHandle("/db/{n}/create", http.MethodPost, PERM_DBADMIN, func(ctx context.Context, in string) (out *types.DBOut, err error) {

		dsn := service.CtxVars(ctx)["n"]
		_, err = service.DataHSet("dsn", dsn, in)
		if err != nil {
			return
		}

		d, err := service.DBN("master")
		if err != nil {
			return
		}

		err = d.Exec(fmt.Sprintf("create database %s", dsn)).Error
		if err != nil {
			return
		}

		db, err := service.DBN(dsn)

		if err != nil {
			return
		}

		out = &types.DBOut{}

		err = db.AutoMigrate(&service.SecPerm{})
		if err != nil {
			return
		}

		err = db.AutoMigrate(&service.SecGroup{})
		if err != nil {
			return
		}

		err = db.AutoMigrate(&service.SecUser{})
		if err != nil {
			return
		}

		user := &service.SecUser{Username: "root"}
		passbs, err := bcrypt.GenerateFromPassword([]byte("Aa1234"), 0)
		if err != nil {
			return
		}
		ptrTrue := true
		err = db.Find(user).First(user).Error
		user.Hash = string(passbs)
		user.Enabled = &ptrTrue

		perm := &service.SecPerm{Name: "*", Val: "*"}
		group := &service.SecGroup{Name: "root", Perms: []*service.SecPerm{perm}}
		user.Groups = []*service.SecGroup{group}

		if err != nil {

			err = db.Create(user).Error
		} else {
			err = db.Save(user).Error
		}
		out.Data = "ok"
		return
	})

	service.HttpHandle("/db/{n}/drop", http.MethodGet, PERM_DBADMIN, func(ctx context.Context, in *service.EMPTY_TYPE) (out *service.EMPTY_TYPE, err error) {

		dsn := service.CtxVars(ctx)["n"]
		_, err = service.DataHSet("dsn", dsn, in)
		if err != nil {
			return
		}

		d, err := service.DB()
		if err != nil {
			return
		}

		err = d.Raw("drop database " + dsn).Error
		if err != nil {
			return
		}
		_, err = service.DataHDel("dsn", dsn)
		return
	})

	service.HttpRun("")
	return nil

}
