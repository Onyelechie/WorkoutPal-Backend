# Sprint 1 – Gym Workout Tracking Application

## 1. Testing Plan
- [Testing Plan for Sprint 1](./misc/testing-plan.md)

## 2. Unit / Integration / Acceptance Testing

### Backend

- API layer: 100% method coverage (every method has at least 1 tested line).

- Logic classes: ≥80% line coverage.

- Integration tests: 100% class coverage, with strong line & method coverage.

### Frontend

- Logic layer (`src/utils`): ≥80% coverage. Direct link: [test coverage](https://github.com/Onyelechie/WorkoutPal-Frontend/blob/main/documentation/tests/sprint_1_test_coverage.png)

- UI tests: 

    - See our [testing plan](https://github.com/Onyelechie/WorkoutPal-Backend/blob/main/docs/misc/testing-plan.md) under 'Acceptance Tests'. For testing the UI that is not yet covered by our acceptance tests, we perform manual testing. When we manual test, we perform black box testing and think of edge cases that may break the UI. After this, we move on to white box testing and look at the logic of the code and attempt to go through every possible branch and get as much coverage as we can.

    - All of our testable logic lives in the [utils folder](https://github.com/Onyelechie/WorkoutPal-Frontend/tree/main/src/utils/__unit_tests__) of our frontend. So far, every utility class is tested except for construction.ts. This file is only temporary and serves as a helper function during development: to show that a feature is not yet implemented. This is not being tested because all it does right now is show an alert with a message “This action is not yet implemented”. Quality of our code is not affected by not testing construction.ts.

    - So far, we haven't had to do anything unusual or unique with our testing approach. We feel like our features and code structure are not complex enough at this point to warrant any deviations from traditional testing approaches (unit testing and acceptance testing).


---
### Coverage Report

[Backend coverage report](../coverage.txt)

[Frontend coverage report](https://github.com/Onyelechie/WorkoutPal-Frontend/blob/main/documentation/tests/sprint_1_test_coverage.png)
---

## 3. Testing Importance

### Unit tests:

1. Authentication

File: src/internal/handler/auth_handler_test.go
Test: TestAuthHandler_Login_OK
Description: 
Verifies a successful login flow
Lowercases incoming email
AuthService is called correctly
Sets access_token cookie with HttpOnly
Returns authenticated user JSON

---

2. User CRUD

File: src/internal/handler/user_handler_test.go
Test: TestUserHandler_UpdateUser_OK
Description: 
Ensures user updates are handled properly
ID is propagated from request context
Service receives correct payload
Returns the updated user

---

3. Validation Utilities

File: src/internal/util/validation_test.go
Test: TestValidateEmail
Description: 
Validates email input behavior
Rejects missing/invalid formats
Ensures early input sanitization

### Integration tests:

Integration Tests:

1. Me Test

File: src/test/e2e/auth_test.go
Test: testEndToEnd_Me
Description:
The /me endpoint returns the authenticated user
Session state is preserved via cookies
Authorization logic works correctly

---

2. Login Test

File: src/test/e2e/auth_test.go
Test: testEndToEnd_Login
Description:
The /login endpoint verifies the user exists in the system
and that the provided email was valid. Returns the user object
along with a set-cookie header

---

3. User Update

File: src/test/e2e/user_test.go
Test: testEndToEnd_UpdateByID
Description:
Creates a user, and then patches it.
Tests both the user creation process
and manipulating existing users
While ensuring the return id stayed consistent

### Acceptance Tests:

Since our automated acceptance tests currently only tests authentication, the story it tests is "**Given** a new user, **when** they register, **then** they must complete all required fields (name, handle, height, weight) within 5 minutes and receive confirmation within 2 seconds". Taken from [sprint0.md](sprint0.md), under user profile feature.
1. Only one test so far which tests the authentication of our system: [auth.cy.ts](https://github.com/Onyelechie/WorkoutPal-Frontend/tree/main/cypress/e2e). The top 3 test functions here would be:
    1. user can register, login and logout - line 47
    2. user cannot have duplicate username and email - line 98
    3. user cannot login with invalid password for an existing account - line 153


## 4. Reproducible Environments

### Group 5: Table track

[Group 5 run instructions](./misc/screenshots/reproducible-environments/group-5-run-instructions.png)

#### Build time [4 mins]

#### Integration test

[Group 5 integration test coverage](./misc/screenshots/reproducible-environments/group-5-integration-test.png)

[Group 5 integration test error](./misc/screenshots/reproducible-environments/group-5-integration-test-error.png)

#### Unit test

[Group 5 unit test coverage](./misc/screenshots/reproducible-environments/group-5-unit-test.png)

#### Running the webserver [5 mins]

[Group 5 home page](./misc/screenshots/reproducible-environments/group-5-web-server.png)

#### Clarity of documentation

- Building through docker and following the instructions is straightforward.
- Instructions to run the frontend is vague. See [screenshot](./misc/screenshots/reproducible-environments/group-5-frontend-instructions.png)
    - No step by step instructions on how to load the webpage using vscode’s live server extension.
    - Had to look up online on how to start the live server.
- Documentation needs to be more detailed and clear. Add step by step instructions to start the frontend application.
- A direct link to the front end start instructions from the README.md can be helpful.

### Summary

- Build time was 4 mins and running the website was 5 mins.
- Building the server is straight forward (using docker).
- Clarity of documentation can be improved. More detailed steps for running the frontend can be beneficial, especially for someone not familiar with VSCode.

## **Sprint 1 Quick Checklist**

- [x]  Link to testing plan.
- [x]  Backend: API 100% method coverage, logic classes ≥80% line coverage.
- [x]  Backend integration tests: 100% class coverage.
- [x]  Frontend logic coverage ≥80% (if applicable).
- [x]  UI test approach described.
- [x]  Missing tests explained (if any).
- [x]  Coverage report included.
- [x]  Top 3 unit tests linked & described.
- [x]  Top 3 integration tests linked & described.
- [x]  Top 3 acceptance tests linked & described.
- [x]  Environment setup for another team tested (screenshots + notes).