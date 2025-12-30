<template>
  <div class="sidebar">
    
    <header class="sidebar-header">
      <div class="user-profile-section" @click="$emit('openProfile')" title="Editar mi perfil">
        <img :src="avatarUrl" @error="handleImageError" class="user-avatar" alt="Avatar" />
        <div class="my-info">
          <span class="my-name">{{ userName }}</span>
        </div>
      </div>
      
      <div class="header-actions">
        <span class="icon" @click="showGroupModal = true" title="Nuevo Grupo">➕👥</span>
        <span class="icon logout-btn" @click.stop="handleLogout" title="Cerrar Sesión">🚪</span>
      </div>
    </header>

    <div class="search-bar">
      <div class="search-input-wrapper" @click="showSearchModal = true">
        <span class="search-icon">🔍</span>
        <span class="fake-input-text">Buscar usuarios</span>
      </div>
    </div>

    <div class="chats-container">
      <div v-if="loading && conversations.length === 0" class="state-msg">Cargando...</div>
      <div v-else-if="conversations.length === 0" class="state-msg"><p>No tienes conversaciones.</p></div>
      <div v-else v-for="chat in conversations" :key="chat.id" class="chat-item" :class="{ active: selectedChatId === chat.id }" @click="selectChat(chat)">
        <img :src="chat.photo || DEFAULT_AVATAR" @error="handleImageError" class="chat-avatar" />
        <div class="chat-info">
          <div class="chat-row-top">
            <span class="chat-name">{{ chat.name || chat.Name || 'Chat' }}</span>
            <span class="chat-time" :class="{ 'green-time': chat.unreadCount > 0 }">{{ formatTime(chat.lastMessageAt || chat.LastMessageAt) }}</span>
          </div>
          <div class="chat-row-bottom">
            <div class="preview-wrapper">
              <span v-if="isMe(chat)" class="status-icon">
                 <span v-if="chat.lastMessageStatus === 'read'" style="color: #53bdeb;">✓✓</span>
                 <span v-else>✓</span>
              </span>
              <span v-else-if="chat.type === 'group' && chat.lastMessageSenderName" class="sender-prefix">{{ chat.lastMessageSenderName }}:</span>
              <span class="chat-preview-text">{{ chat.lastMessagePreview || chat.LastMessagePreview || '...' }}</span>
            </div>
            <div v-if="chat.unreadCount > 0" class="unread-badge">{{ chat.unreadCount }}</div>
          </div>
        </div>
      </div>
    </div>

    <Teleport to="body">
      <UserSearchModal 
        v-if="showSearchModal" 
        @close="showSearchModal = false"
        @chatStarted="onChatCreated" 
      />

      <CreateGroupModal
        v-if="showGroupModal"
        @close="showGroupModal = false"
        @groupCreated="onChatCreated" 
      />
    </Teleport>

  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, computed } from 'vue';
import { getConversations } from '../services/api'; 
import { store } from '../store.js'; 
import UserSearchModal from './UserSearchModal.vue'; 
import CreateGroupModal from './CreateGroupModal.vue';

const DEFAULT_AVATAR = 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png';
const emit = defineEmits(['chatSelected', 'openProfile']);

const conversations = ref([]);
const loading = ref(true);
const selectedChatId = ref(null);
const showSearchModal = ref(false);
const showGroupModal = ref(false); 
let refreshInterval = null;

const userName = computed(() => store.currentUser?.name || store.currentUser?.Name || 'Usuario');
const avatarUrl = computed(() => store.currentUser?.photo || store.currentUser?.Photo || DEFAULT_AVATAR);

function handleImageError(e) { e.target.src = DEFAULT_AVATAR; }

function isMe(chat) {
  const myId = store.currentUser?.id || localStorage.getItem('userId');
  return chat.lastMessageSenderId === myId;
}

function formatTime(dateStr) {
  if (!dateStr || dateStr === '' || dateStr.startsWith('0001')) return '';
  const safeDateStr = dateStr.replace(' ', 'T');
  let date = new Date(safeDateStr);
  if (isNaN(date.getTime())) date = new Date(safeDateStr + 'Z');
  if (isNaN(date.getTime())) return '';
  const now = new Date();
  const isToday = date.getDate() === now.getDate() && date.getMonth() === now.getMonth() && date.getFullYear() === now.getFullYear();
  return isToday 
    ? date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
    : date.toLocaleDateString([], { day: '2-digit', month: '2-digit', year: '2-digit' });
}

function handleLogout() {
  if (confirm('¿Cerrar sesión?')) store.logout();
}

async function loadConversations() {
  try {
    const res = await getConversations();
    conversations.value = res || [];
  } catch (error) {
    console.error("Error chats:", error);
  } finally {
    loading.value = false;
  }
}

async function onChatCreated(chat) {
  showSearchModal.value = false;
  showGroupModal.value = false; 

  const existing = conversations.value.find(c => c.id === chat.id);
  if (!existing) {
    conversations.value.unshift(chat);
  }
  selectChat(chat);
  await loadConversations();
}

function selectChat(chat) {
  selectedChatId.value = chat.id;
  emit('chatSelected', chat);
}

onMounted(() => {
  loadConversations();
  refreshInterval = setInterval(loadConversations, 4000);
});

onUnmounted(() => {
  if (refreshInterval) clearInterval(refreshInterval);
});

defineExpose({ refreshList: loadConversations });
</script>

<style scoped>
/* AÑADIDO UN MARGEN ENTRE ICONOS */
.header-actions { display: flex; gap: 20px; align-items: center; } 

/* (El resto de estilos es IDÉNTICO, solo asegúrate de tener el css previo) */
.sidebar { display: flex; flex-direction: column; height: 100%; background-color: #111b21; border-right: 1px solid #2f3b43; color: #e9edef; position: relative; }
.sidebar-header { height: 60px; background-color: #202c33; padding: 0 16px; display: flex; align-items: center; justify-content: space-between; flex-shrink: 0; }
.user-profile-section { display: flex; align-items: center; cursor: pointer; padding: 5px; border-radius: 8px; max-width: 80%; transition: 0.2s; }
.user-profile-section:hover { background-color: #2a3942; }
.user-avatar { width: 40px; height: 40px; border-radius: 50%; object-fit: cover; margin-right: 10px; background-color: #dfe5e7; }
.my-name { font-weight: 600; font-size: 1rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.icon { color: #aebac1; font-size: 1.4rem; cursor: pointer; }
.icon:hover { color: #fff; }
.logout-btn:hover { color: #f15c6d; }
.search-bar { padding: 8px 12px; border-bottom: 1px solid #202c33; }
.search-input-wrapper { background-color: #202c33; border-radius: 8px; padding: 8px 12px; display: flex; align-items: center; cursor: pointer; transition: background 0.2s; }
.search-input-wrapper:hover { background-color: #2a3942; }
.search-icon { margin-right: 10px; color: #8696a0; font-size: 0.9rem; }
.fake-input-text { color: #8696a0; font-size: 0.95rem; user-select: none; }
.chats-container { flex: 1; overflow-y: auto; }
.state-msg { padding: 20px; text-align: center; color: #888; }
.chat-item { display: flex; align-items: center; padding: 0 15px; height: 72px; cursor: pointer; border-bottom: 1px solid #222; }
.chat-item:hover { background-color: #202c33; }
.chat-item.active { background-color: #2a3942; }
.chat-avatar { width: 49px; height: 49px; border-radius: 50%; margin-right: 15px; object-fit: cover; background-color: #dfe5e7; }
.chat-info { flex: 1; overflow: hidden; display: flex; flex-direction: column; justify-content: center; gap: 4px; }
.chat-row-top { display: flex; justify-content: space-between; align-items: center; }
.chat-row-bottom { display: flex; justify-content: space-between; align-items: center; }
.chat-name { font-weight: 500; font-size: 17px; color: #e9edef; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; margin-right: 10px; }
.chat-time { font-size: 12px; color: #8696a0; white-space: nowrap; flex-shrink: 0; }
.chat-time.green-time { color: #00a884; font-weight: bold; }
.preview-wrapper { display: flex; align-items: center; overflow: hidden; flex: 1; margin-right: 10px; font-size: 14px; color: #8696a0; }
.status-icon { margin-right: 3px; font-size: 0.8rem; min-width: 16px; display: flex; align-items: center; }
.sender-prefix { margin-right: 4px; color: #e9edef; font-weight: 500; white-space: nowrap; }
.chat-preview-text { white-space: nowrap; overflow: hidden; text-overflow: ellipsis; flex: 1; }
.unread-badge { background-color: #00a884; color: #111b21; font-size: 0.75rem; font-weight: bold; border-radius: 50%; min-width: 20px; height: 20px; display: flex; align-items: center; justify-content: center; padding: 0 5px; flex-shrink: 0; animation: popIn 0.2s; }
@keyframes popIn { from { transform: scale(0); } to { transform: scale(1); } }
.chats-container::-webkit-scrollbar { width: 5px; }
.chats-container::-webkit-scrollbar-thumb { background-color: #374045; }
</style>