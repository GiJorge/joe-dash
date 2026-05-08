<script setup>
import { ref, onMounted } from 'vue';
import NavBar from './components/NavBar.vue';
import AppCard from './components/AppCard.vue';
import SettingsItem from './components/SettingsItem.vue'; // 1. Import it
import SettingsForm from './components/SettingsForm.vue'; // 2. Import the new SettingsForm component
import ThemeToggle from './components/ThemeToggle.vue';

const API_URL = '/api/apps';
const apps = ref([]);
const showSettings = ref(false);
const editingId = ref(null);
const editForm = ref({ name: '', url: '', icon: '' });
const newApp = ref({ name: '', url: '', icon: '' });

// API Actions
const fetchApps = async () => {
  const res = await fetch(API_URL)
  apps.value = await res.json()
}


const addApp = async () => {
  if (!newApp.value.name || !newApp.value.url) return
  const res = await fetch(API_URL, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(newApp.value)
  })
  const saved = await res.json()
  apps.value.unshift(saved)
  newApp.value = { name: '', url: '', icon: '' }
}



// const deleteApp = async (id) => {
//   await fetch(`/api/apps/${id}`, { method: 'DELETE' });
//   fetchApps();
// };

const deleteApp = async (id) => {
  if (!confirm("Remove this shortcut?")) return
  await fetch(`${API_URL}/${id}`, { method: 'DELETE' })
  apps.value = apps.value.filter(a => a.id !== id)
  fetchApps();
}


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

onMounted(fetchApps);
</script>

<template>
<div class="min-h-screen bg-slate-260 text-slate-900 dark:bg-[#0f172a] dark:text-slate-200 transition-colors duration-300">
    <NavBar :isSettings="showSettings" @toggle-settings="showSettings = !showSettings" >
      <ThemeToggle />
      </NavBar>
   
<main class="max-w-6xl mx-auto px-6 py-10">

  <!-- MAIN DASHBOARD -->
  <div v-if="!showSettings" class="animate-in">
    <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3">
      <AppCard v-for="app in apps" :key="app.id" v-bind="app" />
      
      <!-- Quick Add Button -->
      <button @click="showSettings = true" class="p-6 border-2 border-dashed border-slate-800 rounded-3xl flex flex-col items-center justify-center text-slate-600 hover:text-blue-400 hover:border-blue-400/50 transition-colors">
        <span class="text-4xl mb-1">+</span>
        <span class="text-xs font-bold uppercase tracking-widest">New App</span>
      </button>
    </div>
  </div>

  <!-- SETTINGS / MANAGEMENT -->
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