# Forum

## Overview
This project is a web-based forum that facilitates user communication through posts and comments, supports category-based filtering, and allows users to like or dislike posts and comments. The forum is built using Go, SQLite, and Docker while following best practices in web development, authentication, and database management.

## Features
- **User Authentication**: Registration, login, and session management using cookies.
- **Post & Comment System**: Users can create posts and comment on existing posts.
- **Category Association**: Posts can be categorized for better organization.
- **Like & Dislike System**: Registered users can like or dislike posts and comments.
- **Filtering Mechanism**:
  - By category (acts as subforums)
  - By user-created posts
  - By user-liked posts
- **Database Management**:
  - SQLite as the primary database
  - Implementation of at least one `SELECT`, `CREATE`, and `INSERT` query
  - ER Diagram for structured database design
- **Security Enhancements** (Bonus Tasks):
  - Password encryption using bcrypt
  - Session management using UUID
- **Error Handling**:
  - HTTP status handling
  - Technical error management
- **Dockerization**:
  - Containerized application for better dependency management and deployment

## Technologies Used
- **Go** (Standard Go packages)
- **SQLite** (sqlite3 package)
- **Docker** (for containerization)
- **bcrypt** (for password encryption)
- **UUID** (for unique session management)
- **HTML** (Frontend markup)
- **HTTP** (For client-server communication)

## Installation & Setup
### Prerequisites
Ensure you have the following installed:
- Go
- Docker
- SQLite

### Steps
1. Clone the repository:
   ```sh
   git clone https://learn.zone01kisumu.ke/git/forum
   cd forum
   ```
2. Build and run the application:
   ```sh
   go build -o forum .
   ./forum
   ```
3. Alternatively, run the application directly:
   ```sh
   go run .
   ```
4. Run the application using Docker:
   ```sh
   docker build -t forum-app .
   docker run -p 8080:8080 forum-app
   ```

## Usage
1. **Register/Login** to access the forum.
2. **Create Posts & Comments** (Only for registered users).
3. **Browse Posts** (Public visibility for all users).
4. **Like/Dislike** posts and comments (Only for registered users).
5. **Filter Posts** by categories, created posts, or liked posts.

## Learning Outcomes
This project helps in understanding:
- Web development basics (HTML, HTTP, Sessions, and Cookies)
- Authentication and security best practices
- Database management and SQL queries
- Containerization with Docker
- Structuring a Go-based web application
- Implementing filtering and category management
- Best practices in error handling and testing


## Best Practices Followed
- Secure authentication mechanisms
- Proper error handling
- Database structuring with ER diagrams
- Unit testing for core functionalities

## Contributors
- **@dochiel**
- **@ramuiruri**
- **@wonyango**
- **@najwang**
- **@shfana**

## Issues & Contributions
Found a bug? Have a feature request? Submit an issue or contribute to the project by creating a pull request.

## License
This project is open-source and licensed under [MIT License](LICENSE).
