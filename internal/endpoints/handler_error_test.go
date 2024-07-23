package endpoints

import (
	internalerrors "emailn/internal/internal-errors"
	"encoding/json"
	"errors"
	"testing"

	"net/http"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func TestHanlderErrorWhenEndpointRetunrsInternalError(t *testing.T) {
	assert := assert.New(t)

	endpoints := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		return nil, http.StatusInternalServerError, internalerrors.ErrInternal
	}

	handlerFunc := HandlerError(endpoints)
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	handlerFunc.ServeHTTP(res, req)

	assert.Equal(http.StatusInternalServerError, res.Code)
}

func TestHanlderErrorWhenEndpointRetunrsDomainError(t *testing.T) {
	assert := assert.New(t)

	endpoints := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		return nil, http.StatusBadRequest, errors.New("Domain Error")
	}

	handlerFunc := HandlerError(endpoints)
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	handlerFunc.ServeHTTP(res, req)

	assert.Equal(http.StatusBadRequest, res.Code)
}

func TestHanlderErrorWhenEndpointRetunrsObjectAndStatus(t *testing.T) {
	assert := assert.New(t)

	type bodyForTest struct {
		Id int
	}
	objExpected := bodyForTest{Id: 2}
	endpoints := func(w http.ResponseWriter, r *http.Request) (interface{}, int, error) {
		return objExpected, http.StatusCreated, nil
	}

	handlerFunc := HandlerError(endpoints)
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	handlerFunc.ServeHTTP(res, req)

	objReturned := bodyForTest{}
	json.Unmarshal(res.Body.Bytes(), &objReturned)

	assert.Equal(http.StatusCreated, res.Code)
	assert.Equal(objExpected, objReturned)
}
