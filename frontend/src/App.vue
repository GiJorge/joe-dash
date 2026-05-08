<script setup>
import { ref, onMounted } from 'vue'

//const API_URL = 'http://localhost:8080/api/apps'
const API_URL = '/api/apps';
const apps = ref([])
const showSettings = ref(false)

// Form States
const newApp = ref({ name: '', url: '', icon: '' })
const editingId = ref(null)
const editForm = ref({ name: '', url: '', icon: '' })

// Helper to check if string is a URL for an image
const isImg = (str) => str && (str.startsWith('http') || str.startsWith('/img') || str.endsWith('.png') || str.endsWith('.svg'))

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

const deleteApp = async (id) => {
  if (!confirm("Remove this shortcut?")) return
  await fetch(`${API_URL}/${id}`, { method: 'DELETE' })
  apps.value = apps.value.filter(a => a.id !== id)
}

const startEdit = (app) => {
  editingId.value = app.id
  editForm.value = { ...app }
}

const saveEdit = async (id) => {
  const res = await fetch(`${API_URL}/${id}`, {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(editForm.value)
  })
  const updated = await res.json()
  const index = apps.value.findIndex(a => a.id === id)
  apps.value[index] = updated
  editingId.value = null
}

onMounted(fetchApps)
</script>

<template>
  <div class="min-h-screen bg-[#0f172a] text-slate-200 font-sans selection:bg-blue-500/30">
    
    <!-- Top Navigation -->
    <nav class="border-b border-slate-800 bg-slate-900/50 backdrop-blur-md sticky top-0 z-50">
      <div class="max-w-6xl mx-auto px-6 h-16 flex items-center justify-between">
        <h1 class="text-2xl font-bold text-blue-400 croissant-one-regular  tracking-wide">Joe Dash</h1>
        <button @click="showSettings = !showSettings" 
          class="flex items-center gap-2 px-5 py-2 rounded-full bg-slate-800 hover:bg-slate-700 border border-slate-700 transition-all text-sm font-medium">
          <span>{{ showSettings ? '🏠 Home' : '⚙️ Settings' }}</span>
        </button>
      </div>
    </nav>

    <main class="max-w-6xl mx-auto px-6 py-10">

      <!-- MAIN DASHBOARD -->
      <div v-if="!showSettings" class="animate-in">
        <div class="grid grid-cols-2 md:grid-cols-3 lg:grid-cols-5 gap-3">
          <a v-for="app in apps" :key="app.id" :href="app.url" target="_blank"
            class="group p-6 bg-slate-800/30 rounded-3xl border border-slate-700/50 hover:border-blue-500/50 hover:bg-slate-800/80 transition-all duration-300 hover:-translate-y-1">
            <div class="flex flex-col items-center text-center space-y-4">
              <div class="w-20 h-20 flex items-center justify-center bg-slate-900 rounded-2xl text-4xl shadow-inner group-hover:scale-105 transition-transform overflow-hidden p-2">
                <img v-if="isImg(app.icon)" :src="app.icon" class="w-full h-full object-contain" />
                <span v-else>{{ app.icon || '🌐' }}</span>
              </div>
              <div>
                <h3 class="font-bold text-lg text-slate-100 jj">{{ app.name }}</h3>
                <p class="text-[10px] text-slate-500 uppercase tracking-[0.2em] mt-1">Open Service</p>
              </div>
            </div>
          </a>

          <button @click="showSettings = true" class="p-6 border-2 border-dashed border-slate-800 rounded-3xl flex flex-col items-center justify-center text-slate-600 hover:text-blue-400 hover:border-blue-400/50 transition-colors">
            <span class="text-4xl mb-1">+</span>
            <span class="text-xs font-bold uppercase tracking-widest">New App</span>
          </button>
        </div>
      </div>
      <!-- END OF MAIN DASHBOARD -->


      <!-- SETTINGS / MANAGEMENT -->
      <div v-else class="max-w-2xl mx-auto animate-in">
        
        <div class="bg-slate-800/50 rounded-3xl p-8 border border-slate-700 shadow-2xl mb-8">
          <h2 class="text-xl font-bold mb-6 flex items-center gap-2 font-croissant">Add Application</h2>
          <div class="grid grid-cols-1 md:grid-cols-1 gap-4">
            <input v-model="newApp.name" placeholder="Name" class="input-style">
            <div class="flex gap-1">
              <input v-model="newApp.icon" placeholder="Emoji or Image URL" class="input-style flex-1">
              <div class="w-20 h-11  bg-slate-900 rounded-xl flex items-center justify-center border border-slate-700 overflow-hidden text-xl">
                 <img v-if="isImg(newApp.icon)" :src="newApp.icon" class="w-full h-full object-contain" />
                 <span v-else>{{ newApp.icon || '?' }}</span>
              </div>
            </div>
            <input v-model="newApp.url" placeholder="URL (http://...)" class="input-style md:col-span-2">
            <button @click="addApp" class="md:col-span-2 bg-blue-600 hover:bg-blue-500 py-3 rounded-xl font-bold transition-all shadow-lg shadow-blue-900/20">Save to Dashboard</button>
          </div>
        </div>

<div  class="space-y-2 max-h-[500px] overflow-y-auto pr-2 custom-scrollbar">


        <div class="space-y-3">
          <div v-for="app in apps" :key="app.id" class="bg-slate-900/50 border border-slate-800 rounded-2xl p-4">
            <div v-if="editingId === app.id" class="space-y-3">
              <div class="grid grid-cols-2 gap-2">
                <input v-model="editForm.name" class="input-style text-sm">
                <input v-model="editForm.icon" class="input-style text-sm">
              </div>
              <input v-model="editForm.url" class="input-style text-sm w-full">
              <div class="flex justify-end gap-3">
                <button @click="saveEdit(app.id)" class="text-xs font-bold text-green-400">UPDATE</button>
                <button @click="editingId = null" class="text-xs font-bold text-slate-500">CANCEL</button>
              </div>
            </div>
            <div v-else class="flex items-center justify-between">
              <div class="flex items-center gap-4">
                <div class="w-12 h-12 bg-slate-800 rounded-xl flex items-center justify-center overflow-hidden p-1">
                  <img v-if="isImg(app.icon)" :src="app.icon" class="w-full h-full object-contain" />
                  <span v-else class="text-xl">{{ app.icon || '🔗' }}</span>
                </div>
                <div>
                  <h4 class="font-bold text-slate-200">{{ app.name }}</h4>
                  <p class="text-[10px] text-slate-500 truncate w-40">{{ app.url }}</p>
                </div>
              </div>
              <div class="flex gap-4">
                <button @click="startEdit(app)" class="text-slate-500 hover:text-blue-400">Edit</button>
                <button @click="deleteApp(app.id)" class="text-slate-500 hover:text-red-500">Delete</button>
              </div>
            </div>
          </div>
        </div>
    </div>    
      </div>
      
      <!--END OF SETTINGS / MANAGEMENT -->
    </main>
  </div>
</template>

<!-- <style scoped>
@reference "./style.css";

.font-croissant { font-family: 'Croissant One', cursive; }

.input-style {
  @apply bg-slate-900 border border-slate-700 rounded-xl px-4 py-2 outline-none focus:border-blue-500 transition-all text-slate-200 placeholder:text-slate-600;
}

.animate-in {
  animation: fadeIn 0.4s ease-out;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(8px); }
  to { opacity: 1; transform: translateY(0); }
}
</style> -->