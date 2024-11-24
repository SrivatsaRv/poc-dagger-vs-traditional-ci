import pytest
from app import app

@pytest.fixture
def client():
    app.testing = True
    with app.test_client() as client:
        yield client

def test_health_check(client):
    response = client.get('/healthcheck')
    assert response.status_code == 200
    assert response.get_json() == {"message": "REST API Healthy"}  # Matches the actual response

def test_greet(client):
    response = client.get('/')
    assert response.status_code == 200
    assert response.get_json() == {"message": "Hello, World!"}  # Matches the actual response
