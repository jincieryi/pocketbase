package forms

import (
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/pocketbase/pocketbase/core"
	"github.com/pocketbase/pocketbase/models"
)

type SqlCollectionUpsert struct {
	*CollectionUpsert
}

func NewSqlCollectionUpsert(app core.App, collection *models.Collection) *SqlCollectionUpsert {
	return &SqlCollectionUpsert{NewCollectionUpsert(app, collection)}
}

// Validate makes the form validatable by implementing [validation.Validatable] interface.
func (form *SqlCollectionUpsert) Validate() error {
	return validation.ValidateStruct(form,
		validation.Field(
			&form.Id,
			validation.When(
				form.collection.IsNew(),
				validation.Length(models.DefaultIdLength, models.DefaultIdLength),
				validation.Match(idRegex),
			).Else(validation.In(form.collection.Id)),
		),
		validation.Field(
			&form.System,
			validation.By(form.ensureNoSystemFlagChange),
		),
		validation.Field(
			&form.Name,
			validation.Required,
			validation.Length(1, 255),
			validation.Match(collectionNameRegex),
			validation.By(form.ensureNoSystemNameChange),
			validation.By(form.checkUniqueName),
		),
		// 跳过schema的校验
		validation.Field(
			&form.Schema,
			validation.Skip,
		),
		validation.Field(&form.ListRule, validation.By(form.checkRule)),
		validation.Field(&form.ViewRule, validation.By(form.checkRule)),
		validation.Field(&form.CreateRule, validation.By(form.checkRule)),
		validation.Field(&form.UpdateRule, validation.By(form.checkRule)),
		validation.Field(&form.DeleteRule, validation.By(form.checkRule)),
	)
}

// Submit validates the form and upserts the form's Collection model.
//
// On success the related record table schema will be auto updated.
//
// You can optionally provide a list of InterceptorFunc to further
// modify the form behavior before persisting it.
func (form *SqlCollectionUpsert) Submit(interceptors ...InterceptorFunc) error {
	if err := form.Validate(); err != nil {
		return err
	}

	if form.collection.IsNew() {
		// system flag can be set only on create
		form.collection.System = form.System

		// custom insertion id can be set only on create
		if form.Id != "" {
			form.collection.MarkAsNew()
			form.collection.SetId(form.Id)
		}
	}

	// system collections cannot be renamed
	if form.collection.IsNew() || !form.collection.System {
		form.collection.Name = form.Name
	}

	form.collection.Schema = form.Schema
	form.collection.ListRule = form.ListRule
	form.collection.ViewRule = form.ViewRule
	form.collection.CreateRule = form.CreateRule
	form.collection.UpdateRule = form.UpdateRule
	form.collection.DeleteRule = form.DeleteRule

	return runInterceptors(func() error {
		return form.config.Dao.SaveSqlCollection(form.collection)
	}, interceptors...)
}
