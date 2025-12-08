<template>
  <div>
    <LoginForm v-if="!loggedIn" @loginSuccess="onLogin" />

    <Dashboard 
      v-else-if="!showProfile && currentUser.id" 
      :userId="userId"
      @openProfile="showProfile = true"
    />

    <Profile 
      v-else-if="showProfile && currentUser.id"
      :userId="userId"
      @back="showProfile = false"
    />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { currentUser } from './store.js'
import { getCurrentUser } from './services/api'

import LoginForm from './components/LoginForm.vue'
import Dashboard from './components/Dashboard.vue'
import Profile from './views/Profile.vue'

const loggedIn = ref(false)
const showProfile = ref(false)
const userId = ref('')

async function onLogin({ name, userId: id }) {
  loggedIn.value = true
  userId.value = id

  try {
    const user = await getCurrentUser(id)
    currentUser.value = user
  } catch (err) {
    console.error('Error cargando usuario al iniciar sesión:', err)
    loggedIn.value = false
    userId.value = ''
  }
}
</script>
