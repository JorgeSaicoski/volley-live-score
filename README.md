# Volleyball Live Score Tracker

## Project Description
This web application is designed to help parents of a volleyball team stay updated with live match results. When the team travels, a parent on-site will update the scores live, and the rest of the parents can view the results in real-time. Additionally, the app features a blog section where users can post match reports or team updates.

## Features
- **Live Score Updates**: One parent updates the match score in real-time, and others view it.
- **Match and Sets Tracking**: View details of ongoing and past matches, including each set's score.
- **User and Guest Roles**: Users can post blog articles and update match results, while guests can view results and read the blog.
- **Concurrency Control**: Ensure that only one user updates a match score at a time to prevent conflicts.

## Technology Stack
- **Backend**: Golang with GORM for ORM and SQLite as the database.
- **Frontend**: React for building the user interface.
- **Authentication**: JWT for handling user logins and role-based access control.
  
## ER Diagram Overview
- **Guest**: Can view results and blog posts, login to become a user.
- **User**: Can update results, manage matches, and post blog articles.
- **Match**: Contains details about matches including the adversary, date, and win status.
- **Sets**: Stores each set’s score for both teams and determines the winner.
  
## Future Enhancements
- Implement notifications to alert parents when the score is updated.
- Add a comments section to the blog posts for parent interaction.
- Include a dashboard to display past match statistics and player performance.

## How to Run the Project
### Not yet defined.

## Contribution
Feel free to submit issues or pull requests to improve the project. All contributions are welcome!

## License
This project is licensed under the MIT License.

