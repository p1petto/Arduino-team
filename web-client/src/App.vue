<script setup lang="ts">
import { ref } from "vue";
import Navbar from "@/components/UI/Navbar.vue";
import Cell from "@/components/Game/Cell/index.vue";
import type { ICell } from "@/components/Game/Cell/types.ts";

const cells = ref<ICell[]>([]);
for (let i = 0; i < 24 * 12; i++) {
  cells.value.push({
    color: "foo",
  });
}

cells.value[0].color = "foo";
</script>

<template>
  <Navbar>
    <a href="#">foo</a>
  </Navbar>
  <div
    class="grid grid-rows-24 grid-cols-12 grid-flow-row h-full w-full"
  >
    <Cell :data="cell" v-for="cell of cells" />
  </div>
</template>

<style scoped></style>
