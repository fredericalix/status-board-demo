package api

import (
	"backend-server/config"
	"backend-server/model"
	"backend-server/sse"
	"database/sql"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/labstack/echo/v4"
)

// @Summary Récupérer le status de l'API
// @Description Récupérer le status de l'API
// @Produce  text/plain
// @Success 200 {string} string "hello from clever cloud"
// @Router / [get]
func ApiStatus(c echo.Context) error {
	return c.String(http.StatusOK, "hello from clever cloud")
}

// @Summary Créer un nouveau status
// @Description Ajoute un nouveau status à la base de données
// @Accept  json
// @Produce  json
// @Param status body model.Status true "Status à ajouter"
// @Success 201 {object} model.Status
// @Router /status [post]
func CreateStatus(b *sse.SSEBroker) echo.HandlerFunc {
	return func(c echo.Context) error {
		d := new(model.Status)
		if err := c.Bind(d); err != nil {
			return err
		}

		db, err := config.ConnectToDB()
		if err != nil {
			return err
		}
		defer db.Close()

		query := `INSERT INTO status (designation, state) VALUES ($1, $2) RETURNING id`
		err = db.QueryRow(query, d.Designation, d.State).Scan(&d.ID)
		if err != nil {
			return err
		}

		b.NotifyClients("Status créé")

		return c.JSON(http.StatusCreated, d)
	}
}

// @Summary Modifier un status existant
// @Description Met à jour un status existant dans la base de données
// @Accept  json
// @Produce  json
// @Param id path int true "ID du status à modifier"
// @Param status body model.Status true "Nouvelles données du status"
// @Success 200 {object} model.Status
// @Router /status/{id} [put]
func UpdateStatus(b *sse.SSEBroker) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		d := new(model.Status)
		if err := c.Bind(d); err != nil {
			return err
		}

		db, err := config.ConnectToDB()
		if err != nil {
			return err
		}
		defer db.Close()

		query := `UPDATE status SET designation = $1, state = $2 WHERE id = $3 RETURNING id`
		err = db.QueryRow(query, d.Designation, d.State, id).Scan(&d.ID)
		if err != nil {
			return err
		}

		b.NotifyClients("Status mis à jour")

		return c.JSON(http.StatusOK, d)
	}
}

// @Summary Supprimer un status
// @Description Supprime un status de la base de données
// @Produce  json
// @Param id path int true "ID du status à supprimer"
// @Success 204 {string} string ""
// @Router /status/{id} [delete]
func DeleteStatus(b *sse.SSEBroker) echo.HandlerFunc {
	return func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			return err
		}

		db, err := config.ConnectToDB()
		if err != nil {
			return err
		}
		defer db.Close()

		query := `DELETE FROM status WHERE id = $1`
		res, err := db.Exec(query, id)
		if err != nil {
			return err
		}

		affected, err := res.RowsAffected()
		if err != nil {
			return err
		}

		if affected == 0 {
			return echo.NewHTTPError(http.StatusNotFound, "Status not found")
		}

		b.NotifyClients("Status supprimé")

		return c.NoContent(http.StatusNoContent)
	}
}

// @Summary Récupérer un status par son ID
// @Description Récupère un status spécifique par son ID
// @Produce  json
// @Param id path int true "ID du status à récupérer"
// @Success 200 {object} model.Status
// @Router /status/{id} [get]
func GetStatusByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return err
	}

	db, err := config.ConnectToDB()
	if err != nil {
		return err
	}
	defer db.Close()

	d := new(model.Status)
	query := `SELECT id, designation, state FROM status WHERE id = $1`
	err = db.QueryRow(query, id).Scan(&d.ID, &d.Designation, &d.State)
	if err != nil {
		if err == sql.ErrNoRows {
			return echo.NewHTTPError(http.StatusNotFound, "Status not found")
		}
		return err
	}

	return c.JSON(http.StatusOK, d)
}

// @Summary Récupérer tous les status
// @Description Récupère la liste de tous les status
// @Produce  json
// @Success 200 {array} model.Status
// @Router /status [get]
func GetAllStatus(c echo.Context) error {
	db, err := config.ConnectToDB()
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT id, designation, state FROM status`)
	if err != nil {
		return err
	}
	defer rows.Close()

	statusList := make([]model.Status, 0)
	for rows.Next() {
		d := new(model.Status)
		err := rows.Scan(&d.ID, &d.Designation, &d.State)
		if err != nil {
			return err
		}
		statusList = append(statusList, *d)
	}

	if err = rows.Err(); err != nil {
		return err
	}

	return c.JSON(http.StatusOK, statusList)
}

// @Summary Provoque un crash du programme
// @Description Cette route provoque un crash de l'application avec un code de retour 1 après un délai de 3 secondes.
// @ID crash
// @Produce  json
// @Success 500 {object} model.CrashResponse
// @Router /crash [get]
func CrashApp(c echo.Context) error {
	response := &model.CrashResponse{
		Message: "Crash de l'appli dans 3 secondes",
	}

	go func() {
		time.Sleep(3 * time.Second)
		fmt.Println("Le programme va crasher avec un code de retour 1.")
		os.Exit(1)
	}()

	return c.JSON(http.StatusInternalServerError, response)
}
