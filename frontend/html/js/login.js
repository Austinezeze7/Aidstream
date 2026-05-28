document.addEventListener("DOMContentLoaded", () => {

  // --- LOGIN ---
  document.getElementById("loginForm").addEventListener("submit", async (e) => {
    e.preventDefault();
    const email = document.getElementById("loginEmail").value;
    const password = document.getElementById("loginPassword").value;

    try {
      const res = await fetch("/api/login", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email, password })
      });
      const data = await res.json();
      if (data.status === "success") {
        localStorage.setItem("user", email);
        window.location.href = "/donate.html";
      } else {
        alert(data.message || "Login failed");
      }
    } catch (err) {
      console.error(err);
      alert("Server error");
    }
  });

  // --- SIGNUP ---
  document.getElementById("signUpForm").addEventListener("submit", async (e) => {
    e.preventDefault();

    const password = document.getElementById("signUpPassword").value;
    const confirmPassword = document.getElementById("confirmPassword").value;

    if (password !== confirmPassword) {
      alert("Passwords do not match");
      return;
    }

    const payload = {
      firstName:   document.getElementById("firstName").value,
      middleName:  document.getElementById("middleName").value,
      lastName:    document.getElementById("lastName").value,
      email:       document.getElementById("signUpEmail").value,
      password:    password,
      accountType: document.querySelector('input[name="accountType"]:checked').value
    };

    try {
      const res = await fetch("/api/register", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(payload)
      });
      const data = await res.json();
      if (data.status === "success") {
        alert("Account created! Please log in.");
        switchTab("login");
      } else {
        alert(data.message || "Registration failed");
      }
    } catch (err) {
      console.error(err);
      alert("Server error");
    }
  });

});