// src/services/api.js
const API_URL = 'http://localhost:3000'; // Asegúrate de que este es el puerto de tu Go

async function request(endpoint, options = {}) {
  const headers = { 'Content-Type': 'application/json', ...options.headers };
  
  // Si ya tenemos token, lo inyectamos (para el futuro)
  const token = localStorage.getItem('userId');
  if (token) headers['Authorization'] = `Bearer ${token}`;

  const config = { ...options, headers };
  
  try {
    const response = await fetch(`${API_URL}${endpoint}`, config);
    if (!response.ok) {
      // Intentamos leer el error que manda tu Go
      const errorData = await response.json().catch(() => ({}));
      throw new Error(errorData.error || 'Error en el servidor');
    }
    return response.json();
  } catch (err) {
    console.error("API Error:", err);
    throw err;
  }
}

// Esta es la función que usaremos ahora
export async function doLogin(name) {
  // Tu backend espera POST /session con body { "name": "..." }
  return request('/session', {
    method: 'POST',
    body: JSON.stringify({ name })
  });
}

// Para obtener el perfil una vez logueado
export async function getCurrentUser() {
  return request('/me');
}