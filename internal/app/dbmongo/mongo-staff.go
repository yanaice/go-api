package dbmongo

import (
	"context"
	"go-starter-project/internal/app/database"
	"go-starter-project/pkg/auth"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type staffDatabaseImpl struct{}

func GetStaffDatabsase() database.StaffDatabase {
	staffDB := &staffDatabaseImpl{}
	return staffDB
}

func (s *staffDatabaseImpl) GetUserByName(username string) (auth.UserStaff, error) {
	var user auth.UserStaff
	var userWithID bson.M

	res := mdb.Collection("user_staffs").FindOne(context.TODO(), bson.M{"username": username})
	err := res.Decode(&user)
	if err != nil {
		return auth.UserStaff{}, returnError(err)
	}
	err = res.Decode(&userWithID)
	if err != nil {
		return auth.UserStaff{}, returnError(err)
	}
	user.ID = userWithID["_id"].(primitive.ObjectID).Hex()
	return user, nil
}

func (s *staffDatabaseImpl) CreateUser(username, password string, roleLevel int) (auth.UserStaff, error) {
	var user auth.UserStaff

	payload := bson.D{
		{"username", username},
		{"password", password},
		{"role_level", roleLevel},
	}

	res, err := mdb.Collection("user_staffs").InsertOne(context.TODO(), payload)
	if err != nil {
		return auth.UserStaff{}, returnError(err)
	}
	user.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return user, nil
}
