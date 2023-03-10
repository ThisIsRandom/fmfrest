package main

import (
	"github.com/gofiber/fiber/v2"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/thisisrandom/fmfrest/controllers"
	"github.com/thisisrandom/fmfrest/database"
	"github.com/thisisrandom/fmfrest/internal"
	"fmt"
	"os"
)

var (
	fakeItTillUMakeIt = "root:root@tcp(127.0.0.1:3306)/fmf?charset=utf8mb4&parseTime=True&loc=Local"
)

func main() {
	app := fiber.New()

	api := app.Group("api")

	imageStore, err := internal.NewCloudinaryStorage(
		&internal.CloudinaryStorageConfig{
			Cloud:  "zanzanzan",
			Key:    "748773632958652",
			Secret: "a5puHSHwEyy12RtXBz44fPr104s",
		},
	)

	if err != nil {
		panic(err)
	}

	dbConn, err := database.NewDatabaseConnection(
		//fakeItTillUMakeIt,
		fmt.Sprintf(
			"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
			os.Getenv("MYSQLUSER"),
			os.Getenv("MYSQLPASSWORD"),
			os.Getenv("MYSQLHOST"),
			os.Getenv("MYSQLPORT"),
			os.Getenv("MYSQLDATABASE"),
		),
	)

	if err != nil {
		panic(err)
	}
	controllers.RegisterAuthController(api, dbConn)

	app.Use(jwtware.New(jwtware.Config{
		SigningKey: []byte("secret"),
	}))

	controllers.RegisterAdvertisementController(api, dbConn)
	controllers.RegisterUserController(api, dbConn)
	controllers.RegisterTaskController(api, dbConn, imageStore)

	panic(app.Listen( fmt.Sprintf(":%s", os.Getenv("PORT")) /* ":3000" */))
}
