package controllers

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"spendid/app/models"
)

func UserDashboardHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userIDRaw, ok := session.Values["userID"]
	if !ok {
		// Handle jika userID tidak ada dalam session
		http.Error(w, "UserID not found in session", http.StatusInternalServerError)
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		// Handle jika userID tidak dapat di-assert sebagai int
		http.Error(w, "UserID is not of type int", http.StatusInternalServerError)
		return
	}

	userName, ok := session.Values["userName"].(string)
	if !ok {
		log.Println("User Name is not set or invalid")
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	log.Println("Error parsing template:", userID)
	log.Println("Error parsing template:", userName)

	query := `SELECT PengeluaranID, UserID, UmkmID, Nominal, Date, KategoriPengeluaranID FROM Pengeluaran WHERE UserID = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		return
	}
	defer rows.Close()

	var pengeluarans []models.Pengeluaran
	for rows.Next() {
		var pengeluaran models.Pengeluaran
		if err := rows.Scan(&pengeluaran.PengeluaranID, &pengeluaran.UserID, &pengeluaran.UmkmID, &pengeluaran.Nominal, &pengeluaran.Date, &pengeluaran.KategoriPengeluaranID); err != nil {
			return
		}
		pengeluarans = append(pengeluarans, pengeluaran)
	}
	// pengeluarans, err := models.GetPengeluarans(db, userID)
	// if err != nil {
	// 	log.Println("Error fetching spendings:", err)
	// 	http.Error(w, "Error fetching spendings", http.StatusInternalServerError)
	// 	return
	// }

	reminders, err := models.GetReminders(db, userID)
	if err != nil {
		log.Println("Error fetching reminders:", err)
		http.Error(w, "Error fetching reminders", http.StatusInternalServerError)
		return
	}

	data := struct {
		UserName       string
		Spendings      []models.Pengeluaran
		MonthlyIncome  float64
		MonthlyOutcome float64
		Reminders      []models.Reminder
	}{
		UserName:       userName,
		Spendings:      pengeluarans,
		MonthlyIncome:  100000, // Replace with actual data
		MonthlyOutcome: 500000, // Replace with actual data
		Reminders:      reminders,
	}

	tmpl, err := template.ParseFiles("resources/views/templates/user/dashboard.html")
	if err != nil {

		log.Println("Error parsing template:", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func CreatePemasukanUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userIDRaw, ok := session.Values["userID"]
	if !ok {
		// Handle if userID does not exist in session
		http.Error(w, "UserID not found in session", http.StatusInternalServerError)
		return
	}

	userID, ok := userIDRaw.(int)
	if !ok {
		// Handle if userID cannot be asserted as int
		http.Error(w, "UserID is not of type int", http.StatusInternalServerError)
		return
	}

	query := `SELECT UmkmID, UmkmName FROM UMKM WHERE UserID = ?`
	rows, err := db.Query(query, userID)
	if err != nil {
		log.Println("Error fetching UMKMs:", err)
		http.Error(w, "Error fetching UMKMs", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var umkms []models.Umkm
	for rows.Next() {
		var umkm models.Umkm
		if err := rows.Scan(&umkm.UmkmID, &umkm.UmkmName); err != nil {
			log.Println("Error scanning UMKM:", err)
			http.Error(w, "Error scanning UMKM", http.StatusInternalServerError)
			return
		}
		umkms = append(umkms, umkm)
	}

	data := struct {
		UserID int
		UMKMs  []models.Umkm
	}{
		UserID: userID,
		UMKMs:  umkms,
	}

	tmpl, err := template.ParseFiles("resources/views/templates/user/pemasukan_user.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func InsertPemasukanUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	session, _ := store.Get(r, "session")
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID, ok := session.Values["userID"].(int)
	if !ok {
		log.Println("Error retrieving User ID from session")
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
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func CreatePengeluaranUserHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	userID := session.Values["userID"].(int)
	umkms, err := models.GetUMKMsByUserID(db, userID)
	if err != nil {
		log.Println("Error fetching UMKMs:", err)
		http.Error(w, "Error fetching UMKMs", http.StatusInternalServerError)
		return
	}

	kategoriPengeluarans, err := models.GetKategoriPengeluarans(db)
	if err != nil {
		log.Println("Error fetching Kategori Pengeluarans:", err)
		http.Error(w, "Error fetching Kategori Pengeluarans", http.StatusInternalServerError)
		return
	}

	data := struct {
		UserID               int
		UMKMs                []models.Umkm
		KategoriPengeluarans []models.KategoriPengeluaran
	}{
		UserID:               userID,
		UMKMs:                umkms,
		KategoriPengeluarans: kategoriPengeluarans,
	}

	tmpl, err := template.ParseFiles("resources/views/templates/user/pengeluaran_user.html")
	if err != nil {
		log.Println("Error parsing template:", err)
		http.Error(w, "Error parsing template", http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, data)
}

func InsertPengeluaranUserHandler(w http.ResponseWriter, r *http.Request) {
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
	http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
}

func MonthlyFinanceHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Baca file monthlyfinance.html
	data, err := ioutil.ReadFile("resources/views/templates/monthlyfinance.html") // Sesuaikan path
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Set header Content-Type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Tulis konten file ke response
	w.Write(data)
}

func MonthlySpendingHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Baca file monthlyfinance.html
	data, err := ioutil.ReadFile("resources/views/templates/monthlyspending.html") // Sesuaikan path
	if err != nil {
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	// Set header Content-Type
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Tulis konten file ke response
	w.Write(data)
}
