package values

// The Values interface is the primary means of handling a collection of values
// Thes same interface and value types are used for both Series values and Index labels
type Values interface {
	Describe() string
	Len() int

	All() interface{}          // all Value/Null elements regardless of null status
	Vals() []interface{}       // all values regardless of null status
	In([]int) Values           // Value/Null elements at one or more integer positions
	Valid() []int              // integer positions of non-null values
	Null() []int               // integer positions of null values
	Element(int) []interface{} // value element at an integer position, in the form []interface{} {val, null}

	ToFloat() Values
	ToInt() Values
	ToString() Values
	ToBool() Values
	ToDateTime() Values
}
