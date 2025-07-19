package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	gh "github.com/gorilla/handlers"
	_ "github.com/lib/pq"
	"github.com/softilium/elorm"
	"github.com/softilium/mmp-go/models"
	_ "modernc.org/sqlite"
)

var dbc *models.DbContext
var err error

func logError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

func HandleErr(w http.ResponseWriter, status int, err error) {
	if status == 0 {
		status = http.StatusInternalServerError
	}
	errStr := ""
	if err != nil {
		log.Println(err.Error())
		errStr = err.Error()
	}
	http.Error(w, errStr, status)
}

func gracefullShutdown(server *http.Server, quit <-chan os.Signal, done chan<- bool) {

	<-quit
	log.Println("Server is shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Could not gracefully shutdown the server: %v\n", err)
	}
	close(done)

}

func SpaHandler() func(w http.ResponseWriter, r *http.Request) {

	fn := func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, Cfg.FrontEndFolder+"/index.html")
	}
	return http.HandlerFunc(fn)

}

func initServer() *http.Server {

	router := http.NewServeMux()

	fs := http.FileServer(http.Dir(Cfg.FrontEndFolder))

	// SPA routes
	router.Handle("/", http.StripPrefix("/", fs))
	router.HandleFunc("/login", SpaHandler())

	// API routes
	shopsRestApiConfig := elorm.CreateStdRestApiConfig(
		*dbc.ShopDef.EntityDef,
		dbc.LoadShop,
		dbc.ShopDef.SelectEntities,
		dbc.CreateShop)
	shopsRestApiConfig.DefaultPageSize = 0
	router.HandleFunc("/api/shops", elorm.HandleRestApi(shopsRestApiConfig))

	allusersRestApiConfig := elorm.CreateStdRestApiConfig(
		*dbc.UserDef.EntityDef,
		dbc.LoadUser,
		dbc.UserDef.SelectEntities,
		dbc.CreateUser)
	allusersRestApiConfig.DefaultPageSize = 0
	allusersRestApiConfig.EnableDelete = true
	allusersRestApiConfig.EnablePost = false
	allusersRestApiConfig.EnablePut = true

	allusersRestApiConfig.BeforeMiddleware = func(w http.ResponseWriter, r *http.Request) bool {
		user, err := UserFromHttpRequest(r)
		if err != nil {
			HandleErr(w, 0, fmt.Errorf("unauthorized: %v", err))
			return false
		}
		if !user.Admin() {
			HandleErr(w, http.StatusForbidden, fmt.Errorf("user isn't admin"))
			return false
		}
		return true
	}
	router.HandleFunc("/api/admin/allusers", elorm.HandleRestApi(allusersRestApiConfig))

	router.HandleFunc("/api/admin/migrate", Migrate)

	// AUTH routes

	router.HandleFunc("/identity/register", UserRegister)
	router.HandleFunc("/identity/login", UserLogin)
	router.HandleFunc("/identity/logout", UserLogout)
	router.HandleFunc("/identity/myprofile", UserPublicProfile)

	// CORE

	server := &http.Server{
		Addr:    Cfg.ListenAddr,
		Handler: router,
	}
	if Cfg.Debug {
		server.ReadTimeout = 5 * time.Minute
		server.WriteTimeout = 10 * time.Minute
		server.IdleTimeout = 15 * time.Minute
	}

	if !Cfg.Debug {
		server.Handler = gh.RecoveryHandler()(server.Handler)
	}
	server.Handler = gh.LoggingHandler(os.Stdout, server.Handler)

	if Cfg.Debug {
		c1 := gh.AllowedOrigins([]string{"http://localhost:56056"}) //vita dev server addr
		c2 := gh.AllowCredentials()
		c3 := gh.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
		c4 := gh.AllowedHeaders([]string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Set-Cookie"})
		server.Handler = gh.CORS(c1, c2, c3, c4)(server.Handler)
	}

	return server

}

func main() {

	if Cfg.Debug {
		fmt.Printf("Debug mode is ON\n")
		fmt.Printf("Cfg.Dialect            = %s\n", Cfg.dbDialect)
		fmt.Printf("Cfg.DbConnectionString = %s\n", Cfg.DbConnectionString)
		fmt.Printf("Cfg.ListenAddr         = %s\n", Cfg.ListenAddr)
		fmt.Printf("Cfg.FrontEndFolder     = %s\n", Cfg.FrontEndFolder)

	}

	dbc, err = models.CreateDbContext(Cfg.dbDialect, Cfg.DbConnectionString)
	logError(err)

	dbc.AggressiveReadingCache = true

	err = dbc.SetHandlers()
	logError(err)

	err = dbc.EnsureDBStructure()
	logError(err)
	fmt.Println("Database structure ensured successfully.")

	users, _, err := dbc.UserDef.SelectEntities(nil, nil, 0, 0)
	if err != nil {
		logError(err)
	}
	if len(users) == 0 {
		fmt.Println("No users found, creating default admin user...")
		adminUser, err := dbc.CreateUser()
		if err != nil {
			logError(err)
		}
		adminUser.SetEmail(Cfg.AdminEmail)
		adminUser.SetUsername(Cfg.AdminEmail)
		adminUser.SetPassword(Cfg.AdminPassword)
		adminUser.SetAdmin(true)
		adminUser.SetShopManager(true)
		err = adminUser.Save()
		if err != nil {
			logError(err)
		}
	}

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)

	server := initServer()
	go gracefullShutdown(server, quit, done)

	fmt.Printf("Listening on http://%s\n", Cfg.ListenAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("could not listen on %s: %v\n", Cfg.ListenAddr, err.Error())
	}

	<-done

}
