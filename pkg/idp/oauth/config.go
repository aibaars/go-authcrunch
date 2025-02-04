// Copyright 2022 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oauth

import (
	"fmt"
	"github.com/greenpau/go-authcrunch/pkg/errors"
	"net/url"
	"regexp"
	"strings"
)

// Config holds the configuration for the IdentityProvider.
type Config struct {
	Name              string `json:"name,omitempty" xml:"name,omitempty" yaml:"name,omitempty"`
	Realm             string `json:"realm,omitempty" xml:"realm,omitempty" yaml:"realm,omitempty"`
	Driver            string `json:"driver,omitempty" xml:"driver,omitempty" yaml:"driver,omitempty"`
	DomainName        string `json:"domain_name,omitempty" xml:"domain_name,omitempty" yaml:"domain_name,omitempty"`
	ClientID          string `json:"client_id,omitempty" xml:"client_id,omitempty" yaml:"client_id,omitempty"`
	ClientSecret      string `json:"client_secret,omitempty" xml:"client_secret,omitempty" yaml:"client_secret,omitempty"`
	ServerID          string `json:"server_id,omitempty" xml:"server_id,omitempty" yaml:"server_id,omitempty"`
	ServerName        string `json:"server_name,omitempty" xml:"server_name,omitempty" yaml:"server_name,omitempty"`
	AppSecret         string `json:"app_secret,omitempty" xml:"app_secret,omitempty" yaml:"app_secret,omitempty"`
	TenantID          string `json:"tenant_id,omitempty" xml:"tenant_id,omitempty" yaml:"tenant_id,omitempty"`
	IdentityTokenName string `json:"identity_token_name,omitempty" xml:"identity_token_name,omitempty" yaml:"identity_token_name,omitempty"`

	Scopes []string `json:"scopes,omitempty" xml:"scopes,omitempty" yaml:"scopes,omitempty"`

	// The number if seconds to wait before getting key material
	// from an OAuth 2.0 identity provider.
	DelayStart int `json:"delay_start,omitempty" xml:"delay_start,omitempty" yaml:"delay_start,omitempty"`
	// The number of the retry attempts getting key material
	// from an OAuth 2.0 identity provider.
	RetryAttempts int `json:"retry_attempts,omitempty" xml:"retry_attempts,omitempty" yaml:"retry_attempts,omitempty"`
	// The number of seconds to wait until the retrying.
	RetryInterval int `json:"retry_interval,omitempty" xml:"retry_interval,omitempty" yaml:"retry_interval,omitempty"`

	UserRoleMapList []map[string]interface{} `json:"user_roles,omitempty" xml:"user_roles,omitempty" yaml:"user_roles,omitempty"`

	// The URL to OAuth 2.0 Custom Authorization Server.
	BaseAuthURL string `json:"base_auth_url,omitempty" xml:"base_auth_url,omitempty" yaml:"base_auth_url,omitempty"`

	// The URL to OAuth 2.0 metadata related to your Custom Authorization Server.
	MetadataURL string `json:"metadata_url,omitempty" xml:"metadata_url,omitempty" yaml:"metadata_url,omitempty"`

	// The regex filters for user groups extracted via IdP API.
	UserGroupFilters []string `json:"user_group_filters,omitempty" xml:"user_group_filters,omitempty" yaml:"user_group_filters,omitempty"`
	// The regex filters for user orgs extracted via IdP API.
	UserOrgFilters []string `json:"user_org_filters,omitempty" xml:"user_org_filters,omitempty" yaml:"user_org_filters,omitempty"`

	// Disables metadata discovery via public metadata URL.
	MetadataDiscoveryDisabled bool `json:"metadata_discovery_disabled,omitempty" xml:"metadata_discovery_disabled,omitempty" yaml:"metadata_discovery_disabled,omitempty"`

	KeyVerificationDisabled bool `json:"key_verification_disabled,omitempty" xml:"key_verification_disabled,omitempty" yaml:"key_verification_disabled,omitempty"`
	PassGrantTypeDisabled   bool `json:"pass_grant_type_disabled,omitempty" xml:"pass_grant_type_disabled,omitempty" yaml:"pass_grant_type_disabled,omitempty"`
	ResponseTypeDisabled    bool `json:"response_type_disabled,omitempty" xml:"response_type_disabled,omitempty" yaml:"response_type_disabled,omitempty"`
	NonceDisabled           bool `json:"nonce_disabled,omitempty" xml:"nonce_disabled,omitempty" yaml:"nonce_disabled,omitempty"`
	ScopeDisabled           bool `json:"scope_disabled,omitempty" xml:"scope_disabled,omitempty" yaml:"scope_disabled,omitempty"`

	AcceptHeaderEnabled bool `json:"accept_header_enabled,omitempty" xml:"accept_header_enabled,omitempty" yaml:"accept_header_enabled,omitempty"`

	JsCallbackEnabled bool `json:"js_callback_enabled,omitempty" xml:"js_callback_enabled,omitempty" yaml:"js_callback_enabled,omitempty"`

	ResponseType []string `json:"response_type,omitempty" xml:"response_type,omitempty" yaml:"response_type,omitempty"`

	AuthorizationURL string `json:"authorization_url,omitempty" xml:"authorization_url,omitempty" yaml:"authorization_url,omitempty"`
	TokenURL         string `json:"token_url,omitempty" xml:"token_url,omitempty" yaml:"token_url,omitempty"`

	RequiredTokenFields []string `json:"required_token_fields,omitempty" xml:"required_token_fields,omitempty" yaml:"required_token_fields,omitempty"`

	TLSInsecureSkipVerify bool `json:"tls_insecure_skip_verify,omitempty" xml:"tls_insecure_skip_verify,omitempty" yaml:"tls_insecure_skip_verify,omitempty"`

	// The predefined public RSA based JWKS keys.
	JwksKeys map[string]string `json:"jwks_keys,omitempty" xml:"jwks_keys,omitempty" yaml:"jwks_keys,omitempty"`

	// Disables the check for the presence of email field in a token.
	EmailClaimCheckDisabled bool `json:"email_claim_check_disabled,omitempty" xml:"email_claim_check_disabled,omitempty" yaml:"email_claim_check_disabled,omitempty"`

	IconName  string `json:"icon_name,omitempty" xml:"icon_name,omitempty" yaml:"icon_name,omitempty"`
	IconText  string `json:"icon_text,omitempty" xml:"icon_text,omitempty" yaml:"icon_text,omitempty"`
	IconColor string `json:"icon_color,omitempty" xml:"icon_color,omitempty" yaml:"icon_color,omitempty"`

	UserInfoFields []string `json:"user_info_fields,omitempty" xml:"user_info_fields,omitempty" yaml:"user_info_fields,omitempty"`
}

// Validate validates identity store configuration.
func (cfg *Config) Validate() error {
	if cfg.Name == "" {
		return errors.ErrIdentityProviderConfigureNameEmpty
	}

	if cfg.Realm == "" {
		return errors.ErrIdentityProviderConfigureRealmEmpty
	}

	if cfg.ClientID == "" {
		return errors.ErrIdentityProviderConfig.WithArgs("client id not found")
	}

	if cfg.ClientSecret == "" {
		return errors.ErrIdentityProviderConfig.WithArgs("client secret not found")
	}

	if cfg.DelayStart > 0 {
		if cfg.RetryAttempts < 1 {
			cfg.RetryAttempts = 2
		}
		if cfg.RetryInterval == 0 {
			cfg.RetryInterval = cfg.DelayStart
		}
	}

	if cfg.RetryAttempts > 0 && cfg.DelayStart == 0 {
		if cfg.RetryInterval == 0 {
			cfg.RetryInterval = 5
		}
	}

	if len(cfg.Scopes) < 1 {
		switch cfg.Driver {
		case "facebook":
			cfg.Scopes = []string{
				// "public_profile",
				"email",
			}
		case "github":
			cfg.Scopes = []string{"read:user"}
		case "nextcloud":
			cfg.Scopes = []string{"email"}
		case "google":
			cfg.Scopes = []string{"openid", "email", "profile"}
		default:
			cfg.Scopes = []string{"openid", "email", "profile"}
		}
	}

	switch cfg.IdentityTokenName {
	case "":
		cfg.IdentityTokenName = "id_token"
	case "id_token", "access_token":
	default:
		return errors.ErrIdentityProviderConfig.WithArgs(
			fmt.Errorf("identity token name %q is unsupported", cfg.IdentityTokenName),
		)
	}

	switch cfg.Driver {
	case "okta":
		if cfg.ServerID == "" {
			return errors.ErrIdentityProviderConfig.WithArgs("server id not found")
		}
		if cfg.DomainName == "" {
			return errors.ErrIdentityProviderConfig.WithArgs("domain name not found")
		}
		if cfg.BaseAuthURL == "" {
			cfg.BaseAuthURL = fmt.Sprintf(
				"https://%s/oauth2/%s/",
				cfg.DomainName, cfg.ServerID,
			)
			cfg.MetadataURL = cfg.BaseAuthURL + ".well-known/openid-configuration?client_id=" + cfg.ClientID
		}
	case "google":
		if cfg.BaseAuthURL == "" {
			cfg.BaseAuthURL = "https://accounts.google.com/o/oauth2/v2/"
			cfg.MetadataURL = "https://accounts.google.com/.well-known/openid-configuration"
		}
		// If Google client_id does not contains domain name, append with
		// the default of .apps.googleusercontent.com.
		if !strings.Contains(cfg.ClientID, ".") {
			cfg.ClientID = cfg.ClientID + ".apps.googleusercontent.com"
		}
	case "github":
		if cfg.BaseAuthURL == "" {
			cfg.BaseAuthURL = "https://github.com/login/oauth/"
		}
		cfg.RequiredTokenFields = []string{"access_token"}
		cfg.AuthorizationURL = "https://github.com/login/oauth/authorize"
		cfg.TokenURL = "https://github.com/login/oauth/access_token"
	case "gitlab":
		if cfg.DomainName == "" {
			cfg.DomainName = "gitlab.com"
		}
		if cfg.BaseAuthURL == "" {
			cfg.BaseAuthURL = fmt.Sprintf("https://%s/", cfg.DomainName)
			cfg.MetadataURL = cfg.BaseAuthURL + ".well-known/openid-configuration"
		}
	case "azure":
		if cfg.TenantID == "" {
			cfg.TenantID = "common"
		}
		if cfg.BaseAuthURL == "" {
			cfg.BaseAuthURL = "https://login.microsoftonline.com/" + cfg.TenantID + "/oauth2/v2.0/"
			cfg.MetadataURL = "https://login.microsoftonline.com/" + cfg.TenantID + "/v2.0/.well-known/openid-configuration"
		}
	case "facebook":
		if cfg.BaseAuthURL == "" {
			cfg.BaseAuthURL = "https://www.facebook.com/v12.0/dialog/"
		}
		cfg.RequiredTokenFields = []string{"access_token"}
		cfg.AuthorizationURL = "https://www.facebook.com/v12.0/dialog/oauth"
		cfg.TokenURL = "https://graph.facebook.com/v12.0/oauth/access_token"
	case "nextcloud":
		cfg.AuthorizationURL = fmt.Sprintf("%s/apps/oauth2/authorize", cfg.BaseAuthURL)
		cfg.TokenURL = fmt.Sprintf("%s/apps/oauth2/api/v1/token", cfg.BaseAuthURL)
	case "generic":
	case "":
		return errors.ErrIdentityProviderConfig.WithArgs("driver name not found")
	default:
		return errors.ErrIdentityProviderConfig.WithArgs(
			fmt.Errorf("driver %q is unsupported", cfg.Driver),
		)
	}

	if len(cfg.RequiredTokenFields) < 1 {
		cfg.RequiredTokenFields = []string{"access_token", "id_token"}
	}

	if cfg.BaseAuthURL == "" {
		if cfg.MetadataURL == "" {
			return errors.ErrIdentityProviderConfig.WithArgs("base authentication url not found")
		}
	}

	// Validate metadata URL, i.e. endpoint discovery.
	switch cfg.Driver {
	case "github":
	case "facebook":
	case "nextcloud":
	default:
		if len(cfg.JwksKeys) > 0 && cfg.AuthorizationURL != "" && cfg.TokenURL != "" {
			for kid, fp := range cfg.JwksKeys {
				if _, err := NewJwksKeyFromRSAPublicKeyPEM(kid, fp); err != nil {
					return errors.ErrIdentityProviderConfig.WithArgs(
						fmt.Errorf("failed loading kid %q: %v", kid, err),
					)
				}
			}
		} else {
			if cfg.MetadataURL == "" {
				return errors.ErrIdentityProviderConfig.WithArgs("metadata url not found")
			}
		}
	}

	parsedBaseAuthURL, err := url.Parse(cfg.BaseAuthURL)
	if err != nil {
		return errors.ErrIdentityProviderConfig.WithArgs(
			fmt.Errorf("failed to parse base auth url %q: %v", cfg.BaseAuthURL, err),
		)
	}
	cfg.ServerName = parsedBaseAuthURL.Host

	if len(cfg.ResponseType) < 1 {
		cfg.ResponseType = []string{"code"}
	}

	// Configure user group filters, if any.
	for _, pattern := range cfg.UserGroupFilters {
		if _, err := regexp.Compile(pattern); err != nil {
			return errors.ErrIdentityProviderConfig.WithArgs(
				fmt.Errorf("invalid user group pattern %q: %v", pattern, err),
			)
		}
	}

	// Configure user org filters, if any.
	for _, pattern := range cfg.UserOrgFilters {
		if _, err := regexp.Compile(pattern); err != nil {
			return errors.ErrIdentityProviderConfig.WithArgs(
				fmt.Errorf("invalid user org pattern %q: %v", pattern, err),
			)
		}
	}

	// Configure UI Icons.
	if cfg.IconName == "" {
		switch cfg.Driver {
		case "google":
			cfg.IconName = "google"
		case "facebook":
			cfg.IconName = "facebook"
		case "linkedin":
			cfg.IconName = "linkedin"
		case "github":
			cfg.IconName = "github"
		case "gitlab":
			cfg.IconName = "gitlab"
		case "windows":
			cfg.IconName = "windows"
		case "azure":
			cfg.IconName = "windows"
		case "aws":
			cfg.IconName = "aws"
		case "amazon":
			cfg.IconName = "amazon"
		default:
			cfg.IconName = "codepen"
		}
	}

	if cfg.IconText == "" {
		switch cfg.Driver {
		case "google":
			cfg.IconText = "Google"
		case "facebook":
			cfg.IconText = "Facebook"
		case "linkedin":
			cfg.IconText = "LinkedIn"
		case "github":
			cfg.IconText = "Github"
		case "windows":
			cfg.IconText = "Microsoft"
		case "azure":
			cfg.IconText = "Azure"
		case "aws":
			cfg.IconText = "AWS"
		case "amazon":
			cfg.IconText = "Amazon"
		default:
			cfg.IconText = cfg.Realm
		}
	}

	if cfg.IconColor == "" {
		switch cfg.Driver {
		case "google":
			cfg.IconColor = "red darken-1"
		case "facebook":
			cfg.IconColor = "blue darken-4"
		case "linkedin":
			cfg.IconColor = "blue darken-1"
		case "github":
			cfg.IconColor = "grey darken-3"
		case "windows":
			cfg.IconColor = "orange darken-1"
		case "azure":
			cfg.IconColor = "blue"
		case "aws", "amazon":
			cfg.IconColor = "blue-grey darken-2"
		default:
			cfg.IconColor = "grey darken-3"
		}
	}

	return nil
}
