<script setup lang="ts">
defineOptions({
  inheritAttrs: false,
});
defineProps(["icon"]);
</script>
<template>
  <div class="flex flex-col">
    <div
      class="flex justify-center align-middle transition rounded p-2 m-auto hover:bg-slate-100"
      v-bind="$attrs"
    >
      <component :is="icon" class="size-6 text-slate-500" />
    </div>
  </div>
</template>
