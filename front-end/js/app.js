// Ong Jia Yuan / S10227735B
// /front-end/js/app.js

document.addEventListener('DOMContentLoaded', function () {
    console.log('DOM fully loaded and parsed');
    
    var registerForm = document.getElementById('registerForm');
    var loginForm = document.getElementById('loginForm');

    if (registerForm) {
        // Function to check password requirements
        document.getElementById('registerPassword').addEventListener('input', function(e) {
            var value = e.target.value;
            var lengthRequirementMet = value.length >= 8;
            var specialCharRequirementMet = /[!@#$%^&*(),.?":{}]/.test(value);

            document.getElementById('lengthRequirement').classList.toggle('requirement-met', lengthRequirementMet);
            document.getElementById('specialCharRequirement').classList.toggle('requirement-met', specialCharRequirementMet);
        });
        registerForm.onsubmit = function (e) {
            e.preventDefault();

            var password = document.getElementById('registerPassword').value;
            var confirmPassword = document.getElementById('confirmPassword').value;

            // Check if passwords match
            var passwordMatchElement = document.getElementById('passwordMatch');
            if (passwordMatchElement) {
                if (password !== confirmPassword) {
                    passwordMatchElement.style.display = 'block';
                    return false; // Stop the form from submitting
                } else {
                    passwordMatchElement.style.display = 'none';
                }
            }

            // If passwords match, proceed with the form submission
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
            })
            .then(response => {
                if (response.ok) {
                    // Assuming the response is not expected to have a body, or it's not important
                    alert("You have successfully registered!");
                    window.location.href = 'login.html';
                } else {
                    // Handle HTTP errors
                    throw new Error(`Server responded with status: ${response.status}`);
                }
            })
            .catch(error => {
                console.error('Error:', error);
            });
        };
    }

    if (loginForm) {
        loginForm.onsubmit = function (e) {
            e.preventDefault();
            var formData = {
                username: document.getElementById('loginUsername').value,
                password: document.getElementById('loginPassword').value
            };
            fetch('http://localhost:5000/api/login', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData)
            })
            .then(response => {
                if (response.ok) {
                    // Assuming the response is not expected to have a body, or it's not important
                    alert("You have successfully login!");
                    window.location.href = '/index.html';
                } else {
                    // Handle HTTP errors
                    throw new Error(`Server responded with status: ${response.status}`);
                }
            })
            .catch(error => {
                console.error('Error:', error);
            });
        };
    }

    // Function to log the current user's details
    function logUserDetails() {
        fetch('http://localhost:5000/api/current-user', {
            method: 'GET',
            credentials: 'include' // Ensure cookies are sent with the request
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Could not fetch user session: ${response.statusText}`);
            }
            return response.json();
        })
        .then(userDetails => {
            console.log('User details:', userDetails);
        })
        .catch(error => {
            console.error('Error:', error);
        });
    }

    // Call the function to log the user details after successful login
    logUserDetails();
});