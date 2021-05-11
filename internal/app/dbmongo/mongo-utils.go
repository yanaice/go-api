package dbmongo

import (
	"encoding/hex"
	"go-starter-project/pkg/derror"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func returnError(err error) error {
	if err == mongo.ErrNoDocuments {
		return derror.ErrItemNotFound
	}
	if err == primitive.ErrInvalidHex {
		return derror.ErrInputValidationFailed
	}
	if _, ok := err.(hex.InvalidByteError); ok {
		return derror.ErrInputValidationFailed
	}
	return err
}
