package exceptions

// InterfaceException lo que tienen en comun las excepciones controladas (FieldsException, VideosException): un
// mensaje para el usuario y el codigo HTTP equivalente.
type InterfaceException interface {
	error
	GetCode() int
}
