//
//  mgo.go
//
//  Created by Arka Mukherjee on 27/02/20.
//
//

package main

import (
	"fmt"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type TestMongoStruct struct {
	Foo string `bson:"foo"`
	Bar string `bson:"bar"`
	Baz int    `bson:"baz"`
}

var isDropMe = false

func main() {
	//const storing DB details
	const (
		Database   = "test"
		Collection = "users"
	)
	session, err := mgo.Dial("127.0.0.1:27017")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	testMongoData := TestMongoStruct{
		Foo: "alice",
		Bar: "bob",
		Baz: 123,
	}

	if isDropMe {
		err = session.DB("test").DropDatabase()
		if err != nil {
			panic(err)
		}
	}

	// Collection data
	c := session.DB(Database).C(Collection)

	/* CRUD CALLS */

	// Insert operation
	if err := c.Insert(testMongoData); err != nil {
		panic(err)
	}

	//Find operation
	var model TestMongoStruct
	err = c.Find(bson.M{"foo": "alice"}).One(&model)
	if err != nil {
		panic(err)
	}

	//Find and Modify
	change := mgo.Change{
		Update:    bson.M{"$set": bson.M{"foo": "berncastel"}},
		ReturnNew: true,
	}
	_, err = c.Find(bson.M{"bar": "bob"}).Apply(change, &model)
	if err != nil {
		panic(err)
	}

	//Update One
	var update bson.M
	update = bson.M{
		"$inc": bson.M{
			"baz": 123,
		},
	}
	err = c.Update(bson.M{"foo": "alice"}, update)
	if err != nil {
		panic(err)
	}

	//Find one (with projection)
	err = c.Find(bson.M{"foo": "alice"}).Select(bson.M{"baz": 1}).One(&model)
	if err != nil {
		panic(err)
	}
	fmt.Println("Value of baz is", model)

	//ObjectIdFromHex Illustration

	oidHex := "XXXXXXXXXXXXXXXX"
	if !bson.IsObjectIdHex(oidHex) {
		fmt.Println("userID is not a valid hex")
	}

	update = bson.M{
		"$inc": bson.M{
			"baz": 123,
		},
	}
	err = c.Update(bson.M{"_id": bson.ObjectIdHex(oidHex)}, update)
	if err != nil {
		panic(err)
	}

	//Delete the entry
	err = c.Remove(bson.M{"foo": "alice"})
	if err != nil {
		panic(err)
	}
}
