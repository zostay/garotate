package rotate

import (
	"context"
	"time"

	"github.com/zostay/garotate/pkg/secret"
)

// Storage is the client interface implemented by storage client plugins. A
// storage plugin must provide a means of checking the last time a given secret
// was stored and a means of updating the stored password upon rotation.
type Storage interface {
	// Name must return the value in the configuration for reference to the
	// adminstrator running the rotation service. This will be used in errors
	// and logs.
	Name() string

	// LastSaved must return the timestamp the given secret was last updated in
	// the storage or an error. If there is no last saved date,
	// that is an error. The plugin should return secret.ErrNotFound.
	//
	// The context provides a logger via context tools in the config package.
	//
	// The secret.Storage describes the information about the secret as it
	// pertains to the storage client.
	//
	// The final string value is the individual key as a given secret might have
	// multiple values. For example, an account might have a username and password
	// that changes each time or it might have an API key and a secret key.
	// This method will be called for each.
	LastSaved(context.Context, secret.Storage, string) (time.Time, error)

	// SaveKeys will be called a single time for each rotation. It must perform
	// storage of the secret following a fresh secret rotation or return an
	// error.
	//
	// The context provides a logger via context tools in the config package.
	//
	// The secret.Storage describes information about the secret as it pertains
	// to the storage client.
	//
	// The map is the set of values to be stored. This will already be remapped
	// from the values output by the rotation client into the values configured
	// for storage.
	SaveKeys(context.Context, secret.Storage, secret.Map) error
}

// Client is the interface implemented by rotation plugins. These are plugins
// responsible for performing the rotation of secrets.
type Client interface {
	// Name must return the value in the configuration for reference to the
	// administrator running the rotation service. This will be used in errors
	// and logs.
	Name() string

	// Keys must return the keys that the rotation plugin will return when
	// RotateSecret() is called. The values are ignored. It is assumed that
	// every key returned here is required for complete rotation, so any missing
	// key will trigger a fresh rotation. There are no optional keys.
	Keys() secret.Map

	// LastRotated must return the date of the most recent rotation of the given
	// secret or an error.
	//
	// The context provides a logger via context tools in the config package.
	//
	// The secret.Info describes the secret to be rotated.
	LastRotated(context.Context, secret.Info) (time.Time, error)

	// RotateSecret must immediately rotate the secret and return a map
	// containing all the new values. The keys returned should be carefully
	// documented and be consistent so the administrator running the service can
	// remap them to storages as required using static names in the
	// configuration. It is recommended that the names be the most natural names
	// for the accounting system being rotated.
	//
	// If rotation cannot be performed, an error must be returned.
	//
	// The context provides a logger via context tools in the config package.
	//
	// The secret.Info describes the secret to be rotated.
	RotateSecret(context.Context, secret.Info) (secret.Map, error)
}
