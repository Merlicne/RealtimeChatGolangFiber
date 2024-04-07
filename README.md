# RealtimeChatGolangFiber
This repo is about implementing Real Time chat (multiple room) with Golang using Fiber and websocket.

Server
- Serve websocket to user that request connection
- Navigate user to the room that user want to join
          
Room
- Act as a room for user to chating with other
- Manipulate user connections 
- Receive message from every user in the room
- Broadcast message to erery user in the room

Client (user)
- Represent user connection
- Write message and send to the room
- Read message from room's broadcasting
 
user -> server  -> chatroom1 // room will be create if not exist
                -> chatroom2 
                -> chatroom3  
                -> chatroom4
                .
                .
                .
