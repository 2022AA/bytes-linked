import { defineStore } from 'pinia';

const routerStore = defineStore({
    // 这里的id必须为唯一ID
    id: 'router',
    state: () => {
        return {
            routerFrom: '/',
            routerName: '/home',
        };
    },
    // 等同于vuex的getter
    getters: {
        getRouterFrom: (state) => state.routerFrom,
        getRouterName: (state) => state.routerName,
    },
    // pinia 放弃了 mutations 只使用 actions
    actions: {
        // actions可以用async做成异步形式
        setName(type: string) {
            this.routerName = type;
        },
    },
});

export default routerStore;
