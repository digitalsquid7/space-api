package querymodifiers

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

type Fields struct {
	fields []Field
}

func (f Fields) GetFieldByAPIName(name string) (Field, bool) {
	for i := range f.fields {
		if f.fields[i].APIName == name {
			return f.fields[i], true
		}
	}

	return Field{}, false
}
