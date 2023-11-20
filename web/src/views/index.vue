<template>
    <a-layout-header class="header clearfix">
        <div class="left fl">
            <a-button @click="showModal" style="margin-left: 8px">设置Token</a-button>
            <a-modal v-model:open="open" title="设置Token" @ok="setToken">
                <a-textarea v-model:value="token" placeholder="请输入" style="height: 200px;" />
            </a-modal>
            <a-select ref="select" v-model:value="defaultAddress" style="width: 220px;margin-left: 20px;" @focus="focus"
                @change="handleChange" placeholder="请选择收货地址">
                <a-select-option v-for="item in address" :value="item.id">{{ item.address }}</a-select-option>
            </a-select>
        </div>
        <div class="right fr" v-if="userInfo">
            <a-avatar :src="userInfo.avatar_url" style="margin-right: 8px;"/>
            <a-tag>昵称：{{ userInfo.nick_name }}</a-tag>
            <a-tag style="margin-right: 8px">当前积分：{{ userInfo.currentScoreSum }}</a-tag>
        </div>
    </a-layout-header>
    <a-layout-content class="content">
        <a-table class="table" :columns="columns" :data-source="goods" bordered :pagination="false" :loading="loadingTable">
            <template #bodyCell="{ column, text, record }">
                <template v-if="column.dataIndex === 'name'">
                    <a>{{ text }}</a>
                </template>
                <template v-if="column.dataIndex === 'point'">
                    <a-tag>{{ text }}</a-tag>
                </template>
                <template v-if="column.dataIndex === 'mainImage'">
                    <a-image :width="58" :src="text" />
                </template>

                <template v-else-if="column.dataIndex === 'operation'">
                    <a-input-number id="inputNumber" v-model:value="record.num" :min="0" :max="10" @change="numChange" />
                </template>
            </template>
        </a-table>

        <div class="submit">
            <a-button class="right-button" @click="exchangeGoods" :loading="loadingSubmit" danger>立即兑换</a-button>
        </div>
    </a-layout-content>

    <!-- <a-layout-footer class="footer">
    </a-layout-footer> -->
</template>
  
<script lang="ts" setup>
import { onMounted, ref } from 'vue';
import { getGoods, getAddress, exchange, getUser } from '@/api/index';
import type { AxiosPromise } from "axios";
import { message } from 'ant-design-vue';
import { UserOutlined } from '@ant-design/icons-vue';

const token = ref<string>('')
const open = ref<boolean>(false);
const goods = ref<any>();
const address = ref<any>();
const defaultAddress = ref<string>();
const addressId = ref<string>()
const userInfo = ref<any>()
const loadingSubmit = ref<boolean>(false)
const loadingTable = ref<boolean>(false)

const showModal = () => {
    open.value = true;
};

onMounted(() => {
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
        message.error(error.message);
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
            const row = { sku_id: item.sku_id, num: item.num }
            data.push(row)
        }
    }
    const response: AxiosPromise<any[]> = exchange({ address_id: addressId.value, product: data }, token.value);
    response.then((res: any) => {
        console.log(res)
        message.success("兑换成功！");
    }).catch((error: any) => {
        message.error(error.message);
    }).finally(() => {
        loadingSubmit.value = false
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
  