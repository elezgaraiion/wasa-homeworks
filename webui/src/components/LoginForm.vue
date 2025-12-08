<template>
  <div class="login">
    <h2>Ingresa tu nombre</h2>
    <input v-model="name" placeholder="Tu nombre" />
    <button @click="handleEnter">Comenzar</button>
    <p v-if="error" style="color:red">{{ error }}</p>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { loginOrRegister } from '../services/api'

const name = ref('')
const error = ref(null)
const emit = defineEmits(['loginSuccess'])

async function handleEnter() {
  if (name.value.trim().length < 3) {
    error.value = 'El nombre debe tener al menos 3 caracteres'
    return
  }

  try {
    const data = await loginOrRegister(name.value.trim())
    emit('loginSuccess', { name: name.value.trim(), userId: data.identifier })
  } catch (e) {
    error.value = e.message
  }
}
</script>

<style>
.login { display:flex; flex-direction:column; align-items:center; justify-content:center; height:100vh; }
input { margin-top:0.5rem; padding:0.5rem; font-size:1rem; }
button { margin-top:1rem; padding:0.5rem 1rem; }
</style>
