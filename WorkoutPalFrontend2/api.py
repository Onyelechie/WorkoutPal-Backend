import requests

BASE_URL = "http://localhost:8080"
session = requests.Session()

def api_get(path, params=None):
    url = f"{BASE_URL}{path}"
    res = session.get(url, params=params)
    return res

def api_post(path, json=None):
    url = f"{BASE_URL}{path}"
    res = session.post(url, json=json)
    return res

def api_patch(path, json=None):
    url = f"{BASE_URL}{path}"
    res = session.patch(url, json=json)
    return res

def api_put(path, json=None):
    url = f"{BASE_URL}{path}"
    res = session.put(url, json=json)
    return res