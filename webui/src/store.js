import { reactive } from 'vue';

export const store = reactive({
  currentUser: null,
  isAuthenticated: false,

  // Guardar usuario y token
  login(user, token) {
    this.currentUser = user; // Aquí se guarda el {id, name, photo} que viene del backend
    this.isAuthenticated = true;
    localStorage.setItem('userId', token);
  },

  // Cerrar sesión
  logout() {
    this.currentUser = null;
    this.isAuthenticated = false;
    localStorage.removeItem('userId');
    // Forzamos recarga para limpiar cualquier estado residual en memoria
    window.location.reload(); 
  },
  
  // Actualizar datos parciales (ej: si cambias nombre o foto)
  updateUser(updates) {
    if (this.currentUser) {
      Object.assign(this.currentUser, updates);
    }
  }
});