# Group 4: Testing Plan

This document outlines the testing strategy, tools, and quality assurance approach for our project.  

---

## Testing Goals and Scope  
**Unit test**

In our frontend, one important component to be unit tested is the utility class responsible for sending HTTP requests to the backend (apiRequests.ts). This is to ensure that the communication to the backend is functional and the functions correctly append the passed endpoint to the backend URL. See the [unit test](https://github.com/Onyelechie/WorkoutPal-Frontend/blob/3079c1eea150dfd967f67af3617622ff2c012460/src/utils/__unit_tests__/apiRequests.test.ts) and the [utility class](https://github.com/Onyelechie/WorkoutPal-Frontend/blob/3079c1eea150dfd967f67af3617622ff2c012460/src/utils/apiRequests.ts).

**Acceptance test**

Due to time constraints, we decided to only test the authentication feature as it is the only fully functional feature in this sprint.

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
    - Vitest is a modern alternative to Jest.
    - Works seamlessly with Vite without needing extra setup or configuration.
    - Built in support for mocks, and measuring coverage.
    - Faster to run tests.

**Acceptance Test**

We initially chose **React Testing Library (RTL)** with **Vitest** to perform our acceptance tests. However, due to the limitations of RTL, we had to mock certain parts of the system, such as routing (page navigation). Therefore, we could not perform proper acceptance tests on our system with RTL. By switching to **Cypress + Vitest**, we were able to test the application end-to-end without relying on mocks, resulting in more reliable and realistic acceptance testing.

- Why Cypress?
    - Ability to perform real end-to-end testing compared to RTL.
    - Has test isolation, test retries and automatic waiting to [minimize test flakyness](https://www.cypress.io/app#flake_resistance).
    - Test isolation also helps us create more deterministic tests.
    - Has a GUI compared to RTL. The GUI shows us the exact flow of the tests, making the [debugging process of the tests themselves easier](https://www.cypress.io/app#visual_debugging).

### Backend

We chose Go’s built-in testing package because it provides a simple, efficient solution for both unit and integration testing without needing external libraries. Additionally Go also makes mocking straightforward, eliminating the need for third-party mocking frameworks. This approach keeps our test suite lightweight, consistent, and easy to manage.

Go’s testing package

- Ease of use.
- Quick and effective / no need to configure external packages.
- Provides everything we need for testing (e.g. mocks).


## Test Organization and Structure  
### Backend

**Unit Test**

Tests are located in the same directory as the module they test. This makes it easier to locate and access them quickly. The module is suffixed with `_test` to indicate it is testing a module.

- Example:  `handler/exercise_handler.go`
    - Test file: `handler/exercise_handler_test.go`

**Integration test**

Integration tests are located in `test/e2e`.The same naming convention is used as the unit test to indicate what is being tested.

### Frontend

- [**Unit Test**](https://github.com/Onyelechie/WorkoutPal-Frontend/tree/main/src/utils/__unit_tests__): `src/utils/__unit_tests__`
    - The test files are named `[fileBeingTested].test.ts` which is similar to the backend naming conventions
- [**Acceptance Test**](https://github.com/Onyelechie/WorkoutPal-Frontend/tree/main/cypress/e2e): `cypress/e2e`
    - The test files are name `[fileBeingTested].cy.ts`.


## Coverage Targets  
Our testing coverage goals are designed to ensure reliability, maintainability, and alignment with sprint requirements.

### Frontend Target
- **Logic:** ≥80% coverage.
- **UI:** Login, User Profiles (not met), Dashboard (not met)
  
### Backend Target
- **Logic:** ≥80% line coverage
- **API Layer:** 100% method coverage to target all endpoints
- **Integration:** 100% class coverage

## Running Tests  
### Frontend Unit Testing

```bash
# Run frontend unit tests
npm run test

# Run test with coverage
npm run test:coverage
```

### Frontend Acceptance Testing

Some assumptions, requirements and notes, under 'Acceptance Tests', are outlined in the [README.md](https://github.com/Onyelechie/WorkoutPal-Frontend/blob/3079c1eea150dfd967f67af3617622ff2c012460/README.md).
The tests will fail without the prerequisites and assumptions.

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

[Coverage screenshot](https://github.com/Onyelechie/WorkoutPal-Frontend/blob/3079c1eea150dfd967f67af3617622ff2c012460/documentation/tests/sprint_1_test_coverage.png). This screenshot shows a line coverage of 92.98% for our utility classes `src/utils`, which contains our frontend logic.

### Backend Test Report
[Coverage txt file](/coverage.txt)


## Test Data and Environment Setup

All the setup required to run our tests are explained in both our [Frontend](https://github.com/Onyelechie/WorkoutPal-Frontend/blob/main/README.md) and [Backend](https://github.com/Onyelechie/WorkoutPal-Backend/blob/main/README.md) README.md. There are special requirements outlined in the frontend for acceptance tests, which is already highlighted in the README.md, under 'Acceptance Tests'.

## Quality Assurance and Exceptions

We maintain overall quality by ensuring consistent code reviews. On every pull request we create, we will have someone other than the creator of the pull request comb through the files and ensure that we are following consistent coding standards and using existing constants and utility functions. On the final stage of our sprint, we perform another code review and refactor any inconsistencies with variable naming. Our continuous integration tests also lints through the code and fails when there are unused imports or variables.

The exception for untested components has already been explained in the [sprint 1 worksheet](../sprint1.md), under "Frontend" of "2. Unit / Integration / Acceptance Testing".  

## Continuous Integration

All tests currently run automatically in a CI pipeline for both backend and frontend. Our frontend code tests itself against the backend in our acceptance tests as well. This ensures that all code is working as expected and removes the need to worry about running the tests last minute and seeing some fail. The frontend tests also prove that new frontend code interacts correctly with the current backend.



