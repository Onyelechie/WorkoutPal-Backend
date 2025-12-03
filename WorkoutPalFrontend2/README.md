# THIS PROJECT WAS MADE USING CHATGPT

## 1. Requirements

To run the WorkoutPal Python client, ensure you have the following installed:

---

#### Python 3.10+

Check your version:

```bash
python3 --version
```

---

### pip (Python package manager)

Usually included with Python:

```bash
pip --version
```

---

### Running the WorkoutPal Go backend API

The Python client communicates with the backend, so make sure the API is running locally.
Before starting the backend, create a `.env` file in the root directory of the Go project and add:

```dotenv
COOKIE_SECURE=false
```

---

## 2. Installation

### Step 1 — Create a Virtual Environment

A virtual environment keeps your Python dependencies isolated.

Mac / Linux:

```bash
python3 -m venv venv
source venv/bin/activate
```

Windows (PowerShell):

```bash
python -m venv venv
.\venv\Scripts\activate
```

---

### 3. Install Dependencies

Once inside the virtual environment, install all required packages:

```bash
pip install -r requirements.txt
```

---

### 4. Run the Client

Start the WorkoutPal Python CLI:

```bash
python main.py
```

---

### 5. Follow Instructions in the Client

Once running enter the `create user` command and follow the prompts

Then you can `login` with you newly created user

Follow the other tool tips or enter `help` at any point to see what commands are available