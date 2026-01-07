<template>
  <div class="profile-overlay" @click.self="$emit('close')">
    <div class="profile-card">
      
      <div class="card-header">
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>

      <div class="avatar-section">
        <div class="avatar-wrapper">
          <img :src="resolveUrl(user.photo) || 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png'" class="avatar-img" />
        </div>
      </div>

      <div class="info-section">
        <h3 class="user-name">{{ user.name || user.Name }}</h3>
        <p class="user-bio">Hi! I am using Wasa Web.</p>
      </div>

      <div class="actions">
        <button class="btn-chat" @click="handleEnterChat" :disabled="loading">
          <span v-if="loading" class="spinner-mini"></span>
          <span v-else class="icon">💬</span> 
          {{ loading ? 'Opening...' : 'Send Message' }}
        </button>
      </div>

    </div>
  </div>
</template>

<script>
export default {
  name: "UserProfileReadOnly",
  props: ['user'],
  emits: ['close', 'chatStarted'],
  data() {
    return {
      loading: false
    };
  },
  methods: {
    resolveUrl(path) {
        if (!path) return '';
        if (path.startsWith('http')) return path;
        return this.$axios.defaults.baseURL + path;
    },
    async handleEnterChat() {
      this.loading = true;
      try {
        const targetId = this.user.id || this.user.ID;
        
        const response = await this.$axios.post('/chats', { targetUserId: targetId });
        const chat = response.data;

        this.$emit('chatStarted', chat);
        this.$emit('close');
      } catch (e) {
        alert("Error opening chat: " + e.message);
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>

<style scoped>
.profile-overlay { position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.7); backdrop-filter: blur(5px); z-index: 10001; display: flex; justify-content: center; align-items: center; }
.profile-card { width: 320px; background: #111b21; border-radius: 16px; padding: 2rem; display: flex; flex-direction: column; align-items: center; gap: 1.5rem; border: 1px solid #333; }
.card-header { width: 100%; display: flex; justify-content: flex-end; }
.btn-close { background: none; border: none; color: #888; font-size: 1.5rem; cursor: pointer; }
.avatar-wrapper { width: 140px; height: 140px; border-radius: 50%; overflow: hidden; border: 4px solid #202c33; margin-bottom: 10px; }
.avatar-img { width: 100%; height: 100%; object-fit: cover; }
.info-section { text-align: center; color: white; }
.user-name { font-size: 1.6rem; margin: 0; }
.user-bio { color: #888; margin-top: 5px; }
.btn-chat { width: 100%; padding: 12px; border-radius: 25px; background: #00a884; border: none; font-weight: bold; cursor: pointer; display: flex; align-items: center; justify-content: center; gap: 8px; font-size: 1rem; }
.btn-chat:hover { background: #008f6f; }
.btn-chat:disabled { background: #333; color: #666; cursor: wait; }
.spinner-mini { width: 16px; height: 16px; border: 2px solid #000; border-top-color: transparent; border-radius: 50%; animation: spin 1s linear infinite; }
@keyframes spin { to { transform: rotate(360deg); } }
</style>