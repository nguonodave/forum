const userData = JSON.parse(localStorage.getItem("user"));

let name
name = userData ? userData.username : null;
export const homePageContext = `
    <div class="main-left-div">
        <div id="friendsList">
            <h5>Online Users</h5>
        </div>
        <div id="userList"></div> 
        
    </div
                <!-- Create Post -->
            <div class="create-post">
                <form id="createPostForm">
                    <label id="title"><h4>Add a title</h4></label>
                    <input type="hidden" name="identity" value="post">
                    <input type="text" name="title" id="title" placeholder="Add a title" required style="width: 40%; padding: 6px; border-radius: 4px; outline: none; font-size: medium;">
                    <textarea id="textarea" placeholder="What's on your mind?" name="content" required></textarea>
                    <div class="post-actions">
                        <ul>
                            <div class="choice">
                                <li><input type="checkbox" name="entertainment" id="entertainment" value="entertainment"></li>
                                <label for="entertainment">entertainment</label>
                            </div>
                            <div class="choice">
                                <li><input type="checkbox" name="tech" id="tech" value="tech"></li>
                                <label for="tech">tech</label>
                            </div>
                            <div class="choice">
                                <li><input type="checkbox" name="lifestyle" id="lifestyle" value="lifestyle"></li>
                                <label for="lifestyle">lifestyle</label>
                            </div>
                            <div class="choice">
                                <li><input type="checkbox" name="sports" id="sports" value="sports"></li>
                                <label for="sports">sports</label>
                            </div>
                            <div class="choice">
                                <li><input type="checkbox" name="games" id="games" value="games"></li>
                                <label for="games">games</label>
                            </div>
                        </ul>
                    </div>
                    <button type="submit"><i class="fas fa-paper-plane"></i> Post</button>
                </form>
            </div>

        <div class="posts-section">
                <!-- Posts Feed -->           
        </div>

    <div class="userPanel" >
                <div class="userdetails">
                    <img src="https://ui-avatars.com/api/?name=${name}&background=FEB5B1&color=fff" alt="Avatar" class="avatar">
                    <strong id="userNickName"></strong>
                </div>
                <div class="userfilter">
                    <h3 class="sidetitle">Profile</h3>
                    <ul class="sidelist">
                        <li id="user-info">
                            <div class="user-details">
                            <p><i class="fas fa-user"></i> <span id="user-fullname"></span></p>
                            <p><i class="fas fa-birthday-cake"></i> <span id="user-age"></span></p>
                            <p><i class="fas fa-envelope"></i> <span id="user-email"></span></p>
                            </div>
                         </li>
                    </ul>
                </div>
                <div class="userFunctions">
                    <ul class="sidelist">
                        <li class="nav-link-logout">
                             <a class="nav-link-out"><i class="fas fa-sign-out-alt" id="out"></i>Logout</a>
                        </li>
                    </ul>
                </div>
            </div>
        <div class="createPost">
        <button id="showPostDiv"><i class="fas fa-plus"></i></button>
    </div>
    `