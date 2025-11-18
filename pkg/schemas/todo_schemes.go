package schemas

type TodoCreate struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type TodoPatchUpdate struct {
	Title   *string `json:"title,omitempty"`
	Content *string `json:"content,omitempty"`
}
