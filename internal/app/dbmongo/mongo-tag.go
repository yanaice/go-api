package dbmongo

import (
	"context"
	"errors"
	"go-starter-project/internal/app/database"
	"go-starter-project/internal/app/model"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type tagDatabaseImpl struct{}

func GetTagDatabase() database.TagDatabase {
	tagDB := &tagDatabaseImpl{}
	return tagDB
}

func (db *tagDatabaseImpl) CreateTag(tag model.Tag) error {
	_, err := mdb.Collection("tag").InsertOne(context.TODO(), tag)
	if err != nil {
		return err
	}
	return nil
}

func (db *tagDatabaseImpl) ReadTag(tagID string) (model.Tag, error) {
	objID, err := primitive.ObjectIDFromHex(tagID)
	if err != nil {
		return model.Tag{}, err
	}

	var tag model.Tag
	res := mdb.Collection("tag").FindOne(context.TODO(), bson.M{"_id": objID})
	if err := res.Decode(&tag); err != nil {
		return model.Tag{}, nil
	}

	return tag, nil
}

func (db *tagDatabaseImpl) ReadTags() ([]model.Tag, error) {
	cur, err := mdb.Collection("tag").Find(context.TODO(), bson.M{})
	if err != nil {
		return nil, err
	}

	var tags []model.Tag
	for cur.Next(context.TODO()) {
		var tag model.Tag
		if err := cur.Decode(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

func (db *tagDatabaseImpl) UpdateTag(tagID string, name string) error {
	objID, err := primitive.ObjectIDFromHex(tagID)
	result, err := mdb.Collection("tag").UpdateOne(context.TODO(), bson.M{"_id": objID}, bson.M{"$set": bson.M{"name": name}})
	if err != nil {
		return err
	}
	if result.MatchedCount == 0 {
		return errors.New("not matched")
	}
	return nil
}

func (db *tagDatabaseImpl) DeleteTag(tagID string) error {
	objID, err := primitive.ObjectIDFromHex(tagID)
	result, err := mdb.Collection("tag").DeleteOne(context.TODO(), bson.M{"_id": objID})
	if err != nil {
		return err
	}

	if result.DeletedCount == 0 {
		return errors.New("not matched")
	}
	return nil
}
