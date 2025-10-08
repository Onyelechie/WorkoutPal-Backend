#!/bin/bash

echo "=== Testing User API ==="

echo "1. Creating user..."
curl -X POST http://localhost:8080/users \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "name": "Test User",
    "password": "password123",
    "age": 25,
    "height": 175.5,
    "heightMetric": "cm",
    "weight": 70.0,
    "weightMetric": "kg"
  }'

echo -e "\n\n2. Getting all users..."
curl http://localhost:8080/users

echo -e "\n\n3. Getting user by ID..."
curl http://localhost:8080/users/1

echo -e "\n\n4. Creating goal..."
curl -X POST http://localhost:8080/users/1/goals \
  -H "Content-Type: application/json" \
  -d '{
    "type": "weight",
    "targetValue": 65.0,
    "targetDate": "2024-12-31"
  }'

echo -e "\n\n5. Getting user goals..."
curl http://localhost:8080/users/1/goals

echo -e "\n\nDone!"