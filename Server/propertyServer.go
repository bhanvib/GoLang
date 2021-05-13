package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"project.com/proto/propertypb"
)

var collection *mongo.Collection

type propertyItem struct {
	Propid       string `bson:"_propid,omitempty"`
	Propertyname string `bson:"propertyname"`
	Address      string `bson:"address"`
	City         string `bson:"city"`
}
type Server struct {
}
//CreateProperty
//func(server Server) CreateProperty ( ctx context.Context, request *propertypb.CreatePropertyRequest) (*propertypb.CreatePropertyRequest, error){
func (server Server) CreateProperty(ctx context.Context, request *propertypb.CreatePropertyRequest) (*propertypb.CreatePropertyResponse, error) {

	fmt.Println("Create property request")
	property := request.GetProperty()

	data := propertyItem{Propid: property.GetPropid(),
		Propertyname: property.GetPropertyname(),
		Address:      property.GetAddress(),
		City:         property.GetCity(),
	}

	res, err := collection.InsertOne(context.Background(), data)
	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Internal error: %v", err),
		)
	}
	oid, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot convert to OID"),
		)
	}

	return &propertypb.CreatePropertyResponse{
		Property: &propertypb.Property{
			Propid:       oid.Hex(),
			Propertyname: property.GetPropertyname(),
			Address:      property.GetAddress(),
			City:         property.GetCity(),
		},
	}, nil

}


//GETPROPERTY
func (server Server) GetPropertyById(ctx context.Context, request *propertypb.GetPropertyByIdRequest) (*propertypb.GetPropertyByIdResponse, error) {
	fmt.Println("Read blog request")

	propertyPropid := request.GetPropertyPropid()
	oid, err := primitive.ObjectIDFromHex(propertyPropid)
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	data := &propertyItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}

	return &propertypb.GetPropertyByIdResponse{
		Property: dataToPropertyPb(data),
	}, nil
}
func dataToPropertyPb(data *propertyItem) *propertypb.Property {
	return &propertypb.Property{
		Propid:       data.Propid,
		Propertyname: data.Propertyname,
		Address:      data.Address,
		City:         data.City,
	}
}

//UPDATEPROPERTY
func (server Server) UpdateProperty(ctx context.Context, request *propertypb.UpdatePropertyRequest) (*propertypb.UpdatePropertyResponse, error) {
	fmt.Println("Update blog request")
	property := request.GetProperty()
	oid, err := primitive.ObjectIDFromHex(property.GetPropid())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}
data := &propertyItem{}
	filter := bson.M{"_id": oid}

	res := collection.FindOne(context.Background(), filter)
	if err := res.Decode(data); err != nil {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog with specified ID: %v", err),
		)
	}
	data.Propid = property.GetPropid()
	data.Propertyname = property.GetPropertyname()
	data.City = property.GetCity()
	data.Address = property.GetAddress()

	_, updateErr := collection.ReplaceOne(context.Background(), filter, data)
	if updateErr != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot update object in MongoDB: %v", updateErr),
		)
	}

	return &propertypb.UpdatePropertyResponse{
		Property: dataToPropertyPb(data),
	}, nil

}


//DELETEPROPERTY
func (server Server) DeleteProperty(ctx context.Context, request *propertypb.DeletePropertyRequest) (*propertypb.DeletePropertyResponse, error) {
	fmt.Println("Delete blog request")
	oid, err := primitive.ObjectIDFromHex(request.GetPropertyPropid())
	if err != nil {
		return nil, status.Errorf(
			codes.InvalidArgument,
			fmt.Sprintf("Cannot parse ID"),
		)
	}

	filter := bson.M{"_propid": oid}

	res, err := collection.DeleteOne(context.Background(), filter)

	if err != nil {
		return nil, status.Errorf(
			codes.Internal,
			fmt.Sprintf("Cannot delete object in MongoDB: %v", err),
		)
	}

	if res.DeletedCount == 0 {
		return nil, status.Errorf(
			codes.NotFound,
			fmt.Sprintf("Cannot find blog in MongoDB: %v", err),
		)
	}

	return &propertypb.DeletePropertyResponse{PropertyPropid: request.GetPropertyPropid()}, nil
}

//GETALLPROPERTY
func (server Server) GetProperty(req *propertypb.GetPropertyRequest, stream propertypb.PropertyService_GetPropertyServer) error {
	fmt.Println("List blog request")

	cur, err := collection.Find(context.Background(), primitive.D{{}})
	if err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	defer cur.Close(context.Background())
	for cur.Next(context.Background()) {
		data := &propertyItem{}
		err := cur.Decode(data)
		if err != nil {
			return status.Errorf(
				codes.Internal,
				fmt.Sprintf("Error while decoding data from MongoDB: %v", err),
			)

		}
		stream.Send(&propertypb.GetPropertyResponse{Property: dataToPropertyPb(data)})
	}
	if err := cur.Err(); err != nil {
		return status.Errorf(
			codes.Internal,
			fmt.Sprintf("Unknown internal error: %v", err),
		)
	}
	return nil
}

func main() {
	
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	fmt.Println("Connecting to MongoDB")
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	err = client.Connect(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Property Service Started")
	collection = client.Database("mydb").Collection("properties")

	lis, err := net.Listen("tcp", "0.0.0.0:50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	opts := []grpc.ServerOption{}
	s := grpc.NewServer(opts...)
	propertypb.RegisterPropertyServiceServer(s, &Server{})

	reflection.Register(s)

	go func() {
		fmt.Println("Starting Server...")
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)
	<-ch
	fmt.Println("Closing MongoDB Connection")
	if err := client.Disconnect(context.TODO()); err != nil {
		log.Fatalf("Error on disconnection with MongoDB : %v", err)
	}
	fmt.Println("Closing the listener")
	if err := lis.Close(); err != nil {
		log.Fatalf("Error on closing the listener : %v", err)
	}
	fmt.Println("Stopping the server")
	s.Stop()
	fmt.Println("End of Program")
}
