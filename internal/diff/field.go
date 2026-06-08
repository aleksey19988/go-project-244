package diff

type Field struct {
	Name     string
	Status   string
	OldValue any
	NewValue any
	Depth    int
	Children []Field
}
