# Group 4: Testing Plan

This document outlines the testing strategy, tools, and quality assurance approach for your project.  

**Please visit our [Notion](https://www.notion.so/Testing-Plan-Draft-285b607cf385801ab46cc54968a2c7ef) page to see the draft and the description for each outlines.** Feel free to contribute to atleast one of the outlines listed in this template.

---

## Testing Goals and Scope  


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



## Reporting and Results  


## Test Data and Environment Setup  



## Quality Assurance and Exceptions  



## Continuous Integration



