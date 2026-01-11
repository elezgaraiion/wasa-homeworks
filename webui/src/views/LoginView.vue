<template>
  <div class="login-container">
    
    <div class="intro-layer" :class="{ 'blurred': showInput }">
      <h1 class="giant-title">GUATXAP<br></h1>
      
      <transition name="fade">
        <button v-if="!showInput" @click="startInteraction" class="btn-start">
          START
        </button>
      </transition>
    </div>

    <div class="hidden-interface" :class="{ 'slide-down': showInput }">
      <div class="interface-content">
        <h2>Identification</h2>
        <p>Enter your codename to access</p>
        
        <input 
          v-model="username" 
          type="text" 
          placeholder="Username..." 
          @keyup.enter="handleLogin"
          :disabled="loading"
          ref="inputField"
        />

        <button @click="handleLogin" :disabled="loading || !username.trim()" class="btn-enter">
          {{ loading ? 'ACCESSING...' : 'ENTER' }}
        </button>
        
        <p v-if="error" class="error-msg">{{ error }}</p>
        
        <span class="cancel-link" @click="showInput = false">Cancel</span>
      </div>
    </div>

  </div>
</template>

<script>
export default {
  name: "LoginView",
  data() {
    return {
      showInput: false,
      username: "",
      loading: false,
      error: ""
    };
  },
  methods: {
    startInteraction() {
      this.showInput = true;
      
      this.$nextTick(() => {
        setTimeout(() => {
          if (this.$refs.inputField) {
            this.$refs.inputField.focus();
          }
        }, 500);
      });
    },

    async handleLogin() {
      this.error = "";
      if (this.username.trim().length < 3) {
        this.error = "Name is too short (min 3 letters).";
        return;
      }

      this.loading = true;

      try {
        const response = await this.$axios.post("/session", { 
            name: this.username.trim() 
        });
        
        const data = response.data;
        const userId = data.id || data.identifier;

        localStorage.setItem("userId", userId);

        this.$emit("loginSuccess");

        this.$router.push("/");

      } catch (e) {
        if (e.response && e.response.data) {
          this.error = e.response.data.error || e.message;
        } else {
          this.error = e.message;
        }
      } finally {
        this.loading = false;
      }
    }
  }
};
</script>

<style scoped>
.login-container {
  position: relative;
  width: 100vw;
  height: 100vh;
  background-color: #0b141a; 
  color: #00a884; 
  overflow: hidden;
  font-family: 'Helvetica Neue', Helvetica, Arial, sans-serif;
}

.intro-layer {
  height: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  transition: filter 0.5s ease;
}

.intro-layer.blurred {
  filter: blur(5px) brightness(0.5); 
}

.giant-title {
  font-size: 5rem;
  text-align: center;
  line-height: 1;
  letter-spacing: -2px;
  margin-bottom: 3rem;
  text-shadow: 0 0 20px rgba(0, 168, 132, 0.4);
  animation: pulse 3s infinite;
}

.btn-start {
  background: transparent;
  border: 2px solid #00a884;
  color: #00a884;
  padding: 1rem 3rem;
  font-size: 1.5rem;
  letter-spacing: 2px;
  cursor: pointer;
  transition: all 0.3s;
  text-transform: uppercase;
}

.btn-start:hover {
  background: #00a884;
  color: #0b141a;
  box-shadow: 0 0 30px #00a884;
}

.hidden-interface {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 50vh; 
  background: #111b21;
  border-bottom: 2px solid #00a884;
  box-shadow: 0 10px 50px rgba(0,0,0,0.8);
  
  transform: translateY(-100%); 
  transition: transform 0.6s cubic-bezier(0.22, 1, 0.36, 1); 
  
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 10;
}

.hidden-interface.slide-down {
  transform: translateY(0);
}

.interface-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 80%;
  max-width: 400px;
}

.interface-content h2 {
  color: #e9edef;
  margin-bottom: 0.5rem;
}

.interface-content p {
  color: #8696a0;
  margin-bottom: 2rem;
}

input {
  width: 100%;
  padding: 1rem;
  background: #202c33;
  border: none;
  border-bottom: 2px solid #005c4b;
  color: white;
  font-size: 1.2rem;
  text-align: center;
  margin-bottom: 1.5rem;
  outline: none;
  transition: border-color 0.3s;
}

input:focus {
  border-bottom-color: #00a884;
}

.btn-enter {
  width: 100%;
  padding: 1rem;
  background: #00a884;
  color: white;
  border: none;
  font-weight: bold;
  font-size: 1rem;
  cursor: pointer;
}

.btn-enter:disabled {
  background: #202c33;
  color: #8696a0;
  cursor: not-allowed;
}

.error-msg {
  color: #f15c6d;
  margin-top: 1rem;
}

.cancel-link {
  margin-top: 1.5rem;
  color: #8696a0;
  cursor: pointer;
  font-size: 0.9rem;
  text-decoration: underline;
}

@keyframes pulse {
  0%, 100% { opacity: 1; }
  50% { opacity: 0.7; }
}

.fade-enter-active, .fade-leave-active {
  transition: opacity 0.3s;
}
.fade-enter-from, .fade-leave-to {
  opacity: 0;
}
</style>