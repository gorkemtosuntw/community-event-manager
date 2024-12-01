const express = require('express');
const bodyParser = require('body-parser');
const { v4: uuidv4 } = require('uuid');
const cors = require('cors');

const app = express();

// Middleware
app.use(cors());
app.use(bodyParser.json());

// Error handling middleware
app.use((err, req, res, next) => {
    console.error(err.stack);
    res.status(500).json({ error: 'Something went wrong!' });
});

// In-memory user storage (would be replaced by a database in production)
const users = {};

// Health check endpoint
app.get('/health', (req, res) => {
    res.json({ status: 'healthy' });
});

// User registration
app.post('/users', (req, res) => {
    try {
        const { username, email, password } = req.body;

        // Basic validation
        if (!username || !email || !password) {
            return res.status(400).json({ error: 'Missing required fields' });
        }

        // Check if user already exists
        const existingUser = Object.values(users).find(
            user => user.username === username || user.email === email
        );

        if (existingUser) {
            return res.status(409).json({ error: 'User already exists' });
        }

        // Create new user
        const userId = uuidv4();
        const newUser = {
            id: userId,
            username,
            email,
            // In a real app, NEVER store plain text passwords
            password,
            createdAt: new Date().toISOString()
        };

        users[userId] = newUser;

        // Remove password before sending response
        const { password: _, ...userResponse } = newUser;
        res.status(201).json(userResponse);
    } catch (error) {
        console.error('Error in user registration:', error);
        res.status(500).json({ error: 'Internal server error' });
    }
});

// User authentication (simplified)
app.post('/users/login', (req, res) => {
    try {
        const { username, password } = req.body;

        const user = Object.values(users).find(
            u => u.username === username && u.password === password
        );

        if (!user) {
            return res.status(401).json({ error: 'Invalid credentials' });
        }

        res.json({
            message: 'Login successful',
            userId: user.id
        });
    } catch (error) {
        console.error('Error in user login:', error);
        res.status(500).json({ error: 'Internal server error' });
    }
});

// Get user profile
app.get('/users/:userId', (req, res) => {
    try {
        const user = users[req.params.userId];

        if (!user) {
            return res.status(404).json({ error: 'User not found' });
        }

        const { password, ...userProfile } = user;
        res.json(userProfile);
    } catch (error) {
        console.error('Error in getting user profile:', error);
        res.status(500).json({ error: 'Internal server error' });
    }
});

// Graceful shutdown handling
process.on('SIGTERM', () => {
    console.log('Received SIGTERM. Performing graceful shutdown...');
    process.exit(0);
});

process.on('SIGINT', () => {
    console.log('Received SIGINT. Performing graceful shutdown...');
    process.exit(0);
});

const PORT = process.env.PORT || 3000;
app.listen(PORT, '0.0.0.0', () => {
    console.log(`User Service running on port ${PORT}`);
});
