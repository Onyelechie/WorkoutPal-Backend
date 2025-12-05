# Sprint 3 Worksheet

## Load Testing

### Describe your load testing environment.
Load testing was performed on a local Windows 11 machine using Apache Jmeter. All test cases were run using 100 threads and 10 loops with a ramp up of 100 seconds, with the exeption of Login where 1000 threads and one loop were used. Test cases are based on commonly expected usage of the frontend.

### Test Cases

Test cases include:

* **Login** - Simulates logging in to WorkoutPal
  
  * Login (POST /login)
  
* **Activity** - Simulates the frontend's home page activity feed, including fetching posts and achievements, as well as creating and commenting on posts
  
  * Get Posts (GET /posts) -
  * Create Post (POST /posts)
  * Get All Unlocked Achievements (GET /achievements/feed)
  * Comment On Post (POST /posts/comment)
  
* **Profile** - Simulates the frontend's profile page, including fetching and updating the user's profile, as well as their followed users and followers
  
  * Get User Profile (GET /me)
  * Update User Profile (PATCH /users/1)
  * Get Followers (GET /user/id/followers)
  * Get Following (GET /user/id/following)
  
* **Achievements** - Simulates the frontend's achievements page, where users view their locked and unlocked achievements
  
  * Get All Achievements (GET /achievements)
  * Get Unlocked Achievements (GET /achievements/unlocked)
  
* **Routines** - Simulates the frontend's routines page, including fetching routines, exercises and schedules as well as creating new routines and scheduling them
  
  * Get Routines (GET /users/1/routines)
  * Get Exercises (GET /exercises)
  * Create Routine (POST /users/1/routines)
  * Get Schedules (GET /schedules)
  * Schedule Routine (POST /schedules)

### Provide the **test report**.

[Test Report](https://html-preview.github.io/?url=https://github.com/Onyelechie/WorkoutPal-Backend/blob/df09112602e5061d45c1cdfb5363f6707a32a52d/docs/load_testing/report/index.html)

### Discuss **one bottleneck** found.
Testing with JMeter was largely successful with low sample sizes, but with high sample sizes JMeter reported a high rate of error. Errors were mostly split between HTTP 500 and 503 errors, which at first were difficult to find the cause. Upon further digging, a bottleneck was found in which the postgresql database was refusing connections as it's maximum amount of clients had been reached while under the heavy server load. Connection pooling through tools such as Pgpool and PGBouncer can be introduced to mitigate this bottleneck.

### State whether you met your **non-functional requirements**.
Our initial non-functional requirements were that "The system must handle at least **100 users with up to 1000 concurrent requests per minute** without significant performance degradation."  
Unfortunately due to the aformentioned bottleneck, we were unable to meet these requirements. These requirements could be achieved through the introduction of connection pooling, as well as more powerful server hardware.

### JMeter files
[.jmx file](https://github.com/Onyelechie/WorkoutPal-Backend/blob/df09112602e5061d45c1cdfb5363f6707a32a52d/docs/load_testing/loadtesting.jmx)

![Load Testing Summary](https://github.com/Onyelechie/WorkoutPal-Backend/blob/df09112602e5061d45c1cdfb5363f6707a32a52d/docs/load_testing/summary.png)

![Load Testing Statistics](https://github.com/Onyelechie/WorkoutPal-Backend/blob/df09112602e5061d45c1cdfb5363f6707a32a52d/docs/load_testing/statistics.png)

![Load Testing Errors](https://github.com/Onyelechie/WorkoutPal-Backend/blob/df09112602e5061d45c1cdfb5363f6707a32a52d/docs/load_testing/errors.png)


## Security Analysis

### Describe your chosen **security analysis tool** and how you ran it.
gosec was chosen as our security analysis tool as it best met our requirements. Installation was performed by running the following command:
```bash
go install github.com/securego/gosec/v2/cmd/gosec@latest
```
Analysis was performed by running gosec with the following command which recursively scans the current directory running all rules:
```bash
gosec ./...
```

### Attach **static analysis report** as an appendix.
[Static Analysis Report](https://github.com/Onyelechie/WorkoutPal-Backend/blob/df09112602e5061d45c1cdfb5363f6707a32a52d/docs/security_analysis/gosec_results.txt)

### Randomly select **5 detected problems** and discuss what you see.
The results of gosec's scan detected 10 lines of code with problems within our backend, which can be boiled down to 3 main problems:

1. Use of net/http serve function that has no support for setting timeouts - Since go's default http.Server settings don't support timeouts, the server is potentially vulnerable to attacks.

2. Unhandled error when calling Rows.Close() - the function has the potential to return an error when failing to close the database, which is left unhandled.

3.  Unhandled error when calling json.NewEncoder(w).Encode(response) - the function has the potential to return an error when failing to encode a json response, which is left unhandled.

### **Required:** Handle or mitigate all **Critical** and **High** vulnerabilities.

No critical or high vulnerabilities were found, nothing to handle.

### If no critical/high vulnerabilities: Discuss **2 other problems** found.
No other problems to discuss.

## Continuous Integration & Deployment
### Continuous Integration

All tests currently run automatically in a CI pipeline for both backend and frontend. Our frontend code tests itself against the backend in our acceptance tests as well. This ensures that all code is working as expected and removes the need to worry about running the tests last minute and seeing some fail. The frontend tests also prove that new frontend code interacts correctly with the current backend.

### Continuous Deployment (CD)

WorkoutPal-Backend uses GitHub Actions to automatically build, push, and deploy Docker images to Azure App Service whenever changes are merged into the main branch.

#### How It Works

- Build: On every push to main, GitHub Actions checks out the code and builds a Docker image of the backend.

- Push: The image is pushed to Docker Hub (ilightlysaltedi/workoutpal-backend:latest).

- Deploy: Azure App Service is configured to pull the latest image automatically, updating the running backend without downtime. deployed backend can be found [here](workoutpal-api-daghb9augub5g9ez.canadacentral-01.azurewebsites.net)

## Reflections

### Design Changes

#### In one paragraph (as a group): What would you change about the **design** of your project now that you’ve been through development?

Better planning of core feature requirements and UI design is what we would like to change. In retrospect, more frequent meetings and planning early in the development cycle could have provided some much needed improvements to our app, whether it be through sketching potential UI designs or brainstorming of ideas to better flesh out our core features.

#### In one paragraph (as a group): What would you change about the **course/project setup**?

The current requirements for the course project felt comfortable. More frequent check-ins might be helpful to keep teams on track for sprint release dates, as well as clearer requirements for testing plans such as load testing and security analysis could benifit teams.

### Individual AI / External Resource  Reflection

### Ebere
During Sprint 3, I leveraged AI extensively to optimize the avatar storage system by implementing PostgreSQL BYTEA binary storage. When I discovered that our base64 TEXT storage was inefficient (33% size overhead), I used ChatGPT to understand PostgreSQL BYTEA implementation patterns and the proper Go interfaces needed. AI helped me design the custom `ByteaData` type that implements `driver.Valuer` and `sql.Scanner` interfaces for seamless PostgreSQL integration. 

The most challenging part was ensuring the conversion between base64 (frontend API compatibility) and binary storage (database optimization) worked correctly. AI assisted me in writing comprehensive test cases to validate binary data conversion accuracy and helped troubleshoot SQL query patterns with proper `::bytea` casting. When repository tests failed due to the BYTEA implementation changes, AI guided me through updating test expectations for binary data handling while maintaining test coverage.

The result was a 33% storage reduction for avatar data with zero frontend changes required, demonstrating how AI can accelerate complex database optimizations while maintaining system compatibility.

[Code Implementation](https://github.com/Onyelechie/WorkoutPal-Backend/blob/Ebere/src/internal/repository/user_repository.go)

### Taren
when setting up the CI/CD pipelines I relied heavily on AI to give me the overall structure of the GitHub actions YML as well as have it give me options for what commands I can use especially for the frontend CD since I had never used GitHub actions to push to azure before, it eventually gave me Azure/static-web-apps-deploy@v1 which worked great but it also failed to tell me I needed to attach my environment variables at build time and that injecting them in azure would not be good enough, so after solving that bug I had to add the environment variables myself. When designing the DB I also used ChatGPT to validate my idea for the database. I cam up with what I thought was good tables and relations and data types using DBML, I then fed that DBML to ChatGPT and it pointed out a few things, such as storing exercise settings in their own table linked to routines and exercises instead of having it only for each user, since users may want to have different settings for different routines. 

### Kurt
During sprint 2, I was creating an add schedule modal where I wanted a `<select>` where a user can select a set of unique routines, no more than once. So I had prompted ChatGPT to create exactly that. The main thing I injected in were the state variables I had created for the payload that will be sent to the backend. This helped speed up the process of implementing this feature imensely as the logic was a little more complex, especially without sample code to follow prior. As I continued to work through the modal, ChatGPT had also introduced me to the use if inputMode="numeric" for `<input>` in html as type="number" was not restrictive enough (it was allowing 'e', negative numbers, and always had a '0' typed out even if the input is supposed to be empty). Because of this, I refactored the rest of my number inputs to match this style and the restrictions I was aiming for.

[Code sample](https://github.com/Onyelechie/WorkoutPal-Frontend/blob/main/src/components/Workouts/Routines/RoutineScheduler/CreateScheduleModal.tsx)

### Christian
In sprint1, my goal was to create a modal dialog that will be used for displaying error messages and actions that require confirmation such as deleting an item. Since this dialog will be used everywhere in the app, I needed to make this code easily accessible and reusable. ChatGPT suggested to use React's Context / Provider pattern. The code it generated was all in one single file. I read the documentation about using the Context API to gain more understanding of its function. Then, I refactored the code by rearranging it and regrouping related code. This custom Dialog API was constantly being maintained and refactored throughout the whole project.

[Dialog Implementation](https://github.com/Onyelechie/WorkoutPal-Frontend/tree/main/src/components/Common/Dialogs)

### Ivory
In terms of code, I mostly used AI for simple boilerplate react components and css styling such as the login and create account cards. I did however use AI extensively to better understand typescript and react, in particular when writing the activity feed component as I struggled to set up the comment and achievement card mapping using typescript's type guard system. I turned to AI when running load testing and security analysis as a means to summarize their results and explain potential causes. This was particularily helpful in the case of security analysis with gosec, as I wasn't familiar with go myself. For testing code withing our frontend, Microsoft Copilot was the convenient choice as it was already well integratied into VScode. For individual files as in gosec and JMeter's results, I used ChatGPT.

### Max
I used AI extensively for the development of the API Backend. Initial project setup was completed manually to create the model definitions and the different layers of the backend architecture. After which AI was used to flesh out each portion of the app that was essentially repetition. I would manually define the interfaces for each layer and provide them as context to ChatGPT to fill in the implementation. This was done mainly to speed up the process, I knew what needed to be written so I was able to quickly verify if Chat understood correctly and if not I would add additional context to adjust the generation. The other piece where I used chat was writing tests for each layer. These tend to be quite verbose so I would once again feed chat the context of the interface signatures and which kinds of tests I needed written. I would the verify those as well or ask for adjustments. Overall it sped up development of the API significantly. 
