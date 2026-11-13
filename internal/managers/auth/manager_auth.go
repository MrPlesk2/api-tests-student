package auth

import (
	"encoding/json"
	"io"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"

	"api-tests-template/internal/client/http/auth"
	"api-tests-template/internal/managers/auth/models"
	"api-tests-template/internal/utils"
)

// Login выполняет авторизацию пользователя и возвращает токен с данными пользователя, ожидает 200
func Login(t *testing.T, username, password string) models.LoginOkResponse {
	t.Helper()
	require.NotEmpty(t, username, "Отсутствует логин пользователя")
	require.NotEmpty(t, password, "Отсутствует пароль пользователя")

	request := models.LoginRequest{
		Email:    username,
		Password: password,
	}

	loginRawResponse := auth.HttpPostAuthLogin(t, request)
	defer loginRawResponse.Body.Close()

	respArray, errResponse := io.ReadAll(loginRawResponse.Body)
	require.NoError(t, errResponse)

	loginResponse := models.LoginOkResponse{}
	err := json.Unmarshal(respArray, &loginResponse)
	require.NoError(t, err)

	return loginResponse
}

// RegisterOrLogin provisions an account once and reuses it on subsequent scenario runs.
func RegisterOrLogin(t *testing.T, request models.RegisterRequest) models.LoginOkResponse {
	t.Helper()
	response := auth.HttpPostAuthRegister(t, request)
	body := utils.ReadResponseBody(t, response)
	if response.StatusCode == http.StatusConflict {
		return Login(t, request.Email, request.Password)
	}
	require.Equal(t, http.StatusCreated, response.StatusCode, "Cannot register favorites account; response: %s", body)
	var result models.LoginOkResponse
	require.NoError(t, json.Unmarshal(body, &result), "Invalid registration JSON: %s", body)
	return result
}
