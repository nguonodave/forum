const signInForm = document.getElementById("signInForm")
const container = document.querySelector(".container");
const header = document.querySelector(".header");

import { welcomePage } from "./pageContext/welcome.js"
import { headerContext } from "./pageContext/header.js"
import { signinContext } from "./pageContext/login.js"
import { signupContext } from "./pageContext/signup.js"
import { homePageContext } from "./pageContext/home.js"
import { loadPosts, posts } from "./loadPosts.js"
import { notify } from "./pageContext/notification.js"

export let ws;

const handleRoutes = () => {

    // console.log("route", window.location.pathname);

    switch (window.location.pathname) {
        case "/login":
            loadSignIn(false);
            break;
        case "/":
            const user = JSON.parse(localStorage.getItem("user"));
            if (user) {
                if (!ws || ws.readyState !== WebSocket.OPEN) {
                    const ws = new WebSocket("ws://localhost:8080/ws");

                    ws.onopen = () => {
                        console.log("WebSocket connected");
                        if (user) {
                            ws.send(JSON.stringify({ type: "user", name: user.username, id: user.id }));
                        }
                    };

                    loadHomePage(ws, user);
                } else {

                    loadHomePage(ws, user);
                }           
            } else {
                loadWelcomePage(false);
            }
            break;
        case "/signup":
            loadSignUp(false);
            break;
        default:
            loadWelcomePage(false);
            break;
    }
};

const loadWelcomePage = (pushState = true) => {
    if (pushState) history.pushState({}, "", "/welcome");
    header.innerHTML = headerContext;
    container.innerHTML = welcomePage;
    attachWelcomePageListeners();
};

const loadSignIn = (pushState = true) => {
    if (pushState) history.pushState({}, "", "/login");
    header.innerHTML = headerContext;
    container.innerHTML = signinContext;
    attachListeners();
};

const loadSignUp = (pushState = true) => {
    console.log("Lod sjajs")
    if (pushState) history.pushState({}, "", "/signup");
    header.innerHTML = headerContext;
    container.innerHTML = signupContext;
    attachListeners();
};

document.addEventListener("DOMContentLoaded", handleRoutes);
window.onpopstate = handleRoutes;

const attachWelcomePageListeners = () => {
    const login = document.getElementById("loginBtn");
    const signUp = document.getElementById("signUpBtn");


    if (login) login.addEventListener("click", loadSignIn);
    if (signUp) signUp.addEventListener("click", loadSignUp);
};

const attachListeners = () => {
    const signUpBtn = document.getElementById("signUpBtn");
    const login = document.getElementById("loginBtn");
    const loginform = document.querySelector("#loginForm");
    const signupform = document.getElementById("signUpForm");
    
    if (signUpBtn) {
        signUpBtn.addEventListener("click", loadSignUp);
    }
    if (login) {
        login.addEventListener("click", loadSignIn);
    }

    if (loginform) loginform.addEventListener("submit", handleLogin);
    if (signupform) signupform.addEventListener("submit", handleSignup);
};

const handleLogin = async (e) => {
    e.preventDefault();

    const name = document.querySelector("input[name='email']").value;
    const password = document.querySelector("input[name='password']").value;


    try {
        const response = await fetch("/login", {
            method: 'POST',
            body: JSON.stringify({ name, password }),
            headers: { 'Content-Type': 'application/json' },
        });

        const data = await response.json();
        if (response.status != 200) {

            notify(data.message, "var(--danger-color)");

        } else if (response.status == 200) {
            history.pushState({}, "", "/");
            ws = new WebSocket(`ws://localhost:8080/ws`)
            ws.onopen = () => {
                const user = { type: "user", name: data.username, id: data.id };
                ws.send(JSON.stringify(user));
            };

            localStorage.setItem("sessionToken", data.session_token);
            localStorage.setItem("user", JSON.stringify(data));

            loadHomePage(ws, data);

    }
        
    } catch (err) {
        console.error("Login failed:", err);
    }

};

function updateFriendList(message, data) {
    const user = JSON.parse(localStorage.getItem("user"));
    let friendsUnsorted = JSON.parse(localStorage.getItem("friends") || "[]");

    const otherUser = message.data.sender === user.id ? message.data.receiver : message.data.sender;

    const existingIndex = friendsUnsorted.findIndex(friend => friend.user === otherUser);

    if (existingIndex !== -1) {
        friendsUnsorted[existingIndex].date = message.data.created;
    } else {
        friendsUnsorted.push({ user: otherUser, date: message.data.reated });
    }

    localStorage.setItem("friends", JSON.stringify(friendsUnsorted));
    const storedusers = JSON.parse(localStorage.getItem("userList") || "[]");
    updateUserList(storedusers, data);
}

let messageCount = 0;
const notificationElem = document.querySelector(".newMessage-notification");

if (notificationElem) {
    
    messageCount = parseInt(notificationElem.textContent.trim(), 10) || 0;
}

const updateUserList = (msg, data) => {
    localStorage.setItem("userList", JSON.stringify(msg));
    const homeUser = JSON.parse(localStorage.getItem("user"));
    document.getElementById("userList").innerHTML = "";

    const allusers = msg.data.filter((u) => u.id !== data.id);
    
    const unsortedusers = Array.from(new Map(allusers.map(user => [user.id, user])).values());
    
    const friendsUnsorted = JSON.parse(localStorage.getItem("friends") || "[]");
  
    const sortedFriends = friendsUnsorted.sort((a, b) => new Date(b.date) - new Date(a.date));
    
    const friends = [];
    const nonFriends = [];
    
    unsortedusers.forEach(user => {
        const isFriend = sortedFriends.some(friend => friend.user === user.id);
        if (isFriend) {
            friends.push(user);
        } else {
            nonFriends.push(user);
        }
    });

    const sortedNonFriends = nonFriends.sort((a, b) => a.name.localeCompare(b.name));
    
    const finalUserList = Array.from(new Map(
        [...sortedFriends.map(friend => friends.find(u => u.id === friend.user)).filter(Boolean), ...sortedNonFriends]
        .map(user => [user.id, user]) 
    ).values());

    if (msg.data.length !== 1) {
        document.getElementById("userList").innerHTML = finalUserList.map(user => `
            <li>
                <button>
                    <div class="status-indicator"></div>
                    <span class="user-name" data-userid="${user.id}">${user.name} online</span>
                </button>
                <form id="privateMessageForm" class="privateMessageForm">
                    <input type="hidden" name="receiver" value="${user.id}">
                    <input class="privateMessageTextarea" placeholder="chat.." name="message" required>
                    <input type="hidden" name="username" value="${homeUser.username}">
                    <input type="hidden" name="userid" value="${homeUser.id}">
                    <input type="submit" value="Send" class="btn">
                </form>
            </li>
        `).join("");
    } else {
        document.getElementById("userList").innerHTML = `<li>1 Online</li>`;
    }
    
}

const handleWS = (ws, data) => {
    ws.onmessage = (event) => {
        try {
            const msg = JSON.parse(event.data);
            const user = JSON.parse(localStorage.getItem("user"));
            if (msg.type === "new_post") {
                appendNewPost(msg.data);
                notify("New Post", "var(--success-color)");
            }

            if (msg.type === "user_list") {
                updateUserList(msg, data);
            }                  

            if (msg.type === "reaction") {
                const postId = msg.data.data.id;

                if (msg.data.data.reaction === "like") {
                    const form = document.querySelector(`#postLikesForm input[name='postid'][value='${postId}']`)?.closest("#postLikesForm");
        
                    if (form) {
                        const btn = form.querySelector("#likes-count");
                        const dislikeBtn = form.closest(".post-actions").querySelector("#dislike-count");
        
                        btn.textContent = msg.data.data.Likes;
        
                        if (msg.data.data.Dislikes >= 0) {
                            dislikeBtn.textContent = msg.data.data.Dislikes;
                        }
                    } else {
                        console.warn(`No form found for post ID ${postId}`);
                    }
                }
                if (msg.data.data.reaction === "dislike") {
                    const form = document.querySelector(`#postDislikeBtn input[name='postid'][value='${postId}']`)?.closest("#postDislikeBtn");
        
                    if (form) {
                        const btn =form.querySelector("#dislike-count");
                        const likeBtn = form.closest(".post-actions").querySelector("#likes-count");
                        
                        btn.textContent = msg.data.data.Dislikes;
        
                        if (msg.data.data.Likes >= 0) {
                            likeBtn.textContent = msg.data.data.Likes;
                        }
                    } else {
                        console.warn(`No form found for post ID ${postId}`);
                    }
                }
            }

            if(msg.type === "comment") {
                addNewComment(msg.data);
            }

            if (msg.type === "message") {
                const chatContainer = document.getElementById("chat");
                const newMessage = document.createElement("div");
                const messageClass = msg.data.receiver === user.id ? "receiver" : "sender";
                newMessage.classList.add("message", messageClass);
    
                const senderName = msg.data.receiver === user.id ? msg.data.sendername : "You";
            
                newMessage.innerHTML = `
                    <div class="meta">${senderName} • ${msg.data.created}</div>
                    ${msg.data.message}
                `;
            
                chatContainer.appendChild(newMessage);
            
                chatContainer.scrollTop = chatContainer.scrollHeight;

                const chatboxid = document.querySelector(".chat-wrapper .input-container .receiverId");
                if (msg.data.receiver === user.id) {
                    chatboxid.value = msg.data.sender;
                    const notificationElem = document.querySelector(".newMessage-notification");
                    messageCount++;
                    notificationElem.textContent =  messageCount.toString();
                }

                fetchUserMessages(user.id);
                updateFriendList(msg, data);
                
            }

            if (msg.type === "typing") {
                const users = document.querySelectorAll("[data-userid]");
                let found = false;
                
                users.forEach((user) => {
                
                    if (user.dataset.userid === msg.senderid) {
                        user.textContent = `${msg.sendername} is typing...`;
                        user.style.color = "green"
                        found = true;
                    }
                });
                
                if (!found) {
                    console.warn("User not found for typing:", msg.senderid);
                }
                
            }
            

        } catch (err) {
            console.error("Error parsing WebSocket message:", err);
        }
    };
};

const loadHomePage = (ws, data) => {
    let name , firstName, lastName, age, email, fullUserName;
    header.innerHTML = headerContext;
    document.querySelector("#user-inbox").innerHTML = ` <i class="fas fa-envelope" ><span class="newMessage-notification">0</span></i>`
    document.querySelector(".nav-link").innerHTML = `<div id="profile"><i class="fa-regular fa-user"></i></div>`;
    container.innerHTML = homePageContext  ;
    document.getElementById("userNickName").innerHTML = `${data.username}`;
    const searchContainer =document.querySelector(".search-wrapper")
    searchContainer.innerHTML = `
                        <div class="search-input-container">
                        <i class="fas fa-search search-icon"></i>
                        <input type="text" placeholder="Search ...">
                    </div>
    `

    name = data ? data.username : null;
    firstName = data ? data.firstname : "";
    lastName = data ? data.lastname : "";
    age = data ? data.age : "";
    email = data ? data.email : "";
    fullUserName = `${firstName}  ${lastName}`;

    document.getElementById("user-fullname").textContent = fullUserName;
    document.getElementById("user-email").textContent = email;
    document.getElementById("user-age").textContent = age;
    handleWS(ws, data);           
    loadPosts(ws).then(() => {console.log(`posts have been loaded`)});
    createPost(ws);
    loadHomePageListeners(ws);
    handleMessages(ws);
    fetchUserMessages(data.id).then(() => {console.log("messages have been fetched")});

};

const loadHomePageListeners = (ws) => {

    document.addEventListener("focusin", (event) => {
        if (event.target.matches("input[name='message']")) {
            const form = event.target.closest("form"); 
            if (form) {
                const username = form.querySelector("input[name='username']");
                const receiver = form.querySelector("input[name='receiver']");
                const userid = form.querySelector("input[name='userid']");
                const data = {
                    type: "typing",
                    sendername: username.value,
                    receiverid : receiver.value,
                    senderid: userid.value,
                };
                if (username) {
                    if(ws && ws.readyState === WebSocket.OPEN) {
                        ws.send(JSON.stringify(data));
                    }
                } else {
                    console.log("Username not found.");
                }
            }
        }
    });
    
    const logoutBtn = document.querySelector(".nav-link-logout");

    if (logoutBtn) {
        logoutBtn.addEventListener("click", async () => {
            try {
                const res = await fetch("/logout", {
                    credentials: "include",
                });
                
                if (!res.ok) {
                    localStorage.clear();
                    window.location.href = "/welcome";
                    throw new Error(`Logout failed: ${res.status}`);
                }
                
                const data = await res.json();
                console.log("/logout",data.message);
    
                localStorage.clear();
                window.location.href = "/welcome";
            } catch (error) {
                console.error("Logout error:", error.message);
            }
        });
    } else {
        console.warn("Logout button not found!");
    }
    

    document.addEventListener('click', function (e) {
        const button = e.target.closest('.comment-btn');
        if (!button) return;
    
        e.preventDefault();
        const post = button.closest('.post');
        const commentsSection = post.querySelector('.comments-section');
    
        if (!commentsSection) return;
    
        if (commentsSection.classList.contains('open')) {
            commentsSection.style.maxHeight = '0px';
            commentsSection.style.opacity = '0';
            commentsSection.classList.remove('open');
        } else {
            commentsSection.style.maxHeight = commentsSection.scrollHeight + 'px';
            commentsSection.style.opacity = '1';
            commentsSection.classList.add('open');
    
        }
        const addCommentForms = document.querySelectorAll(".addCommentForm");

    addCommentForms.forEach((form) => {
            form.addEventListener('submit', async (e) => {
                e.preventDefault();
                const comment = form.querySelector("input[name='comment']").value.trim()
                const identity = "comment";
                const postId = form.querySelector("input[name='postId']").value.trim();
                form.querySelector("input[name='comment']").value = "";
                createComment(ws, comment, identity, postId);
                
            })
        })
    });

            const displayProfile = document.querySelector(".user-profile-nav-link");
            const userPanel = document.querySelector(".userPanel");
            userPanel.style.display = "none";

            displayProfile.addEventListener("click", (e) => {
                e.stopPropagation();
                userPanel.style.display = (userPanel.style.display === "none") ? "block" : "none";
            });
            

            
            document.addEventListener("click", () => {
                userPanel.style.display = "none";
            });


            const showPostDiv = document.getElementById("showPostDiv");
            const create_Post = document.querySelector(".create-post");
            create_Post.style.display = "none";

        
    showPostDiv.addEventListener("click", () => {
                if (create_Post.style.display === "none") {
                    create_Post.style.display = "block";
                } else {
                    create_Post.style.display = "none";
                }
            });

    window.addEventListener("click", (e) => {
                const create_Post = document.querySelector(".create-post");
                const showPostDiv = document.getElementById("showPostDiv");
            
                // Check if the clicked element is NOT the toggle button or inside the create post form
                if (
                    create_Post.style.display === "block" &&
                    !create_Post.contains(e.target) &&
                    e.target !== showPostDiv
        
                ) {
                    create_Post.style.display = "none";
                }
        
                if (
                    userPanel.style.display === "block" &&
                    !userPanel.contains(e.target) 
                    // e.target !== displayProfile
                ){ userPanel.style.display = "none";}
            });

    
            const userList = document.querySelector("#userList");
        
    userList.addEventListener('click', (e) => {
                const button = e.target.closest("button");
            if (button) {
                const siblings = Array.from(e.target.parentNode.children);
                const form = button.nextElementSibling;
        
                if (form && form.classList.contains("privateMessageForm")) {``
                    form.style.display = (form.style.display === "none" || !form.style.display) ? "flex" : "none";
                } else {
                    console.error("No form found next to the button");
                }
            }
            });
};

const handleMessages = (ws) => {
    const privateMessageForm = document.querySelector("#userList");
    const user = JSON.parse(localStorage.getItem("user"));
    if (privateMessageForm) {
        privateMessageForm.addEventListener("submit", async (e) => {
                e.preventDefault();
                const form = e.target.closest("form");
                if (form) {
                    const msg = form.querySelector("input[name='message']").value.trim();
                    const receiver = form.querySelector("input[name='receiver']").value.trim();
                    const sender = user.id;
                    const sendername = user.username;


                    if (msg && receiver) {
                        const response = await fetch("/privateMessage", {
                            method: "POST",
                            headers: {"Content-Type" : "application/json"},
                            body: JSON.stringify({
                                message: msg,
                                receiver: receiver,
                                sender: sender,
                                sendername: sendername,
                                seen: false,
                                created: new Date().toLocaleString('en-US', {hour12: false})
                            })
                        });
                        
                        const resBody = await response.json();

                        if (response.status == 200) {
                            ws.send(JSON.stringify({"type" : "message", "data" : resBody.data}))
                            form.querySelector("input[name='message']").value = "";
                            form.style.display = "none" 
                            fetchUserMessages(user.id);
                        }

                    } else {
                        notify("Please fill out both the message and receiver fields.");
                    }
                }
        });
    } else {
        console.log("Private Message Form not Found.")
    }

    const inbox = document.querySelector("#user-inbox");
    const hideMessageThread = document.querySelector(".messageThread");
    const chatWrapper = document.querySelector(".chat-wrapper");
    chatWrapper.style.display = "none";
    hideMessageThread.style.display = "none";
    if (inbox) {
        inbox.addEventListener('click', (e) => {
            e.stopPropagation()
            if (hideMessageThread.style.display === "none") {
                hideMessageThread.style.display = "flex";
            } else {
                hideMessageThread.style.display = "none";
                chatWrapper.style.display = "none";
            }
        });
    }

    // Prevent closing when clicking inside the message thread or chat wrapper
        hideMessageThread.addEventListener("click", (event) => {
            event.stopPropagation();
        });
        chatWrapper.addEventListener("click", (event) => {
            event.stopPropagation();
        });

        // Close when clicking outside
        document.addEventListener("click", () => {
            hideMessageThread.style.display = "none";
            chatWrapper.style.display = "none";
        });

    document.addEventListener('click', () => {
        hideMessageThread.style.display = "none";
        chatWrapper.style.display = "none";
    });

    hideMessageThread.addEventListener('click', (e) => {
        const thread = e.target.closest(".thread");
        if (!thread) return;
    
        e.preventDefault();
        const threads = document.querySelectorAll(".thread");

        if (chatWrapper.style.display === "none" || !thread.classList.contains("active")) {
            threads.forEach((t) => t.classList.remove("active"));
            thread.classList.add("active");
            chatWrapper.style.display = "flex";
            if (messageCount > 0) {
                const notificationElem = document.querySelector(".newMessage-notification");
                messageCount--; 
                notificationElem.textContent = messageCount.toString(); 
            }
        }
    
        const chatContainer = document.getElementById("chat");
        if (chatContainer) {
            chatContainer.scrollTo({
                top: chatContainer.scrollHeight,
                behavior: 'smooth'
            });
        }

        const messageInput = document.getElementById("messageInput");
        if (messageInput) messageInput.focus();
    });
    
    

    const sendBtn = document.querySelector(".input-container #sendBtn");
    const messageInputElem = document.querySelector(".input-container #messageInput");
    
    if (sendBtn && messageInputElem) {
        const sendMessage = async () => {
            const chatReceiverId = document.querySelector(".input-container .receiverId").value;
            const messageInput = messageInputElem.value.trim();
    
            if (!messageInput) return;
    

            const response = await fetch("/privateMessage", {
                method: "POST",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({
                    message: messageInput,
                    receiver: chatReceiverId,
                    sender: user.id,
                    sendername: user.username,
                    seen: false,
                    created: new Date().toLocaleString('en-US', { hour12: false })
                })
            });
    
            const resBody = await response.json();
    
            if (response.status === 200) {
                ws.send(JSON.stringify({ "type": "message", "data": resBody.data }));
                messageInputElem.value = "";
                fetchUserMessages(user.id);
            }
    
        };

        sendBtn.addEventListener("click", sendMessage);

        messageInputElem.addEventListener("keydown", (event) => {
            if (event.key === "Enter" && !event.shiftKey) {
                event.preventDefault();
                sendMessage();
                fetchUserMessages(user.id);
            }
        });
    
    } else {
        console.log("sendBtn and messageInput not found");
    }   
}

const throttle = (func, limit) => {
    let inThrottle;
    return function () {
        if (!inThrottle) {
            func.apply(this, arguments);
            inThrottle = true;
            setTimeout(() => (inThrottle = false), limit);
        }
    };
};

let friendList = [];
const fetchUserMessages = async (userid) => {
    try {
        const res = await fetch(`/getusermessages?userId=${userid}`);
        if (!res.ok) throw new Error("Network response was not ok");

        const messages = await res.json();
        const messageThreadDiv = document.querySelector(".messageThread");

        if (!messages.data || !Array.isArray(messages.data)) {
            messageThreadDiv.textContent = "You Have No Messages";
            return;
        }

        messageThreadDiv.replaceChildren();

        messages.data.forEach((message) => {

            console.log("msg",message);
            let parsed;
            try {
                parsed = JSON.parse(message);
                console.log("parsed->", parsed);
                if (!Array.isArray(parsed) || parsed.length === 0) return;
            } catch (error) {
                console.error("Error parsing message JSON:", error);
                return;
            }

            const latestMessage = parsed[parsed.length - 1];
            console.log("Latest Message: ", latestMessage);
            const thread = document.createElement("div");
            thread.classList.add("thread");
            console.log("lastmsg by jones",latestMessage);
            console.log("name is ", latestMessage.sendername)
            const senderName = latestMessage.sender === userid ? "You" : latestMessage.sendername|| "Unknown User";
            thread.textContent = `${senderName}: ${latestMessage.message}`;

            if (latestMessage.receiver !== userid) {
                friendList.push({ user: latestMessage.receiver, date: latestMessage.created });
            } else {
                friendList.push({ user: latestMessage.sender, date: latestMessage.created });
            }
            
            document.addEventListener("DOMContentLoaded", loadMessages(parsed, userid))
            thread.addEventListener("click", () => loadMessages(parsed, userid));

            messageThreadDiv.appendChild(thread);
        });
    } catch (error) {
        console.error("Fetch error:", error);
        document.querySelector(".messageThread").innerHTML = "<p>Failed to load messages. Please try again.</p>";
    }

};

const loadMessages = (parsed, userid) => {
    try {
        const chatContainer = document.getElementById("chat");
        const chatboxid = document.querySelector(".chat-wrapper .input-container .receiverId");

        chatContainer.replaceChildren();
        const fragment = document.createDocumentFragment();
        const reversedParsed = parsed.reduce((acc, obj) => [obj, ...acc], []);
        const chunkedMessages = chunkArray(reversedParsed);
        let count = 0;

        const appendMessages = (index) => {
            if (index < 0 || index >= chunkedMessages.length) return;

            const messageFragment = document.createDocumentFragment();
            chunkedMessages[index].reduce((acc, obj) => [obj, ...acc], []).forEach((msg) => {
                console.log("message now>>", msg)
                const newMessage = document.createElement("div");
                newMessage.classList.add("message", msg.sender === userid ? "sender" : "receiver");

                const displayName = (msg.sender === userid) ? "You" : (msg.sendername ?? "Unknown User");

                chatboxid.value = msg.sender === userid ? msg.receiver : msg.sender;

                const metaDiv = document.createElement("div");
                metaDiv.classList.add("meta");
                metaDiv.textContent = `${displayName} • ${msg.created}`;

                const textDiv = document.createElement("div");
                textDiv.classList.add("text");
                textDiv.textContent = msg.message;

                newMessage.appendChild(metaDiv);
                newMessage.appendChild(textDiv);
                messageFragment.appendChild(newMessage);

            });

            chatContainer.prepend(messageFragment);
            
            // Maintain scroll position after loading new messages
            if (index !== count) {
                chatContainer.scrollTop = chatContainer.scrollHeight / 4;
            }
        };
        // console.log("FRIENDS:", friendList);
        localStorage.setItem("friends", JSON.stringify(friendList));

        appendMessages(count);

        chatContainer.addEventListener("scroll", throttle(() => {
            if (chatContainer.scrollTop === 0) {
                count++;
                if (count >= chunkedMessages.length) {
                    console.log("No more messages to load.");
                    return;
                }
                appendMessages(count);
            }
        }, 1));

    } catch (error) {
        console.error("Error displaying chat thread:", error);
    }
};

function chunkArray(array, size = 10) {
    const result = [];
    for (let i = 0; i < array.length; i += size) {
        result.push(array.slice(i, i + size));
    }
    return result;
}

const handleSignup = async (e) => {
    e.preventDefault();

    const username = document.querySelector("input[name='username']").value;
    const firstName = document.querySelector("input[name='firstName']").value;
    const lastName = document.querySelector("input[name='lastName']").value;

    // gender resolves to male or female depending on what the user selects
    const gender = document.querySelector('select.input').value.trim().toLowerCase();

    const age = parseInt(document.querySelector("input[name='age']").value);
    const email = document.querySelector("input[name='email']").value;
    const password = document.querySelector("input[name='password']").value;

    function validatePassword(password) {
        const lengthCheck = password.length >= 8;
        if (!lengthCheck) return "Password must be at least 8 characters long.";
        return true;
    }

    const passwordValidation = validatePassword(password);
    if (passwordValidation !== true) {
        notify(passwordValidation); // Show error as a popup
        return;
    }


    console.log(`--- signup payload ${username} ${firstName} ${lastName} ${age} ${email} ${password}----`);

    try {
        const res = await fetch("/signup", {
            method: "POST",
            body: JSON.stringify({
                Username: username, 
                FirstName: firstName, 
                LastName: lastName, 
                Age: age,
                Gender: gender, 
                Email: email, 
                Password: password
            }),
            headers: {"Content-Type": "application/json"},
        });

        const data = await res.json();
        if (res.status === 201) {
            notify(data.message, "var(--success-color)");
            loadSignIn();
        } else {
            notify(data.message, "var(--danger-color)");
            console.log(data.message);
        }
        
    } catch (error) {
        console.log("SignUp Failed", error);
    }
};

const createPost = (ws) => {
    const postForm = document.getElementById("createPostForm");
        // Clone the form to remove old event listeners
    const newPostForm = postForm.cloneNode(true);
    postForm.parentNode.replaceChild(newPostForm, postForm)

    newPostForm.addEventListener("submit", async (e) => {
        e.preventDefault();
        
                const title = document.querySelector("input[name='title']").value;
                const content = document.querySelector("#textarea").value;
                const choices = document.querySelectorAll(".post-actions .choice input:checked");
                const identity = "post";
                const categories = Array.from(choices).map(choice => choice.value).join(" ");
        
                // console.log(title, content, choices.length, categories);
        
                // Reset checked checkboxes
                choices.forEach(choice => choice.checked = false);
        
                try { 
                    const response = await fetch("/createPost", {
                        method: "POST",
                        headers: { "Content-Type": "application/json" },
                        body: JSON.stringify({ title, content, categories, identity }),
                    });
        
                    const respData = await response.json();

                    // console.log(response.status)
                    
                    if (response.status === 201) {
                        notify("Post created successfully!", "var(--success-color)");
                        const user = JSON.parse(localStorage.getItem("user"));
                        const newPost = {
                            type: "new_post",
                            id: respData.postId,
                            name : user.username,
                            title: title,
                            content: content,
                            categories: categories,
                            timestamp: new Date().toISOString(),
                        };
                        // console.log("RECEIVED RESPONSE: ", newPost);

                        if (ws && ws.readyState === WebSocket.OPEN) {
                            ws.send(JSON.stringify(newPost));
                            // console.log("WEBSOCKET IS OPEN TO SEND POST OVER ", newPost);
                            // createPost(ws);
                        }
                    } else if (response.status ===  401) {
                        notify(respData.message || "Failed to create post", "var(--danger-color)");
                        setTimeout( async() => {
                            const res = await fetch("/logout");
                            if (res.status === 200) {
                                localStorage.clear();
                                window.location.href = "/welcome";
                            }
                        },3000)
                        
                    }
        
                } catch (error) {
                    console.error("Error creating post:", error);
                    notify("An error occurred. Try again.", "var(--danger-color)");
                }
        
                // Reset input fields **after** submission
                document.querySelector("input[name='title']").value = "";
                document.querySelector("#textarea").value = "";
                const create_Post = document.querySelector(".create-post");
                create_Post.style.display = "none";
                console.log("Post Submitted");
        });
};

const appendNewPost = (post) => {
    // console.log("AVAILABLE POSTS LENGTH: ", posts.length);
    let pageData = "";
    const postsContainer = document.querySelector(".posts-section");// Adjust selector as needed
    const postDiv = document.createElement("div");
    postDiv.classList.add("post");
    postDiv.setAttribute("postid", post.id)

    pageData += `
        <div class="post-header">
            <div class="avatar">
            <i class="fa-regular fa-user"></i>
            </div>
            <div class="post-info">
                <h4>${post.name}</h4>
                <span class="post-time">${post.timestamp}</span>
                <span class="post-category">Category(ies): ${post.categories}</span>
            </div>
        </div>
        <div class="post-content">
            <h6>${post.title}</h6>
            <p>${post.content}</p>
        </div>
        <div class="post-actions">
            <form id="postLikesForm">
                <input type="hidden" name="postid" value="${post.id}">
                <input type="hidden" name="reaction" value="like">
                <button type="submit" class="like-btn"><i class="fas fa-thumbs-up"></i><span id="likes-count">0</span></button>
            </form>
            <form id="postDislikeBtn">
                <input type="hidden" name="postid" value="${post.id}">
                <input type="hidden" name="reaction" value="dislike">
                <button type="submit" class="like-btn"><i class="fas fa-thumbs-down"></i> <span id="dislike-count">0 <span> </button>
            </form>
            <button class="comment-btn"><i class="fas fa-comment"></i> 0</button>
        </div>
        <div class="comments-section">
            <div class="add-comment">
                <div class="avatar">
                    <img src="https://avatars.dicebear.com/api/adventurer/johndoe.svg" alt="image">
                </div>
                <form class="addCommentForm">
                    <input type="hidden" name="identity" value="comment">
                    <input type="hidden" name="postId" value="${post.id}">
                    <input type="text" placeholder="Write a comment..." name="comment" required>
                    <button type="submit"><i class="fas fa-paper-plane"></i></button>
                </form>
            </div>
        </div>
`;

    postDiv.innerHTML = pageData;
    if (posts.length < 1) {
        postsContainer.textContent = "";
        postsContainer.prepend(postDiv);
        window.location.reload();
    }
    postsContainer.prepend(postDiv);
};

const createComment = async (ws, comment, identity, postId) => {

    try {
        const res = await fetch("/createComment", {
            method: "POST",
            headers: {"Content-Type" : "application/json"},
            body: JSON.stringify({comment, identity, postId}),
        });

        const resBody = await res.json();

        const user = JSON.parse(localStorage.getItem("user"));
        const newComment = {
            type: "comment",
            commentId: resBody.commentId,
            name : user.username,
            postId: resBody.postId,
            context: resBody.comment,
            timestamp: resBody.time,
            comments_count: resBody.comments_count,
        };

        if (res.status == 201) {
            if(ws && ws.readyState == WebSocket.OPEN) {
                ws.send(JSON.stringify(newComment));
            }
            console.log(resBody);
        } else {
            console.log(resBody);
        }

    } catch(error) {
        console.log("Unable to Send Comment: ", error);
    }
};

const addNewComment = (comment) => {
    console.log("RECEIVED COMMENTS CONTENT FROM WS", comment.postId);

    const commentsContainer = document.querySelector(`.post[postid="${comment.postId}"] .comments-section`);
    const commentCountUpdate = document.querySelector(`.post[postid="${comment.postId}"] .comment-btn`);

    if (!commentsContainer) {
        console.error("Comments container not found for post:", comment.post_Id);
        return;
    }

    const newComment = document.createElement("div");
    newComment.classList.add("comment");

    newComment.innerHTML = `
        <div class="avatar">
        <i class="fa-regular fa-user"></i>
        </div>
        <div class="comment-content">
            <div class="comment-header">
                <h4>${comment.name}</h4>
                <span class="comment-time">${comment.timestamp}</span>
            </div>
            <p class="comment-text">${comment.context.comment}</p>
         <!---   <form class="commentLikeForm">
                <input type="hidden" name="commentid" value="${comment.id}">
                <input type="hidden" name="reaction" value="like">
                <button class="like-btn"><i class="fas fa-thumbs-up"></i>0</button>
            </form>
            <form class="commentDislikeForm">
                <input type="hidden" name="commentid" value="${comment.id}">
                <input type="hidden" name="reaction" value="dislike">
                <button class="like-btn"><i class="fas fa-thumbs-down"></i>0</button>
            </form> --->
        </div>
    `;

    commentsContainer.prepend(newComment);
    commentCountUpdate.innerHTML = `<i class="fas fa-comment"></i>${comment.comments_count}`;
};