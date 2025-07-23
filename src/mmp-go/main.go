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

var err error

func logError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
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
	router.HandleFunc("/shop/{shopref}", SpaHandler())

	// API routes

	//// shops

	shopsRestApiConfig := elorm.CreateStdRestApiConfig(
		*models.Dbc.ShopDef.EntityDef,
		models.Dbc.LoadShop,
		models.Dbc.ShopDef.SelectEntities,
		models.Dbc.CreateShop)
	shopsRestApiConfig.DefaultPageSize = 0
	shopsRestApiConfig.AdditionalFilter = func(r *http.Request) []*elorm.Filter {
		res := []*elorm.Filter{}
		res = append(res, elorm.AddFilterEQ(models.Dbc.ShopDef.IsDeleted, false))
		return res
	}
	shopsRestApiConfig.DefaultSorts = func(r *http.Request) []*elorm.SortItem {
		return []*elorm.SortItem{
			{Field: models.Dbc.ShopDef.Caption, Asc: true},
		}
	}
	router.HandleFunc("/api/shops", elorm.HandleRestApi(shopsRestApiConfig))

	//// goods

	goodsRestApiConfig := elorm.CreateStdRestApiConfig(
		*models.Dbc.GoodDef.EntityDef,
		models.Dbc.LoadGood,
		models.Dbc.GoodDef.SelectEntities,
		models.Dbc.CreateGood)
	goodsRestApiConfig.DefaultPageSize = 0
	goodsRestApiConfig.AdditionalFilter = func(r *http.Request) []*elorm.Filter {
		res := []*elorm.Filter{}
		res = append(res, elorm.AddFilterEQ(models.Dbc.GoodDef.IsDeleted, false))
		shopref := r.URL.Query().Get("shopref")
		if shopref != "" {
			res = append(res, elorm.AddFilterEQ(models.Dbc.GoodDef.OwnerShop, shopref))
		}
		return res
	}
	goodsRestApiConfig.DefaultSorts = func(r *http.Request) []*elorm.SortItem {
		return []*elorm.SortItem{
			{Field: models.Dbc.GoodDef.OrderInShop, Asc: true},
		}
	}
	router.HandleFunc("/api/goods", elorm.HandleRestApi(goodsRestApiConfig))

	//// allusers

	allusersRestApiConfig := elorm.CreateStdRestApiConfig(
		*models.Dbc.UserDef.EntityDef,
		models.Dbc.LoadUser,
		models.Dbc.UserDef.SelectEntities,
		models.Dbc.CreateUser)
	allusersRestApiConfig.DefaultPageSize = 0
	allusersRestApiConfig.EnableDelete = true
	allusersRestApiConfig.EnablePost = false
	allusersRestApiConfig.EnablePut = true
	allusersRestApiConfig.BeforeMiddleware = func(w http.ResponseWriter, r *http.Request) bool {
		user, err := models.UserFromHttpRequest(r)
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
	allusersRestApiConfig.Context = models.HttpUserContext
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

	models.Dbc, err = models.CreateDbContext(Cfg.dbDialect, Cfg.DbConnectionString)
	logError(err)

	models.Dbc.AggressiveReadingCache = true

	err = models.Dbc.SetHandlers()
	logError(err)

	err = models.Dbc.EnsureDBStructure()
	logError(err)
	fmt.Println("Database structure ensured successfully.")

	users, _, err := models.Dbc.UserDef.SelectEntities(nil, nil, 0, 0)
	if err != nil {
		logError(err)
	}
	if len(users) == 0 {
		fmt.Println("No users found, creating default admin user...")
		adminUser, err := models.Dbc.CreateUser()
		if err != nil {
			logError(err)
		}
		adminUser.SetEmail(Cfg.AdminEmail)
		adminUser.SetUsername(Cfg.AdminEmail)
		adminUser.SetPassword(Cfg.AdminPassword)
		adminUser.SetAdmin(true)
		adminUser.SetShopManager(true)
		err = adminUser.Save(context.Background())
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

//TODO profile
//TODO orders+lines
//TODO goods images
//TODO telegram middleware
//TODO unmarshal json, autoexpand support
