package app

// Brand strings — single source of truth for product-facing identifiers.
// Reference these from chrome (topbar) and any view that surfaces the host or
// the workspace/config roots, so renaming the product is a one-file change.
const (
	Name           = "autopus"
	Tagline        = "agent daemon"
	AppHost        = "app.autopus.ai"
	WSHost         = "api.autopus.ai/ws"
	WorkspacesRoot = "~/autopus_workspaces"
	ConfigRoot     = "~/.autopus"
)
