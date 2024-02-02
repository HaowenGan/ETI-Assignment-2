// Ong Jia Yuan / S10227735B
// /front-end/js/app.js

document.addEventListener('DOMContentLoaded', function () {
    var registerForm = document.getElementById('registerForm');
    var loginForm = document.getElementById('loginForm');

    if (registerForm) {
        registerForm.onsubmit = function (e) {
            e.preventDefault();
            var password = document.getElementById('password').value;
            var confirmPassword = document.getElementById('confirmPassword').value;

            // Check if passwords match
            if (password !== confirmPassword) {
                // If passwords do not match, display an error message
                document.getElementById('passwordMatch').style.display = 'block';
                return false; // Stop the form from submitting
            }

            // If passwords match, proceed with the form submission
            document.getElementById('passwordMatch').style.display = 'none';
            var formData = {
                firstName: document.getElementById('firstName').value,
                lastName: document.getElementById('lastName').value,
                email: document.getElementById('email').value,
                username: document.getElementById('username').value,
                password: password
            };
            fetch('http://localhost:5000/api/register', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            }).then(response => response.json())
              .then(data => console.log(data))
              .catch(error => console.error('Error:', error));
        };
    }

    if (loginForm) {
        loginForm.onsubmit = function (e) {
            e.preventDefault();
            var formData = {
                username: document.getElementById('username').value,
                password: document.getElementById('password').value
            };
            fetch('http://localhost:5000/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            }).then(response => response.json())
              .then(data => console.log(data))
              .catch(error => console.error('Error:', error));
        };
    }
});