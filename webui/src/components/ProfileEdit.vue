<template>
  <div class="profile-overlay">
    <div class="profile-card">
      
      <div class="card-header">
        <h2>Editar Perfil</h2>
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>

      <div class="avatar-section">
        <label class="avatar-wrapper" title="Click para cambiar foto">
          <img 
            :src="previewPhoto || currentAvatar" 
            @error="handleImageError"
            class="avatar-img" 
            alt="Perfil"
          />
          
          <div class="avatar-overlay">
            <span class="camera-icon">📷</span>
          </div>

          <input type="file" @change="handleFileChange" accept="image/jpeg" hidden />
        </label>
        
        <p class="avatar-hint">Toca la imagen para cambiarla</p>
      </div>

      <div class="form-section">
        <div class="input-group">
          <label>Tu Nombre</label>
          <input 
            v-model="newName" 
            type="text" 
            placeholder="Escribe tu nombre público" 
            @keyup.enter="saveName"
          />
          <span class="focus-border"></span>
        </div>

        <div class="info-group">
          <label>ID de Usuario</label>
          <div class="id-box">
            {{ store.currentUser?.id || '...' }}
            <span class="copy-icon" title="Copiar ID">📋</span>
          </div>
        </div>
      </div>

      <div class="actions">
        <button class="btn-cancel" @click="$emit('close')">Cancelar</button>
        <button 
          class="btn-save" 
          @click="saveName" 
          :disabled="loading || !isNameChanged"
        >
          {{ loading ? 'Guardando...' : 'Guardar Cambios' }}
        </button>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue';
import { store } from '../store.js';
import { updateUserName, updatePhoto } from '../services/api';

const DEFAULT_AVATAR = 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png';

const emit = defineEmits(['close']);

const newName = ref(store.currentUser?.name || store.currentUser?.Name || '');
const previewPhoto = ref(null);
const loading = ref(false);

const currentAvatar = computed(() => {
  const u = store.currentUser;
  return u?.photo || u?.Photo || DEFAULT_AVATAR;
});

function handleImageError(e) {
  e.target.src = DEFAULT_AVATAR;
}

const isNameChanged = computed(() => {
  const currentName = store.currentUser?.name || store.currentUser?.Name;
  return newName.value.trim() !== '' && newName.value !== currentName;
});

async function saveName() {
  if (!isNameChanged.value) return;
  
  loading.value = true;
  try {
    const updatedUser = await updateUserName(newName.value);
    store.updateUser(updatedUser);
    alert('Nombre actualizado correctamente');
  } catch (e) {
    alert('Error: ' + e.message);
  } finally {
    loading.value = false;
  }
}

async function handleFileChange(event) {
  const file = event.target.files[0];
  if (!file) return;

  previewPhoto.value = URL.createObjectURL(file);
  loading.value = true;

  try {
    const updatedUser = await updatePhoto(file);
    store.updateUser(updatedUser);
  } catch (e) {
    alert('Error al subir imagen: ' + e.message);
    previewPhoto.value = null; 
  } finally {
    loading.value = false;
  }
}
</script>

<style scoped>
/* FONDO */
.profile-overlay {
  position: fixed; top: 0; left: 0; width: 100vw; height: 100vh;
  background: rgba(0, 0, 0, 0.85); backdrop-filter: blur(8px);
  z-index: 1000; display: flex; justify-content: center; align-items: center;
  animation: fadeIn 0.3s ease;
}

/* TARJETA */
.profile-card {
  width: 100%; max-width: 450px; background: #111b21; border: 1px solid #333;
  border-radius: 16px; padding: 2rem; box-shadow: 0 25px 50px -12px rgba(0, 0, 0, 0.5);
  display: flex; flex-direction: column; gap: 1.5rem;
  animation: slideUp 0.4s cubic-bezier(0.16, 1, 0.3, 1);
}

.card-header { display: flex; justify-content: space-between; align-items: center; margin-bottom: 1rem; }
.card-header h2 { margin: 0; color: #e9edef; font-weight: 300; font-size: 1.8rem; }
.btn-close { background: transparent; border: none; color: #8696a0; font-size: 1.5rem; cursor: pointer; }
.btn-close:hover { color: #f15c6d; }

/* AVATAR */
.avatar-section { display: flex; flex-direction: column; align-items: center; }

.avatar-wrapper {
  position: relative; width: 160px; height: 160px; border-radius: 50%;
  overflow: hidden; border: 4px solid #202c33; box-shadow: 0 0 20px rgba(0, 168, 132, 0.2);
  transition: transform 0.3s; cursor: pointer;
  /* Fondo para que la silueta gris se vea bien */
  background-color: #dfe5e7; 
}

.avatar-wrapper:hover { transform: scale(1.05); border-color: #00a884; }

.avatar-img { width: 100%; height: 100%; object-fit: cover; }

.avatar-overlay {
  position: absolute; top: 0; left: 0; width: 100%; height: 100%;
  background: rgba(0, 0, 0, 0.4); display: flex; justify-content: center; align-items: center;
  opacity: 0; transition: opacity 0.2s;
}

.avatar-wrapper:hover .avatar-overlay { opacity: 1; }
.camera-icon { font-size: 2rem; filter: drop-shadow(0 2px 4px rgba(0,0,0,0.5)); }
.avatar-hint { margin-top: 10px; color: #8696a0; font-size: 0.9rem; }

/* INPUTS */
.input-group { position: relative; margin-bottom: 1.5rem; }
.input-group label { display: block; color: #00a884; font-size: 0.85rem; font-weight: bold; text-transform: uppercase; margin-bottom: 8px; letter-spacing: 1px; }
.input-group input { width: 100%; background: #202c33; border: none; border-radius: 8px; padding: 12px 16px; color: #e9edef; font-size: 1.1rem; outline: none; transition: background 0.2s; }
.input-group input:focus { background: #2a3942; box-shadow: 0 0 0 2px #00a884; }

/* ID BOX */
.info-group label { color: #8696a0; font-size: 0.85rem; margin-bottom: 5px; display: block; }
.id-box { background: rgba(255,255,255,0.05); padding: 10px; border-radius: 6px; color: #8696a0; font-family: monospace; font-size: 0.9rem; display: flex; justify-content: space-between; }

/* BOTONES */
.actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem; }
.btn-cancel { background: transparent; color: #e9edef; border: 1px solid #333; padding: 10px 20px; border-radius: 24px; cursor: pointer; font-weight: 500; transition: all 0.2s; }
.btn-cancel:hover { background: rgba(255,255,255,0.1); }
.btn-save { background: #00a884; color: #111b21; border: none; padding: 10px 24px; border-radius: 24px; font-weight: bold; cursor: pointer; box-shadow: 0 4px 15px rgba(0, 168, 132, 0.4); transition: transform 0.2s, box-shadow 0.2s; }
.btn-save:hover:not(:disabled) { transform: translateY(-2px); box-shadow: 0 6px 20px rgba(0, 168, 132, 0.6); }
.btn-save:disabled { background: #333; color: #666; box-shadow: none; cursor: default; }

/* ANIMACIONES */
@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideUp { from { transform: translateY(50px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
</style>