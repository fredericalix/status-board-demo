package api_test

import (
	"backend-server/api"
	"backend-server/model"
	"backend-server/sse"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
)

var createdStatusID int



func TestCreateStatus(t *testing.T) {
	e := echo.New()
	status := model.Status{
		Designation: "test",
		State:       "green",
	}
	reqBody, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPost, "/status", bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	b := sse.NewSSEBroker()
	go b.Start()

	err := api.CreateStatus(b)(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, rec.Code)

	var createdStatus model.Status
	json.Unmarshal(rec.Body.Bytes(), &createdStatus)
	assert.Equal(t, "test", createdStatus.Designation)
	assert.Equal(t, "green", createdStatus.State)

	createdStatusID = createdStatus.ID
}

func TestGetStatusByID(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/status/%d", createdStatusID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/status/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(createdStatusID))

	err := api.GetStatusByID(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var status model.Status
	json.Unmarshal(rec.Body.Bytes(), &status)
	assert.Equal(t, createdStatusID, status.ID)
}

func TestUpdateStatus(t *testing.T) {
	e := echo.New()
	status := model.Status{
		Designation: "testmaj",
		State:       "red",
	}
	reqBody, _ := json.Marshal(status)
	req := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/status/%d", createdStatusID), bytes.NewReader(reqBody))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/status/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(createdStatusID))

	b := sse.NewSSEBroker()
	go b.Start()

	err := api.UpdateStatus(b)(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, rec.Code)

	var updatedStatus model.Status
	json.Unmarshal(rec.Body.Bytes(), &updatedStatus)
	assert.Equal(t, "testmaj", updatedStatus.Designation)
	assert.Equal(t, "red", updatedStatus.State)
}

func TestDeleteStatus(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/status/%d", createdStatusID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetPath("/status/:id")
	c.SetParamNames("id")
	c.SetParamValues(strconv.Itoa(createdStatusID))

	b := sse.NewSSEBroker()
	go b.Start()

	err := api.DeleteStatus(b)(c)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, rec.Code)
}
