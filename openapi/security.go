package openapi

import (
	"slices"

	openapi "github.com/getkin/kin-openapi/openapi3"
)

// OAuthAuthorizePath is the host's OAuth2 authorization endpoint advertised by
// the oauth2 security scheme. It must equal host-manager's
// auth.OAuthAuthorizePath, the path the host actually serves.
const OAuthAuthorizePath = "/-/oauth/authorize"

// TokenPath is the host's OAuth2 token endpoint, relative to the host issuer.
// Functions also POST RFC 8693 token exchanges to it (FUNCTION_CONTRACT.md §8).
const TokenPath = "/-/token"

const schemeTypeAPIKey = "apiKey"

// AllPathTypes is every route type, in the order of KDexHost
// spec.openapi.typesToInclude's default.
var AllPathTypes = []PathType{BackendPathType, FunctionPathType, PagePathType, SystemPathType}

// NewOneOffBuilder returns the Builder of a host with auth enabled and the
// default spec.openapi.typesToInclude -- the configuration client-side codegen
// assumes when it renders a KDexFunction's one-off document without a running
// host (fngogen local / build-spec).
func NewOneOffBuilder() *Builder {
	return &Builder{
		SecuritySchemes: SecuritySchemes(),
		TypesToInclude:  slices.Clone(AllPathTypes),
	}
}

func oauthFlowScopes() map[string]string {
	return map[string]string{
		"openid":  "standard oidc scope",
		"profile": "user profile info",
	}
}

// SecuritySchemes is the registry of security schemes an auth-enabled host
// defines. A KDexFunction's spec.api refers to them by name; the one-off
// document carries the ones its operations use. A host with auth disabled
// defines none.
func SecuritySchemes() *openapi.SecuritySchemes {
	req := &openapi.SecuritySchemes{}

	(*req)["apiKeyCookie"] = &openapi.SecuritySchemeRef{
		Value: &openapi.SecurityScheme{
			Description: "Stateless API Token in 'X-API-TOKEN' cookie",
			Type:        schemeTypeAPIKey,
			In:          "cookie",
			Name:        "X-API-TOKEN",
		},
	}

	(*req)["apiKeyHeader"] = &openapi.SecuritySchemeRef{
		Value: &openapi.SecurityScheme{
			Description: "Stateless API Token in 'X-API-TOKEN' header",
			Type:        schemeTypeAPIKey,
			In:          "header",
			Name:        "X-API-TOKEN",
		},
	}

	(*req)["apiKeyQuery"] = &openapi.SecuritySchemeRef{
		Value: &openapi.SecurityScheme{
			Description: "Stateless API Token in 'api_token' query parameter",
			Type:        schemeTypeAPIKey,
			In:          "query",
			Name:        "api_token",
		},
	}

	(*req)["bearer"] = &openapi.SecuritySchemeRef{
		Value: &openapi.SecurityScheme{
			Description:  "Bearer Token - This is the default scheme",
			Type:         "http",
			Scheme:       "bearer",
			BearerFormat: "JWT",
		},
	}

	(*req)["oauth2"] = &openapi.SecuritySchemeRef{
		Value: &openapi.SecurityScheme{
			Description: "OAuth2 authentication using Password, Authorization Code, or Client Credentials Flow",
			Flows: &openapi.OAuthFlows{
				AuthorizationCode: &openapi.OAuthFlow{
					AuthorizationURL: OAuthAuthorizePath,
					Scopes:           oauthFlowScopes(),
					TokenURL:         TokenPath,
				},
				ClientCredentials: &openapi.OAuthFlow{
					Scopes:   oauthFlowScopes(),
					TokenURL: TokenPath,
				},
				Password: &openapi.OAuthFlow{
					Scopes:   oauthFlowScopes(),
					TokenURL: TokenPath,
				},
			},
			Type: "oauth2",
		},
	}

	(*req)["oidc"] = &openapi.SecuritySchemeRef{
		Value: &openapi.SecurityScheme{
			Description:      "OpenID Connect discovery",
			OpenIdConnectUrl: "/.well-known/openid-configuration",
			Type:             "openIdConnect",
		},
	}

	return req
}
