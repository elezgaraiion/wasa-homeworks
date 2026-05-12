# WASAText Frontend

Frontend implementation of WASAText, a web-based instant messaging application inspired by WhatsApp Web.
This project was developed as part of a university assignment using the provided OpenAPI specification.

The application allows users to communicate through private and group conversations, exchange messages, react with emojis, and manage their personal profile.

## Features

### Authentication
- Simple login/registration system using username-based authentication.
- Session persistence using Bearer Token authentication.

### User Profile
- View current user information.
- Update username.
- Upload and change profile picture.

### Conversations
- View all conversations sorted by recent activity.
- Start private chats with other users.
- Create group conversations.
- Display unread message counters and message previews.

### Messaging
- Send text messages.
- Send photo messages.
- Reply to messages.
- Forward messages between conversations.
- Delete own messages.
- Pagination support for message history.

### Reactions
- Add emoji reactions to messages.
- Update existing reactions.
- Remove reactions.

### Groups
- Create group chats.
- Add users to groups.
- Leave groups.
- Change group name.
- Update group picture.

### Read Status
- Mark conversations as seen.
- Display message delivery/read status.

## Tech Stack
- Vue.js 3
- Vue Router
- Axios
- Bootstrap / Custom CSS
- OpenAPI-based REST API
- Vite

## Project Structure
```text
src/
├── assets/          # Static assets
├── components/      # Reusable UI components
├── pages/           # Application pages/views
├── router/          # Vue Router configuration
├── services/        # API communication layer
├── stores/          # State management
├── utils/           # Utility functions
└── App.vue
```

## API Integration

The frontend communicates with the WASAText backend through a REST API defined using the OpenAPI 3.0 specification.

Main API areas:
- Session management
- User management
- Conversations
- Messages
- Reactions
- Groups

Authentication is handled using Bearer Tokens.

Example request:
```http
GET /conversations
Authorization: Bearer <user-id>
```

## Main Screens

### Login Page
- Username-based authentication.
- Automatic registration if the username does not exist.

### Chat List
- Displays all conversations.
- Shows latest message preview and unread counts.

### Chat View
- Real-time style messaging interface.
- Supports text, photos, replies, reactions, and forwarding.

### Profile Settings
- Change username and profile picture.

### Group Management
- Create and manage group conversations.

## Installation

Clone the repository:
```bash
git clone <repository-url>
cd wasatext-frontend
```

Install dependencies:
```bash
npm install
```

Run the development server:
```bash
npm run dev
```

Build for production:
```bash
npm run build
```

## Design Goals

The main objective of this project was to recreate the core user experience of modern messaging platforms such as WhatsApp Web while following:
- RESTful API principles
- Component-based frontend architecture
- Responsive UI design
- Clean state management
- Modular and maintainable code structure

## Possible Improvements
- WebSocket real-time communication
- Push notifications
- Typing indicators
- Voice messages
- End-to-end encryption
- Dark mode support

## Academic Context

This project was developed for educational purposes as part of a university course focused on:
- Web development
- REST APIs
- Frontend engineering
- Client-server communication
- OpenAPI integration