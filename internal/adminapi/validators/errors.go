package validators

import "fmt"

const uniqueConstraintViolationHeader = "unique constraints violation"

// UniqueConstraintViolationError is an error that indicates that a unique
// constraint was violated for a specific key. The error will indicate the
// type of the key as well as the reference Kubernetes object name in the
// error output.
type UniqueConstraintViolationError struct {
	ObjectType      string
	ObjectName      string
	ObjectNamespace string
	Type            string
	Key             string
}

func (u UniqueConstraintViolationError) Error() string {
	return fmt.Sprintf(
		"%s for type %s on key %s for object type %s named %s in namespace %s",
		uniqueConstraintViolationHeader,
		u.Type, u.Key,
		u.ObjectType, u.ObjectName, u.ObjectNamespace,
	)
}
