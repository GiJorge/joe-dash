<script setup>
import { ref, onMounted, onUnmounted } from 'vue';
import NavBar from './components/NavBar.vue';
import AppCard from './components/AppCard.vue';
import SettingsItem from './components/SettingsItem.vue'; 
import SettingsForm from './components/SettingsForm.vue'; 
import ThemeToggle from './components/ThemeToggle.vue';

const API_URL = '/api/apps';
const apps = ref([]);
const showSettings = ref(false);
const editingId = ref(null);
const editForm = ref({ name: '', url: '', icon: '' });
const newApp = ref({ name: '', url: '', icon: '' });

// WebSocket Reference
let socket = null;

// Real-Time WebSocket Logic
const connectWebSocket = () => {
  // Automatically switch between ws:// and wss:// based on your environment
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
  
  // Connects to your Go port 10000 /ws route
  const wsUrl = `${protocol}//${window.location.hostname}:10000/ws`;

  socket = new WebSocket(wsUrl);

  socket.onmessage = (event) => {
    try {
      const updates = JSON.parse(event.data); // Expects: [{id: 1, status: 'online'}, ...]
      
      // Update only the status field inside your reactive apps list
      updates.forEach(update => {
        const app = apps.value.find(a => a.id === update.id);
        if (app) {
          app.status = update.status; // Updates the dot instantly!
        }
      });
    } catch (err) {
      console.error("Error parsing WebSocket update:", err);
    }
  };

  socket.onclose = () => {
    console.log("Dashboard status pipe disconnected. Retrying in 5 seconds...");
    setTimeout(connectWebSocket, 5000); // Clean auto-reconnect if server restarts
  };

  socket.onerror = (err) => {
    console.error("WebSocket Error:", err);
  };
};

// API Actions
const fetchApps = async () => {
  const res = await fetch(API_URL);
  apps.value = await res.json();
};

const addApp = async () => {
  if (!newApp.value.name || !newApp.value.url) return;
  const res = await fetch(API_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(newApp.value)
  });
  const saved = await res.json();
  apps.value.unshift(saved);
  newApp.value = { name: '', url: '', icon: '' };
};

const deleteApp = async (id) => {
  if (!confirm("Remove this shortcut?")) return;
  await fetch(`${API_URL}/${id}`, { method: 'DELETE' });
  apps.value = apps.value.filter(a => a.id !== id);
  fetchApps();
};

const startEdit = (app) => {
  editingId.value = app.id;
  editForm.value = { ...app };
};

const saveEdit = async (formData) => {
  await fetch(`/api/apps/${editingId.value}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(formData)
  });
  editingId.value = null;
  fetchApps();
};

onMounted(async () => {
  // 1. Grab initial application state from database
  await fetchApps();
  
  // 2. Open up the permanent real-time data pipe
  connectWebSocket();
});

// Lifecycle cleanup to close open sockets if the component unmounts
onUnmounted(() => {
  if (socket) socket.close();
});
</script>

<template>
  <div class="min-h-screen bg-slate-200 text-slate-900 dark:bg-[#0f172a] dark:text-slate-200 transition-colors duration-300">
    <NavBar :isSettings="showSettings" @toggle-settings="showSettings = !showSettings" @go-home="showSettings = false">
      <ThemeToggle />
    </NavBar>
   
    <main class="max-w-6xl mx-auto px-6 py-10">

      <div v-if="!showSettings" class="animate-in">
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3">
          <AppCard v-for="app in apps" :key="app.id" v-bind="app" />
          
          <button @click="showSettings = true" class="p-6 border-2 border-dashed border-slate-800 rounded-3xl flex flex-col items-center justify-center text-slate-600 hover:text-blue-400 hover:border-blue-400/50 transition-colors">
            <span class="text-4xl mb-1">+</span>
            <span class="text-xs font-bold uppercase tracking-widest">New App</span>
          </button>
        </div>
      </div>

      <SettingsForm 
        v-else
        :apps="apps"
        :newApp="newApp"
        :editingId="editingId"
        :editForm="editForm"
        @add="addApp"
        @edit="startEdit"
        @save="saveEdit"
        @cancel="editingId = null"
        @delete="deleteApp"
      />

    </main>
  </div>
</template>