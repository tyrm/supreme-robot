package models

// Delete a struct
func (c *Client) Delete(obj interface{}) error {
	switch obj := obj.(type) {
	case *Domain:
		return obj.delete(c)
	default:
		return errUnknownType
	}
}
