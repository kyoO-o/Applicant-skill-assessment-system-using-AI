<template>
  <div class="space-y-6">
    <section
      class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm"
    >
      <p class="text-sm font-medium text-muted-foreground">Settings</p>
      <h1 class="mt-2 text-3xl font-semibold tracking-tight">
        Workspace preferences
      </h1>
      <p class="mt-2 max-w-2xl text-sm leading-6 text-muted-foreground">
        The page now matches the new layout language while leaving the deeper
        account behavior untouched.
      </p>
    </section>

    <section class="grid gap-4 lg:grid-cols-2">
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Notifications</CardTitle>
          <CardDescription>Interview updates, new applicants, and reminders.</CardDescription>
        </CardHeader>
        <CardContent class="space-y-3">
          <div class="rounded-2xl border border-border bg-muted/30 p-4 text-sm text-muted-foreground">
            Placeholder panel ready for notification preferences.
          </div>
        </CardContent>
      </Card>
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Appearance</CardTitle>
          <CardDescription>Aligned with the shared tailwind theme tokens.</CardDescription>
        </CardHeader>
        <CardContent class="space-y-3">
          <div class="rounded-2xl border border-border bg-muted/30 p-4 text-sm text-muted-foreground">
            Placeholder panel ready for density and layout options.
          </div>
        </CardContent>
      </Card>
    </section>
  </div>
</template>
