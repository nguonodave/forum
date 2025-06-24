package controllers

const users_table string = `CREATE TABLE IF NOT EXISTS users(
	id VARCHAR(255) PRIMARY KEY NOT NULL,
	firstname TEXT NOT NULL,
	lastname TEXT NOT NULL, 
	username TEXT NOT NULL, 
	gender TEXT NOT NULL,
	age INTEGER NOT NULL,
	email VARCHAR(255) UNIQUE NOT NULL, 
	password VARCHAR(255) NOT NULL);`
const post_table string = `CREATE TABLE IF NOT EXISTS posts(
	id VARCHAR(255) PRIMARY KEY NOT NULL, 
	userid VARCHAR(255) NOT NULL,
	username VARCHAR(255) NOT NULL,
	title VARCHAR(255) NOT NULL,
	content VARCHAR(255) NOT NULL,
	createdAt VARCHAR(255) NOT NULL,
	category VARCHAR(255), 
	likes INTEGER NOT NULL, 
	dislikes INTEGER NOT NULL, 
	FOREIGN KEY (userid) REFERENCES users(id));`
const comment_table string = `CREATE TABLE IF NOT EXISTS comments(
	id VARCHAR(255) PRIMARY KEY, 
	userid VARCHAR(255) ,
	username VARCHAR(255) ,
	postid VARCHAR(255) ,
	content VARCHAR ,
	createdAt VARCHAR(255) ,
	likes INTEGER,
	dislikes INTEGER ,
	FOREIGN KEY (postid) REFERENCES posts(id),
	FOREIGN KEY (userid) REFERENCES users(id));`
const session_table string = `CREATE TABLE IF NOT EXISTS sessions(
	id VARCHAR(255) PRIMARY KEY,
	sessionToken VARCHAR(255),
	userid VARCHAR(255) ,
	expires VARCHAR(255),
	FOREIGN KEY (userid) REFERENCES users(id));`
const likes_table string = `CREATE TABLE IF NOT EXISTS likes(
	id VARCHAR(255) PRIMARY KEY,
	userid VARCHAR(255),
	postid VARCHAR(255) ,
	reaction VARCHAR(25),
	FOREIGN KEY (postid) REFERENCES posts(id),
	FOREIGN KEY (userid) REFERENCES users(id));`

const comment_likes_table string = `CREATE TABLE IF NOT EXISTS commentlikes(
	id VARCHAR(255) PRIMARY KEY,
	userid VARCHAR(255),
	commentid VARCHAR(255),
	reaction VARCHAR(25),
	FOREIGN KEY (userid) REFERENCES users(id) ,
	FOREIGN KEY (commentid) REFERENCES comments(id) 
);`

const Message_thread_table string = `CREATE TABLE IF NOT EXISTS thread (
	id VARCHAR(255) PRIMARY KEY,
	messages VARCHAR
	);`

const message_table string = `CREATE TABLE IF NOT EXISTS messages (
	id VARCHAR(255) PRIMARY KEY,
	user1Id VARCHAR(255),
	user2Id VARCHAR(255),
	threadId VARCHAR(255),
	FOREIGN KEY (threadId) REFERENCES thread(id)
);`

var Tables = []string{
	users_table,
	Message_thread_table,
	message_table,
	post_table,
	comment_table,
	session_table,
	likes_table,
	comment_likes_table,
}
