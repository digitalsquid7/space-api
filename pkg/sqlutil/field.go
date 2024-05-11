package sqlutil

type FieldType string

const (
	String  FieldType = "string"
	Integer FieldType = "integer"
	Date    FieldType = "date"
)

type Field struct {
	SQLName string
	APIName string
	Type    FieldType
}

type Fields []Field

func (f Fields) GetFieldByAPIName(name string) (Field, bool) {
	for i := range f {
		if f[i].APIName == name {
			return f[i], true
		}
	}

	return Field{}, false
}
