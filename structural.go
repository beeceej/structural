package structural

type (
	// JSON signals that structural should interpret the field as JSON
	// This is only used a a type level identifier; a note for structural
	// so that it can build the correct code
	JSON interface{}

	// XML signals that structural should interpret the field as XML
	// This is only used a a type level identifier; a note for structural
	// so that it can build the correct code
	XML interface{}

	// Route signals that structural should consider a type as a route
	Route interface{}

	// API signals that structural should consider a type as an API
	API interface{}
)
