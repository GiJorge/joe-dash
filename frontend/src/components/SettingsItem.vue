<template>
  <div class="bg-slate-200/50 border dark:bg-slate-900/50  border-slate-800 rounded-2xl p-4">
    <!-- EDIT MODE -->
    <div v-if="isEditing" class="space-y-3">
      <div class="grid grid-cols-2 gap-2">
        <input v-model="editForm.name" class="input-style text-sm" placeholder="Name">
        <input v-model="editForm.icon" class="input-style text-sm" placeholder="Icon/URL">
      </div>
      <input v-model="editForm.url" class="input-style text-sm w-full" placeholder="URL">
      <div class="flex justify-end gap-3">
        <button @click="$emit('save', editForm)" class="text-xs font-bold text-green-400">UPDATE</button>
        <button @click="$emit('cancel')" class="text-xs font-bold text-slate-500">CANCEL</button>
      </div>
    </div>

    <!-- DISPLAY MODE -->
    <div v-else class="flex items-center justify-between">
      <div class="flex items-center gap-4">
        <div class="w-12 h-12 bg-slate-100 dark:bg-slate-300  rounded-xl flex items-center justify-center overflow-hidden p-1">
          <img v-if="isImg(app.icon)" :src="app.icon" class="w-full h-full object-contain" />
          <span v-else class="text-xl">{{ app.icon || '🔗' }}</span>
        </div>
        <div>
          <h4 class="font-bold text-slate-800 dark:text-amber-300 kelly-slab-regular ">{{ app.name }}</h4>
          <!-- Using your requested URL shortening logic here -->
          <p class="text-[13px] text-gray-800 dark:text-amber-50 truncate w-40">
           {{ formatUrl(app.url) }}
          </p>
        </div>
      </div>
      <div class="flex gap-4">
        <button @click="$emit('edit', app)" class="text-blue-900 dark:text-blue-400 hover:text-blue-600">Edit</button>
        <button @click="$emit('delete', app.id)" class="text-red-900 dark:text-red-400 hover:text-red-300">Delete</button>
      </div>
    </div>
  </div>
</template>

<script setup>
defineProps(['app', 'isEditing', 'editForm']);
defineEmits(['edit', 'save', 'cancel', 'delete']);

const isImg = (str) => str && (str.startsWith('http') || str.startsWith('/'));

const formatUrl = (url) => {
  if (!url) return '';
  return url
    .replace('http://', '')
    .replace('https://', '')
    .replace(/\/$/, '');
};
</script>