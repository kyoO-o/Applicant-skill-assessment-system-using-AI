<template>
  <div class="space-y-6">
    <section
      class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm"
    >
      <p class="text-sm font-medium text-muted-foreground">My Zone</p>
      <h1 class="mt-2 text-3xl font-semibold tracking-tight">
        Personal workspace
      </h1>
      <p class="mt-2 max-w-2xl text-sm leading-6 text-muted-foreground">
        This keeps the wireframe structure in place for applicants without
        replacing your existing business logic with mock-only behavior.
      </p>
    </section>

    <section class="grid gap-4 md:grid-cols-3">
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Applications</CardTitle>
          <CardDescription>Track submitted jobs in one place.</CardDescription>
        </CardHeader>
      </Card>
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Assessments</CardTitle>
          <CardDescription>Resume the next invited skill assessment.</CardDescription>
        </CardHeader>
      </Card>
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Documents</CardTitle>
          <CardDescription>Keep CV and portfolio versions organized.</CardDescription>
        </CardHeader>
      </Card>
    </section>
  </div>
</template>
