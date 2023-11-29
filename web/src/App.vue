<template>
  <a-config-provider :locale="locale" :theme="{
      token: {
        colorPrimary: '#3b5999',
      },
    }">
    <RouterView></RouterView>
  </a-config-provider>
</template>

<script lang="ts" setup>
import zhCN from "ant-design-vue/es/locale/zh_CN";
import dayjs from "dayjs"
import "dayjs/locale/zh-cn"
dayjs.locale("zh-cn")
const locale = zhCN

</script>
<style scoped></style>
