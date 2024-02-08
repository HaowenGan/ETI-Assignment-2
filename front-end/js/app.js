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

// Function to update user details in the main content area
function updateUserDetailsInMainContent(userDetails) {
    // Assuming your HTML structure and class names, target the elements for user's name and user type
    const userNameElement = document.querySelector('#main-content .text-2xl.font-semibold');
    const userTypeElement = document.querySelector('#main-content .text-gray-600');

    // Update the text content of these elements with userDetails
    if (userNameElement) userNameElement.textContent = userDetails.firstName + ' ' + userDetails.lastName;
    if (userTypeElement) userTypeElement.textContent = userDetails.usertype.charAt(0).toUpperCase() + userDetails.usertype.slice(1); // Capitalize the first letter
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
            updateUserDetailsInMainContent(userDetails);
        })
        .catch(error => {
            console.error('There has been a problem with your fetch operation:', error);
        });
    }

    var registerForm = document.getElementById('registerForm');
    var loginForm = document.getElementById('loginForm');
    var reviewForm = document.getElementById('reviewForm');

    if (registerForm) {
        var registerPassword = document.getElementById('registerPassword');
        var lengthRequirement = document.getElementById('lengthRequirement');
        var specialCharRequirement = document.getElementById('specialCharRequirement');
        var confirmPassword = document.getElementById('confirmPassword');
        var passwordMatchElement = document.getElementById('passwordMatch');

        // Check if all elements exist
        if (!registerPassword || !lengthRequirement || !specialCharRequirement || !confirmPassword || !passwordMatchElement) {
            console.error('One or more elements do not exist in the DOM.');
            return; // Stop the execution if elements are missing
        }

        registerPassword.addEventListener('input', function (e) {
            var value = e.target.value;
            var lengthRequirementMet = value.length >= 8;
            var specialCharRequirementMet = /[!@#$%^&*(),.?":{}]/.test(value);

            lengthRequirement.classList.toggle('requirement-met', lengthRequirementMet);
            specialCharRequirement.classList.toggle('requirement-met', specialCharRequirementMet);
        });

        registerForm.onsubmit = function (e) {
            e.preventDefault();

            var password = registerPassword.value;
            var confirmPasswordValue = confirmPassword.value;
            var errors = [];

            // Check password requirements and accumulate error messages
            if (password.length < 8) {
                errors.push("Password must be at least 8 characters long.");
            }
            if (!/[!@#$%^&*(),.?":{}]/.test(password)) {
                errors.push("Password must contain at least one special character.");
            }
            if (password !== confirmPasswordValue) {
                errors.push("Passwords do not match.");
                passwordMatchElement.style.display = 'block'; // Show the password match error
            } else {
                passwordMatchElement.style.display = 'none'; // Hide the password match error
            }

            // If there are any errors, show an alert and stop form submission
            if (errors.length > 0) {
                alert("Please correct the following errors before submitting:\n\n" + errors.join("\n"));
                return; // Stop the form from submitting
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
                credentials: 'include',
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
                    alert('Invalid credentials. Please check if you have entered the correct credentials.');
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
    
            fetch('http://localhost:5001/api/submit-review', {
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
                window.location.href = 'ViewReviews.html';
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

    // Fetch reviews from the API and populate the table
    fetch('http://localhost:5001/api/get-reviews', {
        method: 'GET',
        credentials: 'include', // Ensure cookies are sent with the request
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        } else if (!response.headers.get("content-type")?.includes("application/json")) {
            throw new Error("Not a JSON response");
        }
        return response.json();
    })
    .then(reviews => {
        console.log(reviews);
        const tableBody = document.querySelector('#reviewsTable tbody');

        reviews.forEach(review => {
            console.log(review); // Log the review object to inspect its properties
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${review.id}</td>
                <td>${review.courseId}</td>
                <td>${review.rating}</td>
                <td>${review.comment}</td>
                <td><button onclick="editReview(${review.id})">Edit</button></td>
                <td><button onclick="deleteReview(${review.id})">Delete</button></td>
            `;
            tableBody.appendChild(row);
        });
    })
    .catch(error => console.error('Error fetching reviews:', error));

    // Find the profile update form by its ID or class (assuming the form has an id="profileUpdateForm")
    var profileUpdateForm = document.getElementById('profileUpdateForm');

    if (profileUpdateForm) {
        profileUpdateForm.onsubmit = function(e) {
            e.preventDefault(); // Prevent the default form submission
    
            // Object to hold formData
            var formData = {};
    
            // Add only non-empty fields to formData
            if (document.getElementById('first-name').value) formData.firstName = document.getElementById('first-name').value;
            if (document.getElementById('last-name').value) formData.lastName = document.getElementById('last-name').value;
            if (document.getElementById('email').value) formData.email = document.getElementById('email').value;
    
            // Send AJAX request with formData
            fetch('http://localhost:5000/api/update-profile', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify(formData),
                credentials: 'include' // Include cookies for session management
            })
            .then(response => {
                if (response.ok) {
                    alert("Profile changed successfully! Logging out to reflect the changes.");
                    return fetch('http://localhost:5000/api/logout', {
                        method: 'POST',
                        credentials: 'include' // Include cookies for session management
                    });
                } else {
                    alert("Failed to update profile. Please try again.");
                }
            })
            .then(response => {
                if (response.ok) {
                    // Redirect to login page or show success message
                    window.location.href = 'login.html';
                } else {
                    // handle error
                }
            })
            .catch(error => {
                console.error('Error:', error);
            });
        };
    }

    var passwordChangeForm = document.getElementById('passwordChangeForm');

    if (passwordChangeForm) {
        passwordChangeForm.addEventListener('submit', function(e) {
            e.preventDefault();

            // Validate the new passwords match
            var newPassword = document.getElementById('new-password').value;
            var confirmNewPassword = document.getElementById('confirm-new-password').value;

            if (newPassword !== confirmNewPassword) {
                alert("New passwords do not match.");
                return;
            }

            var currentPassword = document.getElementById('current-password').value;

            // Send AJAX request to the backend to update the password
            fetch('http://localhost:5000/api/change-password', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json'
                },
                body: JSON.stringify({
                    currentPassword: currentPassword,
                    newPassword: newPassword
                }),
                credentials: 'include' // Include cookies for session management
            })
            .then(response => {
                if (response.ok) {
                    alert("Password updated successfully! Logging out for security reasons.");
                    return fetch('http://localhost:5000/api/logout', {
                        method: 'POST',
                        credentials: 'include' // Include cookies for session management
                    });
                } else {
                    // handle error
                }
            })
            .then(response => {
                if (response.ok) {
                    // Redirect to login page or show success message
                    window.location.href = 'login.html';
                } else {
                    // handle error
                }
            })
            .catch(error => {
                console.error('Error:', error);
            });
        });
    }
});

// Define the deleteReview function
function deleteReview(reviewId) {
    // Make a DELETE request to the API with the reviewId
    fetch(`http://localhost:5001/api/delete-review/${reviewId}`, {
        method: 'DELETE',
        credentials: 'include', // Ensure cookies are sent with the request
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        // Assuming successful deletion, you might want to handle it based on your requirements
        console.log(`Review with ID ${reviewId} deleted successfully`);

         // Show a pop-up indicating successful deletion
         window.alert('Review deleted successfully');
        
         // Reload the page after deletion
         window.location.reload();

        // You can call your fetchReviews function here or update the UI accordingly
    })
    .catch(error => console.error('Error deleting review:', error));
}

// Function to handle editReview button click
function editReview(reviewid) {
    // Redirect to the edit-review.html page with the review ID
    window.location.href = `EditReview.html?id=${reviewid}`;
}

const urlParams = new URLSearchParams(window.location.search);
const reviewidString = urlParams.get('id');
const reviewid = parseInt(reviewidString, 10);

// Fetch the existing review details based on the review ID
fetch(`http://localhost:5001/api/get-review/${reviewid}`, {
    method: 'GET',
    credentials: 'include', // Ensure cookies are sent with the request
})
.then(response => {
    if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
    } else if (!response.headers.get("content-type")?.includes("application/json")) {
        throw new Error("Not a JSON response");
    }
    return response.json();
})
.then(review => {
    // Populate the textarea with the existing comment
    document.getElementById('editedComment').value = review.comment;
})
.catch(error => console.error('Error fetching review details:', error));


// Function to save the edited review
function saveEditedReview() {
    const editedComment = document.getElementById('editedComment').value;
    const selectedRating = document.querySelector('input[name="rating"]:checked');

    if (!selectedRating) {
        console.error('Please select a rating');
        return;
    }

    const ratingValue = parseInt(selectedRating.value, 10); // Convert to integer
    console.log(reviewid)

    // Make a PATCH request to update the review's comment and rating
    fetch(`http://localhost:5001/api/edit-review/${reviewid}`, {
        method: 'PATCH',
        headers: {
            'Content-Type': 'application/json',
        },
        credentials: 'include', // Ensure cookies are sent with the request
        body: JSON.stringify({
            Id: reviewid,
            comment: editedComment,
            rating: ratingValue,
        }),
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        // Assuming successful update, you might want to handle it based on your requirements
        console.log('Review updated successfully');
        alert("Review Updated Successfully!")
        window.location.href = `ViewReviews.html`
        // Redirect back to the main reviews page or handle navigation as needed
    })
    .catch(error => console.error('Error updating review:', error));
}


