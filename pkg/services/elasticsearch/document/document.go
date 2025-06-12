package document

type Document interface {
	IndexName() string
	DocumentID() int
}
