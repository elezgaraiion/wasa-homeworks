<template>
  <div class="dashboard">
    <header class="dashboard-header">
      <h2 v-if="currentUser.name">Bienvenido, {{ currentUser.name }}</h2>
      <h2 v-else>Cargando usuario...</h2>
      <button @click="$emit('openProfile')">Mi Perfil</button>
    </header>

    <h3>Conversaciones</h3>
    <div v-if="loading">Cargando conversaciones...</div>
    <ul v-else>
      <li v-for="conv in conversations" :key="conv.ID" class="conversation">
        <img v-if="conv.Photo" :src="conv.Photo" class="avatar" />
        <div class="info">
          <strong>{{ conv.Name || 'Sin nombre' }}</strong>
          <p>{{ conv.LastMessagePreview || 'Sin mensajes aún' }}</p>
        </div>
        <span class="time">{{ formatDate(conv.LastMessageAt) }}</span>
      </li>
    </ul>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getConversations } from '../services/api'
import { currentUser } from '@/store.js'

const props = defineProps({ userId: String })
const conversations = ref([])
const loading = ref(true)

function formatDate(dateStr) {
  if (!dateStr) return ''
  return new Date(dateStr).toLocaleTimeString([], { hour:'2-digit', minute:'2-digit' })
}

async function loadConversations() {
  loading.value = true
  try {
    conversations.value = await getConversations(props.userId)
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

onMounted(loadConversations)
</script>

<style>
.dashboard { display:flex; flex-direction:column; padding:1rem; }
.dashboard-header { display:flex; justify-content:space-between; align-items:center; }
ul { list-style:none; padding:0; margin:0; }
.conversation { display:flex; align-items:center; padding:0.5rem 0; border-bottom:1px solid #eee; }
.avatar { width:40px; height:40px; border-radius:50%; margin-right:0.5rem; }
.avatar-large { width:80px; height:80px; border-radius:50%; margin-top:0.5rem; }
.info { flex:1; }
.time { font-size:0.8rem; color:#888; }
</style>
