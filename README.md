# RealtimeChatGolangFiber
This repo is about implementing Real Time chat (multiple room) with Golang using Fiber and websocket.

Server
- Serve websocket to user that request connection
- Navigate user to the room which user want to join
          
Room
- Act as a room for user to chatting with other
- Manipulate user connections 
- Receive message from every user in the room
- Broadcast message to every user in the room

Client (user)
- Represent user connection
- Write a message and send to the room
- Read incoming message from the room

Message
- Represent a message written by a user
- Contain text_message and Client who write the message
