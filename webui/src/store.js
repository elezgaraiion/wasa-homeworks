import { reactive } from 'vue';

export const store = reactive({
  currentUser: null,
  isAuthenticated: false,

  login(user, token) {
    this.currentUser = user; 
    this.isAuthenticated = true;
    localStorage.setItem('userId', token);
  },

  logout() {
    this.currentUser = null;
    this.isAuthenticated = false;
    localStorage.removeItem('userId');
    window.location.reload(); 
  },
  
  updateUser(updates) {
    if (this.currentUser) {
      Object.assign(this.currentUser, updates);
    }
  }
});