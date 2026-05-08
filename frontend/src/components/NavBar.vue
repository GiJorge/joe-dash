<template>
  <nav class="border-b border-slate-200 dark:border-slate-800 bg-white/50 dark:bg-slate-900/50 backdrop-blur-md sticky top-0 z-50">
    <div class="max-w-6xl mx-auto px-6 h-16 flex items-center justify-between">
      <h1 class="text-2xl font-bold text-blue-600 dark:text-blue-400 tracking-wide croissant-one-regular ">Joe Dash</h1>
      
      <div class="flex items-center gap-4">
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
defineProps(['isSettings'])
defineEmits(['toggle-settings'])
</script>