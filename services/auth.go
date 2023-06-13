package services

import "github.com/OhMinsSup/notes-server-go/models"

func (service *Service) GetExistsUserByUsername(username string) bool {
	db := service.store.DataDB()
	has, err := db.Exist(&models.User{
		Username: username,
	})
	if err != nil {
		return false
	}
	return has
}

func (service *Service) GetUserByUsername(username string) *models.User {
	db := service.store.DataDB()
	var user = models.User{
		Username: username,
	}
	has, err := db.Get(&user)
	if err != nil {
		return nil
	}
	if !has {
		return nil
	}
	return &user
}

func (service *Service) GetUserByEmail(email string) *models.User {
	db := service.store.DataDB()
	var user = models.User{
		Email: email,
	}
	has, err := db.Get(&user)
	if err != nil {
		return nil
	}
	if !has {
		return nil
	}
	return &user
}
