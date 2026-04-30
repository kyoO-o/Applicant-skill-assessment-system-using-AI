<script setup lang="ts">
const { user } = useAuth();
</script>

<template>
  <div class="space-y-6">
    <section
      class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm"
    >
      <p class="text-sm font-medium text-muted-foreground">Profile</p>
      <h1 class="mt-2 text-3xl font-semibold tracking-tight">
        {{ user?.name || "Your profile" }}
      </h1>
      <p class="mt-2 text-sm leading-6 text-muted-foreground">
        This page keeps your current user model intact and presents it with the
        wireframe-inspired card layout.
      </p>
    </section>

    <section class="grid gap-4 lg:grid-cols-[1.3fr_0.9fr]">
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Account details</CardTitle>
          <CardDescription>Synced from the authenticated backend user.</CardDescription>
        </CardHeader>
        <CardContent class="grid gap-4 sm:grid-cols-2">
          <div class="rounded-2xl border border-border bg-muted/30 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Email
            </p>
            <p class="mt-2 text-sm font-medium">{{ user?.email || "-" }}</p>
          </div>
          <div class="rounded-2xl border border-border bg-muted/30 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Role
            </p>
            <p class="mt-2 text-sm font-medium capitalize">{{ user?.role || "-" }}</p>
          </div>
          <div class="rounded-2xl border border-border bg-muted/30 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Company
            </p>
            <p class="mt-2 text-sm font-medium">{{ user?.company_name || "-" }}</p>
          </div>
          <div class="rounded-2xl border border-border bg-muted/30 p-4">
            <p class="text-xs uppercase tracking-[0.18em] text-muted-foreground">
              Position
            </p>
            <p class="mt-2 text-sm font-medium">{{ user?.position || "-" }}</p>
          </div>
        </CardContent>
      </Card>

      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Next</CardTitle>
          <CardDescription>Quick path based on your role.</CardDescription>
        </CardHeader>
        <CardContent class="space-y-3">
          <NuxtLink
            to="/jobs"
            class="block rounded-2xl border border-border bg-muted/30 p-4 text-sm font-medium transition hover:bg-muted/50"
          >
            {{ user?.role === "recruiter" ? "Manage job posts" : "Browse job matches" }}
          </NuxtLink>
          <NuxtLink
            to="/settings"
            class="block rounded-2xl border border-border bg-muted/30 p-4 text-sm font-medium transition hover:bg-muted/50"
          >
            Review account settings
          </NuxtLink>
        </CardContent>
      </Card>
    </section>
  </div>
</template>
