const baseUrl = "http://localhost:8080/api";


async function addUser() {
    console.log("fahof");
    const nickname = document.getElementById('nickname').value;
    const password = document.getElementById('password').value;

    const response = await fetch(`${baseUrl}/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nickname, password })
    });

    if (response.ok) {
        alert("User added successfully!");
        getUsers();
    } else {
        alert("Failed to add user.");
    }
}

async function getUsers() {

    const response = await fetch(`${baseUrl}/users`,{
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },

    });
    if (response.ok) {
        alert("Working");
    }
    else{
        alert("qwer");
    }
    const users = await response.json();
    const userList = document.getElementById('userList');
    userList.innerHTML = "";
    users.forEach(user => {
        const li = document.createElement('li');
        li.textContent = `ID: ${user.id}, Nickname: ${user.nickname}, Password: ${user.password}`;
        userList.appendChild(li);
    });
}

async function updateUser() {
    const id = document.getElementById('userId').value;
    const nickname = document.getElementById('newNickname').value;
    const password = document.getElementById('newPassword').value;

    const response = await fetch(`${baseUrl}/user?id=${id}`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nickname, password })
    });

    if (response.ok) {
        alert("User updated successfully!");
        getUsers();
    } else {
        alert("Failed to update user.");
    }
}

async function deleteUser() {
    const id = document.getElementById('userId').value;

    const response = await fetch(`${baseUrl}/user?id=${id}`, {
        method: 'DELETE'
    });

    if (response.ok) {
        alert("User deleted successfully!");
        getUsers();
    } else {
        alert("Failed to delete user.");
    }
}
async function findUser() {
    // Get the user ID from the input field
    const id = document.getElementById('userId').value;

    if (!id) {
        alert("Please enter a user ID.");
        return;
    }

    try {
        // Fetch the user data from the backend
        const response = await fetch(`${baseUrl}/user?id=${encodeURIComponent(id)}`, {
            method: 'GET',
        });

        if (!response.ok) {
            throw new Error(`Failed to find user. Status: ${response.status}`);
        }

        const userfind = await response.json();

        // Display the user information
        const userList = document.getElementById('singleusr');
        userList.innerHTML = ""; // Clear existing content

        // Create and append user details
        const li = document.createElement('li');
        li.textContent = `ID: ${userfind.id}, Nickname: ${userfind.nickname}, Password: ${userfind.password}`;
        userList.appendChild(li);

        alert("User found!");
    } catch (error) {
        console.error("Error fetching user:", error);
        alert("Failed to find user. Please check the console for details.");
    }
}
