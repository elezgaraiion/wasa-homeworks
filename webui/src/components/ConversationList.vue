<template>
  <div class="sidebar">
    
    <header class="sidebar-header">
      <div class="user-profile-section" @click="$emit('openProfile')" title="Edit my profile">
        <img :src="getMyAvatarUrl()" @error="handleImageError" class="user-avatar" alt="Avatar" />
        <div class="my-info">
          <span class="my-name">{{ userName }}</span>
        </div>
      </div>
      
      <div class="header-actions">
        <span class="icon" @click="showGroupModal = true" title="New Group">➕👥</span>
        <span class="icon logout-btn" @click.stop="handleLogout" title="Log Out">🚪</span>
      </div>
    </header>

    <div class="search-bar">
      <div class="search-input-wrapper" @click="showSearchModal = true">
        <span class="search-icon">🔍</span>
        <span class="fake-input-text">Search users</span>
      </div>
    </div>

    <div class="chats-container">
      <div v-if="loading && conversations.length === 0" class="state-msg">Loading...</div>
      <div v-else-if="conversations.length === 0" class="state-msg"><p>You have no conversations.</p></div>
      <div v-else v-for="chat in conversations" :key="chat.id" class="chat-item" :class="{ active: selectedChatId === chat.id }" @click="selectChat(chat)">
        <img :src="getChatAvatarUrl(chat)" @error="handleImageError" class="chat-avatar" />
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

  </div>
</template>

<script>
import UserSearchModal from './UserSearchModal.vue';
import CreateGroupModal from './CreateGroupModal.vue';

export default {
  name: "ConversationList",
  components: {
    UserSearchModal,
    CreateGroupModal
  },
  props: ['currentUser'],
  data() {
    return {
      conversations: [],
      loading: true,
      selectedChatId: null,
      showSearchModal: false,
      showGroupModal: false,
      refreshInterval: null,
      DEFAULT_AVATAR: 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png'
    };
  },
  computed: {
    userName() {
      return this.currentUser?.name || this.currentUser?.Name || 'User';
    }
  },
  methods: {
    handleImageError(e) {
      e.target.src = this.DEFAULT_AVATAR;
    },
    resolveUrl(path) {
        if (!path) return '';
        if (path.startsWith('http')) return path;
        return this.$axios.defaults.baseURL + path;
    },
    getMyAvatarUrl() {
        const path = this.currentUser?.photo || this.currentUser?.Photo;
        if (!path) return this.DEFAULT_AVATAR;
        const fullUrl = this.resolveUrl(path);
        return fullUrl; 
    },
    getChatAvatarUrl(chat) {
        const path = chat.photo || chat.Photo;
        if (!path) return this.DEFAULT_AVATAR;
        return this.resolveUrl(path);
    },
    isMe(chat) {
      const myId = this.currentUser?.id || localStorage.getItem('userId');
      return chat.lastMessageSenderId === myId;
    },
    formatTime(dateStr) {
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
    },
    handleLogout() {
      if (confirm('Log out?')) {
          localStorage.removeItem('userId');
          this.$router.push('/login');
      }
    },
    async loadConversations() {
      try {
        const res = await this.$axios.get('/conversations');
        this.conversations = res.data || [];
      } catch (error) {
        console.error("Error chats:", error);
      } finally {
        this.loading = false;
      }
    },
    async onChatCreated(chat) {
      this.showSearchModal = false;
      this.showGroupModal = false;

      const index = this.conversations.findIndex(c => c.id === chat.id);
      if (index !== -1) {
        this.conversations.splice(index, 1);
        this.conversations.unshift(chat);
      } else {
        this.conversations.unshift(chat);
      }
      this.selectChat(chat);
    },
    selectChat(chat) {
      this.selectedChatId = chat.id;
      this.$emit('chatSelected', chat);
    }
  },
  mounted() {
    this.loadConversations();
    this.refreshInterval = setInterval(this.loadConversations, 4000);
  },
  beforeUnmount() {
    if (this.refreshInterval) clearInterval(this.refreshInterval);
  }
};
</script>

<style scoped>
.header-actions { display: flex; gap: 20px; align-items: center; } 
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