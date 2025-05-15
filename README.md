# PingMe — Real-Time Forum

**PingMe** is a real-time social platform based on our previous forum project — now with WebSocket-powered live messaging. Users can create posts, comment, and chat privately in real time.

---

## 🚀 Features

### 1. User Authentication

* Users must register before accessing the forum.
* **Registration fields include:**

  * Nickname
  * Age
  * Gender
  * First Name
  * Last Name
  * Email
  * Password
* Users can log in using either **nickname** or **email**.
* Logout is accessible from any page.

---

### 2. Posts & Comments

* Create posts categorized by topic.
* View posts in a feed layout.
* Comment on posts (visible only upon post click).

---

### 3. Private Messaging System (WebSocket-Based)

* **User list** shows online/offline users:

  * Sorted by most recent conversation or alphabetically if no messages exist.
* **Real-time messaging** with:

  * Username and timestamp per message.
  * Automatic loading of 10 past messages on chat open.
  * Infinite scroll up loads 10 more messages (optimized with Throttle/Debounce).
* **Live updates** without page refresh.

---

## 🛠️ Tech Stack

| Layer     | Tech Used                             |
| --------- | ------------------------------------- |
| Front-End | HTML, CSS, JavaScript                 |
| Back-End  | Go (Golang), SQLite                   |
| Real-Time | WebSockets (Gorilla WebSocket for Go) |

---

## 📦 Usage

### 1. Installation

Ensure **Go** and **Node.js (for npm)** are installed:

```bash
sudo apt install golang
npm install
```

Then clone the repository:

```bash
git clone https://learn.zone01kisumu.ke/git/ramuiruri/realtime-forum
```

### 2. Running the Project

```bash
cd pingme
make
```

This runs the Go server and launches the real-time forum.

---

## 👥 Authors

* [shfana](https://learn.zone01kisumu.ke/git/shfana)
* [ramuiruri](https://learn.zone01kisumu.ke/git/ramuiruri)

---

## 🤝 Contributing

We welcome contributions!

To contribute:

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/chat-enhancements`)
3. Commit your changes (`git commit -m 'Add typing indicator'`)
4. Push to the branch (`git push origin feature/chat-enhancements`)
5. Open a Pull Request

---

## 🧠 What You'll Learn

* Frontend & backend integration
* Working with WebSockets in Go & JavaScript
* SPA (Single Page Application) architecture using vanilla JS
* Building real-time UIs
* Using SQLite with Go

---

## ❓ Something’s Wrong?

Submit an [issue here](https://learn.zone01kisumu.ke/git/ramuiruri/realtime-forum/issues) and we’ll check it out.

---

Let me know if you’d like badges, screenshots, or GitHub Actions setup for CI/CD added too.
