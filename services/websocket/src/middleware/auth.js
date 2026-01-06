// TODO: Auth middleware to validate JWT tokens
// Will call Go auth service to verify tokens

const jwt = require('jsonwebtoken')
const axios = require('axios')

const AUTH_SERVICE_URL = process.env.AUTH_SERVICE_URL || 'http://localhost:8081'

async function authenticateWebSocket(token) {
  // TODO: Implement JWT validation
  // TODO: Call Go auth service if needed
  return null
}

module.exports = {
  authenticateWebSocket,
}
