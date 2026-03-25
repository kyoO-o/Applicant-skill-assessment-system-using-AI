<script setup lang="ts">
const props = defineProps({
  error: Object,
})
const handleError = () => location.reload()
</script>

<template>
  <div class="container">
    <v-container min-width="100%" class="mt-2 app-bar">
      <v-card rounded="pill" border="sm" flat class="mx-5">
        <template #prepend>
          <img src="/icons/logo.svg" width="130" class="ml-3" />
        </template>
      </v-card>
    </v-container>

    <div class="d-flex flex-column align-center body pt-16" v-if="props.error?.statusCode == 404">
      <img src="/not-found.png" style="max-width: 280px" alt="" />
      <span class="font-weight-bold" style="font-size: 80px">Өөөө.... 404</span>
      <span class="text-h5 text-secondary">Таны хайсан хуудас байхгүй байна. 😅</span>
      <span class="text-h7 text-primary">Та бага тэнэнэ үү.</span>
      <v-btn prepend-icon="fr-home" class="mt-6" href="/" color="primary">
        Нүүр хуудасруу очих</v-btn
      >
      <a
        href="https://www.reddit.com/r/Catmemes/comments/12zxy5g/missing_cat/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button"
        class="text-body-2 mt-4"
        color="#2060EB"
        target="_blank"
      >
        Энэ муурны талаар мэдээлэл авахыг хүсвэл энд дарна уу.
      </a>
    </div>
    <div v-else class="text-center body pt-16 mt-16">
      <h2 class="text-h3 font-weight-light">
        <template v-if="props.error?.statusCode == 400">Оролтын алдаа.</template>
        <!-- <template v-else-if="props.error?.statusCode == 404">Хуудас олдсонгүй.</template> -->
        <template v-else-if="props.error?.statusCode == 500">Серверийн алдаа.</template>
        <template v-else>Тодорхойгүй алдаа.</template>
      </h2>
      <h2 class="text-h1">{{ props.error?.statusCode }}</h2>
      <h2 class="text-h6 font-weight-regular text-info mt-4">
        <template v-if="props.error?.statusCode == 400"
          >Та оролтын утгуудаа шалгаад дахин оролдоно уу.</template
        >
        <template v-else-if="props.error?.statusCode == 404"
          >Энэ хуудас устсан эсвэл огт байхгүй гэсэн үг. Та линкээ шалгана уу.</template
        >
        <template v-else-if="props.error?.statusCode == 500"
          >Та түр хүлээгээд дахин оролдоно уу. Үүний дараа ч болохгүй байвал админд мэдэгдэнэ
          үү.</template
        >
        <template v-else
          >Та түр хүлээгээд дахин оролдоно уу.<br />Дахин оролдсон ч болохгүй байгаа бол админд
          мэдэгдэнэ үү.</template
        >
      </h2>
      <v-btn class="mt-6" color="primary" @click="handleError"> Дахин ачааллах </v-btn>
    </div>
  </div>
</template>
<style scoped lang="scss">
.container {
  height: 100vh;
  width: 100vw;
  overflow: hidden;
}
.app-bar {
  // border-bottom: 1px solid #dfe1e0;
  background-color: #f7f8f8;
  height: 64px;
  width: 100%;
}

.body {
  background-color: #f7f8f8;
  height: calc(100vh - 64px);
}
</style>
