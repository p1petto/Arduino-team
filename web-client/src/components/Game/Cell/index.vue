<script setup lang="ts">
import type { ICell } from "./index";

const props = withDefaults(defineProps<ICell>(), {
  color: 'white'
});

</script>

<template>
  <div
    class="border hover:border-black size-10"
    :style="{ backgroundColor: props.color }"
  ></div>
</template>
