package protocol

type CompletionItem struct {
    Label         string      `json:"label"`
    Kind          int         `json:"kind,omitempty"`
    Detail        string      `json:"detail,omitempty"`
    Documentation MarkupContent `json:"documentation,omitempty"`
    InsertText    string      `json:"insertText,omitempty"`
}
