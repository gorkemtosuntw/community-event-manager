const request = require('supertest');
const app = require('./app');

describe('User Service', () => {
    test('Health check endpoint should return healthy status', async () => {
        const response = await request(app).get('/health');
        expect(response.statusCode).toBe(200);
        expect(response.body).toEqual({ status: 'healthy' });
    });

    test('Should create a new user', async () => {
        const userData = {
            username: 'testuser',
            email: 'test@example.com',
            password: 'password123'
        };

        const response = await request(app)
            .post('/users')
            .send(userData);

        expect(response.statusCode).toBe(201);
        expect(response.body).toHaveProperty('id');
        expect(response.body.username).toBe(userData.username);
        expect(response.body.email).toBe(userData.email);
        expect(response.body).not.toHaveProperty('password');
    });
});
