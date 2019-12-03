import requests
import json

params = {
    'email':'celvis@gmail.com',
    'password':'12345',
    'username':'elvito',
    'id':'3',
    'comment':"elis datum fresi"
}
params = json.dumps(params)
resp = requests.post('http://127.0.0.1:5000/post',data=params)
print(resp.content)