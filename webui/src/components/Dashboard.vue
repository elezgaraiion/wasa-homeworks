<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <h3>Hola, {{ store.currentUser?.name }}</h3>
      <div>
        <button @click="$emit('openProfile')" class="btn-small">Perfil</button>
        <button @click="$emit('logout')" class="btn-small btn-logout">Salir</button>
      </div>
    </header>

    <div class="conv-list">
      <h3>Conversaciones</h3>
      <div v-if="loading">Cargando...</div>
      <div v-else-if="conversations.length === 0">No tienes conversaciones.</div>
      
      <ul v-else>
        <li v-for="conv in conversations" :key="conv.id" class="conversation">
          <img :src="conv.photo || 'https://via.placeholder.com/40'" class="avatar" />
          <div class="info">
            <strong>{{ conv.name || 'Chat sin nombre' }}</strong>
            <p class="preview">{{ conv.lastMessagePreview || '...' }}</p>
          </div>
          <span class="time">{{ formatDate(conv.lastMessageAt) }}</span>
        </li>
      </ul>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { getConversations } from '../services/api';
import { store } from '../store.js';

const conversations = ref([]);
const loading = ref(true);

function formatDate(dateStr) {
  if (!dateStr) return '';
  return new Date(dateStr).toLocaleTimeString([], { hour:'2-digit', minute:'2-digit' });
}

async function loadConversations() {
  loading.value = true;
  try {
    // Ya no pasamos userId, api.js lo maneja
    conversations.value = await getConversations();
  } catch (err) {
    console.error(err);
  } finally {
    loading.value = false;
  }
}

onMounted(loadConversations);
</script>

<style scoped>
.dashboard { display:flex; flex-direction:column; }
.dashboard-header { display:flex; justify-content:space-between; align-items:center; padding: 1rem; background: #eee; }
.conv-list { padding: 1rem; }
ul { list-style:none; padding:0; margin:0; }
.conversation { display:flex; align-items:center; padding:0.8rem 0; border-bottom:1px solid #f0f0f0; cursor: pointer; }
.conversation:hover { background-color: #f9f9f9; }
.avatar { width:50px; height:50px; border-radius:50%; margin-right:1rem; object-fit: cover; }
.info { flex:1; }
.preview { margin: 0; color: #666; font-size: 0.9rem; }
.time { font-size:0.75rem; color:#999; }
.btn-small { margin-left: 0.5rem; padding: 0.3rem 0.6rem; }
.btn-logout { background-color: #d32f2f; color: white; border: none; }
</style>