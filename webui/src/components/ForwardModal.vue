<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="forward-card">
      <div class="card-header">
        <h2>Forward message to...</h2>
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>

      <div class="chats-list">
        <div v-if="loading" class="status">Loading chats...</div>
        
        <div 
          v-else 
          v-for="chat in conversations" 
          :key="chat.id" 
          class="chat-row" 
          @click="handleForward(chat)"
        >
          <img :src="chat.photo || DEFAULT_AVATAR" class="chat-avatar"/>
          <div class="chat-name">{{ chat.name || chat.Name }}</div>
          <button class="btn-send">Send</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { getConversations, forwardMessage } from '../services/api';

const props = defineProps(['message', 'sourceChatId']);
const emit = defineEmits(['close', 'forwarded']);

const DEFAULT_AVATAR = 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png';
const conversations = ref([]);
const loading = ref(true);

onMounted(async () => {
  try {
    const res = await getConversations();
    conversations.value = res || [];
  } catch (e) {
    console.error(e);
  } finally {
    loading.value = false;
  }
});

async function handleForward(targetChat) {
  if (!confirm(`Forward to ${targetChat.name}?`)) return;
  
  try {
    const msgId = props.message.id || props.message.ID;
    
    await forwardMessage(props.sourceChatId, msgId, targetChat.id);
    
    alert("Message forwarded");
    emit('forwarded'); 
    emit('close');
  } catch (e) {
    alert("Error forwarding: " + e.message);
  }
}
</script>

<style scoped>
.modal-backdrop { position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.7); z-index: 10000; display: flex; justify-content: center; align-items: center; }
.forward-card { width: 400px; max-width: 90%; background: #111b21; border-radius: 12px; border: 1px solid #333; display: flex; flex-direction: column; height: 60vh; }
.card-header { padding: 15px; background: #202c33; display: flex; justify-content: space-between; align-items: center; color: white; }
.btn-close { background: none; border: none; color: #aaa; font-size: 1.5rem; cursor: pointer; }
.chats-list { flex: 1; overflow-y: auto; padding: 10px; }
.chat-row { display: flex; align-items: center; padding: 10px; border-bottom: 1px solid #222; cursor: pointer; }
.chat-row:hover { background: #202c33; }
.chat-avatar { width: 40px; height: 40px; border-radius: 50%; margin-right: 10px; }
.chat-name { flex: 1; color: #e9edef; }
.btn-send { background: #00a884; border: none; padding: 5px 12px; border-radius: 4px; color: #111; font-weight: bold; cursor: pointer; }
.status { color: #888; text-align: center; margin-top: 20px; }
</style>