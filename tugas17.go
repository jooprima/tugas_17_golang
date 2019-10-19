package main

import "encoding/json"
import "net/http"
import "fmt"
import _ "mysql-master"
import "database/sql"

//struct sebagai penampung data hasil query
type data_karyawan struct {
	ID      int
	Nama    string
	Umur    int
	Jabatan string
}

//fungsi koneksi ke database mysql
func koneksi() (*sql.DB, error) {
	db, err := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/db_tugas16golang")

	if err != nil {
		return nil, err
	}
	return db, nil
}

var data []data_karyawan

// fungsi utama
func main() {
	ambil_data()
	http.HandleFunc("/karyawan", ambil_karyawan)
	http.HandleFunc("/cari_karyawan", cari_karyawan)

	fmt.Println("running web server on localhost:8080")
	http.ListenAndServe(":8080", nil)
}

//fungsi menampilkan data langsung di postman
func ambil_karyawan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var result, err = json.Marshal(data)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Write(result)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

//fungsi mencari data di postman
func cari_karyawan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "POST" {
		var nama = r.FormValue("Nama")
		var result []byte
		var err error

		for _, each := range data {
			if each.Nama == nama {
				result, err = json.Marshal(each)

				if err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
				w.Write(result)
				return
			}
		}
		http.Error(w, "data karyawan tidak tersedia", http.StatusBadRequest)
		return
	}
	http.Error(w, "", http.StatusBadRequest)
}

//fungsi mengambil data dari database ke postman
func ambil_data() {
	db, err := koneksi()

	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer db.Close()

	rows, err := db.Query("select * from tbl_karyawan")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer rows.Close()

	for rows.Next() {
		var each = data_karyawan{}
		var err = rows.Scan(&each.ID, &each.Nama, &each.Umur, &each.Jabatan)

		if err != nil {
			fmt.Println(err.Error())
			return
		}
		data = append(data, each)
	}
	if err = rows.Err(); err != nil {
		fmt.Println(err.Error())
		return
	}
}
