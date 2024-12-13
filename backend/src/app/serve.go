package app

import (
	"net/http"

	"github.com/cmo7/folly4/src/app/models"
	"github.com/cmo7/folly4/src/app/services"
	"github.com/cmo7/folly4/src/data/database"
	"github.com/cmo7/folly4/src/lib/generics"
	"github.com/cmo7/folly4/src/lib/generics/controller"
	"github.com/cmo7/folly4/src/lib/generics/router"
	"gorm.io/gorm"
)

func Serve() {
	db, err := database.ConnectWithConfig(
		&gorm.Config{},
		&database.ConnectionData{
			File: ":memory:",
		},
	)
	if err != nil {
		panic(err)
	}

	// Migrate the schema
	db.AutoMigrate(
		&models.RoleEntity{},
		&models.PermissionEntity{},
		&models.UserEntity{},
	)

	userController := controller.NewController(
		services.GetUserService(db),
		generics.NewGenericMapperExcluding[*models.UserEntity, *models.UserEntity]([]string{"Password"}),
	)

	userRouter := router.NewRouter(userController)

	router := http.NewServeMux()
	router.Handle(userRouter.GetBaseRoute(), userRouter)
	http.ListenAndServe(":8080", router)
}
