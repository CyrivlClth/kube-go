package handler

import (
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/CyrivlClth/kube-go/app/query"
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
	query.SetDefault(db)
	router := gin.Default()
	Register(router, db)
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/-/reload", nil)
	router.ServeHTTP(w, req)
	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, "{\"data\":[\"cc.yaml\",\"values.yaml\"]}", w.Body.String())
}

func TestDeploy_AddApp(t *testing.T) {
	db, err := gorm.Open(sqlite.Open(":memory:"))
	require.NoError(t, err)
	db = db.Debug()
	query.SetDefault(db)
	router := gin.Default()
	Register(router, db)
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/-/reload", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
	}

	s := `{"name":"gateway-service","maxCPUCount":2,"maxMemoryGB":2,"description":"A gateway service","preCmd":["java"],"args":null,"postCmd":["-jar","./app.jar"],"nodeSelector":{"t":"p"},"replicas":2}`

	{
		var b bytes.Buffer
		b.WriteString(s)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/app-config", &b)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"data\":%s}", s), w.Body.String())
	}
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/app-config", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"data\":[%s]}", s), w.Body.String())
	}
	ds := `{"appName":"gateway-service","envName":"values.yaml","image":"test","tag":"v1"}`
	{
		var b bytes.Buffer
		b.WriteString(ds)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/app-deploy", &b)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"data\":%s}", ds), w.Body.String())
	}
	ds = `{"appName":"gateway-service","envName":"values.yaml","image":"test","tag":"v2"}`
	{
		var b bytes.Buffer
		b.WriteString(ds)
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("POST", "/api/app-deploy", &b)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, fmt.Sprintf("{\"data\":%s}", ds), w.Body.String())
	}
	{
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/api/app-config", nil)
		router.ServeHTTP(w, req)
		assert.Equal(t, 200, w.Code)
		assert.JSONEq(t, `{"data":
		[{
		"name":"gateway-service",
		"maxCPUCount":2,
		"maxMemoryGB":2,
		"description":"A gateway service",
		"preCmd":["java"],
		"args":null,
		"postCmd":["-jar","./app.jar"],
		"nodeSelector":{"t":"p"},
		"replicas":2,
		"deploy":[{"appName":"gateway-service","envName":"values.yaml","image":"test","tag":"v2"}]
		}]}`, w.Body.String())
	}
}
