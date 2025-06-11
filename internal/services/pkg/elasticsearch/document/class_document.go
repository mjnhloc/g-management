package document

const ClassIndexName = "classes"

type ClassDocument struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Trainer     string  `json:"trainer"`
	Schedule    string  `json:"schedule"`
	Duration    int     `json:"duration"`
	MaxCapacity int     `json:"max_capacity"`
	Description *string `json:"description,omitempty"`
}

func (c *ClassDocument) IndexName() string {
	return ClassIndexName
}

func (c *ClassDocument) DocumentID() string {
	return string(c.ID)
}
