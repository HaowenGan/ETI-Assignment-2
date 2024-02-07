// Ong Jia Yuan / S10227735B
// /front-end/js/app.js

// Function to update user details in the navigation bar
function updateUserDetailsInNavBar(userDetails) {
    // Find the elements in the DOM
    const userNameSpan = document.querySelector('#userMenuButton .user-name');
    const userTypeSpan = document.querySelector('#userMenuButton .user-type');

    // Set the inner text of the elements to the user details
    if (userNameSpan) userNameSpan.textContent = userDetails.firstName + ' ' + userDetails.lastName;
    if (userTypeSpan) userTypeSpan.textContent = userDetails.usertype.charAt(0).toUpperCase() + userDetails.usertype.slice(1);
}

document.addEventListener('DOMContentLoaded', function () {
    console.log('DOM fully loaded and parsed');

    // Only run the authentication check on pages other than 'login.html' and 'register.html'.
    if (!['/login.html', '/register.html', '/', '/index.html'].includes(window.location.pathname)) {
        console.log('Checking user authentication status...');
        fetch('http://localhost:5000/api/current-user', {
            method: 'GET',
            credentials: 'include' // Ensure cookies are sent with the request.
        })
        .then(response => {
            if (response.status === 401) {
                alert("Please login to view this page!");

                window.location.href = 'login.html';
            } else if (!response.ok) {
                // Other HTTP errors
                throw new Error('Network response was not ok.');
            }
            return response.json(); // If authorized, proceed to handle the response.
        })
        .then(userDetails => {
            // Now that we have the user details, update the navigation bar
            updateUserDetailsInNavBar(userDetails);
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });
    }

    var registerForm = document.getElementById('registerForm');
    var loginForm = document.getElementById('loginForm');
    var reviewForm = document.getElementById('reviewForm');

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
                    window.location.href = '/dashboard.html';
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

    if (reviewForm) {
        reviewForm.onsubmit = function (e) {
            e.preventDefault();
            
            var selectedRating = document.querySelector('input[name="rating"]:checked');
    
            if (!selectedRating) {
                alert('Please select a rating before submitting the review.');
                return;
            }
    
            var formData = {
                // No need to get the username from session, it will be handled server-side
                courseId: parseInt(document.getElementById('courseId').value),
                rating: parseInt(selectedRating.value),
                comment: document.getElementById('comment').value
            };
    
            fetch('http://localhost:8080/api/submit-review', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(formData),
                credentials: 'include' // Include cookies for session management
            })
            .then(response => {
                if (response.ok) {
                    return response.text();
                } else {
                    throw new Error(`Server responded with status: ${response.status}`);
                }
            })
            .then(data => {
                alert('Review submitted successfully');
                console.log(data);
            })
            .catch(error => {
                console.error('Error submitting review:', error);
            });
        };
    }

    // Find the logout button by ID or class and add an event listener
    var logoutButton = document.getElementById('logoutButton');
    if (logoutButton) {
        logoutButton.addEventListener('click', function(e) {
            e.preventDefault();
            fetch('http://localhost:5000/api/logout', {
                method: 'POST',
                credentials: 'include' // Necessary to include the session cookie
            })
            .then(response => {
                if (response.ok) {
                    // Redirect to login page or display a message
                    alert("Successfully logged out! Redirecting to Login page.");
                    window.location.href = 'login.html';
                } else {
                    throw new Error('Logout failed.');
                }
            })
            .catch(error => {
                console.error('Error:', error);
            });
        });
    }

    // Function to log the current user's details
    function logUserDetails() {
        fetch('http://localhost:5000/api/current-user', {
            method: 'GET',
            credentials: 'same-origin' // Ensure cookies are sent with the request
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`Could not fetch user session: ${response.statusText}`);
            }
            return response.json();
        })
        .then(userDetails => {
            console.log('User details:', userDetails);
            return userDetails; // Return the userDetails for further use
        })
        .catch(error => {
            console.error('Error:', error);
        });
    }


    // Call the function to log the user details after successful login
    logUserDetails();
});