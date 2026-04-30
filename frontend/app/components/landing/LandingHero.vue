<template>
  <section class="overflow-hidden">
    <div class="mx-auto grid max-w-7xl gap-14 px-4 py-20 sm:px-6 lg:grid-cols-[1.05fr_0.95fr] lg:items-center lg:px-8 lg:py-24">
      <div class="space-y-8">
        <Badge class="rounded-full px-3 py-1 shadow-sm shadow-blue-200/50">
          AI Recruitment Platform
        </Badge>

        <div class="space-y-5">
          <h1 class="max-w-3xl text-4xl font-semibold tracking-tight text-slate-950 sm:text-5xl lg:text-6xl">
            AI-Powered Applicant Skill Assessment System
          </h1>
          <p class="max-w-2xl text-base leading-8 text-slate-600 sm:text-lg">
            Match candidates with job requirements using CV analysis, skill
            scoring, and explainable AI.
          </p>
        </div>

        <div class="flex flex-col gap-3 sm:flex-row">
          <Button as-child size="lg" class="shadow-sm shadow-blue-200/60">
            <NuxtLink to="/register">Get Started</NuxtLink>
          </Button>
          <Button as-child size="lg" variant="outline">
            <NuxtLink to="/login">Login</NuxtLink>
          </Button>
        </div>

        <div class="grid gap-4 sm:grid-cols-3">
          <div class="rounded-3xl border border-slate-200 bg-white/90 p-5 shadow-sm">
            <p class="text-sm text-slate-500">Assessment</p>
            <p class="mt-2 text-lg font-semibold text-slate-950">Skill scoring</p>
          </div>
          <div class="rounded-3xl border border-slate-200 bg-white/90 p-5 shadow-sm">
            <p class="text-sm text-slate-500">Matching</p>
            <p class="mt-2 text-lg font-semibold text-slate-950">Job alignment</p>
          </div>
          <div class="rounded-3xl border border-slate-200 bg-white/90 p-5 shadow-sm">
            <p class="text-sm text-slate-500">Explainability</p>
            <p class="mt-2 text-lg font-semibold text-slate-950">Transparent AI</p>
          </div>
        </div>
      </div>

      <div class="relative">
        <div
          class="absolute inset-0 -z-10 scale-95 rounded-[2rem] bg-gradient-to-br from-blue-200/50 via-white to-violet-200/50 blur-3xl"
        />

        <Card class="rounded-[2rem] border-white/80 bg-white/90 shadow-2xl shadow-slate-200/70 backdrop-blur">
          <CardHeader class="space-y-4">
            <div class="flex items-center justify-between">
              <Badge variant="secondary" class="rounded-full px-3 py-1">
                Live Preview
              </Badge>
              <div class="flex items-center gap-2">
                <span class="h-2.5 w-2.5 rounded-full bg-blue-400" />
                <span class="h-2.5 w-2.5 rounded-full bg-slate-300" />
                <span class="h-2.5 w-2.5 rounded-full bg-slate-300" />
              </div>
            </div>

            <CardTitle class="text-2xl text-slate-950">
              Candidate assessment overview
            </CardTitle>
            <CardDescription class="text-base leading-7 text-slate-600">
              A recruiter-friendly snapshot of candidate fit, highlighted skills,
              and actionable next steps.
            </CardDescription>
          </CardHeader>

          <CardContent class="space-y-6">
            <div class="rounded-3xl bg-slate-50 p-5">
              <div class="flex items-start justify-between gap-4">
                <div>
                  <p class="text-sm text-slate-500">Candidate</p>
                  <p class="mt-1 text-lg font-semibold text-slate-950">
                    Frontend Engineer Applicant
                  </p>
                </div>
                <Badge class="rounded-full px-3 py-1">86% Match</Badge>
              </div>

              <div class="mt-5 grid gap-3 sm:grid-cols-2">
                <div class="rounded-2xl bg-white p-4 shadow-sm">
                  <p class="text-sm text-slate-500">Top Skills</p>
                  <p class="mt-2 text-sm font-medium text-slate-900">
                    Vue, TypeScript, Tailwind CSS
                  </p>
                </div>
                <div class="rounded-2xl bg-white p-4 shadow-sm">
                  <p class="text-sm text-slate-500">Recommendation</p>
                  <p class="mt-2 text-sm font-medium text-slate-900">
                    Move to technical interview
                  </p>
                </div>
              </div>
            </div>

            <div class="grid gap-4 sm:grid-cols-3">
              <div class="rounded-2xl border border-slate-200 p-4">
                <p class="text-sm text-slate-500">CV Parsing</p>
                <p class="mt-2 text-xl font-semibold text-slate-950">Ready</p>
              </div>
              <div class="rounded-2xl border border-slate-200 p-4">
                <p class="text-sm text-slate-500">AI Score</p>
                <p class="mt-2 text-xl font-semibold text-slate-950">8.6 / 10</p>
              </div>
              <div class="rounded-2xl border border-slate-200 p-4">
                <p class="text-sm text-slate-500">Interview</p>
                <p class="mt-2 text-xl font-semibold text-slate-950">Suggested</p>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>
    </div>
  </section>
</template>
