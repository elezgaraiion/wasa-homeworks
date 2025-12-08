<template>
  <div class="sidebar">
    
    <header class="sidebar-header">
      <div class="user-profile-section" @click="$emit('openProfile')" title="Editar mi perfil">
        <img 
          :src="avatarUrl" 
          @error="handleImageError"
          class="user-avatar" 
          alt="Avatar"
        />
        <div class="my-info">
          <span class="my-name">{{ userName }}</span>
        </div>
      </div>

      <div class="header-actions">
        <div class="icon group-add-icon" title="Crear grupo">
           <svg viewBox="0 0 24 24" width="24" height="24" class="">
            <path fill="currentColor" d="M15 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm-9-2V7H4v3H1v2h3v3h2v-3h3v-2H6zm9 4c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"></path>
          </svg>
        </div>
        <span class="icon logout-btn" @click.stop="handleLogout" title="Cerrar Sesión">🚪</span>
      </div>
    </header>

    <div class="search-bar">
      <div class="search-input-wrapper" @click="openSearch">
        <span class="search-icon">🔍</span>
        <input 
          type="text" 
          placeholder="Buscar o empezar un nuevo chat" 
          readonly 
          class="fake-input"
        />
      </div>
    </div>

    <div class="chats-container">
      <div v-if="loading" class="state-msg">Cargando...</div>
      
      <div v-else-if="conversations.length === 0" class="state-msg">
        <p>No tienes conversaciones.</p>
        <button class="btn-new" @click="openSearch">Buscar gente</button>
      </div>

      <div 
        v-else 
        v-for="chat in conversations" 
        :key="chat.id" 
        class="chat-item"
        :class="{ active: selectedChatId === chat.id }"
        @click="selectChat(chat.id)"
      >
        <img 
          :src="chat.photo || DEFAULT_AVATAR" 
          @error="handleImageError"
          class="chat-avatar" 
        />
        <div class="chat-info">
          <div class="chat-top">
            <span class="chat-name">{{ chat.name || chat.Name || 'Chat' }}</span>
            <span class="chat-time">{{ formatTime(chat.lastMessageAt || chat.LastMessageAt) }}</span>
          </div>
          <div class="chat-bottom">
            <span class="chat-preview">{{ chat.lastMessagePreview || chat.LastMessagePreview || '...' }}</span>
          </div>
        </div>
      </div>
    </div>

    <UserSearchModal 
      v-if="showSearchModal" 
      @close="showSearchModal = false" 
    />

  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import { getConversations } from '../services/api'; 
import { store } from '../store.js'; 
import UserSearchModal from './UserSearchModal.vue'; // <--- VERIFICA QUE ESTE ARCHIVO EXISTE

const DEFAULT_AVATAR = 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png';

const emit = defineEmits(['chatSelected', 'openProfile']);

// ESTADOS
const conversations = ref([]);
const loading = ref(true);
const selectedChatId = ref(null);
const showSearchModal = ref(false); // <--- Variable que abre el modal

// FUNCIÓN PARA ABRIR (CON LOG DE DEPURACIÓN)
function openSearch() {
  console.log("¡Click detectado! Abriendo modal..."); 
  showSearchModal.value = true;
}

const userName = computed(() => {
  const u = store.currentUser;
  return u?.name || u?.Name || 'Usuario';
});

const avatarUrl = computed(() => {
  const u = store.currentUser;
  return u?.photo || u?.Photo || DEFAULT_AVATAR;
});

function handleImageError(e) { e.target.src = DEFAULT_AVATAR; }

function formatTime(dateStr) {
  if (!dateStr || dateStr.startsWith('0001')) return '';
  return new Date(dateStr).toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

function handleLogout() {
  if (confirm('¿Cerrar sesión?')) store.logout();
}

onMounted(async () => {
  try {
    loading.value = true;
    conversations.value = await getConversations();
  } catch (error) {
    console.error("Error chats:", error);
  } finally {
    loading.value = false;
  }
});

function selectChat(id) {
  selectedChatId.value = id;
  emit('chatSelected', id);
}
</script>

<style scoped>
/* ESTRUCTURA */
.sidebar { display: flex; flex-direction: column; height: 100%; background-color: #111b21; border-right: 1px solid #2f3b43; color: #e9edef; position: relative; }

/* CABECERA */
.sidebar-header { height: 60px; background-color: #202c33; padding: 0 16px; display: flex; align-items: center; justify-content: space-between; flex-shrink: 0; }

/* PERFIL */
.user-profile-section { display: flex; align-items: center; cursor: pointer; padding: 5px; border-radius: 8px; max-width: 65%; transition: 0.2s; }
.user-profile-section:hover { background-color: #2a3942; }
.user-avatar { width: 40px; height: 40px; border-radius: 50%; object-fit: cover; margin-right: 10px; background-color: #dfe5e7; }
.my-name { font-weight: 600; font-size: 1rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

/* ICONOS */
.header-actions { display: flex; gap: 20px; align-items: center; }
.icon { color: #aebac1; font-size: 1.4rem; cursor: pointer; display: flex; align-items: center; }
.icon:hover { color: #fff; }
.logout-btn:hover { color: #f15c6d; }

/* BARRA DE BÚSQUEDA (ARREGLADA) */
.search-bar { padding: 8px 12px; border-bottom: 1px solid #202c33; }

.search-input-wrapper { 
  background-color: #202c33; 
  border-radius: 8px; 
  padding: 6px 12px; 
  display: flex; 
  align-items: center; 
  cursor: pointer; /* Mano al pasar ratón */
  transition: background 0.2s;
}
.search-input-wrapper:hover { background-color: #2a3942; }

.search-icon { margin-right: 10px; color: #8696a0; }

/* INPUT "FALSO" QUE DEJA PASAR EL CLICK */
.fake-input { 
  background: transparent; 
  border: none; 
  color: #d1d7db; 
  width: 100%; 
  outline: none; 
  cursor: pointer; 
  pointer-events: none; /* <--- TRUCO: El click pasa al div padre */
}

/* LISTA */
.chats-container { flex: 1; overflow-y: auto; }
.state-msg { padding: 20px; text-align: center; color: #888; }
.btn-new { background: #00a884; border: none; padding: 8px 16px; border-radius: 4px; color: #111; cursor: pointer; font-weight: bold; margin-top: 10px; }

/* ITEMS CHAT */
.chat-item { display: flex; align-items: center; padding: 0 15px; height: 72px; cursor: pointer; border-bottom: 1px solid #222; }
.chat-item:hover { background-color: #202c33; }
.chat-item.active { background-color: #2a3942; }
.chat-avatar { width: 49px; height: 49px; border-radius: 50%; margin-right: 15px; object-fit: cover; background-color: #dfe5e7; }
.chat-info { flex: 1; overflow: hidden; }
.chat-top { display: flex; justify-content: space-between; margin-bottom: 4px; }
.chat-name { font-weight: 500; font-size: 17px; }
.chat-time { font-size: 12px; color: #8696a0; }
.chat-preview { font-size: 14px; color: #8696a0; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }

.chats-container::-webkit-scrollbar { width: 5px; }
.chats-container::-webkit-scrollbar-thumb { background-color: #374045; }
</style>