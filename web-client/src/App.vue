<script setup lang="ts">
import { Cell, type ICell } from "@/components/Game/Cell/index";
import { Button, Navbar, Tooltip } from "@/components/UI/index";

import {
  ArrowPathIcon,
  PaintBrushIcon,
  UserIcon,
} from "@heroicons/vue/24/outline";
import { ref } from "vue";

const cells = ref<ICell[]>([]);
const pallete = ["red", "blue", "green", "white"];
const currentColor = ref(0);

for (let i = 0; i < 16 * 16; i++) {
  cells.value.push({});
}

const initGame = () => {
  cells.value.forEach((c) => (c.color = "white"));
};

const colorizeCell = (idx: number) => {
  cells.value[idx].color = pallete[currentColor.value];
};

const selectColor = (idx: number) => {
  currentColor.value = idx;
  colorPickerVisible.value = false;
};

const colorPickerVisible = ref(false);
</script>

<template>
  <div class="flex flex-col h-screen w-screen bg-[#fafafa]">
    <Navbar class="flex flex-row justify-between">
      <div class="m-auto font-bold">Текущая игра</div>
      <Button>
        <UserIcon class="size-6 text-slate-500" />
      </Button>
    </Navbar>

    <div class="flex flex-grow">
      <!-- Панель с инструментами -->
      <div class="flex flex-col justify-between p-2 bg-white shadow">
        <div>
          <Button
            :style="{ borderBottom: `3px solid ${pallete[currentColor]}` }"
            @click="colorPickerVisible = !colorPickerVisible"
          >
            <PaintBrushIcon class="size-6 text-slate-500" />
          </Button>

          <Tooltip v-if="colorPickerVisible">
            <template #header>Выбрать цвет</template>
            <div class="grid grid-cols-4">
              <div
                class="rounded size-6 border transition-all hover:shadow hover:scale-110"
                :style="{ backgroundColor: color }"
                @click="selectColor(idx)"
                v-for="(color, idx) in pallete"
              ></div>
            </div>
          </Tooltip>
        </div>

        <Button @click="initGame()">
          <ArrowPathIcon class="size-6 text-slate-500" />
        </Button>
      </div>

      <!-- Игровая сетка -->
      <div class="grid grid-cols-[repeat(16,_minmax(0,_1fr))] m-auto">
        <Cell
          :color="color"
          :key="idx"
          @click="colorizeCell(idx)"
          v-for="({ color }, idx) in cells"
        />
      </div>
    </div>
  </div>
</template>
