Just a simple project to try out a websocket in Go

Go handles the calculations and sends updates through a websocket

React then just displays the received data, no calculations are happening there (except the initial randomization)

yeah this could have been done just on the client but this is more fun

how to run this ?

### Backend

- navigate to backend and run `go run ./server`

### Frontend

- navigate to frontend and run `npm start`

the server can crash if you change the step delay too often and happen to send a message at the same time as you are receiving one
