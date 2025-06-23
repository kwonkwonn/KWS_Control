package structure

import "errors"

func ErrCoreNotFound(uuid UUID) error {
	return errors.New("core not found for vm uuid: " + string(uuid))
}

func ErrVmNotFound(uuid UUID) error {
	return errors.New("vm not found of uuid: " + string(uuid))
}
