package app_engine

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestHttpApi(t *testing.T) {
	r := gin.Default()
	RegisterHttpApiRoutes(r)
	ts := httptest.NewServer(r)
	defer ts.Close()
	req, _ := http.NewRequest(http.MethodGet, ts.URL+"/pods?namespace=hello", nil)
	resp, err := http.DefaultClient.Do(req)
	assert.NoError(t, err)
	body, _ := io.ReadAll(resp.Body)
	t.Log(string(body))
}

func TestCreateA(t *testing.T) {
	err := CreateApp(context.Background(), &AbstractApp{
		Namespace: "dev",
		Name:      "test-nginx",
		Containers: []*AbstractContainer{
			{
				Name:  "nginx",
				Image: "nginx:1.14.2",
			},
		},
	})
	assert.NoError(t, err)
}
