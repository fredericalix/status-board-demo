package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	echoSwagger "github.com/swaggo/echo-swagger"

	_ "backend-server/docs"
	"backend-server/sse"

	"github.com/lib/pq"
	_ "github.com/lib/pq"

	"backend-server/api"
	"backend-server/config"
)

func main() {
	// Créer une nouvelle instance d'Echo
	e := echo.New()
	// Ajouter le middleware CORS avec une configuration personnalisée
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
		AllowHeaders: []string{echo.HeaderContentType, echo.HeaderAccept},
	}))
	// Middleware
	e.Use(middleware.Logger())  // Utilisation du middleware de journalisation
	e.Use(middleware.Recover()) // Utilisation du middleware de récupération d'erreurs

	// Se connecter à la base de données et vérifier la connexion
	db, err := config.ConnectToDB() // Connexion à la base de données
	if err != nil {
		fmt.Println("Erreur lors de la connexion à la base de données:", err)
		os.Exit(1)
	}

	// Vérifier la connexion à la base de données
	err = db.Ping()
	if err != nil {
		fmt.Println("La connexion à la base de données a échoué:", err)
		os.Exit(1)
	}
	defer db.Close()

	// S'abonner aux mises à jour de statut dans la base de données
	_, err = db.Exec("LISTEN status_update")
	if err != nil {
		log.Fatal(err)
	}

	// Création d'un canal pour les notifications de mise à jour de statut
	notificationChannel := make(chan *pq.Notification)

	// Création une goroutine pour écouter les notifications de mise à jour de statut
	go func() {
		for {
			// Configuration pour les notifications de mise à jour de statut
			dbURI := os.Getenv("POSTGRESQL_ADDON_URI")

			notificationConn := pq.NewListener(
				dbURI,
				10*time.Second,
				time.Minute,
				func(ev pq.ListenerEventType, err error) {
					if err != nil {
						log.Println(err)
					}
				},
			)
			notificationConn.Notify = notificationChannel

			// Écouter les notifications de mise à jour de statut
			err = notificationConn.Listen("status_update")
			if err != nil {
				log.Println(err)
				time.Sleep(5 * time.Second)
				continue
			}
			defer notificationConn.Close()

			// Attendre les notifications de mise à jour de statut
			select {
			case <-notificationConn.Notify:
			case <-time.After(90 * time.Second):
				if err := notificationConn.Ping(); err != nil {
					log.Println("Ping failed:", err)
					return
				}
			}
		}
	}()

	// Création un broker SSE (Server-Sent Events)
	b := sse.NewSSEBroker()
	go b.Start()

	// Créer une goroutine pour traiter les notifications de mise à jour de statut
	go func() {
		for {
			select {
			case notification := <-notificationChannel:
				log.Println("Received notification:", notification)
				// Traiter le message ici
				in := `{"notif":"maj status"}`
				b.NotifyClients(in)
			}
		}
	}()

	// Les endpoints de notre api
	e.GET("/", api.ApiStatus)
	e.GET("/swagger/*", echoSwagger.WrapHandler)
	e.POST("/status", api.CreateStatus(b))
	e.PUT("/status/:id", api.UpdateStatus(b))
	e.DELETE("/status/:id", api.DeleteStatus(b))
	e.GET("/status/:id", api.GetStatusByID)
	e.GET("/status", api.GetAllStatus)
	e.GET("/events", sse.SSEHandler(b))
	e.GET("/crash", api.CrashApp)

	// Démarrage du serveur. Le port par def est le 8080. Il peut être defini par la variable d'env PORT
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// Lancement serveur sur le port spécifié
	e.Logger.Fatal(e.Start(":" + port))
}
