# Sprint 0 – Gym Workout Tracking Application

## 1. Repository Links
- [Frontend GitHub Repository](https://github.com/Onyelechie/WorkoutPal-Frontend)
- [Backend GitHub Repository](https://github.com/Onyelechie/WorkoutPal-Backend)

## 2. Slides & Presentation
- [Proposal Slides (PDF)](./proposal_slides.pdf) 

## 3. Project Summary & Vision
Our project is a **Gym Workout Tracking Application** designed to help users plan, log, and visualize their fitness progress in a simple and motivating way. The application allows users to create personal profiles, design customized workout routines, and track sessions with built-in timers for rest and recovery.  

The system addresses the common issue of users losing motivation by providing **streaks**, **time-series graphs**, and **social features** that encourage accountability. The vision is to create an engaging platform that combines **fitness tracking** with **gamification** and **social interaction**, enabling users to build healthy habits and stay consistent in their fitness journey.  

---

## 4. Features & User Stories

### Core Features

#### 1. User Profiles
**User Story:**  
As someone who likes to stay organized, I want to store my profile information (name, handle, avatar, bio, height, weight, units, timezone) so that my workout logs are consistent.  

**Acceptance Criteria:**  
- **Given** a new user, **when** they register, **then** they should be able to create a profile with personal details.  
- **Given** an existing user, **when** they update their profile, **then** the changes should be saved and displayed immediately.  
- **Given** a user’s profile, **when** it is viewed, **then** it should display their stats and achievements.  

---

#### 2. Workout Schedule and Routine
**User Story:**  
As a busy student, I want to schedule my workout sessions and stick to a routine so that I can stay consistent throughout the week.  

**Acceptance Criteria:**  
- **Given** a user, **when** they create a new workout plan, **then** it should be saved to their account.  
- **Given** a user with a plan, **when** they view their calendar, **then** their workouts should appear on the correct dates.  
- **Given** a user, **when** they edit or delete a workout plan, **then** the calendar should reflect the changes.  

---

#### 3. Session Runner + Timer
**User Story:**  
As someone who likes strictly timed workout sessions, I want to keep track of how long my sets and overall sessions last.  

**Acceptance Criteria:**  
- **Given** a user starts a session, **when** they hit "start," **then** the global session timer should begin.  
- **Given** a user is in a session, **when** they log a set, **then** the app should store weight, reps, RPE, and notes.  
- **Given** a user finishes a set, **when** the rest timer starts, **then** the countdown should be visible and notify the user when time is up.  
- **Given** a session in progress, **when** the user pauses or resumes, **then** the timers should adjust accordingly.  

---

#### 4. User Relationships
**User Story:**  
As an extroverted gym user, I want to connect with new people and share my workouts with them.  

**Acceptance Criteria:**  
- **Given** a user searches for another user, **when** they send a friend request, **then** the request should appear in the recipient’s notifications.  
- **Given** a user has friends, **when** they complete a workout, **then** their friends should see updates from them.  

---

#### 5. Activity Feed
**User Story:**  
As someone who loses motivation easily, I want to see what other people are doing so that I can get inspiration and push myself to go to the gym.  

**Acceptance Criteria:**  
- **Given** a user’s feed, **when** they open it, **then** they should see recent workouts and posts from friends.  
- **Given** a friend completes a workout, **when** the feed updates, **then** the workout should be displayed.  

---

#### 6. Goal Tracking & Targets
**User Story:**  
As someone who is goal oriented, I want to keep track of set goals for weight, body-fat, and lifts for target dates.  

**Acceptance Criteria:**  
- **Given** a user, **when** they create a new goal, **then** it should be saved with a target date.  
- **Given** a user with goals, **when** they log progress, **then** the system should update their progress towards those goals.  
- **Given** a user’s dashboard, **when** it is viewed, **then** goals should display current progress vs. target.  

---

### Non-Functional Requirement
- The system must handle at least **100 users with up to 1000 concurrent requests per minute** without significant performance degradation.  

---

### Additional Features / Stretch Goals
1. **Custom Workouts**  
   - Advanced users can add personal workouts to their routines.  
2. **AI Custom Workout Builder**  
   - The system suggests personalized workout routines using AI based on user history, goals, and preferences.  

---

## 5. Initial Architecture
![Architecture Diagram](./architecture-diagram.png)

- **Frontend:** React.js  
  - Rich UI components, intuitive data control, large ecosystem.  
- **Backend:** Go  
  - Strong typing, lightweight, high concurrency handling.  
- **Database:** MySQL  
  - Reliable relational database, sufficient for structured user and workout data.  

This architecture works well because it provides **clear separation of concerns**: React manages the client-side experience, Go handles scalable API requests efficiently, and MySQL ensures reliable storage of structured workout data. The stack is lightweight yet scalable enough to support our MVP and future growth.  

---

## 6. Work Division & Coordination
We will divide work by **feature ownership**:  
- One person responsible for **Profiles**, another for **Session Runner**, another for **Schedules**, etc.  
- Shared responsibilities for **testing**, **documentation**, and **integration** but may shift depending on team needs.  

Coordination will be done through:  
- **Agile sprints** with weekly standups.  
- **GitHub Projects** for tracking issues and tasks.  
- **Pull request reviews** to ensure code quality.  
- **Slack/Discord** for day-to-day communication.  

---

## ✅ Sprint 0 Checklist
- [x] `sprint0.md` in repo root with working links  
- [x] Repo URL included  
- [ ] Proposal slides linked  
- [x] Project summary & vision  
- [x] ≥4 core features + 1 non-functional feature  
- [x] ≥2 additional features  
- [x] User stories & acceptance criteria for core features  
- [ ] Initial architecture diagram added  
- [x] Work division & coordination  

---
