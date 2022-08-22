package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"

	"github.com/mmclsntr/lineworks-cli/auth"
)

const DEFAULT_ADDR = "127.0.0.1"
const DEFAULT_PORT = "9876"
const DEFAULT_PATH = "/oauth/callback"

func getClientConfigure(profile string) (*auth.ClientCredential, error) {
	cred := auth.ClientCredential{}

	c, err := cred.ReadConfig(profile)
	if os.IsNotExist(err) {
		return nil, errors.New("profile does not exist.")
	} else if err != nil {
		return nil, err
	}
	return c, nil
}

func setClientConfigure(profile string, clientId string, clientSecret string, scopes string, redirectUrl string, addr string, port string, path string, domainId string) error {
	cred := auth.ClientCredential{
		ClientID:     clientId,
		ClientSecret: clientSecret,
		Scopes:       scopes,
		ListenAddr:   addr,
		ListenPort:   port,
		RedirectPath: path,
		DomainID:     domainId,
	}

	err := cred.WriteConfig(profile)
	if err != nil {
		return err
	}
	return nil
}

func getServiceAccountConfigure(profile string) (*auth.ServiceAccount, error) {
	sa := auth.ServiceAccount{}

	s, err := sa.ReadConfig(profile)
	if os.IsNotExist(err) {
		return nil, errors.New("profile does not exist.")
	} else if err != nil {
		return nil, err
	}

	return s, nil
}

func setServiceAccountConfigure(profile string, serviceAccountId string, privateKeyFile string) error {
	privateKeyData, err := ioutil.ReadFile(privateKeyFile)
	if err != nil {
		return err
	}

	sa := auth.ServiceAccount{
		ServiceAccountID: serviceAccountId,
		PrivateKey:       string(privateKeyData),
	}

	err = sa.WriteConfig(profile)
	if err != nil {
		return err
	}
	return nil
}

var configureCmd = &cobra.Command{
	Use:   "configure",
	Short: "Configure authorization settings for access token.",
}

var configureGetClientCmd = &cobra.Command{
	Use:   "get-client",
	Short: "Get client credentials.",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		cred, err := getClientConfigure(profile)
		if err != nil {
			fmt.Printf("%s", err)
			return nil
		}

		b, err := json.MarshalIndent(cred, "", "    ")
		if err != nil {
			fmt.Printf("%s", err)
			return nil
		}
		fmt.Printf("%s\n", b)
		return nil
	},
}

var configureSetClientCmd = &cobra.Command{
	Use:   "set-client",
	Short: "Set client credentials.",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		client_id, _ := cmd.Flags().GetString("client-id")
		client_secret, _ := cmd.Flags().GetString("client-secret")
		scopes, _ := cmd.Flags().GetString("scopes")
		addr, _ := cmd.Flags().GetString("addr")
		port, _ := cmd.Flags().GetString("port")
		path, _ := cmd.Flags().GetString("path")
		domain_id, _ := cmd.Flags().GetString("domain-id")

		redirect_url := fmt.Sprintf("http://%s:%s%s", addr, port, path)
		setClientConfigure(profile, client_id, client_secret, scopes, redirect_url, addr, port, path, domain_id)

		// View
		cred, err := getClientConfigure(profile)
		if err != nil {
			fmt.Printf("%s", err)
			return nil
		}

		b, err := json.MarshalIndent(cred, "", "    ")
		if err != nil {
			fmt.Printf("%s", err)
			return nil
		}
		fmt.Printf("%s\n", b)
		return nil
	},
}

var configureGetRedirectUrlCmd = &cobra.Command{
	Use:   "get-redirect-url",
	Short: "Get redirect url.",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		cred, err := getClientConfigure(profile)
		if err != nil {
			fmt.Printf("%s", err)
			return nil
		}
		fmt.Printf("%s\n", cred.GetRedirectUrl())
		return nil
	},
}

var configureGetServiceAccountCmd = &cobra.Command{
	Use:   "get-service-account",
	Short: "Get service account settings.",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")

		sa, err := getServiceAccountConfigure(profile)
		if err != nil {
			fmt.Printf("%s", err)
			return nil
		}
		b, err := json.MarshalIndent(sa, "", "    ")
		if err != nil {
			fmt.Printf("%s", err)
			return nil
		}
		fmt.Printf("%s\n", b)
		return nil
	},
}

var configureSetServiceAccountCmd = &cobra.Command{
	Use:   "set-service-account",
	Short: "Set service account settings.",
	RunE: func(cmd *cobra.Command, args []string) error {
		profile, _ := cmd.Flags().GetString("profile")
		serviceAccountId, _ := cmd.Flags().GetString("service-account-id")
		privateKeyFile, _ := cmd.Flags().GetString("private-key-file")

		setServiceAccountConfigure(profile, serviceAccountId, privateKeyFile)

		// View
		sa, err := getServiceAccountConfigure(profile)
		if err != nil {
			fmt.Printf("%s", err)
			return nil
		}
		b, err := json.MarshalIndent(sa, "", "    ")
		if err != nil {
			fmt.Printf("%s", err)
			return nil
		}
		fmt.Printf("%s\n", b)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(configureCmd)
	configureCmd.AddCommand(configureGetClientCmd)
	configureCmd.AddCommand(configureSetClientCmd)
	configureCmd.AddCommand(configureGetRedirectUrlCmd)
	configureCmd.AddCommand(configureGetServiceAccountCmd)
	configureCmd.AddCommand(configureSetServiceAccountCmd)

	configureCmd.PersistentFlags().StringP("profile", "", "", "Profile name")
	configureCmd.MarkPersistentFlagRequired("profile")

	configureSetClientCmd.Flags().StringP("client-id", "", "", "Client ID")
	configureSetClientCmd.MarkFlagRequired("client-id")
	configureSetClientCmd.Flags().StringP("client-secret", "", "", "Client Secret")
	configureSetClientCmd.MarkFlagRequired("client-secret")
	configureSetClientCmd.Flags().StringP("scopes", "", "", "Scopes. Must be comma-delimited format (ex. bot,user.read,board)")
	configureSetClientCmd.Flags().StringP("addr", "", DEFAULT_ADDR, "Listening address of callback server")
	configureSetClientCmd.Flags().StringP("port", "", DEFAULT_PORT, "Listening port of callback server")
	configureSetClientCmd.Flags().StringP("path", "", DEFAULT_PATH, "URL path of callback server")
	configureSetClientCmd.Flags().StringP("domain-id", "", "", "Domain ID")

	configureSetServiceAccountCmd.Flags().StringP("service-account-id", "", "", "Service Account ID")
	configureSetClientCmd.MarkFlagRequired("service-account-id")
	configureSetServiceAccountCmd.Flags().StringP("private-key-file", "", "", "Private Key file path")
	configureSetClientCmd.MarkFlagRequired("private-key-file")
}
