<template>
  <a :href="url" target="_blank"
    class="group p-6 bg-slate-50 dark:bg-slate-900/50 rounded-3xl border border-slate-700/50 hover:border-blue-500/50 hover:bg-slate-300/80 dark:hover:bg-slate-800/80  transition-all duration-300 hover:-translate-y-1">
    <div class="flex flex-col items-center text-center space-y-4">
      <div class="w-20 h-20 flex items-center justify-center bg-slate-100 dark:bg-slate-900/50 rounded-2xl text-4xl shadow-inner group-hover:scale-105 transition-transform overflow-hidden p-2">
        <img v-if="isImg" :src="icon" class="w-full h-full object-contain" />
        <span v-else>{{ icon || '🌐' }}</span>
      </div>
      <div>
        <h3 class="font-bold text-lg text-slate-800 dark:text-amber-300 kelly-slab-regular">{{ name }}</h3>
        <p class="text-[11px] text-blue-900 dark:text-amber-50 uppercase tracking-[0.2em] mt-1"> {{ formatUrl(url) }}</p>
      </div>
    </div>
  </a>
</template>

<script setup>
import { computed } from 'vue'
const props = defineProps(['name', 'url', 'icon'])

const isImg = computed(() => {
  return props.icon && (props.icon.startsWith('http') || props.icon.startsWith('/'))
})


const formatUrl = (url) => {
  if (!url) return '';
  return url
    .replace('https://', '')
    .replace('http://', '')
    .replace('https://', '')
    .replace(/\/$/, '');
};

</script>