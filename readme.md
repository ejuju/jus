# Ju's system for distributing conversation-based applications over the Internet

Let's say you're building an app that the user can interact with through a conversation, for example:

Use cases:
- Customer support
- Conversational AI
- Real-time chat
- Restaurant reservation
- Parcel delivery updates

Ju's was designed as an alternative to the current web ecosystem for 
distributing applications through a conversation-based user interface.

- `Jul`: Ju's UI scripting language
- `JuTP`: Ju's transfer protocol

Architecture:
- UI executes scripts that can write messages and send back data to the server.
- UI can periodically retrieve information from server if the script is installed
- Client-side scripts can persist data on the client in a local key-value store
- Client-side scripts can be persisted in the key-value store and run periodically

### Principles

#### Make it simple
The design of the system was made to be easy to understand and implement.

#### Make data collection obvious to the end user
All data collection from client to server is tied to visible message in the UI.

## Todo

- [ ] Implement HTTP gateway library
- [ ] Implement web-based UI

- Implement examples
    - [ ] Echo server
    - [ ] Restaurant reservation
