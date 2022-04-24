package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/tus/tusd/pkg/filestore"
	tusd "github.com/tus/tusd/pkg/handler"
)

//struct for tokens
type Tokens struct {
	Refresh string `json:"refresh"`
	Access  string `json:"access"`
}

//struct for video id (Computer vision)
type Id struct {
	ID int `json:"id"`
}

func login(w http.ResponseWriter, r *http.Request) (access, refresh string) {

	//values for login
	values, _ := json.Marshal(map[string]string{
		"username": "aleksejpavlovv6@gmail.com",
		"password": "656900Alex_",
	})

	responseBody := bytes.NewBuffer(values)

	//request to login
	resp, err := http.Post("https://api.gcdn.co/auth/jwt/login", "application/json", responseBody)

	if err != nil {
		log.Fatal(err)
	}

	//read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//unmarshalling tokens
	token := Tokens{}
	json.Unmarshal(body, &token)

	defer resp.Body.Close()
	return token.Access, token.Refresh
}

func getAllVideos(w http.ResponseWriter, r *http.Request) {
	//getting access token
	access, _ := login(w, r)

	//setting header
	w.Header().Set("Content-Type", "application/json")

	//prepairing request
	url := "https://api.gcdn.co/vp/api/videos"
	var bearer = "Bearer " + access
	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		log.Fatal(err)
	}

	req.Header.Add("Authorization", bearer)

	//sending request
	client := &http.Client{}

	//reading response
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	//prepare response to print out
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	//response to Postman
	w.Write(body)

	//printing response out
	fmt.Println(string(body))
}

func getVideoId(w http.ResponseWriter, r *http.Request) {
	//getting access token
	access, _ := login(w, r)

	//setting header
	w.Header().Set("Content-Type", "application/json")

	//prepairing request
	id := "638280"
	url := "https://api.gcdn.co/vp/api/videos/" + id
	var bearer = "Bearer " + access

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
	}
	req.Header.Add("Authorization", bearer)

	//sending request
	client := &http.Client{}

	//reading response
	resp, err := client.Do(req)

	if err != nil {
		log.Println(err)
	}

	defer resp.Body.Close()

	//prepare response to print out
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Println(err)
	}

	//response to Postman
	w.Write(body)

	//printing response out
	fmt.Printf("%s\n", string(body))
}

//Create video is not ready
func createVideo(w http.ResponseWriter, r *http.Request) {
	// Create a new FileStore instance which is responsible for
	// storing the uploaded file on disk in the specified directory.
	// This path _must_ exist before tusd will store uploads in it.
	// If you want to save them on a different medium, for example
	// a remote FTP server, you can implement your own storage backend
	// by implementing the tusd.DataStore interface.
	store := filestore.FileStore{
		Path: "./uploads",
	}

	// A storage backend for tusd may consist of multiple different parts which
	// handle upload creation, locking, termination and so on. The composer is a
	// place where all those separated pieces are joined together. In this example
	// we only use the file store but you may plug in multiple.
	composer := tusd.NewStoreComposer()
	store.UseIn(composer)

	// Create a new HTTP handler for the tusd server by providing a configuration.
	// The StoreComposer property must be set to allow the handler to function.
	handler, err := tusd.NewHandler(tusd.Config{
		BasePath:              "/files/",
		StoreComposer:         composer,
		NotifyCompleteUploads: true,
	})
	if err != nil {
		panic(fmt.Errorf("unable to create handler: %s", err))
	}

	// Start another goroutine for receiving events from the handler whenever
	// an upload is completed. The event will contains details about the upload
	// itself and the relevant HTTP request.
	go func() {
		for {
			event := <-handler.CompleteUploads
			fmt.Printf("Upload %s finished\n", event.Upload.ID)
		}
	}()

	// Right now, nothing has happened since we need to start the HTTP server on
	// our own. In the end, tusd will start listening on and accept request at
	// http://localhost:8080/files
	http.Handle("/files/", http.StripPrefix("/files/", handler))
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(fmt.Errorf("unable to listen: %s", err))
	}
}

func computerVisionAddVideo(w http.ResponseWriter, r *http.Request) {
	//getting access token
	access, _ := login(w, r)
	var bearer = "Bearer " + access

	//marshalling task for sending request
	task, _ := json.Marshal(map[string]string{
		"url":             "http://10835-vodtestsftp.ams.origin.gcdn.co/shopstory_newdemo_eng_converted.mp4",
		"type":            "cv",
		"keyframnes_only": "1",
		"stop_objects":    "COVERED_BUTTOCKS",
	})

	//prepairing request
	requestBody := bytes.NewBuffer(task)
	req, err := http.NewRequest("POST", "https://api.gcdn.co/vp/api/tasks.json", requestBody)
	if err != nil {
		log.Fatalf(err.Error())
	}
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	//sending request
	client := &http.Client{}

	//reading response
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	//reading id from response body
	respbody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	//response to Postman
	w.Write(respbody)

	//unmarshalling id
	id := Id{}
	json.Unmarshal(respbody, &id)

	//printing out id
	log.Printf("id.ID - %v\n", id.ID)
	log.Printf("respbody value is - %v\n", string(respbody))

}

func computerVisionResult(w http.ResponseWriter, r *http.Request) {
	//getting access token
	access, _ := login(w, r)
	var bearer = "Bearer " + access

	//prepairing request
	id := "11923494"
	url := "https://api.gcdn.co/vp/api/tasks/" + id + ".json"

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Authorization", bearer)
	req.Header.Set("Content-Type", "application/json")

	//sending request
	client := &http.Client{}
	//reading response
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	//preparing response to print out
	bodyGet, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	//response to Postman
	w.Write(bodyGet)

	//printing response out
	log.Printf("Get message: %v\n", string(bodyGet))
}

func handleRequests() {
	http.HandleFunc("/Get", getAllVideos)
	http.HandleFunc("/a", createVideo)
	http.HandleFunc("/", getVideoId)
	http.HandleFunc("/cvAdd", computerVisionAddVideo)
	http.HandleFunc("/cvRes", computerVisionResult)
	log.Fatal(http.ListenAndServe(":8000", nil))
}

func main() {
	handleRequests()
}
