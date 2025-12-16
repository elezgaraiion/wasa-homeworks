// src/services/api.js
const API_URL = 'http://localhost:3000'; 

async function request(endpoint, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...options.headers };
  
  const token = localStorage.getItem('userId');
  
  // CORRECCIÓN DEFINITIVA: 
  // Enviamos solo el token limpio. 
  // Si tu backend recibía "Bearer ...", ahora recibirá "af1d..."
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
  // Añadimos ?_t=... para forzar al navegador a pedir datos frescos siempre
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
  
  // Aquí también enviamos solo el token limpio
  if (token) headers['Authorization'] = token;

  const response = await fetch(`${API_URL}/me/photo`, {
    method: 'PUT',
    headers: headers,
    body: formData
  });

  if (!response.ok) throw new Error('Error al subir la imagen');
  return response.json();
}
export async function searchUsers(query) {
  // GET /users?q=nombre
  // encodeURIComponent asegura que espacios y acentos viajen bien en la URL
  const queryString = query ? `?q=${encodeURIComponent(query)}` : '';
  return request(`/users${queryString}`);
}
// En src/services/api.js

export async function getConversationMessages(conversationId) {
  // El truco es añadir ?_t=AHORA_MISMO
  // Así el navegador cree que es una petición nueva y no usa la memoria caché (304)
  return request(`/conversations/${conversationId}/messages?limit=50&_t=${Date.now()}`);
}
export async function createPrivateChat(targetUserId) {
  // Cambiamos la URL para coincidir con el backend
  return request('/chats', {
    method: 'POST',
    body: JSON.stringify({ targetUserId })
  });
}
export async function sendMessage(chatId, text, replyToMessageId = null) {
  const payload = { 
      text: text,
      // IMPORTANTE: Asegúrate de que el nombre del campo coincide
      // con lo que espera tu struct de Go (json:"replyToMessageId")
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
  // POST /conversations/:id/seen
  return request(`/conversations/${conversationId}/seen`, {
    method: 'POST'
  });
}
export async function createGroup(name, userIds) {
  // POST /groups
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
  // DELETE /conversations/:id/messages/:msgId
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

// Salir del grupo
export async function leaveGroup(conversationId) {
  return request(`/conversations/${conversationId}/users/me`, {
      method: 'DELETE'
  });
}

// Cambiar nombre del grupo
export async function setGroupName(conversationId, name) {
  return request(`/conversations/${conversationId}/name`, {
      method: 'PUT',
      body: JSON.stringify({ name })
  });
}

// Cambiar foto del grupo (Multipart/Form-data)
export async function setGroupPhoto(conversationId, file) {
  const formData = new FormData();
  formData.append('photoFile', file);

  const token = localStorage.getItem('token');
  const response = await fetch(`${API_URL}/conversations/${conversationId}/photo`, {
      method: 'PUT',
      headers: {
          'Authorization': token // NO poner Content-Type, el navegador lo pone solo con el boundary
      },
      body: formData
  });

  if (!response.ok) {
      throw new Error(await response.text());
  }
  return await response.json();
}