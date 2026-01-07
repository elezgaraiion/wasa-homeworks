<template>
  <div id="app-root">
    <div v-if="isLoading" class="loading-screen">
      <div class="spinner"></div>
      <p>Connecting to WASA...</p>
    </div>
    <router-view v-else />
  </div>
</template>

<script>
export default {
  name: 'App',
  data() {
    return {
      isLoading: true
    }
  },
  async mounted() {
    const token = localStorage.getItem('userId');
    
    if (token) {
      try {
        await this.$axios.get('/me');
      } catch (e) {
        localStorage.removeItem('userId');
        if (this.$route.path !== '/login') {
             this.$router.push('/login');
        }
      }
    } else {
      if (this.$route.path !== '/login') {
           this.$router.push('/login');
      }
    }
    this.isLoading = false;
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