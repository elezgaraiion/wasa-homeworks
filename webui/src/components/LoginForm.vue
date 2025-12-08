<template>
  <div class="login">
    <h1>WASAText</h1>
    <h2>Ingresa tu nombre</h2>
    <input 
      v-model="name" 
      placeholder="Ej: María" 
      @keyup.enter="handleEnter" 
      :disabled="isSubmitting"
    />
    <button @click="handleEnter" :disabled="isSubmitting">
      {{ isSubmitting ? 'Entrando...' : 'Comenzar' }}
    </button>
    <p v-if="error" class="error">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue';
import { loginOrRegister, getCurrentUser } from '../services/api';

const name = ref('');
const error = ref(null);
const isSubmitting = ref(false);

const emit = defineEmits(['loginSuccess']);

async function handleEnter() {
  error.value = null;
  
  if (name.value.trim().length < 3) {
    error.value = 'El nombre debe tener al menos 3 caracteres';
    return;
  }

  isSubmitting.value = true;

  try {
    // 1. Login/Registro
    const sessionData = await loginOrRegister(name.value.trim());
    const token = sessionData.identifier;
    
    // Guardamos el token temporalmente en localStorage para que la siguiente petición funcione
    localStorage.setItem('userId', token);

    // 2. Obtenemos los datos completos del usuario
    const user = await getCurrentUser();

    emit('loginSuccess', { user, token });
    
  } catch (e) {
    error.value = e.message;
    localStorage.removeItem('userId'); // Limpieza si falla
  } finally {
    isSubmitting.value = false;
  }
}
</script>

<style scoped>
.login { display:flex; flex-direction:column; align-items:center; justify-content:center; height:100vh; padding: 2rem; }
input { margin-top:0.5rem; padding:0.5rem; font-size:1rem; border: 1px solid #ccc; border-radius: 4px; }
button { margin-top:1rem; padding:0.5rem 1rem; cursor: pointer; background: #008069; color: white; border: none; border-radius: 4px; }
button:disabled { background: #ccc; }
.error { color: red; margin-top: 1rem; }
</style>