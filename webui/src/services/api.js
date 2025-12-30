const API_URL = 'http://localhost:3000'; 

async function request(endpoint, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...options.headers };
  
  const token = localStorage.getItem('userId');
  
  if (token) {
    headers['Authorization'] = token; 
  }

  const config = { ...options, headers };
  
  try {
    const response = await fetch(`${API_URL}${endpoint}`, config);
    
    if (!response.ok) {
      if (response.status === 401) throw new Error('Unauthorized');
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || response.statusText);
    }

    if (response.status === 204) return null;
    return response.json();
  } catch (err) {
    throw err;
  }
}

export async function doLogin(name) {
  return request('/session', {
    method: 'POST',
    body: JSON.stringify({ name })
  });
}

export async function getCurrentUser() {
  return request('/me');
}

export async function getConversations() {
  return request(`/conversations?_t=${Date.now()}`);
}

export async function updateUserName(newName) {
  return request('/me/username', {
    method: 'PUT',
    body: JSON.stringify({ name: newName })
  });
}

export async function updatePhoto(file) {
  const formData = new FormData();
  formData.append('photoFile', file); 

  const token = localStorage.getItem('userId');
  const headers = {};
  
  if (token) headers['Authorization'] = token;

  const response = await fetch(`${API_URL}/me/photo`, {
    method: 'PUT',
    headers: headers,
    body: formData
  });

  if (!response.ok) throw new Error('Error uploading image');
  return response.json();
}

export async function searchUsers(query) {
  const queryString = query ? `?q=${encodeURIComponent(query)}` : '';
  return request(`/users${queryString}`);
}

export async function getConversationMessages(conversationId) {
  return request(`/conversations/${conversationId}/messages?limit=50&_t=${Date.now()}`);
}

export async function createPrivateChat(targetUserId) {
  return request('/chats', {
    method: 'POST',
    body: JSON.stringify({ targetUserId })
  });
}

export async function sendMessage(chatId, text, replyToMessageId = null) {
  const payload = { 
      text: text,
      replyToMessageId: replyToMessageId 
  };
  
  return request(`/conversations/${chatId}/messages`, {
    method: 'POST',
    body: JSON.stringify(payload)
  });
}

export async function getChatInfo(conversationId) {
  return request(`/conversations/${conversationId}`);
}

export async function markChatAsRead(conversationId) {
  return request(`/conversations/${conversationId}/seen`, {
    method: 'POST'
  });
}

export async function createGroup(name, userIds) {
  return request('/conversations', {
    method: 'POST',
    body: JSON.stringify({ name: name, users: userIds })
  });
}

export async function forwardMessage(sourceChatId, messageId, targetChatId) {
  return request(`/conversations/${sourceChatId}/messages/${messageId}`, {
    method: 'POST',
    body: JSON.stringify({ targetConversationId: targetChatId })
  });
}

export async function deleteMessage(chatId, messageId) {
  return request(`/conversations/${chatId}/messages/${messageId}`, {
    method: 'DELETE'
  });
}

export async function addReaction(chatId, messageId, emoji) {
  return request(`/conversations/${chatId}/messages/${messageId}/reactions`, {
    method: 'POST',
    body: JSON.stringify({ emoji })
  });
}

export async function removeReaction(chatId, messageId, reactionId) {
  return request(`/conversations/${chatId}/messages/${messageId}/reactions/${reactionId}`, {
    method: 'DELETE'
  });
}

export async function addUserToGroup(conversationId, userId) {
  return request(`/conversations/${conversationId}/users`, {
      method: 'POST',
      body: JSON.stringify({ userId })
  });
}

export async function leaveGroup(conversationId) {
  return request(`/conversations/${conversationId}/users/me`, {
      method: 'DELETE'
  });
}

export async function setGroupName(conversationId, name) {
  return request(`/conversations/${conversationId}/name`, {
      method: 'PUT',
      body: JSON.stringify({ name })
  });
}

export async function setGroupPhoto(conversationId, file) {
  const formData = new FormData();
  formData.append('photoFile', file);

  const token = localStorage.getItem('token');
  const response = await fetch(`${API_URL}/conversations/${conversationId}/photo`, {
      method: 'PUT',
      headers: {
          'Authorization': token 
      },
      body: formData
  });

  if (!response.ok) {
      throw new Error(await response.text());
  }
  return await response.json();
}