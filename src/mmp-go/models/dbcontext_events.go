package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/softilium/elorm"
)

const (
	OrderStatusNew            = 100
	OrderStatusInProcess      = 200
	OrderStatusReadyToDeliver = 300
	OrderStatusDelivering     = 400
	OrderStatusDone           = 500
	OrderStatusCanceled       = 600
)

var Dbc *DbContext

// global enrished http request context func
type contextKey string

const userContextKey contextKey = "user"

func HttpUserContext(r *http.Request) context.Context {
	ctx := r.Context()
	user, _ := UserFromHttpRequest(r)
	if user != nil {
		ctx = context.WithValue(ctx, userContextKey, user)
	}
	return ctx
}

func (dbc *DbContext) SetHandlers() error {

	dbc.UserDef.AutoExpandFieldsForJSON = map[*elorm.FieldDef]bool{
		dbc.UserDef.Ref:      true,
		dbc.UserDef.Username: true,
	}

	dbc.ShopDef.AutoExpandFieldsForJSON = map[*elorm.FieldDef]bool{
		dbc.ShopDef.Ref:     true,
		dbc.ShopDef.Caption: true,
	}

	dbc.GoodDef.AutoExpandFieldsForJSON = map[*elorm.FieldDef]bool{
		dbc.GoodDef.Ref:       true,
		dbc.GoodDef.Caption:   true,
		dbc.GoodDef.CreatedBy: true,
	}

	// BusinessObjects fragment
	///////////////////////////

	err := dbc.AddBeforeSaveHandler(BusinessObjectsFragment, func(ctx context.Context, entity any) error {

		mig, ok := ctx.Value("migration").(bool)
		if ok && mig {
			return nil // skip for migration
		}

		user, ok := ctx.Value(userContextKey).(*User)
		if !ok {
			user = nil
		}

		et := entity.(BusinessObjectsFragmentMethods)
		ent := entity.(elorm.IEntity)
		if ent.IsNew() {
			et.SetCreatedAt(time.Now())
			if user != nil {
				et.SetCreatedBy(user)
			}
			if et.CreatedBy() == nil {
				return fmt.Errorf("createdBy is required for new entity %s", ent.Def().ObjectName)
			}
		} else {
			et.SetModifiedAt(time.Now())
			if user != nil {
				et.SetModifiedBy(user)
			}
		}

		fld := ent.GetValues()[elorm.IsDeletedFieldName].(*elorm.FieldValueBool)
		if !fld.Get() && fld.GetOld() {
			et.SetDeletedAt(time.Now())
			if user != nil {
				et.SetDeletedBy(user)
			}
		}

		return nil

	})
	if err != nil {
		return err
	}

	// Users
	////////

	err = dbc.AddFillNewHandler(dbc.UserDef.EntityDef, func(entity any) error {
		user := entity.(*User)
		user.SetIsActive(true)
		return nil
	})
	if err != nil {
		return err
	}

	return nil

}
