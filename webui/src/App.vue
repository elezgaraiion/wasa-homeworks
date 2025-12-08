<template>
  <div id="app">
    <LoginView 
      v-if="!isAuthenticated" 
      @loginSuccess="checkAuth" 
    />

    <div v-else class="app-content">
      <h1>¡Dentro! Aquí cargaremos el chat...</h1>
      <button @click="logout">Salir</button>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue';
import LoginView from './views/LoginView.vue';
import { getCurrentUser } from './services/api';

const isAuthenticated = ref(false);

async function checkAuth() {
  const token = localStorage.getItem('userId');
  if (!token) {
    isAuthenticated.value = false;
    return;
  }

  try {
    // Validamos que el token sea real preguntando al backend
    await getCurrentUser();
    isAuthenticated.value = true;
  } catch (error) {
    console.error("Sesión inválida");
    localStorage.removeItem('userId');
    isAuthenticated.value = false;
  }
}

function logout() {
  localStorage.removeItem('userId');
  isAuthenticated.value = false;
}

// Al cargar la página, verificamos si ya había sesión
onMounted(() => {
  checkAuth();
});
</script>

<style>
body {
  margin: 0;
  background-color: #0b141a;
  font-family: sans-serif;
}
</style>