package auth

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/golang-jwt/jwt/v4"
	"io/ioutil"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

type ClientCredential struct {
	ClientID     string `toml:"client_id" json:"client_id"`
	ClientSecret string `toml:"client_secret" json:"client_secret"`
	Scopes       string `toml:"scopes" json:"scopes"`
	ListenAddr   string `toml:"addr" json:"addr"`
	ListenPort   string `toml:"port" json:"port"`
	RedirectPath string `toml:"path" json:"path"`
	DomainID     string `toml:"domain_id,omitempty" json:"domain_id,omitempty"`
}

type ServiceAccount struct {
	ServiceAccountID string `toml:"service_account_id" json:"service_account_id"`
	PrivateKey       string `toml:"private_key" json:"private_key"`
}

type Token struct {
	AccessToken  string `toml:"access_token" json:"access_token"`
	RefreshToken string `toml:"refresh_token" json:"refresh_token"`
	Scopes       string `toml:"scopes" json:"scopes"`
	ExpiredIn    string `toml:"expired_in" json:"expired_in"`
}

const CONFIG_DIR_NAME = ".config"
const CONFIG_SERVICE_DIR_NAME = "lineworks"
const CONFIG_OAUTH_FILE_NAME = "oauth.toml"
const CONFIG_SERVICE_ACCOUNT_FILE_NAME = "service_account.toml"
const CONFIG_TOKEN_FILE_NAME = "token.toml"
const CONFIG_PATH_ENV_NAME = "LINEWORKS_CONFIG_DIR"

// Generate Authorization Code URL
func (cred *ClientCredential) AuthCodeURL(state string) string {
	u, err := url.Parse(AuthURL)
	if err != nil {
		log.Fatal(err)
	}
	q := u.Query()
	// Client ID
	q.Set("client_id", cred.ClientID)
	// Scope
	q.Set("scope", cred.Scopes)
	// RedirectURL
	q.Set("redirect_uri", cred.GetRedirectUrl())
	// Response Type
	q.Set("response_type", "code")
	// Response Type
	q.Set("state", state)
	// Domain ID
	if cred.DomainID != "" {
		q.Set("domain", cred.DomainID)
	}

	u.RawQuery = q.Encode()

	return u.String()
}

// Get access token
func (cred *ClientCredential) GetAccessToken(code string) Token {
	// Create request body
	req := AccessTokenRequestBody{
		Code:         code,
		GrantType:    "authorization_code",
		ClientID:     cred.ClientID,
		ClientSecret: cred.ClientSecret,
		Domain:       cred.DomainID,
	}

	// Request
	res_body, err := RequestAccessToken(req)
	if err != nil {
		log.Fatal(err)
	}

	return Token{
		AccessToken:  res_body.AccessToken,
		RefreshToken: res_body.RefreshToken,
		Scopes:       res_body.Scopes,
		ExpiredIn:    res_body.ExpiredIn,
	}
}

// Get access token (JWT)
func (cred *ClientCredential) GetAccessTokenJWT(sva ServiceAccount) Token {
	// JWT
	jwt := GenerateJWT(cred.ClientID, sva.ServiceAccountID, sva.PrivateKey)

	// Create request body
	req := AccessTokenJWTRequestBody{
		Assertion:    jwt,
		GrantType:    "urn:ietf:params:oauth:grant-type:jwt-bearer",
		ClientID:     cred.ClientID,
		ClientSecret: cred.ClientSecret,
		Scopes:       cred.Scopes,
	}

	// Request
	res_body, err := RequestAccessTokenJWT(req)
	if err != nil {
		log.Fatal(err)
	}

	return Token{
		AccessToken:  res_body.AccessToken,
		RefreshToken: res_body.RefreshToken,
		Scopes:       res_body.Scopes,
		ExpiredIn:    res_body.ExpiredIn,
	}
}

// Refresh access token
//func (conf Config) RefreshAccessToken(token Token) Token {
//}

// Generate JWT
func GenerateJWT(clientId string, serviceAccountId string, privateKey string) string {
	currentTime := time.Now()
	// Claims object
	claims := jwt.MapClaims{
		"iss": clientId,
		"sub": serviceAccountId,
		"iat": currentTime.Unix(),
		"exp": currentTime.Add(time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims) // RSA SHA-256

	// Key
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		log.Fatal(err)
	}

	// Sign
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Fatal(err)
	}
	return tokenString
}

// Get Redirect URL
func (cred *ClientCredential) GetRedirectUrl() string {
	return fmt.Sprintf("http://%s:%s%s", cred.ListenAddr, cred.ListenPort, cred.RedirectPath)
}

func getConfigBasePath() string {
	if configPath := os.Getenv(CONFIG_PATH_ENV_NAME); configPath != "" {
		return configPath
	}

	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return ""
	}
	configDir := filepath.Join(homeDir, CONFIG_DIR_NAME, CONFIG_SERVICE_DIR_NAME)
	return configDir
}

func getConfigProfileDir(profile string) string {
	configDir := getConfigBasePath()
	configDir = filepath.Join(configDir, profile)

	return configDir
}

func makeConfigProfileDir(profile string) error {
	configDir := getConfigProfileDir(profile)
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return err
	}
	return nil
}

func getConfigFileName(profile string, file_name string) string {
	configDir := getConfigProfileDir(profile)
	configFile := filepath.Join(configDir, file_name)

	return configFile
}

func ListConfigProfiles() []string {
	configDir := getConfigBasePath()
	profiles := []string{}

	files, _ := ioutil.ReadDir(configDir)
	for _, f := range files {
		if f.IsDir() {
			profiles = append(profiles, f.Name())
		}
	}
	return profiles
}

func (cred ClientCredential) ReadConfig(profile string) (*ClientCredential, error) {
	configFile := getConfigFileName(profile, CONFIG_OAUTH_FILE_NAME)
	_, err := os.Stat(configFile)
	if err != nil {
		return nil, err
	}
	fp, err := os.Open(configFile)
	defer fp.Close()
	if err != nil {
		return nil, err
	}

	newCred := ClientCredential{}
	_, err = toml.NewDecoder(fp).Decode(&newCred)
	return &newCred, err
}

func (cred *ClientCredential) WriteConfig(profile string) error {
	err := makeConfigProfileDir(profile)
	if err != nil {
		return err
	}
	configFile := getConfigFileName(profile, CONFIG_OAUTH_FILE_NAME)
	fp, err := os.Create(configFile)
	defer fp.Close()
	if err != nil {
		return err
	}

	err = toml.NewEncoder(fp).Encode(cred)
	return err
}

func (sa ServiceAccount) ReadConfig(profile string) (*ServiceAccount, error) {
	configFile := getConfigFileName(profile, CONFIG_SERVICE_ACCOUNT_FILE_NAME)
	_, err := os.Stat(configFile)
	if err != nil {
		return nil, err
	}
	fp, err := os.Open(configFile)
	defer fp.Close()
	if err != nil {
		return nil, err
	}

	newSa := ServiceAccount{}
	_, err = toml.NewDecoder(fp).Decode(&newSa)
	return &newSa, err
}

func (sa *ServiceAccount) WriteConfig(profile string) error {
	err := makeConfigProfileDir(profile)
	if err != nil {
		return err
	}
	configFile := getConfigFileName(profile, CONFIG_SERVICE_ACCOUNT_FILE_NAME)
	fp, err := os.Create(configFile)
	defer fp.Close()
	if err != nil {
		return err
	}

	err = toml.NewEncoder(fp).Encode(sa)
	return err
}

func (token Token) ReadConfig(profile string) (*Token, error) {
	configFile := getConfigFileName(profile, CONFIG_TOKEN_FILE_NAME)
	_, err := os.Stat(configFile)
	if err != nil {
		return nil, err
	}
	fp, err := os.Open(configFile)
	defer fp.Close()
	if err != nil {
		return nil, err
	}

	newToken := Token{}
	_, err = toml.NewDecoder(fp).Decode(&newToken)
	return &newToken, err
}

func (token *Token) WriteConfig(profile string) error {
	err := makeConfigProfileDir(profile)
	if err != nil {
		return err
	}
	configFile := getConfigFileName(profile, CONFIG_TOKEN_FILE_NAME)
	fp, err := os.Create(configFile)
	defer fp.Close()
	if err != nil {
		return err
	}

	err = toml.NewEncoder(fp).Encode(token)
	return err
}
