package cmd

import (
	"context"
	"fmt"
	"errors"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/pkg/browser"
	"github.com/spf13/cobra"

	"github.com/mmclsntr/lineworks-cli/auth"
)

const TIMEOUT_SEC = 60

func getToken(profile string) (*auth.Token, error) {
	token := auth.Token{}

	t, err := token.ReadConfig(profile)
	if os.IsNotExist(err) {
		return nil, errors.New("profile does not exist.")
	} else if err != nil {
		return nil, err
	}

	return t, nil
}

// User Account Auth
func authUserAccount(profile string, clientCred *auth.ClientCredential) error {
	if clientCred.Scopes == "" {
		return errors.New("Scope does not set.\n")
	}
	ctx := context.Background()

	stateReq, _ := uuid.NewUUID()
	url := clientCred.AuthCodeURL(stateReq.String())
	fmt.Printf("Visit the URL for the auth dialog: %v\n", url)
	time.Sleep(1 * time.Second)
	browser.OpenURL(url)
	time.Sleep(1 * time.Second)

	auth.StartCallbackServer(ctx, clientCred.ListenAddr, clientCred.ListenPort, clientCred.RedirectPath, TIMEOUT_SEC, func(code string, state string) {
		if state != stateReq.String() {
			fmt.Errorf("'state' does not match")
		}

		// Get AccessToken
		tok := clientCred.GetAccessToken(code)
		tok.WriteConfig(profile)
	})

	return nil
}

// Service Account Auth
func authServiceAccount(profile string, clientCred *auth.ClientCredential, serviceAccount *auth.ServiceAccount) error {
	if clientCred.Scopes == "" {
		return errors.New("Scope does not set.\n")
	}

	tok := clientCred.GetAccessTokenJWT(*serviceAccount)
	tok.WriteConfig(profile)
	return nil
}

var authCmd = &cobra.Command{
	Use:   "auth",
	Short: "Authorization",
}

var authUserAccountCmd = &cobra.Command{
	Use:   "user-account",
	Short: "User Account Authorization",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		scopes, _ := cmd.Flags().GetString("scopes")
		addr, _ := cmd.Flags().GetString("addr")
		port, _ := cmd.Flags().GetString("port")
		path, _ := cmd.Flags().GetString("path")
		cred, err := getClientConfigure(profile)
		if err != nil {
			fmt.Printf("%s", err)
            return nil
		}

		if scopes != "" {
			cred.Scopes = scopes
		}
		if addr != "" {
			cred.ListenAddr = addr
		}
		if port != "" {
			cred.ListenPort = port
		}
		if path != "" {
			cred.RedirectPath = path
		}
		err = authUserAccount(profile, cred)
		if err != nil {
			fmt.Printf("%s", err)
            return nil
        }

        fmt.Printf("Success\n")
		return nil
	},
}

var authServiceAccountCmd = &cobra.Command{
	Use:   "service-account",
	Short: "Service Account Authorization",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		scopes, _ := cmd.Flags().GetString("scopes")
		cred, err := getClientConfigure(profile)
		if err != nil {
			fmt.Printf("%s", err)
            return nil
		}

		if scopes != "" {
			cred.Scopes = scopes
		}
		sa, err := getServiceAccountConfigure(profile)
		if err != nil {
			fmt.Printf("%s", err)
            return nil
		}
		err = authServiceAccount(profile, cred, sa)
		if err != nil {
			fmt.Printf("%s", err)
            return nil
		}

        fmt.Printf("Success\n")
		return nil
	},
}

var authGetAccessTokenCmd = &cobra.Command{
	Use:   "get-access-token",
	Short: "Get Access Token",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		token, err := getToken(profile)
		if err != nil {
			fmt.Printf("%s", err)
            return nil
		}
        fmt.Printf("%s", token.AccessToken)
		return nil
	},
}

var authGetScopesCmd = &cobra.Command{
	Use:   "get-scopes",
	Short: "Get scopes",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		token, err := getToken(profile)
		if err != nil {
			fmt.Printf("%s", err)
            return nil
		}
        fmt.Printf("%s", token.Scopes)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(authCmd)
	authCmd.AddCommand(authUserAccountCmd)
	authCmd.AddCommand(authServiceAccountCmd)
	authCmd.AddCommand(authGetAccessTokenCmd)
	authCmd.AddCommand(authGetScopesCmd)

	authCmd.PersistentFlags().StringP("profile", "", "", "Profile name")
	authCmd.MarkPersistentFlagRequired("profile")

	authCmd.PersistentFlags().StringP("scopes", "", "", "List scopes by comma-delimited format (ex. bot,user.read,board)")
	authUserAccountCmd.Flags().StringP("addr", "", "", "Listen address of callback server")
	authUserAccountCmd.Flags().StringP("port", "", "", "Listen port of callback server")
	authUserAccountCmd.Flags().StringP("path", "", "", "URL path of callback server")
}
