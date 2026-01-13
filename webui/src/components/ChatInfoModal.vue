<template>
  <div class="modal-overlay" @click.self="$emit('close')">
    <div class="modal-content">
      <button class="close-btn" @click="$emit('close')">✕</button>
      
      <div class="group-header">
        <div class="photo-container">
            <img 
                :src="resolveUrl(localChat.photo) || 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png'" 
                class="group-avatar" 
                @error="handleImgError"
            />
            <label v-if="isGroup" class="edit-photo-btn">
                📷
                <input type="file" @change="handlePhotoChange" accept="image/*" hidden />
            </label>
        </div>

        <div class="name-container">
            
            <div v-if="isEditingName" class="edit-mode-wrapper">
                <input 
                    v-model="newName" 
                    @keyup.enter="saveName"
                    ref="nameInput"
                    class="name-input"
                    placeholder="New name..."
                />
                <div class="edit-actions">
                    <button @click="saveName" class="btn-confirm">Confirm</button>
                    <button @click="isEditingName = false" class="btn-cancel">Cancel</button>
                </div>
            </div>

            <div v-else class="view-mode-wrapper">
                <h2>{{ localChat.name || localChat.Name || 'Chat' }}</h2>
                <button v-if="isGroup" @click="startEditingName" class="icon-btn">✎</button>
            </div>

        </div>
        
        <p class="participants-count" v-if="isGroup">
            {{ participantsList.length }} participants
        </p>
      </div>

      <hr class="divider" />

      <div v-if="isGroup" class="add-section">
        <h3>Add participants</h3>
        <input 
            v-model="searchQuery" 
            @input="handleSearch" 
            type="text" 
            placeholder="Search users by name..." 
            class="search-input"
        />

        <div class="user-list" v-if="searchResults.length > 0">
            <div v-for="user in searchResults" :key="user.id" class="user-row">
                <img :src="resolveUrl(user.photo) || 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png'" class="user-avatar-small" />
                <span class="user-name">{{ user.name }}</span>
                
                <button v-if="isAlreadyMember(user.id)" class="btn-member" disabled>Already a member</button>
                <button v-else class="btn-add" @click="addToGroup(user)">Add</button>
            </div>
        </div>
        <div v-else-if="searchQuery.length > 2" class="no-results">
            No users found.
        </div>
      </div>

      <div v-if="isGroup" class="danger-zone">
        <button class="btn-leave" @click="handleLeaveGroup">Leave group</button>
      </div>

    </div>
  </div>
</template>

<script>
export default {
  name: "ChatInfoModal",
  props: ['chatId', 'chat'],
  emits: ['close', 'chatUpdated'],
  data() {
    return {
      localChat: {},
      participantsList: [],
      isGroup: false,
      isEditingName: false,
      newName: '',
      searchQuery: '',
      searchResults: []
    };
  },
  watch: {
    chat: {
      immediate: true,
      handler(newChat) {
        if (newChat) {
          this.localChat = { ...newChat };
          this.isGroup = newChat.type === 'group';
          this.participantsList = [...(newChat.participants || [])];
        }
      }
    }
  },
  methods: {
    resolveUrl(path) {
        if (!path) return '';
        if (path.startsWith('http')) return path;
        return this.$axios.defaults.baseURL + path;
    },
    handleImgError(e) {
      e.target.src = 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png';
    },

    isAlreadyMember(userId) {
      return this.participantsList.some(p => {
        const pId = p.id || p.user_id || p; 
        return pId === userId;
      });
    },

    async handleSearch() {
      if (this.searchQuery.length < 2) {
        this.searchResults = [];
        return;
      }
      try {
        const res = await this.$axios.get('/users', { params: { q: this.searchQuery } });
        this.searchResults = res.data || [];
      } catch (e) { console.error(e); }
    },

    async addToGroup(user) {
      try {
        await this.$axios.post(`/conversations/${this.chatId}/users`, { 
            userId: user.id 
        });
        
        this.participantsList.push(user);
        this.searchQuery = '';
        this.searchResults = [];
        this.$emit('chatUpdated');
        alert(`${user.name} added.`);
      } catch (e) { 
        console.error(e);
        alert("Error adding user: " + (e.response?.data?.error || e.message)); 
      }
    },

    startEditingName() {
      this.newName = this.localChat.name || '';
      this.isEditingName = true;
      this.$nextTick(() => {
        if(this.$refs.nameInput) this.$refs.nameInput.focus();
      });
    },

    async saveName() {
      if (!this.newName.trim()) return;
      try {
        const response = await this.$axios.put(`/conversations/${this.chatId}/name`, { name: this.newName });
        const updated = response.data;
        this.localChat.name = updated.name;
        this.isEditingName = false;
        
        this.$emit('chatUpdated', { type: 'name', value: updated.name });
        
      } catch (e) {
        alert("Error changing name.");
      }
    },

    async handlePhotoChange(event) {
      const file = event.target.files[0];
      if (!file) return;
      try {
        const formData = new FormData();
        formData.append('photoFile', file);
        const response = await this.$axios.put(`/conversations/${this.chatId}/photo`, formData);
        const updated = response.data;
        const newPhotoUrl = updated.photo + '?t=' + new Date().getTime();
        this.localChat.photo = newPhotoUrl;
        
        this.$emit('chatUpdated', { type: 'photo', value: newPhotoUrl });
      } catch (e) { alert("Error uploading photo."); }
    },

    async handleLeaveGroup() {
      try {
        await this.$axios.delete(`/conversations/${this.chatId}/users/me`);
        this.$emit('chatUpdated');
        this.$emit('close');
        window.location.reload(); 
      } catch (e) { alert("Error leaving."); }
    }
  }
};
</script>

<style scoped>
.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0,0,0,0.7); display: flex; justify-content: center; align-items: center; z-index: 2000; }
.modal-content { background: #202c33; width: 400px; max-height: 90vh; overflow-y: auto; border-radius: 12px; padding: 20px; position: relative; color: #e9edef; box-shadow: 0 4px 20px rgba(0,0,0,0.5); border: 1px solid #374045; }
.close-btn { position: absolute; top: 10px; right: 15px; background: none; border: none; color: #aebac1; font-size: 1.5rem; cursor: pointer; }

.group-header { display: flex; flex-direction: column; align-items: center; gap: 15px; margin-bottom: 20px; }
.photo-container { position: relative; width: 120px; height: 120px; }
.group-avatar { width: 100%; height: 100%; border-radius: 50%; object-fit: cover; border: 3px solid #00a884; }
.edit-photo-btn { position: absolute; bottom: 5px; right: 5px; background: #00a884; width: 35px; height: 35px; border-radius: 50%; display: flex; align-items: center; justify-content: center; cursor: pointer; box-shadow: 0 2px 5px rgba(0,0,0,0.4); font-size: 1.1rem; transition: transform 0.2s; }
.edit-photo-btn:hover { transform: scale(1.1); }

.name-container { width: 100%; display: flex; justify-content: center; min-height: 40px; }

.view-mode-wrapper { display: flex; align-items: center; gap: 10px; }
.view-mode-wrapper h2 { margin: 0; font-size: 1.5rem; text-align: center; }

.edit-mode-wrapper { display: flex; flex-direction: column; align-items: center; gap: 10px; width: 100%; }
.name-input { background: #2a3942; border: 1px solid #00a884; border-radius: 6px; color: white; font-size: 1.2rem; text-align: center; width: 80%; padding: 8px; outline: none; }

.edit-actions { display: flex; gap: 10px; margin-top: 5px; }
.btn-confirm { background: #00a884; color: white; border: none; padding: 6px 15px; border-radius: 6px; cursor: pointer; font-weight: bold; }
.btn-confirm:hover { background: #008f6f; }
.btn-cancel { background: transparent; border: 1px solid #8696a0; color: #8696a0; padding: 6px 15px; border-radius: 6px; cursor: pointer; }
.btn-cancel:hover { border-color: #f15c6d; color: #f15c6d; }

.icon-btn { background: none; border: none; cursor: pointer; font-size: 1.2rem; color: #aebac1; padding: 5px; }
.icon-btn:hover { color: #00a884; }
.participants-count { color: #8696a0; font-size: 0.9rem; margin-top: -5px; }

.divider { border: 0; border-top: 1px solid #374045; margin: 20px 0; }

.add-section h3 { margin-bottom: 15px; font-size: 1rem; color: #00a884; font-weight: bold; text-transform: uppercase; letter-spacing: 0.5px; }
.search-input { width: 100%; padding: 12px; border-radius: 8px; border: none; background: #2a3942; color: white; margin-bottom: 15px; outline: none; font-size: 0.95rem; }
.search-input:focus { box-shadow: 0 0 0 2px #00a884; }
.no-results { color: #8696a0; text-align: center; font-style: italic; }

.user-list { display: flex; flex-direction: column; gap: 8px; max-height: 250px; overflow-y: auto; padding-right: 5px; }
.user-row { display: flex; align-items: center; background: #111b21; padding: 10px; border-radius: 8px; gap: 12px; transition: background 0.2s; }
.user-row:hover { background: #18242b; }
.user-avatar-small { width: 40px; height: 40px; border-radius: 50%; object-fit: cover; }
.user-name { flex: 1; font-weight: 500; font-size: 0.95rem; }

.btn-add { background: #00a884; color: white; border: none; padding: 6px 14px; border-radius: 20px; cursor: pointer; font-weight: bold; font-size: 0.8rem; transition: background 0.2s; }
.btn-add:hover { background: #008f6f; }
.btn-member { background: transparent; color: #8696a0; border: 1px solid #374045; padding: 6px 14px; border-radius: 20px; font-size: 0.8rem; cursor: default; user-select: none; }

.danger-zone { margin-top: 40px; border-top: 1px solid #374045; padding-top: 20px; display: flex; justify-content: center; }
.btn-leave { background: transparent; color: #f15c6d; border: 1px solid #f15c6d; padding: 10px 20px; border-radius: 8px; cursor: pointer; width: 100%; font-weight: bold; transition: background 0.2s; }
.btn-leave:hover { background: rgba(241, 92, 109, 0.1); }

.user-list::-webkit-scrollbar { width: 6px; }
.user-list::-webkit-scrollbar-thumb { background: #374045; border-radius: 3px; }
.user-list::-webkit-scrollbar-track { background: transparent; }
</style>