# Aurora

    Task management system

    Create/Read/Update a user
    Login a user
    Create/Read a task for a user
    Create comments on a task
    Create/Delete tags on comments

# Routes

    Status - /aurora/api/v1/status

    Create user - /aurora/api/v1/user - POST
    Read user - /aurora/api/v1/user - GET
    Update user - /aurora/api/v1/user - PUT
    
    Login - /aurora/api/v1/login - POST

    Create task - /aurora/api/v1/task - POST
    Read task - /aurora/api/v1/task/{id} - GET

    Create comment - /aurora/api/v1/task/{id}/comment - POST

    Create tag - /aurora/api/v1/task/{id}/comment/{id}/tag - POST
    Delete tag - /aurora/api/v1/task/{id}/comment/{id}/tag/{id} - DELETE

