#!/bin/bash

# WorkoutPal API Testing Script
# Make sure the server is running: go run src/cmd/api/main.go

BASE_URL="http://localhost:8080"

echo "🏋️ WorkoutPal API Testing Script"
echo "================================="
echo ""

# Test 1: Create a user
echo "1️⃣ Creating a test user..."
USER_RESPONSE=$(curl -s -X POST $BASE_URL/users \
  -H "Content-Type: application/json" \
  -d '{"username":"testuser","email":"test@example.com","name":"Test User","password":"password123"}')
echo "Response: $USER_RESPONSE"
echo ""

# Test 2: Get all users
echo "2️⃣ Getting all users..."
curl -s $BASE_URL/users | jq '.'
echo ""

# Test 3: Get all exercises
echo "3️⃣ Getting all exercises..."
curl -s $BASE_URL/exercises | jq '.'
echo ""

# Test 4: Create a goal for user 1
echo "4️⃣ Creating a goal for user 1..."
GOAL_RESPONSE=$(curl -s -X POST $BASE_URL/users/1/goals \
  -H "Content-Type: application/json" \
  -d '{"name":"Weight Loss","description":"Lose 10kg","deadline":"2024-12-31"}')
echo "Response: $GOAL_RESPONSE"
echo ""

# Test 5: Create a routine for user 1
echo "5️⃣ Creating a routine for user 1..."
ROUTINE_RESPONSE=$(curl -s -X POST $BASE_URL/users/1/routines \
  -H "Content-Type: application/json" \
  -d '{"name":"Morning Workout","description":"Daily routine"}')
echo "Response: $ROUTINE_RESPONSE"
echo ""

# Test 6: Add exercise to routine
echo "6️⃣ Adding exercise 1 to routine 1..."
ADD_EXERCISE_RESPONSE=$(curl -s -X POST "$BASE_URL/routines/1/exercises?exercise_id=1")
echo "Response: $ADD_EXERCISE_RESPONSE"
echo ""

# Test 7: Get routine with exercises
echo "7️⃣ Getting routine 1 with exercises..."
curl -s $BASE_URL/routines/1 | jq '.'
echo ""

# Test 8: Get user's goals
echo "8️⃣ Getting user 1's goals..."
curl -s $BASE_URL/users/1/goals | jq '.'
echo ""

# Test 9: Get user's routines
echo "9️⃣ Getting user 1's routines..."
curl -s $BASE_URL/users/1/routines | jq '.'
echo ""

echo "✅ API testing complete!"
echo ""
echo "💡 Tips:"
echo "- Add -v flag to curl commands for verbose output"
echo "- Use jq for pretty JSON formatting"
echo "- Check Swagger docs at: $BASE_URL/swagger/index.html"