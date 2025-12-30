<template>
  <div class="chat-window">
    
    <header class="chat-header" @click="showInfo = true">
      <div class="header-left">
        <img :src="chat.photo || 'https://upload.wikimedia.org/wikipedia/commons/7/7c/Profile_avatar_placeholder_large.png'" class="header-avatar" />
        <div class="header-info">
          <span class="header-name">{{ chat.name || chat.Name || 'Chat' }}</span>
          <span class="header-status">{{ chat.type === 'group' ? 'Group' : 'online' }}</span>
        </div>
      </div>
      <div class="header-actions"><span>🔍</span><span>⋮</span></div>
    </header>

    <div class="messages-area" ref="messagesContainer">
      <div v-if="loading && messages.length === 0" class="loading-msgs"><div class="spinner-small"></div></div>
      <div v-else-if="messages.length === 0" class="empty-chat"><span class="empty-badge">No messages.</span></div>

      <div 
        v-else 
        v-for="msg in messages" 
        :key="msg.id || msg.ID" 
        class="message-row" 
        :class="{ 
            'my-message': isMe(msg),
            'has-menu-open': activeMenuMessageId === (msg.id || msg.ID)
        }"
      >
        <div class="message-bubble">
          <div v-if="shouldShowSenderName(msg)" class="sender-name">
            {{ msg.sender?.name || msg.Sender?.Name || '...' }}
          </div>

          <div v-if="msg.replyTo || msg.ReplyTo" class="reply-quote-block">
            <span class="reply-sender-name">{{ msg.replyTo?.sender?.name || msg.ReplyTo?.Sender?.Name || 'User' }}</span>
            <div class="reply-preview-content">
                <span v-if="msg.replyTo?.photo || msg.ReplyTo?.Photo">📷 Photo</span>
                <span>{{ msg.replyTo?.text || msg.ReplyTo?.Text || '...' }}</span>
            </div>
          </div>

          <div class="message-content">
            <div v-if="msg.photo || msg.Photo" class="photo-wrapper">
              <img :src="msg.photo || msg.Photo" class="message-image" @load="scrollToBottom" />
              <p v-if="msg.text || msg.Text" class="caption-text">{{ msg.text || msg.Text }}</p>
            </div>
            <div v-else class="text-content">{{ msg.text || msg.Text }}</div>
          </div>

          <div class="message-meta">
            <span class="message-time">{{ formatTime(msg.createdAt || msg.CreatedAt) }}</span>
            <span v-if="isMe(msg)" class="check-icon">
               <span v-if="msg.status === 'read' || msg.Status === 'read'" style="color: #53bdeb;">✓✓</span>
               <span v-else>✓</span>
            </span>
          </div>

          <div v-if="msg.reactions && msg.reactions.length > 0" class="reactions-container">
             <button 
                v-for="group in getGroupedReactions(msg)" 
                :key="group.emoji" 
                class="reaction-pill"
                :class="{ 'reaction-mine': group.isMine }"
                @click="handleGroupClick(msg, group)"
             >
                <span class="emoji-char">{{ group.emoji }}</span>
                <span class="reaction-count">{{ group.count }}</span>
             </button>
          </div>

          <div class="msg-options-btn" @click.stop="toggleMenu(msg.id || msg.ID)">⌄</div>

          <div v-if="activeMenuMessageId === (msg.id || msg.ID)" class="msg-dropdown">
            <div class="dropdown-item" @click.stop="openReactionModal(msg)">😀 React</div>
            <div class="dropdown-item" @click="startReply(msg)">↩ Reply</div>
            <div class="dropdown-item" @click="openForwardModal(msg)">↪ Forward</div>
            <div v-if="isMe(msg)" class="dropdown-item delete-item" @click="handleDelete(msg)">🗑 Delete</div>
          </div>
        </div>
      </div>
    </div>

    <footer class="chat-footer">
      <div v-if="replyingToMessage" class="reply-preview-bar">
        <div class="reply-info">
            <span class="reply-target-name">Replying to {{ replyingToMessage.sender?.name || replyingToMessage.Sender?.Name }}</span>
            <p class="reply-target-text">{{ replyingToMessage.text || replyingToMessage.Text || (replyingToMessage.photo ? '📷 Photo' : '...') }}</p>
        </div>
        <button class="close-reply-btn" @click="cancelReply">✕</button>
      </div>
      <div class="input-container">
        <div class="input-wrapper">
            <input v-model="newMessage" @keyup.enter="handleSend" type="text" placeholder="Type a message..." :disabled="sending" ref="inputRef" />
        </div>
        <span class="footer-icon send-btn" @click="handleSend" :class="{ 'disabled': !newMessage.trim() }">➤</span>
      </div>
    </footer>

    <Teleport to="body">
      <ChatInfoModal 
          v-if="showInfo" 
          :chat="chat" 
          :chatId="chat.id" 
          @close="showInfo = false" 
          @chatUpdated="handleChatUpdated" 
      />
      
      <ForwardModal v-if="msgToForward" :message="msgToForward" :sourceChatId="chat.id" @close="msgToForward = null" @forwarded="handleForwardedSuccess" />
      
      <div v-if="msgToReact" class="modal-overlay" @click="closeReactionModal">
        <div class="reaction-modal" @click.stop>
            <h3>Choose a reaction</h3>
            <div class="reaction-grid">
                <button 
                    v-for="emoji in availableEmojis" 
                    :key="emoji" 
                    class="emoji-btn"
                    @click="handleReactFromModal(emoji)"
                >
                    {{ emoji }}
                </button>
            </div>
            <button class="close-modal-btn" @click="closeReactionModal">Cancel</button>
        </div>
      </div>
    </Teleport>

  </div>
</template>

<script setup>
import { ref, onMounted, onUnmounted, watch, nextTick } from 'vue';
import { getConversationMessages, sendMessage, markChatAsRead, deleteMessage, addReaction, removeReaction } from '../services/api';
import { store } from '../store.js';
import ChatInfoModal from './ChatInfoModal.vue';
import ForwardModal from './ForwardModal.vue'; 

const props = defineProps(['chat']); 
const emit = defineEmits(['refresh']); 

const messages = ref([]);
const loading = ref(false);
const sending = ref(false);
const newMessage = ref('');
const showInfo = ref(false);
const messagesContainer = ref(null);
const inputRef = ref(null);
let pollInterval = null;

const activeMenuMessageId = ref(null); 
const msgToForward = ref(null);
const replyingToMessage = ref(null);
const msgToReact = ref(null);
const availableEmojis = ['👍', '❤️', '😂', '😮', '😢', '🙏', '🔥', '🎉'];

function handleChatUpdated(updateInfo) {
    if (updateInfo) {
        if (updateInfo.type === 'name') props.chat.name = updateInfo.value;
        if (updateInfo.type === 'photo') props.chat.photo = updateInfo.value;
    }

    emit('refresh');
    
    loadMessages(true);
}

function startPolling() {
    if (pollInterval) clearInterval(pollInterval);
    pollInterval = setInterval(() => loadMessages(true), 3000);
}

function stopPolling() {
    if (pollInterval) {
        clearInterval(pollInterval);
        pollInterval = null;
    }
}

function isMe(msg) {
  const myId = store.currentUser?.id || localStorage.getItem('userId');
  const senderId = msg.sender?.id || msg.Sender?.ID || msg.sender_id || msg.SenderID;
  return senderId === myId;
}

function shouldShowSenderName(msg) {
  return props.chat.type === 'group' && !isMe(msg);
}

function formatTime(dateStr) {
  if (!dateStr) return '';
  const safeDate = dateStr.replace(' ', 'T');
  let date = new Date(safeDate.includes('Z') ? safeDate : safeDate + 'Z');
  if (isNaN(date.getTime())) date = new Date(safeDate); 
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
}

function getGroupedReactions(msg) {
    if (!msg.reactions || msg.reactions.length === 0) return [];
    const myId = store.currentUser?.id || localStorage.getItem('userId');
    const groups = {};

    msg.reactions.forEach(r => {
        const emoji = r.emoji;
        if (!groups[emoji]) {
            groups[emoji] = { emoji: emoji, count: 0, isMine: false, myReactionId: null };
        }
        groups[emoji].count++;
        const reactorId = r.user?.id || r.User?.ID;
        if (reactorId === myId) {
            groups[emoji].isMine = true;
            groups[emoji].myReactionId = r.id || r.ID;
        }
    });
    return Object.values(groups).sort((a, b) => b.count - a.count);
}

async function handleGroupClick(msg, group) {
    if (group.isMine) {
        try {
            await removeReaction(props.chat.id, msg.id || msg.ID, group.myReactionId);
            msg.reactions = msg.reactions.filter(r => (r.id || r.ID) !== group.myReactionId);
        } catch (e) { console.error(e); }
    } else {
        handleReact(msg, group.emoji);
    }
}

function openReactionModal(msg) {
    stopPolling();
    msgToReact.value = msg;
    activeMenuMessageId.value = null; 
}

function closeReactionModal() {
    msgToReact.value = null;
    startPolling();
}

async function handleReactFromModal(emoji) {
    if (msgToReact.value) {
        await addReactionInternal(msgToReact.value, emoji);
        closeReactionModal();
    }
}

async function addReactionInternal(msg, emoji) {
    const myId = store.currentUser?.id || localStorage.getItem('userId');
    try {
        const msgId = msg.id || msg.ID;
        if (!msg.reactions) msg.reactions = [];
        const tempReaction = {
            id: 'temp-' + Date.now(),
            emoji: emoji,
            user: { id: myId },
            User: { ID: myId }
        };
        const existingIdx = msg.reactions.findIndex(r => (r.user?.id || r.User?.ID) === myId);
        if (existingIdx !== -1) {
            msg.reactions[existingIdx] = tempReaction;
        } else {
            msg.reactions.push(tempReaction);
        }
        const newReaction = await addReaction(props.chat.id, msgId, emoji);
        const correctIdx = msg.reactions.findIndex(r => r.id === tempReaction.id);
        if (correctIdx !== -1) {
             msg.reactions[correctIdx] = newReaction;
        }
    } catch (e) {
        console.error(e);
        alert("Error reacting");
    }
}

function toggleMenu(msgId) {
    if (activeMenuMessageId.value === msgId) {
        activeMenuMessageId.value = null;
        startPolling(); 
    } else {
        stopPolling(); 
        activeMenuMessageId.value = msgId;
    }
}

window.addEventListener('click', () => { 
    if (activeMenuMessageId.value) {
        activeMenuMessageId.value = null; 
        startPolling(); 
    }
});

function openForwardModal(msg) { msgToForward.value = msg; activeMenuMessageId.value = null; }
function handleForwardedSuccess() {msgToForward.value = null; loadMessages(true);}
async function handleDelete(msg) {
    if (!confirm("Delete?")) return;
    activeMenuMessageId.value = null;
    try {
        await deleteMessage(props.chat.id, msg.id || msg.ID);
        messages.value = messages.value.filter(m => (m.id || m.ID) !== (msg.id || msg.ID));
    } catch (e) { alert("Error"); }
}
function startReply(msg) { replyingToMessage.value = msg; activeMenuMessageId.value = null; nextTick(() => inputRef.value?.focus()); }
function cancelReply() { replyingToMessage.value = null; }
async function handleSend() {
  const text = newMessage.value.trim();
  if (!text) return;
  newMessage.value = '';
  sending.value = true;
  const replyToId = replyingToMessage.value ? (replyingToMessage.value.id || replyingToMessage.value.ID) : null;
  replyingToMessage.value = null; 
  try {
    const sentMsg = await sendMessage(props.chat.id, text, replyToId);
    if (sentMsg) { messages.value.push(sentMsg); scrollToBottom(); } else { await loadMessages(true); }
    nextTick(() => inputRef.value?.focus());
  } catch (e) { newMessage.value = text; } finally { sending.value = false; }
}

async function loadMessages(isBackgroundUpdate = false) {
  if (!props.chat?.id) return;
  if (isBackgroundUpdate && msgToReact.value) return;

  markChatAsRead(props.chat.id).catch(e => {});
  if (!isBackgroundUpdate && messages.length === 0) loading.value = true;
  try {
    const res = await getConversationMessages(props.chat.id);
    const list = res || [];
    const orderedList = list.reverse();
    if (messages.value.length > 0) {
        orderedList.forEach(serverMsg => {
            const localMsg = messages.value.find(m => (m.id || m.ID) === (serverMsg.id || serverMsg.ID));
            if (localMsg && localMsg.reactions && localMsg.reactions.length > 0) {
                if (!serverMsg.reactions || serverMsg.reactions.length === 0) {
                      serverMsg.reactions = localMsg.reactions;
                }
            } else if (!serverMsg.reactions) {
                serverMsg.reactions = [];
            }
        });
    } else {
        orderedList.forEach(m => { if(!m.reactions) m.reactions = []; });
    }
    messages.value = orderedList;
    if (!isBackgroundUpdate) scrollToBottom();
  } catch (e) { } finally { loading.value = false; }
}

function scrollToBottom() { nextTick(() => { if (messagesContainer.value) messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight; }); }

onMounted(() => { 
    loadMessages(); 
    nextTick(() => inputRef.value?.focus()); 
    startPolling(); 
});

onUnmounted(() => { 
    stopPolling(); 
    window.removeEventListener('click', () => {}); 
});

watch(() => props.chat.id, () => { 
    messages.value = []; 
    replyingToMessage.value = null; 
    msgToReact.value = null; 
    loadMessages(); 
    nextTick(() => inputRef.value?.focus()); 
    startPolling(); 
});
</script>

<style scoped>
.chat-header { height: 60px; background-color: #202c33; padding: 0 16px; display: flex; align-items: center; border-left: 1px solid #333; flex-shrink: 0; cursor: pointer; transition: background 0.2s; }
.chat-window { display: flex; flex-direction: column; height: 100%; width: 100%; background-color: #0b141a; }
.header-left { display: flex; align-items: center; }
.header-avatar { width: 40px; height: 40px; border-radius: 50%; margin-right: 15px; object-fit: cover; }
.header-info { display: flex; flex-direction: column; }
.header-name { color: #e9edef; font-weight: bold; }
.header-status { color: #8696a0; font-size: 0.8rem; }
.header-actions { color: #aebac1; font-size: 1.2rem; gap: 20px; display: flex; margin-left: auto; }
.messages-area { flex: 1; overflow-y: auto; padding: 20px; display: flex; flex-direction: column; gap: 5px; background-image: url("https://user-images.githubusercontent.com/15075759/28719144-86dc0f70-73b1-11e7-911d-60d70fcded21.png"); background-repeat: repeat; background-blend-mode: overlay; background-color: #0b141a; }
.loading-msgs { display: flex; justify-content: center; padding: 20px; }
.spinner-small { width: 20px; height: 20px; border: 2px solid rgba(255,255,255,0.3); border-top-color: #00a884; border-radius: 50%; animation: spin 0.8s linear infinite; }
.empty-chat { display: flex; justify-content: center; margin-top: 20px; }
.empty-badge { background: #202c33; color: #ffd279; padding: 5px 12px; border-radius: 8px; font-size: 0.8rem; box-shadow: 0 1px 2px rgba(0,0,0,0.3); }
.message-row { display: flex; width: 100%; margin-bottom: 2px; position: relative; z-index: 1; }
.message-row.my-message { justify-content: flex-end; }
.message-row.has-menu-open { z-index: 100; }
.message-bubble { max-width: 65%; padding: 6px 9px; border-radius: 8px; background-color: #202c33; color: white; border-top-left-radius: 0; position: relative; box-shadow: 0 1px 0.5px rgba(0,0,0,0.13); overflow: visible; }
.my-message .message-bubble { background-color: #005c4b; border-top-left-radius: 8px; border-top-right-radius: 0; }
.sender-name { font-size: 0.8rem; color: #d65c3e; font-weight: 500; margin-bottom: 2px; line-height: 1.2; }
.message-content { font-size: 0.9rem; line-height: 1.3; }
.message-image { max-width: 100%; border-radius: 6px; margin-bottom: 4px; display: block; }
.caption-text { margin: 4px 0 0 0; }
.message-meta { float: right; margin-left: 8px; margin-top: 4px; display: flex; align-items: center; gap: 3px; position: relative; top: 4px; }
.message-time { font-size: 0.65rem; color: rgba(255,255,255,0.6); }
.check-icon { font-size: 0.75rem; color: #8696a0; }
.msg-options-btn { position: absolute; top: 0; right: 0; padding: 0 5px 0 8px; background: linear-gradient(to right, transparent, #202c33 50%); border-top-right-radius: 8px; color: #aebac1; cursor: pointer; opacity: 0; transition: opacity 0.2s; font-weight: bold; font-size: 1.2rem; line-height: 20px; z-index: 10; }
.my-message .msg-options-btn { background: linear-gradient(to right, transparent, #005c4b 50%); }
.message-bubble:hover .msg-options-btn { opacity: 1; }
.msg-dropdown + .msg-options-btn { opacity: 1; }
.msg-dropdown { position: absolute; top: 25px; background: #233138; border-radius: 6px; box-shadow: 0 4px 10px rgba(0,0,0,0.5); z-index: 200; min-width: 140px; padding: 5px 0; animation: fadeIn 0.1s; left: 0; right: auto; transform-origin: top left; }
.my-message .msg-dropdown { left: auto; right: 0; transform-origin: top right; }
.dropdown-item { padding: 10px 15px; color: #e9edef; font-size: 0.9rem; cursor: pointer; }
.dropdown-item:hover { background: #182229; }
.delete-item:hover { color: #f15c6d; }
.reply-quote-block { background-color: rgba(0, 0, 0, 0.15); border-left: 4px solid #00a884; border-radius: 4px; padding: 5px 8px; margin-bottom: 5px; cursor: pointer; font-size: 0.85rem; display: flex; flex-direction: column; }
.reply-sender-name { color: #00a884; font-weight: bold; font-size: 0.8rem; margin-bottom: 2px; }
.my-message .reply-sender-name { color: #d1d5db; }
.reply-preview-content { color: rgba(255, 255, 255, 0.7); white-space: nowrap; overflow: hidden; text-overflow: ellipsis; }
.chat-footer { display: flex; flex-direction: column; background-color: #202c33; padding: 0; min-height: 60px; }
.input-container { display: flex; align-items: center; padding: 10px 16px; width: 100%; }
.reply-preview-bar { background-color: #1d262d; padding: 8px 12px; display: flex; justify-content: space-between; align-items: center; border-left: 5px solid #00a884; border-top: 1px solid #333; animation: slideUp 0.2s; margin: 5px 10px 0 10px; border-radius: 8px; }
.reply-info { flex: 1; overflow: hidden; display: flex; flex-direction: column; }
.reply-target-name { color: #00a884; font-weight: bold; font-size: 0.85rem; margin-bottom: 2px; }
.reply-target-text { color: #aebac1; font-size: 0.85rem; white-space: nowrap; overflow: hidden; text-overflow: ellipsis; margin: 0; }
.close-reply-btn { background: none; border: none; color: #8696a0; font-size: 1.2rem; cursor: pointer; padding: 0 5px; }
.input-wrapper { flex: 1; background: #2a3942; border-radius: 8px; padding: 9px 12px; margin-right: 10px; }
.input-wrapper input { background: transparent; border: none; color: white; width: 100%; outline: none; font-size: 1rem; }
.send-btn { color: #8696a0; font-size: 1.5rem; cursor: pointer; transition: color 0.2s; }
.send-btn:hover { color: #00a884; }
.send-btn.disabled { color: #444; cursor: default; }

.reactions-container {
    display: flex;
    flex-wrap: wrap;
    gap: 5px;
    margin-top: 6px;
    margin-bottom: -4px;
    float: right;
    width: 100%;
    justify-content: flex-end;
}

.reaction-pill {
    background: #233138;
    border: 1px solid #3b4a54;
    border-radius: 16px;
    padding: 3px 8px;
    font-size: 0.85rem;
    color: #e9edef;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 5px;
    user-select: none;
    transition: all 0.2s;
    box-shadow: 0 1px 2px rgba(0,0,0,0.2);
}

.reaction-pill:hover {
    background: #374248;
    transform: translateY(-1px);
}

.reaction-pill.reaction-mine {
    background: #005c4b;
    border-color: #00a884;
    color: white;
}

.emoji-char {
    font-size: 1rem;
    line-height: 1;
}

.reaction-count {
    font-size: 0.75rem;
    font-weight: bold;
    opacity: 0.9;
}

.modal-overlay { position: fixed; top: 0; left: 0; right: 0; bottom: 0; background: rgba(0, 0, 0, 0.7); display: flex; justify-content: center; align-items: center; z-index: 9999; }
.reaction-modal { background: #233138; padding: 25px; border-radius: 16px; width: 320px; text-align: center; box-shadow: 0 10px 30px rgba(0,0,0,0.7); animation: popIn 0.2s; border: 1px solid #333; }
.reaction-modal h3 { color: #e9edef; margin-top: 0; margin-bottom: 20px; font-size: 1.2rem; font-weight: 500; }
.reaction-grid { display: grid; grid-template-columns: repeat(4, 1fr); gap: 15px; margin-bottom: 25px; }
.emoji-btn { background: #374248; border: none; border-radius: 10px; font-size: 2.2rem; padding: 10px 0; cursor: pointer; transition: transform 0.1s, background 0.2s; }
.emoji-btn:hover { background: #4a555b; transform: scale(1.15); }
.close-modal-btn { background: transparent; border: 1px solid #8696a0; color: #8696a0; padding: 10px 30px; border-radius: 20px; cursor: pointer; transition: all 0.2s; font-weight: bold; }
.close-modal-btn:hover { border-color: #f15c6d; color: #f15c6d; background: rgba(241, 92, 109, 0.1); }
@keyframes popIn { from { transform: scale(0.9); opacity: 0; } to { transform: scale(1); opacity: 1; } }
</style>