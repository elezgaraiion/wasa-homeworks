<template>
  <div id="app-root">
    <div v-if="isLoading" class="loading-screen">
      <div class="spinner"></div>
      <p>Conectando con WASA...</p>
    </div>

    <div v-else>
      <LoginView 
        v-if="!store.isAuthenticated" 
        @loginSuccess="onLoginSuccess" 
      />
      
      <Dashboard 
        v-else 
      />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import { store } from './store.js';
import { getCurrentUser } from './services/api';

import LoginView from './views/LoginView.vue';
import Dashboard from './views/Dashboard.vue';

const isLoading = ref(true);

onMounted(async () => {
  const token = localStorage.getItem('userId');
  
  if (token) {
    try {
      const user = await getCurrentUser();
      
      store.login(user, token);
      
    } catch (e) {
      console.error("Token inválido o error de red:", e);
      store.logout(); 
    }
  }
  
  isLoading.value = false;
});

async function onLoginSuccess() {
  isLoading.value = true;
  const token = localStorage.getItem('userId');
  try {
    const user = await getCurrentUser();
    store.login(user, token);
  } catch(e) {
    console.error(e);
  } finally {
    isLoading.value = false;
  }
}
</script>

<style>
body { margin: 0; background-color: #0b141a; font-family: sans-serif; color: white; }
.loading-screen { 
  height: 100vh; display: flex; flex-direction: column; 
  align-items: center; justify-content: center; color: #00a884; 
}
.spinner {
  border: 4px solid rgba(255,255,255,0.1); width: 36px; height: 36px; 
  border-radius: 50%; border-left-color: #00a884; animation: spin 1s linear infinite; margin-bottom: 20px;
}
@keyframes spin { 100% { transform: rotate(360deg); } }
</style>