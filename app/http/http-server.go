package main

import (
	"fmt"
	dbutil "github.com/ChrHan/go-sqlite-utility/dbutil"
	_ "github.com/mattn/go-sqlite3"
	"github.com/vrischmann/envconfig"
	"log"
	"net/http"
	"strconv"
	"syscall"
)

type config struct {
	// Filename of SQLite3 database
	Filename string `envconfig:"default=database.db"`

	// LogLevel is a minimal log severity required for the message to be logged.
	// Valid levels: [debug, info, warn, error, fatal, panic].
	LogLevel string `envconfig:"default=info"`
}

type appContext struct {
	db *dbutil.Dbutil
}

// DatabaseService serves as a sample HTTP Service
type DatabaseService struct {
	ctx appContext
}

// Select returns select result from SQLite3 Database
func (ds *DatabaseService) Select(w http.ResponseWriter, req *http.Request) {
	result, err := ds.ctx.db.Select()
	if err != nil {
		fmt.Println("Error found on select")
		log.Print(err.Error())
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
			fmt.Println(err.Error())
			log.Print(err.Error())
		}
		resultString += fmt.Sprintf("id: %s \nname: %s", strconv.Itoa(id), name)
	}
	err = result.Err()
	if err != nil {
		log.Println(err.Error())
	}
	fmt.Fprintf(w, resultString)
}

// Insert performs insert when called via HTTP using ?id={id}&product_name={product_name} as parameter
func (ds *DatabaseService) Insert(w http.ResponseWriter, req *http.Request) {
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
	err := ds.ctx.db.Insert(id, productName)
	if err != nil {
		log.Print(err.Error())
		w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		fmt.Fprintf(w, id+" "+productName+" is not inserted; duplicate ID found")
		return
	}

	fmt.Fprintf(w, id+" "+productName)
}

// Update performs update when called via HTTP using ?id={id}&product_name={product_name} as parameter
func (ds *DatabaseService) Update(w http.ResponseWriter, req *http.Request) {
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

	err := ds.ctx.db.Update(id, productName)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
	fmt.Fprintf(w, id+" "+productName)
}

// Delete performs delete when called via HTTP using ?id={id} as parameter
func (ds *DatabaseService) Delete(w http.ResponseWriter, req *http.Request) {
	ids, ok := req.URL.Query()["id"]
	if !ok || len(ids) < 1 {
		log.Println("Url Param 'id' is missing")
		return
	}

	id := ids[0]

	err := ds.ctx.db.Delete(id)
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
	fmt.Fprintf(w, id)
}

// DeleteAll performs delete all records when called via HTTP
func (ds *DatabaseService) DeleteAll(w http.ResponseWriter, req *http.Request) {
	err := ds.ctx.db.DeleteAll()
	if err != nil {
		log.Println(err.Error())
		w.WriteHeader(http.StatusInternalServerError) // Proper HTTP response
		return
	}
}

func main() {
	// -> config from env
	cfg := &config{}
	if err := envconfig.InitWithPrefix(&cfg, "APP"); err != nil {
		log.Println("init config: err=%s\n", err)
		syscall.Exit(1)
	}

	localCtx := &appContext{db: dbutil.New(cfg.Filename)}
	dbService := &DatabaseService{ctx: *localCtx}
	http.HandleFunc("/select", dbService.Select)
	http.HandleFunc("/insert", dbService.Insert)
	http.HandleFunc("/update", dbService.Update)
	http.HandleFunc("/deleteAll", dbService.DeleteAll)
	http.HandleFunc("/delete", dbService.Delete)

	fmt.Println("Server is starting....")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
