// Ong Jia Yuan / S10227735B
// /front-end/js/app.js

document.addEventListener('DOMContentLoaded', function () {
    var registerForm = document.getElementById('registerForm');
    var loginForm = document.getElementById('loginForm');

    if (registerForm) {
        registerForm.onsubmit = function (e) {
            e.preventDefault();
            var formData = {
                firstName: document.getElementById('firstName').value,
                lastName: document.getElementById('lastName').value,
                email: document.getElementById('email').value,
                username: document.getElementById('username').value,
                password: document.getElementById('password').value
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