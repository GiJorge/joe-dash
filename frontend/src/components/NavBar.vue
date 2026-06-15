<template>
  <nav class="border-b border-slate-200 dark:border-slate-800 bg-white/50 dark:bg-slate-900/50 backdrop-blur-md sticky top-0 z-50">
    <div class="max-w-6xl mx-auto px-6 h-16 flex items-center justify-between">
        
       <!-- Clickable Logo/Title -->
      <h1 
        @click="$emit('go-home')" 
        class="text-2xl font-bold text-blue-600 dark:text-blue-400 tracking-wide cursor-pointer hover:opacity-80 transition-opacity font-croissant"
      >
        Joe Dash
      </h1>
      <div class="flex gap-4 text-[10px] font-mono">



</div>

      <div class="flex items-center gap-4  ">
        <!-- 1. ADD THE SLOT HERE -->
        <slot /> 

        <button @click="$emit('toggle-settings')" 
          class="flex items-center gap-2 px-5 py-2 rounded-full bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 border border-slate-200 dark:border-slate-700 transition-all text-sm font-medium">
          <span>{{ isSettings ? '🏠 Home' : '⚙️ Settings' }}</span>
        </button>
      </div>
    </div>
  </nav>
</template>

<script setup>
import { ref, onMounted } from 'vue';
defineProps(['isSettings'])
defineEmits(['toggle-settings'])

const stats = ref({ cpu: 0, ram: 0 });

const fetchStats = async () => {
  const res = await fetch('/api/stats');
  stats.value = await res.json();
};

// Update every 5 seconds
onMounted(() => {
  fetchStats();
  setInterval(fetchStats, 2000);
});
</script>