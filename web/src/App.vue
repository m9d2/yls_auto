<template>
  <div class="mian">
    <div class="header clearfix">
      <div class="left fl">
        <a-button @click="showModal">设置Token</a-button>
        <a-modal v-model:open="open" title="设置Token" @ok="setToken">
          <a-textarea v-model:value="token" placeholder="请输入" />
        </a-modal>
        <a-select ref="select" v-model:value="defaultAddress" style="width: 220px;margin-left: 20px;" @focus="focus"
          @change="handleChange" placeholder="请选择收货地址">
          <a-select-option v-for="item in address" :value="item.id">{{ item.address }}</a-select-option>
        </a-select>
      </div>
      <div class="right fr">
        <a-avatar :src="userInfo.avatar_url" style="margin-right: 8px;" />
        <a-tag>昵称：{{ userInfo.nick_name }}</a-tag>
        <a-tag>当前积分：{{ userInfo.currentScoreSum }}</a-tag>
      </div>

    </div>
    <div class="content">
      <a-table class="table" :columns="columns" :data-source="goods" bordered :pagination="false">
        <template #bodyCell="{ column, text, record }">
          <template v-if="column.dataIndex === 'name'">
            <a>{{ text }}</a>
          </template>
          <template v-if="column.dataIndex === 'point'">
            <a-tag>{{ text }}</a-tag>
          </template>
          <template v-if="column.dataIndex === 'mainImage'">
            <a-image :width="100" :src="text" />
          </template>

          <template v-else-if="column.dataIndex === 'operation'">
            <a-input-number id="inputNumber" v-model:value="record.num" :min="0" :max="10" @change="numChange" />
          </template>
        </template>
      </a-table>

      <div class="submit">
        <a-button type="primary" class="right-button" @click="exchangeGoods">立即兑换</a-button>
      </div>
    </div>
  </div>
</template>

<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { getGoods, getAddress, exchange, getUser } from '@/api/index';
import type { AxiosPromise } from "axios";
import { notification } from 'ant-design-vue';

const goodNum = ref<number>(0);
const token = ref<string>('')
const open = ref<boolean>(false);
const goods = ref<any>();
const address = ref<any>();
const defaultAddress = ref<string>();
const addressId = ref<string>()
const userInfo = ref<any>({
  avatar_url: "",
  nick_name: "",
  currentScoreSum: ""
})

const showModal = () => {
  open.value = true;
};


const openNotificationWithIcon = (type: string, message: string) => {
  notification[type]({
    message: message,
    description:
      '',
  });
};

onMounted(() => {
  const tokenStr = localStorage.getItem('token');
  if (tokenStr) {
    token.value = tokenStr
    loadGoods()
    loadUser()
  }

  console.log(goods.value)

  const addressStr = localStorage.getItem('address');
  if (addressStr) {
    const json = JSON.parse(addressStr);
    defaultAddress.value = json.address
    addressId.value = json.id
  }


  const goodsStr = localStorage.getItem('goods');
  // if (goodsStr) {
  //   const json = JSON.parse(goodsStr);
  //   for (const item of goods.value) {
  //     for (const storeItem of json) {
  //       console.log(item)
  //       console.log(storeItem)
  //       if (item.sku_id === storeItem.sku_id) {
  //         item.num = storeItem.num
  //       }
  //     }
  //   }
  // }
})

const loadUser = () => {
  const response: AxiosPromise<any[]> = getUser(token.value);
  response.then((res: any) => {
    userInfo.value = res.data
  }).catch((error: any) => {
    openNotificationWithIcon('warning', error.message)
  })
}

const loadGoods = () => {
  const response: AxiosPromise<any[]> = getGoods(token.value);
  response.then((res: any) => {
    goods.value = res.data[0].list
    for (const item of goods.value) {
      item.num = 0;
    }
  }).catch((error: any) => {
    openNotificationWithIcon('warning', error.message)
  })
}

const loadAddress = () => {
  const response: AxiosPromise<any[]> = getAddress(token.value);
  response.then((res: any) => {
    address.value = res.data.list
  })
}

const exchangeGoods = () => {
  const data = []
  for (const item of goods.value) {
    if (item.num > 0) {
      const row = { sku_id: item.sku_id, num: item.num }
      data.push(row)
    }
  }
  const response: AxiosPromise<any[]> = exchange({ address_id: addressId.value, product: data }, token.value);
  response.then((res: any) => {
    console.log(res)
  }).catch((error: any) => {
    openNotificationWithIcon('warning', error.message)
  })
}

const setToken = (e: MouseEvent) => {
  localStorage.setItem('token', token.value);
  loadGoods()
  open.value = false;
};

const handleChange = (value: string) => {
  const addr = address.value.find((item: { id: string; }) => item.id === value)
  localStorage.setItem('address', JSON.stringify(addr));
};

const focus = () => {
  loadAddress()
};

const numChange = (value: number) => {

  const data = []
  for (const item of goods.value) {
    if (item.num > 0) {
      const row = { sku_id: item.sku_id, num: item.num }
      data.push(row)
    }
  }
  localStorage.setItem('goods', JSON.stringify(data));
}

const columns = [
  {
    title: '名称',
    dataIndex: 'name',
  },
  {
    title: '所需积分',
    dataIndex: 'point',
  },
  {
    title: '图片',
    dataIndex: 'mainImage',
  },
  {
    title: '兑换数量',
    dataIndex: 'operation',
  }
];
</script>
<style scoped>
.mian {
  width: 1180px;
  /* min-width: 1180px; */
  margin: 0 auto;
}

.header {
  /* margin: 20px; */
  padding: 8px 0;
}

th.column-money,
td.column-money {
  text-align: right !important;
}

.content {
  /* width: 60%; */
  margin: 0 auto;
}

.submit {
  overflow: hidden;
  margin-top: 20px;
}

.right-button {
  float: right;
}

.table {
  max-height: 800px;
  overflow: auto;
}
</style>
