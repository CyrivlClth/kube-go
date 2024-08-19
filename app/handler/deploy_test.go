package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestDeploy_Reload(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	require.NoError(t, err)
	db = db.Debug()
	router := gin.Default()
	Register(router, db)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/-/reload", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "{\"data\":[\"cc.yaml\"]}", w.Body.String())
}

func TestDeploy_AddApp(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	require.NoError(t, err)
	db = db.Debug()
	router := gin.Default()
	Register(router, db)
	w := httptest.NewRecorder()
	var b bytes.Buffer
	s:=`{"name":"gatewa-service","maxCPUCount":2,"maxMemoryGB":2,"description":"A gateway service","preCmd":["java"],"args":null,"postCmd":["-jar","./app.jar"],"nodeSelector":{"t":"p"},"replicas":2}`
	b.WriteString(s)
	req, _ := http.NewRequest("POST", "/api/app-config", &b)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, fmt.Sprintf("{\"data\":%s}", s), w.Body.String())
}
