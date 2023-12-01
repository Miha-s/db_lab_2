package main

import (
	"context"
	"encoding/json"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"net/http"
	"time"
)

func registerMongoHandlers(router *http.ServeMux, client *mongo.Client) {
	db := client.Database("company")

	router.HandleFunc("/test/mongo/emploees/insert", 
	func(w http.ResponseWriter, r *http.Request) {
		insertTestEmploeesMongoHandler(w, r, db)
	})
	router.HandleFunc("/test/mongo/emploees/deleteAll", 
	func(w http.ResponseWriter, r *http.Request) {
		deleteAllEmploeesDataMongoHandler(w, r, db)
	})
	router.HandleFunc("/test/mongo/emploees/update", 
	func(w http.ResponseWriter, r *http.Request) {
		updateAllEmploeesPositionMongoHandler(w, r, db)
	})
	router.HandleFunc("/test/mongo/emploees/add_skills", 
	func(w http.ResponseWriter, r *http.Request) {
		addSkillsForEmploeesMongoHandler(w, r, db)
	})
}

func addSkillsForEmploeesMongoHandler(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()

	emploeeCollection := db.Collection("emploees")

	filter := bson.D{} // Your filter criteria to match the documents you want to update
	update := bson.D{
		{"$push", bson.D{
			{"emploee_skills", "NewSkill"},
		}},
	}
	
	_, err := emploeeCollection.UpdateMany(context.TODO(), filter, update)
	if err != nil {
		log.Fatal(err)
	}
	

	duration := time.Since(start).Seconds()

	sendResponse(w, "MongoDB", duration)
}

func updateAllEmploeesPositionMongoHandler(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()

	emploeeCollection := db.Collection("emploees")
	updatePipeline := mongo.Pipeline{
		{
			{"$set", bson.D{{"current_position", "higher"}}},
		},
	}

	
	result, err := emploeeCollection.UpdateMany(context.Background(), bson.M{}, updatePipeline)
	if err != nil {
		http.Error(w, "Error updating", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Seconds()

	sendResponse(w, "MongoDB", duration)
}

func insertTestEmploeesMongoHandler(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method", http.StatusMethodNotAllowed)
		return
	}

	var req UserInsertRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	if req.Count <= 0 {
		http.Error(w, "Count must be a positive integer", http.StatusBadRequest)
		return
	}

	ctx := context.Background()

	emploeeCollection := db.Collection("emploees")

	start := time.Now()


	var userDocuments []interface{}
	for i := 0; i < req.Count; i++ {
		userDocument := bson.D{
			{"name", fmt.Sprintf("Test User %d", i)},
			{"surname", "some_surname"},
			{"telephone", "239011212"},
			{"email", fmt.Sprintf("test%d@example.com", i)},
			{"employment_date", time.Now()},
			{"firing_date", time.Now()},
			{"current_position", "senior"},
			{"emploee_skills", bson.A{"Skill1", "Skill2", "Skill3"}},
			{"current_project", bson.D{
				{"name", "selectedProject"},
				{"start_date", time.Now()},
				{"end_date", time.Now()},
			}},
		}
		userDocuments = append(userDocuments, userDocument)
	}
	var err error
	_, err = emploeeCollection.InsertMany(ctx, userDocuments)
	if err != nil {
		log.Printf("Error inserting: %v\n", err)
		return
	}

	duration := time.Since(start).Seconds()
	sendResponse(w, "MongoDB", duration)
}


func deleteAllEmploeesDataMongoHandler(w http.ResponseWriter, r *http.Request, db *mongo.Database) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method", http.StatusMethodNotAllowed)
		return
	}

	start := time.Now()

	emploeeCollection := db.Collection("emploees")

	if _, err := emploeeCollection.DeleteMany(context.Background(), bson.D{}); err != nil {
		http.Error(w, "Error deleting", http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Seconds()

	sendResponse(w, "MongoDB", duration)
}
