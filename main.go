package main

import (
	"log"
	"net/http"
	"spendid/app/http/controllers"
	"spendid/routes"
)

func main() {

	controllers.InitDB()
	routes.Web()

	// router := gin.Default()

	// router.GET("/login", func(c *gin.Context) {
	// 	c.HTML(200, "login.html", nil)
	// })

	// if err := router.Run(":8080"); err != nil {
	// 	log.Fatal(err)
	// }

	log.Println("Server started at :8800")
	log.Fatal(http.ListenAndServe(":8800", nil))

}
