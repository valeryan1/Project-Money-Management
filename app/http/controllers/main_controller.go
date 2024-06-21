package controllers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"spendid/app/models"

	_ "github.com/go-sql-driver/mysql"
	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

func InitDB() {
	dataSourceName := "root:@tcp(127.0.0.1:3306)/money"
	db, err = sql.Open("mysql", dataSourceName)
	if err != nil {
		log.Fatal(err)
	}

	if err = db.Ping(); err != nil {
		log.Fatal(err)
	}
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	umkms, err := models.GetUmkms(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("resources/views/templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, umkms)
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		username := r.FormValue("username")
		password := r.FormValue("password")

		var dbPassword string
		err := db.QueryRow("SELECT UserPassword FROM user WHERE Username = ?", username).Scan(&dbPassword)
		if err != nil || password != dbPassword {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		var dbId int
		err = db.QueryRow("SELECT UserID FROM user WHERE Username = ?", username).Scan(&dbId)
		if err != nil {
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		session, _ := store.Get(r, "session")
		session.Values["authenticated"] = true
		session.Values["authenticated"] = true
		session.Values["userID"] = dbId
		session.Values["userName"] = username
		err = session.Save(r, w)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		session.Save(r, w)

		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	tmpl, err := template.ParseFiles("resources/views/templates/login.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func LogoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["authenticated"] = false
	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

func CreateHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	categories, err := models.GetUmkmCategories(db)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("resources/views/templates/create.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, categories)
}

func InsertHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	umkmCategoryID, err := strconv.Atoi(r.FormValue("umkmcategory_id"))
	if err != nil {
		http.Error(w, "Invalid UMKM Category ID", http.StatusBadRequest)
		return
	}

	umkm := models.Umkm{
		UmkmName:       r.FormValue("name"),
		UmkmCategoryID: umkmCategoryID,
		UmkmAddress:    r.FormValue("address"),
		UmkmNoTelp:     r.FormValue("notelp"),
	}

	query := `INSERT INTO UMKM (UmkmName, UmkmCategoryID, UmkmAddress, UmkmNoTelp) VALUES (?, ?, ?, ?)`
	_, err = db.Exec(query, umkm.UmkmName, umkm.UmkmCategoryID, umkm.UmkmAddress, umkm.UmkmNoTelp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func EditHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id := r.URL.Query().Get("id")

	var umkm models.Umkm
	row := db.QueryRow("SELECT UmkmID, UmkmName, UmkmCategoryID, UmkmAddress, UmkmNoTelp FROM UMKM WHERE UmkmID = ?", id)
	if err := row.Scan(&umkm.UmkmID, &umkm.UmkmName, &umkm.UmkmCategoryID, &umkm.UmkmAddress, &umkm.UmkmNoTelp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("resources/views/templates/update.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, umkm)
}

func UpdateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	umkmID, err := strconv.Atoi(r.FormValue("id"))
	if err != nil {
		http.Error(w, "Invalid UMKM ID", http.StatusBadRequest)
		return
	}

	umkm := models.Umkm{
		UmkmID:      umkmID,
		UmkmName:    r.FormValue("name"),
		UmkmAddress: r.FormValue("address"),
		UmkmNoTelp:  r.FormValue("notelp"),
	}

	query := `UPDATE UMKM SET UmkmName = ?, UmkmAddress = ?, UmkmNoTelp = ? WHERE UmkmID = ?`
	_, err = db.Exec(query, umkm.UmkmName, umkm.UmkmAddress, umkm.UmkmNoTelp, umkm.UmkmID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	if session.Values["authenticated"] != true {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	id := r.URL.Query().Get("id")

	query := `DELETE FROM UMKM WHERE UmkmID = ?`
	_, err := db.Exec(query, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func CalcHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("resources/views/templates/calc.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Parsing date of birth
		dob, err := time.Parse("2006-01-02", r.FormValue("dob"))
		if err != nil {
			http.Error(w, "Invalid date of birth", http.StatusBadRequest)
			return
		}

		// Creating a user instance from form values
		user := models.User{
			Name:         r.FormValue("name"),
			Username:     r.FormValue("username"),
			UserAddress:  r.FormValue("address"),
			UserEmail:    r.FormValue("email"),
			UserPhone:    r.FormValue("phone"),
			UserDOB:      dob,
			UserPassword: r.FormValue("password"),
		}

		// Check if the username already exists in the database
		var existingUser models.User
		err = db.QueryRow("SELECT UserID, UmkmID, Name, Username, UserAddress, UserEmail, UserPhone, UserDOB, UserPassword FROM user WHERE Username = ?", user.Username).Scan(
			&existingUser.UserID, &existingUser.UmkmID, &existingUser.Name, &existingUser.Username, &existingUser.UserAddress, &existingUser.UserEmail, &existingUser.UserPhone, &existingUser.UserDOB, &existingUser.UserPassword)
		if err != nil {
			if err != sql.ErrNoRows {
				http.Error(w, "awal Server error"+r.FormValue("username"), http.StatusInternalServerError)
				return
			}
		} else {
			http.Error(w, "Username already exists", http.StatusConflict)
			return
		}

		// Insert the new user into the database
		if err := user.Create(db); err != nil {
			http.Error(w, "insert db Server error"+user.UserEmail, http.StatusInternalServerError)
			return
		}

		var userID int
		err = db.QueryRow("SELECT UserID FROM user WHERE Username = ?", user.Username).Scan(&userID)
		if err != nil {
			http.Error(w, "Server error: "+err.Error(), http.StatusInternalServerError)
			return
		}

		// Set session values
		session, _ := store.Get(r, "session")
		session.Values["authenticated"] = true
		session.Values["userName"] = user.Username
		session.Values["userID"] = userID
		session.Save(r, w)

		// Redirect to home page
		http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		return
	}

	// Render the signup template for GET requests
	tmpl, err := template.ParseFiles("resources/views/templates/signup.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ChatHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("resources/views/templates/chat.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	tmpl.Execute(w, nil)
}

func ChatBotHandler(w http.ResponseWriter, r *http.Request) {
	var msg models.ChatResponse
	if r.Method == "POST" {
		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&msg)
		if err != nil {
			http.Error(w, "Error decoding JSON", http.StatusBadRequest)
			return
		}
		fmt.Println("Pesan baru:", msg.Message)
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{"response": "Pesan diterima"}
		json.NewEncoder(w).Encode(response)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("AIzaSyDcCRsjI4Mep_b-VYTr8_LHb5cwsJ5TXmU")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()
	model := client.GenerativeModel("gemini-1.5-flash")
	cs := model.StartChat()

	send := func(msg string) *genai.GenerateContentResponse {
		fmt.Printf("== Me: %s\n== My AI:\n", msg)
		res, err := cs.SendMessage(ctx, genai.Text(msg))
		if err != nil {
			log.Fatal(err)
		}
		return res
	}

	res := send(msg.Message)
	if err != nil {
		log.Fatal(err)
	}
	printResponse(res)
}

func printResponse(resp *genai.GenerateContentResponse) {
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				fmt.Println(part)
			}
		}
	}
	fmt.Println("---")
}
