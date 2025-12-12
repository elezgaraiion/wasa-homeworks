<template>
  <div class="dashboard-layout">
    
    <aside class="left-pane">
      <ConversationList 
        v-if="!showProfileEdit"
        :currentUser="store.currentUser" 
        @chatSelected="handleChatSelected"
        @openProfile="showProfileEdit = true"
      />
      <ProfileEdit v-else @close="showProfileEdit = false" />
    </aside>

    <main class="right-pane">
      
      <div v-if="showProfileEdit" class="welcome-placeholder">
        <img 
          src="https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png" 
          width="120" 
          class="placeholder-img"
        />
        <h2>Tu Perfil</h2>
        <p>Actualiza tu foto y tu nombre en el panel izquierdo.</p>
      </div>

      <ChatWindow 
        v-else-if="currentChat" 
        :chat="currentChat" 
      />

      <div v-else class="welcome-placeholder">
        <img 
          src="https://upload.wikimedia.org/wikipedia/commons/thumb/6/6b/WhatsApp.svg/1200px-WhatsApp.svg.png" 
          width="100" 
          class="placeholder-img"
        />
        <h2>WASA Web</h2>
        <p>Envía y recibe mensajes sin conectar tu teléfono.</p>
        <p class="encrypted-text">🔒 Cifrado de extremo a extremo (mentira, es un proyecto de clase)</p>
      </div>

    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { store } from '../store.js';
import ConversationList from '../components/ConversationList.vue';
import ProfileEdit from '../components/ProfileEdit.vue';
import ChatWindow from '../components/ChatWindow.vue';
import { getCurrentUser } from '../services/api';

const currentChat = ref(null);
const showProfileEdit = ref(false);

function handleChatSelected(chat) {
  // RESETEAMOS LA BOLA VERDE AL MOMENTO
  chat.unreadCount = 0;
  
  currentChat.value = chat;
  showProfileEdit.value = false;
}

onMounted(async () => {
  if (!store.currentUser) {
    try {
      const user = await getCurrentUser();
      store.login(user, localStorage.getItem('userId'));
    } catch (e) {
      console.error(e);
    }
  }
});
</script>

<style scoped>
.dashboard-layout { display: flex; height: 100vh; width: 100vw; background-color: #0b141a; overflow: hidden; }
.left-pane { width: 30%; min-width: 300px; max-width: 450px; height: 100%; border-right: 1px solid #333; display: flex; flex-direction: column; }
.right-pane { flex: 1; background-color: #222e35; border-left: 1px solid #333; position: relative; display: flex; flex-direction: column; }
.welcome-placeholder { width: 100%; height: 100%; display: flex; flex-direction: column; align-items: center; justify-content: center; text-align: center; color: #e9edef; z-index: 10; background-image: url("https://user-images.githubusercontent.com/15075759/28719144-86dc0f70-73b1-11e7-911d-60d70fcded21.png"); background-blend-mode: overlay; background-color: #222e35; }
.placeholder-img { opacity: 0.4; filter: grayscale(100%); margin-bottom: 20px; border-radius: 50%; }
.welcome-placeholder h2 { font-weight: 300; margin-bottom: 10px; font-size: 2rem; }
.welcome-placeholder p { color: #8696a0; }
.encrypted-text { margin-top: 40px; font-size: 0.8rem; color: #667781 !important; }
</style>