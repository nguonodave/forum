-- Seed Users
INSERT INTO users (username, email, password) VALUES
                                                  ('alice', 'alice@example.com', 'password123'),
                                                  ('bob', 'bob@example.com', 'password123');

-- Seed Categories
INSERT INTO categories (id, name) VALUES
                                  (1, 'Technology'),
                                  (2, 'Health'),
                                  (3, 'Lifestyle');

-- Seed Posts
INSERT INTO posts (user_id, title, content) VALUES
                                                ('1', 'Go Programming Basics', 'This is a post about Go programming basics...'),
                                                ('2', 'Healthy Living Tips', 'This is a post about living a healthy life...');

-- Seed Comments
INSERT INTO comments (post_id, user_id, content) VALUES
                                                     ('1', '2', 'This is a comment on the Go Programming Basics post'),
                                                     ('2', '1', 'This is a comment on the Healthy Living Tips post');

-- Seed Votes
INSERT INTO votes (user_id, post_id, type) VALUES
                                               ('1', '1', 'like'),
                                               ('2', '2', 'dislike');

-- Seed Session
INSERT INTO sessions (token, user_id, expires_at) VALUES
                                                      ('sample-token-1', '1', '2025-12-31 23:59:59'),
                                                      ('sample-token-2', '2', '2025-12-31 23:59:59');
