package protocol

type SignatureHelpOptions struct {
    TriggerCharacters      []string `json:"triggerCharacters,omitempty"`
    RetriggerCharacters    []string `json:"retriggerCharacters,omitempty"`
}
