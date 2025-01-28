**plan to implement the forum project** using the MVC structure, divided among **5 people** over a week.

---

### **General Team Roles**
- **Team Lead (Person 1):** Oversees progress, handles database and Docker setup, and integrates features.
- **Backend Developers (Person 2 & 3):** Build the controllers and models.
- **Frontend Developer (Person 4):** Handles HTML templates and CSS for views.
- **Tester/QA (Person 5):** Writes and executes test cases, ensures bug-free delivery.

---

### **Day 1: Planning & Setup**
#### **Tasks:**
1. **Project Setup (Person 1)**:
   - Initialize a Git repository and set up version control.
   - Create the Go project with `go mod init`.
   - Draft the `Dockerfile` and `docker-compose.yml`.

2. **Database Schema (Person 1)**:
   - Define the database schema in `schema.sql` for users, posts, comments, likes, and categories.
   - Write seed data in `seed.sql`.

3. **Folder Structure (Person 2)**:
   - Set up the MVC folder structure (`controllers/`, `models/`, `views/`, etc.).
   - Create a `db.go` file to initialize SQLite.

4. **HTML Template Drafts (Person 4)**:
   - Create basic `base.html` with a navbar and placeholders for dynamic content.
   - Draft `login.html` and `register.html`.

5. **Testing Plan (Person 5)**:
   - Define test cases for authentication, posts, comments, likes, and filtering.
   - Set up a basic `auth_test.go` file.

---

### **Day 2: Authentication System**
#### **Tasks:**
1. **User Model (Person 2)**:
   - Create `user.go` for user registration and login.
   - Add functions to check if an email is taken and validate credentials.

2. **Authentication Controller (Person 3)**:
   - Build `auth_controller.go` for registering and logging in users.
   - Implement session creation with cookies and expiration.

3. **Encryption (Person 2)**:
   - Use `bcrypt` to hash passwords before saving them in the database.

4. **HTML Templates (Person 4)**:
   - Finish `login.html` and `register.html`.
   - Add basic CSS for the authentication pages.

5. **Testing (Person 5)**:
   - Write unit tests for registration and login.
   - Validate password encryption and cookie expiration.

---

### **Day 3: Posts & Categories**
#### **Tasks:**
1. **Post Model (Person 2)**:
   - Create `post.go` to handle post creation, retrieval, and linking posts to categories.
   - Define functions to save posts and fetch posts by category.

2. **Post Controller (Person 3)**:
   - Build `post_controller.go` to handle routes for creating and displaying posts.
   - Add support for associating posts with multiple categories.

3. **Category Model (Person 2)**:
   - Create `category.go` to manage categories (CRUD operations).
   - Seed initial categories in `seed.sql`.

4. **HTML Templates (Person 4)**:
   - Design `create.html` for new posts.
   - Draft `list.html` to display all posts with filtering options.

5. **Testing (Person 5)**:
   - Write tests for post creation, retrieval, and filtering by category.

---

### **Day 4: Comments**
#### **Tasks:**
1. **Comment Model (Person 2)**:
   - Create `comment.go` to handle comment creation and retrieval.
   - Define functions to associate comments with posts and users.

2. **Comment Controller (Person 3)**:
   - Build `comment_controller.go` to manage adding and displaying comments.

3. **HTML Templates (Person 4)**:
   - Design `create.html` for adding comments.
   - Update `single.html` to display post details with comments.

4. **Frontend Refinement (Person 4)**:
   - Improve styling for posts and comments pages.

5. **Testing (Person 5)**:
   - Write tests for creating and retrieving comments.

---

### **Day 5: Likes/Dislikes**
#### **Tasks:**
1. **Like Model (Person 2)**:
   - Create `like.go` to manage likes and dislikes for posts and comments.

2. **Like Controller (Person 3)**:
   - Build `like_controller.go` to handle routes for liking/disliking posts and comments.

3. **HTML Updates (Person 4)**:
   - Add like/dislike buttons to `single.html` (post details).
   - Display like/dislike counts.

4. **Integration (Person 1)**:
   - Connect the like/dislike system to the database and frontend.

5. **Testing (Person 5)**:
   - Write tests for liking and disliking functionality.

---

### **Day 6: Filtering**
#### **Tasks:**
1. **Filter Controller (Person 2)**:
   - Build `filter_controller.go` to filter posts by categories, created posts, and liked posts.

2. **Post Model Updates (Person 2)**:
   - Add database queries for filtering posts.

3. **Frontend Updates (Person 4)**:
   - Update `list.html` to include filtering options (dropdown for categories).

4. **Integration (Person 1)**:
   - Ensure filtering works seamlessly across controllers, models, and views.

5. **Testing (Person 5)**:
   - Write tests for filtering by categories, user-created posts, and liked posts.

---

### **Day 7: Finalization**
#### **Tasks:**
1. **Dockerization (Person 1)**:
   - Finalize the `Dockerfile` and `docker-compose.yml`.
   - Test the app in a containerized environment.

2. **Bug Fixes (All)**:
   - Fix any remaining bugs in the code.
   - Conduct a final walkthrough of all features.

3. **Styling & Responsiveness (Person 4)**:
   - Ensure the forum is visually appealing and mobile-friendly.

4. **End-to-End Testing (Person 5)**:
   - Perform end-to-end testing of all features (authentication, posts, comments, likes, filtering).

5. **Documentation (Person 1)**:
   - Write a `README.md` with setup instructions, features, and team credits.

---

### **Deliverables**
- **Functional forum** with authentication, posts, comments, likes/dislikes, and filtering.
- **Dockerized app** ready to deploy.
- **Test coverage** for all major features.
- Clean and modular **MVC structure**.
