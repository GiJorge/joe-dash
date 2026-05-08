<template>
  <div class="max-w-2xl mx-auto animate-in">
    
    <!-- ADD NEW APP BOX -->
    <div class="bg-slate-300/50 dark:bg-slate-900/50  rounded-3xl p-8 border border-slate-700 shadow-2xl mb-8">
      <h2 class="text-xl font-bold mb-6 flex items-center gap-2 font-croissant">Add Application</h2>
      <div class="grid grid-cols-1 gap-4 ">
        <input v-model="newApp.name" placeholder="Name" class="input-style">
        
        <div class="flex gap-1">
          <input v-model="newApp.icon" placeholder="Emoji or Image URL" class="input-style flex-1">
          <div class="w-20 h-11   rounded-xl flex items-center justify-center border border-slate-700 overflow-hidden text-xl">
             <img v-if="isImg(newApp.icon)" :src="newApp.icon" class="w-full h-full object-contain" />
             <span v-else>{{ newApp.icon || '?' }}</span>
          </div>
        </div>

        <input v-model="newApp.url" placeholder="URL (http://...)" class="input-style ">
        
        <button @click="$emit('add', newApp)" 
          class="bg-blue-300 dark:bg-blue-600 hover:bg-blue-400 py-3 rounded-xl font-bold transition-all shadow-lg shadow-blue-900/20">
          Save to Dashboard
        </button>
      </div>
    </div>

    <!-- SCROLLABLE LIST -->
    <div class="space-y-2 max-h-[500px] overflow-y-auto pr-2 custom-scrollbar">
      <div class="space-y-3">
        <SettingsItem 
          v-for="app in apps" 
          :key="app.id" 
          :app="app"
          :isEditing="editingId === app.id"
          :editForm="editForm"
          @edit="$emit('edit', $event)"
          @save="$emit('save', $event)"
          @cancel="$emit('cancel')"
          @delete="$emit('delete', $event)"
        />
      </div>
    </div>
  </div>
</template>

<script setup>
import SettingsItem from './SettingsItem.vue';

defineProps(['apps', 'newApp', 'editingId', 'editForm']);
defineEmits(['add', 'edit', 'save', 'cancel', 'delete']);

const isImg = (str) => str && (str.startsWith('http') || str.startsWith('/'));
</script>

