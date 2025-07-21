package models

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/softilium/elorm"
)

var Dbc *DbContext

// global enrished http request context func
type contextKey string

const UserContextKey contextKey = "user"

func HttpUserContext(r *http.Request) context.Context {
	ctx := r.Context()
	user, _ := UserFromHttpRequest(r)
	if user != nil {
		ctx = context.WithValue(ctx, UserContextKey, user)
	}
	return ctx
}

func UserFromHttpRequest(r *http.Request) (*User, error) {

	raw := r.Header.Get("Authorization")
	if raw == "" {
		return nil, fmt.Errorf("Unauthorized")
	}
	if len(raw) < 7 || raw[:7] != "Bearer " {
		return nil, fmt.Errorf("Unauthorized")
	}
	token := raw[7:]

	if item, ok := TokensByAT[token]; ok {
		if item.AccessTokenExpiresAt.After(time.Now()) {
			user, err := Dbc.LoadUser(item.UserRef)
			if err != nil {
				return nil, fmt.Errorf("Unauthorized")
			}
			if !user.IsActive() {
				return nil, fmt.Errorf("user isn't active")
			}
			return user, nil
		}
	}
	return nil, fmt.Errorf("Unauthorized")

}

type TokenItem struct {
	UserRef               string
	RefreshToken          string
	AccessTokenExpiresAt  time.Time
	RefreshTokenExpiresAt time.Time
}

var TokensByAT = make(map[string]TokenItem)

func (dbc *DbContext) SetHandlers() error {

	// BusinessObjects fragment
	///////////////////////////

	err := dbc.AddBeforeSaveHandler(BusinessObjectsFragment, func(ctx context.Context, entity any) error {

		mig, ok := ctx.Value("migration").(bool)
		if ok && mig {
			return nil // skip for migration
		}

		user, ok := ctx.Value(UserContextKey).(*User)
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
		} else {
			et.SetModifiedAt(time.Now())
			if user != nil {
				et.SetModifiedBy(user)
			}
		}

		isDeletedV := ent.GetValues()[elorm.IsDeletedFieldName].(*elorm.FieldValueBool)

		if !ent.IsDeleted() && isDeletedV.GetOld() {
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
