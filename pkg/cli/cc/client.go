package cc

import (
	"os"
	"strconv"
)

const (
	defaultDirectorySvc    = "localhost:9292"
	defaultDirectoryKey    = ""
	defaultDirectoryToken  = ""
	defaultAuthorizerSvc   = "localhost:8282"
	defaultAuthorizerKey   = ""
	defaultAuthorizerToken = ""
	defaultTenantID        = ""
	defaultInsecure        = false
	defaultNoCheck         = false
)

func DirectorySvc() string {
	if directorySvc := os.Getenv("TOPAZ_DIRECTORY_SVC"); directorySvc != "" {
		return directorySvc
	}
	return defaultDirectorySvc
}

func DirectoryKey() string {
	if directoryKey := os.Getenv("TOPAZ_DIRECTORY_KEY"); directoryKey != "" {
		return directoryKey
	}
	return defaultDirectoryKey
}

func DirectoryToken() string {
	if directoryToken := os.Getenv("TOPAZ_DIRECTORY_TOKEN"); directoryToken != "" {
		return directoryToken
	}
	return defaultDirectoryToken
}

func AuthorizerSvc() string {
	if authorizerSvc := os.Getenv("TOPAZ_AUTHORIZER_SVC"); authorizerSvc != "" {
		return authorizerSvc
	}
	return defaultAuthorizerSvc
}

func AuthorizerKey() string {
	if authorizerKey := os.Getenv("TOPAZ_AUTHORIZER_KEY"); authorizerKey != "" {
		return authorizerKey
	}
	return defaultAuthorizerKey
}

func AuthorizerToken() string {
	if authorizerToken := os.Getenv("TOPAZ_AUTHORIZER_TOKEN"); authorizerToken != "" {
		return authorizerToken
	}
	return defaultAuthorizerToken
}

func TenantID() string {
	if tenantID := os.Getenv("ASERTO_TENANT_ID"); tenantID != "" {
		return tenantID
	}
	return defaultTenantID
}

func Insecure() bool {
	if insecure := os.Getenv("TOPAZ_INSECURE"); insecure != "" {
		if b, err := strconv.ParseBool(insecure); err == nil {
			return b
		}
	}
	return defaultInsecure
}

func NoCheck() bool {
	if noCheck := os.Getenv("TOPAZ_NO_CHECK"); noCheck != "" {
		if b, err := strconv.ParseBool(noCheck); err == nil {
			return b
		}
	}
	return defaultNoCheck
}
