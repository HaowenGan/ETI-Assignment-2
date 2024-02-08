function createCourse() {
    const title = document.getElementById("title").value;
    const content = document.getElementById("content").value;
    const price = parseFloat(document.getElementById("price").value);
    const sectionTitle = document.getElementById("sectionTitle").value;
    const sectionContent = document.getElementById("sectionContent").value;

    const courseData = {
        title: title,
        content: content,
        price: price,
        sections: [{
            title: sectionTitle,
            content: sectionContent,
        }],
    };

    console.log('Create Course Data:', courseData);

    fetch('http://localhost:8080/courses', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(courseData),
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.json();
    })
    .then(data => {
        alert('Course created with ID: ' + data);
        listCourses(); // Refresh the course list after creating a new course
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

function updateCourse() {
    const courseId = document.getElementById("updateCourseID").value;
    const title = document.getElementById("updateTitle").value;
    const content = document.getElementById("updateContent").value;
    const price = parseFloat(document.getElementById("updatePrice").value);
    const sectionTitle = document.getElementById("updateSectionTitle").value;
    const sectionContent = document.getElementById("updateSectionContent").value;

    const updatedCourseData = {
        title: title,
        content: content,
        price: price,
        sections: [{
            title: sectionTitle,
            content: sectionContent,
        }],
    };

    fetch(`http://localhost:8080/courses/${courseId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(updatedCourseData),
    })
    .then(response => {
        if (response.status === 200) {
            alert('Course updated successfully');
            listCourses(); // Refresh the course list after updating a course
        } else {
            alert('Course not found');
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

function deleteCourse() {
    const courseId = document.getElementById("deleteCourseID").value;

    fetch(`http://localhost:8080/courses/${courseId}`, {
        method: 'DELETE',
    })
    .then(response => {
        if (response.status === 200) {
            alert('Course deleted successfully');
        } else {
            alert('Course not found');
        }
    })
    .catch(error => {
        console.error('Error:', error);
    });
}
function listCourses() {
    fetch('http://localhost:8080/courses', {
        method: 'GET',
        headers: {
            'User-Type': 'student', 
        },
        credentials: 'same-origin', // Ensure cookies are sent with the request
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`Could not fetch courses: ${response.statusText}`);
        }
        return response.json();
    })
    .then(data => {
        const courseListBody = document.getElementById('courseListBody');
        courseListBody.innerHTML = '';

        // Sort courses in descending order based on price
        data.sort((a, b) => b.Price - a.Price);

        data.forEach(course => {
            const row = document.createElement('tr');
            row.innerHTML = `<td>${course.ID}</td><td>${course.Title}</td><td>${course.Content}</td><td>${course.Price}</td>`;
            courseListBody.appendChild(row);
        });
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

function showHideElements(userType) {
    const createForm = document.getElementById("createCourseForm");
    const updateForm = document.getElementById("updateCourseForm");
    const deleteForm = document.getElementById("deleteCourseForm");

    if (userType === "admin") {
        // Show all forms for admins
        createForm.style.display = "block";
        updateForm.style.display = "block";
        deleteForm.style.display = "block";
    } else if (userType === "student") {
        // Hide update and delete forms for students
        createForm.style.display = "none";
        updateForm.style.display = "none";
        deleteForm.style.display = "none";
    } else {
        // Handle other user types or unauthorized access
        alert("Unauthorized access.");
        window.location.href = "/login"; // Redirect to login page or handle as needed
    }
}

function logUserDetails() {
    fetch('http://localhost:5000/api/current-user', {
        method: 'GET',
        credentials: 'same-origin',
    })
    .then(response => {
        if (!response.ok) {
            throw new Error(`Could not fetch user session: ${response.statusText}`);
        }
        return response.json();
    })
    .then(userDetails => {
        console.log('User details:', userDetails);
        showHideElements(userDetails.usertype);
        return userDetails; // Return the userDetails for further use
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

logUserDetails();
listCourses();