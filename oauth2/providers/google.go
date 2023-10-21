package providers

import (
	"context"
	"encoding/json"
	"io"

	"golang.org/x/oauth2"

	"github.com/alayton/papers"
)

const googleUserInfoEndpoint = "https://openidconnect.googleapis.com/v1/userinfo"

var GoogleScopes = []string{"openid", "email"}

type googleIdentityResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

func (g *googleIdentityResponse) GetID() string {
	return g.ID
}

func (g *googleIdentityResponse) GetEmail() string {
	return g.Email
}

func GetGoogleIdentity(ctx context.Context, cfg oauth2.Config, token *oauth2.Token) (papers.OAuth2Identity, error) {
	client := cfg.Client(ctx, token)
	resp, err := client.Get(googleUserInfoEndpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var identity googleIdentityResponse
	if err := json.Unmarshal(buf, &identity); err != nil {
		return nil, err
	}

	return &googleIdentityResponse{
		ID:    identity.ID,
		Email: identity.Email,
	}, nil
}
