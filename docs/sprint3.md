# Sprint 3 Worksheet

## Load Testing

## Security Analysis

## Continuous Integration & Deployment

## Group Reflection

## Individual AI / External Resource  Reflections

### Ebere

### Taren

### Kurt
During sprint 2, I was creating an add schedule modal where I wanted a `<select>` where a user can select a set of unique routines, no more than once. So I had prompted ChatGPT to create exactly that. The main thing I injected in were the state variables I had created for the payload that will be sent to the backend. This helped speed up the process of implementing this feature imensely as the logic was a little more complex, especially without sample code to follow prior. As I continued to work through the modal, ChatGPT had also introduced me to the use if inputMode="numeric" for `<input>` in html as type="number" was not restrictive enough (it was allowing 'e', negative numbers, and always had a '0' typed out even if the input is supposed to be empty). Because of this, I refactored the rest of my number inputs to match this style and the restrictions I was aiming for.

[Code sample](https://github.com/Onyelechie/WorkoutPal-Frontend/blob/main/src/components/Workouts/Routines/RoutineScheduler/CreateScheduleModal.tsx)

### Christian
In sprint1, my goal was to create a modal dialog that will be used for displaying error messages and actions that require confirmation such as deleting an item. Since this dialog will be used everywhere in the app, I needed to make this code easily accessible and reusable. ChatGPT suggested to use React's Context / Provider pattern. The code it generated was all in one single file. I read the documentation about using the Context API to gain more understanding of its function. Then, I refactored the code by rearranging it and regrouping related code. This custom Dialog API was constantly being maintained and refactored throughout the whole project.

### Ivory
