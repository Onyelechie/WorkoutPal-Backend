# Sprint 3 Worksheet

## Load Testing

## Security Analysis

## Continuous Integration & Deployment

## Group Reflection

## Individual AI / External Resource  Reflections

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

### Max
I used AI extensively for the development of the API Backend. Initial project setup was completed manually to create the model definitions and the different layers of the backend architecture. After which AI was used to flesh out each portion of the app that was essentially repetition. I would manually define the interfaces for each layer and provide them as context to ChatGPT to fill in the implementation. This was done mainly to speed up the process, I knew what needed to be written so I was able to quickly verify if Chat understood correctly and if not I would add additional context to adjust the generation. The other piece where I used chat was writing tests for each layer. These tend to be quite verbose so I would once again feed chat the context of the interface signatures and which kinds of tests I needed written. I would the verify those as well or ask for adjustments. Overall it sped up development of the API significantly. 
