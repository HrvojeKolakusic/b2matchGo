package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type User struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	OrgId string `json:"OrgId"`
}

type Event struct {
	Id    string `json:"Id"`
	Name  string `json:"Name"`
	OrgId string `json:"OrgId"`
	Users []User `json:"Users"`
}

type Invitation struct {
	MeetingName string `json:"MeetingName"`
	SenderId    string `json:"SenderId"`
	TargetId    string `json:"TargetId"`
	Date        string `json:"Date"`
	Time        string `json:"Time"`
	Status      bool   `json:"Status"`
}

type Meeting struct {
	SenderId    string       `json:"SenderId"`
	Name        string       `json:"Name"`
	Date        string       `json:"Date"`
	Time        string       `json:"Time"`
	Invitations []Invitation `json:"Invitations"`
}

type NewMeetingData struct {
	Name    string   `json:"Name"`
	Date    string   `json:"Date"`
	Time    string   `json:"Time"`
	UserIds []string `json:"UserIds"`
}

// global arrays are used to simulate a database
var Users []User
var Events []Event
var Meetings []Meeting
var Invitations []Invitation

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the HomePage!")
	fmt.Println("Endpoint Hit: homePage")
}

//called when a POST request is sent to /event/{eventId}/addUser
func addUserToEvent(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: addUserToEvent")

	vars := mux.Vars(r)
	eventId := vars["eventId"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var user User
	json.Unmarshal(reqBody, &user)

	for i, event := range Events {
		if event.Id == eventId {
			if event.OrgId != user.OrgId {
				fmt.Println("User organisation and event organisation don't match")
				break
			}

			event.Users = append(event.Users, user)
			Events[i] = event
			json.NewEncoder(w).Encode(user)
		}
	}
}

//called when a GET request is sent to /allEvents
func returnAllEvents(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint Hit: returnAllEvents")
	for _, event := range Events {
		fmt.Fprintf(w, "Name: "+event.Name+", Users in event: ")
		for j, user := range event.Users {
			fmt.Fprintf(w, user.Name)
			if j+1 < len(event.Users) {
				fmt.Fprintf(w, ", ")
			}
		}
		fmt.Fprintf(w, "\n")
	}
}

//called when a POST request is sent to /event/{eventId}/{userId}/newMeeting
func createNewMeeting(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: createNewMeeting")

	vars := mux.Vars(r)
	eventId := vars["eventId"]
	senderId := vars["userId"]

	reqBody, _ := ioutil.ReadAll(r.Body)
	var newMeeting NewMeetingData
	json.Unmarshal(reqBody, &newMeeting)

	var meeting Meeting
	meeting.SenderId = senderId
	meeting.Name = newMeeting.Name
	meeting.Date = newMeeting.Date
	meeting.Time = newMeeting.Time

	for _, targetId := range newMeeting.UserIds {
		var inv Invitation
		inv.MeetingName = newMeeting.Name
		inv.SenderId = senderId

		var found bool = false
		for _, event := range Events {
			if event.Id == eventId {
				for _, user := range event.Users {
					if user.Id == targetId {
						found = true
						break
					}
				}
				if found {
					break
				}
			}
		}

		if !found {
			fmt.Println("User " + targetId + " not in event " + eventId)
			continue
		}

		inv.TargetId = targetId
		inv.Date = newMeeting.Date
		inv.Time = newMeeting.Time
		inv.Status = false
		Invitations = append(Invitations, inv)
		meeting.Invitations = append(meeting.Invitations, inv)
	}

	Meetings = append(Meetings, meeting)

	json.NewEncoder(w).Encode(meeting)
}

//called when a GET request is sent to /event/{eventId}/{userId}/invitations
func checkUserInvitations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	fmt.Println("Endpoint hit: checkUserInvitations: " + userId)

	for _, inv := range Invitations {
		if inv.TargetId == userId {
			var senderName string
			for _, user := range Users {
				if inv.SenderId == user.Id {
					senderName = user.Name
				}
			}
			var accepted string
			if inv.Status {
				accepted = "YES"
			} else {
				accepted = "NO"
			}
			fmt.Fprintf(w, "Meeting: "+inv.MeetingName+", From: "+senderName+", Date: "+inv.Date+", Time: "+inv.Time+", Accepted: "+accepted+"\n")
		}
	}
}

//called when a GET request is sent to /event/{eventId}/{userId}/meetings
func checkMeetingInvitations(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userId := vars["userId"]

	fmt.Println("Endpoint hit: checkMeetingInvitations: " + userId)

	for _, meeting := range Meetings {
		if meeting.SenderId == userId {
			var numOfYes int = 0
			var numOfNo int = 0
			for _, inv := range meeting.Invitations {
				if inv.Status {
					numOfYes++
				} else {
					numOfNo++
				}
			}

			fmt.Fprintf(w, "Meeting: "+meeting.Name+", Date: "+meeting.Date+", Time: "+meeting.Time+", YES: "+strconv.Itoa(numOfYes)+", NO: "+strconv.Itoa(numOfNo)+"\n")
		}
	}
}

//called when a PUT request is sent to /event/{eventId}/{userId}/answerInvitation
func answerInvitation(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Endpoint hit: answerInvitation")
	reqBody, _ := ioutil.ReadAll(r.Body)
	var ansInv Invitation
	json.Unmarshal(reqBody, &ansInv)

	for i, inv := range Invitations {
		if inv.MeetingName == ansInv.MeetingName {
			Invitations[i].Status = true
		}
	}

	for i, meeting := range Meetings {
		if meeting.Name == ansInv.MeetingName {
			for j, inv := range meeting.Invitations {
				if inv.TargetId == ansInv.TargetId {
					Meetings[i].Invitations[j].Status = true
				}
			}
		}
	}
}

func handleRequests() {
	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", homePage)

	//endpoints
	myRouter.HandleFunc("/allEvents", returnAllEvents)
	myRouter.HandleFunc("/event/{eventId}/addUser", addUserToEvent).Methods("POST")
	myRouter.HandleFunc("/event/{eventId}/{userId}/newMeeting", createNewMeeting).Methods("POST")
	myRouter.HandleFunc("/event/{eventId}/{userId}/invitations", checkUserInvitations)
	myRouter.HandleFunc("/event/{eventId}/{userId}/meetings", checkMeetingInvitations)
	myRouter.HandleFunc("/event/{eventId}/{userId}/answerInvitation", answerInvitation).Methods("PUT")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func main() {
	fmt.Println("b2match task")

	//prepare dummy data
	Users = []User{
		{Id: "1", Name: "User1", OrgId: "1"},
		{Id: "2", Name: "User2", OrgId: "1"},
		{Id: "3", Name: "User3", OrgId: "1"},
		{Id: "4", Name: "User4", OrgId: "2"},
		{Id: "5", Name: "User5", OrgId: "2"},
		{Id: "6", Name: "User6", OrgId: "2"},
	}

	Events = []Event{
		{
			Id:    "1",
			Name:  "Event1",
			OrgId: "1",
			Users: []User{Users[0], Users[1], Users[2]},
		},
		{Id: "2", Name: "Event2", OrgId: "2", Users: []User{}},
	}

	handleRequests()
}
