package main

import (
	"context"
	"fmt"
	"io"

	"log"

	"google.golang.org/grpc"
	"project.com/proto/propertypb"
)

func main() {

	fmt.Println("Property Client")
    opts := grpc.WithInsecure()
    cc, err := grpc.Dial("localhost:50051", opts)
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	defer cc.Close()

	c := propertypb.NewPropertyServiceClient(cc)
	fmt.Println(" before Property")
	property := &propertypb.Property{
		    Propid:     " 10",
			Propertyname:    "hyderabad",
		    Address:  "hyderabad",
		    City:     "hyederabad",
		   
	}
	fmt.Println(" after Property")
	createPropertyRes, err := c.CreateProperty(context.Background(), &propertypb.CreatePropertyRequest{Property: property})
	if err != nil {
		log.Fatalf("Unexpected error: %v", err)
	}
	fmt.Printf("Blog has been created: %v", createPropertyRes)
	PropertyPropid := createPropertyRes.GetProperty().GetPropid()
	fmt.Print(PropertyPropid)

	

	//updateProperty
	fmt.Println("before update")
	newProperty := &propertypb.Property{
		Propid:       PropertyPropid,
		Propertyname: "name",
		City:    "My First Blog g(edited)",
		Address:  "Content",
	}
	updateRes, updateErr := c.UpdateProperty(context.Background(), &propertypb.UpdatePropertyRequest{Property: newProperty})
	if updateErr != nil {
		fmt.Printf("Error happened while updating: %v \n", updateErr)
	}
	fmt.Printf(" after update %v\n", updateRes)

	//getproperty
	fmt.Println("get property")

	_, err2 := c.GetPropertyById(context.Background(), &propertypb.GetPropertyByIdRequest{PropertyPropid: "609ba2bfb18cf9362c1db7ae"})
	if err2 != nil {
		fmt.Printf("Error happened while reading: %v \n", err2)
	}

	GetPropertyByIdReq := &propertypb.GetPropertyByIdRequest{PropertyPropid: PropertyPropid}
	GetPropertyRes, readPropertyErr := c.GetPropertyById(context.Background(), GetPropertyByIdReq)
	if readPropertyErr != nil {
		fmt.Printf("Error happened while reading: %v \n", readPropertyErr)
	}
fmt.Printf("getting property: %v \n", GetPropertyRes)


	//delete
	deleteRes, deleteErr := c.DeleteProperty(context.Background(), &propertypb.DeletePropertyRequest{PropertyPropid: "609ba7e518705db2ad6259b3"})

	if deleteErr != nil {
		fmt.Printf("Error happened while deleting: %v \n", deleteErr)
	}
	fmt.Printf("deleted: %v \n", deleteRes)

	
	//getproperty
	stream, err := c.GetProperty(context.Background(), &propertypb.GetPropertyRequest{})
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Something happened: %v", err)
		}
		fmt.Println(res.GetProperty())
	}
}