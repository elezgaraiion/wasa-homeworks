<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="search-card">
      
      <div class="card-header">
        <h2>Nuevo Mensaje</h2>
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>

      <div class="search-input-container">
        <span class="search-icon">🔍</span>
        <input 
          ref="inputRef"
          v-model="query" 
          type="text" 
          placeholder="Buscar usuarios..." 
          class="modern-input"
        />
      </div>

      <div class="users-list-wrapper">
        <div v-if="loading" class="status-msg">Buscando...</div>
        
        <div v-else-if="users.length === 0" class="status-msg">
          {{ query ? 'No se encontraron usuarios' : 'Escribe para buscar...' }}
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
/* ESTE Z-INDEX DEBE SER MUY ALTO */
.modal-backdrop {
  position: fixed; top: 0; left: 0; width: 100vw; height: 100vh;
  background: rgba(0, 0, 0, 0.6); backdrop-filter: blur(5px);
  z-index: 9999; /* <--- MUY IMPORTANTE */
  display: flex; justify-content: center; padding-top: 50px;
  animation: fadeIn 0.2s ease;
}

.search-card {
  width: 400px; max-width: 90%; background: #111b21; border: 1px solid #333;
  border-radius: 12px; display: flex; flex-direction: column;
  box-shadow: 0 20px 50px rgba(0,0,0,0.5); overflow: hidden; max-height: 80vh;
  animation: slideDown 0.3s cubic-bezier(0.19, 1, 0.22, 1);
}

.card-header { padding: 15px; background: #202c33; display: flex; justify-content: space-between; align-items: center; color: white;}
.btn-close { background: none; border: none; color: #aaa; font-size: 1.2rem; cursor: pointer; }
.search-input-container { padding: 10px; border-bottom: 1px solid #333; display: flex; align-items: center; background: #111b21;}
.search-icon { margin-right: 10px; color: #888; }
.modern-input { width: 100%; background: transparent; border: none; color: white; outline: none; }
.users-list-wrapper { flex: 1; overflow-y: auto; min-height: 100px; }
.status-msg { padding: 20px; text-align: center; color: #888; }
.user-row { display: flex; align-items: center; padding: 10px 15px; border-bottom: 1px solid #222; }
.user-avatar { width: 40px; height: 40px; border-radius: 50%; margin-right: 15px; background: #ddd; }
.user-name { color: white; }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideDown { from { transform: translateY(-20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
</style>