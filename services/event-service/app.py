from flask import Flask, request, jsonify
from flask_cors import CORS
from datetime import datetime
import uuid

app = Flask(__name__)
CORS(app)  # Enable CORS for all routes

# In-memory event storage (would be replaced by a database in production)
events = {}

@app.route('/health', methods=['GET'])
def health_check():
    return jsonify({"status": "healthy"}), 200

@app.route('/events', methods=['POST'])
def create_event():
    event_data = request.json
    event_id = str(uuid.uuid4())
    
    # Validate event data
    required_fields = ['title', 'date']
    if not event_data or not all(field in event_data for field in required_fields):
        return jsonify({"error": "Missing required fields"}), 400
    
    # Create event
    event = {
        'id': event_id,
        'title': event_data['title'],
        'description': event_data.get('description', ''),
        'date': event_data['date'],
        'location': event_data.get('location', ''),
        'created_at': datetime.utcnow().isoformat()
    }
    
    events[event_id] = event
    return jsonify(event), 201

@app.route('/events', methods=['GET'])
def list_events():
    return jsonify(list(events.values())), 200

@app.route('/events/<event_id>', methods=['GET'])
def get_event(event_id):
    event = events.get(event_id)
    if not event:
        return jsonify({"error": "Event not found"}), 404
    return jsonify(event), 200

if __name__ == '__main__':
    # Add error handling
    try:
        app.run(host='0.0.0.0', port=5000)
    except Exception as e:
        print(f"Error starting the server: {e}")
        exit(1)
