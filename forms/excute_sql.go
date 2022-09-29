package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/pocketbase/pocketbase/core"
)

type ExcuteSql struct {
	app    core.App
	Did    string `form:"did" json:"did"`
	RawSql string `form:"rawSql" json:"rawSql"`
}

func NewExcuteSql(app core.App) *ExcuteSql {
	return &ExcuteSql{app: app}
}

func (form *ExcuteSql) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Did,
			validation.Required,
			validation.Length(1, 20),
			is.ASCII,
		),
		validation.Field(
			&form.RawSql,
			validation.Required,
		),
	)
}

//TODO 增加 check datasource的可用性， check sql 必须是select  参见 form 包下 user_upsert.go
