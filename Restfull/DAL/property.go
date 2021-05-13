package DAL

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Property struct {
	PropId       int
	PropertyName string
	Address      string
	City         string
	Bedroom      int
}

func GetProperty() []Property {

	var properties []Property

	session := Connect()

	collection := session.Database("webinardb").Collection("properties")
	cur, _ := collection.Find(context.TODO(), bson.M{})

	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {

		var prop Property

		err := cur.Decode(&prop)

		if err != nil {
			log.Fatal(err)
		}

		properties = append(properties, prop)
	}
	return properties

}
func GetPropertyByPropId(id string) Property {
	var prop Property
	session := Connect()
	collection := session.Database("webinardb").Collection("properties") 
	
	filter := bson.M{"propid": id}
	err := collection.FindOne(context.TODO(), filter).Decode(&prop)
	if err != nil {
		log.Println(err)
	}
	return prop
}
func PutProperty(property Property) {
	session := Connect()

	collection := session.Database("webinardb").Collection("properties") 
	
	filter := bson.M{"propid":property.PropId}

	update := bson.M{
		"$set": bson.M{
		  "propertyname": property.PropertyName,
		  "address": property.Address,
		  "city": property.City,
	},
	  }
	
	result,err := collection.UpdateOne(context.TODO(), filter, update)

	log.Println(result)
	log.Println(err)
}
	

func FilterProperty(city string,bedroom int) Property {
	var prop Property
	session := Connect()
	collection := session.Database("webinardb").Collection("properties") 
	filter := bson.M{"city": city,"Bedroom":bedroom}
	err := collection.FindOne(context.TODO(), filter).Decode(&prop)
	if err != nil {
		log.Println(err)
	}
	return prop
}
func InsertProperty(prop Property) {
	session := Connect()
	collection := session.Database("webinardb").Collection("properties") //
	result, err := collection.InsertOne(context.TODO(), prop)
	log.Println(result)
	log.Println(err)
}
func DeletePropertyByPropId(propid int) {

	session := Connect()
	collection := session.Database("webinardb").Collection("properties") 
	filter := bson.M{"propid": propid}
	result, err := collection.DeleteOne(context.TODO(), filter)
	fmt.Println(result, err)
}
func Connect() *mongo.Client {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

	// Connect to MongoDB
	client, err := mongo.Connect(context.TODO(), clientOptions)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("DB connected...")
	return client
}
