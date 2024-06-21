package routes

import (
	"net/http"
	"spendid/app/http/controllers"
)

func Web() {
	http.HandleFunc("/", controllers.IndexHandler)
	http.HandleFunc("/login", controllers.LoginHandler)
	http.HandleFunc("/signup", controllers.SignupHandler)
	http.HandleFunc("/logout", controllers.LogoutHandler)
	http.HandleFunc("/create", controllers.CreateHandler)
	http.HandleFunc("/insert", controllers.InsertHandler)
	http.HandleFunc("/edit", controllers.EditHandler)
	http.HandleFunc("/update", controllers.UpdateHandler)
	http.HandleFunc("/delete", controllers.DeleteHandler)
	http.HandleFunc("/pemasukan", controllers.CreatePemasukanHandler)
	http.HandleFunc("/insert-pemasukan", controllers.InsertPemasukanHandler)
	http.HandleFunc("/pengeluaran", controllers.CreatePengeluaranHandler)
	http.HandleFunc("/insert-pengeluaran", controllers.InsertPengeluaranHandler)
	http.HandleFunc("/calc", controllers.CalcHandler)
	http.HandleFunc("/chat", controllers.ChatHandler)
	http.HandleFunc("/chatbot", controllers.ChatBotHandler)
	http.HandleFunc("/dashboard", controllers.UserDashboardHandler)
	http.HandleFunc("/pemasukan-user", controllers.CreatePemasukanUserHandler)
	http.HandleFunc("/insertpemasukan-user", controllers.InsertPemasukanUserHandler)
	http.HandleFunc("/pengeluaran-user", controllers.CreatePengeluaranUserHandler)
	http.HandleFunc("/insertpengeluaran-user", controllers.InsertPengeluaranUserHandler)
	http.HandleFunc("/reminder", controllers.InsertPengeluaranUserHandler)
	http.HandleFunc("/monthly-finance", controllers.MonthlyFinanceHandler)
	http.HandleFunc("/monthly-spending", controllers.MonthlySpendingHandler)

}
