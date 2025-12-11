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
  return request('/conversations');
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