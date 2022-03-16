run the app either with an IDE or with the console command: go run main.go

to test the app go to localhost:10000/
--------------------------
to add users to events send a POST request to localhost:10000/event/{eventId}/addUser
e.g. localhost:10000/event/1/addUser with body {Id: "7", Name: "User7", OrgId: "1"}
--------------------------
to get a list of all the events go to localhost:10000/allEvents
--------------------------
to create a new meeting send a POST request to localhost:10000/event/{eventId}/{userId}/newMeeting
e.q. localhost:10000/event/1/1/newMeeting with body: 
{
    "Name": "test_meeting", 
    "Date": "1.1.2000.",
    "Time": "9:00",
    "UserIds": ["2", "3"]
}
--------------------------
to check users invitations go to localhost:10000/event/{eventId}/{userId}/invitations
e.g. localhost:10000/event/1/2/invitations
--------------------------
to check users meetings go to localhost:10000/event/{eventId}/{userId}/meetings
e.g. localhost:10000/event/1/2/meetings
--------------------------
to answer an invitation send a PUT request to localhost:10000/event/{eventId}/{userId}/answerInvitation
e.g. localhost:10000/event/1/2/answerInvitation with body:
{
    "MeetingName": "test_meeting",
    "SenderId": "1",
    "TargetId": "2",
    "Date": "1.1.2000.",
    "Time": "9:00",
    "Status": true
}
