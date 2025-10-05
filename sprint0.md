# Sprint 0 – Gym Workout Tracking Application

## 1. Repository Links
- [Frontend GitHub Repository](https://github.com/Onyelechie/WorkoutPal-Frontend)
- [Backend GitHub Repository](https://github.com/Onyelechie/WorkoutPal-Backend)

## 2. Slides & Presentation
- [Proposal Slides (PDF)](./support_files/proposal_slides.pdf) 

## 3. Project Summary & Vision
Our project is a **Gym Workout Tracking Application** designed to help users plan, log, and visualize their fitness progress in a simple and motivating way. The application allows users to create personal profiles, design customized workout routines, and track sessions with built-in timers for rest and recovery.

### User Scenario
**Meet Sarah, a 22-year-old university student:** Sarah goes to the gym 3-4 times per week but struggles with consistency and progression tracking. She often forgets her previous weights, loses motivation when working out alone, and finds existing apps either too complex or lacking social features. With WorkoutPal, Sarah creates her profile, schedules her weekly workouts, and starts a session with built-in rest timers. She logs her sets with weights and RPE, sees her progress graphs over time, and stays motivated by viewing her friends' workout updates in the activity feed. When she hits a new personal record, her gym buddies celebrate with her through the app, keeping her accountable and engaged.

### Competitive Advantage
While apps like Strong and Jefit focus primarily on logging, and MyFitnessPal emphasizes nutrition, **WorkoutPal uniquely combines**:
- **Social accountability** with real-time friend activity feeds
- **Gamification elements** like streaks and achievement badges
- **Integrated session management** with smart rest timers and RPE tracking
- **Simplified UX** designed specifically for gym environments (large buttons, clear timers)
- **Goal visualization** with progress tracking that motivates continued use

The vision is to create an engaging platform that combines **fitness tracking** with **gamification** and **social interaction**, enabling users to build healthy habits and stay consistent in their fitness journey.  

---

## 4. Features & User Stories

### Core Features

#### 1. User Profiles
**User Story:**  
As someone who likes to stay organized, I want to store my profile information (name, handle, avatar, bio, height, weight, units, timezone) so that my workout logs are consistent.  

**Acceptance Criteria:**  
- **Given** a new user, **when** they register, **then** they must complete all required fields (name, handle, height, weight) within 5 minutes and receive confirmation within 2 seconds
- **Given** an existing user, **when** they update their profile, **then** changes must be saved within 3 seconds and reflected across all app views immediately
- **Given** a user's profile, **when** it is viewed, **then** it must display current stats, total workouts completed, current streak count, and top 3 achievements

---

#### 2. Workout Schedule and Routine
**User Story:**  
As a busy student, I want to schedule my workout sessions and stick to a routine so that I can stay consistent throughout the week.  

**Acceptance Criteria:**  
- **Given** a user, **when** they create a new workout plan with at least 3 exercises, **then** it must be saved within 2 seconds and appear in their workout library
- **Given** a user with a scheduled plan, **when** they view their calendar, **then** workouts must appear on correct dates with exercise count and estimated duration (±5 minutes)
- **Given** a user, **when** they edit or delete a workout plan, **then** the calendar must update within 1 second and send confirmation notification

---

#### 3. Session Runner + Timer
**User Story:**  
As someone who likes strictly timed workout sessions, I want to keep track of how long my sets and overall sessions last.  

**Acceptance Criteria:**  
- **Given** a user starts a session, **when** they hit "start," **then** the global session timer must begin within 0.5 seconds with accuracy to the second
- **Given** a user is in a session, **when** they log a set, **then** the app must store weight (0.1kg precision), reps (1-999 range), RPE (1-10 scale), and notes (max 200 characters) within 1 second
- **Given** a user finishes a set, **when** the rest timer starts, **then** countdown must be visible with audio/vibration notification at 10s, 5s, and 0s remaining
- **Given** a session in progress, **when** the user pauses or resumes, **then** all timers must adjust within 0.2 seconds maintaining accurate time tracking

---

#### 4. User Relationships
**User Story:**  
As an extroverted gym user, I want to connect with new people and share my workouts with them.  

**Acceptance Criteria:**  
- **Given** a user searches by username/handle, **when** they send a friend request, **then** the request must appear in recipient's notifications within 5 seconds with sender's profile preview
- **Given** a user has friends, **when** they complete a workout, **then** friends must see the update in their feed within 10 seconds showing exercise summary and duration

---

#### 5. Activity Feed
**User Story:**  
As someone who loses motivation easily, I want to see what other people are doing so that I can get inspiration and push myself to go to the gym.  

**Acceptance Criteria:**  
- **Given** a user's feed, **when** they open it, **then** they must see the 20 most recent activities from friends within 3 seconds, sorted by timestamp
- **Given** a friend completes a workout, **when** the feed updates, **then** the workout must appear with exercise count, duration, and PR indicators within 10 seconds

---

#### 6. Goal Tracking & Targets
**User Story:**  
As someone who is goal oriented, I want to keep track of set goals for weight, body-fat, and lifts for target dates.  

**Acceptance Criteria:**  
- **Given** a user, **when** they create a new goal with target value and date (max 1 year future), **then** it must be saved within 2 seconds with progress tracking enabled
- **Given** a user with active goals, **when** they log relevant progress data, **then** goal progress must update within 3 seconds showing percentage completion
- **Given** a user's dashboard, **when** it is viewed, **then** goals must display current vs. target values, percentage complete, and days remaining with visual progress bars

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
![Architecture Diagram](./support_files/architecture-diagram.png)

### Technology Stack & Rationale

**Frontend: React.js**
- **Chosen for:** Component reusability, large ecosystem, team familiarity
- **Alternatives considered:** Vue.js (simpler learning curve), Angular (enterprise features)
- **Trade-offs:** React's complexity vs Vue's simplicity; chose React for better job market alignment and extensive third-party libraries

**Backend: Go**
- **Chosen for:** High concurrency, fast compilation, strong typing, minimal memory footprint
- **Alternatives considered:** Node.js (JavaScript consistency), Python (rapid development), Java (enterprise maturity)
- **Trade-offs:** Go's smaller ecosystem vs Node.js familiarity; chose Go for superior performance handling concurrent workout sessions and real-time features

**Database: MySQL**
- **Chosen for:** ACID compliance, mature ecosystem, structured data relationships
- **Alternatives considered:** PostgreSQL (advanced features), MongoDB (document flexibility), Redis (caching)
- **Trade-offs:** MySQL's relational constraints vs MongoDB's flexibility; chose MySQL for data integrity in user relationships and workout tracking

### Architecture Constraints & Risks
- **Performance Risk:** Real-time friend activity updates may strain database with high user counts
- **Scalability Constraint:** Single MySQL instance limits horizontal scaling
- **Development Risk:** Go's smaller talent pool may slow hiring
- **Mitigation:** Implement Redis caching layer, database read replicas, and comprehensive Go training

This architecture provides **clear separation of concerns** while optimizing for our core requirements: real-time workout tracking, social features, and data consistency.

### Repository Structure Decision
We chose **separate repositories** for frontend and backend instead of a monorepo for the following reasons:

- **Independent Development Cycles:** Frontend and backend can be developed, tested, and deployed independently without affecting each other
- **Technology-Specific Tooling:** Each repository can have its own build tools, linting rules, and CI/CD pipelines optimized for React.js and Go respectively
- **Team Specialization:** Developers can focus on their expertise area without needing to understand the entire codebase
- **Deployment Flexibility:** Separate repositories allow for independent scaling and deployment strategies
- **Reduced Complexity:** Smaller, focused repositories are easier to navigate and maintain for a team of our size

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

For this sprint, we all contributed by having a discord call then an in class meeting to discuss what ideas to pursure and what technologies to use. Overall, for this sprint, everyone contributed equally to ideas. Ebere set up the github repositories and presentation while Taren made the high level architecture. 

---

## ✅ Sprint 0 Checklist
- [x] `sprint0.md` in repo root with working links  
- [x] Repo URL included  
- [x] Proposal slides linked  
- [x] Project summary & vision  
- [x] ≥4 core features + 1 non-functional feature  
- [x] ≥2 additional features  
- [x] User stories & acceptance criteria for core features  
- [x] Initial architecture diagram added  
- [x] Work division & coordination  

---