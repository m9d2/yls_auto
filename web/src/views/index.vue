<template>

  <div style="min-width: 1180px">
    <!-- 头 -->
    <a-layout-header class="header clearfix">
      <div class="left fl">
        <a-space>
          <a-button @click="showModal" style="margin-left: 8px">设置Token</a-button>
          <a-modal v-model:open="open" title="设置Token" @ok="setToken">
            <a-textarea v-model:value="token" placeholder="请输入" style="height: 200px;"/>
          </a-modal>
          <a-select ref="select" v-model:value="defaultAddress" style="width: 220px;" @focus="focus"
                    @change="handleChange" placeholder="请选择收货地址">
            <a-select-option v-for="item in address" :value="item.id">{{ item.address }}</a-select-option>
          </a-select>
        </a-space>
      </div>
      <div class="right fr" v-if="userInfo">
        <a-space>
          <a><a-avatar :src="userInfo.avatar_url" @click="openUserInfo = true"/></a>
          <a-tag>昵称：{{ userInfo.nick_name }}</a-tag>
          <a-tag>当前积分：{{ userInfo.currentScoreSum }}</a-tag>
        </a-space>
      </div>
      <a-drawer
          v-model:open="openUserInfo"
          class="custom-class"
          root-class-name="root-class-name"
          :root-style="{ color: 'blue' }"
          title="个人信息"
          :width="420"
          placement="right"
      >
        <a-descriptions :title="userInfo.nick_name" :column="1">
          <a-descriptions-item label=""><a-avatar :src="userInfo.avatar_url"/></a-descriptions-item>
          <a-descriptions-item label="手机号">{{ userInfo.phone_number }}</a-descriptions-item>
          <a-descriptions-item label="生日">{{ userInfo.birthday }}</a-descriptions-item>
          <a-descriptions-item label="会员ID">{{ userInfo.consumerID }}</a-descriptions-item>
          <a-descriptions-item label="会员等级"><a-tag color="#3b5999">{{ userInfo.grade }}</a-tag>（{{userInfo.expireDateText}}）</a-descriptions-item>
          <a-descriptions-item label="积分"><a-tag color="#108ee9">{{ userInfo.currentScoreSum }}</a-tag></a-descriptions-item>
          <a-descriptions-item label="积分有效期">{{ userInfo.expireScoreDate }}</a-descriptions-item>
        </a-descriptions>
      </a-drawer>
    </a-layout-header>

    <!-- 内容区 -->
    <a-layout-content class="content">
      <a-table class="table" :columns="columns" :data-source="goods" bordered :pagination="false"
               :loading="loadingTable">
        <template #bodyCell="{ column, text, record }">
          <template v-if="column.dataIndex === 'name'">
            <a>{{ text }}</a>
          </template>
          <template v-if="column.dataIndex === 'point'">
            <a-tag>{{ text }}</a-tag>
          </template>
          <template v-if="column.dataIndex === 'mainImage'">
            <a-image :width="58" :src="text"/>
          </template>

          <template v-else-if="column.dataIndex === 'num'">
            <a-input-number id="inputNumber" v-model:value="record.num" :min="0" :max="10" @change="numChange"/>
          </template>
        </template>
      </a-table>

      <div class="submit fr">
        <a-space>
          <a-button class="right-button" type="primary" @click="loadGoods">刷新</a-button>
          <a-button class="right-button" type="primary" @click="showDrawer">定时兑换</a-button>
          <a-button class="right-button" @click="exchangeGoods" :loading="loadingSubmit">立即兑换
          </a-button>
        </a-space>
      </div>
    </a-layout-content>

    <!-- 抽屉 -->
    <div class="drawer">
      <a-drawer
          title="定时兑换"
          :width="420"
          :open="drawerStatus"
          :body-style="{ paddingBottom: '80px' }"
          :footer-style="{ textAlign: 'right' }"
          @close="oncloseDrawer"
      >
        <a-form :model="form" layout="vertical">
          <a-form-item label="开始时间" name="time" style="margin-top: 8px">
            <a-time-picker format="HH:mm:ss" @change="onRangeChange"/>
          </a-form-item>
          <a-form-item label="延迟时间（毫秒）" name="time" style="margin-top: 8px">
            <a-input v-model:value="delayTime" style="width: 123px;"/>
          </a-form-item>
          <a-table class="table" :columns="columnsSub" :data-source="goodsSub" bordered :pagination="false"
                   :loading="loadingTable">
            <template #bodyCell="{ column, text, record }">
              <template v-if="column.dataIndex === 'name'">
                <a>{{ text }}</a>
              </template>
            </template>
          </a-table>
        </a-form>
        <template #extra>
          <a-space>
            <a-button type="primary" @click="stop" v-if="start" danger>暂停任务</a-button>
            <a-button type="primary" @click="autoExchange" v-else>开始任务</a-button>
          </a-space>
        </template>

        <div style="margin-top: 20px;" v-if="output.length > 0">
          <a-card title="执行日志" :bordered="true">
            <div ref="outputRef" class="output" style="height: 260px; overflow: auto">
              <p v-for="item in output">
                <span style="color: #1677FF;">{{ item.time }}</span>
                - {{ item.message }}
              </p>
            </div>
          </a-card>
        </div>
      </a-drawer>
    </div>

    <!-- 底部 -->
    <a-layout-footer class="footer">
      <a-button type="link" href="https://github.com/m9d2/ysl_auto" target="_blank" style="padding: 0; color: #000; font-size: 12px;">m9d2/ysl_auto</a-button>
      <a-button type="link" href="https://github.com/m9d2/ysl_auto/releases" target="_blank" style="padding: 0; color: #000; font-size: 12px;">- Release {{version}}</a-button>
    </a-layout-footer>
  </div>
</template>

<script lang="ts" setup>
import {onMounted, ref, reactive} from 'vue';
import {getGoods, getAddress, exchange, getUser, getTags} from '@/api';
import type {AxiosPromise} from "axios";
import {message} from 'ant-design-vue';
import {Dayjs} from 'dayjs';

const token = ref<string>('')
const open = ref<boolean>(false);
const openUserInfo = ref<boolean>(false);
const goods = ref<any>();
const goodsSub = ref<any>();
const address = ref<any>();
const defaultAddress = ref<string>();
const addressId = ref<string>()
const userInfo = ref<any>()
const loadingSubmit = ref<boolean>(false)
const loadingTable = ref<boolean>(false)
const drawerStatus = ref<boolean>(false);
const outputRef = ref(null);
const start = ref<Boolean>(false)
const timeStr = ref<String>()
const delayTime = ref<Number>(500)
const form = reactive({});
const output = ref([])
let version = ref()
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
    dataIndex: 'num',
  }
];
const columnsSub = [
  {
    title: '名称',
    dataIndex: 'name',
  },
  {
    title: '兑换数量',
    dataIndex: 'num',
  }
];

onMounted(() => {
  getLastTag()
  const tokenStr = localStorage.getItem('token');
  if (tokenStr) {
    token.value = tokenStr
    loadGoods()
    loadUser()
  }

  const addressStr = localStorage.getItem('address');
  if (addressStr) {
    const json = JSON.parse(addressStr);
    defaultAddress.value = json.address
    addressId.value = json.id
    console.log(addressId.value)
  }
})

const loadUser = () => {
  const response: AxiosPromise<any[]> = getUser(token.value);
  response.then((res: any) => {
    userInfo.value = res.data
  }).catch((error: any) => {
    //
  })
}

const getLastTag = () => {
  const response: AxiosPromise<any[]> = getTags();
  response.then((res: any) => {
    if (res) {
      version.value = res[0].name
    }
  }).catch((error) => {
    //
  })
}

const loadGoods = () => {
  loadingTable.value = true
  const response: AxiosPromise<any[]> = getGoods(token.value);
  response.then((res: any) => {
    goods.value = res.data[0].list
    for (const item of goods.value) {
      item.num = 0;
    }

    const goodsStr = localStorage.getItem('goods');
    if (goodsStr && goods.value) {
      const json = JSON.parse(goodsStr);
      goodsSub.value = json;
      for (const item of goods.value) {
        for (const storeItem of json) {
          if (item.sku_id === storeItem.sku_id) {
            item.num = storeItem.num
          }
        }
      }
    }
  }).catch((error: any) => {
    message.error(error.message);
    throw error
  }).finally(() => {
    loadingTable.value = false
  })
}

const loadAddress = () => {
  const response: AxiosPromise<any[]> = getAddress(token.value);
  response.then((res: any) => {
    address.value = res.data.list
  }).catch((error: any) => {
    message.error(error.message);
  })
}

const exchangeGoods = () => {
  loadingSubmit.value = true
  const data = []
  for (const item of goods.value) {
    if (item.num > 0) {
      const row = {sku_id: item.sku_id, num: item.num}
      data.push(row)
    }
  }
  const response: AxiosPromise<any[]> = exchange({address_id: addressId.value, product: data}, token.value);
  response.then((res: any) => {
    console.log(res)
    message.success("兑换成功！");
  }).catch((error: any) => {
    message.error(error.message);
  }).finally(() => {
    loadingSubmit.value = false
  })
}

let eventSource

const autoExchange = () => {
  if (!addressId.value) {
    message.error("请选择地址");
    return
  }
  const data = []
  for (const item of goods.value) {
    if (item.num > 0) {
      const row = {sku_id: item.sku_id, num: item.num}
      data.push(row)
    }
  }
  if (data.length === 0) {
    message.error("请选择物品");
    return;
  }

  let baseUrl = import.meta.env.MODE == "release" ? window.location.origin : import.meta.env.VITE_SERVER
  baseUrl = baseUrl + "/v1"
  const order = {address_id: addressId.value, product: data}
  eventSource = new EventSource(baseUrl + "/exchange-auto?token=" + token.value + "&time=" + timeStr.value + "&data=" + JSON.stringify(order) + "&delayTime=" + delayTime.value);

  eventSource.addEventListener("message", function (event) {
    const data = event.data;
    output.value.push({time: formatTime(new Date()), message: data})
    if (output.value.length > 99) {
      output.value.shift();
    }
    scrollToBottom()
  });

  eventSource.addEventListener("error", function () {
    message.success("SSE连接错误");
  });

  eventSource.addEventListener("stop", function () {
    message.success("任务已结束");
    eventSource.close()
    start.value = false
  });
  message.success("任务已开始");
  start.value = true
}

function scrollToBottom() {
  if (outputRef.value) {
    outputRef.value.scrollTop = outputRef.value.scrollHeight;
  }
}

const stop = () => {
  eventSource.close()
  message.success("任务已暂停");
  start.value = false
}

const setToken = (e: MouseEvent) => {
  localStorage.setItem('token', token.value);
  loadGoods()
  loadUser()
  open.value = false;
};

const handleChange = (value: string) => {
  let addr = address.value.find((item: { id: string; }) => item.id === value)
  addressId.value = addr.id
  console.log(addressId.value)
  localStorage.setItem('address', JSON.stringify(addr));
};

const focus = () => {
  loadAddress()
};

const numChange = (value: number) => {
  const data = []
  for (const item of goods.value) {
    if (item.num > 0) {
      const row = {name: item.name, sku_id: item.sku_id, num: item.num}
      data.push(row)
    }
  }
  goodsSub.value = data
  localStorage.setItem('goods', JSON.stringify(data));
}

const onRangeChange = (value: [Dayjs, Dayjs], dateString: [string, string]) => {
  timeStr.value = dateString;
};

const showDrawer = () => {
  drawerStatus.value = true;
};

const oncloseDrawer = () => {
  drawerStatus.value = false;
};

const showModal = () => {
  open.value = true;
};

const formatTime = (date) => {
  const hours = date.getHours().toString().padStart(2, '0');
  const minutes = date.getMinutes().toString().padStart(2, '0');
  const seconds = date.getSeconds().toString().padStart(2, '0');
  return hours + ':' + minutes + ':' + seconds;
}

</script>
<style scoped>
.header {
  padding: 8px 0;
  box-shadow: 0 1px 2px 0 rgba(0, 0, 0, 0.03), 0 1px 6px -1px rgba(0, 0, 0, 0.02), 0 2px 4px 0 rgba(0, 0, 0, 0.02);
}

.content {
  width: 80%;
  margin: 8px auto;
}

.footer {
  text-align: center;
  position: fixed;
  bottom: 0;
  right: 0;
  width: 100%;
  opacity: .55;
  font-size: 12px;
  color: #000;
  height: 30px;
  line-height: 30px;
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
  