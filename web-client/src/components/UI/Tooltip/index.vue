<template>
  <div
    class="absolute left-14 p-4 rounded-lg bg-white border-2 border-slate-200 shadow"
  >
    <h4 class="font-extrabold">
      <slot name="header" />
    </h4>
    <slot />
  </div>
</template>

<script setup lang="ts"></script>
