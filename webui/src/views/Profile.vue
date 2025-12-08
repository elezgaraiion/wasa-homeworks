<template>
  <div class="profile">
    <h2>Mi Perfil</h2>

    <div class="avatar-container">
      <img v-if="currentUser.photo" :src="currentUser.photo" class="avatar-large" />
    </div>

    <p><strong>ID:</strong> {{ currentUser.id }}</p>

    <div class="update-name">
      <label for="name">Nombre:</label>
      <input id="name" v-model="newName" placeholder="Introduce nuevo nombre" />
      <button @click="updateName">Actualizar</button>
      <p v-if="error" style="color:red">{{ error }}</p>
      <p v-if="success" style="color:green">{{ success }}</p>
    </div>

    <button @click="$emit('back')">Volver</button>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { currentUser } from '@/store.js'
import { getCurrentUser, updateUserName } from '../services/api'

const props = defineProps({ userId: String })
const emit = defineEmits(['back'])

const newName = ref('')
const error = ref('')
const success = ref('')

onMounted(async () => {
  try {
    const user = await getCurrentUser(props.userId)
    currentUser.value = user
    newName.value = user.name
  } catch (err) {
    console.error(err)
  }
})

async function updateName() {
  error.value = ''
  success.value = ''

  if (!newName.value.trim()) {
    error.value = 'El nombre no puede estar vacío'
    return
  }

  try {
    const updatedUser = await updateUserName(props.userId, newName.value.trim())
    currentUser.value.name = updatedUser.name
    success.value = 'Nombre actualizado correctamente'
  } catch (err) {
    error.value = err.message
  }
}
</script>

<style>
.profile { display:flex; flex-direction:column; align-items:center; padding:1rem; }
.avatar-large { width:80px; height:80px; border-radius:50%; margin-bottom:1rem; }
.update-name { margin-top:1rem; display:flex; flex-direction:column; align-items:center; }
.update-name input { padding:0.5rem; margin-top:0.5rem; width:200px; }
.update-name button { margin-top:0.5rem; padding:0.5rem 1rem; }
button { margin-top:1rem; padding:0.5rem 1rem; }
</style>
