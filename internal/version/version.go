// Package version holds the application version string embedded at build time
// via -ldflags. When no version is injected the sentinel value "dev" is used.
package version

// Version is the application version. It is set at link time via:
//
//	-ldflags "-X github.com/igorzel/mytets/internal/version.Version=X.Y.Z"
//
// When the binary is built without -ldflags the value defaults to "dev".
var Version = "dev"
