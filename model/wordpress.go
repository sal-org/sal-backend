package model

type Rendered struct {
	Rendered string `json:"rendered"`
}

type WordPressPage struct {
	ID      int      `json:"id"`
	Slug    string   `json:"slug"`
	Link    string   `json:"link"`
	Title   Rendered `json:"title"`
	Content Rendered `json:"content"`
}
