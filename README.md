# Real Time Forum
Real time forum is a project based on the Forum project which we did previously. 



Its a social platform where people can be able to post there findings, comment on the posts and send private messages.

## Features
### 1. User Authentication

Users must register before using the forum.

Required registration fields:

Nickname

Age

Gender

First Name

Last Name

Email

Password

Users can log in using either their nickname or email with their password.

Users can log out from any page.

### 2. Posts and Comments

Users can create posts, which will be categorized.

Users can comment on posts.

Posts are displayed in a feed format.

Comments are only visible when clicking on a post.

### 3. Private Messaging System

Users can send real-time private messages using WebSockets.

The chat system includes:

A user list showing online/offline users, sorted by last message sent (or alphabetically if no messages exist).

A chat window that loads past messages when selecting a user.

Messages formatted with:

Timestamp

Sender's username

Automatic loading of the last 10 messages with more messages loaded when scrolling up (optimized with Throttle/Debounce).

Real-time updates ensure messages are received instantly without refreshing the page.

## Usage
### 1. Installation
You need to install both golang and javascript in your local machine and this is how you do it.
For the go installation `sudo apt install go` and for javascript `npm install`

After this you need to clone the project
```
git clone https://learn.zone01kisumu.ke/git/wnjuguna/real-time-forum
```

### 2. Running
Change directory to the repository with the project.
`cd real-time-forum`

Then run `make` on the command line. This should run the server and make it visible.

## Authors
This program was built and maintained by
* [wnjuguna](https://learn.zone01kisumu.ke/git/wnjuguna)

* [abrakingoo](https://learn.zone01kisumu.ke/git/abrakingoo)