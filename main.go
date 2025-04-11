package main

import (
	"encoding/json"
	"html/template"
	"net/http"
	"os"
)

type ContactInfo struct {
	Telegram string `json:"telegram"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
}

var contactInfo ContactInfo
var adminPassword = "admin123"

func main() {
	loadContacts()
//
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/admin", adminHandler)
	http.HandleFunc("/update-contacts", updateContactsHandler)

	http.ListenAndServe(":8080", nil)
}
// Путь к файлу с контактами
func loadContacts() {
	data, err := os.ReadFile("contacts.json")
	if err == nil {
		json.Unmarshal(data, &contactInfo)
	}
}
// Сохранение контактов в файл
func saveContacts() {
	data, _ := json.MarshalIndent(contactInfo, "", "  ")
	_ = os.WriteFile("contacts.json", data, 0644)
}
// Функция для рендеринга шаблонов
func renderTemplate(w http.ResponseWriter, tmpl string, data any) {
	tmplPath := "templates/" + tmpl
	t, _ := template.ParseFiles(tmplPath)
	t.Execute(w, data)
}
// Обработчик главной страницы
func homeHandler(w http.ResponseWriter, r *http.Request) {
	renderTemplate(w, "index.html", contactInfo)
}
// Обработчик страницы авторизации
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		pass := r.FormValue("password")
		if pass == adminPassword {
			http.SetCookie(w, &http.Cookie{Name: "auth", Value: "true"})
			http.Redirect(w, r, "/admin", http.StatusSeeOther)
			return
		}
	}
	renderTemplate(w, "login.html", nil)
}
// Обработчик страницы администратора
func adminHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("auth")
	if err != nil || cookie.Value != "true" {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}
	renderTemplate(w, "admin.html", contactInfo)
}
// Обработчик обновления контактов
func updateContactsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		contactInfo.Telegram = r.FormValue("telegram")
		contactInfo.Email = r.FormValue("email")
		contactInfo.Phone = r.FormValue("phone")
		saveContacts()
		http.Redirect(w, r, "/admin", http.StatusSeeOther)
	}
}
