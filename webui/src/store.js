// src/store.js
import { reactive } from 'vue';

// Usamos reactive para que sea un objeto profundo
export const store = reactive({
  currentUser: null,
  isAuthenticated: false,

  // Acción para iniciar sesión
  login(user, token) {
    this.currentUser = user;
    this.isAuthenticated = true;
    localStorage.setItem('userId', token); // Guardamos el token
  },

  // Acción para cerrar sesión
  logout() {
    this.currentUser = null;
    this.isAuthenticated = false;
    localStorage.removeItem('userId');
  },
  
  // Actualizar datos del usuario localmente
  updateUser(updates) {
    if (this.currentUser) {
      Object.assign(this.currentUser, updates);
    }
  }
});