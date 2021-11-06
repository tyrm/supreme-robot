package models

// Create a struct
func (c *Client) Create(obj interface{}) error {
	switch obj := obj.(type) {
	case Domain:
		return obj.create(c)
	case *Domain:
		return obj.create(c)
	case Record:
		return obj.create(c)
	case *Record:
		return obj.create(c)
	case User:
		return obj.create(c)
	case *User:
		return obj.create(c)
	default:
		return errUnknownType
	}
}