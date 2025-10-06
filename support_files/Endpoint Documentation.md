### Authorization
##### Login 
Login is required for any subsequent API calls. Successful login returns an HTTP-Only cookie which needs to be appended for future data retrieval or manipulation.
`POST /login` 
```json
{
	"email" : 'email@email.email',
	"password" : 'supersecure'
}
```

##### Logout
Logout expires the cookie.
`POST /logout`

---
### User

##### Read One
`GET /users/{id}`
```json
{
	"id" : 1,
	"username" : 'bigworkoutman',
	"email" : 'email@email.email',
	"password" : 'supersecure',
	"name" : 'Sagev Wen Outfall',
	"age" : 57,
	"heightValue" : 3.7,
	"heightMetric" : 'meters',
	"weightValue" : 281,
	"weightMetric" : 'kg'
}
```

##### Read Many
`GET /users`
```json
[
	{
		"id" : 1,
		"username" : 'bigworkoutman',
		"email" : 'email@email.email',
		"password" : 'supersecure',
		"name" : 'Sagev Wen Outfall',
		"age" : 57,
		"heightValue" : 3.7,
		"heightMetric" : 'meters',
		"weightValue" : 281,
		"weightMetric" : 'kg'
	}
]
```

##### Creation
`POST /users`
```json
{
	"username" : 'bigworkoutman',
	"email" : 'email@email.email',
	"password" : 'supersecure',
	"name" : 'Sagev Wen Outfall',
	"age" : 57,
	"heightValue" : 3.7,
	"heightMetric" : 'meters',
	"weightValue" : 281,
	"weightMetric" : 'kg'
}
```

##### Updating
`PATCH /users/{id}`
```json
{
	"username" : 'bigworkoutman',
	"email" : 'email@email.email',
	"password" : 'supersecure',
	"name" : 'Sagev Wen Outfall',
	"age" : 57,
	"heightValue" : 3.7,
	"heightMetric" : 'meters',
	"weightValue" : 281,
	"weightMetric" : 'kg'
}
```

##### Delete
`DELETE /users/{id}`

---
### Relationships

##### Follow 
By following another user their public posts and achievements will be included in your 'posts feed'.
`POST /relationships/follow`
```json
{
	"userID" : 2
}
```

##### Unfollow
By unfollowing another user their public posts and achievements will no-longer be included in your 'posts feed'.
`POST /relationships/unfollow`
```json
{
	"userID" : 2
}
```

---
### Posts
##### Read
Reads posts with a flag to filter 'followed' posts or random global posts.
`GET /posts?followed=true`
```json
[
	{
		"id" : 1,
		"postedBy" : 'bigworkoutman',
		"title" : 'Big Workout',
		"caption" : 'Bigger gains',
		"date" : '01/01/0001',
		"content" : 'https://some-cdn-link.com',
		"likes" : 100320234,
		"comments": [
			{
				"commentedBy" : 'bigworkoutman',
				"comment" : 'sheeeeesh',
				"date" : '01/02/0003'
			}
		]
	}
]
```

##### Create Post
Private flag determines whether this is visible to all users or only your friends (followers that you are following).
`POST /posts`
```json
{
	"postedBy" : 'bigworkoutman',
	"title" : 'Big Workout',
	"caption" : 'Bigger gains',
	"date" : '01/01/0001',
	"content" : 'https://some-cdn-link.com',
	"likes" : 100320234,
	"private" : true
}
```

##### Comment on Post
`POST /posts/{id}`
```json
{
	"comment" : 'sheeeeesh'
}
```

##### Like a Post
`POST /posts/{id}`

---
### Exercises
##### Read
Can be filtered by target, intensity and expertise and custom
`GET /exercises?type=chest&intensity=high&expertise=beginner?custom=false`
```json
[
	{
		"id" : 1,
		"name" : 'pushups',
		"targets" : ['chest','shoulder','arms'],
		"intesity" : 'high',
		"expertise" : 'beginner',
		"image" : "https://some-cdn-link.com",
		"demo" : "https://some-cdn-link.com",
		"recommendedCount" : 5000,
		"recommendedSets" : 15,
		"recommendedDuration" : 1000,
		"custom" : false
	}
]
```

##### Create Custom exercise
`POST /exercises`
```json
{
	"name" : 'pushups',
	"targets" : ['chest','shoulder','arms'],
	"intesity" : 'high',
	"expertise" : 'beginner',
	"image" : "https://some-cdn-link.com",
	"demo" : "https://some-cdn-link.com",
	"recommendedCount" : 5000,
	"recommendedSets" : 15,
	"recommendedDuration" : 1000,
}
```

---
### Workouts

##### Read
`GET /workouts`
```json
[
	{
		"id" : 1,
		"name" : 'Big Workout',
		"frequency" : 'Daily',
		"nextRound" : '01/02/0004',
		"exercises" : [
			{
				"startTime" : '00:00',
				"endTime" : '01:30',
				"id" : 1,
				"name" : 'pushups',
				"targets" : ['chest','shoulder','arms'],
				"intesity" : 'high',
				"expertise" : 'beginner',
				"image" : "https://some-cdn-link.com",
				"demo" : "https://some-cdn-link.com",
				"count" : 5000,
				"sets" : 15,
				"duration" : 1000,
			}
		]
	}
]
```

##### Create
`POST /workouts`
```json
{
	"name" : 'Big Workout',
	"frequency" : 'Daily',
	"nextRound" : '01/02/0004',
	"exercises" : [
		{
			"startTime" : '00:00',
			"endTime" : '01:30',
			"id" : 1,
			"count" : 5000,
			"sets" : 15,
			"duration" : 1000,
		}		
	]
}
```

##### Update
`POST /workouts/{id}`
```json
{
	"name" : 'Big Workout',
	"frequency" : 'Daily',
	"nextRound" : '01/02/0004',
	"exercises" : [
		{
			"startTime" : '00:00',
			"endTime" : '01:30',
			"id" : 1,
			"count" : 5000,
			"sets" : 15,
			"duration" : 1000,
		}		
	]
}
```

