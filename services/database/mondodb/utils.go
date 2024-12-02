package mondodb

import "go.mongodb.org/mongo-driver/bson/primitive"

func StringToPrimitiveObjectID(id string) (primitive.ObjectID, error) {
	primitiveObjectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return primitive.ObjectID{}, err
	}
	return primitiveObjectID, nil
}

func PrimitiveObjectIDToString(primitiveObjectID primitive.ObjectID) string {
	return primitiveObjectID.Hex()
}
