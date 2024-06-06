<script setup lang="ts">
defineOptions({
  inheritAttrs: false,
});
</script>
<template>
  <div class="flex flex-col">
    <div
      class="flex justify-center align-middle transition rounded p-2 m-auto hover:bg-slate-100"
      v-bind="$attrs"
    >
      <slot />
    </div>
  </div>
</template>
