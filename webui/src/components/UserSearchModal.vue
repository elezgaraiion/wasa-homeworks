<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="search-card">
      
      <div class="card-header">
        <h2>Nuevo Chat</h2>
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>

      <div class="search-input-container">
        <span class="search-icon">🔍</span>
        <input 
          ref="inputRef"
          v-model="query" 
          type="text" 
          placeholder="Busca un usuario..." 
          class="modern-input"
        />
      </div>

      <div class="users-list-wrapper">
        <div v-if="loading" class="status-msg">Buscando...</div>
        
        <div v-else-if="users.length === 0" class="status-msg">
          {{ query ? 'No hay resultados' : 'Escribe un nombre...' }}
        </div>

        <div 
          v-else 
          v-for="user in users" 
          :key="user.id" 
          class="user-row"
        >
          <img 
            :src="user.photo || user.Photo || 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png'" 
            class="user-avatar" 
          />
          <span class="user-name">{{ user.name || user.Name }}</span>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, nextTick } from 'vue';
import { searchUsers } from '../services/api';

const emit = defineEmits(['close']);
const query = ref('');
const users = ref([]);
const loading = ref(false);
const inputRef = ref(null);

// Enfocar automáticamente al abrir
onMounted(() => {
  nextTick(() => inputRef.value?.focus());
});

watch(query, async (newVal) => {
  if (!newVal.trim()) {
    users.value = [];
    return;
  }
  
  loading.value = true;
  try {
    users.value = await searchUsers(newVal);
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
/* FONDO BORROSO (GLASSMORPHISM) */
.modal-backdrop {
  position: fixed;
  top: 0; left: 0; width: 100vw; height: 100vh;
  background: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(5px); /* Desenfoque elegante */
  z-index: 9999;
  display: flex;
  justify-content: center;
  align-items: flex-start;
  padding-top: 60px; /* Un poco separado del techo */
  animation: fadeIn 0.2s ease;
}

/* TARJETA FLOTANTE */
.search-card {
  width: 400px;
  max-width: 90%;
  background-color: #111b21;
  border: 1px solid #333;
  border-radius: 12px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 50px rgba(0,0,0,0.5);
  animation: slideDown 0.3s cubic-bezier(0.19, 1, 0.22, 1);
  overflow: hidden;
  max-height: 80vh;
}

.card-header {
  padding: 15px 20px;
  background-color: #202c33;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid #2f3b43;
}

.card-header h2 { margin: 0; font-size: 1.1rem; color: #e9edef; font-weight: 600; }
.btn-close { background: none; border: none; color: #aebac1; font-size: 1.2rem; cursor: pointer; }
.btn-close:hover { color: #f15c6d; }

.search-input-container {
  padding: 12px 20px;
  border-bottom: 1px solid #2f3b43;
  display: flex;
  align-items: center;
  background: #111b21;
}
.search-icon { margin-right: 10px; color: #8696a0; }
.modern-input { width: 100%; background: transparent; border: none; color: #e9edef; font-size: 1rem; outline: none; }

.users-list-wrapper { flex: 1; overflow-y: auto; min-height: 150px; }
.status-msg { padding: 30px; text-align: center; color: #8696a0; font-size: 0.9rem; }

.user-row {
  display: flex; align-items: center; padding: 12px 20px;
  border-bottom: 1px solid #222; cursor: default;
}
.user-row:hover { background-color: #202c33; }

.user-avatar { width: 45px; height: 45px; border-radius: 50%; object-fit: cover; margin-right: 15px; background-color: #dfe5e7; }
.user-name { color: #e9edef; font-weight: 500; }

.users-list-wrapper::-webkit-scrollbar { width: 5px; }
.users-list-wrapper::-webkit-scrollbar-thumb { background-color: #374045; }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideDown { from { transform: translateY(-30px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
</style>