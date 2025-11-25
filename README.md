# Spelling Bee Game (Assignment 1)
Student: Dmytro Brodar  
Student ID: R00274472  
Course: Distributed Systems Programming (SOFT8023)

### About The Project
This is my first distributed word game made in Go. I used gRPC for communication between client and server.  
The idea of the game is that the player gets 7 letters and one of them is a center letter.   
Every word must use only these seven letters, must include the center letter and must have at least four letters.  
If a word uses all seven letters it gets a bonus.  
 
When I started the game, the client connects to the server. The server sends random letters from pangrams.json.  
Then I type words in the client and the server checks if they are valid using the dictionary.   
If the word is good it gives points and shows total score. If not, it tells why it's wrong.  

### Example of play:
##### Client side:
PS D:\MTU Studying\Year 2\Distributed Systems Programming\Assignment 1\spellingbee> go run ./cmd/client  
---------- Spelling Bee ----------  
Letters: [o n c l g i a] (center i)  
Type words (Ctrl + C to quit).> coin  
Word is valid (+1) | Total: 1  
again  
Word is valid (+5) | Total: 6  
clinic  
Word is valid (+6) | Total: 12  
clon  
Invalid word: Word must include the center letter | TOTAL: 12  
calling  
Word is valid (+7) | Total: 19  
exit status 0xc000013a  

##### Server side:
PS D:\MTU Studying\Year 2\Distributed Systems Programming\Assignment 1\spellingbee> go run ./cmd/server  
2025/10/19 10:22:32 Server running on :50051  
[LOGGING] dictionary check: 'coin' -> true  
[LOGGING] dictionary check: 'again' -> true  
[LOGGING] dictionary check: 'clinic' -> true  
[LOGGING] dictionary check: 'calling' -> true  

### Design Patterns Used
I used 4 patterns in this project.   
Singleton is for GameManager, it keeps only one game on the server.  
Factory creates a new gamefrom pangrams.json, so I don't need to write all logic every time.  
Proxy is for caching dictionary lookups, it remembers words that were already checked.  
Decorator is for logging, it prints which words were checked and doesn't change dictionary code.  

### How to Run the Program
First I start the server with:  
`go run ./cmd/server`

Then in another terminal I start the client with:  
`go run ./cmd/client`

The server response `"Server running on: 50051"` and the client shows seven letters with a center letter.  
Then I can start typing words.

### gRPC Explanation
The gRPC file is game.proto. It has two methods: StartGame and SubmitWord.   
StartGame gives letters, SubmitWord checks the word and returns message and score.   
The project has folders for client, server, game logic, manager and dictionary.   
Dictionary also has proxy and decorator code.  
I used files pangrams.json and words_dictionary.json for data.  

### Materials Used
I used materials from my Distributed System Programming Labs 1-4 and lectures PDFs about Design Patterns and Threading  
I also used the words_dictionary.json and pangrams.json files that were provided by the lecturer.  
I worked in GoLand and used own notes from a class.  
For better learning and understand some labs and lecture materials I used AI (ChatGPT 5) to get explanations in simple language.  
I completed the Gen AI for Students course on the MTU website and already uploaded my certificate on Canvas.  


### Summary
THis project helped me understand how client and server can talk using gRPC.  
I also learned how to use design patterns in real code.  
It was made using Goland and tested on my laptop  