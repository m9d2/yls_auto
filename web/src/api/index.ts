import request from "@/plugins/axios";
import type { AxiosPromise } from "axios";

enum API {
    GOODS_URL = '/goods',
    ADDRESS_URL = '/address',
    EXCHANGE_URL = '/exchange',
    USER_URL = '/user',
}

export const getGoods = (token: string): AxiosPromise<any[]> => {
    return request.get<any[]>(API.GOODS_URL + "?token=" + token);
};

export const getAddress = (token: string): AxiosPromise<any[]> => {
    return request.get<any[]>(API.ADDRESS_URL + "?token=" + token);
};

export const exchange = (data: any, token: string): AxiosPromise<any[]> => {
    return request.post<any[]>(API.EXCHANGE_URL + "?token=" + token, data);
};

export const getUser = (token: string): AxiosPromise<any[]> => {
    return request.get<any[]>(API.USER_URL + "?token=" + token);
};

export const getTags = (): AxiosPromise<any[]> => {
    return request.get<any[]>("https://api.github.com/repos/m9d2/ysl_auto/tags");
};