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
      // CRÍTICO: Pedimos al backend quiénes somos usando el token guardado
      const user = await getCurrentUser();
      
      // Si el backend responde, guardamos el usuario COMPLETO (nombre, foto, id)
      store.login(user, token);
      
    } catch (e) {
      console.error("Token inválido o error de red:", e);
      // Si falla, borramos token y mandamos al login
      store.logout(); 
    }
  }
  
  // Quitamos la pantalla de carga pase lo que pase
  isLoading.value = false;
});

async function onLoginSuccess() {
  // Cuando el LoginView nos avisa que entró, volvemos a pedir los datos frescos
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