package github

import (
	"context"
	"os"
	"reflect"

	"github.com/google/go-github/v42/github"
	"github.com/zostay/garotate/pkg/config"
	"github.com/zostay/garotate/pkg/plugin"
	"golang.org/x/oauth2"
)

// builder implements the plugin.Builder interface and provides the
// factory method for constructing a Client.
type builder struct{}

// Build constructs and returns a github client.
func (b *builder) Build(ctx context.Context, c *config.Plugin) (plugin.Instance, error) {
	token := os.Getenv("GITHUB_TOKEN")
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{
			AccessToken: token,
		},
	)
	oc := oauth2.NewClient(ctx, ts)
	gc := github.NewClient(oc)
	return &Client{gc}, nil
}

// init registers the plugin.
func init() {
	pkg := reflect.TypeOf(Client{}).PkgPath()
	plugin.Register(pkg, new(builder))
}
