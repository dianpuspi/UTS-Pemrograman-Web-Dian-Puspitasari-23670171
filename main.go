package main

import (
	"html/template"
	"log"
	"net/http"
	"strconv"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Kategori struct {
	ID           uint   `gorm:"column:id_kategori"`
	NamaKategori string `gorm:"column:nama_kategori"`
}

type Produk struct {
	IDProduk   uint     `gorm:"primaryKey;column:id_produk"`
	NamaProduk string   `gorm:"column:nama_produk"`
	Harga      int      `gorm:"column:harga"`
	KategoriID uint     `gorm:"column:id_kategori"`
	Kategori   Kategori `gorm:"foreignKey:KategoriID;references:ID"`
}

func connectDB() (*gorm.DB, error) {
	dsn := "root:@tcp(127.0.0.1:3306)/minimarket?charset=utf8mb4&parseTime=True&loc=Local"
	return gorm.Open(mysql.Open(dsn), &gorm.Config{})
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		http.Error(w, "Gagal koneksi ke database", http.StatusInternalServerError)
		return
	}

	var produk []Produk
	var kategori []Kategori
	db.Preload("Kategori").Find(&produk)
	db.Find(&kategori)

	tmpl := template.Must(template.ParseFiles("template/index.html"))
	tmpl.Execute(w, map[string]interface{}{
		"Produk":   produk,
		"Kategori": kategori,
	})
}

func tambahHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nama := r.FormValue("nama_produk")
		harga, _ := strconv.Atoi(r.FormValue("harga"))
		kategoriID, _ := strconv.Atoi(r.FormValue("kategori_id"))

		db, err := connectDB()
		if err != nil {
			http.Error(w, "Gagal koneksi ke database", http.StatusInternalServerError)
			return
		}

		db.Create(&Produk{NamaProduk: nama, Harga: harga, KategoriID: uint(kategoriID)})
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func editHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		id, _ := strconv.Atoi(r.FormValue("id"))
		nama := r.FormValue("nama_produk")
		harga, _ := strconv.Atoi(r.FormValue("harga"))
		kategoriID, _ := strconv.Atoi(r.FormValue("kategori_id"))

		db, err := connectDB()
		if err != nil {
			http.Error(w, "Gagal koneksi ke database", http.StatusInternalServerError)
			return
		}

		db.Model(&Produk{}).Where("id_produk = ?", id).Updates(Produk{
			NamaProduk: nama,
			Harga:      harga,
			KategoriID: uint(kategoriID),
		})
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func hapusHandler(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Query().Get("id"))

	db, err := connectDB()
	if err != nil {
		http.Error(w, "Gagal koneksi ke database", http.StatusInternalServerError)
		return
	}

	db.Delete(&Produk{}, id)
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func main() {
	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/tambah", tambahHandler)
	http.HandleFunc("/edit", editHandler)
	http.HandleFunc("/hapus", hapusHandler)

	log.Println("Server running at :8080")
	http.ListenAndServe(":8080", nil)
}
