<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="info-card">
      
      <div class="card-header">
        <h2>{{ chatType === 'group' ? 'Info del Grupo' : 'Info del Contacto' }}</h2>
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>

      <div class="info-body" v-if="loading">
        <div class="spinner"></div>
      </div>

      <div class="info-body" v-else>
        <div class="profile-header">
          <img :src="info.photo || info.Photo || DEFAULT_AVATAR" class="big-avatar" />
          <h3 class="chat-title">{{ info.name || info.Name }}</h3>
          <p class="chat-subtitle" v-if="chatType === 'group'">
            Grupo · {{ info.participants?.length || 0 }} participantes
          </p>
        </div>

        <div class="section-title">Participantes</div>
        <div class="participants-list">
          <div v-for="user in info.participants" :key="user.id" class="participant-row">
            <img :src="user.photo || DEFAULT_AVATAR" class="part-avatar" />
            <div class="part-info">
              <span class="part-name">{{ user.name }}</span>
              <span v-if="user.id === myId" class="you-badge">Tú</span>
            </div>
          </div>
        </div>
      </div>

    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { getChatInfo } from '../services/api';
import { store } from '../store.js';

const props = defineProps(['chatId']);
const emit = defineEmits(['close']);

const DEFAULT_AVATAR = 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png';
const myId = store.currentUser?.id || localStorage.getItem('userId');

const info = ref({});
const loading = ref(true);
const chatType = ref('direct');

onMounted(async () => {
  try {
    const res = await getChatInfo(props.chatId);
    info.value = res;
    chatType.value = res.type || res.Type;
  } catch (e) {
    console.error(e);
    alert("Error cargando info");
    emit('close');
  } finally {
    loading.value = false;
  }
});
</script>

<style scoped>
.modal-backdrop { position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.6); z-index: 2000; display: flex; justify-content: end; animation: fadeIn 0.2s; }
.info-card { width: 400px; max-width: 90%; height: 100%; background: #111b21; border-left: 1px solid #333; display: flex; flex-direction: column; animation: slideLeft 0.3s; }

.card-header { height: 60px; background: #202c33; display: flex; align-items: center; padding: 0 20px; color: #e9edef; font-weight: 500; justify-content: space-between; flex-shrink: 0; }
.btn-close { background: none; border: none; color: #aebac1; cursor: pointer; font-size: 1.2rem; }

.info-body { flex: 1; overflow-y: auto; padding: 20px; display: flex; flex-direction: column; align-items: center; }

.profile-header { display: flex; flex-direction: column; align-items: center; margin-bottom: 30px; width: 100%; }
.big-avatar { width: 200px; height: 200px; border-radius: 50%; object-fit: cover; margin-bottom: 15px; border: 4px solid #202c33; }
.chat-title { color: #e9edef; font-size: 1.5rem; margin: 0; font-weight: 400; text-align: center; }
.chat-subtitle { color: #8696a0; margin-top: 5px; font-size: 0.9rem; }

.section-title { width: 100%; color: #00a884; font-size: 0.9rem; margin-bottom: 10px; font-weight: 500; align-self: flex-start; }
.participants-list { width: 100%; display: flex; flex-direction: column; gap: 10px; }
.participant-row { display: flex; align-items: center; padding: 8px; border-radius: 8px; transition: background 0.2s; }
.participant-row:hover { background: #202c33; }
.part-avatar { width: 40px; height: 40px; border-radius: 50%; margin-right: 12px; object-fit: cover; }
.part-info { display: flex; flex-direction: column; }
.part-name { color: #e9edef; font-size: 1rem; }
.you-badge { font-size: 0.75rem; color: #8696a0; }

.spinner { border: 3px solid rgba(255,255,255,0.1); border-top: 3px solid #00a884; border-radius: 50%; width: 30px; height: 30px; animation: spin 1s linear infinite; margin-top: 50px; }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideLeft { from { transform: translateX(100%); } to { transform: translateX(0); } }
@keyframes spin { 100% { transform: rotate(360deg); } }
</style>