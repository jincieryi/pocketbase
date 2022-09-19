package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
)

type TestConnectDB struct {
	app core.App

	Name string `form:"name" json:"name"`
	Dsn  string `form:"dsn" json:"dsn"`
	Type string `form:"type" json:"type"`
}

func NewTestConnectDB(app core.App) *TestConnectDB {
	return &TestConnectDB{app: app}
}

func (form *TestConnectDB) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Dsn,
			validation.Required,
			validation.Length(1, 255),
			is.DNSName,
		))
}
