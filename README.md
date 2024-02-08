# ETI-Assignment-2
![Home Page](front-end/images/Index.jpg)

## Database Setup Instructions
1. Run `ETI_ASG2_SQL_Setup.sql`
2. Execute query in mySQL provided by the SQL script

## Website Setup Instructions
Before proceeding, ensure the database have the necessary tables. SQL script is provided.
1. Navigate to `/user-service` directory in terminal
2. Run `go run main.go` in terminal to start user & web services from user-service directory
3. Navigate to `/feedback-service` directory in terminal
4. Run `go run main.go` in terminal for feedback-service to be able to access the Review system in the web front-end
5. Open web browser, type in URL - `localhost:5000` and it should show the index page

## Architecture

The platform follows a microservices architecture to enhance scalability, maintainability, and overall efficiency. Here's an overview of the key components:

1. **User Service:** Manages user authentication, authorization, and profile information.
2. **Payment Service:** Handles secure payment transactions and integrates with various payment gateways.
3. **Feedback Service:** Collects and analyzes user feedback to improve the platform and course offerings.
4. **Course Management Service:** Manages course content, including creation, updates, and version control.
5. **Enrollment Service:** Facilitates user enrollment in courses and ensures a smooth enrollment workflow.



