package routes

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/joho/godotenv"
	handler "github.com/krishnakantha1/to-do-list-backend/Handler"
)

//Application structure
type App struct {
	Routes *mux.Router
	DB     *pgxpool.Pool
}

//call to initilize the DB and Routes
func (a *App) InitializeRoutes() {

	err := godotenv.Load()
	if err != nil {
		log.Printf("godotenv cant find .env file")
	}

	db, err := pgxpool.Connect(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	a.DB = db
	a.Routes = mux.NewRouter()
	a.setRouter()
}

func (a *App) setRouter() {
	//Test route
	a.post("/test", a.wrapHandleFunc(handler.Test))

	//Routes to handle Category
	a.post("/add-category", a.wrapHandleFunc(handler.CreateCategory))
	a.get("/get-category-entry-count", a.wrapHandleFunc(handler.GetCategoryEntryCount))
	a.delete("/remove-category", a.wrapHandleFunc(handler.DeleteCategory))

	//Route to handle to do item List
	a.put("/item-toggle-done", a.wrapHandleFunc(handler.ItemMarkAsDone))
	a.delete("/item-delete", a.wrapHandleFunc(handler.DeleteItem))

	//Route to add to to tasks on to the DB
	a.post("/add-to-do-items", a.wrapHandleFunc(handler.AddToDoItemsToDB))

	//Routes to handle data request upon initial page load
	a.get("/", a.wrapHandleFunc(handler.InitializePage))
}

//Wraper to handle get methods
func (a *App) get(path string, f http.HandlerFunc) {
	a.Routes.HandleFunc(path, f).Methods("GET", "OPTIONS")
}

//Wraper to handle put methods
func (a *App) put(path string, f http.HandlerFunc) {
	a.Routes.HandleFunc(path, f).Methods("PUT", "OPTIONS")
}

//Wraper to handle post methods
func (a *App) post(path string, f http.HandlerFunc) {
	a.Routes.HandleFunc(path, f).Methods("POST", "OPTIONS")
}

//Wraper to handle delete methods
func (a *App) delete(path string, f http.HandlerFunc) {
	a.Routes.HandleFunc(path, f).Methods("DELETE", "OPTIONS")
}

//call to start the Application
func (a *App) AppStart() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}

	log.Println("Listening at port :" + port)
	if err := http.ListenAndServe(":"+port, a.Routes); err != nil {
		panic(err)
	}
}

//Wrapper to handle the handler func
type MyHandleFunc func(*pgxpool.Pool, http.ResponseWriter, *http.Request)

func (a *App) wrapHandleFunc(f MyHandleFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Content-Type", "text/html; charset=utf-8")

		if r.Method == "OPTIONS" {
			return
		}

		f(a.DB, w, r)
	}

}
