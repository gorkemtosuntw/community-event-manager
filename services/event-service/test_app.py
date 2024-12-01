import pytest
from app import app
import json

@pytest.fixture
def client():
    app.config['TESTING'] = True
    with app.test_client() as client:
        yield client

def test_create_event(client):
    event_data = {
        "title": "Test Event",
        "description": "Test Description",
        "date": "2024-12-31",
        "location": "Test Location"
    }
    response = client.post('/events', 
                         data=json.dumps(event_data),
                         content_type='application/json')
    assert response.status_code == 201
    data = json.loads(response.data)
    assert data['title'] == event_data['title']
    assert 'id' in data

def test_health_check(client):
    response = client.get('/health')
    assert response.status_code == 200
    assert json.loads(response.data)['status'] == 'healthy'
