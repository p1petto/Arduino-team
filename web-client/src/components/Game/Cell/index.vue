<script setup lang="ts">
import type { ICell } from "./types.ts";

interface ICellProps {
  data: ICell;
}

const props = defineProps<ICellProps>();
</script>

<template>
  <div
    class="border hover:border-black"
    :style="{ backgroundColor: props.data.color }"
  ></div>
</template>
