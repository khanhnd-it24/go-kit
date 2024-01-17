package mgocvt

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func StrToObjectId(str string) (primitive.ObjectID, error) {
	objectId, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return primitive.NilObjectID, err
	}
	return objectId, nil
}

func NewObjectId(str string) primitive.ObjectID {
	objectId, err := primitive.ObjectIDFromHex(str)
	if err != nil {
		return primitive.NewObjectID()
	}
	return objectId
}

func ObjectIdToStr(obj primitive.ObjectID) string {
	return obj.Hex()
}
