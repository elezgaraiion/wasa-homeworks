<template>
  <div class="chat-window">
    
    <header class="chat-header">
      <div class="header-left">
        <img :src="chat.photo || 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png'" class="header-avatar" />
        <div class="header-info">
          <span class="header-name">{{ chat.name || chat.Name || 'Chat' }}</span>
          <span class="header-status">
            {{ chat.type === 'group' ? 'Tú y otros...' : 'en línea' }}
          </span>
        </div>
      </div>
      <div class="header-actions">
        <span>🔍</span>
        <span>⋮</span>
      </div>
    </header>

    <div class="messages-area" ref="messagesContainer">
      
      <div v-if="loading && messages.length === 0" class="loading-msgs">
        <div class="spinner-small"></div>
      </div>
      
      <div v-else-if="messages.length === 0" class="empty-chat">
        <span class="empty-badge">No hay mensajes aún. Envía el primero.</span>
      </div>

      <div v-else v-for="msg in messages" :key="msg.id || msg.ID" class="message-row" :class="{ 'my-message': isMe(msg) }">
        <div class="message-bubble">
          
          <div v-if="shouldShowSenderName(msg)" class="sender-name">
            {{ msg.sender?.name || msg.Sender?.Name || '...' }}
          </div>

          <div class="message-content">
            <div v-if="msg.photo || msg.Photo" class="photo-wrapper">
              <img :src="msg.photo || msg.Photo" class="message-image" @load="scrollToBottom" />
              <p v-if="msg.text || msg.Text" class="caption-text">{{ msg.text || msg.Text }}</p>
            </div>
            <div v-else class="text-content">
              {{ msg.text || msg.Text }}
            </div>
          </div>

          <div class="message-meta">
            <span class="message-time">{{ formatTime(msg.createdAt || msg.CreatedAt) }}</span>
            <span v-if="isMe(msg)" class="check-icon">
               <span :style="{ color: (msg.status === 'read' || msg.Status === 'read') ? '#53bdeb' : 'inherit' }">✓✓</span>
            </span>
          </div>

        </div>
      </div>
    </div>

    <footer class="chat-footer">
      <div class="input-wrapper">
        <input 
            v-model="newMessage" 
            @keyup.enter="handleSend" 
            type="text" 
            placeholder="Escribe un mensaje..." 
            :disabled="sending"
            ref="inputRef"
        />
      </div>
      <span class="footer-icon send-btn" @click="handleSend" :class="{ 'disabled': !newMessage.trim() }">➤</span>
    </footer>

  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue';
// IMPORTAMOS LA FUNCIÓN DE API
import { getConversationMessages, sendMessage } from '../services/api';
import { store } from '../store.js';

const props = defineProps(['chat']); 
const messages = ref([]);
const loading = ref(false);
const sending = ref(false);
const newMessage = ref('');
const messagesContainer = ref(null);
const inputRef = ref(null);
let pollInterval = null;

function isMe(msg) {
  const myId = store.currentUser?.id || localStorage.getItem('userId');
  const senderId = msg.sender?.id || msg.Sender?.ID || msg.sender_id || msg.SenderID;
  return senderId === myId;
}

function shouldShowSenderName(msg) {
  return props.chat.type === 'group' && !isMe(msg);
}

function formatTime(dateStr) {
  if (!dateStr) return '';
  const safeDate = dateStr.replace(' ', 'T');
  let date = new Date(safeDate.includes('Z') ? safeDate : safeDate + 'Z');
  if (isNaN(date.getTime())) date = new Date(safeDate); 
  if (isNaN(date.getTime())) return '';
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

async function loadMessages(isBackgroundUpdate = false) {
  if (!props.chat?.id) return;
  
  // Solo mostramos loading si no hay mensajes (primera carga)
  if (!isBackgroundUpdate && messages.length === 0) loading.value = true;
  
  try {
    const res = await getConversationMessages(props.chat.id);
    const list = res || [];
    // Invertimos porque la API suele devolver del más nuevo al más viejo
    // y en el chat queremos los viejos arriba y los nuevos abajo
    const orderedList = list.reverse();

    // Solo hacemos scroll si han llegado mensajes nuevos
    const shouldScroll = messages.value.length !== orderedList.length;
    
    messages.value = orderedList;
    
    if (shouldScroll) scrollToBottom();

  } catch (e) {
    // console.error(e);
  } finally {
    loading.value = false;
  }
}

// --- FUNCIÓN DE ENVIAR ---
async function handleSend() {
  const text = newMessage.value.trim();
  if (!text) return;

  // 1. Limpieza visual inmediata (UX rápida)
  newMessage.value = '';
  sending.value = true;

  try {
    // 2. Llamada a la API
    const sentMsg = await sendMessage(props.chat.id, text);

    // 3. Añadir a la lista localmente (para no esperar al polling)
    // El backend nos devuelve el mensaje creado, lo metemos al array
    if (sentMsg) {
        messages.value.push(sentMsg);
        scrollToBottom();
    } else {
        // Fallback por si acaso
        await loadMessages(true);
    }

    // Enfocar de nuevo el input
    nextTick(() => inputRef.value?.focus());

  } catch (e) {
    console.error("Error envío:", e);
    alert("No se pudo enviar");
    newMessage.value = text; // Restaurar texto si falló
  } finally {
    sending.value = false;
  }
}

function scrollToBottom() {
  nextTick(() => {
    if (messagesContainer.value) {
        messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
    }
  });
}

// CICLO DE VIDA
onMounted(() => {
  loadMessages();
  // Focusear el input al entrar
  nextTick(() => inputRef.value?.focus());
  // Polling cada 3 segundos para ver mensajes nuevos
  pollInterval = setInterval(() => loadMessages(true), 3000);
});

onUnmounted(() => {
  if (pollInterval) clearInterval(pollInterval);
});

// Si cambiamos de chat sin desmontar el componente
watch(() => props.chat.id, () => {
  messages.value = [];
  loadMessages();
  nextTick(() => inputRef.value?.focus());
});
</script>

<style scoped>
.chat-window { display: flex; flex-direction: column; height: 100%; width: 100%; background-color: #0b141a; }
.chat-header { height: 60px; background-color: #202c33; padding: 0 16px; display: flex; align-items: center; border-left: 1px solid #333; flex-shrink: 0; }
.header-left { display: flex; align-items: center; }
.header-avatar { width: 40px; height: 40px; border-radius: 50%; margin-right: 15px; object-fit: cover; }
.header-info { display: flex; flex-direction: column; }
.header-name { color: #e9edef; font-weight: bold; }
.header-status { color: #8696a0; font-size: 0.8rem; }
.header-actions { color: #aebac1; font-size: 1.2rem; gap: 20px; display: flex; margin-left: auto; }

.messages-area { 
  flex: 1; 
  overflow-y: auto; 
  padding: 20px; 
  display: flex; 
  flex-direction: column; 
  gap: 5px; 
  background-image: url("https://user-images.githubusercontent.com/15075759/28719144-86dc0f70-73b1-11e7-911d-60d70fcded21.png");
  background-repeat: repeat;
  background-blend-mode: overlay;
  background-color: #0b141a; 
}

.loading-msgs { display: flex; justify-content: center; padding: 20px; }
.spinner-small { width: 20px; height: 20px; border: 2px solid rgba(255,255,255,0.3); border-top-color: #00a884; border-radius: 50%; animation: spin 0.8s linear infinite; }

.empty-chat { display: flex; justify-content: center; margin-top: 20px; }
.empty-badge { background: #202c33; color: #ffd279; padding: 5px 12px; border-radius: 8px; font-size: 0.8rem; box-shadow: 0 1px 2px rgba(0,0,0,0.3); }

.message-row { display: flex; width: 100%; margin-bottom: 2px; }
.message-row.my-message { justify-content: flex-end; }
.message-bubble { max-width: 65%; padding: 6px 9px; border-radius: 8px; background-color: #202c33; color: white; border-top-left-radius: 0; position: relative; box-shadow: 0 1px 0.5px rgba(0,0,0,0.13); }
.my-message .message-bubble { background-color: #005c4b; border-top-left-radius: 8px; border-top-right-radius: 0; }

.sender-name { font-size: 0.8rem; color: #d65c3e; font-weight: 500; margin-bottom: 2px; line-height: 1.2; }
.message-content { font-size: 0.9rem; line-height: 1.3; }
.message-image { max-width: 100%; border-radius: 6px; margin-bottom: 4px; display: block; }
.caption-text { margin: 4px 0 0 0; }

.message-meta { float: right; margin-left: 8px; margin-top: 4px; display: flex; align-items: center; gap: 3px; position: relative; top: 4px; }
.message-time { font-size: 0.65rem; color: rgba(255,255,255,0.6); }
.check-icon { font-size: 0.75rem; color: #8696a0; }

.chat-footer { height: 60px; background-color: #202c33; display: flex; align-items: center; padding: 0 16px; }
.input-wrapper { flex: 1; background: #2a3942; border-radius: 8px; padding: 9px 12px; margin-right: 10px; }
.input-wrapper input { background: transparent; border: none; color: white; width: 100%; outline: none; font-size: 1rem; }
.send-btn { color: #8696a0; font-size: 1.5rem; cursor: pointer; transition: color 0.2s; }
.send-btn:hover { color: #00a884; }
.send-btn.disabled { color: #444; cursor: default; }

@keyframes spin { 100% { transform: rotate(360deg); } }
</style>