//
//  mongo-go-driver.go
//
//  Created by Arka Mukherjee on 27/02/20.
//
//

package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/x/bsonx"
)

type TestMongoStruct struct {
	Foo string `bson:"foo"`
	Bar string `bson:"bar"`
	Baz int    `bson:"baz"`
}

func main() {
	//const storing DB details
	const (
		Database   = "test"
		Collection = "users"
	)

	ctx := context.TODO()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}

	testMongoData := TestMongoStruct{
		Foo: "alice",
		Bar: "bob",
		Baz: 123,
	}

	// Collection data
	c := client.Database("test").Collection("users")

	/* CRUD CALLS */

	//Inserting Data
	_, err = c.InsertOne(ctx, &testMongoData)
	if err != nil {
		panic(err)
	}

	//Finding Data
	var model TestMongoStruct
	err = c.FindOne(ctx, bson.M{"foo": "alice"}).Decode(&model)

	//Find and Modify
	a, err := primitive.ObjectIDFromHex("XXXXXXXXXX")
	if err == nil {
		_ = c.FindOneAndUpdate(ctx, bson.M{"_id": a}, bson.M{"$set": bson.M{"foo": "berncastel"}})
	}

	//Update
	update := bson.M{}
	update = bson.M{
		"$inc": bson.M{
			"baz": 123,
		},
	}
	_, err = c.UpdateOne(ctx, bson.M{"_id": a}, update)

	//Find with projection
	err = c.FindOne(ctx, bson.M{
		"baz": 246,
	}, options.FindOne().SetProjection(bsonx.Doc{{"baz", bsonx.Int32(1)}})).Decode(model)

	//Delete
	_, err = c.DeleteOne(ctx, bson.M{"_id": a})
	if err != nil {
		panic(err)
	}
}
