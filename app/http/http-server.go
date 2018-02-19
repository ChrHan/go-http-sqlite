package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	dbutil "github.com/ChrHan/go-sqlite-utility/dbutil"
	_ "github.com/mattn/go-sqlite3"
)

type appContext struct {
	val int
	db  *dbutil.Dbutil
}

// ServiceA serves as a sample HTTP Service
type ServiceA struct {
	ctx appContext
	a   int
}

// Foo returns sa.a
func (sa *ServiceA) Foo(w http.ResponseWriter, req *http.Request) {
	sa.a = 1
	fmt.Fprintf(w, "sa is now %s", strconv.Itoa(sa.a))
}

// Bar returns sa.a
func (sa *ServiceA) Bar(w http.ResponseWriter, req *http.Request) {
	sa.a = 2
	fmt.Fprintf(w, "sa is now %s", strconv.Itoa(sa.a))
}

// Select returns select result from SQLite3 Database
func (sa *ServiceA) Select(w http.ResponseWriter, req *http.Request) {
	result, err := sa.ctx.db.Select()
	if err != nil {
		fmt.Println("Error found on select")
		log.Fatal(err.Error())
		w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
	var resultString string
	defer result.Close()
	for result.Next() {
		var id int
		var name string
		err := result.Scan(&id, &name)
		if err != nil {
			log.Fatal(err.Error())
		}
		resultString += fmt.Sprintf("id: %s \nname: %s", strconv.Itoa(id), name)
	}
	err = result.Err()
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Fprintf(w, resultString)
}

// Insert performs insert when called via HTTP using ?id={id}&product_name={product_name} as parameter
func (sa *ServiceA) Insert(w http.ResponseWriter, req *http.Request) {
	ids, ok := req.URL.Query()["id"]
	if !ok || len(ids) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}
	productNames, ok := req.URL.Query()["product_name"]
	if !ok || len(productNames) < 1 {
		log.Println("Url Param 'product_name' is missing")
		return
	}

	id := ids[0]
	productName := productNames[0]

	err := sa.ctx.db.Insert(id, productName)
	if err != nil {
		log.Fatal(err.Error())
		w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}

	fmt.Fprintf(w, id+" "+productName)
}

// Update performs update when called via HTTP using ?id={id}&product_name={product_name} as parameter
func (sa *ServiceA) Update(w http.ResponseWriter, req *http.Request) {
	ids, ok := req.URL.Query()["id"]
	if !ok || len(ids) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}
	productNames, ok := req.URL.Query()["product_name"]
	if !ok || len(productNames) < 1 {
		log.Println("Url Param 'product_name' is missing")
		return
	}

	id := ids[0]
	productName := productNames[0]

	err := sa.ctx.db.Update(id, productName)
	if err != nil {
		log.Fatal(err.Error())
		w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
	fmt.Fprintf(w, id+" "+productName)
}

// Delete performs delete when called via HTTP using ?id={id} as parameter
func (sa *ServiceA) Delete(w http.ResponseWriter, req *http.Request) {
	ids, ok := req.URL.Query()["id"]
	if !ok || len(ids) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}

	id := ids[0]

	err := sa.ctx.db.Delete(id)
	if err != nil {
		log.Fatal(err.Error())
		w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
	fmt.Fprintf(w, id)
}

// DeleteAll performs delete all records when called via HTTP
func (sa *ServiceA) DeleteAll(w http.ResponseWriter, req *http.Request) {
	err := sa.ctx.db.DeleteAll()
	if err != nil {
		log.Fatal(err.Error())
		w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
}

func main() {
	localCtx := &appContext{val: 42, db: dbutil.New("database.db")}
	a := &ServiceA{ctx: *localCtx, a: 28}
	http.HandleFunc("/a/foo", a.Foo)
	http.HandleFunc("/a/bar", a.Bar)
	http.HandleFunc("/select", a.Select)
	http.HandleFunc("/insert", a.Insert)
	http.HandleFunc("/update", a.Update)
	http.HandleFunc("/deleteAll", a.DeleteAll)
	http.HandleFunc("/delete", a.Delete)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
