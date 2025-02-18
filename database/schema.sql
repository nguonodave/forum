-- Enable foreign key constraints
PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS posts (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL,
    content TEXT NOT NULL,
    image_url TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS categories (
    name TEXT UNIQUE NOT NULL
);

CREATE TABLE IF NOT EXISTS post_categories (
    post_id TEXT NOT NULL,
    category TEXT NOT NULL,
    PRIMARY KEY (post_id, category),
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (category) REFERENCES categories(name)
);

CREATE TABLE IF NOT EXISTS comments (
    id TEXT PRIMARY KEY,
    post_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    content TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE TABLE IF NOT EXISTS votes (
    id TEXT PRIMARY KEY, -- Unique ID for the vote
    user_id TEXT NOT NULL, -- ID of the user who voted
    post_id TEXT, -- If vote is for a post
    comment_id TEXT, -- If vote is for a comment
    type TEXT CHECK(type IN ('like', 'dislike')), -- Only allows 'like' or 'dislike'
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP, -- Timestamp when vote was made

    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (post_id) REFERENCES posts(id) ON DELETE CASCADE,
    FOREIGN KEY (comment_id) REFERENCES comments(id) ON DELETE CASCADE,

    -- Ensures that a vote is either for a post OR a comment, never both
    CHECK (
        (post_id IS NOT NULL AND comment_id IS NULL) OR 
        (comment_id IS NOT NULL AND post_id IS NULL)
    )
);

-- Ensure a user can only vote once per post
CREATE UNIQUE INDEX IF NOT EXISTS unique_user_vote_post 
ON votes(user_id, post_id) WHERE post_id IS NOT NULL;

-- Ensure a user can only vote once per comment
CREATE UNIQUE INDEX IF NOT EXISTS unique_user_vote_comment 
ON votes(user_id, comment_id) WHERE comment_id IS NOT NULL;


CREATE TABLE IF NOT EXISTS sessions (
    token TEXT PRIMARY KEY,
    user_id TEXT NOT NULL UNIQUE,
    expires_at DATETIME NOT NULL,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
