<template>
  <a :href="status === 'offline' ? '#' : url" 
     :target="status === 'offline' ? '_self' : '_blank'"
     :class="[
       'group p-6 bg-slate-50 dark:bg-slate-900/50 rounded-3xl border transition-all duration-300',
       status === 'offline' 
         ? 'border-red-500/30 opacity-50 cursor-not-allowed' 
         : 'border-slate-700/50 hover:border-blue-500/50 hover:bg-slate-300/80 dark:hover:bg-slate-800/80 hover:-translate-y-1'
     ]">
    
    <div class="flex flex-col items-center text-center space-y-4">
      <div class="relative w-20 h-20 flex items-center justify-center bg-slate-100 dark:bg-slate-900/50 rounded-2xl text-4xl shadow-inner transition-transform overflow-hidden p-2"
           :class="{ 'group-hover:scale-105': status !== 'offline' }">
        
        <img v-if="isImg" :src="icon" class="w-full h-full object-contain" />
        <span v-else>{{ icon || '🌐' }}</span>

        <span class="absolute top-1 right-1 w-3 h-3 rounded-full border-2 border-slate-50 dark:border-slate-900"
              :class="status === 'offline' ? 'bg-red-500 animate-pulse' : 'bg-green-500'">
        </span>
      </div>

      <div>
        <h3 class="font-bold text-lg text-slate-800 dark:text-amber-300 kelly-slab-regular">{{ name }}</h3>
        <p v-if="status === 'offline'" class="text-[9px] text-red-500 font-semibold uppercase tracking-wider mt-1">
          Offline
        </p>
        <p v-else class="text-[10px] text-blue-900 dark:text-amber-50 uppercase tracking-[0.2em] mt-1 truncate w-40"> 
          {{ formatUrl(url) }}
        </p>
      </div>
    </div>
  </a>
</template>

<script setup>
import { computed } from 'vue'

// Added 'status' to expected props (Go passes "online" or "offline")
const props = defineProps(['name', 'url', 'icon', 'status'])

const isImg = computed(() => {
  return props.icon && (props.icon.startsWith('http') || props.icon.startsWith('/'))
})

const formatUrl = (url) => {
  if (!url) return '';
  return url
    .replace('https://', '')
    .replace('http://', '')
    .replace(/\/$/, '');
};

const limitString = (text, limit = 20) => {
  if (!text) return '';
  return text.length > limit ? text.substring(0, limit) + '...' : text;
}
</script>