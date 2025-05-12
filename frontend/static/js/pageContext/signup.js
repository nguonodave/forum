export const signupContext = `
<div class="Form-div">
  <form class="loggin siggnup" id="signUpForm">
    <label>Signup</label>

    <input type="text" name="username" class="input" placeholder="Username" required>
    <input type="text" name="firstName" class="input" placeholder="First Name" required>
    <input type="text" name="lastName" class="input" placeholder="Last Name" required>

    <div class="gender-group">
      <label>Gender:</label>
      <label class="gender"><input type="radio" name="gender" value="Male" required> Male</label>
      <label class="gender"><input type="radio" name="gender" value="Female"> Female</label>
    </div>

    <input type="number" name="age" class="input" placeholder="Age" min="1" required>
    <input type="email" name="email" class="input" placeholder="Email Address" required>
    <input type="password" name="password" class="input" placeholder="Password" required>
    <input type="password" name="password2" class="input" placeholder="Confirm Password" required>

    <input type="submit" value="Sign-up">
  </form>

  <p>Have an account? <button id="loginBtn">Login</button></p>

`