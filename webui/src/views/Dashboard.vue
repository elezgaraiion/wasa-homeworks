<template>
  <div class="dashboard-layout">
    <aside class="left-pane">
      
      <ConversationList 
        v-if="!showProfileEdit"
        :currentUser="store.currentUser" 
        @chatSelected="handleChatSelected"
        @openProfile="showProfileEdit = true"
      />

      <ProfileEdit 
        v-else 
        @close="showProfileEdit = false"
      />

    </aside>

    <main class="right-pane">
      
      <div v-if="showProfileEdit" class="welcome-placeholder">
        <img 
          src="https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png" 
          width="120" 
          style="opacity: 0.3; filter: grayscale(100%); margin-bottom: 20px; border-radius: 50%;" 
        />
        <h2>Tu Perfil</h2>
        <p>Actualiza tu foto y tu nombre en el panel izquierdo.</p>
      </div>

      <div v-else-if="!currentChatId" class="welcome-placeholder">
        <img 
          src="https://upload.wikimedia.org/wikipedia/commons/thumb/6/6b/WhatsApp.svg/1200px-WhatsApp.svg.png" 
          width="100" 
          style="opacity: 0.4; filter: grayscale(100%); margin-bottom: 20px;" 
        />
        <h2>WASA Web</h2>
        <p>Envía y recibe mensajes sin conectar tu teléfono.</p>
        <p class="encrypted-text">🔒 Cifrado de extremo a extremo (mentira, es un proyecto de clase)</p>
      </div>

      <div v-else class="chat-placeholder">
        <h3>Cargando chat ID: {{ currentChatId }}...</h3>
        <p>(Aquí implementaremos los mensajes en el siguiente paso)</p>
      </div>
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { store } from '../store.js';
import ConversationList from '../components/ConversationList.vue';
import ProfileEdit from '../components/ProfileEdit.vue';
import { getCurrentUser } from '../services/api';

const currentChatId = ref(null);
const showProfileEdit = ref(false); 

function handleChatSelected(chatId) {
  currentChatId.value = chatId;
}

onMounted(async () => {
  if (!store.currentUser) {
    try {
      const user = await getCurrentUser();
      store.login(user, localStorage.getItem('userId'));
    } catch (e) {
      console.error("Error obteniendo usuario", e);
    }
  }
});
</script>

<style scoped>
.dashboard-layout {
  display: flex;
  height: 100vh;
  width: 100vw;
  background-color: #0b141a; 
  overflow: hidden;
}

.left-pane {
  width: 30%; 
  min-width: 300px;
  max-width: 450px;
  height: 100%;
  border-right: 1px solid #333; 
  display: flex;      
  flex-direction: column;
}

.right-pane {
  flex: 1; 
  background-color: #222e35;
  border-left: 1px solid #333;
  display: flex;
  align-items: center;
  justify-content: center;
  position: relative;
  /* Fondo de chat (opcional) */
  background-image: url("https://user-images.githubusercontent.com/15075759/28719144-86dc0f70-73b1-11e7-911d-60d70fcded21.png");
  background-blend-mode: overlay;
  background-size: contain;
}

/* REUTILIZAMOS LA CLASE PARA AMBOS ESTADOS (PERFIL Y BIENVENIDA) */
.welcome-placeholder {
  text-align: center;
  color: #e9edef;
  animation: fadeIn 0.3s ease;
  z-index: 10; /* Asegura que esté por encima del fondo */
}

.welcome-placeholder h2 {
  font-weight: 300; /* Fuente fina como pediste */
  margin-bottom: 10px;
  font-size: 2rem;
}

.welcome-placeholder p {
  color: #8696a0;
}

.encrypted-text {
  margin-top: 40px;
  font-size: 0.8rem;
  color: #667781 !important;
}

.chat-placeholder {
  color: white;
  background: rgba(0,0,0,0.6);
  padding: 20px;
  border-radius: 8px;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(10px); }
  to { opacity: 1; transform: translateY(0); }
}
</style>