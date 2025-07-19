package models

func (dbc *DbContext) SetHandlers() error {

	err := dbc.AddFillNewHandler(dbc.UserDef.EntityDef, func(entity any) error {
		user := entity.(*User)
		user.SetIsActive(true)
		return nil
	})
	if err != nil {
		return err
	}

	return nil

}
