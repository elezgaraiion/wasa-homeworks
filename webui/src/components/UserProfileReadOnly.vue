<template>
  <div class="profile-overlay" @click.self="$emit('close')">
    <div class="profile-card">
      
      <div class="card-header">
        <h2>Info del Contacto</h2>
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>

      <div class="avatar-section">
        <div class="avatar-wrapper">
          <img 
            :src="user.photo || user.Photo || 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png'" 
            class="avatar-img" 
          />
        </div>
      </div>

      <div class="info-section">
        <h3 class="user-name">{{ user.name || user.Name }}</h3>
        
        <div class="id-group">
          <label>ID de Usuario</label>
          <div class="id-box">
            {{ user.id || user.ID }}
          </div>
        </div>

        <p class="user-bio">¡Hola! Estoy usando Wasa Web.</p>
      </div>

      <div class="actions">
        <button class="btn-chat" @click="$emit('startChat', user)">
          <span class="icon">💬</span> Enviar Mensaje
        </button>
      </div>

    </div>
  </div>
</template>

<script setup>
// Recibimos el objeto 'user' completo
defineProps(['user']);
defineEmits(['close', 'startChat']);
</script>

<style scoped>
/* REUTILIZAMOS EL ESTILO GLASSMORPHISM */
.profile-overlay {
  position: fixed; top: 0; left: 0; width: 100vw; height: 100vh;
  background: rgba(0, 0, 0, 0.7); backdrop-filter: blur(8px);
  z-index: 10001; /* ENCIMA DEL BUSCADOR */
  display: flex; justify-content: center; align-items: center;
  animation: fadeIn 0.2s ease;
}

.profile-card {
  width: 350px; max-width: 90%; background: #111b21; border: 1px solid #333;
  border-radius: 16px; padding: 2rem; box-shadow: 0 25px 50px rgba(0,0,0,0.5);
  display: flex; flex-direction: column; align-items: center; gap: 1.5rem;
  animation: scaleUp 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.card-header { width: 100%; display: flex; justify-content: space-between; align-items: center; margin-bottom: 0.5rem; }
.card-header h2 { margin: 0; color: #e9edef; font-weight: 300; font-size: 1.4rem; }
.btn-close { background: transparent; border: none; color: #8696a0; font-size: 1.5rem; cursor: pointer; }
.btn-close:hover { color: #f15c6d; }

/* AVATAR */
.avatar-wrapper {
  width: 150px; height: 150px; border-radius: 50%; overflow: hidden;
  border: 4px solid #202c33; box-shadow: 0 0 20px rgba(0,0,0,0.3);
}
.avatar-img { width: 100%; height: 100%; object-fit: cover; }

/* INFO */
.info-section { text-align: center; width: 100%; }
.user-name { font-size: 1.8rem; color: #e9edef; margin: 0 0 10px 0; font-weight: 500; }
.user-bio { color: #8696a0; font-style: italic; font-size: 0.9rem; margin-top: 15px; }

.id-group { text-align: left; margin-top: 15px; background: rgba(255,255,255,0.03); padding: 10px; border-radius: 8px; }
.id-group label { display: block; color: #00a884; font-size: 0.75rem; font-weight: bold; text-transform: uppercase; margin-bottom: 5px; }
.id-box { color: #d1d7db; font-family: monospace; font-size: 0.85rem; word-break: break-all; }

/* BOTÓN */
.btn-chat {
  width: 100%; background-color: #00a884; color: #111; border: none;
  padding: 12px; border-radius: 24px; font-weight: bold; font-size: 1rem;
  cursor: pointer; display: flex; justify-content: center; align-items: center; gap: 8px;
  transition: transform 0.2s;
}
.btn-chat:hover { background-color: #008f6f; transform: translateY(-2px); }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes scaleUp { from { transform: scale(0.9); opacity: 0; } to { transform: scale(1); opacity: 1; } }
</style>