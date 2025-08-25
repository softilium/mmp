// THIS IS GENERATED CODE. DO NOT MODIFY IT
package main

import (
	"github.com/softilium/elorm"
	"time"
)

// CustomerOrder class
//////

type CustomerOrderDefStruct struct {
	*elorm.EntityDef
	Ref         *elorm.FieldDef
	IsDeleted   *elorm.FieldDef
	DataVersion *elorm.FieldDef

	Sender *elorm.FieldDef

	Status *elorm.FieldDef

	Qty *elorm.FieldDef

	Sum *elorm.FieldDef

	OrderDiscountSum *elorm.FieldDef

	SenderComment *elorm.FieldDef

	CustomerComment *elorm.FieldDef

	ExpectedDeliveryDate *elorm.FieldDef

	CreatedBy *elorm.FieldDef

	CreatedAt *elorm.FieldDef

	ModifiedBy *elorm.FieldDef

	ModifiedAt *elorm.FieldDef

	DeletedBy *elorm.FieldDef

	DeletedAt *elorm.FieldDef
}

func (T *CustomerOrderDefStruct) SelectEntities(filters []*elorm.Filter, sorts []*elorm.SortItem, pageNo int, pageSize int) (result []*CustomerOrder, pages int, err error) {

	res, total, err := T.EntityDef.SelectEntities(filters, sorts, pageNo, pageSize)
	if err != nil {
		return nil, 0, err
	}

	res2 := make([]*CustomerOrder, 0, len(res))

	for _, r := range res {
		if r == nil {
			continue
		}
		rt := T.Wrap(r)
		res2 = append(res2, rt.(*CustomerOrder))
	}

	return res2, total, nil

}

type CustomerOrder struct {
	*elorm.Entity

	field_Sender               *elorm.FieldValueRef
	field_Status               *elorm.FieldValueInt
	field_Qty                  *elorm.FieldValueNumeric
	field_Sum                  *elorm.FieldValueNumeric
	field_OrderDiscountSum     *elorm.FieldValueNumeric
	field_SenderComment        *elorm.FieldValueString
	field_CustomerComment      *elorm.FieldValueString
	field_ExpectedDeliveryDate *elorm.FieldValueDateTime
	field_CreatedBy            *elorm.FieldValueRef
	field_CreatedAt            *elorm.FieldValueDateTime
	field_ModifiedBy           *elorm.FieldValueRef
	field_ModifiedAt           *elorm.FieldValueDateTime
	field_DeletedBy            *elorm.FieldValueRef
	field_DeletedAt            *elorm.FieldValueDateTime
}

func (T *CustomerOrder) Sender() *User {
	if T.field_Sender == nil {
		T.field_Sender = T.Values["Sender"].(*elorm.FieldValueRef)
	}
	r, err := T.field_Sender.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *CustomerOrder) SetSender(newValue *User) {
	if T.field_Sender == nil {
		T.field_Sender = T.Values["Sender"].(*elorm.FieldValueRef)
	}
	err := T.field_Sender.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *CustomerOrder) Status() int64 {
	if T.field_Status == nil {
		T.field_Status = T.Values["Status"].(*elorm.FieldValueInt)
	}
	return T.field_Status.Get()
}

func (T *CustomerOrder) SetStatus(newValue int64) {
	if T.field_Status == nil {
		T.field_Status = T.Values["Status"].(*elorm.FieldValueInt)
	}
	T.field_Status.Set(newValue)
}

func (T *CustomerOrder) Qty() float64 {
	if T.field_Qty == nil {
		T.field_Qty = T.Values["Qty"].(*elorm.FieldValueNumeric)
	}
	return T.field_Qty.Get()
}

func (T *CustomerOrder) SetQty(newValue float64) {
	if T.field_Qty == nil {
		T.field_Qty = T.Values["Qty"].(*elorm.FieldValueNumeric)
	}
	T.field_Qty.Set(newValue)
}

func (T *CustomerOrder) Sum() float64 {
	if T.field_Sum == nil {
		T.field_Sum = T.Values["Sum"].(*elorm.FieldValueNumeric)
	}
	return T.field_Sum.Get()
}

func (T *CustomerOrder) SetSum(newValue float64) {
	if T.field_Sum == nil {
		T.field_Sum = T.Values["Sum"].(*elorm.FieldValueNumeric)
	}
	T.field_Sum.Set(newValue)
}

func (T *CustomerOrder) OrderDiscountSum() float64 {
	if T.field_OrderDiscountSum == nil {
		T.field_OrderDiscountSum = T.Values["OrderDiscountSum"].(*elorm.FieldValueNumeric)
	}
	return T.field_OrderDiscountSum.Get()
}

func (T *CustomerOrder) SetOrderDiscountSum(newValue float64) {
	if T.field_OrderDiscountSum == nil {
		T.field_OrderDiscountSum = T.Values["OrderDiscountSum"].(*elorm.FieldValueNumeric)
	}
	T.field_OrderDiscountSum.Set(newValue)
}

func (T *CustomerOrder) SenderComment() string {
	if T.field_SenderComment == nil {
		T.field_SenderComment = T.Values["SenderComment"].(*elorm.FieldValueString)
	}
	return T.field_SenderComment.Get()
}

func (T *CustomerOrder) SetSenderComment(newValue string) {
	if T.field_SenderComment == nil {
		T.field_SenderComment = T.Values["SenderComment"].(*elorm.FieldValueString)
	}
	T.field_SenderComment.Set(newValue)
}

func (T *CustomerOrder) CustomerComment() string {
	if T.field_CustomerComment == nil {
		T.field_CustomerComment = T.Values["CustomerComment"].(*elorm.FieldValueString)
	}
	return T.field_CustomerComment.Get()
}

func (T *CustomerOrder) SetCustomerComment(newValue string) {
	if T.field_CustomerComment == nil {
		T.field_CustomerComment = T.Values["CustomerComment"].(*elorm.FieldValueString)
	}
	T.field_CustomerComment.Set(newValue)
}

func (T *CustomerOrder) ExpectedDeliveryDate() time.Time {
	if T.field_ExpectedDeliveryDate == nil {
		T.field_ExpectedDeliveryDate = T.Values["ExpectedDeliveryDate"].(*elorm.FieldValueDateTime)
	}
	return T.field_ExpectedDeliveryDate.Get()
}

func (T *CustomerOrder) SetExpectedDeliveryDate(newValue time.Time) {
	if T.field_ExpectedDeliveryDate == nil {
		T.field_ExpectedDeliveryDate = T.Values["ExpectedDeliveryDate"].(*elorm.FieldValueDateTime)
	}
	T.field_ExpectedDeliveryDate.Set(newValue)
}

func (T *CustomerOrder) CreatedBy() *User {
	if T.field_CreatedBy == nil {
		T.field_CreatedBy = T.Values["CreatedBy"].(*elorm.FieldValueRef)
	}
	r, err := T.field_CreatedBy.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *CustomerOrder) SetCreatedBy(newValue *User) {
	if T.field_CreatedBy == nil {
		T.field_CreatedBy = T.Values["CreatedBy"].(*elorm.FieldValueRef)
	}
	err := T.field_CreatedBy.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *CustomerOrder) CreatedAt() time.Time {
	if T.field_CreatedAt == nil {
		T.field_CreatedAt = T.Values["CreatedAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_CreatedAt.Get()
}

func (T *CustomerOrder) SetCreatedAt(newValue time.Time) {
	if T.field_CreatedAt == nil {
		T.field_CreatedAt = T.Values["CreatedAt"].(*elorm.FieldValueDateTime)
	}
	T.field_CreatedAt.Set(newValue)
}

func (T *CustomerOrder) ModifiedBy() *User {
	if T.field_ModifiedBy == nil {
		T.field_ModifiedBy = T.Values["ModifiedBy"].(*elorm.FieldValueRef)
	}
	r, err := T.field_ModifiedBy.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *CustomerOrder) SetModifiedBy(newValue *User) {
	if T.field_ModifiedBy == nil {
		T.field_ModifiedBy = T.Values["ModifiedBy"].(*elorm.FieldValueRef)
	}
	err := T.field_ModifiedBy.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *CustomerOrder) ModifiedAt() time.Time {
	if T.field_ModifiedAt == nil {
		T.field_ModifiedAt = T.Values["ModifiedAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_ModifiedAt.Get()
}

func (T *CustomerOrder) SetModifiedAt(newValue time.Time) {
	if T.field_ModifiedAt == nil {
		T.field_ModifiedAt = T.Values["ModifiedAt"].(*elorm.FieldValueDateTime)
	}
	T.field_ModifiedAt.Set(newValue)
}

func (T *CustomerOrder) DeletedBy() *User {
	if T.field_DeletedBy == nil {
		T.field_DeletedBy = T.Values["DeletedBy"].(*elorm.FieldValueRef)
	}
	r, err := T.field_DeletedBy.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *CustomerOrder) SetDeletedBy(newValue *User) {
	if T.field_DeletedBy == nil {
		T.field_DeletedBy = T.Values["DeletedBy"].(*elorm.FieldValueRef)
	}
	err := T.field_DeletedBy.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *CustomerOrder) DeletedAt() time.Time {
	if T.field_DeletedAt == nil {
		T.field_DeletedAt = T.Values["DeletedAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_DeletedAt.Get()
}

func (T *CustomerOrder) SetDeletedAt(newValue time.Time) {
	if T.field_DeletedAt == nil {
		T.field_DeletedAt = T.Values["DeletedAt"].(*elorm.FieldValueDateTime)
	}
	T.field_DeletedAt.Set(newValue)
}

// Good class
//////

type GoodDefStruct struct {
	*elorm.EntityDef
	Ref         *elorm.FieldDef
	IsDeleted   *elorm.FieldDef
	DataVersion *elorm.FieldDef

	OwnerShop *elorm.FieldDef

	Caption *elorm.FieldDef

	Article *elorm.FieldDef

	Url *elorm.FieldDef

	Description *elorm.FieldDef

	Price *elorm.FieldDef

	OrderInShop *elorm.FieldDef

	CreatedBy *elorm.FieldDef

	CreatedAt *elorm.FieldDef

	ModifiedBy *elorm.FieldDef

	ModifiedAt *elorm.FieldDef

	DeletedBy *elorm.FieldDef

	DeletedAt *elorm.FieldDef
}

func (T *GoodDefStruct) SelectEntities(filters []*elorm.Filter, sorts []*elorm.SortItem, pageNo int, pageSize int) (result []*Good, pages int, err error) {

	res, total, err := T.EntityDef.SelectEntities(filters, sorts, pageNo, pageSize)
	if err != nil {
		return nil, 0, err
	}

	res2 := make([]*Good, 0, len(res))

	for _, r := range res {
		if r == nil {
			continue
		}
		rt := T.Wrap(r)
		res2 = append(res2, rt.(*Good))
	}

	return res2, total, nil

}

type Good struct {
	*elorm.Entity

	field_OwnerShop   *elorm.FieldValueRef
	field_Caption     *elorm.FieldValueString
	field_Article     *elorm.FieldValueString
	field_Url         *elorm.FieldValueString
	field_Description *elorm.FieldValueString
	field_Price       *elorm.FieldValueNumeric
	field_OrderInShop *elorm.FieldValueInt
	field_CreatedBy   *elorm.FieldValueRef
	field_CreatedAt   *elorm.FieldValueDateTime
	field_ModifiedBy  *elorm.FieldValueRef
	field_ModifiedAt  *elorm.FieldValueDateTime
	field_DeletedBy   *elorm.FieldValueRef
	field_DeletedAt   *elorm.FieldValueDateTime
}

func (T *Good) OwnerShop() *Shop {
	if T.field_OwnerShop == nil {
		T.field_OwnerShop = T.Values["OwnerShop"].(*elorm.FieldValueRef)
	}
	r, err := T.field_OwnerShop.Get()
	if err != nil {
		panic(err)
	}
	return r.(*Shop)
}

func (T *Good) SetOwnerShop(newValue *Shop) {
	if T.field_OwnerShop == nil {
		T.field_OwnerShop = T.Values["OwnerShop"].(*elorm.FieldValueRef)
	}
	err := T.field_OwnerShop.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *Good) Caption() string {
	if T.field_Caption == nil {
		T.field_Caption = T.Values["Caption"].(*elorm.FieldValueString)
	}
	return T.field_Caption.Get()
}

func (T *Good) SetCaption(newValue string) {
	if T.field_Caption == nil {
		T.field_Caption = T.Values["Caption"].(*elorm.FieldValueString)
	}
	T.field_Caption.Set(newValue)
}

func (T *Good) Article() string {
	if T.field_Article == nil {
		T.field_Article = T.Values["Article"].(*elorm.FieldValueString)
	}
	return T.field_Article.Get()
}

func (T *Good) SetArticle(newValue string) {
	if T.field_Article == nil {
		T.field_Article = T.Values["Article"].(*elorm.FieldValueString)
	}
	T.field_Article.Set(newValue)
}

func (T *Good) Url() string {
	if T.field_Url == nil {
		T.field_Url = T.Values["Url"].(*elorm.FieldValueString)
	}
	return T.field_Url.Get()
}

func (T *Good) SetUrl(newValue string) {
	if T.field_Url == nil {
		T.field_Url = T.Values["Url"].(*elorm.FieldValueString)
	}
	T.field_Url.Set(newValue)
}

func (T *Good) Description() string {
	if T.field_Description == nil {
		T.field_Description = T.Values["Description"].(*elorm.FieldValueString)
	}
	return T.field_Description.Get()
}

func (T *Good) SetDescription(newValue string) {
	if T.field_Description == nil {
		T.field_Description = T.Values["Description"].(*elorm.FieldValueString)
	}
	T.field_Description.Set(newValue)
}

func (T *Good) Price() float64 {
	if T.field_Price == nil {
		T.field_Price = T.Values["Price"].(*elorm.FieldValueNumeric)
	}
	return T.field_Price.Get()
}

func (T *Good) SetPrice(newValue float64) {
	if T.field_Price == nil {
		T.field_Price = T.Values["Price"].(*elorm.FieldValueNumeric)
	}
	T.field_Price.Set(newValue)
}

func (T *Good) OrderInShop() int64 {
	if T.field_OrderInShop == nil {
		T.field_OrderInShop = T.Values["OrderInShop"].(*elorm.FieldValueInt)
	}
	return T.field_OrderInShop.Get()
}

func (T *Good) SetOrderInShop(newValue int64) {
	if T.field_OrderInShop == nil {
		T.field_OrderInShop = T.Values["OrderInShop"].(*elorm.FieldValueInt)
	}
	T.field_OrderInShop.Set(newValue)
}

func (T *Good) CreatedBy() *User {
	if T.field_CreatedBy == nil {
		T.field_CreatedBy = T.Values["CreatedBy"].(*elorm.FieldValueRef)
	}
	r, err := T.field_CreatedBy.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *Good) SetCreatedBy(newValue *User) {
	if T.field_CreatedBy == nil {
		T.field_CreatedBy = T.Values["CreatedBy"].(*elorm.FieldValueRef)
	}
	err := T.field_CreatedBy.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *Good) CreatedAt() time.Time {
	if T.field_CreatedAt == nil {
		T.field_CreatedAt = T.Values["CreatedAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_CreatedAt.Get()
}

func (T *Good) SetCreatedAt(newValue time.Time) {
	if T.field_CreatedAt == nil {
		T.field_CreatedAt = T.Values["CreatedAt"].(*elorm.FieldValueDateTime)
	}
	T.field_CreatedAt.Set(newValue)
}

func (T *Good) ModifiedBy() *User {
	if T.field_ModifiedBy == nil {
		T.field_ModifiedBy = T.Values["ModifiedBy"].(*elorm.FieldValueRef)
	}
	r, err := T.field_ModifiedBy.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *Good) SetModifiedBy(newValue *User) {
	if T.field_ModifiedBy == nil {
		T.field_ModifiedBy = T.Values["ModifiedBy"].(*elorm.FieldValueRef)
	}
	err := T.field_ModifiedBy.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *Good) ModifiedAt() time.Time {
	if T.field_ModifiedAt == nil {
		T.field_ModifiedAt = T.Values["ModifiedAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_ModifiedAt.Get()
}

func (T *Good) SetModifiedAt(newValue time.Time) {
	if T.field_ModifiedAt == nil {
		T.field_ModifiedAt = T.Values["ModifiedAt"].(*elorm.FieldValueDateTime)
	}
	T.field_ModifiedAt.Set(newValue)
}

func (T *Good) DeletedBy() *User {
	if T.field_DeletedBy == nil {
		T.field_DeletedBy = T.Values["DeletedBy"].(*elorm.FieldValueRef)
	}
	r, err := T.field_DeletedBy.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *Good) SetDeletedBy(newValue *User) {
	if T.field_DeletedBy == nil {
		T.field_DeletedBy = T.Values["DeletedBy"].(*elorm.FieldValueRef)
	}
	err := T.field_DeletedBy.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *Good) DeletedAt() time.Time {
	if T.field_DeletedAt == nil {
		T.field_DeletedAt = T.Values["DeletedAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_DeletedAt.Get()
}

func (T *Good) SetDeletedAt(newValue time.Time) {
	if T.field_DeletedAt == nil {
		T.field_DeletedAt = T.Values["DeletedAt"].(*elorm.FieldValueDateTime)
	}
	T.field_DeletedAt.Set(newValue)
}

// GoodTag class
//////

type GoodTagDefStruct struct {
	*elorm.EntityDef
	Ref         *elorm.FieldDef
	IsDeleted   *elorm.FieldDef
	DataVersion *elorm.FieldDef

	Good *elorm.FieldDef

	Tag *elorm.FieldDef
}

func (T *GoodTagDefStruct) SelectEntities(filters []*elorm.Filter, sorts []*elorm.SortItem, pageNo int, pageSize int) (result []*GoodTag, pages int, err error) {

	res, total, err := T.EntityDef.SelectEntities(filters, sorts, pageNo, pageSize)
	if err != nil {
		return nil, 0, err
	}

	res2 := make([]*GoodTag, 0, len(res))

	for _, r := range res {
		if r == nil {
			continue
		}
		rt := T.Wrap(r)
		res2 = append(res2, rt.(*GoodTag))
	}

	return res2, total, nil

}

type GoodTag struct {
	*elorm.Entity

	field_Good *elorm.FieldValueRef
	field_Tag  *elorm.FieldValueRef
}

func (T *GoodTag) Good() *Good {
	if T.field_Good == nil {
		T.field_Good = T.Values["Good"].(*elorm.FieldValueRef)
	}
	r, err := T.field_Good.Get()
	if err != nil {
		panic(err)
	}
	return r.(*Good)
}

func (T *GoodTag) SetGood(newValue *Good) {
	if T.field_Good == nil {
		T.field_Good = T.Values["Good"].(*elorm.FieldValueRef)
	}
	err := T.field_Good.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *GoodTag) Tag() *Tag {
	if T.field_Tag == nil {
		T.field_Tag = T.Values["Tag"].(*elorm.FieldValueRef)
	}
	r, err := T.field_Tag.Get()
	if err != nil {
		panic(err)
	}
	return r.(*Tag)
}

func (T *GoodTag) SetTag(newValue *Tag) {
	if T.field_Tag == nil {
		T.field_Tag = T.Values["Tag"].(*elorm.FieldValueRef)
	}
	err := T.field_Tag.Set(newValue)
	if err != nil {
		panic(err)
	}
}

// OrderLine class
//////

type OrderLineDefStruct struct {
	*elorm.EntityDef
	Ref         *elorm.FieldDef
	IsDeleted   *elorm.FieldDef
	DataVersion *elorm.FieldDef

	Shop *elorm.FieldDef

	CustomerOrder *elorm.FieldDef

	Good *elorm.FieldDef

	Qty *elorm.FieldDef

	LineDiscountSum *elorm.FieldDef

	Sum *elorm.FieldDef

	CreatedBy *elorm.FieldDef

	CreatedAt *elorm.FieldDef

	ModifiedBy *elorm.FieldDef

	ModifiedAt *elorm.FieldDef

	DeletedBy *elorm.FieldDef

	DeletedAt *elorm.FieldDef
}

func (T *OrderLineDefStruct) SelectEntities(filters []*elorm.Filter, sorts []*elorm.SortItem, pageNo int, pageSize int) (result []*OrderLine, pages int, err error) {

	res, total, err := T.EntityDef.SelectEntities(filters, sorts, pageNo, pageSize)
	if err != nil {
		return nil, 0, err
	}

	res2 := make([]*OrderLine, 0, len(res))

	for _, r := range res {
		if r == nil {
			continue
		}
		rt := T.Wrap(r)
		res2 = append(res2, rt.(*OrderLine))
	}

	return res2, total, nil

}

type OrderLine struct {
	*elorm.Entity

	field_Shop            *elorm.FieldValueRef
	field_CustomerOrder   *elorm.FieldValueRef
	field_Good            *elorm.FieldValueRef
	field_Qty             *elorm.FieldValueNumeric
	field_LineDiscountSum *elorm.FieldValueNumeric
	field_Sum             *elorm.FieldValueNumeric
	field_CreatedBy       *elorm.FieldValueRef
	field_CreatedAt       *elorm.FieldValueDateTime
	field_ModifiedBy      *elorm.FieldValueRef
	field_ModifiedAt      *elorm.FieldValueDateTime
	field_DeletedBy       *elorm.FieldValueRef
	field_DeletedAt       *elorm.FieldValueDateTime
}

func (T *OrderLine) Shop() *Shop {
	if T.field_Shop == nil {
		T.field_Shop = T.Values["Shop"].(*elorm.FieldValueRef)
	}
	r, err := T.field_Shop.Get()
	if err != nil {
		panic(err)
	}
	return r.(*Shop)
}

func (T *OrderLine) SetShop(newValue *Shop) {
	if T.field_Shop == nil {
		T.field_Shop = T.Values["Shop"].(*elorm.FieldValueRef)
	}
	err := T.field_Shop.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *OrderLine) CustomerOrder() *CustomerOrder {
	if T.field_CustomerOrder == nil {
		T.field_CustomerOrder = T.Values["CustomerOrder"].(*elorm.FieldValueRef)
	}
	r, err := T.field_CustomerOrder.Get()
	if err != nil {
		panic(err)
	}
	return r.(*CustomerOrder)
}

func (T *OrderLine) SetCustomerOrder(newValue *CustomerOrder) {
	if T.field_CustomerOrder == nil {
		T.field_CustomerOrder = T.Values["CustomerOrder"].(*elorm.FieldValueRef)
	}
	err := T.field_CustomerOrder.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *OrderLine) Good() *Good {
	if T.field_Good == nil {
		T.field_Good = T.Values["Good"].(*elorm.FieldValueRef)
	}
	r, err := T.field_Good.Get()
	if err != nil {
		panic(err)
	}
	return r.(*Good)
}

func (T *OrderLine) SetGood(newValue *Good) {
	if T.field_Good == nil {
		T.field_Good = T.Values["Good"].(*elorm.FieldValueRef)
	}
	err := T.field_Good.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *OrderLine) Qty() float64 {
	if T.field_Qty == nil {
		T.field_Qty = T.Values["Qty"].(*elorm.FieldValueNumeric)
	}
	return T.field_Qty.Get()
}

func (T *OrderLine) SetQty(newValue float64) {
	if T.field_Qty == nil {
		T.field_Qty = T.Values["Qty"].(*elorm.FieldValueNumeric)
	}
	T.field_Qty.Set(newValue)
}

func (T *OrderLine) LineDiscountSum() float64 {
	if T.field_LineDiscountSum == nil {
		T.field_LineDiscountSum = T.Values["LineDiscountSum"].(*elorm.FieldValueNumeric)
	}
	return T.field_LineDiscountSum.Get()
}

func (T *OrderLine) SetLineDiscountSum(newValue float64) {
	if T.field_LineDiscountSum == nil {
		T.field_LineDiscountSum = T.Values["LineDiscountSum"].(*elorm.FieldValueNumeric)
	}
	T.field_LineDiscountSum.Set(newValue)
}

func (T *OrderLine) Sum() float64 {
	if T.field_Sum == nil {
		T.field_Sum = T.Values["Sum"].(*elorm.FieldValueNumeric)
	}
	return T.field_Sum.Get()
}

func (T *OrderLine) SetSum(newValue float64) {
	if T.field_Sum == nil {
		T.field_Sum = T.Values["Sum"].(*elorm.FieldValueNumeric)
	}
	T.field_Sum.Set(newValue)
}

func (T *OrderLine) CreatedBy() *User {
	if T.field_CreatedBy == nil {
		T.field_CreatedBy = T.Values["CreatedBy"].(*elorm.FieldValueRef)
	}
	r, err := T.field_CreatedBy.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *OrderLine) SetCreatedBy(newValue *User) {
	if T.field_CreatedBy == nil {
		T.field_CreatedBy = T.Values["CreatedBy"].(*elorm.FieldValueRef)
	}
	err := T.field_CreatedBy.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *OrderLine) CreatedAt() time.Time {
	if T.field_CreatedAt == nil {
		T.field_CreatedAt = T.Values["CreatedAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_CreatedAt.Get()
}

func (T *OrderLine) SetCreatedAt(newValue time.Time) {
	if T.field_CreatedAt == nil {
		T.field_CreatedAt = T.Values["CreatedAt"].(*elorm.FieldValueDateTime)
	}
	T.field_CreatedAt.Set(newValue)
}

func (T *OrderLine) ModifiedBy() *User {
	if T.field_ModifiedBy == nil {
		T.field_ModifiedBy = T.Values["ModifiedBy"].(*elorm.FieldValueRef)
	}
	r, err := T.field_ModifiedBy.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *OrderLine) SetModifiedBy(newValue *User) {
	if T.field_ModifiedBy == nil {
		T.field_ModifiedBy = T.Values["ModifiedBy"].(*elorm.FieldValueRef)
	}
	err := T.field_ModifiedBy.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *OrderLine) ModifiedAt() time.Time {
	if T.field_ModifiedAt == nil {
		T.field_ModifiedAt = T.Values["ModifiedAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_ModifiedAt.Get()
}

func (T *OrderLine) SetModifiedAt(newValue time.Time) {
	if T.field_ModifiedAt == nil {
		T.field_ModifiedAt = T.Values["ModifiedAt"].(*elorm.FieldValueDateTime)
	}
	T.field_ModifiedAt.Set(newValue)
}

func (T *OrderLine) DeletedBy() *User {
	if T.field_DeletedBy == nil {
		T.field_DeletedBy = T.Values["DeletedBy"].(*elorm.FieldValueRef)
	}
	r, err := T.field_DeletedBy.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *OrderLine) SetDeletedBy(newValue *User) {
	if T.field_DeletedBy == nil {
		T.field_DeletedBy = T.Values["DeletedBy"].(*elorm.FieldValueRef)
	}
	err := T.field_DeletedBy.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *OrderLine) DeletedAt() time.Time {
	if T.field_DeletedAt == nil {
		T.field_DeletedAt = T.Values["DeletedAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_DeletedAt.Get()
}

func (T *OrderLine) SetDeletedAt(newValue time.Time) {
	if T.field_DeletedAt == nil {
		T.field_DeletedAt = T.Values["DeletedAt"].(*elorm.FieldValueDateTime)
	}
	T.field_DeletedAt.Set(newValue)
}

// Shop class
//////

type ShopDefStruct struct {
	*elorm.EntityDef
	Ref         *elorm.FieldDef
	IsDeleted   *elorm.FieldDef
	DataVersion *elorm.FieldDef

	Caption *elorm.FieldDef

	Description *elorm.FieldDef

	DeliveryConditions *elorm.FieldDef

	DiscountPercent *elorm.FieldDef

	CreatedBy *elorm.FieldDef

	CreatedAt *elorm.FieldDef

	ModifiedBy *elorm.FieldDef

	ModifiedAt *elorm.FieldDef

	DeletedBy *elorm.FieldDef

	DeletedAt *elorm.FieldDef
}

func (T *ShopDefStruct) SelectEntities(filters []*elorm.Filter, sorts []*elorm.SortItem, pageNo int, pageSize int) (result []*Shop, pages int, err error) {

	res, total, err := T.EntityDef.SelectEntities(filters, sorts, pageNo, pageSize)
	if err != nil {
		return nil, 0, err
	}

	res2 := make([]*Shop, 0, len(res))

	for _, r := range res {
		if r == nil {
			continue
		}
		rt := T.Wrap(r)
		res2 = append(res2, rt.(*Shop))
	}

	return res2, total, nil

}

type Shop struct {
	*elorm.Entity

	field_Caption            *elorm.FieldValueString
	field_Description        *elorm.FieldValueString
	field_DeliveryConditions *elorm.FieldValueString
	field_DiscountPercent    *elorm.FieldValueNumeric
	field_CreatedBy          *elorm.FieldValueRef
	field_CreatedAt          *elorm.FieldValueDateTime
	field_ModifiedBy         *elorm.FieldValueRef
	field_ModifiedAt         *elorm.FieldValueDateTime
	field_DeletedBy          *elorm.FieldValueRef
	field_DeletedAt          *elorm.FieldValueDateTime
}

func (T *Shop) Caption() string {
	if T.field_Caption == nil {
		T.field_Caption = T.Values["Caption"].(*elorm.FieldValueString)
	}
	return T.field_Caption.Get()
}

func (T *Shop) SetCaption(newValue string) {
	if T.field_Caption == nil {
		T.field_Caption = T.Values["Caption"].(*elorm.FieldValueString)
	}
	T.field_Caption.Set(newValue)
}

func (T *Shop) Description() string {
	if T.field_Description == nil {
		T.field_Description = T.Values["Description"].(*elorm.FieldValueString)
	}
	return T.field_Description.Get()
}

func (T *Shop) SetDescription(newValue string) {
	if T.field_Description == nil {
		T.field_Description = T.Values["Description"].(*elorm.FieldValueString)
	}
	T.field_Description.Set(newValue)
}

func (T *Shop) DeliveryConditions() string {
	if T.field_DeliveryConditions == nil {
		T.field_DeliveryConditions = T.Values["DeliveryConditions"].(*elorm.FieldValueString)
	}
	return T.field_DeliveryConditions.Get()
}

func (T *Shop) SetDeliveryConditions(newValue string) {
	if T.field_DeliveryConditions == nil {
		T.field_DeliveryConditions = T.Values["DeliveryConditions"].(*elorm.FieldValueString)
	}
	T.field_DeliveryConditions.Set(newValue)
}

func (T *Shop) DiscountPercent() float64 {
	if T.field_DiscountPercent == nil {
		T.field_DiscountPercent = T.Values["DiscountPercent"].(*elorm.FieldValueNumeric)
	}
	return T.field_DiscountPercent.Get()
}

func (T *Shop) SetDiscountPercent(newValue float64) {
	if T.field_DiscountPercent == nil {
		T.field_DiscountPercent = T.Values["DiscountPercent"].(*elorm.FieldValueNumeric)
	}
	T.field_DiscountPercent.Set(newValue)
}

func (T *Shop) CreatedBy() *User {
	if T.field_CreatedBy == nil {
		T.field_CreatedBy = T.Values["CreatedBy"].(*elorm.FieldValueRef)
	}
	r, err := T.field_CreatedBy.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *Shop) SetCreatedBy(newValue *User) {
	if T.field_CreatedBy == nil {
		T.field_CreatedBy = T.Values["CreatedBy"].(*elorm.FieldValueRef)
	}
	err := T.field_CreatedBy.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *Shop) CreatedAt() time.Time {
	if T.field_CreatedAt == nil {
		T.field_CreatedAt = T.Values["CreatedAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_CreatedAt.Get()
}

func (T *Shop) SetCreatedAt(newValue time.Time) {
	if T.field_CreatedAt == nil {
		T.field_CreatedAt = T.Values["CreatedAt"].(*elorm.FieldValueDateTime)
	}
	T.field_CreatedAt.Set(newValue)
}

func (T *Shop) ModifiedBy() *User {
	if T.field_ModifiedBy == nil {
		T.field_ModifiedBy = T.Values["ModifiedBy"].(*elorm.FieldValueRef)
	}
	r, err := T.field_ModifiedBy.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *Shop) SetModifiedBy(newValue *User) {
	if T.field_ModifiedBy == nil {
		T.field_ModifiedBy = T.Values["ModifiedBy"].(*elorm.FieldValueRef)
	}
	err := T.field_ModifiedBy.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *Shop) ModifiedAt() time.Time {
	if T.field_ModifiedAt == nil {
		T.field_ModifiedAt = T.Values["ModifiedAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_ModifiedAt.Get()
}

func (T *Shop) SetModifiedAt(newValue time.Time) {
	if T.field_ModifiedAt == nil {
		T.field_ModifiedAt = T.Values["ModifiedAt"].(*elorm.FieldValueDateTime)
	}
	T.field_ModifiedAt.Set(newValue)
}

func (T *Shop) DeletedBy() *User {
	if T.field_DeletedBy == nil {
		T.field_DeletedBy = T.Values["DeletedBy"].(*elorm.FieldValueRef)
	}
	r, err := T.field_DeletedBy.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *Shop) SetDeletedBy(newValue *User) {
	if T.field_DeletedBy == nil {
		T.field_DeletedBy = T.Values["DeletedBy"].(*elorm.FieldValueRef)
	}
	err := T.field_DeletedBy.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *Shop) DeletedAt() time.Time {
	if T.field_DeletedAt == nil {
		T.field_DeletedAt = T.Values["DeletedAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_DeletedAt.Get()
}

func (T *Shop) SetDeletedAt(newValue time.Time) {
	if T.field_DeletedAt == nil {
		T.field_DeletedAt = T.Values["DeletedAt"].(*elorm.FieldValueDateTime)
	}
	T.field_DeletedAt.Set(newValue)
}

// Tag class
//////

type TagDefStruct struct {
	*elorm.EntityDef
	Ref         *elorm.FieldDef
	IsDeleted   *elorm.FieldDef
	DataVersion *elorm.FieldDef

	Name *elorm.FieldDef

	Color *elorm.FieldDef
}

func (T *TagDefStruct) SelectEntities(filters []*elorm.Filter, sorts []*elorm.SortItem, pageNo int, pageSize int) (result []*Tag, pages int, err error) {

	res, total, err := T.EntityDef.SelectEntities(filters, sorts, pageNo, pageSize)
	if err != nil {
		return nil, 0, err
	}

	res2 := make([]*Tag, 0, len(res))

	for _, r := range res {
		if r == nil {
			continue
		}
		rt := T.Wrap(r)
		res2 = append(res2, rt.(*Tag))
	}

	return res2, total, nil

}

type Tag struct {
	*elorm.Entity

	field_Name  *elorm.FieldValueString
	field_Color *elorm.FieldValueString
}

func (T *Tag) Name() string {
	if T.field_Name == nil {
		T.field_Name = T.Values["Name"].(*elorm.FieldValueString)
	}
	return T.field_Name.Get()
}

func (T *Tag) SetName(newValue string) {
	if T.field_Name == nil {
		T.field_Name = T.Values["Name"].(*elorm.FieldValueString)
	}
	T.field_Name.Set(newValue)
}

func (T *Tag) Color() string {
	if T.field_Color == nil {
		T.field_Color = T.Values["Color"].(*elorm.FieldValueString)
	}
	return T.field_Color.Get()
}

func (T *Tag) SetColor(newValue string) {
	if T.field_Color == nil {
		T.field_Color = T.Values["Color"].(*elorm.FieldValueString)
	}
	T.field_Color.Set(newValue)
}

// Token class
//////

type TokenDefStruct struct {
	*elorm.EntityDef
	Ref         *elorm.FieldDef
	IsDeleted   *elorm.FieldDef
	DataVersion *elorm.FieldDef

	User *elorm.FieldDef

	AccessToken *elorm.FieldDef

	AccessTokenExpiresAt *elorm.FieldDef

	RefreshToken *elorm.FieldDef

	RefreshTokenExpiresAt *elorm.FieldDef
}

func (T *TokenDefStruct) SelectEntities(filters []*elorm.Filter, sorts []*elorm.SortItem, pageNo int, pageSize int) (result []*Token, pages int, err error) {

	res, total, err := T.EntityDef.SelectEntities(filters, sorts, pageNo, pageSize)
	if err != nil {
		return nil, 0, err
	}

	res2 := make([]*Token, 0, len(res))

	for _, r := range res {
		if r == nil {
			continue
		}
		rt := T.Wrap(r)
		res2 = append(res2, rt.(*Token))
	}

	return res2, total, nil

}

type Token struct {
	*elorm.Entity

	field_User                  *elorm.FieldValueRef
	field_AccessToken           *elorm.FieldValueString
	field_AccessTokenExpiresAt  *elorm.FieldValueDateTime
	field_RefreshToken          *elorm.FieldValueString
	field_RefreshTokenExpiresAt *elorm.FieldValueDateTime
}

func (T *Token) User() *User {
	if T.field_User == nil {
		T.field_User = T.Values["User"].(*elorm.FieldValueRef)
	}
	r, err := T.field_User.Get()
	if err != nil {
		panic(err)
	}
	return r.(*User)
}

func (T *Token) SetUser(newValue *User) {
	if T.field_User == nil {
		T.field_User = T.Values["User"].(*elorm.FieldValueRef)
	}
	err := T.field_User.Set(newValue)
	if err != nil {
		panic(err)
	}
}

func (T *Token) AccessToken() string {
	if T.field_AccessToken == nil {
		T.field_AccessToken = T.Values["AccessToken"].(*elorm.FieldValueString)
	}
	return T.field_AccessToken.Get()
}

func (T *Token) SetAccessToken(newValue string) {
	if T.field_AccessToken == nil {
		T.field_AccessToken = T.Values["AccessToken"].(*elorm.FieldValueString)
	}
	T.field_AccessToken.Set(newValue)
}

func (T *Token) AccessTokenExpiresAt() time.Time {
	if T.field_AccessTokenExpiresAt == nil {
		T.field_AccessTokenExpiresAt = T.Values["AccessTokenExpiresAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_AccessTokenExpiresAt.Get()
}

func (T *Token) SetAccessTokenExpiresAt(newValue time.Time) {
	if T.field_AccessTokenExpiresAt == nil {
		T.field_AccessTokenExpiresAt = T.Values["AccessTokenExpiresAt"].(*elorm.FieldValueDateTime)
	}
	T.field_AccessTokenExpiresAt.Set(newValue)
}

func (T *Token) RefreshToken() string {
	if T.field_RefreshToken == nil {
		T.field_RefreshToken = T.Values["RefreshToken"].(*elorm.FieldValueString)
	}
	return T.field_RefreshToken.Get()
}

func (T *Token) SetRefreshToken(newValue string) {
	if T.field_RefreshToken == nil {
		T.field_RefreshToken = T.Values["RefreshToken"].(*elorm.FieldValueString)
	}
	T.field_RefreshToken.Set(newValue)
}

func (T *Token) RefreshTokenExpiresAt() time.Time {
	if T.field_RefreshTokenExpiresAt == nil {
		T.field_RefreshTokenExpiresAt = T.Values["RefreshTokenExpiresAt"].(*elorm.FieldValueDateTime)
	}
	return T.field_RefreshTokenExpiresAt.Get()
}

func (T *Token) SetRefreshTokenExpiresAt(newValue time.Time) {
	if T.field_RefreshTokenExpiresAt == nil {
		T.field_RefreshTokenExpiresAt = T.Values["RefreshTokenExpiresAt"].(*elorm.FieldValueDateTime)
	}
	T.field_RefreshTokenExpiresAt.Set(newValue)
}

// User class
//////

type UserDefStruct struct {
	*elorm.EntityDef
	Ref         *elorm.FieldDef
	IsDeleted   *elorm.FieldDef
	DataVersion *elorm.FieldDef

	Username *elorm.FieldDef

	Email *elorm.FieldDef

	PasswordHash *elorm.FieldDef

	IsActive *elorm.FieldDef

	ShopManager *elorm.FieldDef

	Admin *elorm.FieldDef

	TelegramUsername *elorm.FieldDef

	TelegramCheckCode *elorm.FieldDef

	TelegramVerified *elorm.FieldDef

	TelegramChatId *elorm.FieldDef

	Description *elorm.FieldDef

	RecoverCodeEmail *elorm.FieldDef

	RecoverCodeDeadline *elorm.FieldDef
}

func (T *UserDefStruct) SelectEntities(filters []*elorm.Filter, sorts []*elorm.SortItem, pageNo int, pageSize int) (result []*User, pages int, err error) {

	res, total, err := T.EntityDef.SelectEntities(filters, sorts, pageNo, pageSize)
	if err != nil {
		return nil, 0, err
	}

	res2 := make([]*User, 0, len(res))

	for _, r := range res {
		if r == nil {
			continue
		}
		rt := T.Wrap(r)
		res2 = append(res2, rt.(*User))
	}

	return res2, total, nil

}

type User struct {
	*elorm.Entity

	field_Username            *elorm.FieldValueString
	field_Email               *elorm.FieldValueString
	field_PasswordHash        *elorm.FieldValueString
	field_IsActive            *elorm.FieldValueBool
	field_ShopManager         *elorm.FieldValueBool
	field_Admin               *elorm.FieldValueBool
	field_TelegramUsername    *elorm.FieldValueString
	field_TelegramCheckCode   *elorm.FieldValueString
	field_TelegramVerified    *elorm.FieldValueBool
	field_TelegramChatId      *elorm.FieldValueInt
	field_Description         *elorm.FieldValueString
	field_RecoverCodeEmail    *elorm.FieldValueString
	field_RecoverCodeDeadline *elorm.FieldValueDateTime
}

func (T *User) Username() string {
	if T.field_Username == nil {
		T.field_Username = T.Values["Username"].(*elorm.FieldValueString)
	}
	return T.field_Username.Get()
}

func (T *User) SetUsername(newValue string) {
	if T.field_Username == nil {
		T.field_Username = T.Values["Username"].(*elorm.FieldValueString)
	}
	T.field_Username.Set(newValue)
}

func (T *User) Email() string {
	if T.field_Email == nil {
		T.field_Email = T.Values["Email"].(*elorm.FieldValueString)
	}
	return T.field_Email.Get()
}

func (T *User) SetEmail(newValue string) {
	if T.field_Email == nil {
		T.field_Email = T.Values["Email"].(*elorm.FieldValueString)
	}
	T.field_Email.Set(newValue)
}

func (T *User) PasswordHash() string {
	if T.field_PasswordHash == nil {
		T.field_PasswordHash = T.Values["PasswordHash"].(*elorm.FieldValueString)
	}
	return T.field_PasswordHash.Get()
}

func (T *User) SetPasswordHash(newValue string) {
	if T.field_PasswordHash == nil {
		T.field_PasswordHash = T.Values["PasswordHash"].(*elorm.FieldValueString)
	}
	T.field_PasswordHash.Set(newValue)
}

func (T *User) IsActive() bool {
	if T.field_IsActive == nil {
		T.field_IsActive = T.Values["IsActive"].(*elorm.FieldValueBool)
	}
	return T.field_IsActive.Get()
}

func (T *User) SetIsActive(newValue bool) {
	if T.field_IsActive == nil {
		T.field_IsActive = T.Values["IsActive"].(*elorm.FieldValueBool)
	}
	T.field_IsActive.Set(newValue)
}

func (T *User) ShopManager() bool {
	if T.field_ShopManager == nil {
		T.field_ShopManager = T.Values["ShopManager"].(*elorm.FieldValueBool)
	}
	return T.field_ShopManager.Get()
}

func (T *User) SetShopManager(newValue bool) {
	if T.field_ShopManager == nil {
		T.field_ShopManager = T.Values["ShopManager"].(*elorm.FieldValueBool)
	}
	T.field_ShopManager.Set(newValue)
}

func (T *User) Admin() bool {
	if T.field_Admin == nil {
		T.field_Admin = T.Values["Admin"].(*elorm.FieldValueBool)
	}
	return T.field_Admin.Get()
}

func (T *User) SetAdmin(newValue bool) {
	if T.field_Admin == nil {
		T.field_Admin = T.Values["Admin"].(*elorm.FieldValueBool)
	}
	T.field_Admin.Set(newValue)
}

func (T *User) TelegramUsername() string {
	if T.field_TelegramUsername == nil {
		T.field_TelegramUsername = T.Values["TelegramUsername"].(*elorm.FieldValueString)
	}
	return T.field_TelegramUsername.Get()
}

func (T *User) SetTelegramUsername(newValue string) {
	if T.field_TelegramUsername == nil {
		T.field_TelegramUsername = T.Values["TelegramUsername"].(*elorm.FieldValueString)
	}
	T.field_TelegramUsername.Set(newValue)
}

func (T *User) TelegramCheckCode() string {
	if T.field_TelegramCheckCode == nil {
		T.field_TelegramCheckCode = T.Values["TelegramCheckCode"].(*elorm.FieldValueString)
	}
	return T.field_TelegramCheckCode.Get()
}

func (T *User) SetTelegramCheckCode(newValue string) {
	if T.field_TelegramCheckCode == nil {
		T.field_TelegramCheckCode = T.Values["TelegramCheckCode"].(*elorm.FieldValueString)
	}
	T.field_TelegramCheckCode.Set(newValue)
}

func (T *User) TelegramVerified() bool {
	if T.field_TelegramVerified == nil {
		T.field_TelegramVerified = T.Values["TelegramVerified"].(*elorm.FieldValueBool)
	}
	return T.field_TelegramVerified.Get()
}

func (T *User) SetTelegramVerified(newValue bool) {
	if T.field_TelegramVerified == nil {
		T.field_TelegramVerified = T.Values["TelegramVerified"].(*elorm.FieldValueBool)
	}
	T.field_TelegramVerified.Set(newValue)
}

func (T *User) TelegramChatId() int64 {
	if T.field_TelegramChatId == nil {
		T.field_TelegramChatId = T.Values["TelegramChatId"].(*elorm.FieldValueInt)
	}
	return T.field_TelegramChatId.Get()
}

func (T *User) SetTelegramChatId(newValue int64) {
	if T.field_TelegramChatId == nil {
		T.field_TelegramChatId = T.Values["TelegramChatId"].(*elorm.FieldValueInt)
	}
	T.field_TelegramChatId.Set(newValue)
}

func (T *User) Description() string {
	if T.field_Description == nil {
		T.field_Description = T.Values["Description"].(*elorm.FieldValueString)
	}
	return T.field_Description.Get()
}

func (T *User) SetDescription(newValue string) {
	if T.field_Description == nil {
		T.field_Description = T.Values["Description"].(*elorm.FieldValueString)
	}
	T.field_Description.Set(newValue)
}

func (T *User) RecoverCodeEmail() string {
	if T.field_RecoverCodeEmail == nil {
		T.field_RecoverCodeEmail = T.Values["RecoverCodeEmail"].(*elorm.FieldValueString)
	}
	return T.field_RecoverCodeEmail.Get()
}

func (T *User) SetRecoverCodeEmail(newValue string) {
	if T.field_RecoverCodeEmail == nil {
		T.field_RecoverCodeEmail = T.Values["RecoverCodeEmail"].(*elorm.FieldValueString)
	}
	T.field_RecoverCodeEmail.Set(newValue)
}

func (T *User) RecoverCodeDeadline() time.Time {
	if T.field_RecoverCodeDeadline == nil {
		T.field_RecoverCodeDeadline = T.Values["RecoverCodeDeadline"].(*elorm.FieldValueDateTime)
	}
	return T.field_RecoverCodeDeadline.Get()
}

func (T *User) SetRecoverCodeDeadline(newValue time.Time) {
	if T.field_RecoverCodeDeadline == nil {
		T.field_RecoverCodeDeadline = T.Values["RecoverCodeDeadline"].(*elorm.FieldValueDateTime)
	}
	T.field_RecoverCodeDeadline.Set(newValue)
}

// BusinessObjects fragment

const BusinessObjectsFragment = "BusinessObjects"

type BusinessObjectsFragmentMethods interface {
	IsDeleted() bool
	SetIsDeleted(newValue bool)

	CreatedBy() *User
	SetCreatedBy(newValue *User)

	CreatedAt() time.Time
	SetCreatedAt(newValue time.Time)

	ModifiedBy() *User
	SetModifiedBy(newValue *User)

	ModifiedAt() time.Time
	SetModifiedAt(newValue time.Time)

	DeletedBy() *User
	SetDeletedBy(newValue *User)

	DeletedAt() time.Time
	SetDeletedAt(newValue time.Time)
}

// DbContext core
//////

type DbContext struct {
	*elorm.Factory
	CustomerOrderDef CustomerOrderDefStruct
	GoodDef          GoodDefStruct
	GoodTagDef       GoodTagDefStruct
	OrderLineDef     OrderLineDefStruct
	ShopDef          ShopDefStruct
	TagDef           TagDefStruct
	TokenDef         TokenDefStruct
	UserDef          UserDefStruct
}

func CreateDbContext(dbDialect string, connectionString string) (*DbContext, error) {

	var err error
	frg := []string{}
	_ = frg // to avoid unused variable error if no fragments are defined

	r := &DbContext{}
	r.Factory, err = elorm.CreateFactory(dbDialect, connectionString)
	if err != nil {
		return nil, err
	}

	r.CustomerOrderDef.EntityDef, err = r.CreateEntityDef("CustomerOrder", "CustomerOrders")
	if err != nil {
		return nil, err
	}

	r.CustomerOrderDef.UseSoftDelete = false

	r.CustomerOrderDef.Fragments = make([]string, 0)
	frg = []string{}

	frg = append(frg, "BusinessObjects")

	r.CustomerOrderDef.Fragments = frg

	r.GoodDef.EntityDef, err = r.CreateEntityDef("Good", "Goods")
	if err != nil {
		return nil, err
	}

	r.GoodDef.UseSoftDelete = true

	r.GoodDef.Fragments = make([]string, 0)
	frg = []string{}

	frg = append(frg, "BusinessObjects")

	r.GoodDef.Fragments = frg

	r.GoodTagDef.EntityDef, err = r.CreateEntityDef("GoodTag", "GoodTags")
	if err != nil {
		return nil, err
	}

	r.GoodTagDef.UseSoftDelete = false

	r.OrderLineDef.EntityDef, err = r.CreateEntityDef("OrderLine", "OrderLines")
	if err != nil {
		return nil, err
	}

	r.OrderLineDef.UseSoftDelete = false

	r.OrderLineDef.Fragments = make([]string, 0)
	frg = []string{}

	frg = append(frg, "BusinessObjects")

	r.OrderLineDef.Fragments = frg

	r.ShopDef.EntityDef, err = r.CreateEntityDef("Shop", "Shops")
	if err != nil {
		return nil, err
	}

	r.ShopDef.UseSoftDelete = true

	r.ShopDef.Fragments = make([]string, 0)
	frg = []string{}

	frg = append(frg, "BusinessObjects")

	r.ShopDef.Fragments = frg

	r.TagDef.EntityDef, err = r.CreateEntityDef("Tag", "Tags")
	if err != nil {
		return nil, err
	}

	r.TagDef.UseSoftDelete = false

	r.TokenDef.EntityDef, err = r.CreateEntityDef("Token", "Tokens")
	if err != nil {
		return nil, err
	}

	r.TokenDef.UseSoftDelete = false

	r.UserDef.EntityDef, err = r.CreateEntityDef("User", "Users")
	if err != nil {
		return nil, err
	}

	r.UserDef.UseSoftDelete = false

	// CustomerOrder
	//////

	r.CustomerOrderDef.Ref = r.CustomerOrderDef.FieldDefByName("Ref")
	r.CustomerOrderDef.IsDeleted = r.CustomerOrderDef.FieldDefByName("IsDeleted")
	r.CustomerOrderDef.DataVersion = r.CustomerOrderDef.FieldDefByName("DataVersion")

	r.CustomerOrderDef.Sender, _ = r.CustomerOrderDef.AddRefFieldDef("Sender", r.UserDef.EntityDef)
	r.CustomerOrderDef.Status, _ = r.CustomerOrderDef.AddIntFieldDef("Status")
	r.CustomerOrderDef.Qty, _ = r.CustomerOrderDef.AddNumericFieldDef("Qty", 15, 2)
	r.CustomerOrderDef.Sum, _ = r.CustomerOrderDef.AddNumericFieldDef("Sum", 15, 2)
	r.CustomerOrderDef.OrderDiscountSum, _ = r.CustomerOrderDef.AddNumericFieldDef("OrderDiscountSum", 15, 2)
	r.CustomerOrderDef.SenderComment, _ = r.CustomerOrderDef.AddStringFieldDef("SenderComment", 200)
	r.CustomerOrderDef.CustomerComment, _ = r.CustomerOrderDef.AddStringFieldDef("CustomerComment", 200)
	r.CustomerOrderDef.ExpectedDeliveryDate, _ = r.CustomerOrderDef.AddDateTimeFieldDef("ExpectedDeliveryDate")
	r.CustomerOrderDef.CreatedBy, _ = r.CustomerOrderDef.AddRefFieldDef("CreatedBy", r.UserDef.EntityDef)
	r.CustomerOrderDef.CreatedAt, _ = r.CustomerOrderDef.AddDateTimeFieldDef("CreatedAt")
	r.CustomerOrderDef.ModifiedBy, _ = r.CustomerOrderDef.AddRefFieldDef("ModifiedBy", r.UserDef.EntityDef)
	r.CustomerOrderDef.ModifiedAt, _ = r.CustomerOrderDef.AddDateTimeFieldDef("ModifiedAt")
	r.CustomerOrderDef.DeletedBy, _ = r.CustomerOrderDef.AddRefFieldDef("DeletedBy", r.UserDef.EntityDef)
	r.CustomerOrderDef.DeletedAt, _ = r.CustomerOrderDef.AddDateTimeFieldDef("DeletedAt")

	r.CustomerOrderDef.Wrap = func(source *elorm.Entity) any { return &CustomerOrder{Entity: source} }

	// Good
	//////

	r.GoodDef.Ref = r.GoodDef.FieldDefByName("Ref")
	r.GoodDef.IsDeleted = r.GoodDef.FieldDefByName("IsDeleted")
	r.GoodDef.DataVersion = r.GoodDef.FieldDefByName("DataVersion")

	r.GoodDef.OwnerShop, _ = r.GoodDef.AddRefFieldDef("OwnerShop", r.ShopDef.EntityDef)
	r.GoodDef.Caption, _ = r.GoodDef.AddStringFieldDef("Caption", 100)
	r.GoodDef.Article, _ = r.GoodDef.AddStringFieldDef("Article", 50)
	r.GoodDef.Url, _ = r.GoodDef.AddStringFieldDef("Url", 500)
	r.GoodDef.Description, _ = r.GoodDef.AddStringFieldDef("Description", 4096)
	r.GoodDef.Price, _ = r.GoodDef.AddNumericFieldDef("Price", 10, 2)
	r.GoodDef.OrderInShop, _ = r.GoodDef.AddIntFieldDef("OrderInShop")
	r.GoodDef.CreatedBy, _ = r.GoodDef.AddRefFieldDef("CreatedBy", r.UserDef.EntityDef)
	r.GoodDef.CreatedAt, _ = r.GoodDef.AddDateTimeFieldDef("CreatedAt")
	r.GoodDef.ModifiedBy, _ = r.GoodDef.AddRefFieldDef("ModifiedBy", r.UserDef.EntityDef)
	r.GoodDef.ModifiedAt, _ = r.GoodDef.AddDateTimeFieldDef("ModifiedAt")
	r.GoodDef.DeletedBy, _ = r.GoodDef.AddRefFieldDef("DeletedBy", r.UserDef.EntityDef)
	r.GoodDef.DeletedAt, _ = r.GoodDef.AddDateTimeFieldDef("DeletedAt")

	r.GoodDef.Wrap = func(source *elorm.Entity) any { return &Good{Entity: source} }

	err = r.GoodDef.AddIndex(false,
		r.GoodDef.OwnerShop,
	)
	if err != nil {
		return nil, err
	}

	// GoodTag
	//////

	r.GoodTagDef.Ref = r.GoodTagDef.FieldDefByName("Ref")
	r.GoodTagDef.IsDeleted = r.GoodTagDef.FieldDefByName("IsDeleted")
	r.GoodTagDef.DataVersion = r.GoodTagDef.FieldDefByName("DataVersion")

	r.GoodTagDef.Good, _ = r.GoodTagDef.AddRefFieldDef("Good", r.GoodDef.EntityDef)
	r.GoodTagDef.Tag, _ = r.GoodTagDef.AddRefFieldDef("Tag", r.TagDef.EntityDef)

	r.GoodTagDef.Wrap = func(source *elorm.Entity) any { return &GoodTag{Entity: source} }

	// OrderLine
	//////

	r.OrderLineDef.Ref = r.OrderLineDef.FieldDefByName("Ref")
	r.OrderLineDef.IsDeleted = r.OrderLineDef.FieldDefByName("IsDeleted")
	r.OrderLineDef.DataVersion = r.OrderLineDef.FieldDefByName("DataVersion")

	r.OrderLineDef.Shop, _ = r.OrderLineDef.AddRefFieldDef("Shop", r.ShopDef.EntityDef)
	r.OrderLineDef.CustomerOrder, _ = r.OrderLineDef.AddRefFieldDef("CustomerOrder", r.CustomerOrderDef.EntityDef)
	r.OrderLineDef.Good, _ = r.OrderLineDef.AddRefFieldDef("Good", r.GoodDef.EntityDef)
	r.OrderLineDef.Qty, _ = r.OrderLineDef.AddNumericFieldDef("Qty", 15, 2)
	r.OrderLineDef.LineDiscountSum, _ = r.OrderLineDef.AddNumericFieldDef("LineDiscountSum", 15, 2)
	r.OrderLineDef.Sum, _ = r.OrderLineDef.AddNumericFieldDef("Sum", 15, 2)
	r.OrderLineDef.CreatedBy, _ = r.OrderLineDef.AddRefFieldDef("CreatedBy", r.UserDef.EntityDef)
	r.OrderLineDef.CreatedAt, _ = r.OrderLineDef.AddDateTimeFieldDef("CreatedAt")
	r.OrderLineDef.ModifiedBy, _ = r.OrderLineDef.AddRefFieldDef("ModifiedBy", r.UserDef.EntityDef)
	r.OrderLineDef.ModifiedAt, _ = r.OrderLineDef.AddDateTimeFieldDef("ModifiedAt")
	r.OrderLineDef.DeletedBy, _ = r.OrderLineDef.AddRefFieldDef("DeletedBy", r.UserDef.EntityDef)
	r.OrderLineDef.DeletedAt, _ = r.OrderLineDef.AddDateTimeFieldDef("DeletedAt")

	r.OrderLineDef.Wrap = func(source *elorm.Entity) any { return &OrderLine{Entity: source} }

	// Shop
	//////

	r.ShopDef.Ref = r.ShopDef.FieldDefByName("Ref")
	r.ShopDef.IsDeleted = r.ShopDef.FieldDefByName("IsDeleted")
	r.ShopDef.DataVersion = r.ShopDef.FieldDefByName("DataVersion")

	r.ShopDef.Caption, _ = r.ShopDef.AddStringFieldDef("Caption", 100)
	r.ShopDef.Description, _ = r.ShopDef.AddStringFieldDef("Description", 300)
	r.ShopDef.DeliveryConditions, _ = r.ShopDef.AddStringFieldDef("DeliveryConditions", 300)
	r.ShopDef.DiscountPercent, _ = r.ShopDef.AddNumericFieldDef("DiscountPercent", 4, 1)
	r.ShopDef.CreatedBy, _ = r.ShopDef.AddRefFieldDef("CreatedBy", r.UserDef.EntityDef)
	r.ShopDef.CreatedAt, _ = r.ShopDef.AddDateTimeFieldDef("CreatedAt")
	r.ShopDef.ModifiedBy, _ = r.ShopDef.AddRefFieldDef("ModifiedBy", r.UserDef.EntityDef)
	r.ShopDef.ModifiedAt, _ = r.ShopDef.AddDateTimeFieldDef("ModifiedAt")
	r.ShopDef.DeletedBy, _ = r.ShopDef.AddRefFieldDef("DeletedBy", r.UserDef.EntityDef)
	r.ShopDef.DeletedAt, _ = r.ShopDef.AddDateTimeFieldDef("DeletedAt")

	r.ShopDef.Wrap = func(source *elorm.Entity) any { return &Shop{Entity: source} }

	// Tag
	//////

	r.TagDef.Ref = r.TagDef.FieldDefByName("Ref")
	r.TagDef.IsDeleted = r.TagDef.FieldDefByName("IsDeleted")
	r.TagDef.DataVersion = r.TagDef.FieldDefByName("DataVersion")

	r.TagDef.Name, _ = r.TagDef.AddStringFieldDef("Name", 20)
	r.TagDef.Color, _ = r.TagDef.AddStringFieldDef("Color", 20)

	r.TagDef.Wrap = func(source *elorm.Entity) any { return &Tag{Entity: source} }

	err = r.TagDef.AddIndex(true,
		r.TagDef.Name,
	)
	if err != nil {
		return nil, err
	}

	// Token
	//////

	r.TokenDef.Ref = r.TokenDef.FieldDefByName("Ref")
	r.TokenDef.IsDeleted = r.TokenDef.FieldDefByName("IsDeleted")
	r.TokenDef.DataVersion = r.TokenDef.FieldDefByName("DataVersion")

	r.TokenDef.User, _ = r.TokenDef.AddRefFieldDef("User", r.UserDef.EntityDef)
	r.TokenDef.AccessToken, _ = r.TokenDef.AddStringFieldDef("AccessToken", 50)
	r.TokenDef.AccessTokenExpiresAt, _ = r.TokenDef.AddDateTimeFieldDef("AccessTokenExpiresAt")
	r.TokenDef.RefreshToken, _ = r.TokenDef.AddStringFieldDef("RefreshToken", 50)
	r.TokenDef.RefreshTokenExpiresAt, _ = r.TokenDef.AddDateTimeFieldDef("RefreshTokenExpiresAt")

	r.TokenDef.Wrap = func(source *elorm.Entity) any { return &Token{Entity: source} }

	err = r.TokenDef.AddIndex(false,
		r.TokenDef.User,
	)
	if err != nil {
		return nil, err
	}

	err = r.TokenDef.AddIndex(true,
		r.TokenDef.AccessToken, r.TokenDef.AccessTokenExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	err = r.TokenDef.AddIndex(true,
		r.TokenDef.RefreshToken, r.TokenDef.RefreshTokenExpiresAt,
	)
	if err != nil {
		return nil, err
	}

	// User
	//////

	r.UserDef.Ref = r.UserDef.FieldDefByName("Ref")
	r.UserDef.IsDeleted = r.UserDef.FieldDefByName("IsDeleted")
	r.UserDef.DataVersion = r.UserDef.FieldDefByName("DataVersion")

	r.UserDef.Username, _ = r.UserDef.AddStringFieldDef("Username", 100)
	r.UserDef.Email, _ = r.UserDef.AddStringFieldDef("Email", 100)
	r.UserDef.PasswordHash, _ = r.UserDef.AddStringFieldDef("PasswordHash", 200)
	r.UserDef.IsActive, _ = r.UserDef.AddBoolFieldDef("IsActive")
	r.UserDef.ShopManager, _ = r.UserDef.AddBoolFieldDef("ShopManager")
	r.UserDef.Admin, _ = r.UserDef.AddBoolFieldDef("Admin")
	r.UserDef.TelegramUsername, _ = r.UserDef.AddStringFieldDef("TelegramUsername", 100)
	r.UserDef.TelegramCheckCode, _ = r.UserDef.AddStringFieldDef("TelegramCheckCode", 100)
	r.UserDef.TelegramVerified, _ = r.UserDef.AddBoolFieldDef("TelegramVerified")
	r.UserDef.TelegramChatId, _ = r.UserDef.AddIntFieldDef("TelegramChatId")
	r.UserDef.Description, _ = r.UserDef.AddStringFieldDef("Description", 300)
	r.UserDef.RecoverCodeEmail, _ = r.UserDef.AddStringFieldDef("RecoverCodeEmail", 30)
	r.UserDef.RecoverCodeDeadline, _ = r.UserDef.AddDateTimeFieldDef("RecoverCodeDeadline")

	r.UserDef.Wrap = func(source *elorm.Entity) any { return &User{Entity: source} }

	err = r.UserDef.AddIndex(true,
		r.UserDef.Username,
	)
	if err != nil {
		return nil, err
	}

	err = r.UserDef.AddIndex(true,
		r.UserDef.Email,
	)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func (T *DbContext) CreateCustomerOrder() (*CustomerOrder, error) {
	r, err := T.CreateEntityWrapped(T.CustomerOrderDef.EntityDef)
	if err != nil {
		return nil, err
	}
	rt := r.(*CustomerOrder)
	return rt, nil
}

func (T *DbContext) LoadCustomerOrder(Ref string) (*CustomerOrder, error) {
	r, err := T.LoadEntityWrapped(Ref)
	if err != nil {
		return nil, err
	}
	rt := r.(*CustomerOrder)
	return rt, nil
}

func (T *DbContext) CreateGood() (*Good, error) {
	r, err := T.CreateEntityWrapped(T.GoodDef.EntityDef)
	if err != nil {
		return nil, err
	}
	rt := r.(*Good)
	return rt, nil
}

func (T *DbContext) LoadGood(Ref string) (*Good, error) {
	r, err := T.LoadEntityWrapped(Ref)
	if err != nil {
		return nil, err
	}
	rt := r.(*Good)
	return rt, nil
}

func (T *DbContext) CreateGoodTag() (*GoodTag, error) {
	r, err := T.CreateEntityWrapped(T.GoodTagDef.EntityDef)
	if err != nil {
		return nil, err
	}
	rt := r.(*GoodTag)
	return rt, nil
}

func (T *DbContext) LoadGoodTag(Ref string) (*GoodTag, error) {
	r, err := T.LoadEntityWrapped(Ref)
	if err != nil {
		return nil, err
	}
	rt := r.(*GoodTag)
	return rt, nil
}

func (T *DbContext) CreateOrderLine() (*OrderLine, error) {
	r, err := T.CreateEntityWrapped(T.OrderLineDef.EntityDef)
	if err != nil {
		return nil, err
	}
	rt := r.(*OrderLine)
	return rt, nil
}

func (T *DbContext) LoadOrderLine(Ref string) (*OrderLine, error) {
	r, err := T.LoadEntityWrapped(Ref)
	if err != nil {
		return nil, err
	}
	rt := r.(*OrderLine)
	return rt, nil
}

func (T *DbContext) CreateShop() (*Shop, error) {
	r, err := T.CreateEntityWrapped(T.ShopDef.EntityDef)
	if err != nil {
		return nil, err
	}
	rt := r.(*Shop)
	return rt, nil
}

func (T *DbContext) LoadShop(Ref string) (*Shop, error) {
	r, err := T.LoadEntityWrapped(Ref)
	if err != nil {
		return nil, err
	}
	rt := r.(*Shop)
	return rt, nil
}

func (T *DbContext) CreateTag() (*Tag, error) {
	r, err := T.CreateEntityWrapped(T.TagDef.EntityDef)
	if err != nil {
		return nil, err
	}
	rt := r.(*Tag)
	return rt, nil
}

func (T *DbContext) LoadTag(Ref string) (*Tag, error) {
	r, err := T.LoadEntityWrapped(Ref)
	if err != nil {
		return nil, err
	}
	rt := r.(*Tag)
	return rt, nil
}

func (T *DbContext) CreateToken() (*Token, error) {
	r, err := T.CreateEntityWrapped(T.TokenDef.EntityDef)
	if err != nil {
		return nil, err
	}
	rt := r.(*Token)
	return rt, nil
}

func (T *DbContext) LoadToken(Ref string) (*Token, error) {
	r, err := T.LoadEntityWrapped(Ref)
	if err != nil {
		return nil, err
	}
	rt := r.(*Token)
	return rt, nil
}

func (T *DbContext) CreateUser() (*User, error) {
	r, err := T.CreateEntityWrapped(T.UserDef.EntityDef)
	if err != nil {
		return nil, err
	}
	rt := r.(*User)
	return rt, nil
}

func (T *DbContext) LoadUser(Ref string) (*User, error) {
	r, err := T.LoadEntityWrapped(Ref)
	if err != nil {
		return nil, err
	}
	rt := r.(*User)
	return rt, nil
}
