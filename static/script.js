const baseUrl = "http://localhost:8080/";


async function addUser() {

    const nickname = document.getElementById('nickname').value;
    const password = document.getElementById('password').value;

    const response = await fetch(`${baseUrl}/add_user`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ nickname, password })
    });
    console.log(JSON.stringify({ nickname, password }));


    if (response.ok) {
        alert("User added successfully!");
        getUsers();
    } else {
        alert("Failed to add user.");
    }
}

async function getUsers() {

    const response = await fetch(`${baseUrl}/all_users`,{
        method: 'GET',
        headers: { 'Content-Type': 'application/json' },

    });
    if (response.ok) {
        alert("Users fetched successfully!");
    }
    else{
        alert("Error in the server");
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

    const response = await fetch(`${baseUrl}/update_user`, {
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

    const response = await fetch(`${baseUrl}/delete_user?id=${id}`, {
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
        const response = await fetch(`${baseUrl}/users/get?id=${id}`, {
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
async function filterUsers() {
    const filter = document.getElementById('filterInput').value;

    // Формируем URL с параметром фильтра
    const url = `${baseUrl}/filter_user?filter=${filter}`;

    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
        });

        if (response.ok) {
            const users = await response.json();
            const userList = document.getElementById('userList');
            userList.innerHTML = ""; // Очистить список

            users.forEach(user => {
                const li = document.createElement('li');
                li.textContent = `ID: ${user.id}, Nickname: ${user.nickname}, Password: ${user.password}`;
                userList.appendChild(li);
            });
        } else {
            alert("Failed to fetch users: " + response.status);
        }
    } catch (error) {
        console.error("Error fetching users:", error);
        alert("Error occurred while fetching users.");
    }
}
async function sortUsers(){
    const sort = document.getElementById('sortInput').value;
    const url = `${baseUrl}/sort_user?sort=${sort}`;
    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
        });
        if (response.ok) {
            const users = await response.json();
        }
    }
    catch (error) {
        console.error("Error fetching users:", error);
        alert("Error occurred while fetching users.");
    }

}

document.addEventListener("DOMContentLoaded", () => {
    const userId = "67890"; // Укажите ID пользователя

    fetch(`/interactions?userId=${userId}`)
        .then(response => {
            if (!response.ok) {
                console.error("Failed to fetch interactions");
                return [];
            }
            return response.json();
        })
        .then(interactions => {
            const historyList = document.getElementById('interactionHistory');
            historyList.innerHTML = ""; // Очистка списка перед добавлением новых элементов
            interactions.forEach(interaction => {
                const li = document.createElement('li');
                li.textContent = `${interaction.Type} at ${interaction.Details.Timestamp}`;
                historyList.appendChild(li);
            });
        })
        .catch(err => console.error("Error loading interactions:", err));
});

