document.getElementById('regAcc').addEventListener('click', async () => {
    const Nickname = document.getElementById('nickname').value;
    const Password = document.getElementById('password').value;

    if (!Nickname || !Password) {
        alert('Please enter both nickname and password.');
        return;
    }

    try {
        const response = await fetch('/add_user', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ Nickname, Password })
        });

        const result = await response.json();
        alert(result.message || 'User added successfully!');
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to add user.');
    }
});

document.querySelector('button[type="button"]').addEventListener('click', async () => {
    try {
        const response = await fetch('/users');
        const users = await response.json();

        alert(JSON.stringify(users, null, 2)); // Display users in a readable format
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to fetch users.');
    }
});

document.getElementById('updUsr').addEventListener('click', async () => {
    const userId = document.getElementById('idUsr').value;

    if (!userId) {
        alert('Please enter a user ID.');
        return;
    }

    try {
        const response = await fetch(`/user`);
        const user = await response.json();

        if (response.ok) {
            alert(`User found: ${JSON.stringify(user, null, 2)}`);
        } else {
            alert(user.message || 'User not found.');
        }
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to find user.');
    }
});

document.getElementById('dltUsr').addEventListener('click', async () => {
    const userId = document.getElementById('idUsr').value;

    if (!userId) {
        alert('Please enter a user ID.');
        return;
    }

    try {
        const response = await fetch(`/delete-user`, {
            method: 'DELETE'
        });

        const result = await response.json();
        alert(result.message || 'User deleted successfully!');
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to delete user.');
    }
});

document.getElementById('nwpassword').addEventListener('click', async () => {
    const userId = document.getElementById('idUsr').value;
    const newPassword = document.getElementById('editpassword').value;

    if (!userId || !newPassword) {
        alert('Please enter both user ID and new password.');
        return;
    }

    try {
        const response = await fetch(`/update-user`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({password: newPassword})
        });

        const result = await response.json();
        alert(result.message || 'Password updated successfully!');
    } catch (error) {
        console.error('Error:', error);
        alert('Failed to update password');
    }
});
