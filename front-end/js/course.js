function createCourse() {
    const title = document.getElementById("title").value;
    const content = document.getElementById("content").value;
    const sectionTitle = document.getElementById("sectionTitle").value;
    const sectionContent = document.getElementById("sectionContent").value;

    fetch('http://localhost:8080/courses.html', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            title: title,
            content: content,
            sections: [{
                title: sectionTitle,
                content: sectionContent,
            }],
        }),
    })
    .then(response => response.json())
    .then(data => {
        alert('Course created with ID: ' + data);
    })
    .catch(error => {
        console.error('Error:', error);
    });
}

function updateCourse() {
    const courseId = document.getElementById("updateCourseID").value;
    const title = document.getElementById("updateTitle").value;
    const content = document.getElementById("updateContent").value;
    const sectionTitle = document.getElementById("updateSectionTitle").value;
    const sectionContent = document.getElementById("updateSectionContent").value;

    fetch(`http://localhost:8080/courses/${courseId}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            title: title,
            content: content,
            sections: [{
                title: sectionTitle,
                content: sectionContent,
            }],
        }),
    })
    .then(response => {
        if (response.status === 200) {
            alert('Course updated successfully');
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