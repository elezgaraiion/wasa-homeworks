<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="group-card">
      
      <div class="card-header">
        <h2>New Group</h2>
        <button class="btn-close" @click="$emit('close')">✕</button>
      </div>

      <div class="input-section">
        <input 
          v-model="groupName" 
          type="text" 
          placeholder="Group name..." 
          class="main-input"
          maxlength="30"
        />
      </div>

      <div class="selected-area" v-if="selectedUsers.length > 0">
        <div v-for="user in selectedUsers" :key="user.id" class="user-chip">
          <img :src="user.photo || DEFAULT_AVATAR" class="chip-avatar">
          <span>{{ user.name }}</span>
          <span class="remove-btn" @click="removeUser(user.id)">✕</span>
        </div>
      </div>

      <div class="search-section">
        <span class="search-icon">🔍</span>
        <input 
          v-model="query" 
          type="text" 
          placeholder="Add participants..." 
          class="search-input"
        />
      </div>

      <div class="results-list">
        <div v-if="loading" class="status-msg">Searching...</div>
        
        <div 
          v-else 
          v-for="user in searchResults" 
          :key="user.id" 
          class="user-row"
          @click="selectUser(user)"
          :class="{ 'disabled': isSelected(user.id) }"
        >
          <img :src="user.photo || DEFAULT_AVATAR" class="user-avatar" />
          <span class="user-name">{{ user.name }}</span>
          <span v-if="isSelected(user.id)" class="check-mark">✓</span>
        </div>
      </div>

      <div class="footer-actions">
        <button 
          class="btn-create" 
          @click="handleCreateGroup" 
          :disabled="creating || !groupName.trim() || selectedUsers.length === 0"
        >
          {{ creating ? 'Creating...' : 'Create Group' }}
        </button>
      </div>

    </div>
  </div>
</template>

<script>
export default {
  name: "CreateGroupModal",
  emits: ['close', 'groupCreated'],
  data() {
    return {
      groupName: '',
      query: '',
      searchResults: [],
      selectedUsers: [],
      loading: false,
      creating: false,
      DEFAULT_AVATAR: 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png'
    };
  },
  watch: {
    async query(newVal) {
      if (!newVal.trim()) {
        this.searchResults = [];
        return;
      }
      this.loading = true;
      try {
        const response = await this.$axios.get('/users', { params: { q: newVal } });
        this.searchResults = response.data || [];
      } catch (e) {
        console.error(e);
      } finally {
        this.loading = false;
      }
    }
  },
  methods: {
    isSelected(userId) {
      return this.selectedUsers.some(u => (u.id || u.ID) === userId);
    },

    selectUser(user) {
      const uid = user.id || user.ID;
      if (this.isSelected(uid)) return; 
      this.selectedUsers.push(user);
      this.query = ''; 
      this.searchResults = [];
    },

    removeUser(userId) {
      this.selectedUsers = this.selectedUsers.filter(u => (u.id || u.ID) !== userId);
    },

    async handleCreateGroup() {
      if (!this.groupName.trim() || this.selectedUsers.length === 0) return;
      
      this.creating = true;
      try {
        const userIds = this.selectedUsers.map(u => u.id || u.ID);
        
        const response = await this.$axios.post('/conversations', { 
            name: this.groupName, 
            users: userIds 
        });
        
        const newGroupChat = response.data;
        
        this.$emit('groupCreated', newGroupChat);
        this.$emit('close');

      } catch (e) {
        alert('Error creating group: ' + e.message);
      } finally {
        this.creating = false;
      }
    }
  }
};
</script>

<style scoped>
.modal-backdrop { position: fixed; top: 0; left: 0; width: 100vw; height: 100vh; background: rgba(0,0,0,0.7); z-index: 9999; display: flex; justify-content: center; align-items: center; animation: fadeIn 0.2s; }
.group-card { width: 450px; max-width: 95%; background: #111b21; border: 1px solid #333; border-radius: 12px; display: flex; flex-direction: column; overflow: hidden; height: 80vh; animation: slideUp 0.3s; }

.card-header { padding: 15px 20px; background: #202c33; display: flex; justify-content: space-between; align-items: center; }
.card-header h2 { margin: 0; font-size: 1.1rem; color: #e9edef; }
.btn-close { background: none; border: none; color: #aebac1; font-size: 1.5rem; cursor: pointer; }

.input-section { padding: 20px; border-bottom: 1px solid #2f3b43; }
.main-input { width: 100%; background: transparent; border: none; border-bottom: 2px solid #00a884; color: white; font-size: 1.1rem; padding: 5px 0; outline: none; }

.selected-area { padding: 10px 20px; display: flex; flex-wrap: wrap; gap: 8px; max-height: 100px; overflow-y: auto; background: #0b141a; }
.user-chip { background: #202c33; color: #e9edef; padding: 4px 8px; border-radius: 16px; display: flex; align-items: center; gap: 6px; font-size: 0.9rem; border: 1px solid #333; }
.chip-avatar { width: 20px; height: 20px; border-radius: 50%; }
.remove-btn { cursor: pointer; color: #f15c6d; font-weight: bold; margin-left: 4px; }

.search-section { padding: 10px 20px; display: flex; align-items: center; background: #111b21; border-bottom: 1px solid #2f3b43; }
.search-icon { color: #8696a0; margin-right: 10px; }
.search-input { width: 100%; background: transparent; border: none; color: white; outline: none; }

.results-list { flex: 1; overflow-y: auto; }
.status-msg { padding: 20px; text-align: center; color: #8696a0; }
.user-row { display: flex; align-items: center; padding: 10px 20px; cursor: pointer; border-bottom: 1px solid #222; }
.user-row:hover { background: #202c33; }
.user-row.disabled { opacity: 0.5; cursor: default; }
.user-avatar { width: 40px; height: 40px; border-radius: 50%; margin-right: 15px; }
.user-name { color: #e9edef; flex: 1; }
.check-mark { color: #00a884; font-weight: bold; }

.footer-actions { padding: 15px; background: #202c33; display: flex; justify-content: flex-end; }
.btn-create { background: #00a884; color: #111; border: none; padding: 10px 24px; border-radius: 20px; font-weight: bold; cursor: pointer; }
.btn-create:disabled { background: #333; color: #666; cursor: not-allowed; }

@keyframes fadeIn { from { opacity: 0; } to { opacity: 1; } }
@keyframes slideUp { from { transform: translateY(20px); opacity: 0; } to { transform: translateY(0); opacity: 1; } }
</style>