<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="search-card">
      
      <div class="card-header">
        <h2>Users</h2>
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>

      <div class="search-input-container">
        <span class="search-icon">🔍</span>
        <input 
          ref="inputRef"
          v-model="query" 
          type="text" 
          placeholder="Filter users..." 
          class="modern-input"
        />
      </div>

      <div class="users-list-wrapper">
        <div v-if="loading" class="status-msg">Loading contacts...</div>
        
        <div v-else-if="users.length === 0" class="status-msg">
           No users found.
        </div>

        <div 
          v-else 
          v-for="user in users" 
          :key="user.id" 
          class="user-row"
          @click="inspectUser(user)" 
        >
          <img 
            :src="user.photo || user.Photo || 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png'" 
            class="user-avatar" 
            @error="handleImageError"
          />
          <span class="user-name">{{ user.name || user.Name }}</span>
        </div>
      </div>

    </div>

    <UserProfileReadOnly 
      v-if="selectedUser" 
      :user="selectedUser"
      @close="selectedUser = null"
      @chatStarted="handleChatStarted" 
    />

  </div>
</template>

<script>
import UserProfileReadOnly from './UserProfileReadOnly.vue';

export default {
  name: "UserSearchModal",
  components: {
    UserProfileReadOnly
  },
  emits: ['close', 'chatStarted'],
  data() {
    return {
      query: '',
      users: [],
      loading: false,
      selectedUser: null
    };
  },
  watch: {
    // Cada vez que escribes (o borras), buscamos de nuevo
    query() {
      this.fetchUsers();
    }
  },
  methods: {
    async fetchUsers() {
      this.loading = true;
      try {
        // Al enviar query vacío, el backend devuelve TODOS (ordenados por nombre)
        // Al enviar texto, el backend filtra por ese texto
        const res = await this.$axios.get('/users', { params: { q: this.query } });
        this.users = res.data || [];
      } catch (e) {
        console.error(e);
      } finally {
        this.loading = false;
      }
    },

    inspectUser(user) {
      this.selectedUser = user;
    },

    handleChatStarted(chat) {
      this.$emit('chatStarted', chat);
      this.$emit('close'); 
    },

    handleImageError(e) {
      e.target.src = 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png';
    }
  },
  mounted() {
    // 1. Poner el foco
    this.$nextTick(() => {
        if(this.$refs.inputRef) this.$refs.inputRef.focus();
    });
    
    // 2. Cargar la lista completa nada más abrir
    this.fetchUsers();
  }
};
</script>

<style scoped>
.modal-backdrop {
  position: fixed; top: 0; left: 0; width: 100vw; height: 100vh;
  background: rgba(0, 0, 0, 0.6); backdrop-filter: blur(5px);
  z-index: 9999; display: flex; justify-content: center; align-items: flex-start;
  padding-top: 60px; animation: fadeIn 0.2s ease;
}
.search-card {
  width: 400px; max-width: 90%; background-color: #111b21; border: 1px solid #333;
  border-radius: 12px; display: flex; flex-direction: column;
  box-shadow: 0 20px 50px rgba(0,0,0,0.5); animation: slideDown 0.3s cubic-bezier(0.19, 1, 0.22, 1);
  overflow: hidden; max-height: 80vh; /* Limita la altura para que el scroll funcione */
}
.card-header {
  padding: 15px 20px; background-color: #202c33; display: flex;
  justify-content: space-between; align-items: center; border-bottom: 1px solid #2f3b43;
}
.card-header h2 { margin: 0; font-size: 1.1rem; color: #e9edef; font-weight: 600; }
.btn-close { background: none; border: none; color: #aebac1; font-size: 1.2rem; cursor: pointer; }
.btn-close:hover { color: #f15c6d; }
.search-input-container {
  padding: 12px 20px; border-bottom: 1px solid #2f3b43;
  display: flex; align-items: center; background: #111b21;
}
.search-icon { margin-right: 10px; color: #8696a0; }
.modern-input { width: 100%; background: transparent; border: none; color: #e9edef; font-size: 1rem; outline: none; }

/* Wrapper de la lista con scroll */
.users-list-wrapper { 
  flex: 1; 
  overflow-y: auto; 
  min-height: 150px; /* Altura mínima estética */
}

.status-msg { padding: 30px; text-align: center; color: #8696a0; font-size: 0.9rem; }
.user-row {
  display: flex; align-items: center; padding: 12px 20px;
  border-bottom: 1px solid #222; cursor: pointer; transition: background 0.2s;
}
.user-row:hover { background-color: #202c33; }
.user-avatar { width: 45px; height: 45px; border-radius: 50%; object-fit: cover; margin-right: 15px; background-color: #dfe5e7; }
.user-name { color: #e9edef; font-weight: 500; font-size: 1rem; }

/* Scrollbar personalizado */
.users-list-wrapper::-webkit-scrollbar { width: 6px; }
.users-list-wrapper::-webkit-scrollbar-thumb { background-color: #374045; border-radius: 3px; }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideDown { from { transform: translateY(-30px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
</style>