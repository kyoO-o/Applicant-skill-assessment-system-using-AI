<script setup lang="ts">
import type { HTMLAttributes } from "vue"
import { cn } from "@/lib/utils"

const props = defineProps<{
  class?: HTMLAttributes["class"]
}>()
</script>

<template>
  <li
    data-slot="sidebar-menu-item"
    data-sidebar="menu-item"
    :class="cn('group/menu-item relative group-data-[collapsible=icon]:flex group-data-[collapsible=icon]:justify-center', props.class)"
  >
    <slot />
  </li>
</template>
