package lsp

type WorkspaceFoldersServerCapabilities struct {
	Supported           bool `json:"supported,omitempty"`
	ChangeNotifications bool `json:"changeNotifications,omitempty"`
}

type WorkspaceCapabilities struct {
	WorkspaceFolders *WorkspaceFoldersServerCapabilities `json:"workspaceFolders,omitempty"`
}

type DidChangeConfigurationParams struct {
	Settings interface{} `json:"settings"`
}

type DidChangeWatchedFilesParams struct {
	Changes []FileEvent `json:"changes"`
}

type FileEvent struct {
	URI  DocumentUri `json:"uri"`
	Type int         `json:"type"`
}
