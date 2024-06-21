package controllers

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"spendid/app/models"
	"strconv"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/sessions"
)

var db *sql.DB
var err error
var store = sessions.NewCookieStore([]byte("secret-key"))

func CreatePemasukanHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")

	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	pemasukan, err := models.GetPemasukans(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("resources/views/templates/pemasukan.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, pemasukan)
}

func InsertPemasukanHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userID, err := strconv.Atoi(r.FormValue("userid"))
	if err != nil {
		log.Println("Error parsing User ID:", err)
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	umkmID, err := strconv.Atoi(r.FormValue("umkmid"))
	if err != nil {
		log.Println("Error parsing UMKM ID:", err)
		http.Error(w, "Invalid UMKM ID", http.StatusBadRequest)
		return
	}

	nominal, err := strconv.ParseFloat(r.FormValue("nominal"), 64)
	if err != nil {
		log.Println("Error parsing Nominal:", err)
		http.Error(w, "Invalid Nominal", http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Println("Error parsing Date:", err)
		http.Error(w, "Invalid Date", http.StatusBadRequest)
		return
	}

	pemasukan := models.Pemasukan{
		UserID:  userID,
		UmkmID:  umkmID,
		Nominal: nominal,
		Date:    parsedDate.Format("2006-01-02"),
	}

	query := `INSERT INTO Pemasukan (UserID, UmkmID, Nominal, date) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, pemasukan.UserID, pemasukan.UmkmID, pemasukan.Nominal, pemasukan.Date)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Pemasukan successfully inserted:", pemasukan)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CreatePengeluaranHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	KategoriPengeluaranID, err := models.GetKategoriPengeluarans(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("resources/views/templates/pengeluaran.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, KategoriPengeluaranID)
}

func InsertPengeluaranHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	userID, err := strconv.Atoi(r.FormValue("userid"))
	if err != nil {
		log.Println("Error parsing User ID:", err)
		http.Error(w, "Invalid User ID", http.StatusBadRequest)
		return
	}

	umkmID, err := strconv.Atoi(r.FormValue("umkmid"))
	if err != nil {
		log.Println("Error parsing UMKM ID:", err)
		http.Error(w, "Invalid UMKM ID", http.StatusBadRequest)
		return
	}

	nominal, err := strconv.ParseFloat(r.FormValue("nominal"), 64)
	if err != nil {
		log.Println("Error parsing Nominal:", err)
		http.Error(w, "Invalid Nominal", http.StatusBadRequest)
		return
	}

	date := r.FormValue("date")
	parsedDate, err := time.Parse("2006-01-02", date)
	if err != nil {
		log.Println("Error parsing Date:", err)
		http.Error(w, "Invalid Date", http.StatusBadRequest)
		return
	}

	kategoriID, err := strconv.Atoi(r.FormValue("kategori"))
	if err != nil {
		log.Println("Error parsing Kategori ID:", err)
		http.Error(w, "Invalid Kategori Pengeluaran ID", http.StatusBadRequest)
		return
	}

	pengeluaran := models.Pengeluaran{
		UserID:                userID,
		UmkmID:                umkmID,
		Nominal:               nominal,
		Date:                  parsedDate.Format("2006-01-02"),
		KategoriPengeluaranID: kategoriID,
	}

	query := `INSERT INTO Pengeluaran (UserID, UmkmID, Nominal, date, KategoriPengeluaranID) VALUES (?, ?, ?, ?, ?)`
	_, err = db.Exec(query, pengeluaran.UserID, pengeluaran.UmkmID, pengeluaran.Nominal, pengeluaran.Date, pengeluaran.KategoriPengeluaranID)
	if err != nil {
		log.Println("Error executing query:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Pengeluaran successfully inserted:", pengeluaran)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}
