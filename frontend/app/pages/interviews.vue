<script setup lang="ts">
const { user } = useAuth();

const recruiterItems = [
  {
    title: "Frontend Developer - Technical Screen",
    subtitle: "Tomorrow at 10:00",
    status: "Scheduled",
  },
  {
    title: "UI/UX Designer - Portfolio Review",
    subtitle: "Friday at 14:30",
    status: "Pending confirmation",
  },
];

const applicantItems = [
  {
    title: "Tech screening with MatchHire recruiter",
    subtitle: "Awaiting invitation",
    status: "No upcoming interviews",
  },
];
</script>

<template>
  <div class="space-y-6">
    <section
      class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm"
    >
      <div class="flex flex-col gap-3 md:flex-row md:items-end md:justify-between">
        <div class="space-y-2">
          <p class="text-sm font-medium text-muted-foreground">Interviews</p>
          <h1 class="text-3xl font-semibold tracking-tight">
            {{ user?.role === "recruiter" ? "Interview pipeline" : "Interview updates" }}
          </h1>
          <p class="max-w-2xl text-sm leading-6 text-muted-foreground">
            Wireframe-inspired scheduling space is connected into the app shell.
            The actions are intentionally lightweight here so we do not alter your
            existing hiring flow.
          </p>
        </div>
      </div>
    </section>

    <section class="grid gap-4 lg:grid-cols-2">
      <Card
        v-for="item in user?.role === 'recruiter' ? recruiterItems : applicantItems"
        :key="item.title"
        class="rounded-3xl border-border shadow-sm"
      >
        <CardHeader class="space-y-2">
          <CardTitle class="text-xl">{{ item.title }}</CardTitle>
          <CardDescription>{{ item.subtitle }}</CardDescription>
        </CardHeader>
        <CardContent>
          <Badge variant="secondary" class="rounded-full px-3 py-1">
            {{ item.status }}
          </Badge>
        </CardContent>
      </Card>
    </section>
  </div>
</template>
