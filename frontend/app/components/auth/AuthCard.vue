<script setup lang="ts">
interface Props {
  title: string;
  description: string;
  footerText: string;
  footerLinkText: string;
  footerLinkTo: string;
}

defineProps<Props>();
</script>

<template>
  <div>
    <!-- Heading -->
    <h1 class="text-[26px] font-semibold tracking-[-0.6px] text-foreground">{{ title }}</h1>
    <p class="mt-1.5 text-[14px] text-muted-foreground">{{ description }}</p>

    <!-- Form content -->
    <div class="mt-7">
      <slot />
    </div>

    <!-- Footer link -->
    <p class="mt-6 text-center text-[13px] text-muted-foreground">
      {{ footerText }}
      <NuxtLink
        :to="footerLinkTo"
        class="ml-1 font-semibold text-primary hover:text-primary/80"
      >
        {{ footerLinkText }}
      </NuxtLink>
    </p>
  </div>
</template>
