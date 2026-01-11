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

<script>
export default {
  name: "ForwardModal",
  props: ['message', 'sourceChatId'],
  emits: ['close', 'forwarded'],
  data() {
    return {
      conversations: [],
      loading: true,
      DEFAULT_AVATAR: 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png'
    };
  },
  async mounted() {
    try {
      const res = await this.$axios.get('/conversations');
      this.conversations = res.data || [];
    } catch (e) {
      console.error(e);
    } finally {
      this.loading = false;
    }
  },
  methods: {
    async handleForward(targetChat) {      
      try {
        const msgId = this.message.id || this.message.ID;
        
        // We use the same endpoint logic we defined previously
        await this.$axios.post(`/conversations/${this.sourceChatId}/messages/${msgId}`, {
            targetConversationId: targetChat.id
        });
        
        this.$emit('forwarded'); 
        this.$emit('close');
      } catch (e) {
        const errorMsg = e.response?.data?.error || e.message;
        alert("Error forwarding: " + errorMsg);
      }
    }
  }
};
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