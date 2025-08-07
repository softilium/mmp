package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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

var OrderStatusCaptions map[int64]string = map[int64]string{
	OrderStatusNew:            "Новый",
	OrderStatusInProcess:      "В работе",
	OrderStatusReadyToDeliver: "Готов к доставке",
	OrderStatusDelivering:     "Доставляется",
	OrderStatusDone:           "Доставлен",
	OrderStatusCanceled:       "Отменен",
}

var DB *DbContext

// global enrished http request context func
type contextKey string

const userContextKey contextKey = "user"

func AddUserContext(ctx context.Context, user *User) context.Context {
	if user != nil {
		ctx = context.WithValue(ctx, userContextKey, user)
	}
	return ctx
}

func LoadUserFromHttpToContext(r *http.Request) context.Context {
	user, _, err := UserFromHttpRequest(r)
	if err == nil {
		return AddUserContext(r.Context(), user)
	}
	return r.Context()
}

func (dbc *DbContext) SetHandlers() error {

	dbc.CustomerOrderDef.ExpectedDeliveryDate.DateTimeJSONFormat = time.DateOnly

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

	dbc.TagDef.AutoExpandFieldsForJSON = map[*elorm.FieldDef]bool{
		dbc.TagDef.Ref:   true,
		dbc.TagDef.Name:  true,
		dbc.TagDef.Color: true,
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

	//Tags
	//////

	err = dbc.AddBeforeDeleteHandler(dbc.TagDef.EntityDef, func(ctx context.Context, entity any) error {
		tg := entity.(*elorm.Entity)
		usages, _, err := dbc.GoodTagDef.SelectEntities([]*elorm.Filter{
			elorm.AddFilterEQ(dbc.GoodTagDef.Tag, tg)}, nil, 0, 0)
		if err != nil {
			return fmt.Errorf("error selecting good tags for tag: %v", err)
		}
		for _, usage := range usages {
			err = dbc.DeleteEntity(ctx, usage.RefString())
			if err != nil {
				return fmt.Errorf("error deleting good tag %s: %v", usage.RefString(), err)
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

	err = dbc.AddAfterSaveHandler(dbc.CustomerOrderDef.EntityDef, func(ctx context.Context, entity any) error {

		mig, ok := ctx.Value("migration").(bool)
		if ok && mig {
			return nil // skip for migration
		}

		order := entity.(*CustomerOrder)
		orderID := order.CreatedAt().Format(time.DateTime)

		cu, _ := ctx.Value(userContextKey).(*User)
		_ = cu // unused, but can be used for logging or other purposes

		if order.IsNew() {
			if order.Sender().TelegramVerified() {
				_, _ = Bot.Send(
					tgbotapi.NewMessage(order.Sender().TelegramChatId(),
						fmt.Sprintf("Новый заказ для вашей обработки. Заказчик [%s, Заказ [%s].", order.CreatedBy().Username(), orderID)))
			}
			if order.CreatedBy().TelegramVerified() {
				_, _ = Bot.Send(
					tgbotapi.NewMessage(order.CreatedBy().TelegramChatId(),
						fmt.Sprintf("Заказ [%s] отправлен вледельцу витрины. Бот будет уведомлять вас о всех будущих изменениях.", orderID)))
			}
		} else {
			if order.Status() != order.field_Status.GetOld() {
				if order.CreatedBy().TelegramVerified() {
					_, _ = Bot.Send(
						tgbotapi.NewMessage(order.CreatedBy().TelegramChatId(),
							fmt.Sprintf("Статус заказа [%s] изменен владельцем [%s] с [%s] на [%s].",
								orderID,
								order.Sender().Username(),
								OrderStatusCaptions[order.field_Status.GetOld()],
								OrderStatusCaptions[order.Status()],
							)))
				}
			}
			if order.CustomerComment() != order.field_CustomerComment.GetOld() {
				if order.Sender().TelegramVerified() {
					_, _ = Bot.Send(
						tgbotapi.NewMessage(order.Sender().TelegramChatId(),
							fmt.Sprintf("Заказчик указал примечание к заказу [%s]\n\r\n\r%s", orderID, order.CustomerComment())))
				}
			}
			if order.SenderComment() != order.field_SenderComment.GetOld() {
				if order.CreatedBy().TelegramVerified() {
					_, _ = Bot.Send(
						tgbotapi.NewMessage(order.CreatedBy().TelegramChatId(),
							fmt.Sprintf("Владелец витрины указал примечание к заказу [%s]\n\r\n\r%s", orderID, order.SenderComment())))
				}
			}
			if order.ExpectedDeliveryDate() != order.field_ExpectedDeliveryDate.GetOld() {
				if order.CreatedBy().TelegramVerified() {
					_, _ = Bot.Send(
						tgbotapi.NewMessage(order.CreatedBy().TelegramChatId(),
							fmt.Sprintf("Владелец витрины указал дату доставки к заказу [%s]\n\r\n\r%s", orderID, order.ExpectedDeliveryDate().Format(time.DateTime))))
				}
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	return nil

}
