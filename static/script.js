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
async function sortUsers() {
    const sort = document.getElementById('sortingInput').value;
    console.log(`Sorting criteria: ${sort}`);

    const url = `${baseUrl}/sort_user?sort=${sort}`;
    console.log(`Request URL: ${url}`);

    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
        });
        console.log(`Response status: ${response.status}`);

        if (response.ok) {
            const users = await response.json();
            console.log("Fetched users:", users);

            const userList = document.getElementById('userList');
            userList.innerHTML = "";

            users.forEach(user => {
                const li = document.createElement('li');
                li.textContent = `ID: ${user.id}, Nickname: ${user.nickname}, Password: ${user.password}`;
                userList.appendChild(li);
            });
        } else {
            console.error("Failed to fetch users. Status:", response.status);
            alert("Failed to fetch users.");
        }
    } catch (error) {
        console.error("Error fetching users:", error);
        alert("Error occurred while fetching users.");
    }
}
let currentPage = 1;
const limit = 4; //

async function fetchUsers(page = 1) {
    const url = `${baseUrl}/users?page=${page}&limit=${limit}`;
    try {
        const response = await fetch(url, {
            method: 'GET',
            headers: { 'Content-Type': 'application/json' },
        });

        if (response.ok) {
            const result = await response.json();
            const users = result.data; // Пользователи
            const totalPages = result.totalPages; // Всего страниц
            currentPage = result.currentPage; // Текущая страница

            // Обновить интерфейс
            const userList = document.getElementById('userList');
            userList.innerHTML = ""; // Очистить текущий список
            users.forEach(user => {
                const li = document.createElement('li');
                li.textContent = `ID: ${user.id}, Nickname: ${user.nickname}`;
                userList.appendChild(li);
            });

            // Обновить кнопки пагинации
            updatePaginationControls(totalPages);
        } else {
            console.error("Failed to fetch users. Status:", response.status);
        }
    } catch (error) {
        console.error("Error fetching users:", error);
    }
}

function updatePaginationControls(totalPages) {
    const prevButton = document.querySelector('#paginationControls button:nth-child(1)');
    const nextButton = document.querySelector('#paginationControls button:nth-child(3)');
    const currentPageSpan = document.getElementById('currentPage');

    prevButton.disabled = currentPage <= 1;
    nextButton.disabled = currentPage >= totalPages;
    currentPageSpan.textContent = currentPage;
}

function nextPage() {
    currentPage++;
    fetchUsers(currentPage);
}

function prevPage() {
    currentPage--;
    fetchUsers(currentPage);
}

// Первоначальная загрузка
fetchUsers();

