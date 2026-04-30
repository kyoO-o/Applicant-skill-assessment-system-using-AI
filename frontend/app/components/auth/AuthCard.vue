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
  <Card class="mx-auto w-full max-w-xl rounded-[2rem] border-white/80 bg-white/90 shadow-2xl shadow-slate-200/70 backdrop-blur">
    <CardHeader class="space-y-3">
      <Badge variant="secondary" class="w-fit rounded-full px-3 py-1">
        Secure access
      </Badge>
      <CardTitle class="text-2xl text-slate-950 sm:text-3xl">
        {{ title }}
      </CardTitle>
      <CardDescription class="text-base leading-7 text-slate-600">
        {{ description }}
      </CardDescription>
    </CardHeader>

    <CardContent class="space-y-6">
      <slot />

      <Separator />

      <p class="text-center text-sm text-slate-600">
        {{ footerText }}
        <NuxtLink
          :to="footerLinkTo"
          class="font-semibold text-blue-700 transition hover:text-blue-800"
        >
          {{ footerLinkText }}
        </NuxtLink>
      </p>
    </CardContent>
  </Card>
</template>
