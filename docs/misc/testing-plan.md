# Group 4: Testing Plan

This document outlines the testing strategy, tools, and quality assurance approach for our project.  

---

## Testing Goals and Scope  
**Unit test**

In our frontend, one important component to be unit tested is the module responsible for making HTTP request (apiRequests.ts). This is to ensure that the communication to the backend is always present and functional. The purpose of this test verify the behavior the API request functions and make sure that the response is being returned correctly to the caller.

**Acceptance test**

Due to time constraints we decided to only test the Login feature. It is the only fully functional feature in this sprint. The test will cover the logging in, registering and logging out aspect of this feature.

This test will verify that:

- A new user can successfully register with unique credentials.
- Registered users can log in with correct credentials.
- Invalid login attempts are handled properly and clearly communicated to the user
- Authenticated users can log out successfully.

## Testing Frameworks and Tools  
### Frontend

**Unit Tests**

**Jest** is a popular framework for React unit testing, which is what the team originally planned to use. Although, we chose to use **Vitest** instead. This is because [Jest is not supported by Vite](https://jestjs.io/docs/getting-started), which is our chosen build tool. 

- Why Vitest?
    - Vitest is a modern alternative to Jest
    - Works seamlessly with Vite without needing extra setup or configuration
    - Built in support for mocks, and measuring coverage
    - Faster to run tests

**Acceptance Test**

We initially chose **React Testing Library** with **Vitest** to perform our acceptance tests. However, due to the limitations of RTL, we had to mock certain parts of the system, such as routing (page navigation), which compromised the accuracy of our tests. By switching to **Cypress + Vitest**, we were able to test the application end-to-end without relying on mocks, resulting in more reliable and realistic acceptance testing.

- Why Cypress?
    - Easy to use
    - End-to-end testing to test real user interaction

### Backend

We chose Go’s built-in testing package because it provides a simple, efficient solution for both unit and integration testing without needing external libraries. Additionally Go also makes mocking straightforward, eliminating the need for third-party mocking frameworks. This approach keeps our test suite lightweight, consistent, and easy to manage.

Go’s testing package

- ease of use
- quick and effective / no need to configure external packages
- provides everything we need for testing (e.g. mocks)


## Test Organization and Structure  



## Coverage Targets  



## Running Tests  
### Frontend Unit Testing

```bash
# Run frontend unit tests
npm run test

# Run test with coverage
npm run test:coverage
```

### Frontend Acceptance Testing

```bash
# Run frontend acceptance test
npm run cy:run

# Run frontend acceptance test using Cypress GUI
npm run cy:open
```

### Backend
```bash
# Run all tests with coverage:
go test ./... -coverpkg=./... -covermode=atomic -coverprofile=coverage.out

# Run specific test files, Example:
go test ./src/internal/test/handler_test.go -v
```

## Reporting and Results  
### Frontend Test Report
Our coverage results can be found in `documentation/tests/sprint_1_test_coverage.png`

[Coverage screenshot](https://github.com/Onyelechie/WorkoutPal-Frontend/blob/main/documentation/tests/sprint_1_test_coverage.png)
This shows that ≥ 80% of lines have been tested.

### Backend Test Report



## Test Data and Environment Setup  



## Quality Assurance and Exceptions  



## Continuous Integration
All tests currently run automatically in a CI pipeline for both backend and frontend. Our frontend code tests itself against the backend in our acceptance tests as well. This ensures that all code is working as expected and removes the need to worry about running the tests last minute and seeing some fail. The frontend tests also prove that new frontend code interacts correctly with the current backend.



