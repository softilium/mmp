package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	gh "github.com/gorilla/handlers"
	_ "github.com/lib/pq"
	"github.com/softilium/elorm"
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
	SPAroutes := []string{
		"/login", "/edit-shop", "/shop/", "/edit-good", "/checkout", "/search",
		"/orders", "/myprofile", "/profile/", "/set-roles", "/inc-orders", "/good/",
	}
	for _, route := range SPAroutes {
		router.HandleFunc(route, SpaHandler())
	}

	router.Handle("/", http.StripPrefix("/", fs))

	// API routes

	shopsRestApiConfig := elorm.CreateStdRestApiConfig(
		*DB.ShopDef.EntityDef,
		DB.LoadShop,
		DB.ShopDef.SelectEntities,
		DB.CreateShop)
	shopsRestApiConfig.DefaultPageSize = 0
	shopsRestApiConfig.DefaultSorts = func(r *http.Request) ([]*elorm.SortItem, error) {
		return []*elorm.SortItem{{Field: DB.ShopDef.Caption, Asc: true}}, nil
	}
	shopsRestApiConfig.Context = LoadUserFromHttpToContext
	shopsRestApiConfig.BeforeMiddleware = UserRequiredForEdit
	router.HandleFunc("/api/shops", elorm.HandleRestApi(shopsRestApiConfig))

	//// goods

	goodsRestApiConfig := elorm.CreateStdRestApiConfig(
		*DB.GoodDef.EntityDef,
		DB.LoadGood,
		DB.GoodDef.SelectEntities,
		DB.CreateGood)
	goodsRestApiConfig.DefaultPageSize = 0
	goodsRestApiConfig.AdditionalFilter = func(r *http.Request) ([]*elorm.Filter, error) {
		res := []*elorm.Filter{}
		shopref := r.URL.Query().Get("shopref")
		if shopref != "" {
			res = append(res, elorm.AddFilterEQ(DB.GoodDef.OwnerShop, shopref))
		}
		return res, nil
	}
	goodsRestApiConfig.DefaultSorts = func(r *http.Request) ([]*elorm.SortItem, error) {
		return []*elorm.SortItem{
			{Field: DB.GoodDef.OrderInShop, Asc: true},
			{Field: DB.GoodDef.Caption, Asc: true},
		}, nil
	}
	goodsRestApiConfig.Context = LoadUserFromHttpToContext
	router.HandleFunc("/api/goods", elorm.HandleRestApi(goodsRestApiConfig))

	//// allusers

	allusersRestApiConfig := elorm.CreateStdRestApiConfig(
		*DB.UserDef.EntityDef,
		DB.LoadUser,
		DB.UserDef.SelectEntities,
		DB.CreateUser)
	allusersRestApiConfig.DefaultPageSize = 0
	allusersRestApiConfig.EnablePost = false
	allusersRestApiConfig.BeforeMiddleware = UserAdminRequired
	allusersRestApiConfig.Context = LoadUserFromHttpToContext
	router.HandleFunc("/api/admin/allusers", elorm.HandleRestApi(allusersRestApiConfig))

	initRouterImages(router)
	initRouterAuth(router)
	initRouterBasket(router)
	initRouterOrders(router)
	initRouterSearchGoods(router)
	initRouterGoodTags(router)

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

	server.Handler = TelegramMiddleware(server.Handler)

	if Cfg.Debug {
		c1 := gh.AllowedOrigins(Cfg.CORSAlloweedHosts)
		c2 := gh.AllowCredentials()
		c3 := gh.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "OPTIONS"})
		c4 := gh.AllowedHeaders([]string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "Set-Cookie"})
		server.Handler = gh.CORS(c1, c2, c3, c4)(server.Handler)
	}

	return server

}

func createDefaultUser() {
	users, _, err := DB.UserDef.SelectEntities(nil, nil, 0, 0)
	if err != nil {
		logError(err)
	}
	if len(users) == 0 {
		fmt.Println("No users found, creating default admin user...")
		adminUser, err := DB.CreateUser()
		if err != nil {
			logError(err)
		}
		admPwdHash, _ := HashPassword(Cfg.AdminPassword)
		adminUser.SetEmail(Cfg.AdminEmail)
		adminUser.SetUsername(Cfg.AdminEmail)
		adminUser.SetPasswordHash(admPwdHash)
		adminUser.SetAdmin(true)
		adminUser.SetShopManager(true)
		err = adminUser.Save(context.Background())
		if err != nil {
			logError(err)
		}
	}
}

func removeOrphanedImages() {
	rows, err := DB.Query("select ref from goods")
	if err != nil {
		logError(err)
	}
	defer func() { _ = rows.Close() }()

	goodRefs := make(map[string]bool)
	for rows.Next() {
		var ref string
		err = rows.Scan(&ref)
		if err != nil {
			logError(err)
		}
		goodRefs[ref] = true
	}
	ifiles, err := os.ReadDir(Cfg.ImagesFolder)
	if err != nil {
		logError(err)
	}
	deletedImages := 0
	for _, file := range ifiles {

		fn := file.Name()
		if file.IsDir() {
			continue
		}
		tokens := strings.Split(fn, "-")
		if len(tokens) != 3 || tokens[0] != "goodImage" {
			continue
		}
		if _, ok := goodRefs[tokens[1]]; !ok {
			err = os.Remove(Cfg.ImagesFolder + "/" + file.Name())
			if err != nil {
				logError(err)
			}
			deletedImages++
		}
	}
	fmt.Printf("Deleted %d images for goods that don't exist\n", deletedImages)
}

func removeOldTokens() {
	tokensToDelete, _, err := DB.TokenDef.SelectEntities([]*elorm.Filter{
		elorm.AddFilterLT(DB.TokenDef.RefreshTokenExpiresAt, time.Now()),
	}, nil, 0, 0)
	if err != nil {
		logError(err)
	} else {
		for _, token := range tokensToDelete {
			err = DB.DeleteEntity(context.Background(), token.RefString())
			if err != nil {
				logError(err)
			}
		}
		fmt.Printf("Deleted %d expired tokens\n", len(tokensToDelete))
	}
}

func UserAdminRequired(w http.ResponseWriter, r *http.Request) bool {
	user, _, err := UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
		return false
	}
	if !user.Admin() {
		HandleErr(w, http.StatusForbidden, fmt.Errorf("user isn't admin"))
		return false
	}
	return true
}

func UserRequired(w http.ResponseWriter, r *http.Request) bool {
	_, _, err := UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
		return false
	}
	return true
}

func UserRequiredForEdit(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodGet {
		return true
	}
	_, _, err := UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
		return false
	}
	return true
}

func AdminRequiredForEdit(w http.ResponseWriter, r *http.Request) bool {
	if r.Method == http.MethodGet {
		return true
	}
	u, _, err := UserFromHttpRequest(r)
	if err != nil {
		HandleErr(w, http.StatusUnauthorized, fmt.Errorf("unauthorized: %v", err))
		return false
	}
	if !u.Admin() {
		HandleErr(w, http.StatusForbidden, fmt.Errorf("user isn't admin"))
	}

	return true
}

func main() {

	if Cfg.Debug {
		fmt.Printf("Debug mode is ON\n")
		fmt.Printf("Cfg.Dialect            = %s\n", Cfg.dbDialect)
		fmt.Printf("Cfg.DbConnectionString = %s\n", Cfg.DbConnectionString)
		fmt.Printf("Cfg.ListenAddr         = %s\n", Cfg.ListenAddr)
		fmt.Printf("Cfg.FrontEndFolder     = %s\n", Cfg.FrontEndFolder)
	}

	DB, err = CreateDbContext(Cfg.dbDialect, Cfg.DbConnectionString)
	logError(err)

	DB.AggressiveReadingCache = true

	err = DB.SetHandlers()
	logError(err)

	err = DB.EnsureDBStructure()
	logError(err)
	fmt.Println("Database structure ensured successfully.")

	removeOldTokens()
	removeOrphanedImages()
	createDefaultUser()

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
