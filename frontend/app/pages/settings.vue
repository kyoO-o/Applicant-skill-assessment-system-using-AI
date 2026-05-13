<template>
  <div class="space-y-6">
    <section
      class="rounded-3xl border border-border bg-card px-6 py-6 shadow-sm"
    >
      <p class="text-sm font-medium text-muted-foreground">Тохиргоо</p>
      <h1 class="mt-2 text-3xl font-semibold tracking-tight">
        Ажлын орчны тохиргоо
      </h1>
      <p class="mt-2 max-w-2xl text-sm leading-6 text-muted-foreground">
        Бүртгэлийн болон дэлгэцийн тохиргоог энд хийнэ үү.
      </p>
    </section>

    <section class="grid gap-4 lg:grid-cols-2">
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Мэдэгдэл</CardTitle>
          <CardDescription>Ярилцлага, шинэ горилогч болон сануулгын тохиргоо.</CardDescription>
        </CardHeader>
        <CardContent class="space-y-3">
          <div class="rounded-2xl border border-border bg-muted/30 p-4 text-sm text-muted-foreground">
            Мэдэгдлийн тохиргоо удахгүй нэмэгдэнэ.
          </div>
        </CardContent>
      </Card>
      <Card class="rounded-3xl border-border shadow-sm">
        <CardHeader>
          <CardTitle>Харагдац</CardTitle>
          <CardDescription>Дэлгэцийн загвар болон өнгөний тохиргоо.</CardDescription>
        </CardHeader>
        <CardContent class="space-y-3">
          <div class="rounded-2xl border border-border bg-muted/30 p-4 text-sm text-muted-foreground">
            Харагдацын тохиргоо удахгүй нэмэгдэнэ.
          </div>
        </CardContent>
      </Card>
    </section>
  </div>
</template>
