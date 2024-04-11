from flask import Flask, request, render_template
from typing import List
import json

app = Flask(__name__)

sessions = []

@app.route('/sessions', methods=['POST'])
def receive_sessions():
    data = request.get_json()
    if isinstance(data, List):
        sessions.extend(data)
    else:
        sessions.append(data)
    return 'OK', 200

@app.route('/')
def show_sessions():
    return render_template('sessions.html', sessions=sessions)

if __name__ == '__main__':
    app.run(host='0.0.0.0', port=8080)