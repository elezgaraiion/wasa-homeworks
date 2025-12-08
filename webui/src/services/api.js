const API_URL = __API_URL__  // definir en Vite env

export async function loginOrRegister(name) {
  const res = await fetch(`${API_URL}/session`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({ name })
  })

  if (!res.ok) throw new Error('Error al conectar con el backend')
  return res.json() // { identifier: "id-del-usuario" }
}

export async function getCurrentUser(userId) {
  const res = await fetch(`${API_URL}/me`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': userId
    }
  })

  if (!res.ok) throw new Error('No se pudo obtener el perfil del usuario')
  return res.json()
}

export async function updateUserName(userId, newName) {
  const res = await fetch(`${API_URL}/me/username`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': userId // solo UUID, sin "Bearer"
    },
    body: JSON.stringify({ name: newName })
  })

  if (!res.ok) {
    const text = await res.text()
    throw new Error(text)
  }

  return res.json()
}

export async function getConversations(userId) {
  const res = await fetch(`${API_URL}/conversations`, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
      'Authorization': userId
    }
  })

  if (!res.ok) throw new Error('Error al obtener conversaciones')
  return res.json()
}
