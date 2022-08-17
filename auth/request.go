package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const AuthURL = "https://auth.worksmobile.com/oauth2/v2.0/authorize"
const TokenURL = "https://auth.worksmobile.com/oauth2/v2.0/token"

type AccessTokenRequestBody struct {
	Code         string `json:"code"`
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Domain       string `json:"domain,omitempty"`
}

type AccessTokenJWTRequestBody struct {
	Assertion    string `json:"assertion"`
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Scopes       string `json:"scope"`
}

type AccessTokenResponseBody struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Scopes       string `json:"scope"`
	ExpiredIn    string `json:"expires_in"`
	TokenType    string `json:"token_type"`
}

type RefreshTokenRequestBody struct {
	RefreshToken string `json:"refresh_token"`
	GrantType    string `json:"grant_type"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type RefreshTokenResponseBody struct {
	AccessToken string `json:"access_token"`
	Scopes      string `json:"scope"`
	ExpiredIn   string `json:"expires_in"`
	TokenType   string `json:"token_type"`
}

// Get AccessToken
func RequestAccessToken(req_body AccessTokenRequestBody) (AccessTokenResponseBody, error) {
	req_body_json, _ := json.Marshal(req_body)

	return requestAccessToken(req_body_json)
}

// Get AccessToken (JWT)
func RequestAccessTokenJWT(req_body AccessTokenJWTRequestBody) (AccessTokenResponseBody, error) {
	req_body_json, _ := json.Marshal(req_body)

	return requestAccessToken(req_body_json)
}

func requestAccessToken(req_body_json []byte) (AccessTokenResponseBody, error) {
	//log.Printf("%s", req_body_json)

	res_body := AccessTokenResponseBody{}

	mapData := map[string]string{}
	if err := json.Unmarshal(req_body_json, &mapData); err != nil {
		return res_body, err
	}
	req_body_data := url.Values{}
	for k, v := range mapData {
		req_body_data.Add(k, v)
	}

	// Request
	res, _ := http.PostForm(TokenURL, req_body_data)
	defer res.Body.Close()

	//log.Printf("%v", res)
	body, _ := io.ReadAll(res.Body)

	if res.StatusCode != 200 {
		return res_body, errors.New(fmt.Sprintf("Error: status code %d, body %s", res.StatusCode, body))
	}

	if err := json.Unmarshal(body, &res_body); err != nil {
		return res_body, err
	}

	return res_body, nil
}

// Refresh AccessToken
//func RequestRefreshAccessToken(req_body RefreshTokenRequestBody) (RefreshTokenResponseBody, error) {
//}
