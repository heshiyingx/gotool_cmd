package {{.pkg}}


type (
	// {{.upperStartCamelObject}}DBInterface is an interface to be customized, add more methods here,
	// and implement the added methods in custom{{.upperStartCamelObject}}Model.
	{{.upperStartCamelObject}}DBInterface interface {
		{{.lowerStartCamelObject}}Model
	}


)

